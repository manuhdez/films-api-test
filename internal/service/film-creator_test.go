package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestFilmCreator_Create(t *testing.T) {
	t.Parallel()

	t.Run("returns no error if film is created successfully", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		creator := NewFilmCreator(repo)

		testFilm := factories.Film()
		repo.On("Save", mock.Anything, testFilm).Return(nil).Once()

		err := creator.Create(testFilm)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("returns error if cannot create film", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		creator := NewFilmCreator(repo)

		testFilm := factories.Film()
		testErr := errors.New("cannot save film")
		repo.On("Save", mock.Anything, testFilm).Return(testErr).Once()

		err := creator.Create(testFilm)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
