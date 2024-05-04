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

	"github.com/manuhdez/films-api-test/internal/service"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

func TestLoginUser(t *testing.T) {
	t.Parallel()

	testUser := factories.User()

	body, marshalErr := json.Marshal(UserLoginRequest{
		Username: testUser.Username,
		Password: testUser.Password,
	})

	assert.NoError(t, marshalErr)

	t.Run("with valid user credentials", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		repo := new(mocks.MockUserRepository)
		hasher := new(mocks.MockPasswordHasher)
		srv := service.NewUserLogin(repo, hasher)
		tokenSrv := new(mocks.MockTokenGenerator)
		handler := NewLoginUser(srv, tokenSrv)

		repo.On("SearchByUsername", mock.Anything, testUser.Username).Return(testUser, nil)
		hasher.On("Compare", testUser.Password, testUser.Password).Return(true)
		tokenSrv.On("Generate", testUser.ID.String()).Return("mock-token", nil)

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
		tokenSrv.AssertExpectations(t)
	})

	t.Run("with invalid user credentials", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		repo := new(mocks.MockUserRepository)
		hasher := new(mocks.MockPasswordHasher)
		srv := service.NewUserLogin(repo, hasher)
		tokenSrv := new(mocks.MockTokenGenerator)
		handler := NewLoginUser(srv, tokenSrv)

		repo.On("SearchByUsername", mock.Anything, testUser.Username).Return(testUser, nil)
		hasher.On("Compare", testUser.Password, testUser.Password).Return(false)

		err := handler.Handle(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
		tokenSrv.AssertExpectations(t)
	})
}
