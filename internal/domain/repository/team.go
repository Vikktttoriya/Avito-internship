package repository

import (
	"context"
	"test/internal/domain/model"
)

type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	GetByName(ctx context.Context, name string) (*model.Team, error)
}
