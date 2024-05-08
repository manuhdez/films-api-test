package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestFilmFinder_Find(t *testing.T) {
	t.Parallel()

	t.Run("returns a film", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		finder := NewFilmFinder(repo)

		testFilm := factories.Film()
		repo.On("Find", mock.Anything, mock.Anything).Return(testFilm, nil)

		id := uuid.New()
		f, err := finder.Find(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, f, testFilm)
		repo.AssertExpectations(t)
	})

	t.Run("returns an error if the film cannot be found", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		finder := NewFilmFinder(repo)
		repo.On("Find", mock.Anything, mock.Anything).Return(film.Film{}, ErrFilmNotFound)

		id := uuid.New()
		_, err := finder.Find(context.Background(), id)
		assert.ErrorIs(t, err, ErrFilmNotFound)
		repo.AssertExpectations(t)
	})
}
