package proxy

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
)

//go:generate mockery --name ZitadelProxy --output ./mocks --case underscore
type ZitadelProxy interface {
	VerifyTokenWithBackendIdentity(ctx context.Context, accessToken string) (*model.User, error)
	SyncMeta2Zitadel(ctx context.Context, user *model.User) error
}

type zitadelProxyImpl struct {
	apiKey      *ZitadelAppKey
	mngKey      *ZitadelServiceKey
	accessToken *ZitadelAccessToken
	lock        sync.Mutex
}

var (
	zitadelProxyInstance *zitadelProxyImpl
	zitadelProxyOnce     sync.Once
)

func InitZitadel() {
	GetZitadelProxy()
	err := zitadelProxyInstance.loadKey()
	if err != nil {
		panic(err)
	}
	if zitadelProxyInstance.apiKey == nil || zitadelProxyInstance.mngKey == nil {
		panic("failed to load ZITADEL keys from environment variables")
	}
}

func GetZitadelProxy() *zitadelProxyImpl {
	zitadelProxyOnce.Do(func() {
		zitadelProxyImpl := &zitadelProxyImpl{}
		zitadelProxyInstance = zitadelProxyImpl
	})
	return zitadelProxyInstance
}

type ZitadelServiceKey struct {
	UserID string `json:"userId"`
	KeyID  string `json:"keyId"`
	Key    string `json:"key"` // RSA private key in PEM format
}

type ZitadelAppKey struct {
	KeyID    string `json:"keyId"`
	Key      string `json:"key"` // RSA private key in PEM format
	AppId    string `json:"appId"`
	ClientId string `json:"clientId"`
}

type MetadataRequest struct {
	Metadata []MetadataItem `json:"metadata"`
}

type MetadataItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ZitadelAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (t *ZitadelAccessToken) IsExpired() bool {
	return time.Now().Unix() >= t.ExpiresAt
}

func (z *zitadelProxyImpl) loadKey() error {
	appKey := &ZitadelAppKey{}
	err := loadKeyWithEnvName("ZITADEL_APP_API_KEY", appKey)
	if err != nil {
		return err
	}
	z.apiKey = appKey
	serviceKey := &ZitadelServiceKey{}
	err = loadKeyWithEnvName("ZITADEL_SERVICE_API_KEY", serviceKey)
	if err != nil {
		return err
	}
	z.mngKey = serviceKey
	return nil
}

func loadKeyWithEnvName(envName string, apiKey any) error {
	keyData := os.Getenv(envName)
	if keyData == "" {
		return fmt.Errorf("environment variable %s is not set", envName)
	}
	err := json.Unmarshal([]byte(keyData), apiKey)
	return err
}

func (z *zitadelProxyImpl) VerifyTokenWithBackendIdentity(ctx context.Context, accessToken string) (*model.User, error) {
	introspectURL := config.Config.ZitadelConfig.Host + "/oauth/v2/introspect"

	// gen signed JWT
	assertionToken, err := generateAssersionToken(z.apiKey.ClientId, z.apiKey.Key, z.apiKey.KeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}

	form := url.Values{}
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", assertionToken)
	form.Add("token", accessToken)

	// request ZITADEL
	resp, err := http.Post(introspectURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Logger.Errorf("failed to close response body: %v", err)
		}
	}()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("%d Error Detail: %s\n", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("failed to introspect token, status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	active, ok := result["active"].(bool)
	if !ok || !active {
		return nil, fmt.Errorf("invalid token: token is not active")
	}
	sub, ok := result["sub"].(string)
	if !ok || sub == "" {
		return nil, fmt.Errorf("invalid token: sub claim is missing or not a string")
	}
	email, ok := result["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("invalid token: email claim is missing or not a string")
	}
	return &model.User{
		ZitadelSub: sub,
		Email:      email,
	}, nil
}

func (z *zitadelProxyImpl) SyncMeta2Zitadel(ctx context.Context, user *model.User) error {
	escapedSub := url.PathEscape(user.ZitadelSub)
	apiURL := fmt.Sprintf("%s/v2/users/%s/metadata", config.Config.ZitadelConfig.Host, escapedSub)

	payload, _ := json.Marshal(MetadataRequest{
		Metadata: []MetadataItem{
			{
				Key:   "local_userid",
				Value: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", user.ID))),
			},
		},
	})

	accessToken, err := z.getActualAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %v", err)
	}
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				log.Logger.Errorf("failed to close response body: %v", err)
			}
		}
	}()
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Logger.Errorf("%d Error Detail: %s\n", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("failed to sync metadata, status code: %d", resp.StatusCode)
	}
	return nil
}

func (z *zitadelProxyImpl) getActualAccessToken() (string, error) {
	if z.accessToken != nil && !z.accessToken.IsExpired() {
		return z.accessToken.AccessToken, nil
	}
	z.lock.Lock()
	defer z.lock.Unlock()
	if z.accessToken != nil && !z.accessToken.IsExpired() {
		return z.accessToken.AccessToken, nil
	}
	assertionToken, err := generateAssersionToken(z.mngKey.UserID, z.mngKey.Key, z.mngKey.KeyID)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	data.Set("assertion", assertionToken) // 这就是你签名的那个 JWT
	data.Set("scope", "openid profile urn:zitadel:iam:org:project:id:zitadel:aud")

	resp, err := http.PostForm(fmt.Sprintf("%s/oauth/v2/token", config.Config.ZitadelConfig.Host), data)
	defer func() {
		if resp != nil && resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				log.Logger.Errorf("failed to close response body: %v", err)
			}
		}
	}()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to get access token, status code: %d, detail: %s\n", resp.StatusCode, string(bodyBytes))
		return "", fmt.Errorf("failed to get access token, status code: %d", resp.StatusCode)
	}

	// 2. 解析返回的 Access Token
	var result ZitadelAccessToken
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if result.AccessToken != "" {
		z.accessToken = &result
		z.accessToken.ExpiresAt = time.Now().Unix() + int64(result.ExpiresIn) - 10 // 10s before actual expiration to avoid edge case
	}
	return result.AccessToken, nil
}

func generateAssersionToken(sub, key, keyId string) (string, error) {
	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"iss": sub,
		"sub": sub,
		"aud": config.Config.ZitadelConfig.Host,
		"exp": now + 3600,
		"iat": now,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyId

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	return token.SignedString(signKey)
}
