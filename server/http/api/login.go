package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sw5005-sus/ceramicraft-user-mservice/common/bo"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/proxy"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/redis"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/service"

	"github.com/gin-gonic/gin"
)

const tokenExpireDuration = 3600 * 24 * 365 // 1 year

// OAuthLogin initiates the admin OAuth login flow.
//
// @Summary Admin OAuthLogin
// @Description Initiates the OAuth login process for admin users by generating a state, setting it in a cookie, and redirecting to the OAuth provider's authorization URL.
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 302 "Redirects to OAuth provider for authentication"
// @Router /user-ms/v1/merchant/oauth-login [get]
func OAuthLogin(c *gin.Context) {
	state, err := generateState(16)
	if err != nil {
		log.Logger.Errorf("Failed to generate state: %v", err)
		c.JSON(http.StatusInternalServerError, data.BaseResponse{ErrMsg: "Failed to generate state"})
		return
	}
	c.SetCookie("oauth_state", state, 300, "/", getCookieDomain(c.Request.Host), false, true)
	authURL := proxy.GetZitadelProxy().GetAuthCodeURL(state)
	c.Redirect(http.StatusFound, authURL)
}

func generateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

// AdminLoginCallback handles admin oauth login callback.
//
// @Summary Admin OAuthLoginCallback
// @Description Handles the OAuth callback for admin login, exchanges the code for a token, and sets it in a cookie.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from OAuth provider"
// @Param state query string true "State parameter for CSRF protection"
// @Success 200	{object} data.BaseResponse{data=string} "Login successful, returns auth token in cookie"
// @Failure 500 {object} data.BaseResponse{data=string}
// @Router /user-ms/v1/merchant/login-callback [get]
func AdminLoginCallback(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, data.BaseResponse{ErrMsg: "State parameter is required"})
		return
	}
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || cookieState != state {
		c.JSON(http.StatusBadRequest, data.BaseResponse{ErrMsg: "Invalid state parameter"})
		return
	}
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, data.BaseResponse{ErrMsg: "Code parameter is required"})
		return
	}
	token, err := service.GetOAuthService().ZitadelCallback(c.Request.Context(), code)
	if err != nil {
		log.Logger.Errorf("OAuth callback error: %v", err)
		c.JSON(http.StatusInternalServerError, data.BaseResponse{ErrMsg: "Failed to authenticate with Zitadel"})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     bo.OAuthTokenCookieName,
		Value:    token,
		Domain:   getCookieDomain(c.Request.Host),
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	c.Redirect(http.StatusFound, config.Config.SystemConfig.HomeUrl)
}

// UserLogin handles user login requests.
//
// @Summary User Login
// @Description Authenticates a user with their email and password and returns a token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body data.UserLoginVO true "User login information"
// @Param client path string true "Client identifier" Enums(customer, merchant)
// @Success 200	{object} data.BaseResponse{data=string} "Login successful, returns auth token in cookie"
// @Failure 400 {object} data.BaseResponse{data=string}
// @Failure 500 {object} data.BaseResponse{data=string}
// @Router /user-ms/v1/{client}/login [post]
func UserLogin(c *gin.Context) {
	user := &data.UserLoginVO{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, data.BaseResponse{ErrMsg: err.Error()})
		return
	}
	token, err := service.GetLoginService().Login(c.Request.Context(), user.Email, user.Password)
	if err != nil {
		log.Logger.Errorf("Login error: %v", err)
		c.JSON(http.StatusInternalServerError, data.BaseResponse{ErrMsg: err.Error()})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth-token",
		Value:    token,
		Path:     "/",
		Domain:   c.Request.Host,
		Expires:  time.Now().Add(time.Duration(tokenExpireDuration) * time.Second),
		Secure:   false,
		HttpOnly: true,
	})
	c.JSON(http.StatusOK, data.BaseResponse{Data: "Login successful"})
}

// UserLogout handles user logout requests.
//
// @Summary User Logout
// @Description invalidates the user's auth token cookie.
// @Tags Authentication
// @Param client path string true "Client identifier" Enums(customer, merchant)
// @Success 200 object data.BaseResponse{data=string} "Logout successful"
// @Router /user-ms/v1/{client}/logout [post]
func UserLogout(c *gin.Context) {
	// Invalidate the auth-token cookie by setting its MaxAge to -1
	c.SetCookie("auth-token", "", -1, "/", c.Request.Host, true, true)
	c.JSON(http.StatusOK, data.BaseResponse{Data: "Logout successful"})
}

// OAuthLogout handles admin logout requests.
// @Summary OAuth Logout
// @Description Logs out the user from the OAuth provider and invalidates the session.
// @Tags Authentication
// @Success 302 "Redirects to OAuth provider logout endpoint"
// @Router /user-ms/v1/merchant/oauth-logout [get]
func OAuthLogout(c *gin.Context) {
	userId := c.GetInt("userID")
	cookieDomain := getCookieDomain(c.Request.Host)
	// Invalidate the auth-token cookie by setting its MaxAge to -1
	c.SetCookie(bo.OAuthTokenCookieName, "", -1, "/", cookieDomain, true, true)
	c.SetCookie("oauth_state", "", -1, "/", cookieDomain, true, true)
	zitadelLogoutURL := config.Config.ZitadelConfig.Host + "/oidc/v1/end_session"

	postLogoutRedirectURI := config.Config.SystemConfig.HomeUrl
	idToken, _ := c.Cookie(bo.OAuthTokenCookieName)
	finalURL := fmt.Sprintf("%s?id_token_hint=%s&post_logout_redirect_uri=%s",
		zitadelLogoutURL, url.QueryEscape(idToken), url.QueryEscape(postLogoutRedirectURI))
	err := redis.GetUserSessionDao().DelSession(c.Request.Context(), userId)
	if err != nil {
		log.Logger.Warnf("Failed to delete user session for user ID %d: %v", userId, err)
	}
	c.Redirect(http.StatusFound, finalURL)
}

func getCookieDomain(host string) string {
	// Extract the domain from the host (remove port if present)
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	// Check if the host is an IP address
	if net.ParseIP(host) != nil {
		return host // Return the IP address as is
	}

	// Split the host into parts
	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		// Return the last two parts (b.c)
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return host // Return the host if it's already a top-level domain
}
