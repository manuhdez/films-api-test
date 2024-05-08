package service

import (
	"context"
	"errors"

	"github.com/manuhdez/films-api-test/internal/domain/film"
)

var (
	ErrUnableToGetFilms = errors.New("unable to get films")
)

type FilmsGetter struct {
	repository film.Repository
}

func NewFilmsGetter(repo film.Repository) FilmsGetter {
	return FilmsGetter{repository: repo}
}

func (f FilmsGetter) Get(ctx context.Context, filter film.Filter) ([]film.Film, error) {
	films, err := f.repository.All(ctx, filter)
	if err != nil {
		return nil, ErrUnableToGetFilms
	}

	return films, nil
}
