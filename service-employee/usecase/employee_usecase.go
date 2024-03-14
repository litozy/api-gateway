package usecase

import (
	"service-employee/helpers"
	"service-employee/model"
	"service-employee/repository"

	"github.com/gofiber/fiber/v2"
)

type EmployeeUsecase interface {
	CreateEmployee(emp *model.Employee, c *fiber.Ctx) error
}

type empUsecase struct {
	repository repository.EmployeeRepository
}

func NewEmployeeUsecase(repository repository.EmployeeRepository) EmployeeUsecase {
	return &empUsecase{repository: repository}
}

func (usecase *empUsecase) CreateEmployee(emp *model.Employee, c *fiber.Ctx) error {
	if emp.Name == "" {
		return &helpers.WebResponse{
			Code: 401,
			Status: "BAD_REQUEST",
			Data: "name cannot be empty",
		}
	}
	return usecase.repository.CreateEmployee(emp)
}