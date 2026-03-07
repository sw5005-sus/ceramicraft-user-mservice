package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/bo"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/proxy"
)

// OAuthTokenValidate validates jwt_token.
//
// @Summary Validate OAuth Token
// @Description This endpoint validates the provided JWT token. If the token is valid, it sets response headers.
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

	setHeaders(c, authUser.LocalUserId)

	c.Status(http.StatusOK)
}

func setHeaders(c *gin.Context, userId int) {
	c.Writer.Header().Set(bo.OAuthHeaderUserId, fmt.Sprint(userId))
}
