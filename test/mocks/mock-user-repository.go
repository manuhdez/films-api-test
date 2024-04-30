package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/user"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(c context.Context, u user.User) (user.User, error) {
	args := m.Called(c, u)
	return args.Get(0).(user.User), args.Error(1)
}
