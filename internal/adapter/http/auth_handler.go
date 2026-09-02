package http

import (
	"fmt"
	"go-chat/internal/app"
	"go-chat/internal/domain"
	"go-chat/internal/shared/convert"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService *app.UserService
	authService *app.AuthService
}

func NewAuthHandler(rg *gin.RouterGroup, userService *app.UserService, authService *app.AuthService) {
	h := &AuthHandler{userService: userService, authService: authService}

	rg.POST("/validate/phone-number", h.ValidatePhoneNumber)
	rg.POST("/login", h.Login)
	rg.POST("/register", h.Register)
	rg.POST("/register/employee", h.RegisterEmployee)

	protected := rg.Group("/")
	protected.Use(AuthMiddleware(authService))
	{
		adminOnly := protected.Group("/admin")
		adminOnly.Use(RequireRole("admin", "root"))
		{
			adminOnly.GET("/user", h.GetMe)
		}

		protected.GET("/user", h.GetMe)
		protected.POST("/logout", h.Logout)
	}
}

func (h *AuthHandler) ValidatePhoneNumber(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	err := c.ShouldBind(&req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Nomor hp harus diisi!",
			"data": gin.H{
				"or": err.Error(),
			},
		})
		return
	}

	_, err = h.userService.GetUserByPhoneNumber(c.Request.Context(), req.PhoneNumber)

	if err != nil {
		// log.Println("DB Error:", err.Error())
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

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Nomor hp yang diberikan berhasil divalidasi",
		"data":    nil,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"login" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Platform    string `json:"platform" binding:"required"`
	}

	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Request payload tidak valid",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	user, err := h.userService.GetUserByPhoneNumber(c.Request.Context(), req.PhoneNumber)

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

		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "Nomor hp atau pin Anda salah!",
			"data":    nil,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "Nomor hp atau pin Anda salah!",
			"data":    nil,
		})
		return
	}

	plainTextToken, err := h.authService.GetUserNewToken(c, uint64(user.ID))

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menyimpan token",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "login berhasil",
		"data": gin.H{
			"user": gin.H{
				"id":           user.ID,
				"name":         user.Name,
				"phone_number": user.PhoneNumber,
			},
			"token": plainTextToken,
			"role":  user.Roles[0].Name,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		UserID int `json:"user_id" binding:"required"`
	}

	err := c.ShouldBind(&req)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Payload user_id required!",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	// err = h.authService.DeleteWebToken(c, req.UserID)

	// if err != nil {
	// 	log.Println("DB Error:", err.Error())

	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status":  "failed",
	// 		"message": "Terjadi kesalahan",
	// 		"data":    nil,
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User successfully logged out",
		"data":    nil,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		PhoneNumber string  `json:"phone_number" binding:"required"`
		Platform    string  `json:"platform" binding:"required"`
		Name        string  `json:"name" binding:"required"`
		Email       string  `json:"email" binding:"required,email"`
		Role        string  `json:"role" binding:"required"`
		NIK         *string `json:"nik" `
		Password    string  `json:"password" binding:"required"`
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

	user, _ := h.userService.GetUserByPhoneNumber(c.Request.Context(), req.PhoneNumber)

	if user != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "failed",
			"message": "Nomor hp telah didaftarkan!",
			"data":    nil,
		})
		return
	}

	user, _ = h.userService.GetUserByEmail(c.Request.Context(), req.Email)

	if user != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "failed",
			"message": "Email telah didaftarkan!",
			"data":    nil,
		})
		return
	}

	newUser := domain.User{
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password,
	}

	var newEmployee domain.Employee

	if req.Role == "employee" {
		if *req.NIK == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "NIK harus diisi!",
				"data":    nil,
			})
			return
		}

		newEmployee = domain.Employee{
			UniqueNumber: *req.NIK,
		}
	}

	err = h.userService.RegisterNewUser(c, &newUser, &newEmployee, req.Role)

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Terjadi kesalahan",
			"data":    nil,
		})
		return
	}

	plainTextToken, err := h.authService.GetUserNewToken(c, uint64(newUser.ID))

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menyimpan token",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User berhasil dibuat",
		"data": gin.H{
			"id":           newUser.ID,
			"name":         newUser.Name,
			"phone_number": newUser.PhoneNumber,
			"nik":          newUser.Employee.UniqueNumber,
			"role":         newUser.Roles[0].Name,
			"token":        plainTextToken,
		},
	})
}

func (h *AuthHandler) RegisterEmployee(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		NIK         string `json:"nik" binding:"required"`
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

	phoneNumber := convert.NormalizePhoneNumber(req.PhoneNumber)

	user, _ := h.userService.GetUserByPhoneNumber(c.Request.Context(), phoneNumber)

	if user != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "failed",
			"message": "Nomor hp telah didaftarkan!",
			"data":    nil,
		})
		return
	}

	if req.Email != "" {
		user, _ = h.userService.GetUserByEmail(c.Request.Context(), req.Email)

		if user != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "failed",
				"message": "Email telah didaftarkan!",
				"data":    nil,
			})
			return
		}
	}

	nik := req.NIK

	password := fmt.Sprintf("%s%s", nik[len(nik)-3:], phoneNumber[len(phoneNumber)-3:])

	newUser := domain.User{
		// Email:       req.Email,
		PhoneNumber: phoneNumber,
		Name:        req.Name,
		Password:    password,
	}

	if req.Email != "" {
		newUser.Email = req.Email
	}

	var newEmployee domain.Employee

	newEmployee = domain.Employee{
		UniqueNumber: req.NIK,
	}

	err = h.userService.RegisterNewUser(c, &newUser, &newEmployee, "employee")

	if err != nil {
		c.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Terjadi kesalahan",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User berhasil dibuat",
		"data": gin.H{
			"id":           newUser.ID,
			"name":         newUser.Name,
			"email":        newUser.Email,
			"phone_number": newUser.PhoneNumber,
			"nik":          newUser.Employee.UniqueNumber,
			"role":         newUser.Roles[0].Name,
		},
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	// Use "currentUser" is the exact key set in the middleware
	userData, exists := c.Get("currentUser")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to retrieve user context",
			"data":    nil,
		})
		return
	}

	// Type-cast the data back into your domain.User struct
	user := userData.(*domain.User)

	// Return the user
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User data successfully retrived",
		"data": gin.H{
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"nik":          user.Employee.UniqueNumber,
			"role":         user.Roles[0].Name,
			"updated_at":   user.UpdatedAt,
		},
	})
}
