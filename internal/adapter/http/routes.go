package http

import (
	"go-chat/internal/app"

	"github.com/gin-gonic/gin"
)

// Business logic services
type Services struct {
	UserService *app.UserService
	AuthService *app.AuthService
}

func RegisterRoute(api *gin.RouterGroup, services Services) {
	NewUserHandler(api, services.UserService)
	NewTestHandler(api, services.UserService)
	NewAuthHandler(api, services.UserService, services.AuthService)
}
