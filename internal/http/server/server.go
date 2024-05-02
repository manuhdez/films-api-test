package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/http/handler"
	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/internal/service"
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

	userRepo := infra.NewPostgresUserRepository(db)
	passwordHasher := infra.NewBcryptPasswordHasher()
	tokenGenerator := infra.NewJWTGenerator()
	userCreator := service.NewUserRegister(userRepo, passwordHasher)
	userLogger := service.NewUserLogin(userRepo, passwordHasher)

	api := e.Group("/api")
	api.POST("/register", handler.NewRegisterUser(userCreator).Handle)
	api.POST("/login", handler.NewLoginUser(userLogger, tokenGenerator).Handle)

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
