package service

import (
	"context"
	"errors"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/internal/infra"
)

var (
	ErrWrongCredentials = errors.New("invalid credentials")
)

type UserLogin struct {
	repository user.Repository
	hasher     PasswordHasher
}

func NewUserLogin(r user.Repository, h PasswordHasher) UserLogin {
	return UserLogin{
		repository: r,
		hasher:     h,
	}
}

func (srv UserLogin) Login(username, password string) (user.User, error) {
	c := context.Background()
	u, searchErr := srv.repository.SearchByUsername(c, username)
	if searchErr != nil {
		if errors.Is(searchErr, infra.ErrUsernameNotFound) {
			return user.User{}, ErrWrongCredentials
		}

		return user.User{}, searchErr
	}

	if match := srv.hasher.Compare(u.Password, password); match == false {
		return user.User{}, ErrWrongCredentials
	}

	return u, nil
}
