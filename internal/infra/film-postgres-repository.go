package infra

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/manuhdez/films-api-test/internal/domain/film"
)

type PostgresFilmRepository struct {
	db *sql.DB
}

func NewPostgresFilmRepository(db *sql.DB) *PostgresFilmRepository {
	return &PostgresFilmRepository{db: db}
}

func (r *PostgresFilmRepository) All(c context.Context) ([]film.Film, error) {
	query := `SELECT id, title, director, release_date, genre, synopsis, casting, created_by FROM films`

	rows, queryErr := r.db.QueryContext(c, query)
	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			slog.Info("no films found")
			return []film.Film{}, nil
		}
		return nil, queryErr
	}

	defer rows.Close()

	var films []film.Film
	for rows.Next() {
		var f PostgresFilm
		scanErr := rows.Scan(
			&f.ID,
			&f.Title,
			&f.Director,
			&f.ReleaseDate,
			&f.Genre,
			&f.Synopsis,
			&f.Casting,
			&f.CreatedBy,
		)

		if scanErr != nil {
			if errors.Is(scanErr, sql.ErrNoRows) {
				return []film.Film{}, nil
			}

			slog.Error("failed to scan film row", "error", scanErr.Error())
			return nil, scanErr
		}

		films = append(films, f.ToDomain())
	}

	return films, nil
}
