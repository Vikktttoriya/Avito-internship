package service

import (
	"context"
	"test/internal/domain/domain_errors"
	"test/internal/domain/model"
	"test/internal/domain/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) SetIsActive(ctx context.Context, id string, active bool) (*model.User, error) {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, domain_errors.ErrUserNotFound
	}

	if active {
		u.Activate()
	} else {
		u.Deactivate()
	}

	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}
