package command

import (
	"github.com/pingcap/community/bot/config"
	"github.com/pingcap/community/bot/pkg/operator"
)

type CommandRedeliver struct {
	repo *config.RepoConfig
	opr  *operator.Operator
}

// Init create cherry pick middleware instance
func Init(repo *config.RepoConfig, opr *operator.Operator) *CommandRedeliver {
	return &CommandRedeliver{
		repo: repo,
		opr:  opr,
	}
}
