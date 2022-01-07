package user

import "github.com/uCibar/bootcamp-radio/entity"

type Service struct {
	userRepository Repository
}

func NewService(userRepository Repository) *Service {
	return &Service{userRepository: userRepository}
}

func (service *Service) CreateUser(u *entity.User) error {
	if service.EmailExist(u.Email) {
		return entity.ErrEmailExist
	}

	if service.UsernameExist(u.Username) {
		return entity.ErrUsernameExist
	}

	return service.userRepository.Create(u)
}

func (service *Service) GetByID(id int64) (*entity.User, error) {
	return service.userRepository.GetByID(id)
}

func (service *Service) GetByUsername(username string) (*entity.User, error) {
	return service.userRepository.GetByUsername(username)
}

func (service *Service) EmailExist(email string) bool {
	u, _ := service.userRepository.GetByEmail(email)
	if u == nil {
		return false
	}
	return true
}

func (service *Service) UsernameExist(username string) bool {
	u, _ := service.userRepository.GetByUsername(username)
	if u == nil {
		return false
	}
	return true
}
