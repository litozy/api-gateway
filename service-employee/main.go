package main

import (
	"fmt"
	"service-employee/config"
	"service-employee/controller"
	"service-employee/repository"
	"service-employee/usecase"

	"github.com/gofiber/fiber/v2"
)


func main() {
	empInfra := config.NewInfraManager().GetDB()
	empRepository := repository.NewEmployeeRepository(empInfra)
	empUsecase := usecase.NewEmployeeUsecase(empRepository)
	empController := controller.NewEmployeeController(empUsecase)
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hi from service-employee")
	})
	app.Post("/employee", empController.CreateEmployee)
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Headers", "access_token")
		return c.Next()
	})

	port := 3002
	fmt.Printf("Service employee is running on :%d...\n", port)

	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error starting Service employee: %v\n", err)
	}
}
