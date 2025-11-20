package mapper

import (
	"test/internal/api"
	"test/internal/domain/model"
)

func ToAPIUser(u *model.User) api.User {
	if u == nil {
		return api.User{}
	}

	return api.User{
		UserId:   u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func ToAPITeamMember(u *model.User) api.TeamMember {
	if u == nil {
		return api.TeamMember{}
	}

	return api.TeamMember{
		UserId:   u.ID,
		Username: u.Username,
		IsActive: u.IsActive,
	}
}

func ToAPITeam(team *model.Team) api.Team {
	if team == nil {
		return api.Team{}
	}

	apiMembers := make([]api.TeamMember, 0, len(team.Members))
	for _, m := range team.Members {
		apiMembers = append(apiMembers, ToAPITeamMember(m))
	}

	return api.Team{
		TeamName: team.Name,
		Members:  apiMembers,
	}
}

func ToAPIPullRequest(pr *model.PullRequest) api.PullRequest {
	if pr == nil {
		return api.PullRequest{}
	}

	status := api.PullRequestStatusOPEN
	if pr.Status == model.StatusMerged {
		status = api.PullRequestStatusMERGED
	}

	return api.PullRequest{
		PullRequestId:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorId:          pr.AuthorID,
		AssignedReviewers: pr.AssignedReviewers,
		Status:            status,
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func ToAPIPullRequestShort(pr *model.PullRequest) api.PullRequestShort {
	if pr == nil {
		return api.PullRequestShort{}
	}

	status := api.PullRequestShortStatusOPEN
	if pr.Status == model.StatusMerged {
		status = api.PullRequestShortStatusMERGED
	}

	return api.PullRequestShort{
		PullRequestId:   pr.ID,
		PullRequestName: pr.Name,
		AuthorId:        pr.AuthorID,
		Status:          status,
	}
}
