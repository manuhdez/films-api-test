package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type Server struct {
	engine *echo.Echo
	db     *sql.DB
}

func New(db *sql.DB) Server {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	return Server{
		engine: e,
		db:     db,
	}
}

func (s *Server) Start() error {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		panic("not APP_PORT has been defined")
	}

	return s.engine.Start(fmt.Sprintf(":%s", port))
}
