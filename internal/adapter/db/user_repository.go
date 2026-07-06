package db

import (
	"context"
	"errors"
	"fmt"
	"go-chat/internal/domain"
	port "go-chat/internal/port"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	result := r.db.WithContext(ctx).Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", result.Error)
	}

	return users, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
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

func (r *UserRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	var user domain.User

	result := r.db.WithContext(ctx).
		Where("phone_number = ?", phoneNumber).
		Preload("Roles").
		First(&user)

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
