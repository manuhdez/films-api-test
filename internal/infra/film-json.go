package infra

import "github.com/manuhdez/films-api-test/internal/domain/film"

type FilmJSON struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Director    string   `json:"director"`
	ReleaseDate int      `json:"release_date"`
	Genre       string   `json:"genre"`
	Casting     []string `json:"casting"`
	Synopsis    string   `json:"synopsis"`
	CreatedBy   string   `json:"created_by"`
}

func NewFilmJSON(f film.Film) FilmJSON {
	return FilmJSON{
		ID:          f.ID.String(),
		Title:       f.Title,
		Director:    f.Director,
		ReleaseDate: f.ReleaseDate,
		Genre:       f.Genre,
		Casting:     f.Casting,
		Synopsis:    f.Synopsis,
		CreatedBy:   f.CreatedBy.String(),
	}
}
