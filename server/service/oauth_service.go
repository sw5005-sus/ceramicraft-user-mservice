package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/sw5005-sus/ceramicraft-user-mservice/common/bo"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/proxy"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/redis"
)

type OAuthService interface {
	ZitadelCallback(ctx context.Context, code string) (string, error)
	CustomerVerify(ctx context.Context, token string) (headerSet map[string]string, cookieSet map[string]string, err error)
	MerchantVerify(ctx context.Context, token string) (headerSet map[string]string, cookieSet map[string]string, err error)
}

var (
	oAuthServiceInst OAuthService
	oAuthSyncOnce    sync.Once
)

func GetOAuthService() OAuthService {
	oAuthSyncOnce.Do(func() {
		oAuthServiceInst = &oAuthServiceImpl{
			zitadelProxy:   proxy.GetZitadelProxy(),
			userSessionDao: redis.GetUserSessionDao(),
		}
	})
	return oAuthServiceInst
}

type oAuthServiceImpl struct {
	zitadelProxy   proxy.ZitadelProxy
	userSessionDao redis.UserSessionDAO
}

// CustomerVerify implements [OAuthService].
func (o *oAuthServiceImpl) CustomerVerify(ctx context.Context, token string) (headerSet map[string]string, cookieSet map[string]string, err error) {
	authUser, err := proxy.GetZitadelProxy().ValidateToken(ctx, token, proxy.APP_ZITADEL_USER_KEY, config.Config.ZitadelConfig.AppClientId)
	if err != nil {
		return nil, nil, err
	}
	headerSet = map[string]string{
		bo.OAuthHeaderUserId: fmt.Sprintf("%d", authUser.LocalUserId),
	}
	return headerSet, map[string]string{}, nil
}

// MerchantVerify implements [OAuthService].
func (o *oAuthServiceImpl) MerchantVerify(ctx context.Context, token string) (headerSet map[string]string, cookieSet map[string]string, err error) {
	authUser, err := proxy.GetZitadelProxy().ValidateToken(ctx, token, proxy.ADMIN_ZITADEL_USER_KEY, config.Config.ZitadelConfig.AdminClientId)
	if authUser == nil {
		log.Logger.Infof("Merchant token validation failed: %v", err)
		return nil, nil, err
	}
	if err == nil {
		return map[string]string{
			bo.OAuthHeaderUserId: fmt.Sprintf("%d", authUser.LocalUserId),
		}, map[string]string{}, err
	}
	// If token is invalid, try to refresh the session and validate again
	userSession, err := o.userSessionDao.GetSession(ctx, authUser.LocalUserId)
	if err != nil {
		log.Logger.Errorf("Failed to get user session from Redis: %v", err)
		return nil, nil, err
	}
	newSession, err := o.zitadelProxy.RefreshUserSession(ctx, userSession.RefreshToken)
	if err != nil {
		log.Logger.Errorf("Failed to refresh user session: %v", err)
		return nil, nil, err
	}
	err = o.userSessionDao.SetSession(ctx, newSession)
	if err != nil {
		return nil, nil, err
	}
	return map[string]string{
			bo.OAuthHeaderUserId: fmt.Sprintf("%d", authUser.LocalUserId),
		}, map[string]string{
			"oauth_token": newSession.IDToken,
		}, nil
}

// ZitadelCallback implements [OAuthService].
func (o *oAuthServiceImpl) ZitadelCallback(ctx context.Context, code string) (string, error) {
	userSession, err := o.zitadelProxy.AuthCallback(ctx, code)
	if err != nil {
		log.Logger.Errorf("Failed to handle Zitadel callback: %v", err)
		return "", err
	}
	err = o.userSessionDao.SetSession(ctx, userSession)
	if err != nil {
		log.Logger.Errorf("Failed to set user session in Redis: %v", err)
		return "", err
	}
	return userSession.IDToken, nil
}
