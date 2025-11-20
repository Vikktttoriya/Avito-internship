package pg_model

import "time"

type PullRequestDb struct {
	ID        string
	Name      string
	AuthorID  string
	Status    string
	CreatedAt time.Time
	MergedAt  *time.Time
}
