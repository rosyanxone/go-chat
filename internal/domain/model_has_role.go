package domain

type ModelHasRole struct {
	RoleID    uint   `gorm:"primaryKey;column:role_id"`
	ModelType string `gorm:"primaryKey;column:model_type"`
	ModelID   uint   `gorm:"primaryKey;column:model_id"`
}
