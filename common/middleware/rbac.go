package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/bo"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/utils"
)

var (
	RbacProjectId = os.Getenv("RBAC_PROJECT_ID")
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetTokenFromHeader(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized: token not provided"})
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized: token parse error"})
			c.Abort()
			return
		}
		roleKey := bo.ZitadelGlobalRoleKey
		if RbacProjectId != "" {
			roleKey = fmt.Sprintf(bo.ZitadelRoleKey, RbacProjectId)
		}
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
