package handler

import (
	dto "auth-service/internal/dto"
	"auth-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) PostNewUser(c *gin.Context) {
	var newUserRequest dto.RegisterRequest

	if err := c.ShouldBindJSON(&newUserRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest,
			gin.H{
				"error": "Неверный формат данных",
			})
		return
	}

	newUserResponse, err := h.authService.RegisterNewUser(newUserRequest)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			})
		return
	}

	c.IndentedJSON(http.StatusCreated, newUserResponse)
}
