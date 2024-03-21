package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"service-employee/model"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type EmployeeUsecase interface {
	CreateEmployee(user *model.Employee, ctx *MockContext) error
}

type EmployeeUsecaseMock struct {
	mock.Mock
}

type MockContext struct {
	App      *fiber.App
	Response *http.Response
}

func (c *EmployeeUsecaseMock) CreateEmployee(user *model.Employee, ctx *MockContext) error {
	args := c.Called(user, ctx)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func TestCreateEmployee(t *testing.T) {
	app := fiber.New()
	mockUsecase := &EmployeeUsecaseMock{}
	requestBody := model.Employee{
		Id:   "C001",
		Name: "John Doe",
	}
	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/employee", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&http.Cookie{Name: "access_token", Value: "mock_access_token"})

	userAuthHandler := func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	}
	app.Get("/user/auth", userAuthHandler)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ctx := &MockContext{
		App:      app,
		Response: resp,
	}

	err = mockUsecase.CreateEmployee(&requestBody, ctx)
	assert.NoError(t, err)
}
