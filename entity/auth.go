package entity

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Register struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

type AuthToken struct {
	Token string
}

func NewAuthToken(u *User) (*AuthToken, error) {
	j := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{"user_id": u.ID, "exp": time.Now().Add(time.Minute * 60).Unix()},
	)

	s, err := j.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	return &AuthToken{Token: s}, nil
}
