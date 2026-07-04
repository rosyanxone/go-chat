package models

import (
	"time"
)

// User represents the database schema for a user
type User struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Email       string `gorm:"uniqueIndex;not null"`
	PhoneNumber string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// This single line bridges Go to Laravel's model_has_roles table.
	// - many2many: tells GORM the pivot table name
	// - polymorphic: tells GORM to look for `model_id` and `model_type` columns
	// - polymorphicValue: injects Laravel's exact namespace string
	Roles []Role `gorm:"many2many:model_has_roles;joinForeignKey:model_id;joinReferences:role_id"`
}

func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}
