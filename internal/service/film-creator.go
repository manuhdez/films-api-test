package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/manuhdez/films-api-test/internal/domain"
	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/internal/infra"
)

type FilmCreatedEvent struct {
	Film film.Film
}

func (FilmCreatedEvent) Key() string {
	return "api.films.created"
}

func (e FilmCreatedEvent) Data() []byte {
	data, err := json.Marshal(infra.NewFilmJSON(e.Film))
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
		return nil
	}

	return data
}

type FilmCreator struct {
	repository film.Repository
	eventBus   domain.EventBus
}

func NewFilmCreator(r film.Repository, b domain.EventBus) FilmCreator {
	return FilmCreator{repository: r, eventBus: b}
}

func (fc FilmCreator) Create(ctx context.Context, f film.Film) error {
	err := fc.repository.Save(ctx, f)
	if err != nil {
		return errors.New("failed to save film")
	}

	event := FilmCreatedEvent{Film: f}
	err = fc.eventBus.Publish(ctx, event)
	if err != nil {
		log.Println("failed to publish film created event")
	}

	return nil
}
