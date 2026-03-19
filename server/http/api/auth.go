package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/bo"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/service"
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
	client := getClient(c)
	if client == "customer" {
		token := readTokenFromHeader(c)
		if token == "" {
			return
		}
		headers, _, err := service.GetOAuthService().CustomerVerify(c.Request.Context(), token)
		if err != nil {
			log.Logger.Errorf("Customer token validation failed: %v", err)
			return
		}
		setHeaders(c, headers)
	} else {
		token := readTokenFromCookie(c)
		if token == "" {
			return
		}
		headers, cookies, err := service.GetOAuthService().MerchantVerify(c.Request.Context(), token)
		if err != nil {
			log.Logger.Errorf("Merchant token validation failed: %v", err)
			return
		}
		setHeaders(c, headers)
		setCookies(c, cookies)
	}

	c.Status(http.StatusOK)
}

func setHeaders(c *gin.Context, headers map[string]string) {
	for k, v := range headers {
		c.Writer.Header().Set(k, v)
	}
}

func setCookies(c *gin.Context, cookies map[string]string) {
	for k, v := range cookies {
		c.SetCookie(k, v, tokenExpireDuration, "/", getCookieDomain(c.Request.Host), false, true)
	}
}

func readTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
		log.Logger.Infof("no auth header found")
		return ""
	}
	return authHeader[7:]
}

func readTokenFromCookie(c *gin.Context) string {
	token, err := c.Cookie(bo.OAuthTokenCookieName)
	if err != nil {
		log.Logger.Infof("no auth cookie found")
		return ""
	}
	return token
}

func getClient(c *gin.Context) string {
	requestUri := c.GetHeader("X-Forwarded-Uri")
	if strings.Contains(requestUri, "/customer/") {
		return "customer"
	} else {
		return "merchant"
	}
}
