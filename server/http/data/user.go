package data

import "strings"

type UserLoginVO struct {
	ID       int    `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

type UserActivateReq struct {
	Code string `json:"code" binding:"required,min=6,max=6"`
}

type UserProfileVO struct {
	ID             int            `json:"id"`
	Email          string         `json:"email"`
	Name           string         `json:"name" binding:"max=50"`
	Avatar         string         `json:"avatar"`
	DefaultAddress *UserAddressVO `json:"default_address,omitempty"`
}

type UserAddressVO struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	ZipCode      string `json:"zip_code" binding:"required,alphanum,max=10"`
	Country      string `json:"country" binding:"required,max=50"`
	Province     string `json:"province" binding:"max=50"`
	City         string `json:"city" binding:"max=50"`
	Detail       string `json:"detail" binding:"required,max=255"`
	FirstName    string `json:"first_name" binding:"required,max=50"`
	LastName     string `json:"last_name" binding:"required,max=50"`
	ContactPhone string `json:"contact_phone" binding:"required,e164"`
	IsDefault    bool   `json:"is_default"`
}

func MaskUserProfile(profile *UserProfileVO) {
	if profile == nil {
		return
	}
	profile.Email = maskEmail(profile.Email)
}

func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	username := parts[0]
	runes := []rune(username)
	if len(runes) <= 1 {
		return "*@" + parts[1]
	}
	return string(runes[0]) + "****@" + parts[1]
}
