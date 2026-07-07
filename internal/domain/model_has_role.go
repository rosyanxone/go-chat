package domain

type ModelHasRole struct {
	RoleID    uint   `gorm:"primaryKey;column:role_id"`
	ModelType string `gorm:"primaryKey;column:model_type;default:'App\\Models\\User'"`
	ModelID   uint   `gorm:"primaryKey;column:model_id"`
}
