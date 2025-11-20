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

type TeamRepository struct {
	db       *sql.DB
	sb       sq.StatementBuilderType
	userRepo *UserRepository
}

func NewTeamRepository(db *sql.DB, userRepo *UserRepository) *TeamRepository {
	return &TeamRepository{
		db:       db,
		sb:       sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		userRepo: userRepo,
	}
}

func (r *TeamRepository) Create(ctx context.Context, team *model.Team) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	teamDb := pg_mapper.MapTeamToTeamDb(team)

	query, args, err := r.sb.Insert("teams").
		Columns("name").
		Values(teamDb.Name).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	for _, u := range team.Members {
		u.TeamName = team.Name
		if err = r.userRepo.SaveTx(ctx, tx, u); err != nil {
			return err
		}
	}

	return nil
}

func (r *TeamRepository) GetByName(ctx context.Context, name string) (*model.Team, error) {
	query, args, err := r.sb.Select("name").
		From("teams").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRowContext(ctx, query, args...)
	var teamDb pg_model.TeamDb

	if err := row.Scan(&teamDb.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	users, err := r.userRepo.GetByTeam(ctx, name)
	if err != nil {
		return nil, err
	}

	team := pg_mapper.MapTeamDbToTeam(&teamDb, users)

	return team, nil
}
