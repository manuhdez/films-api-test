package database

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	ErrMissingConnURL   = errors.New("missing connection url")
	ErrConnectionFailed = errors.New("database connection failed")
)

// NewPostgresConnection - connects with the database and returns a connection pool
func NewPostgresConnection() (*sql.DB, error) {
	dbUrl, ok := os.LookupEnv("DB_URL")
	if !ok {
		log.Println("DB_URL environment variable not set")
		return nil, ErrMissingConnURL
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil, ErrConnectionFailed
	}

	return conn, nil
}
