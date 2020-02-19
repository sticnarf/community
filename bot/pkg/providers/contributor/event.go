package contributor

import (
	"github.com/pingcap/community/bot/util"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

func (c *Contributor) ProcessPullRequestEvent(event *github.PullRequestEvent) {
	var err error

	if *event.Action == "opened" || *event.Action == "reopened" {
		err = c.processOpenedPR(event.PullRequest)
	}

	if err != nil {
		util.Error(errors.Wrap(err, "cherry picker process pull request event"))
	}
}

func (c *Contributor) processOpenedPR(pull *github.PullRequest) error {
	if c.cfg.ContributorLabel == "" {
		return nil
	}
	login := pull.GetUser().GetLogin()
	if c.opr.Member.IfMember(login) {
		if isReviewer, err := c.isReviewer(login); err != nil {
			return errors.Wrap(err, "process opened PR")
		} else if !isReviewer {
			// is a member and is not a reviewer -> employee
			return nil
		}
	}

	return c.addContributor(pull)
}
