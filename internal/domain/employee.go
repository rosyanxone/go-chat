package domain

import "time"

type Employee struct {
	ID           uint `gorm:"primaryKey"`
	UserID       uint // UserID foreign key
	UniqueNumber string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
