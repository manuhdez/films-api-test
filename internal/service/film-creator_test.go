package service

import (
	"context"
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
		bus := new(mocks.MockEventBus)
		creator := NewFilmCreator(repo, bus)

		testFilm := factories.Film()
		repo.On("Save", mock.Anything, testFilm).Return(nil).Once()
		bus.On("Publish", mock.Anything, mock.Anything).Return(nil).Once()

		err := creator.Create(context.Background(), testFilm)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("returns error if cannot create film", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		bus := new(mocks.MockEventBus)
		creator := NewFilmCreator(repo, bus)

		testFilm := factories.Film()
		testErr := errors.New("cannot save film")
		repo.On("Save", mock.Anything, testFilm).Return(testErr).Once()

		err := creator.Create(context.Background(), testFilm)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
