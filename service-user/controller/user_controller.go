package controller

import (
	"errors"
	"fmt"
	"service-user/helpers"
	"service-user/model"
	"service-user/usecase"

	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type userController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(userUsecase usecase.UserUsecase) UserController {
    return &userController{userUsecase: userUsecase}
}

func (controller *userController) Register(c *fiber.Ctx) error {
	var requestBody model.User
	err := c.BodyParser(&requestBody)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":      false,
			"errorMessage": "Invalid JSON data",
		})
	}

	err = controller.userUsecase.Register(&requestBody, c)
	if err != nil {
		webResponse := &helpers.WebResponse{}
		if errors.As(err, &webResponse) {
			fmt.Printf("user.Register() 1: %v", err.Error())
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success":      false,
				"errorMessage": webResponse.Error(),
			})
		} else {
			fmt.Printf("user.Register() 2: %v", err.Error())
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success":      false,
				"errorMessage": "An error occurred while saving user data",
			})
		}
		return c.JSON(&helpers.WebResponse{
			Code: 401,
			Status: "BAD_REQUEST",
			Data: webResponse.Error(),
		})
	}

	return c.JSON(&helpers.WebResponse{
		Code: 201,
		Status: "OK",
		Data: requestBody.Email,
	})
}

func (controller *userController) Login(c *fiber.Ctx) error {
	var requestBody model.User
 
	err := c.BodyParser(&requestBody)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":      false,
			"errorMessage": "Invalid JSON data",
		})
	}

	data, err := controller.userUsecase.Login(&requestBody, c)
	if err != nil {
		webResponse := &helpers.WebResponse{}
		if errors.As(err, &webResponse) {
			fmt.Printf("user.Login() 1: %v", err.Error())
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success":      false,
				"errorMessage": webResponse.Error(),
			})
		} else {
			fmt.Printf("user.Login() 2: %v", err.Error())
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success":      false,
				"errorMessage": "An error occurred while saving user data",
			})
		}
		return c.JSON(&helpers.WebResponse{
			Code: 401,
			Status: "BAD_REQUEST",
			Data: webResponse.Error(),
		})
	}
	
	// if data == nil {
	// 	c.JSON(WebResponse{
	// 		Code: 401,
	// 		Status: "BAD_REQUEST",
	// 		Data: data,
	// 	})
	// }

	access_token := helpers.SignToken(requestBody.Email)
	c.Cookie(&fiber.Cookie{
        Name:  "access_token",
        Value: access_token,
        Path:  "/", 
    })

	return c.JSON(struct{
		Code int 
		Status string
		AccessToken string
		Data interface{}
	}{
		Code: 200,
		Status: "OK",
		AccessToken: access_token,
		Data: data,
	})
}

func Auth(c *fiber.Ctx) error {
	return c.JSON("OK")
}
