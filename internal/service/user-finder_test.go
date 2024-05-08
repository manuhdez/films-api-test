package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestUserFinder_Find(t *testing.T) {
	t.Parallel()

	t.Run("returns a user", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		finder := NewUserFinder(repo)

		fakeUser := factories.User()
		repo.On("Find", mock.Anything, mock.Anything).Return(fakeUser, nil)

		usr, err := finder.Find(context.Background(), fakeUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, fakeUser, usr)
		repo.AssertExpectations(t)
	})

	t.Run("returns an error", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		finder := NewUserFinder(repo)

		repo.On("Find", mock.Anything, mock.Anything).Return(user.User{}, errors.New("not found"))

		_, err := finder.Find(context.Background(), uuid.New())
		assert.ErrorIs(t, err, ErrUserNotFound)
		repo.AssertExpectations(t)
	})
}
