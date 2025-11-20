package pg_mapper

import (
	"test/internal/domain/model"
	"test/internal/infrastructure/persistence/postgres/pg_model"
)

func MapUserToUserDb(u *model.User) *pg_model.UserDb {
	return &pg_model.UserDb{
		ID:       u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func MapUserDbToUser(u *pg_model.UserDb) *model.User {
	return &model.User{
		ID:       u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func MapTeamToTeamDb(t *model.Team) *pg_model.TeamDb {
	return &pg_model.TeamDb{
		Name: t.Name,
	}
}

func MapTeamDbToTeam(t *pg_model.TeamDb, members []*model.User) *model.Team {
	return &model.Team{
		Name:    t.Name,
		Members: members,
	}
}

func MapPrToPrDb(pr *model.PullRequest) *pg_model.PullRequestDb {
	return &pg_model.PullRequestDb{
		ID:        pr.ID,
		Name:      pr.Name,
		AuthorID:  pr.AuthorID,
		Status:    string(pr.Status),
		CreatedAt: pr.CreatedAt,
		MergedAt:  pr.MergedAt,
	}
}

func MapPrDbToPr(prDb *pg_model.PullRequestDb, reviewers []*pg_model.UserDb) *model.PullRequest {
	reviewerIDs := make([]string, len(reviewers))
	for i, r := range reviewers {
		reviewerIDs[i] = r.ID
	}
	return &model.PullRequest{
		ID:                prDb.ID,
		Name:              prDb.Name,
		AuthorID:          prDb.AuthorID,
		Status:            model.Status(prDb.Status),
		CreatedAt:         prDb.CreatedAt,
		MergedAt:          prDb.MergedAt,
		AssignedReviewers: reviewerIDs,
	}
}
