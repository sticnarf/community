package operator

import (
	"github.com/pingcap/community/bot/config"
	"github.com/pingcap/community/bot/pkg/db"
	githubPkg "github.com/pingcap/community/bot/pkg/github"
	"github.com/pingcap/community/bot/pkg/member"
	"github.com/pingcap/community/bot/pkg/slack"
	"github.com/pingcap/community/bot/util"

	"github.com/google/go-github/v29/github"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Operator contains pkg instances
type Operator struct {
	Config *config.Config
	DB     *gorm.DB
	Github *github.Client
	Slack  slack.Bot
	Member *member.Member
}

// InitOperator create context from config
func InitOperator(cfg *config.Config) *Operator {
	githubClient := githubPkg.GetGithubClient(cfg.Github)
	slackClient, err := slack.GetSlackClient(cfg.Slack, cfg.Repos)
	if err != nil {
		util.Fatal(errors.Wrap(err, "init context"))
	}
	dbConnect := db.CreateDbConnect(cfg.Database)
	m := member.New(githubClient)

	return &Operator{
		Config: cfg,
		DB:     dbConnect,
		Github: githubClient,
		Slack:  slackClient,
		Member: m,
	}
}
