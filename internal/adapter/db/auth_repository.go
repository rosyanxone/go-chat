package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-chat/internal/domain"
	port "go-chat/internal/port"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) port.AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) UpdateToken(ctx context.Context, personalToken *domain.PersonalAccessToken) error {
	tx := r.db.WithContext(ctx).Begin()

	deleteResult := tx.Where("name = 'web' AND tokenable_id = ?", personalToken.TokenableID).
		Delete(&domain.PersonalAccessToken{})

	if deleteResult.Error != nil {
		tx.Rollback()
		return fmt.Errorf("database delete token error: %w", deleteResult.Error)
	}

	createResult := tx.Create(personalToken)

	if createResult.Error != nil {
		tx.Rollback()
		return fmt.Errorf("database create token error: %w", createResult.Error)
	}

	return tx.Commit().Error
}

func (r *AuthRepository) GetUserByToken(ctx context.Context, tokenID string, tokenHash string) (*domain.User, error) {
	var token domain.PersonalAccessToken

	// Verify the token exists and the hash matches
	query := r.db.WithContext(ctx).Where("token = ?", tokenHash)

	if tokenID != "" {
		query = query.Where("id = ?", tokenID)
	}

	err := query.First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired token")
		}

		return nil, fmt.Errorf("database error during token lookup: %w", err)
	}

	// Fetch the user associated with this token's UserID
	var user domain.User

	userResult := r.db.WithContext(ctx).
		Preload("Roles").
		First(&user, token.TokenableID)

	if userResult.Error != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (r *AuthRepository) GetTokenByHash(ctx context.Context, tokenHash string) (*domain.PersonalAccessToken, error) {
	var token domain.PersonalAccessToken

	err := r.db.WithContext(ctx).Where("token = ?", tokenHash).First(&token).Error

	if err != nil {
		return nil, fmt.Errorf("database error during token lookup: %w", err)
	}

	return &token, nil
}

func (r *AuthRepository) DeleteWebTokenByUserID(ctx context.Context, userID string) error {
	result := r.db.WithContext(ctx).
		Where("name = 'web' AND tokenable_id = ?", userID).
		Delete(&domain.PersonalAccessToken{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete user tokens: %w", result.Error)
	}

	return nil
}

func (r *AuthRepository) UpdateLastUsedToken(ctx context.Context, tokenID string) error {
	now := time.Now()

	err := r.db.Model(&domain.PersonalAccessToken{}).
		Where("id = ?", tokenID).
		Update("last_used_at", now).
		Error

	return err
}
