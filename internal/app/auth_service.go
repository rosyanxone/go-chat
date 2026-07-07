package app

import (
	"context"
	"fmt"
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

func (s *AuthService) GetUserNewToken(ctx context.Context, userID uint64) (string, error) {
	plainToken, err := auth.GeneratePlainToken()

	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	tokenHash := auth.HashToken(plainToken)

	personalToken := domain.PersonalAccessToken{
		TokenableID: userID,
		Name:        "web",
		Token:       tokenHash,
		ExpiresAt:   nil, // tidak ada expiry
	}

	err = s.repo.UpdateToken(ctx, &personalToken)

	if err != nil {
		return "", fmt.Errorf("failed to save token to database: %w", err)
	}

	plainTextToken := auth.BuildPlainTextToken(uint64(personalToken.ID), plainToken)

	return plainTextToken, nil
}

// middleware purpose
func (s *AuthService) UpdateLastUsedToken(ctx context.Context, tokenID string) error {
	return s.repo.UpdateLastUsedToken(ctx, tokenID)
}

func (s *AuthService) GetUserByToken(ctx context.Context, tokenID string, tokenHash string) (*domain.User, error) {
	return s.repo.GetUserByToken(ctx, tokenID, tokenHash)
}

func (s *AuthService) GetTokenByBearer(ctx context.Context, bearerToken string) (string, string) {
	parts := strings.Split(bearerToken, "|")

	var tokenID string
	var plainToken string

	if len(parts) != 2 {
		tokenID = ""
		plainToken = parts[0]
	} else {
		tokenID = parts[0]
		plainToken = parts[1]
	}

	tokenHash := auth.HashToken(plainToken)

	return tokenID, tokenHash
}
