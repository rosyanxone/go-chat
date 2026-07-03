package http

import (
	"go-chat/internal/app"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *app.UserService
}

func NewUserHandler(rg *gin.RouterGroup, service *app.UserService) {
	h := &UserHandler{service: service}

	rg.GET("/users", h.getUsers)
	rg.GET("/user/email", h.getUserByEmail)
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
