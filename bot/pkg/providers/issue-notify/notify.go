package notify

import (
	"github.com/pingcap/community/bot/config"
	"github.com/pingcap/community/bot/pkg/operator"
)

type Notify struct {
	owner    string
	repo     string
	ready    bool
	notify   bool
	channel  string
	notifyID string
	opr      *operator.Operator
	cfg      *config.RepoConfig
}

// Init create cherry pick middleware instance
func Init(repo *config.RepoConfig, opr *operator.Operator) *Notify {
	n := Notify{
		owner:    repo.Owner,
		repo:     repo.Repo,
		ready:    false,
		notify:   repo.IssueSlackNotice,
		channel:  repo.IssueSlackNoticeChannel,
		notifyID: repo.IssueSlackNoticeNotify,
		opr:      opr,
		cfg:      repo,
	}
	return &n
}

func (c *Notify) Ready() {
	c.ready = true
}
