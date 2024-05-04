package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/service"
)

var (
	ErrCannotDeleteFilm = errors.New("cannot delete film")
)

type DeleteFilm struct {
	deleter service.FilmDeleter
	finder  service.FilmFinder
}

func NewDeleteFilm(d service.FilmDeleter, f service.FilmFinder) DeleteFilm {
	return DeleteFilm{deleter: d, finder: f}
}

func (df DeleteFilm) Handle(c echo.Context) error {
	filmID, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(ErrInvalidID))
	}

	film, findErr := df.finder.Find(filmID)
	if findErr != nil {
		return c.JSON(http.StatusNotFound, NewErrorResponse(ErrFilmNotFound))
	}

	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok || userID != film.CreatedBy {
		return c.JSON(http.StatusUnauthorized, NewErrorResponse(ErrUnauthorized))
	}

	deleteErr := df.deleter.Delete(filmID)
	if deleteErr != nil {
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(ErrCannotDeleteFilm))
	}

	return c.JSON(http.StatusOK, nil)
}
