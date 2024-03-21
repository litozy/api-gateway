package controller

import (
	"errors"
	"fmt"
	"net/http"
	"service-employee/helpers"
	"service-employee/model"
	"service-employee/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var user_uri string = "http://service-user:3001/user"

type WebResponse struct {
	Code int
	Status string
	Data interface{}
}

type EmployeeController interface {
	CreateEmployee(c *fiber.Ctx) error
}

type empController  struct {
	usecase usecase.EmployeeUsecase
}

func NewEmployeeController(usecase usecase.EmployeeUsecase) EmployeeController {
	return &empController{usecase: usecase}
}

func (emp *empController) CreateEmployee(c *fiber.Ctx) error {
	var requestBody model.Employee
	requestBody.Id = uuid.New().String()
	err := c.BodyParser(&requestBody)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":      false,
			"errorMessage": "Invalid JSON data",
		})
	}

	access_token := c.Cookies("access_token")
	if len(access_token) == 0 {
		return c.Status(401).SendString("Invalid token: Access token missing")
	}

	req, err := http.NewRequest("GET", user_uri + "/auth", nil)
	fmt.Println(req)
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic(err)
	}

	// Set headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("access_token", access_token)

	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error sending request")
	}
	defer resp.Body.Close()

	// Print the response
	// fmt.Println("Response Status:", resp.Status)
	// fmt.Println("Response Headers:", resp.Header)

	if resp.Status != "200 OK" {
		c.Status(401).SendString("invalid token")
	}

	err = emp.usecase.CreateEmployee(&requestBody, c)
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
		Data: requestBody,
	})
}