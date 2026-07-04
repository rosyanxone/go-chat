package http

import (
	"go-chat/internal/adapter/models"
	"go-chat/internal/app"
	"go-chat/internal/app/auth"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService *app.UserService
	authService *app.AuthService
}

func NewAuthHandler(rg *gin.RouterGroup, userService *app.UserService, authService *app.AuthService) {
	userServiceH := &AuthHandler{userService: userService, authService: authService}

	rg.POST("/login", userServiceH.Login)

	protected := rg.Group("/")
	protected.Use(AuthMiddleware(authService))
	{
		protected.GET("/user", userServiceH.GetMe)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"login" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Platform    string `json:"platform" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("DB Error:", err.Error())

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Request payload tidak valid",
			"data":    nil,
		})
		return
	}

	user, err := h.userService.GetUserByPhoneNumber(c.Request.Context(), req.PhoneNumber)

	if err != nil {
		log.Println("DB Error:", err.Error())

		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "Nomor hp atau pin Anda salah!",
			"data":    nil,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		log.Println("DB Error:", err.Error())

		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "Nomor hp atau pin Anda salah!",
			"data":    nil,
		})

		return
	}

	plainToken, err := auth.GeneratePlainToken()

	if err != nil {
		log.Println("DB Error:", err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal membuat token",
			"data":    nil,
		})

		return
	}

	tokenHash := auth.HashToken(plainToken)

	personalToken := models.PersonalAccessToken{
		TokenableID: uint64(user.ID),
		Name:        req.Platform,
		Token:       tokenHash,
		ExpiresAt:   nil, // tidak ada expiry
	}

	err = h.authService.UpdateNewToken(c, &personalToken)

	if err != nil {
		log.Println("DB Error:", err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menyimpan token",
			"data":    nil,
		})
		return
	}

	plainTextToken := auth.BuildPlainTextToken(uint64(personalToken.ID), plainToken)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "login berhasil",
		"data": gin.H{
			"id":           user.ID,
			"name":         user.Name,
			"phone_number": user.PhoneNumber,
			"token":        plainTextToken,
			"role":         user.Roles[0].Name,
		},
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	// Use "currentUser" is the exact key set in the middleware
	userData, exists := c.Get("currentUser")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to retrieve user context",
			"data":    nil,
		})
		return
	}

	// Type-cast the data back into your models.User struct
	user := userData.(*models.User)

	// Return the user
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User data successfully retrived",
		"data": gin.H{
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"role":         user.Roles[0].Name,
			"updated_at":   user.UpdatedAt,
		},
	})
}
