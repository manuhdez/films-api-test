package service

import (
	"context"
	"errors"
	"log"

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

func (srv UserLogin) Login(ctx context.Context, username, password string) (user.User, error) {
	u, searchErr := srv.repository.SearchByUsername(ctx, username)
	if searchErr != nil {
		if errors.Is(searchErr, infra.ErrUsernameNotFound) {
			return user.User{}, ErrWrongCredentials
		}

		return user.User{}, searchErr
	}
	if username != u.Username {
		log.Println("username mismatch", u.Username)
		return user.User{}, ErrWrongCredentials
	}

	if match := srv.hasher.Compare(u.Password, password); !match {
		return user.User{}, ErrWrongCredentials
	}

	return u, nil
}
