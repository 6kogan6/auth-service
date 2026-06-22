package handler

import (
	"github.com/gin-gonic/gin"
)

func sendError(c *gin.Context, statusCode int, textError string) {
	c.IndentedJSON(statusCode,
		gin.H{
			"error": textError,
		})
}
