package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		panic("not APP_PORT has been defined")
	}

	e.Logger.Fatal(e.Start(":" + port))
}
