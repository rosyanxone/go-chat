package app

import (
	"go-chat/internal/domain"
	port "go-chat/internal/port/db"
)

type UserService struct {
	repo port.UserPort
}

func NewUserService(r port.UserPort) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetUsers() ([]domain.User, error) {
	return s.repo.GetAll()
}
