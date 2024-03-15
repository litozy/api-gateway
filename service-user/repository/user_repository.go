package repository

import (
	"database/sql"
	"fmt"
	"service-user/model"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	FindOneByEmail(email string) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (userRepo *userRepository) CreateUser(user *model.User) error {
	qry := "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)"
	_, err := userRepo.db.Exec(qry, &user.Id, &user.Email, &user.Password)
	if err != nil {
		return fmt.Errorf("error on userRepository.CreateEmployee() : %w", err)
	}
	return nil
}

func (userRepo *userRepository) FindOneByEmail(email string) (*model.User, error) {
	qry := "SELECT id, email, password FROM users WHERE email = $1"

	user := &model.User{}
	err := userRepo.db.QueryRow(qry, email).Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error on userRepository.FindOneByEmail() : %w", err)
	}
	return user, nil
}