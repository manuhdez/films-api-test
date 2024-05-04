package seeders

import (
	"database/sql"
	"log"

	"github.com/manuhdez/films-api-test/database"
	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/internal/infra"
)

type Repositories struct {
	userRepository user.Repository
	filmRepository film.Repository
}

type Seeder struct {
	db    *sql.DB
	repos Repositories
}

func NewSeeder() Seeder {
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("failed to connect to postgres: %e", err)
	}

	repos := Repositories{
		userRepository: infra.NewPostgresUserRepository(db),
		filmRepository: infra.NewPostgresFilmRepository(db),
	}

	return Seeder{
		db:    db,
		repos: repos,
	}
}

func (s Seeder) Seed() {
	users := s.SeedUsers()
	s.SeedFilms(users)
}
