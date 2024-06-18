package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain"
)

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, evt domain.Event) error {
	args := m.Called(ctx, evt)
	return args.Error(0)
}
