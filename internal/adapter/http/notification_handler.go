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
	web.Use(AuthMiddleware(authService))
	{
		web.POST("/subscribe", h.Subscribe)
		web.POST("/unsubscribe", h.Unsubscribe)
	}
}

func (h *NotificationHandler) Subscribe(c *gin.Context) {
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

	err = h.service.UpdateUserSubscription(c.Request.Context(), req.UserID, req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to save subscription",
			"data":    nil,
		})
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
