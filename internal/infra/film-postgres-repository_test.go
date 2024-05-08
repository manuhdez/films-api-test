package infra

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/manuhdez/films-api-test/internal/domain/film"
	"github.com/manuhdez/films-api-test/test/containers"
	"github.com/manuhdez/films-api-test/test/factories"
)

func TestPostgresFilmRepository_All(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := containers.CreatePostgresContainer(ctx)
	assert.NoError(t, err, "should create postgres container")

	db, err := sql.Open("postgres", pgContainer.ConnectionString)
	assert.NoError(t, err, "failed to connect to postgres")

	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS("../../database/migrations"))
	assert.NoError(t, err, "failed to init goose provider")

	_, err = provider.Up(ctx)
	assert.NoError(t, err, "failed to run migrations")

	repo := NewPostgresFilmRepository(db)
	userRepo := NewPostgresUserRepository(db)

	testUser := factories.User()
	testFilm := factories.Film()
	testFilm.CreatedBy = testUser.ID

	t.Run("should create a film", func(t *testing.T) {
		_, saveUserErr := userRepo.Save(context.Background(), testUser)
		require.NoError(t, saveUserErr, "should not fail saving testUser")

		saveErr := repo.Save(ctx, testFilm)
		require.NoError(t, saveErr, "should save the film without error")
	})

	t.Run("should get all stored films", func(t *testing.T) {
		films, allErr := repo.All(ctx, film.Filter{})
		assert.NoError(t, allErr, "should get all films")
		assert.Equal(t, len(films), 1)
	})

	t.Run("should find film by id", func(t *testing.T) {
		found, foundErr := repo.Find(ctx, testFilm.ID)
		assert.NoError(t, foundErr, "should find film by id")
		assert.Equal(t, testFilm, found)
	})

	t.Run("should update film", func(t *testing.T) {
		newFilm := factories.Film()
		newFilm.CreatedBy = testUser.ID
		saveErr := repo.Save(ctx, newFilm)
		require.NoError(t, saveErr, "should save the film without error")

		updatedFilm := newFilm
		updatedFilm.Title = "Buscando a Nemo"
		updateErr := repo.Update(ctx, updatedFilm)
		require.NoError(t, updateErr, "should update the film without error")

		found, foundErr := repo.Find(ctx, newFilm.ID)
		assert.NoError(t, foundErr, "should find film by id")
		assert.Equal(t, updatedFilm.Title, found.Title)
	})

	t.Run("should delete film by id", func(t *testing.T) {
		newFilm := factories.Film()
		newFilm.CreatedBy = testUser.ID
		saveErr := repo.Save(ctx, newFilm)
		require.NoError(t, saveErr, "should save the film without error")

		deleteErr := repo.Delete(ctx, newFilm.ID)
		assert.NoError(t, deleteErr, "should delete the film without error")

		_, foundErr := repo.Find(ctx, newFilm.ID)
		assert.ErrorIs(t, foundErr, sql.ErrNoRows)
	})
}
