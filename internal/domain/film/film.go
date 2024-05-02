package film

import (
	"github.com/google/uuid"
)

type Film struct {
	ID          uuid.UUID
	Title       string
	Director    string
	ReleaseDate int
	Genre       string
	Synopsis    string
	Casting     []string
	CreatedBy   uuid.UUID
}

func New(
	id uuid.UUID,
	title string,
	director string,
	year int,
	genre string,
	synopsis string,
	casting []string,
	createdBy uuid.UUID,
) Film {
	return Film{
		ID:          id,
		Title:       title,
		Director:    director,
		ReleaseDate: year,
		Genre:       genre,
		Synopsis:    synopsis,
		Casting:     casting,
		CreatedBy:   createdBy,
	}
}
