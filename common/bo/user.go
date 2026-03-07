package bo

const (
	OAuthHeaderUserId = "X-Original-User-ID"
)

type UserBO struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
