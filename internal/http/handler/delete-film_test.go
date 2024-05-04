package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/internal/service"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestDeleteFilm_Handle(t *testing.T) {
	t.Parallel()

	t.Run("returns code 200 if film is deleted", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/film/:id", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		testFilm := factories.Film()
		ctx.SetParamNames("id")
		ctx.SetParamValues(testFilm.ID.String())
		ctx.Set("userID", testFilm.CreatedBy)

		repo := new(mocks.MockFilmRepository)
		deleter := service.NewFilmDeleter(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewDeleteFilm(deleter, finder)

		repo.On("Find", mock.Anything, mock.Anything).Return(testFilm, nil)
		repo.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("return a 400 if film id is not valid", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/film/:id", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		repo := new(mocks.MockFilmRepository)
		deleter := service.NewFilmDeleter(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewDeleteFilm(deleter, finder)

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("return a 404 if film does not exist", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/film/:id", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.New().String())

		repo := new(mocks.MockFilmRepository)
		deleter := service.NewFilmDeleter(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewDeleteFilm(deleter, finder)

		repo.On("Find", mock.Anything, mock.Anything).Return(film.Film{}, ErrFilmNotFound)

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("returns code 401 if user is not the film owner", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/film/:id", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.New().String())

		repo := new(mocks.MockFilmRepository)
		deleter := service.NewFilmDeleter(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewDeleteFilm(deleter, finder)

		testFilm := factories.Film()
		repo.On("Find", mock.Anything, mock.Anything).Return(testFilm, nil)

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("returns code 500 if film can't be deleted", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/film/:id", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		testFilm := factories.Film()
		ctx.SetParamNames("id")
		ctx.SetParamValues(testFilm.ID.String())
		ctx.Set("userID", testFilm.CreatedBy)

		repo := new(mocks.MockFilmRepository)
		deleter := service.NewFilmDeleter(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewDeleteFilm(deleter, finder)

		repo.On("Find", mock.Anything, mock.Anything).Return(testFilm, nil)
		repo.On("Delete", mock.Anything, mock.Anything).Return(errors.New("cannot delete film")).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		repo.AssertExpectations(t)
	})
}
