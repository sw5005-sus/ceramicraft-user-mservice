package model

import "time"

type UserSession struct {
	UserID       int       `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshToken string    `json:"refresh_token"`
	IDToken      string    `json:"id_token"`
}
