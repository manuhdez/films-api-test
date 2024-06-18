package rabbit

import (
	"encoding/json"
	"log"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/internal/infra"
)

type UserCreatedEvent struct {
	User user.User
}

func NewUserCreatedEvent(u infra.UserJson) UserCreatedEvent {
	return UserCreatedEvent{User: u.ToDomainUser()}
}

func (UserCreatedEvent) Key() string {
	return "api.users.created"
}

func (e UserCreatedEvent) Data() []byte {
	data, err := json.Marshal(infra.NewUserJson(e.User))
	if err != nil {
		log.Printf("Failed to marshal user data: %v", err)
		return nil
	}
	return data
}
