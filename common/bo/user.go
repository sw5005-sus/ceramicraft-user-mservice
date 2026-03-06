package bo

const (
	OAuthHeaderUserId    = "X-Original-User-ID"
	OAuthHeaderTimestamp = "X-Original-Timestamp"
	OAuthHeaderSign      = "X-Original-Sign"
)

type UserBO struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
