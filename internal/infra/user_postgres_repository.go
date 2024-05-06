package infra

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/manuhdez/films-api-test/internal/domain/user"
)

var (
	ErrFailedToCreateUser   = errors.New("failed to create user")
	ErrUsernameAlreadyInUse = errors.New("username is already in use")
	ErrUsernameNotFound     = errors.New("username does not exist")
)

type PostgresUser struct {
	ID       uuid.UUID `sql:"id"`
	Username string    `sql:"username"`
	Password string    `sql:"password"`
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(c context.Context, u user.User) (user.User, error) {
	q := psql.Insert("users").Columns("id", "username", "password").Values(u.ID, u.Username, u.Password)
	query, _, queryErr := q.ToSql()
	if queryErr != nil {
		slog.Error("failed to generate query", "error", queryErr)
		return user.User{}, queryErr
	}

	_, err := r.db.ExecContext(c, query, u.ID, u.Username, u.Password)
	if err != nil {
		slog.Error("failed to save user", "err", err)

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return user.User{}, ErrUsernameAlreadyInUse
			default:
				return user.User{}, ErrFailedToCreateUser
			}
		}

		return user.User{}, err
	}

	return u, nil
}

func (r *PostgresUserRepository) SearchByUsername(c context.Context, username string) (user.User, error) {
	query, _, err := psql.Select("*").From("users").Where(sq.Eq{"username": username}).ToSql()
	if err != nil {
		slog.Error("failed to generate query", "error", err)
		return user.User{}, err
	}

	row := r.db.QueryRowContext(c, query, username)

	if row.Err() != nil {
		slog.Error("failed to search user by username", "username", username, "err", row.Err())
		return user.User{}, row.Err()
	}

	var u PostgresUser
	scanErr := row.Scan(&u.ID, &u.Username, &u.Password)
	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return user.User{}, ErrUsernameNotFound
		}

		slog.Error("failed to scan user row", "username", username, "err", scanErr)
		return user.User{}, scanErr
	}

	return user.User{ID: u.ID, Username: u.Username, Password: u.Password}, nil
}

func (r *PostgresUserRepository) Find(c context.Context, id uuid.UUID) (user.User, error) {
	query, _, err := psql.Select("id, username").From("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		slog.Error("failed to generate query", "error", err.Error())
		return user.User{}, err
	}

	row := r.db.QueryRowContext(c, query, id)
	if row.Err() != nil {
		slog.Error("failed to find user by id", "id", id, "err", row.Err())
	}

	var u PostgresUser
	scanErr := row.Scan(&u.ID, &u.Username)
	if scanErr != nil {
		slog.Info("no user found by id", "id", id)
		return user.User{}, scanErr
	}

	return user.User{ID: u.ID, Username: u.Username}, nil
}
