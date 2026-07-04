package app

import (
	"context"
	"go-chat/internal/adapter/models"
	port "go-chat/internal/port/db"
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
