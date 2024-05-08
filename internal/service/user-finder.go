package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/manuhdez/films-api-test/internal/domain/user"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserFinder struct {
	repository user.Repository
}

func NewUserFinder(repository user.Repository) UserFinder {
	return UserFinder{repository: repository}
}

func (f UserFinder) Find(ctx context.Context, id uuid.UUID) (user.User, error) {
	u, err := f.repository.Find(ctx, id)
	if err != nil {
		slog.Error("failed to find user", "id", id.String(), "err", err.Error())
		return user.User{}, ErrUserNotFound
	}

	return u, nil
}
