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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/manuhdez/films-api-test/internal/domain/user"
	"github.com/manuhdez/films-api-test/internal/service"
	"github.com/manuhdez/films-api-test/test/factories"
	"github.com/manuhdez/films-api-test/test/mocks"
)

type RegisterUserSuite struct {
	suite.Suite
	recorder *httptest.ResponseRecorder
	userRepo *mocks.MockUserRepository
	hasher   *mocks.MockPasswordHasher
	eventBus *mocks.MockEventBus
	handler  RegisterUser
}

func (suite *RegisterUserSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.hasher = new(mocks.MockPasswordHasher)
	suite.userRepo = new(mocks.MockUserRepository)
	suite.eventBus = new(mocks.MockEventBus)
	suite.handler = NewRegisterUser(
		service.NewUserRegister(suite.userRepo, suite.hasher, suite.eventBus),
	)
}

func (suite *RegisterUserSuite) TestRegisterValidUser() {
	mockUser, userErr := user.Create("m0cku5ern4me", "validPassword")
	require.NoError(suite.T(), userErr)

	body, marshalErr := json.Marshal(RegisterUserRequest{
		Username: mockUser.Username,
		Password: mockUser.Password,
	})
	require.NoError(suite.T(), marshalErr)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/register", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, suite.recorder)

	suite.hasher.On("Hash", mockUser.Password).Return(mockUser.Password, nil).Once()
	suite.userRepo.On("Save", mock.Anything, mock.AnythingOfType("user.User")).Return(user.User{}, nil).Once()
	suite.eventBus.On("Publish", mock.Anything, mock.Anything).Return(nil).Once()
	err := suite.handler.Handle(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, suite.recorder.Code)
	suite.userRepo.AssertExpectations(suite.T())
	suite.hasher.AssertExpectations(suite.T())
}

func (suite *RegisterUserSuite) TestRegisterInvalidUser() {
	mockUser := factories.User()
	body, marshalErr := json.Marshal(RegisterUserRequest{
		Username: "m0cku$ername",
		Password: mockUser.Password,
	})
	require.NoError(suite.T(), marshalErr)

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/register", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(req, suite.recorder)

	err := suite.handler.Handle(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, suite.recorder.Code)
}

func TestRegisterUser(t *testing.T) {
	suite.Run(t, new(RegisterUserSuite))
}
