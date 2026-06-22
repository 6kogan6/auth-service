package middleware

import (
	"auth-service/internal/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if (authHeader == "") || !(strings.HasPrefix(authHeader, "Bearer ")) {
			c.IndentedJSON(http.StatusUnauthorized,
				gin.H{"error": "Токен отсутствует или неверный формат"})
			c.Abort()
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := token.ParseAccessToken(accessToken, jwtSecret)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized,
				gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
