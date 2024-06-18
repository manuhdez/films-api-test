package infra

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/manuhdez/films-api-test/internal/domain/film"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

type PostgresFilmRepository struct {
	db  *sql.DB
	gdb *gorm.DB
}

func NewPostgresFilmRepository(db *sql.DB) *PostgresFilmRepository {
	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: db}),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &PostgresFilmRepository{db: db, gdb: gormDB}
}

func (r *PostgresFilmRepository) All(c context.Context, filter film.Filter) ([]film.Film, error) {
	query := psql.Select("*").From("films")

	var params []interface{}
	if filter.Title != "" {
		query = query.Where("LOWER(title) = LOWER(?)")
		params = append(params, filter.Title)
	}
	if filter.Director != "" {
		query = query.Where("LOWER(director) = LOWER(?)")
		params = append(params, filter.Director)
	}
	if filter.Genre != "" {
		query = query.Where("LOWER(genre) = LOWER(?)")
		params = append(params, filter.Genre)
	}
	if filter.ReleaseDate != 0 {
		query = query.Where("release_date = ?")
		params = append(params, filter.ReleaseDate)
	}

	sqlQuery, _, sqlErr := query.ToSql()
	if sqlErr != nil {
		slog.Error("failed to build get films query", "error", sqlErr.Error())
		return []film.Film{}, sqlErr
	}

	rows, queryErr := r.db.QueryContext(c, sqlQuery, params...)
	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			slog.Error("no films to return")
			return []film.Film{}, nil
		}

		slog.Error("failed to execute query", "error", queryErr)
		return nil, queryErr
	}

	defer rows.Close()

	var films []film.Film
	for rows.Next() {
		var f PostgresFilm
		scanErr := rows.Scan(
			&f.ID, &f.Title, &f.Director, &f.ReleaseDate, &f.Casting, &f.Genre, &f.Synopsis, &f.CreatedBy,
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

func (r *PostgresFilmRepository) Find(c context.Context, id uuid.UUID) (film.Film, error) {
	query := psql.Select("*").From("films").Where(sq.Eq{"id": id})

	sqlQuery, _, sqlErr := query.ToSql()
	if sqlErr != nil {
		slog.Error("failed to build get films query", "error", sqlErr.Error())
		return film.Film{}, sqlErr
	}

	row := r.db.QueryRowContext(c, sqlQuery, id)
	if row.Err() != nil {
		slog.Error("film not found", "id", id.String())
		return film.Film{}, row.Err()
	}

	var f PostgresFilm
	scanErr := row.Scan(&f.ID, &f.Title, &f.Director, &f.ReleaseDate, &f.Casting, &f.Genre, &f.Synopsis, &f.CreatedBy)
	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			slog.Error("film not found, empty row value", "id", id.String())
		}
		return film.Film{}, scanErr
	}

	return f.ToDomain(), nil
}

func (r *PostgresFilmRepository) Save(c context.Context, f film.Film) error {
	query, _, err := psql.Insert("films").Columns(
		"id", "title", "director", "release_date", "genre", "synopsis", "casting", "created_by",
	).Values("$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8").ToSql()

	if err != nil {
		slog.Error("failed to build sql query", "error", err.Error())
		return err
	}

	_, insertErr := r.db.ExecContext(
		c, query, f.ID, f.Title, f.Director, f.ReleaseDate, f.Genre, f.Synopsis, pq.Array(f.Casting), f.CreatedBy,
	)
	if insertErr != nil {
		slog.Error("failed to insert film", "error", insertErr.Error())
		return insertErr
	}

	return nil
}

func (r *PostgresFilmRepository) Delete(c context.Context, id uuid.UUID) error {
	query, _, queryErr := psql.Delete("films").Where(sq.Eq{"id": id}).ToSql()
	if queryErr != nil {
		slog.Error("failed to build sql query", "error", queryErr.Error())
		return queryErr
	}

	_, err := r.db.ExecContext(c, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresFilmRepository) Update(c context.Context, f film.Film) error {
	query, _, queryErr := psql.Update("films").SetMap(map[string]interface{}{
		"title":        f.Title,
		"director":     f.Director,
		"release_date": f.ReleaseDate,
		"genre":        f.Genre,
		"synopsis":     f.Synopsis,
		"casting":      f.Casting,
	}).Where(sq.Eq{"id": f.ID}).ToSql()

	slog.Info("update film query", "query", query)

	if queryErr != nil {
		slog.Error("failed to build sql query", "error", queryErr.Error())
		return queryErr
	}

	casting := pq.Array(f.Casting)
	_, err := r.db.ExecContext(c,
		query,
		casting,
		f.Director,
		f.Genre,
		f.ReleaseDate,
		f.Synopsis,
		f.Title,
		f.ID)
	if err != nil {
		slog.Error("failed to update film", "film", f.ID.String(), "error", err)
		return err
	}

	return nil
}
