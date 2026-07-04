package port

import (
	"context"
	"go-chat/internal/adapter/models"
)

type AuthPort interface {
	UpdateToken(ctx context.Context, personalToken *models.PersonalAccessToken) error
	GetUserByToken(ctx context.Context, tokenID string, tokenHash string) (*models.User, error)
}
