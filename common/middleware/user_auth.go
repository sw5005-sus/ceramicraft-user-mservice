package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := getUserIDFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if userId > 0 {
			c.Set("userID", userId)
			c.Next()
			return
		}
		userId, err = getUserIDFromCookie(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		// Set user ID into the context
		c.Set("userID", userId)

		// Token is valid, proceed to the next handler
		c.Next()
	}
}

func getUserIDFromHeader(c *gin.Context) (int, error) {
	userIdStr := c.GetHeader("X-Original-User-ID")
	if userIdStr == "" {
		return 0, nil
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in header"})
		return -1, err
	}
	return userId, nil
}

func getUserIDFromCookie(c *gin.Context) (int, error) {
	authCookie, err := c.Cookie("auth-token")
	if err != nil {
		return -1, fmt.Errorf("auth token cookie is required")
	}

	if authCookie == "" {
		return -1, fmt.Errorf("authorization is required")
	}
	ret, err := utils.ValidateJWTToken(authCookie)

	if err != nil || ret <= 0 {
		return -1, fmt.Errorf("invalid or expired token")
	}
	return ret, nil
}
