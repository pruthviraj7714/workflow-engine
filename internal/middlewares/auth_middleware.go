package middlewares

import (
	"strings"
	"workflow-engine/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		authToken := strings.Split(authHeader, " ")[1]

		if authToken == "" {
			c.JSON(401, gin.H{"error": "Authorization token is missing"})
			c.Abort()
			return
		}

		claims, err := utils.VerifyToken(authToken)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.ID.String())

		c.Next()
	}
}
