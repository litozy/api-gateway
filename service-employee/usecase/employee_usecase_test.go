package usecase

import (
	"errors"
	"service-employee/model"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyEmployees = []model.Employee{
	{
		Id:   "C001",
		Name: "name1",
	},
	{
		Id:   "C002",
		Name: "name2",
	},
}

type employeeUsecaseTestSuite struct {
	suite.Suite
	repoMock *repoMock
}

func (s *employeeUsecaseTestSuite) SetupTest() {
	s.repoMock = new(repoMock)
}

type repoMock struct {
	mock.Mock
}

func (r *repoMock) CreateEmployee(newEmployee *model.Employee) error {
	args := r.Called(newEmployee)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *employeeUsecaseTestSuite) TestEmployeeRegister_Success() {
	dummyEmployee := &model.Employee{
		Id:       dummyEmployees[0].Id,
		Name:    dummyEmployees[0].Name,
	}
	ctx := &fiber.Ctx{}
	s.repoMock.On("CreateEmployee", dummyEmployee).Return(nil)
	employeeUsecaseTest := NewEmployeeUsecase(s.repoMock)
	err := employeeUsecaseTest.CreateEmployee(dummyEmployee, ctx)
	assert.Nil(s.T(), err)
}

func (s *employeeUsecaseTestSuite) TestEmployeeRegister_Failed() {
	dummyEmployee := &model.Employee{
		Id:       dummyEmployees[0].Id,
		Name:    dummyEmployees[0].Name,
	}
	ctx := &fiber.Ctx{}
	expectedError := errors.New("failed")
	s.repoMock.On("CreateEmployee", dummyEmployee).Return(expectedError)
	employeeUsecaseTest := NewEmployeeUsecase(s.repoMock)
	err := employeeUsecaseTest.CreateEmployee(dummyEmployee, ctx)
	assert.EqualError(s.T(), err, expectedError.Error())
}

func TestEmployeeUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(employeeUsecaseTestSuite))
}
