package repository

import (
	"database/sql"
	"fmt"
	"log"
	"service-user/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var dummyUsers = []model.User{
	{
		Id:       "C001",
		Email:    "testing1@example.com",
		Password: "password1",
	},
	{
		Id:       "C002",
		Email:    "testing2@example.com",
		Password: "password2",
	},
}

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDb *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	mockDb , mockSql , err := sqlmock.New()
	if err != nil {
		log.Fatal("An Error when opening a stub datbase connection", err)
	}

	suite.mockDb = mockDb
	suite.mockSql = mockSql
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func (suite *UserRepositoryTestSuite) TestUserCreate_Success() {
	dummyUser := &dummyUsers[0]
	suite.mockSql.ExpectExec(`INSERT INTO users \(id, email, password\) VALUES \(\$1, \$2, \$3\)`).WithArgs(dummyUsers[0].Id, dummyUsers[0].Email, dummyUsers[0].Password).WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewUserRepository(suite.mockDb)
	err := repo.CreateUser(dummyUser)

	assert.Nil(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestUserCreate_Failed() {
    dummyUser := &dummyUsers[0]
    suite.mockSql.ExpectExec("^INSERT INTO users \\(id, email, password\\) VALUES \\(\\$1, \\$2, \\$3\\)$").WillReturnError(fmt.Errorf("failed"))
    repo := NewUserRepository(suite.mockDb)
    err := repo.CreateUser(dummyUser)

    assert.Error(suite.T(), err)
    assert.EqualError(suite.T(), err, "error on userRepository.CreateEmployee() : failed")
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
