package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/internal/service"
)

var (
	ErrUpdatingFilm = errors.New("failed to update film")
)

type PutFilmRequest struct {
	Title       string   `json:"title"`
	Director    string   `json:"director"`
	ReleaseDate int      `json:"release_date"`
	Genre       string   `json:"genre"`
	Synopsis    string   `json:"synopsis"`
	Casting     []string `json:"casting"`
}

type PutFilm struct {
	finder  service.FilmFinder
	updater service.FilmUpdater
}

func NewPutFilm(finder service.FilmFinder, updater service.FilmUpdater) PutFilm {
	return PutFilm{finder: finder, updater: updater}
}

func (h PutFilm) Handle(c echo.Context) error {
	filmID, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(ErrInvalidID))
	}

	f, findErr := h.finder.Find(filmID)
	if findErr != nil {
		return c.JSON(http.StatusNotFound, NewErrorResponse(ErrFilmNotFound))
	}

	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok || userID != f.CreatedBy {
		return c.JSON(http.StatusUnauthorized, NewErrorResponse(ErrUnauthorized))
	}

	// TODO: validate request values
	var req PutFilmRequest
	bindErr := c.Bind(&req)
	if bindErr != nil || validateRequest(req) != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.New("missing or invalid request body")))
	}

	updated := film.New(
		f.ID,
		req.Title,
		req.Director,
		req.ReleaseDate,
		req.Genre,
		req.Synopsis,
		req.Casting,
		f.CreatedBy,
	)

	err := h.updater.Update(updated)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(ErrUpdatingFilm))
	}

	return c.JSON(http.StatusOK, infra.NewFilmJSON(updated))
}

// validateRequest checks that all required fields are sent on the request body
func validateRequest(req PutFilmRequest) error {
	if req.Title == "" || req.Director == "" || req.Genre == "" || req.Synopsis == "" || req.ReleaseDate == 0 || req.Casting == nil {
		return errors.New("missing required fields")
	}
	return nil
}
