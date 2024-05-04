package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/internal/service"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestGetFilms(t *testing.T) {
	t.Parallel()

	t.Run("returns a list of films", func(t *testing.T) {
		e := echo.New()
		testFilms := factories.FilmList(5)
		req := httptest.NewRequest(http.MethodGet, "/films", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		repo := new(mocks.MockFilmRepository)
		srv := service.NewFilmsGetter(repo)
		handler := NewGetFilms(srv)

		repo.On("All", mock.Anything).Return(testFilms, nil).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var responseFilms []infra.FilmJSON
		err = json.Unmarshal(rec.Body.Bytes(), &responseFilms)
		assert.NoError(t, err)
		assert.Equal(t, len(testFilms), len(responseFilms))
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/films", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		repo := new(mocks.MockFilmRepository)
		srv := service.NewFilmsGetter(repo)
		handler := NewGetFilms(srv)

		repo.On("All", mock.Anything).Return([]film.Film{}, errors.New("something failed")).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
