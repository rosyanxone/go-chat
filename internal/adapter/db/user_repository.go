package db

import (
	"context"
	"errors"
	"fmt"
	"go-chat/internal/adapter/models"
	port "go-chat/internal/port/db"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) port.UserPort {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	result := r.db.WithContext(ctx).Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", result.Error)
	}

	return users, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)

	if result.Error != nil {
		// sql.ErrNoRows means the query succeeded, but the email doesn't exist
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}

		// A real database error occurred
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return &user, nil
}

func (r *UserRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&user)

	if result.Error != nil {
		// sql.ErrNoRows means the query succeeded, but the email doesn't exist
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}

		// A real database error occurred
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return &user, nil
}
