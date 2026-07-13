package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-chat/internal/adapter/db"
	"go-chat/internal/app"
	"go-chat/internal/app/ai"
	"go-chat/internal/domain"
)

func main() {
	log.Println("Go running...")

	// load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	// store env vars
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	appPort := os.Getenv("APP_PORT")
	vapidPublicKey := os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey := os.Getenv("VAPID_PRIVATE_KEY")
	vapidSubject := os.Getenv("VAPID_SUBJECT")

	// build DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)
	conn, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		log.Fatal(err)
	}

	err = conn.SetupJoinTable(&domain.User{}, "Roles", &domain.ModelHasRole{})
	if err != nil {
		log.Fatal("Failed to setup join table")
	}

	// Configure the underlying connection pool
	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatal("failed to get underlying database: %w", err)
		return
	}

	sqlDB.SetMaxIdleConns(10)  // Max idle connections
	sqlDB.SetMaxOpenConns(100) // Max open connections
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Initialize Repositories
	userRepo := db.NewUserRepository(conn)
	authRepo := db.NewAuthRepository(conn)
	notificationRepo := db.NewNotificationRepository(conn)
	chatRepo := db.NewChatRepository(conn)

	// Inital AI client
	client := ai.Client()
	defer client.Close()

	// Pack all services
	services := Services{
		TestService:         app.NewTestService(client),
		UserService:         app.NewUserService(userRepo),
		AuthService:         app.NewAuthService(authRepo),
		NotificationService: app.NewNotificationService(notificationRepo, vapidPublicKey, vapidPrivateKey, vapidSubject),
		ChatService:         app.NewChatService(chatRepo),
	}

	// Creates a blank router default by Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8000",
			"http://localhost:8081",
			"http://localhost:3000",
			"https://chat.laut-timur.com",
			"https://notif.laut-timur.com",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// NOT using any reverse proxy
	router.SetTrustedProxies(nil)

	// Create the base /api group by Gin
	api := router.Group("/api")

	// Pass the group AND the routes package service
	RegisterRoute(api, services)

	log.Printf("API running on http://localhost:%s\n", appPort)
	err = router.Run(":" + appPort)

	if err != nil {
		panic(err)
	}
}
