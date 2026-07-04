package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-chat/internal/app"
)

func AuthMiddleware(authService *app.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Unauthenticated: Missing token",
				"data":    nil,
			})
			return
		}

		// Validate the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Unauthenticated: Invalid token format",
				"data":    nil,
			})
			return
		}

		// Extract the token
		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := authService.GetUserFromBearerToken(c.Request.Context(), rawToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Unauthenticated: " + err.Error(),
				"data":    nil,
			})
		}
		fmt.Println(user)

		// Store the user object securely inside the Gin context
		// This makes the user data available to any handler that runs after this middleware
		c.Set("currentUser", user)

		// Continue to the next handler
		c.Next()
	}
}
