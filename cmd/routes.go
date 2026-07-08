package main

import (
	"go-chat/internal/adapter/http"
	"go-chat/internal/app"

	"github.com/gin-gonic/gin"
)

// Business logic services
type Services struct {
	UserService         *app.UserService
	AuthService         *app.AuthService
	TestService         *app.TestService
	NotificationService *app.NotificationService
}

func RegisterRoute(api *gin.RouterGroup, services Services) {
	http.NewTestHandler(api, services.UserService, services.TestService)
	http.NewUserHandler(api, services.UserService)
	http.NewAuthHandler(api, services.UserService, services.AuthService)
	http.NewNotificationHandler(api, services.NotificationService, services.AuthService)
}
