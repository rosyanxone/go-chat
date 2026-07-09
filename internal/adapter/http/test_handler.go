package http

import (
	"go-chat/internal/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	userService *app.UserService
	testService *app.TestService
	chatService *app.ChatService
}

func NewTestHandler(rg *gin.RouterGroup, userService *app.UserService, testService *app.TestService, chatService *app.ChatService) {
	h := &TestHandler{userService, testService, chatService} // equals to: {userService: userService, testService: testService}

	rg.POST("/test", h.getTest)
	rg.POST("/prompt", h.tryPrompt)
}

func (h *TestHandler) getTest(c *gin.Context) {
	chatRooms, err := h.chatService.GetRooms(c, "1175")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Request payload tidak valid",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Test successful",
		"data":    chatRooms,
	})
}

func (h *TestHandler) tryPrompt(c *gin.Context) {
	var req struct {
		Body string `json:"body" binding:"required"`
	}

	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Request payload tidak valid",
			"data":    nil,
		})
		return
	}

	result, respond := h.testService.Prompt(req.Body)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Test successful",
		"data": gin.H{
			"answer":  result,
			"respond": respond,
		},
	})
}
