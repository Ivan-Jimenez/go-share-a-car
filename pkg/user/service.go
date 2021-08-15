package user

import (
	"github.com/Ivan-Jimenez/go-share-a-car/pkg/entities"
)

type Service interface {
	NewUser(user *entities.User) (*entities.User, error)
	Login(credentials *entities.LoginCredentials) (*entities.User, string, error)
}

type service struct {
	repository Repository
}

func NewService() Service {
	return &service{repository: NewRepo()}
}

func (s *service) NewUser(user *entities.User) (*entities.User, error) {
	return s.repository.NewUser(user)
}

func (s *service) Login(credentials *entities.LoginCredentials) (*entities.User, string, error) {
	return s.repository.Login(credentials)
}
