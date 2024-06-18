package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/manuhdez/films-api-test/internal/domain/userfilms"
)

type UserFilmsPostgresRepository struct {
	db *gorm.DB
}

func NewUserFilmsPostgresRepository(db *sql.DB) *UserFilmsPostgresRepository {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	return &UserFilmsPostgresRepository{db: gormDB}
}

func (r UserFilmsPostgresRepository) Create(ctx context.Context, userFilms userfilms.UserFilms) error {
	result := r.db.WithContext(ctx).Create(&userFilms)
	if result.Error != nil {
		return fmt.Errorf("failed to create userFilms: %v", result.Error)
	}
	return nil
}

func (r UserFilmsPostgresRepository) Increment(ctx context.Context, userId uuid.UUID) error {
	var userFilms = userfilms.UserFilms{UserId: userId}
	current := r.db.WithContext(ctx).First(&userFilms)
	if current.Error != nil {
		return fmt.Errorf("cannot increment count, user has no counter: %v", current.Error)
	}

	userFilms.Films = userFilms.Films + 1
	response := r.db.WithContext(ctx).Update("films", &userFilms)
	if response.Error != nil {
		return fmt.Errorf("failed to increment count: %v", response.Error)
	}

	return nil
}

func (r UserFilmsPostgresRepository) Count(ctx context.Context, userId uuid.UUID) (userfilms.UserFilms, error) {
	var userFilms = userfilms.UserFilms{UserId: userId}
	current := r.db.WithContext(ctx).First(&userFilms)
	if current.Error != nil {
		return userfilms.UserFilms{}, fmt.Errorf("userfilms entry not found: %v", current.Error)
	}

	return userFilms, nil
}
