package user

import "github.com/uCibar/bootcamp-radio/entity"

type Repository interface {
	Create(user *entity.User) error
	GetByID(id int64) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
}
