package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/internal/service"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestRegisterUser(t *testing.T) {
	// Setup
	e := echo.New()

	mockUser := factories.User()
	body, marshalErr := json.Marshal(RegisterUserRequest{
		Username: mockUser.Username,
		Password: mockUser.Password,
	})
	assert.NoError(t, marshalErr)

	req := httptest.NewRequest(echo.POST, "/register", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	repo := new(mocks.MockUserRepository)
	hasher := new(mocks.MockPasswordHasher)
	srv := service.NewUserRegister(repo, hasher)
	handler := NewRegisterUser(srv)

	hasher.On("Hash", mockUser.Password).Return(mockUser.Password, nil).Once()
	repo.On("Save", mock.Anything, mock.AnythingOfType("user.User")).Return(user.User{}, nil).Once()

	// Assert
	err := handler.Handle(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	repo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}
