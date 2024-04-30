package infra

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/manuhdez/films-api-test/internal/domain/user"
)

var (
	ErrFailedToCreateUser   = errors.New("failed to create user")
	ErrUsernameAlreadyInUse = errors.New("username is already in use")
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
	query := `INSERT INTO users (id, username, password) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, u.ID, u.Username, u.Password)
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
