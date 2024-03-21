// File: controller/user_controller_test.go

package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"service-user/model"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock dari UserUsecase
type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(user *model.User, c *fiber.Ctx) error {
	args := m.Called(user, c)
	return args.Error(0)
}

func (m *MockUserUsecase) Login(user *model.User, c *fiber.Ctx) (*model.User, error) {
	args := m.Called(user, c)
	return args.Get(0).(*model.User), args.Error(1)
}

func NewMockUserUsecase() *MockUserUsecase {
	return &MockUserUsecase{}
}

func TestRegister(t *testing.T) {
    mockUsecase := NewMockUserUsecase()
    controller := NewUserController(mockUsecase)

    reqBody, _ := json.Marshal(map[string]string{
        "email":    "newuser@example.com",
        "password": "password123",
    })
    req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()

    mockUsecase.On("Register", mock.Anything, mock.Anything).Return(nil)

    app := fiber.New()
    app.Post("/user/register", controller.Register)

    app.Use(func(c *fiber.Ctx) error {
        c.Locals("csrf", "token")
        return c.Next()
    })

    app.Test(req, -1)
    assert.Equal(t, http.StatusOK, resp.Code)
    mockUsecase.AssertCalled(t, "Register", mock.Anything, mock.Anything)
}



func TestLogin(t *testing.T) {
    mockUsecase := NewMockUserUsecase()
    controller := NewUserController(mockUsecase)

    reqBody, _ := json.Marshal(map[string]string{
        "email":    "user@example.com",
        "password": "password123",
    })
    req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()

    mockUsecase.On("Login", mock.Anything, mock.Anything).Return(&model.User{
        Email: "user@example.com",
    }, nil)

    app := fiber.New()
    app.Post("/user/login", controller.Login)

    app.Use(func(c *fiber.Ctx) error {
        c.Locals("csrf", "token")
        return c.Next()
    })

    app.Test(req, -1) 
    assert.Equal(t, http.StatusOK, resp.Code)
}


