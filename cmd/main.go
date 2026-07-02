package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"go-chat/internal/adapter/db"
	userHttp "go-chat/internal/adapter/http"
	"go-chat/internal/app"
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
		"%s:%s@tcp(%s:%s)/%s",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)
	conn, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// wiring chat
	repo := db.NewUserRepository(conn)
	service := app.NewUserService(repo)

	r := gin.Default()
	userHttp.NewUserHandler(r, service)

	log.Printf("API running on http://localhost:%s\n", appPort)
	r.Run(":" + appPort)
}
