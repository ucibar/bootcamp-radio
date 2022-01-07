package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/uCibar/bootcamp-radio/entity"
)

type Service struct {
	userService UserService
}

func NewService(userService UserService) *Service {
	return &Service{userService: userService}
}

func (service *Service) AuthenticateWithToken(auth *entity.Auth) (*entity.AuthToken, error) {
	u, err := service.userService.GetByUsername(auth.Username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, entity.ErrUserNotFound
	}

	if !u.Password.Compare(auth.Password) {
		return nil, entity.ErrPasswordIncorrect
	}

	return entity.NewAuthToken(u)
}

func (service *Service) VerifyToken(tokenString string) (*entity.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, entity.ErrAuthTokenInvalid
		}

		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, entity.ErrAuthTokenExpired
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)

	u, err := service.userService.GetByID(int64(userID))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (service *Service) RegisterUser(register *entity.Register) error {
	if register.Password != register.PasswordRepeat {
		return entity.ErrPasswordRepeatIncorrect
	}

	p, err := entity.NewPassword(register.Password)
	if err != nil {
		return err
	}

	u := entity.User{Username: register.Username, Email: register.Email, Password: p}

	return service.userService.CreateUser(&u)
}
