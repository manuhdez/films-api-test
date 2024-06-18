package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/manuhdez/films-api-test/internal/domain"
	"github.com/manuhdez/films-api-test/internal/domain/user"
)

type UserCreatedEvent struct {
	UserID   string `json:"user_id,omitempty"`
	UserName string `json:"username,omitempty"`
}

func NewUserCreatedEvent(id, username string) UserCreatedEvent {
	return UserCreatedEvent{UserID: id, UserName: username}
}

func (UserCreatedEvent) Key() string {
	return "api.users.created"
}

func (e UserCreatedEvent) Data() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		log.Printf("Failed to marshal user data: %v", err)
		return nil
	}
	return data
}

type UserRegister struct {
	repository user.Repository
	eventBus   domain.EventBus
	hasher     PasswordHasher
}

func NewUserRegister(r user.Repository, h PasswordHasher, b domain.EventBus) UserRegister {
	return UserRegister{
		repository: r,
		hasher:     h,
		eventBus:   b,
	}
}

func (r UserRegister) Register(ctx context.Context, u user.User) (user.User, error) {
	hashed, err := r.hasher.Hash(u.Password)
	if err != nil {
		return user.User{}, err
	}

	u.Password = hashed
	savedUser, err := r.repository.Save(ctx, u)
	if err != nil {
		log.Println("error registering user:", err)
		return user.User{}, err
	}

	if pubErr := r.eventBus.Publish(ctx, UserCreatedEvent{
		UserID:   savedUser.ID.String(),
		UserName: savedUser.Username,
	}); pubErr != nil {
		log.Println("error publishing user created event")
	}

	return savedUser, nil
}
