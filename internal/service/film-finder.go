package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/manuhdez/films-api-test/internal/domain/film"
)

var (
	ErrFilmNotFound = errors.New("film not found")
)

type FilmFinder struct {
	repository film.Repository
}

func NewFilmFinder(r film.Repository) FilmFinder {
	return FilmFinder{repository: r}
}

func (finder FilmFinder) Find(filmId uuid.UUID) (film.Film, error) {
	ctx := context.Background()
	f, err := finder.repository.Find(ctx, filmId)
	if err != nil {
		slog.Error("failed to find film", "id", filmId.String(), "error", err.Error())
		return film.Film{}, ErrFilmNotFound
	}

	return f, nil
}
