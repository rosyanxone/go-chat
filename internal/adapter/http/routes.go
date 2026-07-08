package http

import (
	"go-chat/internal/app"

	"github.com/gin-gonic/gin"
)

// Business logic services
type Services struct {
	UserService *app.UserService
	AuthService *app.AuthService
	TestService *app.TestService
}

func RegisterRoute(api *gin.RouterGroup, services Services) {
	NewUserHandler(api, services.UserService)
	NewTestHandler(api, services.UserService, services.TestService)
	NewAuthHandler(api, services.UserService, services.AuthService)
}
