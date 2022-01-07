package repository

import (
	"database/sql"
	"errors"
	"github.com/uCibar/bootcamp-radio/entity"
)

type UserPostgresRepository struct {
	db *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

func (repository *UserPostgresRepository) Create(user *entity.User) error {
	var insertedId int64

	err := repository.db.QueryRow("INSERT INTO users(username, email, password) VALUES($1, $2, $3) RETURNING id;",
		user.Username, user.Email, user.Password.Hashed).Scan(&insertedId)
	if err != nil {
		return err
	}

	user.ID = insertedId

	return err
}

func (repository *UserPostgresRepository) GetByID(id int64) (*entity.User, error) {
	return repository.getUserBy("id", id)
}

func (repository *UserPostgresRepository) GetByEmail(email string) (*entity.User, error) {
	return repository.getUserBy("email", email)
}

func (repository *UserPostgresRepository) GetByUsername(username string) (*entity.User, error) {
	return repository.getUserBy("username", username)
}

func (repository *UserPostgresRepository) getUserBy(column string, value interface{}) (*entity.User, error) {
	var query string

	switch column {
	case "id":
		query = "SELECT * FROM users WHERE id = $1"
	case "email":
		query = "SELECT * FROM users WHERE email = $1"
	case "username":
		query = "SELECT * FROM users WHERE username = $1"
	}

	var u entity.User

	err := repository.db.QueryRow(query, value).Scan(&u.ID, &u.Username, &u.Email, &u.Password.Hashed, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &u, nil
}
