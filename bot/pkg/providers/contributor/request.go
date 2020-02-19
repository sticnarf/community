package contributor

import (
	"context"

	"github.com/google/go-github/v29/github"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type ApproveRecord struct {
	ID        int    `gorm:"column:id"`
	Owner     string `gorm:"column:owner"`
	Repo      string `gorm:"column:repo"`
	Github    string `gorm:"column:github"`
	CreatedAt string `gorm:"column:created_at"`
}

func (c *Contributor) addContributor(pull *github.PullRequest) error {
	return c.labelPull(pull, c.cfg.ContributorLabel)
}

func (c *Contributor) labelPull(pull *github.PullRequest, label string) error {
	if label == "" {
		return nil
	}
	var labels []string

	hasTargetLabelLabel := false
	for _, l := range pull.Labels {
		labels = append(labels, *l.Name)
		if *l.Name == label {
			hasTargetLabelLabel = true
		}
	}
	if !hasTargetLabelLabel {
		labels = append(labels, label)
	}

	_, _, err := c.opr.Github.Issues.AddLabelsToIssue(context.Background(),
		c.owner, c.repo, *pull.Number, labels)
	return errors.Wrap(err, "label PR")
}

func (c *Contributor) isReviewer(login string) (bool, error) {
	model := &ApproveRecord{}
	if err := c.opr.DB.Where("github = ?", login).First(model).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, errors.Wrap(err, "query can approve failed")
	}
	return true, nil
}
