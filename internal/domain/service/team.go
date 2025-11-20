package service

import (
	"context"
	"test/internal/domain/domain_errors"
	"test/internal/domain/model"
	"test/internal/domain/repository"
)

type TeamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
	prRepo   repository.PrRepository
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository, prRepo repository.PrRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, name string, members []*model.User) (*model.Team, error) {
	teamObj, err := s.teamRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if teamObj != nil {
		return nil, domain_errors.ErrTeamExists
	}

	userIDs := make([]string, len(members))
	for i, m := range members {
		userIDs[i] = m.ID
	}

	hasOpenPRs, err := s.prRepo.CheckUserOpenPRs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if hasOpenPRs {
		return nil, domain_errors.ErrUserHasOpenPullRequests
	}

	team := model.NewTeam(name, members)
	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*model.Team, error) {
	team, err := s.teamRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, domain_errors.ErrTeamNotFound
	}
	return team, nil
}
