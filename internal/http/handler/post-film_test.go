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

	"github.com/manuhdez/films-api-test/internal/service"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestPostFilm_Handle(t *testing.T) {
	t.Parallel()

	t.Run("returns a 201 when film is created", func(t *testing.T) {
		e := echo.New()
		testFilm := factories.Film()

		body, marshalErr := json.Marshal(PostFilmRequest{
			Title:       testFilm.Title,
			Director:    testFilm.Director,
			ReleaseDate: testFilm.ReleaseDate,
			Genre:       testFilm.Genre,
			Casting:     testFilm.Casting,
			Synopsis:    testFilm.Synopsis,
		})
		assert.NoError(t, marshalErr)

		req := httptest.NewRequest(http.MethodPost, "/films", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userID", uuid.New())

		repo := new(mocks.MockFilmRepository)
		bus := new(mocks.MockEventBus)
		filmCreator := service.NewFilmCreator(repo, bus)
		handler := NewPostFilm(filmCreator)

		repo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()
		bus.On("Publish", mock.Anything, mock.Anything).Return(nil).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("returns a 401 when user is not authorized", func(t *testing.T) {
		e := echo.New()
		testFilm := factories.Film()

		body, marshalErr := json.Marshal(PostFilmRequest{
			Title:       testFilm.Title,
			Director:    testFilm.Director,
			ReleaseDate: testFilm.ReleaseDate,
			Genre:       testFilm.Genre,
			Casting:     testFilm.Casting,
			Synopsis:    testFilm.Synopsis,
		})
		assert.NoError(t, marshalErr)

		req := httptest.NewRequest(http.MethodPost, "/films", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		repo := new(mocks.MockFilmRepository)
		bus := new(mocks.MockEventBus)
		filmCreator := service.NewFilmCreator(repo, bus)
		handler := NewPostFilm(filmCreator)

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("returns a 500 when film cannot be saved", func(t *testing.T) {
		e := echo.New()
		testFilm := factories.Film()

		body, marshalErr := json.Marshal(PostFilmRequest{
			Title:       testFilm.Title,
			Director:    testFilm.Director,
			ReleaseDate: testFilm.ReleaseDate,
			Genre:       testFilm.Genre,
			Casting:     testFilm.Casting,
			Synopsis:    testFilm.Synopsis,
		})
		assert.NoError(t, marshalErr)

		req := httptest.NewRequest(http.MethodPost, "/films", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userID", uuid.New())

		repo := new(mocks.MockFilmRepository)
		bus := new(mocks.MockEventBus)
		filmCreator := service.NewFilmCreator(repo, bus)
		handler := NewPostFilm(filmCreator)

		testErr := errors.New("something went wrong")
		repo.On("Save", mock.Anything, mock.Anything).Return(testErr).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		repo.AssertExpectations(t)
	})
}
