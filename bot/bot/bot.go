package bot

import (
	"github.com/pingcap/community/bot/config"
	"github.com/pingcap/community/bot/pkg/operator"
	"github.com/pingcap/community/bot/pkg/providers/approve"
	auto_update "github.com/pingcap/community/bot/pkg/providers/auto-update"
	"github.com/pingcap/community/bot/pkg/providers/cherry"
	"github.com/pingcap/community/bot/pkg/providers/contributor"
	notify "github.com/pingcap/community/bot/pkg/providers/issue-notify"
	"github.com/pingcap/community/bot/pkg/providers/label"
	"github.com/pingcap/community/bot/pkg/providers/merge"
	"github.com/pingcap/community/bot/pkg/providers/prlimit"
	"github.com/pingcap/community/bot/pkg/providers/pullstatus"
	issueRedeliver "github.com/pingcap/community/bot/pkg/providers/redeliver"
	command "github.com/pingcap/community/bot/pkg/providers/redeliver-command"
)

// Bot contains main polling process
type Bot interface {
	StartPolling()
	Webhook(event interface{})
	MonthlyCheck() (*map[string]*[]string, error)
	GetMiddleware() Middleware
}

// Middleware defines middleware struct
type Middleware struct {
	cherry           cherry.Cherry
	label            label.Label
	Prlimit          prlimit.PrLimit
	Merge            merge.Merge
	IssueRedeliver   issueRedeliver.Redeliver
	PullStatus       pullstatus.PullStatus
	AutoUpdate       auto_update.AutoUpdate
	CommandRedeliver *command.CommandRedeliver
	Notify           *notify.Notify
	Approve          *approve.Approve
	Contributor      *contributor.Contributor
}

type bot struct {
	owner              string
	repo               string
	interval           int
	fullupdateInterval int
	rule               string
	release            string
	dryrun             bool
	opr                *operator.Operator
	cfg                *config.RepoConfig
	Middleware         Middleware
}

// InitBot return bot instance
func InitBot(repo *config.RepoConfig, opr *operator.Operator) Bot {
	bot := bot{
		owner:              repo.Owner,
		repo:               repo.Repo,
		interval:           repo.Interval,
		fullupdateInterval: repo.Fullupdate,
		rule:               repo.Rule,
		release:            repo.Release,
		dryrun:             repo.Dryrun,
		opr:                opr,
		cfg:                repo,
		Middleware: Middleware{
			cherry:           cherry.Init(repo, opr),
			label:            label.Init(repo, opr),
			Prlimit:          prlimit.Init(repo, opr),
			Merge:            merge.Init(repo, opr),
			IssueRedeliver:   issueRedeliver.Init(repo, opr),
			PullStatus:       pullstatus.Init(repo, opr),
			AutoUpdate:       auto_update.Init(repo, opr),
			CommandRedeliver: command.Init(repo, opr),
			Notify:           notify.Init(repo, opr),
			Approve:          approve.Init(repo, opr),
			Contributor:      contributor.Init(repo, opr),
		},
	}
	return &bot
}

func (b *bot) ready() {
	b.Middleware.cherry.Ready()
	b.Middleware.label.Ready()
}

func (b *bot) GetMiddleware() Middleware {
	return b.Middleware
}
