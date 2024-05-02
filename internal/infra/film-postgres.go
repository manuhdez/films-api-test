package infra

import (
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/manuhdez/films-api-test/internal/domain/film"
)

type PostgresFilm struct {
	ID          uuid.UUID      `sql:"id"`
	Title       string         `sql:"title"`
	Director    string         `sql:"director"`
	ReleaseDate int            `sql:"release_date"`
	Genre       string         `sql:"genre"`
	Synopsis    string         `sql:"synopsis"`
	Casting     pq.StringArray `sql:"casting"`
	CreatedBy   uuid.UUID      `sql:"created_by"`
}

func (f PostgresFilm) ToDomain() film.Film {
	return film.New(
		f.ID,
		f.Title,
		f.Director,
		f.ReleaseDate,
		f.Genre,
		f.Synopsis,
		f.Casting,
		f.CreatedBy,
	)
}
