package port

import (
	"context"
	"go-chat/internal/adapter/models"
)

type UserPort interface {
	GetAll(ctx context.Context) ([]models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
}
