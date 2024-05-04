package middleware

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/http/handler"
)

func LoggedUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ID, err := getUserIDFromContext(c)
		if err != nil {
			return handler.ErrUnauthorized
		}

		uid, parseErr := uuid.Parse(ID)
		if parseErr != nil {
			return handler.ErrInvalidID
		}

		c.Set("userID", uid)
		return next(c)
	}
}

func getUserIDFromContext(c echo.Context) (string, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user"].(string)

	if userID == "" {
		return "", errors.New("missing user token")
	}

	return userID, nil
}
