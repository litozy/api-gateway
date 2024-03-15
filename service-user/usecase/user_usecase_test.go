package usecase

import (
	"errors"
	"fmt"
	"service-user/helpers"
	"service-user/model"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyUserData = []struct {
	id       string
	email    string
	password string
}{
	{
		id:       "C001",
		email:    "testing1@example.com",
		password: "password1",
	},
	{
		id:       "C002",
		email:    "testing2@example.com",
		password: "password2",
	},
}

type userUsecaseTestSuite struct {
	suite.Suite
	repoMock *repoMock
}

func (s *userUsecaseTestSuite) SetupTest() {
	s.repoMock = new(repoMock)
}

type repoMock struct {
	mock.Mock
}

func (r *repoMock) CreateUser(newUser *model.User) error {
	args := r.Called(newUser)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *repoMock) FindOneByEmail(email string) (*model.User, error) {
	args := r.Called(email)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*model.User); ok {
		return user, nil
	}
	return nil, nil
}

func (s *userUsecaseTestSuite) TestUserRegister_Success() {
	dummyUser := &model.User{
		Id:       dummyUserData[0].id,
		Email:    dummyUserData[0].email,
		Password: dummyUserData[0].password,
	}
	ctx := &fiber.Ctx{}
	s.repoMock.On("FindOneByEmail", dummyUser.Email).Return(nil, nil)
	s.repoMock.On("CreateUser", dummyUser).Return(nil)
	userUsecaseTest := NewUserUsecase(s.repoMock)
	err := userUsecaseTest.Register(dummyUser, ctx)
	assert.Nil(s.T(), err)
}

func (s *userUsecaseTestSuite) TestUserRegister_Failed() {
	dummyUser := &model.User{
		Id:       dummyUserData[0].id,
		Email:    dummyUserData[0].email,
		Password: dummyUserData[0].password,
	}
	ctx := &fiber.Ctx{}
	expectedError := errors.New("failed")
	s.repoMock.On("FindOneByEmail", dummyUser.Email).Return(expectedError, expectedError)
	s.repoMock.On("CreateUser", dummyUser).Return(expectedError)
	userUsecaseTest := NewUserUsecase(s.repoMock)
	err := userUsecaseTest.Register(dummyUser, ctx)
	assert.EqualError(s.T(), err, expectedError.Error())
}

func (suite *userUsecaseTestSuite) TestUserLogin_Success() {
	validUser := &model.User{
		Email:    dummyUserData[0].email,
		Password: dummyUserData[0].password,
	}
	suite.repoMock.On("FindOneByEmail", dummyUserData[0].email).Return(validUser, nil)

	dummyUser := &model.User{
		Email:    dummyUserData[0].email,
		Password: dummyUserData[0].password,
	}
	fmt.Println(validUser)
	fmt.Println(dummyUser)
	ctx := &fiber.Ctx{}
	userUsecase := NewUserUsecase(suite.repoMock)
	loggedInUser, err := userUsecase.Login(dummyUser, ctx)
	fmt.Println(loggedInUser)

	suite.Nil(err)
	suite.NotNil(loggedInUser)
	suite.Equal(dummyUser.Email, loggedInUser.Email)
	suite.True(helpers.ComparePassword([]byte(validUser.Password), []byte(dummyUser.Password)))
}

func (s *userUsecaseTestSuite) TestUserLogin_Failed_EmailNotFound() {
    s.repoMock.On("FindOneByEmail", mock.Anything).Return(nil, nil)

    dummyUser := &model.User{
        Email:    "nonexistent@example.com",
        Password: "password123",
    }
	ctx := &fiber.Ctx{}
    userUsecaseTest := NewUserUsecase(s.repoMock)
    _, err := userUsecaseTest.Login(dummyUser, ctx)

    assert.NotNil(s.T(), err)
    assert.EqualError(s.T(), err, "(401) BAD_REQUEST, email is not a valid ")
}

func (s *userUsecaseTestSuite) TestUserLogin_Failed_PasswordMismatch() {
	validUser := &model.User{
		Email:    dummyUserData[0].email,
		Password: "password1", // Kata sandi yang benar
	}
	s.repoMock.On("FindOneByEmail", dummyUserData[0].email).Return(validUser, nil)

	dummyUser := &model.User{
		Email:    dummyUserData[0].email,
		Password: "wrong_password",
	}
	ctx := &fiber.Ctx{}
	userUsecaseTest := NewUserUsecase(s.repoMock)
	_, err := userUsecaseTest.Login(dummyUser, ctx)

	assert.NotNil(s.T(), err)
	assert.EqualError(s.T(), err, "(401) BAD_REQUEST, password does not match ")
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(userUsecaseTestSuite))
}
