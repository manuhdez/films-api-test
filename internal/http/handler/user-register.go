package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/internal/infra"
	"github.com/manuhdez/films-api-test/internal/service"
)

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}

type RegisterUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type RegisterUser struct {
	registerService service.UserRegister
	validator       *validator.Validate
}

func NewRegisterUser(registerService service.UserRegister) RegisterUser {
	return RegisterUser{
		registerService: registerService,
		validator:       validator.New(),
	}
}

func (h RegisterUser) Handle(c echo.Context) error {
	var req RegisterUserRequest
	if bindErr := c.Bind(&req); bindErr != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(bindErr))
	}

	validationErr := h.validator.StructCtx(c.Request().Context(), req)
	if validationErr != nil {
		fmt.Println("validation error", validationErr)
		return c.JSON(http.StatusBadRequest, NewErrorResponse(validationErr))
	}

	newUser, userErr := user.Create(req.Username, req.Password)
	if userErr != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, NewErrorResponse(userErr))
	}

	ctx := c.Request().Context()
	usr, registerErr := h.registerService.Register(ctx, newUser)
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
