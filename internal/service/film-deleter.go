package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/manuhdez/films-api-test/internal/domain/film"
)

var (
	ErrCannotDeleteFilm = errors.New("cannot delete film")
)

type FilmDeleter struct {
	repository film.Repository
}

func NewFilmDeleter(r film.Repository) FilmDeleter {
	return FilmDeleter{repository: r}
}

func (d FilmDeleter) Delete(ctx context.Context, filmID uuid.UUID) error {
	err := d.repository.Delete(ctx, filmID)
	if err != nil {
		return ErrCannotDeleteFilm
	}
	return nil
}
