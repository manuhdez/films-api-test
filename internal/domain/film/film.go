package film

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
)

var (
	ErrInvalidFilterValues = errors.New("invalid filter values")
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

func Create(
	title string,
	director string,
	year int,
	genre string,
	synopsis string,
	casting []string,
	createdBy uuid.UUID,
) Film {
	id := uuid.New()
	return New(id, title, director, year, genre, synopsis, casting, createdBy)
}

type Filter struct {
	Title       string
	Director    string
	Genre       string
	ReleaseDate int
}

func NewFilter(title, director, genre, releaseDate string) Filter {
	filter := Filter{
		Title:    title,
		Director: director,
		Genre:    genre,
	}

	year, err := strconv.Atoi(releaseDate)
	if err == nil {
		filter.ReleaseDate = year
	}

	return filter
}
