package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestFilmsGetter(t *testing.T) {
	t.Parallel()

	t.Run("returns a list of films", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		getter := NewFilmsGetter(repo)
		testFilms := factories.FilmList(5)
		repo.On("All", mock.Anything, mock.Anything).Return(testFilms, nil).Once()

		filter := film.Filter{}
		films, err := getter.Get(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, testFilms, films)
		repo.AssertExpectations(t)
	})

	t.Run("returns error if films cannot be read", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		getter := NewFilmsGetter(repo)
		repo.On("All", mock.Anything, mock.Anything).Return([]film.Film{}, ErrUnableToGetFilms).Once()

		filter := film.Filter{}
		_, err := getter.Get(context.Background(), filter)
		assert.ErrorIs(t, err, ErrUnableToGetFilms)
		repo.AssertExpectations(t)
	})
}
