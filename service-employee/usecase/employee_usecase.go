package usecase

import (
	"service-employee/model"
	"service-employee/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EmployeeUsecase interface {
	CreateEmployee(emp *model.Employee, c *fiber.Ctx) error
}

type empUsecase struct {
	repository repository.EmployeeRepository
}

type WebResponse struct {
	Code int
	Status string
	Data interface{}
}

func NewEmployeeUsecase(repository repository.EmployeeRepository) EmployeeUsecase {
	return &empUsecase{repository: repository}
}

func (usecase *empUsecase) CreateEmployee(emp *model.Employee, c *fiber.Ctx) error {
	emp.Id = uuid.New().String()
	if emp.Name == "" {
		return c.JSON(WebResponse{
			Code: 401,
			Status: "BAD_REQUEST",
			Data: "Name Cannot Be Empty",
		})
	}
	return usecase.repository.CreateEmployee(emp)
}