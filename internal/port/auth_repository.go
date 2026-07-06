package port

import (
	"context"
	"go-chat/internal/domain"
)

type AuthRepository interface {
	UpdateToken(ctx context.Context, personalToken *domain.PersonalAccessToken) error
	GetUserByToken(ctx context.Context, tokenID string, tokenHash string) (*domain.User, error)
}
