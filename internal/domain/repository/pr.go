package repository

import (
	"context"
	"test/internal/domain/model"
)

type PrRepository interface {
	GetByID(ctx context.Context, id string) (*model.PullRequest, error)
	Save(ctx context.Context, pr *model.PullRequest) error
	GetByReviewer(ctx context.Context, reviewerID string) ([]*model.PullRequest, error)
	CheckUserOpenPRs(ctx context.Context, userIDs []string) (bool, error)
}
