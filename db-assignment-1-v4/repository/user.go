package repository

import (
	"a21hc3NpZ25tZW50/model"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	Add(user model.User) error
	CheckAvail(user model.User) error
	FetchByID(id int) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepository {
	return &userRepository{db}
}

func (u *userRepository) Add(user model.User) error {
	err := u.CheckAvail(user)
	if err != nil {
		return err
	}

	_, err = u.db.Exec(
		"INSERT INTO users (username, password) VALUES ($1, $2)",
		user.Username,
		user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) CheckAvail(user model.User) error {
	var exists bool
	err := u.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)",
		user.Username,
	).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("username %s sudah digunakan", user.Username)
	}
	return nil
}

func (u *userRepository) FetchByID(id int) (*model.User, error) {
	row := u.db.QueryRow("SELECT id, username, password FROM users WHERE id = $1", id)

	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
