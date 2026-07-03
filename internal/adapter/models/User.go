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
}
