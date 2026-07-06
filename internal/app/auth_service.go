package app

import (
	"context"
	"errors"
	"go-chat/internal/app/auth"
	"go-chat/internal/domain"
	port "go-chat/internal/port"
	"strings"
)

type AuthService struct {
	repo port.AuthRepository
}

func NewAuthService(r port.AuthRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) UpdateNewToken(ctx context.Context, personalToken *domain.PersonalAccessToken) error {
	return s.repo.UpdateToken(ctx, personalToken)
}

func (s *AuthService) DeleteWebToken(ctx context.Context, userID string) error {
	return s.repo.DeleteWebTokenByUserID(ctx, userID)
}

// middleware purpose
func (s *AuthService) GetUserFromBearerToken(ctx context.Context, bearerToken string) (*domain.User, error) {
	parts := strings.Split(bearerToken, "|")
	if len(parts) != 2 {
		return nil, errors.New("Invalid token format")
	}

	tokenID := parts[0]
	plainToken := parts[1]

	tokenHash := auth.HashToken(plainToken)

	return s.repo.GetUserByToken(ctx, tokenID, tokenHash)
}
