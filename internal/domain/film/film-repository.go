package film

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	All(ctx context.Context) ([]Film, error)
	Find(ctx context.Context, id uuid.UUID) (Film, error)
	Save(ctx context.Context, film Film) error
	Delete(ctx context.Context, id uuid.UUID) error
}
