package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/proxy"
)

// OAuthTokenValidate validates jwt_token.
//
// @Summary Validate OAuth Token
// @Description This endpoint validates the provided JWT token and returns user information if the token is valid.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param        Authorization  header    string  true  "Insert your access token with 'Bearer ' prefix"
// @Success 200
// @Router /oauth/v1/verify [get]
func OAuthTokenValidate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
		log.Logger.Infof("no auth header found")
		return
	}
	token := authHeader[7:]
	authUser, err := proxy.GetZitadelProxy().ValidateToken(c.Request.Context(), token)
	if err != nil {
		log.Logger.Errorf("token validation failed: %v", err)
		return
	}

	c.Writer.Header().Set("X-Original-User-ID", fmt.Sprint(authUser.LocalUserId))
}
