package main

import (
	"context"

	"go-server/config"
	"go-server/handler"

	"go-server/repository"
	"go-server/ws"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func transitionHandler(r *gin.Engine, UserHandler *handler.UserHandler, hub *ws.Hub) {
	r.POST("/api/register", UserHandler.CreateUser)
	r.POST("/api/login", UserHandler.GetUser)
	r.GET("/users/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/users/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})
}

func runServer() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := config.ConnectPostgres()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Ошибка получения sql.DB: %v", err)
	}
	defer sqlDB.Close()

	config.ConnectRedis()

	if err := db.AutoMigrate(&config.Users{}); err != nil {
		log.Fatalf("Ошибка миграции БД: %v", err)
	}

	userRepo, err := repository.NewUserRepo(db)
	if err != nil {
		log.Fatalf("Ошибка инициализации репозитория: %v", err)
	}

	userHandler, err := handler.NewUserHandler(userRepo)
	if err != nil {
		log.Fatalf("Ошибка инициализации хэндлера: %v", err)
	}
	hub := ws.NewHub()
	go hub.Run()
	transitionHandler(r, userHandler, hub)

	srv := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Сервер запущен на :3000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	<-quit
	log.Println("Получен сигнал на завершение работы...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %v", err)
	}

	log.Println("Сервер успешно завершил работу")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	switch os.Getenv("gin_mode") {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
		log.Println("режим отладки не указан по умолчанию сервер запущен с Debug модом")
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Паника в main: %v", r)
		}
	}()

	runServer()
}
