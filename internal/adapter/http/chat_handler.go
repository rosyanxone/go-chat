package http

import (
	"go-chat/internal/adapter/dto"
	"go-chat/internal/app"
	"go-chat/internal/domain"
	"go-chat/internal/shared/convert"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service             *app.ChatService
	authService         *app.AuthService
	userService         *app.UserService
	notificationService *app.NotificationService
}

func NewChatHandler(rg *gin.RouterGroup, service *app.ChatService, authService *app.AuthService, userService *app.UserService, notificationService *app.NotificationService) {
	h := &ChatHandler{service, authService, userService, notificationService}

	chat := rg.Group("/chat")
	chat.Use(AuthMiddleware(authService))
	{
		chat.GET("/rooms", h.GetRooms)
		chat.POST("/messages", h.GetMessages)
		chat.POST("/send", h.SendMessage)
	}
}

func (h *ChatHandler) GetRooms(c *gin.Context) {
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

	page := c.Query("page")

	chatRooms, err := h.service.GetRooms(c, strconv.FormatUint(uint64(user.ID), 10), convert.StringToInt(page))

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

func (h *ChatHandler) GetMessages(c *gin.Context) {
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
		ChatRoomID int `json:"chat_room_id" binding:"required"`
		Page       int `json:"page"`
	}

	err := c.ShouldBind(&req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Request payload harus terpenuhi!",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	chatRoomID := convert.IntToString(req.ChatRoomID)

	messages, err := h.service.GetAndReadMessages(c, chatRoomID, convert.UintToString(user.ID), req.Page)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mengambil dan membaca pesan",
			"data": gin.H{
				"chat_room_id": req.ChatRoomID,
				"user_id":      user.ID,
				"page":         req.Page,
			},
		})
		return
	}

	userTarget, err := h.service.GetMemberInfoByChatRoomId(c, convert.UintToString(user.ID), chatRoomID)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mengambil info anggota chat",
			"data": gin.H{
				"user_id":      user.ID,
				"chat_room_id": req.ChatRoomID,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User messages successfully retrieved",
		"data": gin.H{
			"id":           userTarget.ID,
			"title":        userTarget.Name,
			"phone_number": userTarget.PhoneNumber,
			"messages":     messages,
		},
	})
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	userSenderData, exists := c.Get("currentUser")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mengambil data user",
			"data":    nil,
		})
		return
	}

	userSender := userSenderData.(*domain.User)

	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Message     string `json:"message" binding:"required"`
		IsNotify    bool   `json:"is_notify" binding:"required"`
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

	userTarget, err := h.userService.GetUserByPhoneNumber(c, req.PhoneNumber)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mendapatkan user target",
			"data": gin.H{
				"phone_number": req.PhoneNumber,
			},
		})
		return
	}

	if userSender.ID == userTarget.ID {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "failed",
			"message": "Nomor hp yang diberikan identikal",
			"data": gin.H{
				"user_sender": userSender.PhoneNumber,
				"user_target": userTarget.PhoneNumber,
			},
		})
		return
	}

	chat, err := h.service.GetChat(c, userSender.ID, userTarget.ID)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mendapatkan room chat",
			"data": gin.H{
				"phone_number": req.PhoneNumber,
			},
		})
		return
	}

	if req.IsNotify {
		roomChatUrl, err := h.service.GetRoomChatUrl(c, chat.ID, chat.ChatRoomID)

		if err != nil {
			c.Error(err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Failed to get unread total message",
				"data": gin.H{
					"error":   err.Error(),
					"chat_id": chat.ID,
				},
			})
			return
		}

		payload := dto.PushPayload{
			Title: userSender.Name,
			Body:  req.Message,
			// Icon: example.logo.com,

			Actions: []dto.PushAction{
				{
					Title:  "Buka",
					Action: "open_url",
				},
			},

			Data: dto.PushData{
				URL: *roomChatUrl,
			},
		}

		err = h.notificationService.SendToUser(c.Request.Context(), userTarget.ID, payload)

		if err != nil {
			c.Error(err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Message notification failed to sent",
				"data": gin.H{
					"error":   err.Error(),
					"chat_id": chat.ID,
					"message": req.Message,
				},
			})
			return
		}
	}

	cmd := dto.CreateMessageCommand{
		ChatID:  chat.ID,
		Message: req.Message,
	}

	h.service.CreateNewMessage(c.Request.Context(), cmd)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Pesan berhasil dikirim",
		"data":    cmd,
	})
}
