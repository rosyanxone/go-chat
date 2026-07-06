package app

import (
	"context"
	"go-chat/internal/domain"
	port "go-chat/internal/port"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(r port.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	return s.repo.FindByPhoneNumber(ctx, phoneNumber)
}
