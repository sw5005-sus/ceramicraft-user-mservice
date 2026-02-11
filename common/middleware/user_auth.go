package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authCookie, err := c.Cookie("auth-token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth token cookie is required"})
			c.Abort()
			return
		}

		if authCookie == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization is required"})
			c.Abort()
			return
		}
		ret, err := utils.ValidateJWTToken(authCookie)

		if err != nil || ret <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		// Set user ID into the context
		c.Set("userID", ret)

		// Token is valid, proceed to the next handler
		c.Next()
	}
}
