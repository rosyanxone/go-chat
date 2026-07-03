package app

// func Login(c *gin.Context) {
// 	var req struct {
// 		Email      string `json:"email" binding:"required,email"`
// 		Password   string `json:"password" binding:"required"`
// 		DeviceName string `json:"device_name" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(400, gin.H{"message": "request tidak valid"})
// 		return
// 	}

// 	var user models.User
// 	// if err := db.Query("email = ?", req.Email).First(&user).Error; err != nil {
// 	// 	c.JSON(401, gin.H{"message": "email atau password salah"})
// 	// 	return
// 	// }
// 	// db.

// 	err := bcrypt.CompareHashAndPassword(
// 		[]byte(user.Password),
// 		[]byte(req.Password),
// 	)

// 	if err != nil {
// 		c.JSON(401, gin.H{"message": "email atau password salah"})
// 		return
// 	}

// 	plainToken, err := auth.GeneratePlainToken()
// 	if err != nil {
// 		c.JSON(500, gin.H{"message": "gagal membuat token"})
// 		return
// 	}

// 	tokenHash := auth.HashToken(plainToken)

// 	personalToken := models.PersonalAccessToken{
// 		UserID:    uint64(user.ID),
// 		Name:      req.DeviceName,
// 		TokenHash: tokenHash,
// 		Abilities: `["chat:read","chat:write"]`,
// 		IPAddress: c.ClientIP(),
// 		UserAgent: c.GetHeader("User-Agent"),
// 		ExpiresAt: nil, // tidak ada expiry
// 	}

// 	if err := db.Create(&personalToken).Error; err != nil {
// 		c.JSON(500, gin.H{"message": "gagal menyimpan token"})
// 		return
// 	}

// 	plainTextToken := auth.BuildPlainTextToken(personalToken.ID, plainToken)

// 	c.JSON(200, gin.H{
// 		"message": "login berhasil",
// 		"token":   plainTextToken,
// 		"user": gin.H{
// 			"id":    user.ID,
// 			"name":  user.Name,
// 			"email": user.Email,
// 		},
// 	})
// }
