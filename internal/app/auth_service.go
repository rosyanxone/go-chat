package app

import (
	"context"
	"errors"
	"go-chat/internal/adapter/models"
	"go-chat/internal/app/auth"
	port "go-chat/internal/port/db"
	"strings"
)

type AuthService struct {
	repo port.AuthPort
}

func NewAuthService(r port.AuthPort) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) UpdateNewToken(ctx context.Context, personalToken *models.PersonalAccessToken) error {
	return s.repo.UpdateToken(ctx, personalToken)
}

// middleware purpose
func (s *AuthService) GetUserFromBearerToken(ctx context.Context, bearerToken string) (*models.User, error) {
	parts := strings.Split(bearerToken, "|")
	if len(parts) != 2 {
		return nil, errors.New("Invalid token format")
	}

	tokenID := parts[0]
	plainToken := parts[1]

	tokenHash := auth.HashToken(plainToken)

	return s.repo.GetUserByToken(ctx, tokenID, tokenHash)
}
