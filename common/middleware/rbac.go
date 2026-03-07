package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/bo"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/utils"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetTokenFromHeader(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized: invalid token"})
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized: invalid token"})
			c.Abort()
			return
		}
		clientId, exist := claims["client_id"].(string)
		if !exist {
			c.JSON(401, gin.H{"error": "Unauthorized: client_id not found in token"})
			c.Abort()
			return
		}
		roleKey := fmt.Sprintf(bo.ZitadelRoleKey, clientId)
		userRoles, exist := claims[roleKey].(map[string]interface{})
		if !exist {
			c.JSON(403, gin.H{"error": "Forbidden: no roles found for client"})
			c.Abort()
			return
		}
		if hasRequiredRole(userRoles, roles) {
			c.Next()
			return
		}
		c.JSON(403, gin.H{"error": "Forbidden: insufficient permissions"})
		c.Abort()
	}
}

func hasRequiredRole(userRoles map[string]interface{}, requiredRoles []string) bool {
	for _, role := range requiredRoles {
		if _, exist := userRoles[role]; exist {
			return true
		}
	}
	return false
}
