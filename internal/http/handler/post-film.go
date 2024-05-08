package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/internal/service"
)

type PostFilmRequest struct {
	Title       string   `json:"title"`
	Director    string   `json:"director"`
	ReleaseDate int      `json:"release_date"`
	Genre       string   `json:"genre"`
	Synopsis    string   `json:"synopsis"`
	Casting     []string `json:"casting"`
}

type PostFilm struct {
	creator service.FilmCreator
}

func NewPostFilm(creator service.FilmCreator) PostFilm {
	return PostFilm{creator: creator}
}

func (h PostFilm) Handle(c echo.Context) error {
	var req PostFilmRequest
	bindErr := c.Bind(&req)
	if bindErr != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(bindErr))
	}

	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, NewErrorResponse(ErrUnauthorized))
	}

	ctx := c.Request().Context()
	f := film.Create(req.Title, req.Director, req.ReleaseDate, req.Genre, req.Synopsis, req.Casting, userID)

	createErr := h.creator.Create(ctx, f)
	if createErr != nil {
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(createErr))
	}

	return c.JSON(http.StatusCreated, infra.NewFilmJSON(f))
}
