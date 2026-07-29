package http

import (
	"go-chat/internal/adapter/dto"
	"go-chat/internal/app"
	"go-chat/internal/domain"
	"go-chat/internal/shared/convert"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	service     *app.UserService
	authService *app.AuthService
}

func NewUserHandler(rg *gin.RouterGroup, service *app.UserService, authService *app.AuthService) {
	h := &UserHandler{service: service}

	rg.GET("/users", h.getUsers)
	rg.GET("/user/email", h.getUserByEmail)

	protected := rg.Group("/")
	protected.Use(AuthMiddleware(authService))
	{
		user := protected.Group("/user")

		user.POST("/update", h.updateUserName)
		user.POST("/update/pin", h.updateUserPin)
	}
}

func (h *UserHandler) getUsers(c *gin.Context) {
	users, err := h.service.GetUsers(c.Request.Context())

	if err != nil {
		// Log the real error
		log.Println("DB Error:", err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Users fetched successfully",
		"data":    users,
	})
}

func (h *UserHandler) getUserByEmail(c *gin.Context) {
	email := c.Query("email")

	data, err := h.service.GetUserByEmail(c.Request.Context(), email)

	if err != nil {
		// Log the real error
		log.Println("DB Error:", err.Error())

		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "No user exists with this email",
				"data":    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Internal server error",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Test successful",
		"data":    data,
	})
}

func (h *UserHandler) updateUserName(c *gin.Context) {
	user, exists := c.Get("currentUser")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mengambil data user",
			"data":    nil,
		})
		return
	}

	userData := user.(*domain.User)

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	err := c.ShouldBind(&req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Request payload harus terpenuhi!",
			"data":    nil,
		})
		return
	}

	userResult, err := h.service.UpdateUserName(c, convert.UintToString(userData.ID), req.Name)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal memperbarui data user",
			"data": gin.H{
				"user_id": userData.ID,
				"name":    req.Name,
			},
		})
	}

	userNik := convert.NullIfEmpty(userResult.Employee.UniqueNumber)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "failed",
		"message": "User name updated successfully",
		"data": gin.H{
			"user": dto.UserDataResponse{
				ID:          userResult.ID,
				Name:        userResult.Name,
				NIK:         userNik,
				Email:       userResult.Email,
				PhoneNumber: userResult.PhoneNumber,
				Role:        userResult.Roles[0].Name,
			},
		},
	})
}

func (h *UserHandler) updateUserPin(c *gin.Context) {
	user, exists := c.Get("currentUser")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mengambil data user",
			"data":    nil,
		})
		return
	}

	userData := user.(*domain.User)

	var req struct {
		OldPin          string `json:"old_pin" binding:"required"`
		Pin             string `json:"pin" binding:"required"`
		PinConfirmation string `json:"pin_confirmation" binding:"required"`
	}

	err := c.ShouldBind(&req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Semua kolom wajib diisi!",
			"data":    nil,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(userData.Password),
		[]byte(req.OldPin),
	)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "failed",
			"message": "Pin lama yang diberikan salah!",
			"data":    nil,
		})
		return
	}

	if len(req.Pin) < 6 || len(req.Pin) > 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "failed",
			"message": "Pin harus terdiri dari 6 digit!",
			"data":    nil,
		})
		return
	}

	if req.Pin != req.PinConfirmation {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "failed",
			"message": "Konfirmasi pin baru tidak sesuai!",
			"data":    nil,
		})
		return
	}

	userResult, err := h.service.UpdateUserPin(c, convert.UintToString(userData.ID), req.Pin)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Request payload harus terpenuhi!",
			"data": gin.H{
				"user_id": userData.ID,
				"name":    req.Pin,
			},
		})
	}

	userNik := convert.NullIfEmpty(userResult.Employee.UniqueNumber)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "failed",
		"message": "User pin updated successfully",
		"data": gin.H{
			"user": dto.UserDataResponse{
				ID:          userResult.ID,
				Name:        userResult.Name,
				NIK:         userNik,
				Email:       userResult.Email,
				PhoneNumber: userResult.PhoneNumber,
				Role:        userResult.Roles[0].Name,
			},
		},
	})
}
