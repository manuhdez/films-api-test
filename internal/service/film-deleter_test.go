package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestFilmDeleter_Delete(t *testing.T) {
	t.Parallel()

	t.Run("returns no error if the film has been deleted", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		deleter := NewFilmDeleter(repo)
		repo.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()

		err := deleter.Delete(uuid.New())
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("returns error if the film cannot be deleted", func(t *testing.T) {
		repo := new(mocks.MockFilmRepository)
		deleter := NewFilmDeleter(repo)
		repo.On("Delete", mock.Anything, mock.Anything).Return(ErrCannotDeleteFilm).Once()

		err := deleter.Delete(uuid.New())
		assert.ErrorIs(t, err, ErrCannotDeleteFilm)
		repo.AssertExpectations(t)
	})
}
