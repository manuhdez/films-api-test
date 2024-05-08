package service

import (
	"context"
	"errors"

	"github.com/manuhdez/films-api-test/internal/domain/film"
)

var (
	ErrCannotUpdateFilm = errors.New("cannot update film")
)

type FilmUpdater struct {
	repo film.Repository
}

func NewFilmUpdater(repo film.Repository) FilmUpdater {
	return FilmUpdater{repo: repo}
}

func (f FilmUpdater) Update(ctx context.Context, film film.Film) error {
	err := f.repo.Update(ctx, film)
	if err != nil {
		return ErrCannotUpdateFilm
	}

	return nil
}
