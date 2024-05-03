package handler

import (
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
	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/internal/service"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestFindFilm_Handle(t *testing.T) {
	t.Run("returns a 200 code if the film is found", func(t *testing.T) {
		testFilm := factories.Film()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/films/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(testFilm.ID.String())

		repo := new(mocks.MockFilmRepository)
		srv := service.NewFilmFinder(repo)
		handler := NewFindFilm(srv)

		repo.On("Find", mock.Anything, mock.Anything).Return(testFilm, nil).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var responseFilm infra.FilmJSON
		err = json.Unmarshal(rec.Body.Bytes(), &responseFilm)
		assert.NoError(t, err)
		assert.Equal(t, testFilm, responseFilm.ToDomain())
	})

	t.Run("returns 422 if uuid is not valid", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/films/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("bad-uuid")

		repo := new(mocks.MockFilmRepository)
		srv := service.NewFilmFinder(repo)
		handler := NewFindFilm(srv)

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		repo.AssertNotCalled(t, "Find", mock.Anything, mock.Anything)
	})

	t.Run("returns a 404 if cannot find film", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/films/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.New().String())

		repo := new(mocks.MockFilmRepository)
		srv := service.NewFilmFinder(repo)
		handler := NewFindFilm(srv)
		testErr := errors.New("cannot find film")
		repo.On("Find", mock.Anything, mock.Anything).Return(film.Film{}, testErr).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		repo.AssertExpectations(t)
	})
}
