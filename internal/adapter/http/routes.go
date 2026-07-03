package http

import (
	"go-chat/internal/app"

	"github.com/gin-gonic/gin"
)

// Business logic services
type Services struct {
	UserService *app.UserService
}

func RegisterRoute(api *gin.RouterGroup, services Services) {
	NewUserHandler(api, services.UserService)
	NewTestHandler(api, services.UserService)
}
