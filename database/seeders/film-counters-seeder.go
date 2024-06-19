package seeders

import (
	"context"
	"log"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/internal/domain/userfilms"
)

func (s Seeder) SeedFilmCounters(users []user.User) {
	log.Println("Seeding film counters")

	for _, u := range users {
		c := userfilms.UserFilms{UserId: u.ID, Films: 1}
		_ = s.repos.filmCounterRepository.Create(context.Background(), c)
	}
}
