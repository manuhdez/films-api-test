package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestFilmsGetter(t *testing.T) {
	testFilms := factories.FilmList(5)

	t.Run("returns a list of films", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		getter := NewFilmsGetter(repo)
		repo.On("All", mock.Anything).Return(testFilms, nil).Once()

		films, err := getter.Get()
		assert.NoError(t, err)
		assert.Equal(t, testFilms, films)
		repo.AssertExpectations(t)
	})

	t.Run("returns error if films cannot be read", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		getter := NewFilmsGetter(repo)
		repo.On("All", mock.Anything).Return([]film.Film{}, ErrUnableToGetFilms).Once()

		_, err := getter.Get()
		assert.ErrorIs(t, err, ErrUnableToGetFilms)
		repo.AssertExpectations(t)
	})
}
