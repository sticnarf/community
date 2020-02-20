package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Config is cherry picker config struct
type Config struct {
	Github   *Github
	Slack    *Slack
	Repos    map[string]*RepoConfig
	Database *Database
}

// Redeliver is struct of redeliver rule
type Redeliver struct {
	Label   string
	Keyword string
	Exclude string
	Channel string
}

// PullStatusControlEvent is a rule for pull status control
type PullStatusControlEvent struct {
	Duration int
	Events   string
}

// PullStatusControl is pull status control for a specific label
type PullStatusControl struct {
	Label  string
	Cond   string
	Events []*PullStatusControlEvent `toml:"event"`
}

// RepoConfig is single repo config
type RepoConfig struct {
	// common config
	GithubBotChannel string
	// repo config
	Owner          string
	Repo           string
	Interval       int
	Fullupdate     int
	WebhookSecret  string
	Rule           string
	Release        string
	TypeLabel      string
	IgnoreLabel    string
	RunTestCommand string
	// cherry picker config
	CherryPick          bool
	Dryrun              bool
	CherryPickChannel   string
	ShortCheckDuration  int
	MediumCheckDuration int
	LongCheckDuration   int
	CommonChecker       string
	ChiefChecker        string
	ReplaceLabel        string
	// label check config
	LabelCheck        bool
	LabelCheckChannel string
	DefaultChecker    string
	// pr limit config
	PrLimit          bool
	MaxPrOpened      int
	PrLimitMode      string
	PrLimitOrgs      string
	PrLimitLabel     string
	ContributorLabel string
	// merge config
	Merge                bool
	CanMergeLabel        string
	ReleaseAccessControl bool
	// issue redeliver
	IssueRedeliver   bool
	Redeliver        []*Redeliver
	SignedOffMessage bool
	// pull request status control
	StatusControl     bool `toml:"statusControl"`
	SstatusControl    bool
	LabelOutdated     string
	NoticeChannel     string
	PullStatusControl []*PullStatusControl
	// auto update config
	AutoUpdate        bool
	AutoUpdateChannel string
	UpdateOwner       string
	UpdateRepo        string
	UpdateTargetMap   map[string]string
	UpdateScript      string
	MergeLabel        string
	UpdateAutoMerge   bool
	// issue notify
	IssueSlackNotice        bool
	IssueSlackNoticeChannel string
	IssueSlackNoticeNotify  string
	// approve
	PullApprove bool
	// contributor
	NotifyNewContributorPR bool
}

// Database is db connect config
type Database struct {
	Address  string
	Port     int
	Username string
	Password string
	Dbname   string
}

// Github config
type Github struct {
	Token string
	Bot   string
}

// Slack config
type Slack struct {
	Token     string
	Heartbeat string
	Mute      bool
	Hello     bool
}

type rawConfig struct {
	Github   *Github
	Slack    *Slack
	Repos    []*RepoConfig
	Database *Database
}

// GetConfig read config file
func GetConfig(configPath *string) (*Config, error) {
	rawCfg, err := readConfigFile(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "get config")
	}
	repos := make(map[string]*RepoConfig)
	for _, repo := range rawCfg.Repos {
		key := fmt.Sprintf("%s-%s", repo.Owner, repo.Repo)
		repos[key] = repo
	}
	return &Config{
		Github:   rawCfg.Github,
		Slack:    rawCfg.Slack,
		Repos:    repos,
		Database: rawCfg.Database,
	}, nil
}

func readConfigFile(configPath *string) (*rawConfig, error) {
	var rawCfg rawConfig
	file, err := ioutil.ReadFile(*configPath)
	if err != nil {
		// err
		return nil, errors.Wrap(err, "read config file")
	}
	// json.Unmarshal(file, &rawCfg)
	if _, err := toml.Decode(string(file), &rawCfg); err != nil {
		return nil, errors.Wrap(err, "read config file")
	}
	return &rawCfg, nil
}
