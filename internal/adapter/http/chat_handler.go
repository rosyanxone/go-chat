package http

import (
	"go-chat/internal/app"
	"go-chat/internal/domain"
	"go-chat/internal/shared/convert"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service     *app.ChatService
	authService *app.AuthService
}

func NewChatHandler(rg *gin.RouterGroup, service *app.ChatService, authService *app.AuthService) {
	h := &ChatHandler{service, authService}

	chat := rg.Group("/chat")
	chat.Use(AuthMiddleware(authService))
	{
		chat.GET("/rooms", h.getRooms)
		chat.POST("/messages", h.getMessages)
	}
}

func (h *ChatHandler) getRooms(c *gin.Context) {
	userData, exists := c.Get("currentUser")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mengambil data user",
			"data":    nil,
		})
		return
	}

	user := userData.(*domain.User)

	chatRooms, err := h.service.GetRooms(c, strconv.FormatUint(uint64(user.ID), 10))

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Terjadi kesalahan",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil room chat",
		"data":    chatRooms,
	})
}

func (h *ChatHandler) getMessages(c *gin.Context) {
	userData, exists := c.Get("currentUser")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mengambil data user",
			"data":    nil,
		})
		return
	}

	user := userData.(*domain.User)

	var req struct {
		ChatRoomID string `json:"chat_room_id" binding:"required"`
		Page       string `json:"page"`
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

	page := req.Page

	if page == "" {
		page = "1"
	}

	messages, err := h.service.GetAndReadMessages(c, req.ChatRoomID, convert.UintToString(user.ID), convert.StringToInt(page))

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Terjadi kesalahan",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User messages successfully retrieved",
		"data":    messages,
	})
}
