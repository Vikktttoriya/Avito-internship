package pg_repository

import (
	"context"
	"database/sql"
	"errors"
	"test/internal/domain/model"
	"test/internal/infrastructure/persistence/postgres/pg_mapper"
	"test/internal/infrastructure/persistence/postgres/pg_model"

	sq "github.com/Masterminds/squirrel"
)

type PrRepository struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewPrRepository(db *sql.DB) *PrRepository {
	return &PrRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PrRepository) GetByID(ctx context.Context, id string) (*model.PullRequest, error) {
	query, args, err := r.sb.Select("id", "name", "author_id", "status", "created_at", "merged_at").
		From("pull_requests").
		Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	var dbPR pg_model.PullRequestDb
	if err := row.Scan(&dbPR.ID, &dbPR.Name, &dbPR.AuthorID, &dbPR.Status, &dbPR.CreatedAt, &dbPR.MergedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	reviewerQuery, reviewerArgs, err := r.sb.
		Select("user_id").
		From("pr_reviewers").
		Where(sq.Eq{"pr_id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, reviewerQuery, reviewerArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []*pg_model.UserDb
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, &pg_model.UserDb{ID: uid})
	}

	return pg_mapper.MapPrDbToPr(&dbPR, reviewers), nil
}

func (r *PrRepository) Save(ctx context.Context, pr *model.PullRequest) error {
	dbPR := pg_mapper.MapPrToPrDb(pr)

	query, args, err := r.sb.Insert("pull_requests").
		Columns("id", "name", "author_id", "status", "created_at", "merged_at").
		Values(dbPR.ID, dbPR.Name, dbPR.AuthorID, dbPR.Status, dbPR.CreatedAt, dbPR.MergedAt).
		Suffix("ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, author_id = EXCLUDED.author_id, status = EXCLUDED.status, created_at = EXCLUDED.created_at, merged_at = EXCLUDED.merged_at").
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PrRepository) CheckUserOpenPRs(ctx context.Context, userIDs []string) (bool, error) {
	if len(userIDs) == 0 {
		return false, nil
	}

	authorQuery, authorArgs, err := r.sb.Select("1").
		From("pull_requests").
		Where(sq.Eq{"author_id": userIDs, "status": "OPEN"}).
		Limit(1).ToSql()

	if err != nil {
		return false, err
	}

	row := r.db.QueryRowContext(ctx, authorQuery, authorArgs...)
	var tmp int
	if err := row.Scan(&tmp); err == nil {
		return true, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	reviewerQuery, reviewerArgs, err := r.sb.Select("1").
		From("pr_reviewers AS prr").
		Join("pull_requests AS pr ON pr.id = prr.pr_id").
		Where(sq.Eq{"prr.user_id": userIDs, "pr.status": "OPEN"}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return false, err
	}

	row = r.db.QueryRowContext(ctx, reviewerQuery, reviewerArgs...)
	if err := row.Scan(&tmp); err == nil {
		return true, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else {
		return false, err
	}
}

func (r *PrRepository) GetByReviewer(ctx context.Context, reviewerID string) ([]*model.PullRequest, error) {
	query, args, err := r.sb.
		Select("pr.id", "pr.name", "pr.author_id", "pr.status", "pr.created_at", "pr.merged_at").
		From("pr_reviewers AS prr").
		Join("pull_requests AS pr ON pr.id = prr.pr_id").
		Where(sq.Eq{"prr.user_id": reviewerID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.PullRequest

	for rows.Next() {
		var dbPR pg_model.PullRequestDb
		if err := rows.Scan(&dbPR.ID, &dbPR.Name, &dbPR.AuthorID, &dbPR.Status, &dbPR.CreatedAt, &dbPR.MergedAt); err != nil {
			return nil, err
		}

		reviewerQuery, reviewerArgs, err := r.sb.
			Select("user_id").
			From("pr_reviewers").
			Where(sq.Eq{"pr_id": dbPR.ID}).
			ToSql()
		if err != nil {
			return nil, err
		}

		reviewerRows, err := r.db.QueryContext(ctx, reviewerQuery, reviewerArgs...)
		if err != nil {
			return nil, err
		}

		var reviewers []*pg_model.UserDb
		for reviewerRows.Next() {
			var uid string
			if err := reviewerRows.Scan(&uid); err != nil {
				reviewerRows.Close()
				return nil, err
			}
			reviewers = append(reviewers, &pg_model.UserDb{ID: uid})
		}
		reviewerRows.Close()

		result = append(result, pg_mapper.MapPrDbToPr(&dbPR, reviewers))
	}

	return result, nil
}
