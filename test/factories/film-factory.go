package factories

import (
	"github.com/google/uuid"
	"syreclabs.com/go/faker"

	"github.com/manuhdez/films-api-test/internal/domain/film"
)

func Film() film.Film {
	return film.Film{
		ID:          uuid.New(),
		Title:       faker.Lorem().Sentence(4),
		Director:    faker.Name().Name(),
		ReleaseDate: faker.Time().Birthday(0, 50).Year(),
		Genre:       faker.Lorem().Word(),
		Synopsis:    faker.Lorem().Sentence(25),
		Casting:     filmCasting(),
	}
}

func FilmList(size int) []film.Film {
	var films []film.Film
	for range size {
		films = append(films, Film())
	}
	return films

}

func filmCasting() []string {
	size := faker.Number().Between(1, 5)

	var names []string
	for range size {
		names = append(names, faker.Name().Name())
	}

	return names
}
