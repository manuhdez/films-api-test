package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestUserRegister(t *testing.T) {
	testUser := user.New(uuid.New(), "fake-username", "fake-password")

	t.Run("returns the user saved in the repository", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		hasher := new(mocks.MockPasswordHasher)
		register := NewUserRegister(repo, hasher)

		ctx := context.Background()
		repo.On("Save", ctx, testUser).Return(testUser, nil)
		hasher.On("Hash", mock.Anything).Return(testUser.Password, nil)

		u, err := register.Register(ctx, testUser)
		assert.NoError(t, err)
		assert.Equal(t, testUser, u)
	})

	t.Run("returns an error if the repository fails to save user", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		hasher := new(mocks.MockPasswordHasher)
		register := NewUserRegister(repo, hasher)

		ctx := context.Background()
		testErr := errors.New("failed to save user")
		repo.On("Save", ctx, testUser).Return(user.User{}, testErr)
		hasher.On("Hash", mock.Anything).Return(testUser.Password, nil)

		_, err := register.Register(ctx, testUser)
		assert.ErrorIs(t, err, testErr)
	})

	t.Run("returns an error if password hashing fails", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		hasher := new(mocks.MockPasswordHasher)
		register := NewUserRegister(repo, hasher)

		ctx := context.Background()
		hashErr := errors.New("failed to hash password")
		hasher.On("Hash", mock.Anything).Return("", hashErr)

		_, err := register.Register(ctx, testUser)
		assert.ErrorIs(t, err, hashErr)
	})
}
