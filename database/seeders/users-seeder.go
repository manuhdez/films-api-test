package seeders

import (
	"context"
	"log"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/test/factories"
)

func (s Seeder) SeedUsers() []user.User {
	log.Println("Seeding users...")

	passwordHasher := infra.NewBcryptPasswordHasher()
	hashed, err := passwordHasher.Hash("password")
	if err != nil {
		log.Fatalf("could not hash users password: %v", err)
	}

	users := factories.UserList(5)
	for idx, u := range users {
		if idx == 0 {
			u.Username = "test"
		}

		u.Password = hashed
		_, _ = s.repos.userRepository.Save(context.Background(), u)
	}

	return users
}
