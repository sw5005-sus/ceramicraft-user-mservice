package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	ceramicraftsecure "github.com/sw5005-sus/ceramicraft-secure"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/bo"
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
	userIdStr := c.GetHeader(bo.OAuthHeaderUserId)
	if userIdStr == "" {
		return 0, nil
	}
	if !validateSign(userIdStr, c.GetHeader(bo.OAuthHeaderTimestamp), c.GetHeader(bo.OAuthHeaderSign)) {
		return -1, fmt.Errorf("invalid signature")
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId <= 0 {
		return -1, err
	}
	return userId, nil
}

func validateSign(userId, timestamp, sign string) bool {
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil || time.Now().Unix()-ts > 60 { // 60 seconds time window
		return false
	}
	ret, err := ceramicraftsecure.VerifyHmacSha256(fmt.Sprintf("%s:%s", userId, timestamp), sign)
	if err != nil {
		return false
	}
	return ret
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
