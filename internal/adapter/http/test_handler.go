package http

import (
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

// b := make([]byte, 32) // 32 bytes = 256-bit random
// _, err := rand.Read(b)
// token := base64.RawURLEncoding.EncodeToString(b)

func (h *TestHandler) getTest(c *gin.Context) {
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
