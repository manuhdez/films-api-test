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

func (f FilmsGetter) Get() ([]film.Film, error) {
	ctx := context.Background()
	films, err := f.repository.All(ctx)
	if err != nil {
		return nil, ErrUnableToGetFilms
	}

	return films, nil
}
