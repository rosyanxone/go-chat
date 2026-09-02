package app

import (
	"context"
	"fmt"
	"go-chat/internal/adapter/dto"
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

func (s *UserService) GetContact(ctx context.Context, userID string, search string, page int) ([]dto.UserContactRow, error) {
	limit := 25
	offset := (page - 1) * limit

	return s.repo.GetContact(ctx, userID, search, offset)
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

func (s *UserService) UpdateUserPin(ctx context.Context, userID string, pin string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repo.UpdateUserPin(ctx, userID, string(hashedPassword))
}

func (s *UserService) UpdateUserData(ctx context.Context, nik string, email *string, phoneNumber *string) error {
	return s.repo.UpdateUserData(ctx, nik, email, phoneNumber)
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
