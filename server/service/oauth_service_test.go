package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/bo"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/proxy"
	proxy_mock "github.com/sw5005-sus/ceramicraft-user-mservice/server/proxy/mocks"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
	redis_mock "github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/redis/mocks"
)

func TestCustomVerify(t *testing.T) {
	initEnv()
	ctx := context.Background()
	t.Run("successful CustomerVerify", func(t *testing.T) {
		mockZitadelProxy := new(proxy_mock.ZitadelProxy)
		mockSessionDao := new(redis_mock.UserSessionDAO)

		service := &oAuthServiceImpl{
			zitadelProxy:   mockZitadelProxy,
			userSessionDao: mockSessionDao,
		}
		token := "valid_token"
		mockZitadelProxy.On("ValidateToken", ctx, token, proxy.APP_ZITADEL_USER_KEY, config.Config.ZitadelConfig.AppClientId).Return(&proxy.AuthUser{LocalUserId: 123}, nil)

		headerSet, cookieSet, err := service.CustomerVerify(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{bo.OAuthHeaderUserId: "123"}, headerSet)
		assert.Empty(t, cookieSet)

		mockZitadelProxy.AssertExpectations(t)
	})

	t.Run("failed ValidateToken", func(t *testing.T) {
		mockZitadelProxy := new(proxy_mock.ZitadelProxy)
		mockSessionDao := new(redis_mock.UserSessionDAO)

		service := &oAuthServiceImpl{
			zitadelProxy:   mockZitadelProxy,
			userSessionDao: mockSessionDao,
		}
		token := "invalid_token"
		mockZitadelProxy.On("ValidateToken", ctx, token, proxy.APP_ZITADEL_USER_KEY, config.Config.ZitadelConfig.AppClientId).Return(nil, errors.New("validation error"))

		headerSet, cookieSet, err := service.CustomerVerify(ctx, token)
		assert.Error(t, err)
		assert.Nil(t, headerSet)
		assert.Nil(t, cookieSet)

		mockZitadelProxy.AssertExpectations(t)
	})
}

func TestZitadelCallback(t *testing.T) {
	initEnv()
	ctx := context.Background()
	code := "test_code"

	t.Run("successful callback", func(t *testing.T) {
		mockZitadelProxy := new(proxy_mock.ZitadelProxy)
		mockSessionDao := new(redis_mock.UserSessionDAO)

		service := &oAuthServiceImpl{
			zitadelProxy:   mockZitadelProxy,
			userSessionDao: mockSessionDao,
		}
		mockZitadelProxy.On("AuthCallback", ctx, code).Return(&model.UserSession{IDToken: "test_token"}, nil)
		mockSessionDao.On("SetSession", ctx, mock.Anything).Return(nil)

		token, err := service.ZitadelCallback(ctx, code)
		assert.NoError(t, err)
		assert.Equal(t, "test_token", token)

		mockZitadelProxy.AssertExpectations(t)
		mockSessionDao.AssertExpectations(t)
	})

	t.Run("failed AuthCallback", func(t *testing.T) {
		mockZitadelProxy := new(proxy_mock.ZitadelProxy)
		mockSessionDao := new(redis_mock.UserSessionDAO)

		service := &oAuthServiceImpl{
			zitadelProxy:   mockZitadelProxy,
			userSessionDao: mockSessionDao,
		}
		mockZitadelProxy.On("AuthCallback", ctx, code).Return(nil, errors.New("auth error"))

		token, err := service.ZitadelCallback(ctx, code)
		assert.Error(t, err)
		assert.Empty(t, token)

		mockZitadelProxy.AssertExpectations(t)
	})

	t.Run("failed SetSession", func(t *testing.T) {
		mockZitadelProxy := new(proxy_mock.ZitadelProxy)
		mockSessionDao := new(redis_mock.UserSessionDAO)

		service := &oAuthServiceImpl{
			zitadelProxy:   mockZitadelProxy,
			userSessionDao: mockSessionDao,
		}
		mockZitadelProxy.On("AuthCallback", ctx, code).Return(&model.UserSession{AccessToken: "test_token"}, nil)
		mockSessionDao.On("SetSession", ctx, mock.Anything).Return(errors.New("session error"))

		token, err := service.ZitadelCallback(ctx, code)
		assert.Error(t, err)
		assert.Empty(t, token)

		mockZitadelProxy.AssertExpectations(t)
		mockSessionDao.AssertExpectations(t)
	})
}
func TestMerchantVerify(t *testing.T) {
	initEnv()
	ctx := context.Background()

	t.Run("successful MerchantVerify", func(t *testing.T) {
		mockZitadelProxy := new(proxy_mock.ZitadelProxy)
		mockSessionDao := new(redis_mock.UserSessionDAO)

		service := &oAuthServiceImpl{
			zitadelProxy:   mockZitadelProxy,
			userSessionDao: mockSessionDao,
		}
		token := "valid_merchant_token"
		mockZitadelProxy.On("ValidateToken", ctx, token, proxy.ADMIN_ZITADEL_USER_KEY, config.Config.ZitadelConfig.AdminClientId).Return(&proxy.AuthUser{LocalUserId: 456}, nil)

		headerSet, cookieSet, err := service.MerchantVerify(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{bo.OAuthHeaderUserId: "456"}, headerSet)
		assert.Empty(t, cookieSet)

		mockZitadelProxy.AssertExpectations(t)
	})

	t.Run("failed ValidateToken", func(t *testing.T) {
		mockZitadelProxy := new(proxy_mock.ZitadelProxy)
		mockSessionDao := new(redis_mock.UserSessionDAO)

		service := &oAuthServiceImpl{
			zitadelProxy:   mockZitadelProxy,
			userSessionDao: mockSessionDao,
		}
		token := "invalid_merchant_token"
		mockZitadelProxy.On("ValidateToken", ctx, token, proxy.ADMIN_ZITADEL_USER_KEY, config.Config.ZitadelConfig.AdminClientId).Return(nil, errors.New("validation error"))

		headerSet, cookieSet, err := service.MerchantVerify(ctx, token)
		assert.Error(t, err)
		assert.Nil(t, headerSet)
		assert.Nil(t, cookieSet)

		mockZitadelProxy.AssertExpectations(t)
	})

	t.Run("failed GetSession while token expired", func(t *testing.T) {
		mockZitadelProxy := new(proxy_mock.ZitadelProxy)
		mockSessionDao := new(redis_mock.UserSessionDAO)

		service := &oAuthServiceImpl{
			zitadelProxy:   mockZitadelProxy,
			userSessionDao: mockSessionDao,
		}
		token := "invalid_merchant_token"
		mockZitadelProxy.On("ValidateToken", ctx, token,
			proxy.ADMIN_ZITADEL_USER_KEY,
			config.Config.ZitadelConfig.AdminClientId).Return(&proxy.AuthUser{LocalUserId: 1}, fmt.Errorf("token expired"))
		mockSessionDao.On("GetSession", ctx, mock.Anything).Return(nil, errors.New("session not found"))

		headerSet, cookieSet, err := service.MerchantVerify(ctx, token)
		assert.Error(t, err)
		assert.Nil(t, headerSet)
		assert.Nil(t, cookieSet)

		mockZitadelProxy.AssertExpectations(t)
		mockSessionDao.AssertExpectations(t)
	})

	t.Run("failed RefreshUserSession", func(t *testing.T) {
		mockZitadelProxy := new(proxy_mock.ZitadelProxy)
		mockSessionDao := new(redis_mock.UserSessionDAO)

		service := &oAuthServiceImpl{
			zitadelProxy:   mockZitadelProxy,
			userSessionDao: mockSessionDao,
		}
		token := "invalid_merchant_token"
		mockZitadelProxy.On("ValidateToken", ctx, token, proxy.ADMIN_ZITADEL_USER_KEY, config.Config.ZitadelConfig.AdminClientId).Return(&proxy.AuthUser{LocalUserId: 1}, errors.New("validation error"))
		mockSessionDao.On("GetSession", ctx, mock.Anything).Return(&model.UserSession{RefreshToken: "refresh_token"}, nil)
		mockZitadelProxy.On("RefreshUserSession", ctx, "refresh_token").Return(nil, errors.New("refresh error"))

		headerSet, cookieSet, err := service.MerchantVerify(ctx, token)
		assert.Error(t, err)
		assert.Nil(t, headerSet)
		assert.Nil(t, cookieSet)

		mockZitadelProxy.AssertExpectations(t)
		mockSessionDao.AssertExpectations(t)
	})

	t.Run("successful session refresh", func(t *testing.T) {
		mockZitadelProxy := new(proxy_mock.ZitadelProxy)
		mockSessionDao := new(redis_mock.UserSessionDAO)

		service := &oAuthServiceImpl{
			zitadelProxy:   mockZitadelProxy,
			userSessionDao: mockSessionDao,
		}
		token := "invalid_merchant_token"
		userId := 1
		mockZitadelProxy.On("ValidateToken", ctx, token, proxy.ADMIN_ZITADEL_USER_KEY, config.Config.ZitadelConfig.AdminClientId).Return(&proxy.AuthUser{LocalUserId: userId}, errors.New("validation error"))
		mockSessionDao.On("GetSession", ctx, mock.Anything).Return(&model.UserSession{RefreshToken: "refresh_token"}, nil)
		mockZitadelProxy.On("RefreshUserSession", ctx, "refresh_token").Return(&model.UserSession{IDToken: "new_token"}, nil)
		mockSessionDao.On("SetSession", ctx, mock.Anything).Return(nil)

		headerSet, cookieSet, err := service.MerchantVerify(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{bo.OAuthHeaderUserId: fmt.Sprintf("%d", userId)}, headerSet) // LocalUserId is not set in this case
		assert.Equal(t, map[string]string{bo.OAuthTokenCookieName: "new_token"}, cookieSet)

		mockZitadelProxy.AssertExpectations(t)
		mockSessionDao.AssertExpectations(t)
	})
}
