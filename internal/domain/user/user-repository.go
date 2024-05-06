package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(context.Context, User) (User, error)
	SearchByUsername(context.Context, string) (User, error)
	Find(context.Context, uuid.UUID) (User, error)
}
