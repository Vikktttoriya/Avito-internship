package domain_errors

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrTeamNotFound            = errors.New("team not found")
	ErrPullRequestNotFound     = errors.New("pull request not found")
	ErrPullRequestExists       = errors.New("pull request already exists")
	ErrTeamExists              = errors.New("team already exists")
	ErrPRMerged                = errors.New("pull request already merged")
	ErrReviewerNotAssigned     = errors.New("reviewer is not assigned to pull request")
	ErrNoReplacementCandidate  = errors.New("no active candidate available")
	ErrUserHasOpenPullRequests = errors.New("user has open pull requests")
)
