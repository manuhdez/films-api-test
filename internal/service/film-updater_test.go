package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestFilmUpdater_Update(t *testing.T) {
	t.Run("returns no error if the film has been updated", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		updater := NewFilmUpdater(repo)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		testFilm := factories.Film()
		err := updater.Update(testFilm)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("returns error if the film can't be updated", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		updater := NewFilmUpdater(repo)
		repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("failed")).Once()

		testFilm := factories.Film()
		err := updater.Update(testFilm)
		assert.ErrorIs(t, err, ErrCannotUpdateFilm)
		repo.AssertExpectations(t)
	})
}
