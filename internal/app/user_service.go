package app

import (
	"context"
	"go-chat/internal/adapter/models"
	port "go-chat/internal/port/db"
)

type UserService struct {
	repo port.UserPort
}

func NewUserService(r port.UserPort) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	return s.repo.FindByPhoneNumber(ctx, phoneNumber)
}
