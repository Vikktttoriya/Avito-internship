package model

import "time"

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusMerged Status = "MERGED"
)

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            Status
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          *time.Time
}

func NewPr(id, name, author string) *PullRequest {
	return &PullRequest{
		ID:                id,
		Name:              name,
		AuthorID:          author,
		Status:            StatusOpen,
		CreatedAt:         time.Now(),
		AssignedReviewers: []string{},
	}
}

func (pr *PullRequest) Merge() {
	if pr.Status == StatusMerged {
		return
	}
	pr.Status = StatusMerged
	t := time.Now()
	pr.MergedAt = &t
}

func (pr *PullRequest) ReplaceReviewer(old, new string) {
	for i, r := range pr.AssignedReviewers {
		if r == old {
			pr.AssignedReviewers[i] = new
			return
		}
	}
}
