package bo

const (
	OAuthHeaderUserId  = "X-Original-User-ID"
	ZitadelRoleKey     = "urn:zitadel:iam:org:project:%s:roles"
	ZitadelMetaDataKey = "urn:zitadel:iam:user:metadata"
)

type UserBO struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
