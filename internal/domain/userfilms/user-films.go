package userfilms

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, userFilms UserFilms) error
	Increment(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, id uuid.UUID) (UserFilms, error)
}

type UserFilms struct {
	UserId uuid.UUID
	Films  uint
}

func New(id uuid.UUID, films uint) UserFilms {
	return UserFilms{UserId: id, Films: films}
}
