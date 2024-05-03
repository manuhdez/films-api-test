package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/internal/service"
)

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type RegisterUser struct {
	registerService service.UserRegister
}

func NewRegisterUser(registerService service.UserRegister) RegisterUser {
	return RegisterUser{registerService: registerService}
}

func (h RegisterUser) Handle(c echo.Context) error {
	var req RegisterUserRequest
	if bindErr := c.Bind(&req); bindErr != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(bindErr))
	}

	newUser, userErr := user.Create(req.Username, req.Password)
	if userErr != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, NewErrorResponse(userErr))
	}

	usr, registerErr := h.registerService.Register(context.Background(), newUser)
	if registerErr != nil {
		return handleRegisterError(registerErr)
	}

	return c.JSON(http.StatusCreated, RegisterUserResponse{usr.ID.String(), usr.Username})
}

func handleRegisterError(err error) error {
	if errors.Is(err, infra.ErrUsernameAlreadyInUse) {
		return echo.NewHTTPError(http.StatusBadRequest, NewErrorResponse(err))
	}

	return echo.NewHTTPError(http.StatusInternalServerError, NewErrorResponse(err))
}
