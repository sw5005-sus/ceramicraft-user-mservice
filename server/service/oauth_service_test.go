package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	proxy_mock "github.com/sw5005-sus/ceramicraft-user-mservice/server/proxy/mocks"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
	redis_mock "github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/redis/mocks"
)

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
