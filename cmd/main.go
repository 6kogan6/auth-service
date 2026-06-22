package main

import (
	"auth-service/internal/config"
	"auth-service/internal/db"
	handler "auth-service/internal/handler"
	"auth-service/internal/middleware"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка конфига: %v", err)
	}
	router := gin.Default()

	pool, err := db.ConnectPostgres(cfg.DBURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к бд: %v", err)
	}
	defer pool.Close()
	log.Println("Подключение к бд прошло успешно")

	userRepo := repository.NewUserRepository(pool)
	refreshTokenRepo := repository.NewRefreshTokenRepository(pool)
	authService := service.NewAuthService(userRepo, refreshTokenRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	router.GET("/ping", pingHandler)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.PostRegister)
		authGroup.POST("/login", authHandler.PostLogin)
		authGroup.POST("/refresh", authHandler.PostRefresh)
		authGroup.POST("/logout", authHandler.PostLogout)
	}

	protectedGroup := authGroup.Group("")
	protectedGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protectedGroup.GET("/me", authHandler.GetMe)
		protectedGroup.DELETE("/me", authHandler.DeleteMe)
	}

	log.Printf("Сервер запускается на %s", cfg.ServerAddr)
	if err := router.Run(cfg.ServerAddr); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}

func pingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
