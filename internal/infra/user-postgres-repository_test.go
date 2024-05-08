package infra

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/manuhdez/films-api-test/test/containers"
	"github.com/manuhdez/films-api-test/test/factories"
)

func TestPostgresUserRepository(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := containers.CreatePostgresContainer(ctx)
	require.NoError(t, err, "should create postgres container")

	db, err := sql.Open("postgres", pgContainer.ConnectionString)
	require.NoError(t, err, "should open postgres connection")

	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS("../../database/migrations"))
	require.NoError(t, err, "should setup database")

	_, err = provider.Up(ctx)
	require.NoError(t, err, "should perform database migrations")

	repo := NewPostgresUserRepository(db)
	testUser := factories.User()

	t.Run("should create user", func(t *testing.T) {
		_, saveErr := repo.Save(context.Background(), testUser)
		assert.NoError(t, saveErr, "should save user")
	})

	t.Run("should find user by username", func(t *testing.T) {
		u, findErr := repo.SearchByUsername(context.Background(), testUser.Username)
		assert.NoError(t, findErr, "should find user by username")
		assert.Equal(t, testUser, u)
	})

	t.Run("should find user by id", func(t *testing.T) {
		found, findErr := repo.Find(context.Background(), testUser.ID)
		assert.NoError(t, findErr, "should find user by id")
		assert.Equal(t, testUser.Username, found.Username)
	})
}
