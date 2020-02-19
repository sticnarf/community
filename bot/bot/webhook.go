package bot

import (
	"github.com/pingcap/community/bot/util"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

func (b *bot) Webhook(event interface{}) {
	switch event := event.(type) {
	case *github.PullRequestEvent:
		b.processPullRequestEvent(event)
	case *github.IssueCommentEvent:
		b.processIssueCommentEvent(event)
	case *github.IssuesEvent:
		b.processIssuesEvent(event)
	case *github.PullRequestReviewEvent:
		b.processPullRequestReviewEvent(event)
	case *github.PullRequestReviewCommentEvent:
		b.processPullRequestReviewCommentEvent(event)
	}
}

func (b *bot) processPullRequestEvent(event *github.PullRequestEvent) {
	if *event.Action == "opened" || *event.Action == "labeled" {
		if err := b.createPullRequest(event.PullRequest); err != nil {
			util.Error(errors.Wrap(err, "bot process pull request event"))
		}
	}

	if b.cfg.CherryPick {
		b.Middleware.cherry.ProcessPullRequestEvent(event)
	}

	if b.cfg.LabelCheck {
		b.Middleware.label.ProcessPullRequestEvent(event)
	}

	if b.cfg.PrLimit {
		b.Middleware.Prlimit.ProcessPullRequestEvent(event)
	}

	if b.cfg.Merge {
		b.Middleware.Merge.ProcessPullRequestEvent(event)
	}

	if b.cfg.StatusControl {
		b.Middleware.PullStatus.ProcessPullRequestEvent(event)
	}

	if b.cfg.AutoUpdate {
		b.Middleware.AutoUpdate.ProcessPullRequestEvent(event)
	}

	b.Middleware.Contributor.ProcessPullRequestEvent(event)
}

func (b *bot) processIssuesEvent(event *github.IssuesEvent) {
	util.Println("process issue event in bot", event.GetIssue().GetNumber())
	if b.cfg.IssueRedeliver {
		b.Middleware.IssueRedeliver.ProcessIssuesEvent(event)
	}
	if b.cfg.LabelCheck {
		b.Middleware.label.ProcessIssuesEvent(event)
	}
	if b.cfg.IssueSlackNotice {
		b.Middleware.Notify.ProcessIssuesEvent(event)
	}
}

func (b *bot) processIssueCommentEvent(event *github.IssueCommentEvent) {
	if b.cfg.CherryPick {
		b.Middleware.cherry.ProcessIssueCommentEvent(event)
	}

	if b.cfg.Merge {
		b.Middleware.Merge.ProcessIssueCommentEvent(event)
	}

	if b.cfg.IssueRedeliver {
		b.Middleware.IssueRedeliver.ProcessIssueCommentEvent(event)
	}

	if b.cfg.StatusControl {
		b.Middleware.PullStatus.ProcessIssueCommentEvent(event)
	}

	if b.cfg.IssueSlackNotice {
		b.Middleware.Notify.ProcessIssueCommentEvent(event)
	}

	if b.cfg.PullApprove {
		b.Middleware.Approve.ProcessIssueCommentEvent(event)
	}

	b.Middleware.CommandRedeliver.ProcessIssueCommentEvent(event)
}

func (b *bot) processPullRequestReviewEvent(event *github.PullRequestReviewEvent) {
	if b.cfg.StatusControl {
		b.Middleware.PullStatus.ProcessPullRequestReviewEvent(event)
	}
}

func (b *bot) processPullRequestReviewCommentEvent(event *github.PullRequestReviewCommentEvent) {
	if b.cfg.StatusControl {
		b.Middleware.PullStatus.ProcessPullRequestReviewCommentEvent(event)
	}
}
