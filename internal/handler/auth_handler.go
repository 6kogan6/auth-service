package handler

import (
	"auth-service/internal/dto"
	"auth-service/internal/service"
	"errors"
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

func (h *AuthHandler) PostRegister(c *gin.Context) {
	var newUserRequest dto.RegisterRequest

	if err := c.ShouldBindJSON(&newUserRequest); err != nil {
		sendError(c, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	newUserResponse, err := h.authService.RegisterNewUser(newUserRequest)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			sendError(c, http.StatusConflict, err.Error())
			return
		}

		if (errors.Is(err, service.ErrEmailCheckFailed)) || (errors.Is(err, service.ErrUserCreateFailed)) || (errors.Is(err, service.ErrPasswordHashFailed)) {
			sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}

		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, newUserResponse)
}

func (h *AuthHandler) PostLogin(c *gin.Context) {
	var loginRequest dto.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		sendError(c, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	loginResponse, err := h.authService.Login(loginRequest)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			sendError(c, http.StatusUnauthorized, err.Error())
			return
		}

		if (errors.Is(err, service.ErrTokenCreateFailed)) || (errors.Is(err, service.ErrRefreshTokenCreateFailed)) || (errors.Is(err, service.ErrUserCheckFailed)) {
			sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}

		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, loginResponse)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, err := validateUserID(c)
	if err != nil {
		if errors.Is(err, ErrUserNotFoundContext) {
			sendError(c, http.StatusUnauthorized, err.Error())
			return
		}

		if errors.Is(err, ErrUserBadType) {
			sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}

		sendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	meResponse, err := h.authService.GetMe(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			sendError(c, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, service.ErrUserCheckFailed) {
			sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}

		sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	c.IndentedJSON(http.StatusOK, meResponse)
}

func (h *AuthHandler) PostRefresh(c *gin.Context) {
	var refreshRequest dto.RefreshRequest
	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		sendError(c, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	refreshResponse, err := h.authService.Refresh(refreshRequest)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			sendError(c, http.StatusUnauthorized, err.Error())
			return
		}

		if (errors.Is(err, service.ErrTokenCreateFailed)) || (errors.Is(err, service.ErrRefreshTokenCheckFailed)) || (errors.Is(err, service.ErrUserCheckFailed)) {
			sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}

		sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}
	c.IndentedJSON(http.StatusOK, refreshResponse)
}

func (h *AuthHandler) PostLogout(c *gin.Context) {
	var logoutRequest dto.LogoutRequest
	if err := c.ShouldBindJSON(&logoutRequest); err != nil {
		sendError(c, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	logoutResponse, err := h.authService.Logout(logoutRequest)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			sendError(c, http.StatusUnauthorized, err.Error())
			return
		}

		if errors.Is(err, service.ErrRefreshTokenRevokeFailed) {
			sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}

		sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	c.IndentedJSON(http.StatusOK, logoutResponse)
}

func (h *AuthHandler) DeleteMe(c *gin.Context) {
	userID, err := validateUserID(c)
	if err != nil {
		if errors.Is(err, ErrUserNotFoundContext) {
			sendError(c, http.StatusUnauthorized, err.Error())
			return
		}

		if errors.Is(err, ErrUserBadType) {
			sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}

		sendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	deleteAccountResponse, err := h.authService.DeactivateAccount(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			sendError(c, http.StatusNotFound, err.Error())
			return
		}

		if (errors.Is(err, service.ErrRefreshTokenRevokeFailed)) || (errors.Is(err, service.ErrUserDeactivateFailed)) {
			sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}

		sendError(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	c.IndentedJSON(http.StatusOK, deleteAccountResponse)
}
