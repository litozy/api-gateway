package usecase

import (
	"errors"
	"fmt"
	"service-user/helpers"
	"service-user/model"
	"service-user/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WebResponse struct {
	Code int
	Status string
	Data interface{}
}

type UserUsecase interface {
	Register(user *model.User, c *fiber.Ctx) error
	Login(user *model.User, c *fiber.Ctx) (*model.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (usecase *userUsecase) Register(user *model.User, c *fiber.Ctx) error {
	if user.Email == "" {
		return c.JSON(WebResponse{
			Code: 401,
			Status: "BAD_REQUEST",
			Data: "Email cannot be empty",
		})
	}
	if user.Password == "" {
		return c.JSON(WebResponse{
			Code: 401,
			Status: "BAD_REQUEST",
			Data: "Password cannot be empty",
		})
	}
	existData, _ := usecase.userRepo.FindOneByEmail(user.Email)
	if existData != nil {
		return c.JSON(WebResponse{
			Code: 401,
			Status: "BAD_REQUEST",
			Data: fmt.Sprintf("Email %s already exists", user.Email),
		})
	}
	hashedPassword := helpers.HashPassword([]byte(user.Password))
	user.Id = uuid.New().String()
	user.Password = string(hashedPassword)
	return usecase.userRepo.CreateUser(user)
}

func (usecase *userUsecase) Login(user *model.User, c *fiber.Ctx) (*model.User, error) {
	existData, err := usecase.userRepo.FindOneByEmail(user.Email)
	if err != nil {
		return nil, fmt.Errorf("usecase.Login(): %w", err)
	}
	if existData == nil {
		return nil, c.JSON(WebResponse{
			Code: 401,
			Status: "BAD_REQUEST",
			Data: "Email is not registered",
		})
	}
	checkPassword := helpers.ComparePassword([]byte(existData.Password), []byte(user.Password))
	if !checkPassword {
		return nil, c.JSON(WebResponse{
			Code: 401,
			Status: "BAD_REQUEST",
			Data: errors.New("invalid password").Error(),
		})
	}
	return existData, nil
}