package server

import (
	"database/sql"
	"fmt"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/manuhdez/films-api-test/internal/http/handler"
	middle "github.com/manuhdez/films-api-test/internal/http/middleware"
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
	userFinder := service.NewUserFinder(userRepo)

	filmRepo := infra.NewPostgresFilmRepository(db)
	filmsGetter := service.NewFilmsGetter(filmRepo)
	filmFinder := service.NewFilmFinder(filmRepo)
	filmCreator := service.NewFilmCreator(filmRepo)
	filmDeleter := service.NewFilmDeleter(filmRepo)
	filmUpdater := service.NewFilmUpdater(filmRepo)
	authMiddleware := echojwt.JWT([]byte(os.Getenv("JWT_SECRET_KEY")))

	api := e.Group("/api")

	api.POST("/register", handler.NewRegisterUser(userCreator).Handle)
	api.POST("/login", handler.NewLoginUser(userLogger, tokenGenerator).Handle)

	films := api.Group("/films")
	films.Use(authMiddleware, middle.LoggedUser)
	films.GET("", handler.NewGetFilms(filmsGetter).Handle)
	films.POST("", handler.NewPostFilm(filmCreator).Handle)
	films.GET("/:id", handler.NewFindFilm(filmFinder, userFinder).Handle)
	films.DELETE("/:id", handler.NewDeleteFilm(filmDeleter, filmFinder).Handle)
	films.PUT("/:id", handler.NewPutFilm(filmFinder, filmUpdater).Handle)

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
