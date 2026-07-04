package db

import (
	"context"
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
