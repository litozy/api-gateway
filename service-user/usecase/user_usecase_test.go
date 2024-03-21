// File: usecase/user_usecase_test.go

package usecase

import (
	"service-user/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock dari UserRepository
type MockUserRepository struct {
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	// Implementasi mock CreateUser
	return nil
}

func (m *MockUserRepository) FindOneByEmail(email string) (*model.User, error) {
	// Implementasi mock FindOneByEmail
	if email == "existinguser@example.com" {
		return &model.User{
			Email:    "existinguser@example.com",
			Password: "hashedpassword123",
		}, nil
	}
	return nil, nil
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func TestRegister(t *testing.T) {
	// Persiapan
	userRepo := NewMockUserRepository()
	usecase := NewUserUsecase(userRepo)

	// Pengujian
	user := &model.User{
		Email:    "newuser@example.com",
		Password: "password123",
	}
	err := usecase.Register(user, nil)

	// Pengecekan
	assert.NoError(t, err)
}

func TestLogin(t *testing.T) {
	userRepo := NewMockUserRepository()
	usecase := NewUserUsecase(userRepo)

	user := &model.User{
		Email:    "existinguser@example.com",
		Password: "password123",
	}
	usecase.Login(user, nil)
	
	user = &model.User{
		Email:    "", 
		Password: "password123",
	}
	_, err := usecase.Login(user, nil)
	assert.EqualError(t, err, "(401) BAD_REQUEST, email is not a valid ")

	user = &model.User{
		Email:    "existinguser@example.com",
		Password: "wrongpassword",
	}
	_, err = usecase.Login(user, nil)
	assert.EqualError(t, err, "(401) BAD_REQUEST, password does not match ")
}

