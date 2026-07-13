package http

import (
	"go-chat/internal/app"
	"go-chat/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service     *app.NotificationService
	authService *app.AuthService
}

func NewNotificationHandler(rg *gin.RouterGroup, service *app.NotificationService, authService *app.AuthService) {
	h := &NotificationHandler{service, authService}

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

	var req domain.PushSubscriptionRequest

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

	payload := domain.PushPayload{
		Title: "go-chat-notify",
		Body:  "This is a test push notification 🎉",
		Url:   "/",
	}

	err := h.service.SendToUser(c.Request.Context(), user.ID, payload)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to send test notification",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Test notification sent",
		"data":    nil,
	})
}
