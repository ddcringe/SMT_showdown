package auth

import (
	"net/http"
	"strings"

	"github.com/ddcringe/SMT_showdown/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT токен в заголовке Authorization
func AuthMiddleware(jwtUtil *jwt.TokenUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		claims, err := jwtUtil.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
