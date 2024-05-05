package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/internal/service"
)

type GetFilms struct {
	filmsGetter service.FilmsGetter
}

func NewGetFilms(filmsGetter service.FilmsGetter) GetFilms {
	return GetFilms{filmsGetter: filmsGetter}
}

func (h GetFilms) Handle(c echo.Context) error {
	filter := film.NewFilter(
		c.QueryParam("title"),
		c.QueryParam("director"),
		c.QueryParam("genre"),
		c.QueryParam("release_date"),
	)

	films, err := h.filmsGetter.Get(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(err))
	}

	jsonFilms := make([]infra.FilmJSON, len(films))
	for idx, f := range films {
		jsonFilms[idx] = infra.NewFilmJSON(f)
	}

	return c.JSON(http.StatusOK, jsonFilms)
}
