package handler

import (
	"bytes"
	"encoding/json"
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

func TestPutFilm_Handle(t *testing.T) {
	t.Parallel()

	t.Run("return code 200 if film is updated", func(t *testing.T) {
		e := echo.New()

		testFilm := factories.Film()
		body, marshalErr := json.Marshal(PutFilmRequest{
			Title:       testFilm.Title,
			Director:    testFilm.Director,
			ReleaseDate: testFilm.ReleaseDate,
			Genre:       testFilm.Genre,
			Synopsis:    testFilm.Synopsis,
			Casting:     testFilm.Casting,
		})
		assert.NoError(t, marshalErr)

		req := httptest.NewRequest(http.MethodPut, "/films/:id", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.New().String())
		ctx.Set("userID", testFilm.CreatedBy)

		repo := new(mocks.MockFilmRepository)
		updater := service.NewFilmUpdater(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewPutFilm(finder, updater)

		repo.On("Find", mock.Anything, mock.Anything).Return(testFilm, nil).Once()
		repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("return 400 with an invalid film id", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/films/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		repo := new(mocks.MockFilmRepository)
		updater := service.NewFilmUpdater(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewPutFilm(finder, updater)

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("return 400 if there's a missing field in the request body", func(t *testing.T) {
		e := echo.New()
		testFilm := factories.Film()
		body, marshalErr := json.Marshal(PutFilmRequest{
			Title:       testFilm.Title,
			Director:    testFilm.Director,
			ReleaseDate: testFilm.ReleaseDate,
			Genre:       testFilm.Genre,
		})
		assert.NoError(t, marshalErr)

		req := httptest.NewRequest(http.MethodPut, "/films/:id", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.New().String())
		ctx.Set("userID", testFilm.CreatedBy)

		repo := new(mocks.MockFilmRepository)
		updater := service.NewFilmUpdater(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewPutFilm(finder, updater)

		repo.On("Find", mock.Anything, mock.Anything).Return(testFilm, nil).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("return 404 if film does not exist", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/films/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.New().String())

		repo := new(mocks.MockFilmRepository)
		updater := service.NewFilmUpdater(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewPutFilm(finder, updater)

		repo.On("Find", mock.Anything, mock.Anything).Return(film.Film{}, errors.New("not found")).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("return status 401 if user is not the film creator", func(t *testing.T) {
		e := echo.New()
		testFilm := factories.Film()
		req := httptest.NewRequest(http.MethodPut, "/films/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.New().String())
		ctx.Set("userID", uuid.New())

		repo := new(mocks.MockFilmRepository)
		updater := service.NewFilmUpdater(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewPutFilm(finder, updater)

		repo.On("Find", mock.Anything, mock.Anything).Return(testFilm, nil).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("return 500 if film cannot be updated", func(t *testing.T) {
		e := echo.New()
		testFilm := factories.Film()
		body, marshalErr := json.Marshal(PutFilmRequest{
			Title:       testFilm.Title,
			Director:    testFilm.Director,
			ReleaseDate: testFilm.ReleaseDate,
			Genre:       testFilm.Genre,
			Synopsis:    testFilm.Synopsis,
			Casting:     testFilm.Casting,
		})
		assert.NoError(t, marshalErr)

		req := httptest.NewRequest(http.MethodPut, "/films/:id", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.New().String())
		ctx.Set("userID", testFilm.CreatedBy)

		repo := new(mocks.MockFilmRepository)
		updater := service.NewFilmUpdater(repo)
		finder := service.NewFilmFinder(repo)
		handler := NewPutFilm(finder, updater)

		repo.On("Find", mock.Anything, mock.Anything).Return(testFilm, nil).Once()
		repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("fail to update")).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		repo.AssertExpectations(t)
	})
}
