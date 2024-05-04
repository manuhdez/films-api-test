package seeders

import (
	"context"
	"log"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/test/factories"
)

func (s Seeder) SeedFilms(users []user.User) {
	log.Println("Seeding films...")

	for _, u := range users {
		f := factories.Film()
		f.CreatedBy = u.ID
		_, _ = s.repos.filmRepository.Save(context.Background(), f)
	}
}
