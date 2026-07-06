package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-chat/internal/adapter/db"
	"go-chat/internal/adapter/http"
	"go-chat/internal/app"
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
	// defer conn.Close()

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

	// Pack all services
	services := http.Services{
		UserService: app.NewUserService(userRepo),
		AuthService: app.NewAuthService(authRepo),
	}

	// Creates a blank router default by Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// NOT using any reverse proxy
	router.SetTrustedProxies(nil)

	// Create the base /api group by Gin
	api := router.Group("/api")

	// Pass the group AND the routes package service
	http.RegisterRoute(api, services)

	log.Printf("API running on http://localhost:%s\n", appPort)
	router.Run(":" + appPort)
}
