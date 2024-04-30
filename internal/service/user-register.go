package service

import (
	"context"

	"github.com/manuhdez/films-api-test/internal/domain/user"
)

type UserRegister struct {
	repository user.Repository
	hasher     PasswordHasher
}

func NewUserRegister(r user.Repository, h PasswordHasher) UserRegister {
	return UserRegister{
		repository: r,
		hasher:     h,
	}
}

func (r UserRegister) Register(c context.Context, u user.User) (user.User, error) {
	hashed, err := r.hasher.Hash(u.Password)
	if err != nil {
		return user.User{}, err
	}

	u.Password = hashed
	return r.repository.Save(c, u)
}
