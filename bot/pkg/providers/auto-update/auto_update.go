package auto_update

import (
	"sync"

	"github.com/pingcap/community/bot/config"
	"github.com/pingcap/community/bot/pkg/operator"

	"github.com/google/go-github/v29/github"
)

// AutoUpdate defines methods of auto
type AutoUpdate interface {
	ProcessPullRequestEvent(event *github.PullRequestEvent)
}

type autoUpdate struct {
	owner           string
	watchedRepo     string
	updateOwner     string
	updateRepo      string
	targetMap       map[string]string
	updateScript    string
	updateAutoMerge bool
	opr             *operator.Operator
	cfg             *config.RepoConfig
	sync.Mutex
}

// Init create cherry pick middleware instance
func Init(repo *config.RepoConfig, opr *operator.Operator) AutoUpdate {
	c := autoUpdate{
		owner:           repo.Owner,
		watchedRepo:     repo.Repo,
		updateOwner:     repo.UpdateOwner,
		updateRepo:      repo.UpdateRepo,
		targetMap:       repo.UpdateTargetMap,
		updateScript:    repo.UpdateScript,
		updateAutoMerge: repo.UpdateAutoMerge,
		opr:             opr,
		cfg:             repo,
	}
	return &c
}
