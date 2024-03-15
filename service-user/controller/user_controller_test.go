package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"service-user/model"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyUsers = []model.User{
	{
		Id:       "C001",
		Email:    "testing1@example.com",
		Password: "password1",
	},
	{
		Id:       "C002",
		Email:    "testing2@example.com",
		Password: "password2",
	},
}

type UserUsecaseMock struct {
	mock.Mock
}

type UserControllerTestSuite struct {
	suite.Suite
	routerMock *fiber.App
	useCaseMock *UserUsecaseMock
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.routerMock = fiber.New()
	suite.useCaseMock = new(UserUsecaseMock)
}


func (c *UserUsecaseMock) Login(user *model.User, ctx *fiber.Ctx) (*model.User, error) {
	args := c.Called(user, ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (c *UserUsecaseMock) Register(user *model.User, ctx *fiber.Ctx) error {
	args := c.Called(user, ctx)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

type Response struct {
	Code int 
	Status string
	AccessToken string
	Data interface{}
}

func (suite *UserControllerTestSuite) TestLogin_FailedCase() {
	requestBody := model.User{
		Email:    "testing1@example.com",
		Password: "password1",
	}
	requestBodyJSON, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(requestBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rr, err:= suite.routerMock.Test(req)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusInternalServerError, rr.StatusCode)

	body, err := io.ReadAll(rr.Body)
	assert.Nil(suite.T(), err)
	expectedResponseBody := `{"success":false,"errorMessage":"login failed"}`
	assert.JSONEq(suite.T(), expectedResponseBody, string(body))
}

func (suite *UserControllerTestSuite) TestLogin_SuccessCase() {
	dummyUser := dummyUsers[0]
	app := fiber.New()
	// Stubbing the Login method of UserUsecaseMock to return nil error
	suite.useCaseMock.On("Login", mock.Anything, mock.Anything).Return(&dummyUser, nil)

	requestBody := model.User{
		Email:    "testing1@example.com",
		Password: "password1",
	}
	requestBodyJSON, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(requestBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *UserControllerTestSuite) TestRegisterUserApi_FailedBinding() {
	resp, err := suite.routerMock.Test(httptest.NewRequest(http.MethodPost, "/user/register", nil))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err)

	var errorResponse struct{Err string}
	err = json.Unmarshal(body, &errorResponse)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), errorResponse.Err)
}

func (suite *UserControllerTestSuite) TestRegisterUserApi_SuccessCase() {
	dummyUser := dummyUsers[0]

	reqBody, _ := json.Marshal(dummyUser)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(reqBody))

	resp, err := suite.routerMock.Test(req)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusOK, resp.StatusCode)
	suite.useCaseMock.AssertCalled(suite.T(), "Register", &dummyUser, mock.Anything)
}

func (suite *UserControllerTestSuite) TestRegisterUserApi_FailedUsecase() {
	dummyUser := dummyUsers[0]
	suite.useCaseMock.On("Register", &dummyUser, mock.Anything).Return(errors.New("failed"))

	reqBody, _ := json.Marshal(dummyUser)
	req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(reqBody))
	resp, err := suite.routerMock.Test(req)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err)

	var errorResponse struct{Err string}
	err = json.Unmarshal(body, &errorResponse)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "failed", errorResponse.Err)
}

func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
