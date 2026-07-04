package db

import (
	"context"
	"errors"
	"fmt"
	"go-chat/internal/adapter/models"
	port "go-chat/internal/port/db"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) port.AuthPort {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) UpdateToken(ctx context.Context, personalToken *models.PersonalAccessToken) error {
	result := r.db.WithContext(ctx).
		Where("tokenable_id = ? AND name = ?", personalToken.TokenableID, personalToken.Name).
		Delete(personalToken)

	if result.Error != nil {
		return fmt.Errorf("database delete token error: %w", result.Error)
	}

	result = r.db.Create(personalToken)

	if result.Error != nil {
		return fmt.Errorf("database create token error: %w", result.Error)
	}

	return nil
}

func (r *AuthRepository) GetUserByToken(ctx context.Context, tokenID string, tokenHash string) (*models.User, error) {
	var token models.PersonalAccessToken

	// Verify the token exists and the hash matches
	result := r.db.WithContext(ctx).
		Where("id = ? AND token = ?", tokenID, tokenHash).
		First(&token)

	if result.Error != nil {
		return nil, errors.New("invalid or expired token")
	}

	// Fetch the user associated with this token's UserID
	var user models.User

	userResult := r.db.WithContext(ctx).
		Preload("Roles").
		First(&user, token.TokenableID)

	if userResult.Error != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
