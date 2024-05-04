package factories

import (
	"github.com/google/uuid"
	"syreclabs.com/go/faker"

	"github.com/manuhdez/films-api-test/internal/domain/user"
)

func User() user.User {
	return user.User{
		ID:       uuid.New(),
		Username: faker.Internet().UserName(),
		Password: faker.Internet().Password(
			user.PasswordMinLength,
			user.PasswordMaxLength,
		),
	}
}

func UserList(size int) []user.User {
	var users []user.User
	for range size {
		users = append(users, User())
	}
	return users
}
