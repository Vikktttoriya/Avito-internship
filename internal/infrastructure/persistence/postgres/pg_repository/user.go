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

type UserRepository struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	q := r.sb.
		Select("id", "username", "team_name", "is_active").
		From("users").
		Where(sq.Eq{"id": id})

	row := q.RunWith(r.db).QueryRowContext(ctx)
	var dbUser pg_model.UserDb

	err := row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.TeamName, &dbUser.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return pg_mapper.MapUserDbToUser(&dbUser), nil
}

func (r *UserRepository) Save(ctx context.Context, u *model.User) error {
	dbUser := pg_mapper.MapUserToUserDb(u)

	query, args, err := r.sb.Insert("users").
		Columns("id", "username", "team_name", "is_active").
		Values(dbUser.ID, dbUser.Username, dbUser.TeamName, dbUser.IsActive).
		Suffix("ON CONFLICT (id) DO UPDATE SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active").
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepository) SaveTx(ctx context.Context, tx *sql.Tx, u *model.User) error {
	dbUser := pg_mapper.MapUserToUserDb(u)

	query, args, err := r.sb.Insert("users").
		Columns("id", "username", "team_name", "is_active").
		Values(dbUser.ID, dbUser.Username, dbUser.TeamName, dbUser.IsActive).
		Suffix("ON CONFLICT (id) DO UPDATE SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active").
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepository) GetActiveByTeam(ctx context.Context, team string) ([]*model.User, error) {
	q := r.sb.
		Select("id", "username", "team_name", "is_active").
		From("users").
		Where(sq.Eq{"team_name": team, "is_active": true})

	rows, err := q.RunWith(r.db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var dbUser pg_model.UserDb
		if err := rows.Scan(&dbUser.ID, &dbUser.Username, &dbUser.TeamName, &dbUser.IsActive); err != nil {
			return nil, err
		}
		users = append(users, pg_mapper.MapUserDbToUser(&dbUser))
	}
	return users, nil
}

func (r *UserRepository) GetByTeam(ctx context.Context, team string) ([]*model.User, error) {
	q := r.sb.
		Select("id", "username", "team_name", "is_active").
		From("users").
		Where(sq.Eq{"team_name": team})

	rows, err := q.RunWith(r.db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var dbUser pg_model.UserDb
		if err := rows.Scan(&dbUser.ID, &dbUser.Username, &dbUser.TeamName, &dbUser.IsActive); err != nil {
			return nil, err
		}
		users = append(users, pg_mapper.MapUserDbToUser(&dbUser))
	}
	return users, nil
}
