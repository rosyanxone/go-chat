package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-chat/internal/adapter/db"
	broadcaster "go-chat/internal/adapter/websocket"
	"go-chat/internal/app"
	"go-chat/internal/app/ai"
	"go-chat/internal/domain"
)

var (
	ctx                = context.Background()
	channelNamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{3,64}$`)
	rdb                = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	log.Println("Go running...")

	// load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found, using system environment")
	}

	// connect to redis
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Print("Failed to connect to Redis:", err)
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
		log.Fatalf("failed to get underlying database: %s", err)
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

	// Consuming redis client
	wsBroadcaster := broadcaster.NewBroadcaster(rdb)

	// Pack all services
	services := Services{
		TestService:         app.NewTestService(client),
		UserService:         app.NewUserService(userRepo),
		AuthService:         app.NewAuthService(authRepo),
		NotificationService: app.NewNotificationService(notificationRepo, vapidPublicKey, vapidPrivateKey, vapidSubject),
		ChatService:         app.NewChatService(chatRepo),
		BroadcastService:    app.NewBroadcastService(wsBroadcaster),
	}

	// Creates a blank router default by Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8000",
			"http://localhost:8081",
			"http://localhost:3000",

			"http://127.0.0.1:8000",
			"http://127.0.0.1:8081",
			"http://127.0.0.1:3000",

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

	// assign websocket
	// router.GET("/ws/:channel", handleWebSocket)
	router.GET("/ws/chat", handleWebSocket)

	// NOT using any reverse proxy
	// router.SetTrustedProxies(nil)

	// Create the base /api group by Gin
	api := router.Group("/api")

	// Pass the group AND the routes package service
	RegisterRoute(api, services)

	appUrl := os.Getenv("APP_URL")
	log.Printf("API running on %s\n", appUrl)
	err = router.Run(":" + appPort)

	if err != nil {
		panic(err)
	}
}

type WebSocketRequest struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// One Redis PubSub for this websocket connection.
	pubsub := rdb.Subscribe(ctx)
	defer pubsub.Close()

	redisChannel := pubsub.Channel()

	done := make(chan struct{})

	// Read messages coming FROM frontend.
	go func() {
		defer close(done)

		for {
			var request WebSocketRequest

			err := conn.ReadJSON(&request)
			if err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseNormalClosure,
				) {
					log.Println("WebSocket read error:", err)
				}

				return
			}

			if !channelNamePattern.MatchString(request.Channel) {
				log.Println("Invalid channel:", request.Channel)
				continue
			}

			if request.Event == "" {
				log.Println("Missing websocket event")
				continue
			}

			switch request.Event {

			case "subscribe":
				err := pubsub.Subscribe(ctx, request.Channel)
				if err != nil {
					log.Println("Redis subscribe error:", err)
					continue
				}

				log.Println("Subscribed:", request.Channel)

			case "unsubscribe":
				err := pubsub.Unsubscribe(ctx, request.Channel)
				if err != nil {
					log.Println("Redis unsubscribe error:", err)
					continue
				}

				log.Println("Unsubscribed:", request.Channel)

			default:
				payload := broadcaster.Event{
					Event:   request.Event,
					Channel: request.Channel,
					Data:    request.Data,
				}

				jsonData, err := json.Marshal(payload)
				if err != nil {
					log.Println("WebSocket payload marshal error:", err)
					continue
				}

				err = rdb.Publish(ctx, request.Channel, jsonData).Err()
				if err != nil {
					log.Println("Redis publish error:", err)
					continue
				}

				log.Println("Published:", request.Event, request.Channel)
			}
		}
	}()

	// Receive messages FROM Redis
	// and send them TO frontend.
	for {
		select {
		case <-done:
			return

		case msg, ok := <-redisChannel:
			if !ok {
				return
			}

			log.Println(
				"Redis message:",
				msg.Channel,
				msg.Payload,
			)

			err := conn.WriteMessage(
				websocket.TextMessage,
				[]byte(msg.Payload),
			)

			if err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseNormalClosure,
				) {
					log.Println("WebSocket write error:", err)
				}

				return
			}
		}
	}
}
