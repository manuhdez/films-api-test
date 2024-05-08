package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/service"
)

var (
	ErrInvalidID    = errors.New("invalid uuid")
	ErrFilmNotFound = errors.New("film not found")
)

type FindFilmResponse struct {
	ID          string                  `json:"id"`
	Title       string                  `json:"title"`
	Director    string                  `json:"director"`
	ReleaseDate int                     `json:"release_date"`
	Genre       string                  `json:"genre"`
	Casting     []string                `json:"casting"`
	Synopsis    string                  `json:"synopsis"`
	CreatedBy   FindFilmResponseCreator `json:"created_by"`
}

type FindFilmResponseCreator struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type FindFilm struct {
	filmFinder service.FilmFinder
	userFinder service.UserFinder
}

func NewFindFilm(filmFinder service.FilmFinder, userFinder service.UserFinder) FindFilm {
	return FindFilm{
		filmFinder: filmFinder,
		userFinder: userFinder,
	}
}

func (h FindFilm) Handle(c echo.Context) error {
	id, idErr := uuid.Parse(c.Param("id"))
	if idErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(ErrInvalidID))
	}

	ctx := c.Request().Context()
	film, findErr := h.filmFinder.Find(ctx, id)
	if findErr != nil {
		return c.JSON(http.StatusNotFound, NewErrorResponse(ErrFilmNotFound))
	}

	user, userErr := h.userFinder.Find(ctx, film.CreatedBy)
	if userErr != nil {
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(ErrFilmNotFound))
	}

	return c.JSON(http.StatusOK, FindFilmResponse{
		ID:          film.ID.String(),
		Title:       film.Title,
		Director:    film.Director,
		ReleaseDate: film.ReleaseDate,
		Genre:       film.Genre,
		Synopsis:    film.Synopsis,
		Casting:     film.Casting,
		CreatedBy: FindFilmResponseCreator{
			ID:       user.ID.String(),
			Username: user.Username,
		},
	})
}
