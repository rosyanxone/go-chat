package port

import (
	"context"
	"go-chat/internal/domain"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
	GetRoleIDByName(ctx context.Context, roleName string) (uint, error)
	CreateUser(ctx context.Context, user *domain.User, employee *domain.Employee, roleID uint) error
	UpdateUserName(ctx context.Context, userID string, updatedName string) (*domain.User, error)
	UpdateUserPin(ctx context.Context, userID string, password string) (*domain.User, error)
}
