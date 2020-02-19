package cherry

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pingcap/community/bot/util"

	"github.com/google/go-github/v29/github"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	maxRetryTime = 1
	workDir      = "/tmp"
)

// PullRequest is pull request table structure
type PullRequest struct {
	ID        int       `gorm:"id"`
	PrID      int       `gorm:"pr_id"`
	Owner     string    `gorm:"owner"`
	Repo      string    `gorm:"repo"`
	Title     string    `gorm:"title"`
	Label     string    `gorm:"label"`
	Merge     bool      `gorm:"merge"`
	CreatedAt time.Time `gorm:"created_at"`
}

// CherryPr is cherry pick table structure
type CherryPr struct {
	ID           int    `gorm:"id"`
	PrID         int    `gorm:"pr_id"`
	FromPr       int    `gorm:"from_pr"`
	Owner        string `gorm:"owner"`
	Repo         string `gorm:"repo"`
	Title        string `gorm:"title"`
	Head         string `gorm:"head"`
	Base         string `gorm:"base"`
	Body         string `gorm:"body"`
	CreatedByBot bool   `gorm:"created_by_bot"`
	TryTime      int    `gorm:"try_time"`
	Success      bool   `gorm:"success"`
}

// SlackUser is slack user table structure
type SlackUser struct {
	ID     int    `gorm:"id"`
	Github string `gorm:"github"`
	Email  string `gorm:"email"`
	Slack  string `gorm:"slack"`
}

type labelSlice []string

func parseLabel(labels string) (*labelSlice, error) {
	s := labelSlice{}
	err := json.Unmarshal([]byte(labels), &s)
	if err != nil {
		return nil, errors.Wrap(err, "parse label")
	}
	return &s, nil
}

func (cherry *cherry) createPullRequest(pr *github.PullRequest) error {
	model, err := cherry.getPullRequest(*pr.Number)
	if err != nil {
		return errors.Wrap(err, "create pull request")
	}

	merge := false
	if pr.MergedAt != nil {
		merge = true
	}
	// save new pull request
	if model.PrID == 0 {
		prRecord := PullRequest{
			PrID:      *pr.Number,
			Owner:     cherry.owner,
			Repo:      cherry.repo,
			Title:     *pr.Title,
			Label:     "[]",
			Merge:     merge,
			CreatedAt: *pr.CreatedAt,
		}
		err := cherry.saveModel(&prRecord)
		if err != nil {
			return errors.Wrap(err, "create pull request")
		}
	}
	return nil
}

func (cherry *cherry) createCherryPick(pr *github.PullRequest) error {
	r, _ := regexp.Compile(`\(#([0-9]+)\)$`)
	m := r.FindStringSubmatch(*pr.Title)

	if len(m) < 2 {
		return errors.New("not cherry pick pr")
	}

	from, err := strconv.Atoi(m[1])
	if err != nil {
		return errors.Wrap(err, "get cherry pick")
	}

	model, err := cherry.getCherryPick(from, *pr.Base.Ref)
	if err != nil {
		return errors.Wrap(err, "get cherry pick")
	}

	if model.PrID != 0 {
		// cherry pick already exist
		return nil
	}

	model.PrID = *pr.Number
	model.FromPr = from
	model.Owner = cherry.owner
	model.Repo = cherry.repo
	model.Title = *pr.Title
	model.Head = *pr.Head.Label
	model.Base = *pr.Base.Ref
	model.Body = "" // drop user-edit body, useless
	model.CreatedByBot = false
	model.Success = true
	model.TryTime = 0

	err = cherry.saveModel(&model)
	if err != nil {
		return errors.Wrap(err, "get cherry pick")
	}
	return nil
}

func (cherry *cherry) getTarget(label string) (string, string, error) {
	r, err := regexp.Compile(cherry.rule)
	if err != nil {
		return "", "", errors.Wrap(err, "get target")
	}
	m := r.FindStringSubmatch(label)
	if len(m) < 2 {
		return "", "", errors.New("label parse failed")
	}
	targetVersion := m[1]
	if targetVersion == "master" {
		return targetVersion, targetVersion, nil
	}
	return strings.Replace(cherry.release, "[version]", targetVersion, 1), targetVersion, nil
}

func hasLabel(label string, labels string) (bool, error) {
	existLabel, err := parseLabel(labels)
	if err != nil {
		return false, errors.Wrap(err, "has label ditect")
	}
	hasLabel := false
	for _, l := range *existLabel {
		if l == label {
			hasLabel = true
		}
	}
	return hasLabel, nil
}

func (cherry *cherry) addLabel(model *PullRequest, label string) error {
	existLabel, err := parseLabel(model.Label)
	if err != nil {
		return errors.Wrap(err, "add PR label")
	}
	labels := append(*existLabel, label)
	labelString, err := json.Marshal(labels)
	if err != nil {
		return errors.Wrap(err, "add PR label")
	}
	model.Label = string(labelString)
	if err := cherry.opr.DB.Save(&model).Error; err != nil {
		return errors.Wrap(err, "add PR label")
	}
	return nil
}

func (cherry *cherry) removeLabel(pr *github.PullRequest, label string) error {
	model, err := cherry.getPullRequest(*pr.Number)
	if err != nil {
		return errors.Wrap(err, "remove PR label")
	}
	labels, err := parseLabel(model.Label)
	if err != nil {
		return errors.Wrap(err, "remove PR label")
	}
	var newLabels labelSlice
	for _, l := range *labels {
		if l != label {
			newLabels = append(newLabels, l)
		}
	}
	labelString, err := json.Marshal(labels)
	if err != nil {
		return errors.Wrap(err, "remove PR label")
	}
	model.Label = string(labelString)
	if err := cherry.opr.DB.Save(&model).Error; err != nil {
		return errors.Wrap(err, "add PR label")
	}
	return nil
}

func (cherry *cherry) getPullRequest(prNumber int) (*PullRequest, error) {
	model := &PullRequest{}
	if err := cherry.opr.DB.Where("owner = ? AND repo = ? AND pr_id = ?",
		cherry.owner, cherry.repo, prNumber).First(model).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, errors.Wrap(err, "query pull request failed")
	}
	return model, nil
}

func (cherry *cherry) getCherryPick(from int, base string) (*CherryPr, error) {
	model := &CherryPr{}
	if err := cherry.opr.DB.Where("owner = ? AND repo = ? AND from_pr = ? AND base = ?",
		cherry.owner, cherry.repo, from, base).First(model).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, errors.Wrap(err, "query cherry pick failed")
	}
	return model, nil
}

func (cherry *cherry) saveModel(model interface{}) error {
	ctx := context.Background()
	if err := util.RetryOnError(ctx, maxRetryTime, func() error {
		return cherry.opr.DB.Save(model).Error
	}); err != nil {
		switch model.(type) {
		case *PullRequest:
			return errors.Wrap(err, "save pull request into database failed")
		case *CherryPr:
			return errors.Wrap(err, "save cherry pick into database failed")
		}
	}
	return nil
}

func (cherry *cherry) prepareCherryPick(pr *github.PullRequest, target string) (*github.NewPullRequest, string, error) {
	var newPr github.NewPullRequest
	var message string

	ctx := context.Background()
	err := util.RetryOnError(ctx, maxRetryTime, func() error {
		originRepo := fmt.Sprintf("https://%s:%s@github.com/%s/%s.git", cherry.opr.Config.Github.Bot,
			cherry.opr.Config.Github.Token, cherry.owner, cherry.repo)
		botRepo := fmt.Sprintf("https://%s:%s@github.com/%s/%s.git", cherry.opr.Config.Github.Bot,
			cherry.opr.Config.Github.Token, cherry.opr.Config.Github.Bot, cherry.repo)
		folder := fmt.Sprintf("%s-%s-%s", cherry.owner, cherry.repo, (*pr.MergeCommitSHA)[0:12])
		dir := fmt.Sprintf("%s/%s", workDir, folder)
		newBranch := fmt.Sprintf("%s-%s", target, (*pr.MergeCommitSHA)[0:12])
		patchFile := fmt.Sprintf("%d.patch", *pr.Number)
		patchURI := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", cherry.owner, cherry.repo, *pr.Number)
		commit := fmt.Sprintf("%s (#%d)", *pr.Title, *pr.Number)
		head := fmt.Sprintf("%s:%s", cherry.opr.Config.Github.Bot, newBranch)
		body := fmt.Sprintf("cherry-pick #%d to %s\n\n---\n\n%s", pr.GetNumber(), target, pr.GetBody())
		maintainerCanModify := true
		draft := false

		defer do(workDir, "rm", "-rf", folder)
		if _, err := do(workDir, "git", "clone", originRepo, folder); err != nil {
			message = "git clone failed"
			return errors.Wrap(err, "clone failed")
		}
		if _, err := do(dir, "git", "checkout", target); err != nil {
			message = fmt.Sprintf("branch %s not exist", target)
			return errors.Wrap(err, "checkout failed")
		}
		if _, err := do(dir, "git", "checkout", "-b", newBranch); err != nil {
			message = "create new branch failed"
			return errors.Wrap(err, "checkout to new branch failed")
		}
		if _, err := do(dir, "curl", "-o", patchFile, "-sSL",
			"-H", fmt.Sprintf(`Authorization: token %s`, cherry.opr.Config.Github.Token),
			"-H", `Accept: application/vnd.github.v3.patch`,
			patchURI); err != nil {
			message = "download patch file failed"
			return errors.Wrap(err, "get patch failed")
		}
		if out, err := do(dir, "git", "am", "-3", patchFile); err != nil {
			conflictFileMessage := "Conflict files:"
			sha1Message := "sha1 information is lacking or useless"
			lines := strings.Split(out, "\n")

			var conflictErrs []string
			var sha1Errs []string
			rConflict, _ := regexp.Compile(`Merge conflict in (.*)$`)
			rSha1, _ := regexp.Compile(`sha1 information is lacking or useless \((.*)\).*$`)
			for _, line := range lines {
				mConflict := rConflict.FindStringSubmatch(line)
				if len(mConflict) >= 2 {
					conflictErrs = append(conflictErrs, mConflict[1])
				}
				mSha1 := rSha1.FindStringSubmatch(line)
				if len(mSha1) >= 2 {
					sha1Errs = append(sha1Errs, mSha1[1])
				}
			}

			if len(conflictErrs) > 0 {
				message = fmt.Sprintf("%s%s\n%s\n", message, conflictFileMessage,
					strings.Join(conflictErrs, "\n"))
			}
			if len(sha1Errs) > 0 {
				message = fmt.Sprintf("%s%s\n%s\n", message, sha1Message,
					strings.Join(sha1Errs, "\n"))
			}
			return errors.Wrap(err, "patch failed")
		}
		if _, err := do(dir, "rm", patchFile); err != nil {
			message = "remove patch file failed"
			return errors.Wrap(err, "remove patch file failed")
		}
		if _, err := do(dir, "git", "push", botRepo, newBranch); err != nil {
			message = "git push failed"
			return errors.Wrap(err, "git push failed")
		}

		cherryPickPr := github.NewPullRequest{
			Title:               &commit,
			Head:                &head,
			Base:                &target,
			Body:                &body,
			MaintainerCanModify: &maintainerCanModify,
			Draft:               &draft,
		}
		newPr = cherryPickPr
		return nil
	})
	if err != nil {
		return nil, message, errors.Wrap(err, "prepare pull request")
	}
	return &newPr, "", nil
}

func (cherry *cherry) submitCherryPick(newPr *github.NewPullRequest) (*github.PullRequest, int, error) {
	var (
		resPr   *github.PullRequest
		tryTime int
	)

	if cherry.dryrun {
		number := -1
		return &github.PullRequest{
			Number: &number,
		}, 1, nil
	}

	ctx := context.Background()
	err := util.RetryOnError(ctx, maxRetryTime, func() error {
		p, _, err := cherry.opr.Github.PullRequests.Create(context.Background(),
			cherry.owner, cherry.repo, newPr)
		if err != nil {
			util.Error(err)
			if er, ok := err.(*github.ErrorResponse); ok {
				if er.Message == "Validation Failed" {
					// pull request already exist
					newPr = nil
					tryTime = 1
					return nil
				}
				return errors.Wrap(err, "create github PR failed")
			}
			return nil
		}
		resPr = p
		return nil
	})

	if err != nil {
		return nil, tryTime, errors.Wrap(err, "create github PR failed")
	}
	return resPr, tryTime, nil
}

func (cherry *cherry) addGithubLabel(res *github.PullRequest, from *github.PullRequest, version string) error {
	var labels []string
	r, err := regexp.Compile(cherry.rule)
	if err != nil {
		return errors.Wrap(err, "add github label")
	}
	ignore, err := regexp.Compile(cherry.ignoreLabel)
	if err != nil {
		return errors.Wrap(err, "add github label")
	}

	for _, l := range from.Labels {
		label := l.GetName()
		// skip some labels
		if ignore.MatchString(label) {
			continue
		}
		// copy common labels
		m := r.FindStringSubmatch(label)
		if len(m) < 2 {
			labels = append(labels, label)
			continue
		}
		// convert cherry pick labels
		labelVersion := m[1]
		if version != labelVersion {
			continue
		}
		typeLabel := strings.Replace(cherry.typeLabel, "[version]", labelVersion, 1)
		labels = append(labels, typeLabel)
	}

	_, _, err = cherry.opr.Github.Issues.AddLabelsToIssue(context.Background(),
		cherry.owner, cherry.repo, *res.Number, labels)
	if err != nil {
		return errors.Wrap(err, "add github label")
	}
	return nil
}

func (cherry *cherry) replaceGithubLabel(pr *github.PullRequest, version string) error {
	if cherry.cfg.ReplaceLabel == "" {
		return nil
	}

	if _, _, err := cherry.opr.Github.Issues.AddLabelsToIssue(context.Background(), cherry.owner, cherry.repo,
		pr.GetNumber(), []string{strings.Replace(cherry.cfg.ReplaceLabel, "[version]", version, 1)}); err != nil {
		return errors.Wrap(err, "replace github label")
	}

	r, err := regexp.Compile(cherry.rule)
	if err != nil {
		return errors.Wrap(err, "replace github label")
	}
	for _, l := range pr.Labels {
		label := l.GetName()
		m := r.FindStringSubmatch(label)
		if len(m) < 2 {
			continue
		}
		if version != m[1] {
			continue
		}
		if _, err := cherry.opr.Github.Issues.RemoveLabelForIssue(context.Background(), cherry.owner, cherry.repo,
			pr.GetNumber(), label); err != nil {
			return errors.Wrap(err, "replace github label")
		}
	}

	return nil
}

func (cherry *cherry) prNotice(success bool, target string,
	pr *github.PullRequest, newPr *github.PullRequest, stat string, message string) error {
	if pr == nil || pr.User == nil {
		return errors.Wrap(errors.New("nil pull request"), "send pr notice")
	}

	var channels []string

	slack := cherry.getSlackByGithub(*pr.User.Login)
	if slack == "" {
		for _, e := range strings.Split(cherry.cfg.DefaultChecker, ",") {
			if e != "" {
				channel := cherry.getSlackByGithub(e)
				if channel != "" {
					channels = append(channels, channel)
				}
			}
		}
	} else {
		channels = append(channels, slack)
	}

	if success {
		uri := fmt.Sprintf("https://github.com/%s/%s/pull/%d", cherry.owner, cherry.repo, newPr.GetNumber())
		origin := fmt.Sprintf("https://github.com/%s/%s/pull/%d", cherry.owner, cherry.repo, pr.GetNumber())
		msg := fmt.Sprintf("✅ Create cherry pick pull request from `%s` to `%s`\n%s\nFrom: %s",
			pr.GetHead().GetLabel(), target, uri, origin)

		if err := cherry.opr.Slack.SendMessageWithPr(cherry.cfg.CherryPickChannel,
			msg, newPr, "success"); err != nil {
			return errors.Wrap(err, "send pr notice")
		}
	} else {
		uri := fmt.Sprintf("https://github.com/%s/%s/pull/%d", cherry.owner, cherry.repo, pr.GetNumber())
		msg := fmt.Sprintf("❌ Create cherry pick pull request from `%s` to `%s`\norigin PR\n%s",
			pr.GetHead().GetLabel(), target, uri)
		if message != "" {
			msg = fmt.Sprintf("%s\n%s", msg, message)
		}
		if err := cherry.opr.Slack.SendMessage(cherry.cfg.CherryPickChannel, msg); err != nil {
			return errors.Wrap(err, "send success message")
		}

		for _, channel := range channels {
			err := cherry.opr.Slack.SendMessage(channel, msg)
			if err != nil {
				return errors.Wrap(err, "send pr notice")
			}
		}
	}
	return nil
}

func (cherry *cherry) addGithubReadyComment(pr *github.PullRequest, success bool, target string, to int) error {
	var commentBody string
	if success {
		commentBody = fmt.Sprintf("cherry pick to %s in PR #%d", target, to)
	} else {
		commentBody = fmt.Sprintf("cherry pick to %s failed", target)
	}
	comment := &github.IssueComment{
		Body: &commentBody,
	}
	_, _, err := cherry.opr.Github.Issues.CreateComment(context.Background(),
		cherry.owner, cherry.repo, *pr.Number, comment)
	return errors.Wrap(err, "add github test comment")
}

func (cherry *cherry) addGithubTestComment(pr *github.PullRequest) error {
	commentBody := cherry.cfg.RunTestCommand
	comment := &github.IssueComment{
		Body: &commentBody,
	}
	_, _, err := cherry.opr.Github.Issues.CreateComment(context.Background(),
		cherry.owner, cherry.repo, *pr.Number, comment)
	return errors.Wrap(err, "add github test comment")
}

func (cherry *cherry) addGithubRequestReviews(pr *github.PullRequest, request github.ReviewersRequest) error {
	_, _, err := cherry.opr.Github.PullRequests.RequestReviewers(context.Background(), cherry.owner,
		cherry.repo, pr.GetNumber(), request)
	return errors.Wrap(err, "add github requests reviews")
}

func (cherry *cherry) getReviewers(pr *github.PullRequest) github.ReviewersRequest {
	author := pr.GetUser().GetLogin()
	reviewers := []string{}

	if reviews, _, err := cherry.opr.Github.PullRequests.ListReviews(context.Background(),
		cherry.owner, cherry.repo, pr.GetNumber(), nil); err == nil {
		for _, review := range reviews {
			username := review.GetUser().GetLogin()
			if username != author {
				if !checkExist(username, reviewers) {
					reviewers = append(reviewers, username)
				}
			}
		}
	}

	for _, reviewer := range pr.RequestedReviewers {
		username := reviewer.GetLogin()
		if !checkExist(username, reviewers) {
			reviewers = append(reviewers, username)
		}
	}

	return github.ReviewersRequest{
		Reviewers: reviewers,
	}
}

func (cherry *cherry) getSlackByGithub(github string) string {
	model := SlackUser{}
	if err := cherry.opr.DB.Where("github = ?",
		github).First(&model).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return ""
	}
	if model.Slack != "" {
		return model.Slack
	}
	if model.Email != "" {
		slack, err := cherry.opr.Slack.GetUserByEmail(model.Email)
		if err != nil {
			return ""
		}
		return slack
	}
	return ""
}

func do(dir string, c string, args ...string) (string, error) {
	cmd := exec.Command(c, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func checkExist(item string, slice []string) bool {
	for _, sliceItem := range slice {
		if item == sliceItem {
			return true
		}
	}
	return false
}
