package repository

import (
	"database/sql"
	"fmt"
	"log"
	"service-employee/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var dummyEmployees = []model.Employee{
	{
		Id:       "C001",
		Name:    "name1",
	},
	{
		Id:       "C002",
		Name:    "name2",
	},
}

type EmployeeRepositoryTestSuite struct {
	suite.Suite
	mockDb *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *EmployeeRepositoryTestSuite) SetupTest() {
	mockDb , mockSql , err := sqlmock.New()
	if err != nil {
		log.Fatal("An Error when opening a stub datbase connection", err)
	}

	suite.mockDb = mockDb
	suite.mockSql = mockSql
}

func (suite *EmployeeRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func (suite *EmployeeRepositoryTestSuite) TestEmployeeCreate_Success() {
	dummyEmployee := &dummyEmployees[0]
	suite.mockSql.ExpectExec(`INSERT INTO employees \(id, name\) VALUES \(\$1, \$2\)`).WithArgs(dummyEmployees[0].Id, dummyEmployees[0].Name).WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewEmployeeRepository(suite.mockDb)
	err := repo.CreateEmployee(dummyEmployee)

	assert.Nil(suite.T(), err)
}


func (suite *EmployeeRepositoryTestSuite) TestEmployeeCreate_Failed() {
    dummyEmployee := &dummyEmployees[0]
    suite.mockSql.ExpectExec(`INSERT INTO employees \(id, name\) VALUES \(\$1, \$2\)`).WillReturnError(fmt.Errorf("failed"))
    repo := NewEmployeeRepository(suite.mockDb)
    err := repo.CreateEmployee(dummyEmployee)

    assert.Error(suite.T(), err)
    assert.EqualError(suite.T(), err, "error on employeeRepository.CreateEmployee() : failed")
}

func TestEmployeeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(EmployeeRepositoryTestSuite))
}
