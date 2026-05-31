package main

import (
	"auth-service/internal/db"
	handler "auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	connStr := "postgres://admin:admin@localhost:5432/authdb"
	pool, err := db.ConnectPostgres(connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к бд: %v", err)
	}

	defer pool.Close()

	repo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(repo)
	userHandler := handler.NewAuthHandler(authService)

	log.Println("Подключение к бд прошло успешно")

	router.GET("/ping", checkServerWork)
	router.POST("/api/v1/auth/register", userHandler.PostNewUser)

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Не удалось запустить сервер %v", err)
	}

}

func checkServerWork(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
