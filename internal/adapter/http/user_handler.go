package http

import (
	"go-chat/internal/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *app.UserService
}

func NewUserHandler(rg *gin.RouterGroup, service *app.UserService) {
	h := &UserHandler{service: service}

	rg.GET("/users", h.getUsers)
}

func (h *UserHandler) getUsers(c *gin.Context) {
	users, err := h.service.GetUsers(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err,
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
