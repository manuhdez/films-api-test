package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/service"
)

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type LoginUser struct {
	loginService service.UserLogin
	tokenService service.TokenGenerator
}

func NewLoginUser(loginService service.UserLogin, tokenService service.TokenGenerator) LoginUser {
	return LoginUser{
		loginService: loginService,
		tokenService: tokenService,
	}
}

func (h LoginUser) Handle(c echo.Context) error {
	var req UserLoginRequest
	if bindErr := c.Bind(&req); bindErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": bindErr.Error()})
	}

	user, err := h.loginService.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrWrongCredentials) {
			return c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(err))
		}

		return c.JSON(http.StatusInternalServerError, NewErrorResponse(err))
	}

	token, tokenErr := h.tokenService.Generate(user.ID.String())
	if tokenErr != nil {
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(tokenErr))
	}

	return c.JSON(http.StatusOK, LoginResponse{
		UserID: user.ID.String(),
		Token:  token,
	})
}
