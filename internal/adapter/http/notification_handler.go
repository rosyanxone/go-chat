package http

import (
	"fmt"
	"go-chat/internal/adapter/dto"
	"go-chat/internal/app"
	"go-chat/internal/domain"
	"go-chat/internal/shared/convert"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service          *app.NotificationService
	authService      *app.AuthService
	userService      *app.UserService
	chatService      *app.ChatService
	broadcastService *app.BroadcastService
}

func NewNotificationHandler(
	rg *gin.RouterGroup,
	service *app.NotificationService,
	authService *app.AuthService,
	userService *app.UserService,
	chatService *app.ChatService,
	broadcastService *app.BroadcastService,
) {
	h := &NotificationHandler{service, authService, userService, chatService, broadcastService}

	web := rg.Group("/web")

	// Browser needs, before it can call PushManager.subscribe().
	web.GET("/vapid-public-key", h.VAPIDPublicKey)

	web.Use(AuthMiddleware(authService))
	{
		web.POST("/subscribe", h.Subscribe)
		web.POST("/unsubscribe", h.Unsubscribe)
	}

	rg.Use(AuthMiddleware(authService))
	{
		rg.POST("/notify", h.Notify)
	}
}

func (h *NotificationHandler) VAPIDPublicKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "VAPID public key",
		"data": gin.H{
			"public_key": h.service.VAPIDPublicKey(),
		},
	})
}

func (h *NotificationHandler) Subscribe(c *gin.Context) {
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

	var req dto.PushSubscriptionRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid payload format",
			"data":    nil,
		})
		return
	}

	err = h.service.UpdateUserSubscription(c.Request.Context(), user.ID, req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to save subscription",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Subscription saved successfully",
		"data":    nil,
	})
}

func (h *NotificationHandler) Unsubscribe(c *gin.Context) {
	var req struct {
		Endpoint string `json:"endpoint" binding:"required"`
	}

	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid payload format",
			"data":    nil,
		})
		return
	}

	err = h.service.UserUnsubscribe(c, req.Endpoint)

	if err != nil {
		c.Error(err)

		if err.Error() == "invalid or deleted endpoint" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "failed",
				"message": "Endpoint do not exists",
				"data":    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to do unsubscribe",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Unsubscribe successfully done",
		"data":    nil,
	})
}

// Notify fires a push notification at the logged-in user's own devices.
// Handy for confirming the whole pipeline (VAPID keys, subscription,
// service worker) actually works end to end.
func (h *NotificationHandler) Notify(c *gin.Context) {
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

	var req dto.NotifyRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid payload format",
			"data":    nil,
		})
		return
	}

	userTarget, err := h.userService.GetUserByPhoneNumber(c.Request.Context(), req.PhoneNumber)

	if err != nil {
		c.Error(err)

		if err.Error() == "user not found" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "failed",
				"message": "Nomor hp yang diberikan tidak terdaftar!",
				"data":    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Terjadi kesalahan",
			"data":    nil,
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

	chat, err := h.chatService.GetChat(c, userSender.ID, userTarget.ID)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to get chat data",
			"data": gin.H{
				"user_sender_id": userSender.ID,
				"user_target_id": userTarget.ID,
			},
		})
		return
	}

	cmd := dto.CreateMessageCommand{
		ChatID:     chat.ID,
		Message:    req.Message,
		Url:        req.Url,
		UniqueCode: req.Code,
	}

	h.chatService.CreateNewMessage(c.Request.Context(), cmd)

	roomChatUrl, err := h.chatService.GetRoomChatUrl(c, chat.ID, chat.ChatRoomID)

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

	err = h.service.SendToUser(c.Request.Context(), userTarget.ID, payload)

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

	userNik := convert.NullIfEmpty(userTarget.Employee.UniqueNumber)

	channel := fmt.Sprintf("message.%d", chat.ChatRoomID)
	fmt.Println("Sending message by broadcast...", channel)

	err = h.broadcastService.Send(
		c,
		channel,
		"chat.new.message",
		map[string]interface{}{
			"user_id": userSender.ID,
			"message": req.Message,
		},
	)

	if err != nil {
		fmt.Println("WebSocket broadcast error:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Message successfully sent",
		"data": gin.H{
			"chat_id": chat.ID,
			"message": req.Message,
			"target_user": dto.UserDataResponse{
				ID:          userTarget.ID,
				Name:        userTarget.Name,
				NIK:         userNik,
				Email:       userTarget.Email,
				PhoneNumber: userTarget.PhoneNumber,
				Role:        userTarget.Roles[0].Name,
			},
		},
	})
}
