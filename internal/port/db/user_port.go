package port

import "go-chat/internal/domain"

type UserPort interface {
	GetAll() ([]domain.User, error)
}
