package repository

import (
	"context"
	"test/internal/domain/model"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByTeam(ctx context.Context, team string) ([]*model.User, error)
	GetActiveByTeam(ctx context.Context, team string) ([]*model.User, error)
	Save(ctx context.Context, u *model.User) error
}
