package port

import (
	"context"
	"go-chat/internal/adapter/models"
)

type AuthPort interface {
	UpdateToken(ctx context.Context, personalToken *models.PersonalAccessToken) error
}
