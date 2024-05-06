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
	"github.com/manuhdez/films-api-test/internal/service"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestFindFilm_Handle(t *testing.T) {
	t.Parallel()

	t.Run("returns a 200 code if the film is found", func(t *testing.T) {
		testFilm := factories.Film()
		testUser := factories.User()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/films/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(testFilm.ID.String())

		filmRepo := new(mocks.MockFilmRepository)
		userRepo := new(mocks.MockUserRepository)
		filmFinder := service.NewFilmFinder(filmRepo)
		userFinder := service.NewUserFinder(userRepo)
		handler := NewFindFilm(filmFinder, userFinder)

		filmRepo.On("Find", mock.Anything, mock.Anything).Return(testFilm, nil).Once()
		userRepo.On("Find", mock.Anything, mock.Anything).Return(testUser, nil).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var responseFilm FindFilmResponse
		err = json.Unmarshal(rec.Body.Bytes(), &responseFilm)
		assert.NoError(t, err)
		assert.Equal(t, testFilm.ID.String(), responseFilm.ID)
		assert.Equal(t, testUser.ID.String(), responseFilm.CreatedBy.ID)
	})

	t.Run("returns 422 if uuid is not valid", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/films/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("bad-uuid")

		filmRepo := new(mocks.MockFilmRepository)
		userRepo := new(mocks.MockUserRepository)
		filmFinder := service.NewFilmFinder(filmRepo)
		userFinder := service.NewUserFinder(userRepo)
		handler := NewFindFilm(filmFinder, userFinder)

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("returns a 404 if cannot find film", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/films/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.New().String())

		filmRepo := new(mocks.MockFilmRepository)
		userRepo := new(mocks.MockUserRepository)
		filmFinder := service.NewFilmFinder(filmRepo)
		userFinder := service.NewUserFinder(userRepo)
		handler := NewFindFilm(filmFinder, userFinder)

		testErr := errors.New("cannot find film")
		filmRepo.On("Find", mock.Anything, mock.Anything).Return(film.Film{}, testErr).Once()

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		filmRepo.AssertExpectations(t)
	})
}
