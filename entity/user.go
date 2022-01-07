package entity

import "time"

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  Password
	CreatedAt time.Time
}

func NewUser(username, email string, password Password, createdAt time.Time) *User {
	return &User{
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
	}
}
