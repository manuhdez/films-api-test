package rabbit

import (
	"context"

	"github.com/manuhdez/films-api-test/internal/domain"
)

type Publisher interface {
	Publish(ctx context.Context, event domain.Event) error
}
