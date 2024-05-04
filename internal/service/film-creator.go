package service

import (
	"context"
	"errors"

	"github.com/manuhdez/films-api-test/internal/domain/film"
)

type FilmCreator struct {
	repository film.Repository
}

func NewFilmCreator(r film.Repository) FilmCreator {
	return FilmCreator{repository: r}
}

func (fc FilmCreator) Create(f film.Film) error {
	ctx := context.Background()
	err := fc.repository.Save(ctx, f)
	if err != nil {
		return errors.New("failed to save film")
	}

	return nil
}
