package port

import (
	"context"
	"go-chat/internal/domain"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
}
