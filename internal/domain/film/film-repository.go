package film

import "context"

type Repository interface {
	All(ctx context.Context) ([]Film, error)
}
