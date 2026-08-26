package db

import (
	"context"
	"errors"
	"fmt"
	"go-chat/internal/adapter/dto"
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

func (r *UserRepository) GetContact(ctx context.Context, search string, offset int) ([]dto.UserContactRow, error) {
	var userContact []dto.UserContactRow

	query := r.db.WithContext(ctx).
		Table("users").
		Select("id", "name", "phone_number")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	err := query.Limit(25).
		Offset(offset).
		Find(&userContact).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("database error: %w", err)
	}

	return userContact, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	result := r.db.WithContext(ctx).
		Where("email = ?", email).
		Preload("Roles").
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return &user, nil
}

func (r *UserRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	var user domain.User

	result := r.db.WithContext(ctx).
		Where("phone_number = ?", phoneNumber).
		Preload("Roles").
		Preload("Employee").
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

func (r *UserRepository) GetRoleIDByName(ctx context.Context, roleName string) (uint, error) {
	var role domain.Role

	result := r.db.WithContext(ctx).
		Where("name = ?", roleName).
		First(&role)

	if result.Error != nil {
		return 0, fmt.Errorf("database error: %w", result.Error)
	}

	return role.ID, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User, employee *domain.Employee, roleID uint) error {
	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	newUser := domain.User{
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
	}

	if employee.UniqueNumber != "" {
		newUser.Employee = domain.Employee{
			UniqueNumber: employee.UniqueNumber,
		}
	}

	err := tx.Create(&newUser).Error

	if err != nil {
		return err
	}

	// Manually create the Spatie Pivot Record
	spatiePivot := domain.ModelHasRole{
		RoleID:  roleID,
		ModelID: newUser.ID,
	}

	err = tx.Create(&spatiePivot).Error

	if err != nil {
		return err
	}

	// Fetch the full role from database and attach it to newUser variable
	err = tx.Preload("Roles").First(&newUser, newUser.ID).Error

	if err != nil {
		return err
	}

	// Map the new User back to the original pointer so Handler knows the new ID!
	*user = newUser

	return tx.Commit().Error
}

func (r *UserRepository) UpdateUserName(ctx context.Context, userID string, updatedName string) (*domain.User, error) {
	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return nil, tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	var user domain.User

	tx.Table("users").
		Where("id = ?", userID).
		Preload("Roles").
		Preload("Employee").
		Updates(domain.User{
			Name: updatedName,
		}).
		First(&user)

	return &user, tx.Commit().Error
}

func (r *UserRepository) UpdateUserPin(ctx context.Context, userID string, password string) (*domain.User, error) {
	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return nil, tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	var user domain.User

	tx.Table("users").
		Where("id = ?", userID).
		Preload("Roles").
		Preload("Employee").
		Updates(domain.User{
			Password: password,
		}).
		First(&user)

	return &user, tx.Commit().Error
}

func (r *UserRepository) UpdateUserData(ctx context.Context, nik string, email *string, phoneNumber *string) error {
	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	var employee domain.Employee

	err := tx.Table("employees").
		Select("user_id").
		Where("unique_number = ?", nik).
		First(&employee).
		Error

	if err != nil {
		return err
	}

	updates := map[string]interface{}{}

	if email != nil {
		updates["email"] = *email
	}

	if phoneNumber != nil {
		updates["phone_number"] = *phoneNumber
	}

	// Nothing to update
	if len(updates) == 0 {
		return errors.New("no_data")
	}

	result := tx.Model(&domain.User{}).
		Where("id = ?", employee.UserID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return tx.Commit().Error
}
