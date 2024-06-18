package infra

import (
	"github.com/google/uuid"

	"github.com/manuhdez/films-api-test/internal/domain/user"
)

type UserJson struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func NewUserJson(u user.User) UserJson {
	return UserJson{
		ID:       u.ID.String(),
		Username: u.Username,
	}
}

func (u UserJson) ToDomainUser() user.User {
	return user.User{
		ID:       uuid.MustParse(u.ID),
		Username: u.Username,
	}
}
