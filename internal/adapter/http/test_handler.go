package http

import (
	"go-chat/internal/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	userService *app.UserService
	testService *app.TestService
}

func NewTestHandler(rg *gin.RouterGroup, userService *app.UserService, testService *app.TestService) {
	h := &TestHandler{userService, testService} // equals to: {userService: userService, testService: testService}

	rg.POST("/test", h.getTest)
}

func (h *TestHandler) getTest(c *gin.Context) {
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
