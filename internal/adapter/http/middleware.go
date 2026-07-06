package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-chat/internal/app"
	"go-chat/internal/domain"
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

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user from the context (set previously by AuthMiddleware)
		userData, exists := c.Get("currentUser")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Unauthenticated",
				"data":    nil,
			})
			return
		}

		user := userData.(*domain.User)

		// Check if the user has any of the allowed roles
		hasAccess := false
		for _, requiredRole := range allowedRoles {
			if user.HasRole(requiredRole) {
				hasAccess = true
				break
			}
		}

		// If they don't have the role, block the request with a 403 Forbidden
		if !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "failed",
				"message": "Forbidden: You do not have permission to access this resource",
				"data":    nil,
			})
			return
		}

		// Success! Let them through
		c.Next()
	}
}
