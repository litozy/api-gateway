package middleware

import (
	"fmt"
	"service-user/helpers"
	"service-user/repository"

	"github.com/gofiber/fiber/v2"
)

type Auth interface {
	Authentication(c *fiber.Ctx) error
}

type auth struct {
	userRepo repository.UserRepository
}

func NewAuth(userRepo repository.UserRepository) Auth {
	return &auth{userRepo: userRepo}
}

func (auth *auth) Authentication(c *fiber.Ctx) error {
	access_token := c.Cookies("access_token")

	if len(access_token) == 0 {
		return c.Status(401).SendString("Invalid token: Access token missing")
	}

	checkToken, err := helpers.VerifyToken(access_token)

	if err != nil {
		return c.Status(401).SendString("Invalid token: Failed to verify token")
	}

	fmt.Println(checkToken, "CEKKKK" ,checkToken["email"])

	existData, err := auth.userRepo.FindOneByEmail(checkToken["email"].(string))
	fmt.Println(existData)
	if err != nil {
		return fmt.Errorf("auth.Authentication(): %w", err)
	}
	if existData == nil {
		return &helpers.WebResponse{
			Code: 401,
			Status: "BAD_REQUEST",
			Data: "email is not registered",
		}
	}

	// Set user data in context for future use
	c.Locals("user", existData)

	// Continue processing if user is found
	return c.Next()
}