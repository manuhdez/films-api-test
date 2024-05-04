package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/film"
)

type MockFilmRepository struct {
	mock.Mock
}

func (m *MockFilmRepository) All(c context.Context) ([]film.Film, error) {
	args := m.Called(c)
	return args.Get(0).([]film.Film), args.Error(1)
}

func (m *MockFilmRepository) Find(c context.Context, id uuid.UUID) (film.Film, error) {
	args := m.Called(c, id)
	return args.Get(0).(film.Film), args.Error(1)
}

func (m *MockFilmRepository) Save(c context.Context, f film.Film) error {
	args := m.Called(c, f)
	return args.Error(0)
}

func (m *MockFilmRepository) Delete(c context.Context, id uuid.UUID) error {
	args := m.Called(c, id)
	return args.Error(0)
}
