package domain

import "time"

type Role struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	GuardName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
