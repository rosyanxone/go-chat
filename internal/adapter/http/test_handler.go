package http

import (
	"crypto/rand"
	"encoding/base64"
	"go-chat/internal/app"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	service *app.UserService
}

func NewTestHandler(rg *gin.RouterGroup, userService *app.UserService) {
	userH := &TestHandler{service: userService}

	rg.GET("/test", userH.getTest)
}

func (h *TestHandler) getTest(c *gin.Context) {
	b := make([]byte, 32) // 32 bytes = 256-bit random
	_, err := rand.Read(b)

	if err != nil {
		log.Println("DB Error:", err.Error())
	}

	token := base64.RawURLEncoding.EncodeToString(b)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Test successful",
		"data":    token,
	})
}
