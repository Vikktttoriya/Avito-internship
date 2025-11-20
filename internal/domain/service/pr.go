package service

import (
	"context"
	"math/rand"
	"test/internal/domain/domain_errors"
	"test/internal/domain/model"
	"test/internal/domain/repository"
	"time"
)

type PrService struct {
	prRepo   repository.PrRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPrService(pr repository.PrRepository, u repository.UserRepository, t repository.TeamRepository) *PrService {
	return &PrService{
		prRepo:   pr,
		userRepo: u,
		teamRepo: t,
	}
}

func (s *PrService) CreatePR(ctx context.Context, id, name, authorId string) (*model.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pr != nil {
		return nil, domain_errors.ErrPullRequestExists
	}

	author, err := s.userRepo.GetByID(ctx, authorId)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, domain_errors.ErrUserNotFound
	}

	users, err := s.userRepo.GetActiveByTeam(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	pr = model.NewPr(id, name, authorId)

	var candidates []string
	for _, m := range users {
		if m.ID == authorId {
			continue
		}
		candidates = append(candidates, m.ID)
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	if len(candidates) > 2 {
		candidates = candidates[:2]
	}

	pr.AssignedReviewers = candidates

	if err := s.prRepo.Save(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PrService) Merge(ctx context.Context, id string) (*model.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, domain_errors.ErrPullRequestNotFound
	}

	pr.Merge()

	if err := s.prRepo.Save(ctx, pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *PrService) ReassignReviewer(ctx context.Context, id, oldReviewerId string) (*model.PullRequest, string, error) {
	pr, err := s.prRepo.GetByID(ctx, id)
	if err != nil {
		return nil, "", err
	}
	if pr == nil {
		return nil, "", domain_errors.ErrPullRequestNotFound
	}
	if pr.Status == model.StatusMerged {
		return nil, "", domain_errors.ErrPRMerged
	}

	if !contains(pr.AssignedReviewers, oldReviewerId) {
		return nil, "", domain_errors.ErrReviewerNotAssigned
	}

	oldReviewer, err := s.userRepo.GetByID(ctx, oldReviewerId)
	if err != nil {
		return nil, "", err
	}
	if oldReviewer == nil {
		return nil, "", domain_errors.ErrUserNotFound
	}

	users, err := s.userRepo.GetActiveByTeam(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, "", err
	}

	var candidates []string
	for _, m := range users {
		if m.ID == oldReviewerId || contains(pr.AssignedReviewers, m.ID) {
			continue
		}
		candidates = append(candidates, m.ID)
	}

	if len(candidates) == 0 {
		return nil, "", domain_errors.ErrNoReplacementCandidate
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	newReviewerId := candidates[rnd.Intn(len(candidates))]

	pr.ReplaceReviewer(oldReviewerId, newReviewerId)

	err = s.prRepo.Save(ctx, pr)
	if err != nil {
		return nil, "", err
	}

	return pr, newReviewerId, nil
}

func (s *PrService) GetByReviewer(ctx context.Context, id string) ([]*model.PullRequest, error) {
	return s.prRepo.GetByReviewer(ctx, id)
}

func contains(list []string, v string) bool {
	for _, e := range list {
		if e == v {
			return true
		}
	}
	return false
}
