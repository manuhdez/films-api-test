package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/manuhdez/films-api-test/internal/domain/userfilms"
)

type UserFilmsCounterIncrementer struct {
	repository userfilms.Repository
}

func NewUserFilmsCounterIncrementer(r userfilms.Repository) UserFilmsCounterIncrementer {
	return UserFilmsCounterIncrementer{repository: r}
}

func (u UserFilmsCounterIncrementer) Increment(ctx context.Context, userId uuid.UUID) error {
	log.Printf("Incrementing user film counter for user id %v", userId)
	return u.repository.Increment(ctx, userId)
}
