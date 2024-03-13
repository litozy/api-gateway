package repository

import (
	"database/sql"
	"fmt"
	"service-employee/model"
)

type EmployeeRepository interface {
	CreateEmployee(emp *model.Employee) error
}

type empRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &empRepository{
		db: db,
	}
}

func (repository *empRepository) CreateEmployee(emp *model.Employee) error {
	qry := "INSERT INTO employees (id, name) VALUES ($1, $2)"
	_, err := repository.db.Exec(qry, &emp.Id, &emp.Name)
	if err != nil {
		return fmt.Errorf("error on employeeRepository.CreateEmployee() : %w", err)
	}
	return nil
}

