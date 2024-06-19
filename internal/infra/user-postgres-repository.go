package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

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

var _ user.Repository = (*PostgresUserRepository)(nil)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: db}),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	return &PostgresUserRepository{db: gormDB}
}

func (r *PostgresUserRepository) Save(c context.Context, u user.User) (user.User, error) {
	response := r.db.WithContext(c).Create(&u)
	if response.Error != nil {
		return user.User{}, fmt.Errorf("[%e]:%w", ErrFailedToCreateUser, response.Error)
	}

	return u, nil
}

func (r *PostgresUserRepository) SearchByUsername(c context.Context, username string) (user.User, error) {
	var u user.User
	result := r.db.WithContext(c).Where("username = ?", username).First(&u)
	if result.Error != nil {
		return user.User{}, fmt.Errorf("[Username not found]:%w", result.Error)
	}

	return u, nil
}

func (r *PostgresUserRepository) Find(c context.Context, id uuid.UUID) (user.User, error) {
	var u user.User
	result := r.db.WithContext(c).First(&u, id)
	if result.Error != nil {
		return user.User{}, fmt.Errorf("[UserNotFound]:%w", result.Error)
	}

	return u, nil
}
