package service

import (
	"context"
	"log"

	"github.com/manuhdez/films-api-test/internal/domain/userfilms"
)

type UserFilmsCounterCreator struct {
	repository userfilms.Repository
}

func NewUserFilmsCounterCreator(r userfilms.Repository) UserFilmsCounterCreator {
	return UserFilmsCounterCreator{repository: r}
}

func (c UserFilmsCounterCreator) Create(ctx context.Context, userFilms userfilms.UserFilms) error {
	log.Printf("creating counter for user %s", userFilms.UserId.String())
	return c.repository.Create(ctx, userFilms)
}
