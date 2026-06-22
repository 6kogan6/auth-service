package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrUserNotFoundContext = errors.New("user_id не найден в контексте")
var ErrUserBadType = errors.New("user_id имеет неверный тип")

func validateUserID(c *gin.Context) (uint64, error) {

	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, ErrUserNotFoundContext
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		return 0, ErrUserBadType
	}
	return userID, nil
}
