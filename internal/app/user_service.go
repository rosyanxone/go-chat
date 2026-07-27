package app

import (
	"context"
	"fmt"
	"go-chat/internal/domain"
	port "go-chat/internal/port"

	"golang.org/x/crypto/bcrypt"
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

func (s *UserService) UpdateUserName(ctx context.Context, userID string, name string) (*domain.User, error) {
	return s.repo.UpdateUserName(ctx, userID, name)
}

func (s *UserService) RegisterNewUser(ctx context.Context, user *domain.User, employee *domain.Employee, roleName string) error {
	roleID, err := s.repo.GetRoleIDByName(ctx, roleName)

	if err != nil {
		return fmt.Errorf("failed to get role id: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)

	return s.repo.CreateUser(ctx, user, employee, roleID)
}
