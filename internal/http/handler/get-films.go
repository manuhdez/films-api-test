package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

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
	films, err := h.filmsGetter.Get()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(err))
	}

	var jsonFilms []infra.FilmJSON
	for _, film := range films {
		jsonFilms = append(jsonFilms, infra.NewFilmJSON(film))
	}

	return c.JSON(http.StatusOK, jsonFilms)
}
