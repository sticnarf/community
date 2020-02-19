package merge

import (
	"context"
	"fmt"
	"github.com/pingcap/community/bot/util"
	"time"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

const (
	pollingInterval = 30 * time.Second
	waitForStatus   = 120 * time.Second
	testCommentBody = "/run-all-tests"
	mergeMessage    = "Ready to merge!"
	mergeMethod     = "squash"
)

func (m *merge) processPREvent(event *github.PullRequestEvent) {
	pr := event.GetPullRequest()
	model := AutoMerge{
		PrID:      *pr.Number,
		Owner:     m.owner,
		Repo:      m.repo,
		Status:    false,
		CreatedAt: time.Now(),
	}
	if err := m.saveModel(&model); err != nil {
		util.Error(errors.Wrap(err, "merge process PR event"))
	} else {
		util.Error(m.queueComment(pr))
	}
}

func (m *merge) startJob(mergeJob *AutoMerge) error {
	pr, _, err := m.opr.Github.PullRequests.Get(context.Background(), m.owner, m.repo, (*mergeJob).PrID)
	if err != nil {
		return errors.Wrap(err, "start merge job")
	}
	if pr.MergedAt != nil {
		return nil
	}
	needUpdate, err := m.updateBranch(pr)
	util.Println("update branch", needUpdate, err)
	if err != nil {
		if _, ok := err.(*github.AcceptedError); ok {
			// no need for update branch, continue test
		} else {
			return errors.Wrap(err, "start merge job")
		}
	}
	if needUpdate {
		time.Sleep(waitForStatus)
	}

	commentBody := testCommentBody
	conmment := &github.IssueComment{
		Body: &commentBody,
	}
	_, _, err = m.opr.Github.Issues.CreateComment(context.Background(),
		m.owner, m.repo, *pr.Number, conmment)
	if err != nil {
		return errors.Wrap(err, "start merge job")
	}
	mergeJob.Started = true
	if err := m.saveModel(mergeJob); err != nil {
		util.Error(errors.Wrap(err, "start merge job"))
	}
	time.Sleep(waitForStatus)
	return nil
}

func (m *merge) startPolling() {
	ticker := time.NewTicker(pollingInterval)
	go func() {
		for range ticker.C {
			jobs := m.getMergeJobs()
			if len(jobs) == 0 {
				continue
			}

			var job *AutoMerge
			for _, model := range jobs {
				if model.Started {
					job = model
				}
			}
			if job == nil {
				job = jobs[0]
				m.startJob(job)
			}

			ifComplete := m.checkPR(job)
			if ifComplete {
				job.Status = true
				if err := m.saveModel(job); err != nil {
					util.Error(errors.Wrap(err, "merge polling job"))
				}
			}
		}
	}()
}

func (m *merge) checkPR(mergeJob *AutoMerge) bool {
	pr, _, err := m.opr.Github.PullRequests.Get(context.Background(), m.owner, m.repo, (*mergeJob).PrID)
	if err != nil {
		util.Error(errors.Wrap(err, "checking PR if can be merged"))
		return false
	}
	if pr.MergedAt != nil {
		return true
	}
	// check if still have "can merge" label
	ifHasLabel := false
	for _, l := range pr.Labels {
		if *l.Name == m.cfg.CanMergeLabel {
			ifHasLabel = true
		}
	}
	if !ifHasLabel {
		return true
	}

	// if need update, update branch & re-run test
	needUpdate, err := m.needUpdateBranch(pr)
	if err == nil && needUpdate {
		util.Println("restart job due to branch need update")
		m.startJob(mergeJob)
	}

	success := true
	finish := true
	status, _, err := m.opr.Github.Repositories.GetCombinedStatus(context.Background(), m.owner, m.repo,
		*pr.Head.SHA, nil)
	if err != nil {
		util.Error(errors.Wrap(err, "polling PR status"))
		return false
	}
	if *status.State == "failure" || *status.State == "error" {
		success = false
		util.Println("Tests failed in statuses", status)
	}
	if *status.State == "pending" {
		finish = false
	}

	checks, _, err := m.opr.Github.Checks.ListCheckRunsForRef(context.Background(), m.owner, m.repo,
		*pr.Head.SHA, nil)
	if err != nil {
		util.Error(errors.Wrap(err, "polling PR status"))
		return false
	}
	for _, check := range checks.CheckRuns {
		if *check.Status != "completed" {
			finish = false
		} else {
			if *check.Conclusion != "success" {
				success = false
				util.Println("Tests failed in check-runs", checks)
			}
		}
	}

	if success && finish {
		// send success comment and merge it
		// if err := m.addGithubComment(pr, mergeMessage); err != nil {
		// 	util.Error(errors.Wrap(err, "checking PR"))
		// }
		util.Println("Tests test passed", status, checks)
		// compose commit message
		message := ""
		if m.cfg.SignedOffMessage {
			msg, err := m.getMergeMessage(pr.GetNumber())
			if err != nil {
				util.Error(errors.Wrap(err, "merging PR"))
			} else {
				message = msg
			}
		}

		opt := github.PullRequestOptions{
			CommitTitle: fmt.Sprintf("%s (#%d)", pr.GetTitle(), pr.GetNumber()),
			MergeMethod: mergeMethod,
		}
		_, _, err := m.opr.Github.PullRequests.Merge(context.Background(), m.owner, m.repo,
			pr.GetNumber(), message, &opt)
		if err != nil {
			util.Error(errors.Wrap(err, "checking PR"))
			if err := m.failedMergeSlack(pr); err != nil {
				util.Error(errors.Wrap(err, "checking PR"))
			}
		} else {
			if err := m.successMergeSlack(pr); err != nil {
				util.Error(errors.Wrap(err, "checking PR"))
			}
		}
	} else if !success {
		finish = true
		// send failure comment
		comment := fmt.Sprintf("@%s merge failed.", *pr.User.Login)
		if err := m.addGithubComment(pr, comment); err != nil {
			util.Error(errors.Wrap(err, "checking PR"))
		}
		if err := m.failedMergeSlack(pr); err != nil {
			util.Error(errors.Wrap(err, "checking PR"))
		}
	}

	return finish
}
