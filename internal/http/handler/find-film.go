package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/internal/service"
)

var (
	ErrInvalidID    = errors.New("invalid uuid")
	ErrFilmNotFound = errors.New("film not found")
)

type FindFilm struct {
	finder service.FilmFinder
}

func NewFindFilm(filmFinder service.FilmFinder) FindFilm {
	return FindFilm{finder: filmFinder}
}

func (h FindFilm) Handle(c echo.Context) error {
	fmt.Println("id received", c.Param("id"))
	id, idErr := uuid.Parse(c.Param("id"))
	if idErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(ErrInvalidID))
	}

	film, findErr := h.finder.Find(id)
	if findErr != nil {
		return c.JSON(http.StatusNotFound, NewErrorResponse(ErrFilmNotFound))
	}

	return c.JSON(http.StatusOK, infra.NewFilmJSON(film))
}
