package infra

import (
	"context"
	"database/sql"
	"errors"
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
	current, err := r.Count(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to updated counter: %e", err)
	}

	response := r.db.WithContext(ctx).Model(&current).Where("user_id = ?", userId).Update("films", current.Films+1)
	if response.Error != nil {
		return fmt.Errorf("failed to increment count: %v", response.Error)
	}

	return nil
}

func (r UserFilmsPostgresRepository) Count(ctx context.Context, userId uuid.UUID) (userfilms.UserFilms, error) {
	var counter userfilms.UserFilms
	result := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&counter)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return userfilms.UserFilms{}, fmt.Errorf("counter not found: %v", result.Error)
		}

		return userfilms.UserFilms{}, fmt.Errorf("error retrieving counter: %v", result.Error)
	}

	return counter, nil
}
