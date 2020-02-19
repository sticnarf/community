package bot

import (
	"context"

	"github.com/pingcap/community/bot/pkg/providers/cherry"
	"github.com/pingcap/community/bot/util"

	"github.com/google/go-github/v29/github"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func (b *bot) getPullRequest(prNumber int) (*cherry.PullRequest, error) {
	model := &cherry.PullRequest{}
	if err := b.opr.DB.Where("owner = ? AND repo = ? AND pr_id = ?",
		b.owner, b.repo, prNumber).First(model).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, errors.Wrap(err, "query pull request failed")
	}
	return model, nil
}

func (b *bot) createPullRequest(pr *github.PullRequest) error {
	model, err := b.getPullRequest(*pr.Number)
	if err != nil {
		return errors.Wrap(err, "create pull request")
	}
	// pull request already exist
	if model.PrID != 0 {
		return nil
	}

	merge := false
	if pr.MergedAt != nil {
		merge = true
	}
	// save new pull request
	if model.PrID == 0 {
		prRecord := cherry.PullRequest{
			PrID:      *pr.Number,
			Owner:     b.owner,
			Repo:      b.repo,
			Title:     *pr.Title,
			Label:     "[]",
			Merge:     merge,
			CreatedAt: *pr.CreatedAt,
		}
		err := b.saveModel(&prRecord)
		if err != nil {
			return errors.Wrap(err, "create pull request")
		}
	}
	return nil
}

func (b *bot) saveModel(model interface{}) error {
	ctx := context.Background()
	if err := util.RetryOnError(ctx, maxRetryTime, func() error {
		return b.opr.DB.Save(model).Error
	}); err != nil {
		return errors.Wrap(err, "save pull request into database failed")
	}
	return nil
}
