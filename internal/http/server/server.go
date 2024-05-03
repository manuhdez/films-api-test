package server

import (
	"database/sql"
	"fmt"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	userRepo := infra.NewPostgresUserRepository(db)
	passwordHasher := infra.NewBcryptPasswordHasher()
	tokenGenerator := infra.NewJWTGenerator()
	userCreator := service.NewUserRegister(userRepo, passwordHasher)
	userLogger := service.NewUserLogin(userRepo, passwordHasher)

	api := e.Group("/api")
	api.POST("/register", handler.NewRegisterUser(userCreator).Handle)
	api.POST("/login", handler.NewLoginUser(userLogger, tokenGenerator).Handle)

	filmRepo := infra.NewPostgresFilmRepository(db)
	filmsGetter := service.NewFilmsGetter(filmRepo)
	authMiddleware := echojwt.JWT([]byte(os.Getenv("JWT_SECRET_KEY")))

	films := api.Group("/films")
	films.GET("", handler.NewGetFilms(filmsGetter).Handle, authMiddleware)

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
