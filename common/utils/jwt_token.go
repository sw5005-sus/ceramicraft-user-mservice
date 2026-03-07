package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/bo"
)

var jwtSecret string

func InitJwtSecret() {
	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT secret environment variable JWT_SECRET is not set. Application cannot start.")
	}
}

// Claims structure
type Claims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

const oneDay = 24 * time.Hour

func GenerateJWTToken(user *bo.UserBO) (string, error) {
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT secret is not set")
	}
	// Create a new token object, specifying signing method and the claims
	claims := Claims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(oneDay)), // Token expiration time
			IssuedAt:  jwt.NewNumericDate(time.Now()),             // Token issued time
			NotBefore: jwt.NewNumericDate(time.Now()),             // Token valid from
		},
	}

	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWTToken(token string) (int, error) {
	if jwtSecret == "" {
		return -1, fmt.Errorf("JWT secret is not set")
	}
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return -1, err
	}

	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		return claims.ID, nil
	}

	return -1, errors.New("invalid token")
}

func GetTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
		return "", fmt.Errorf("no auth header found")
	}
	return authHeader[7:], nil
}

func ParseToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	parser := jwt.NewParser()
	_, _, err := parser.ParseUnverified(token, claims)

	if err != nil {
		return nil, fmt.Errorf("could not parse token: %v", err)
	}

	return claims, nil
}
