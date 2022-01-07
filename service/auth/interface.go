package auth

import "github.com/uCibar/bootcamp-radio/entity"

type UserService interface {
	CreateUser(user *entity.User) error
	GetByID(id int64) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
}
