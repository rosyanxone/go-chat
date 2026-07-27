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
	ChatService         *app.ChatService
}

func RegisterRoute(api *gin.RouterGroup, services Services) {
	http.NewTestHandler(api, services.UserService, services.TestService, services.ChatService)
	http.NewUserHandler(api, services.UserService, services.AuthService)
	http.NewAuthHandler(api, services.UserService, services.AuthService)
	http.NewNotificationHandler(api, services.NotificationService, services.AuthService, services.UserService, services.ChatService)
	http.NewChatHandler(api, services.ChatService, services.AuthService)
}
