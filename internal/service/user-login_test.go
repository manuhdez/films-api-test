package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestUserLogin(t *testing.T) {
	t.Run("returns user if credentials are valid", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		hasher := new(mocks.MockPasswordHasher)
		logger := NewUserLogin(repo, hasher)

		ctx := context.Background()
		credentials := user.User{Username: "username", Password: "password"}
		repo.On("SearchByUsername", ctx, credentials.Username).Return(credentials, nil)
		hasher.On("Compare", mock.Anything, credentials.Password).Return(true)

		u, err := logger.Login(credentials.Username, credentials.Password)
		assert.NoError(t, err)
		assert.Equal(t, credentials, u)

		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("returns error if username does not exist", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		hasher := new(mocks.MockPasswordHasher)
		logger := NewUserLogin(repo, hasher)

		ctx := context.Background()
		credentials := user.User{Username: "username", Password: "password"}
		repo.On("SearchByUsername", ctx, credentials.Username).Return(credentials, ErrWrongCredentials)

		_, err := logger.Login(credentials.Username, credentials.Password)
		assert.ErrorIs(t, err, ErrWrongCredentials)

		repo.AssertExpectations(t)
	})

	t.Run("returns error if password does not match", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		hasher := new(mocks.MockPasswordHasher)
		logger := NewUserLogin(repo, hasher)

		ctx := context.Background()
		credentials := user.User{Username: "username", Password: "password"}
		repo.On("SearchByUsername", ctx, credentials.Username).Return(credentials, nil)
		hasher.On("Compare", mock.Anything, credentials.Password).Return(false)

		_, err := logger.Login(credentials.Username, credentials.Password)
		assert.ErrorIs(t, err, ErrWrongCredentials)

		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})
}
