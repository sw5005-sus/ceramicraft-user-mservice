package redis

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
)

//go:generate mockery --name UserSessionDAO --output ./mocks --case underscore
type UserSessionDAO interface {
	SetSession(ctx context.Context, userSession *model.UserSession) error
	GetSession(ctx context.Context, userId int) (*model.UserSession, error)
	DelSession(ctx context.Context, userId int) error
}

type userSessionDaoImpl struct {
}

var (
	userSessionDaoInst UserSessionDAO
	userSessionOnce    sync.Once
)

func GetUserSessionDao() UserSessionDAO {
	userSessionOnce.Do(func() {
		if userSessionDaoInst == nil {
			userSessionDaoInst = &userSessionDaoImpl{}
		}
	})
	return userSessionDaoInst
}

func (dao *userSessionDaoImpl) SetSession(ctx context.Context, userSession *model.UserSession) error {
	ret := redisClient.HMSet(ctx, getUserSessionKey(userSession.UserID),
		"refresh_token", userSession.RefreshToken,
		"expires_at", userSession.ExpiresAt.Unix(),
		"id_token", userSession.IDToken,
		"access_token", userSession.AccessToken)
	if ret.Err() != nil {
		return ret.Err()
	}
	log.WithContext(ctx).Infof("Set session for user ID %d done", userSession.UserID)
	return nil
}

// GetSession implements [UserSessionDAO].
func (dao *userSessionDaoImpl) GetSession(ctx context.Context, userId int) (*model.UserSession, error) {
	ret := redisClient.HGetAll(ctx, getUserSessionKey(userId))
	if ret.Err() != nil {
		return nil, ret.Err()
	}
	data, err := ret.Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	expiresAtUnix, err := strconv.ParseInt(data["expires_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid expires_at value: %v", err)
	}
	return &model.UserSession{
		UserID:       userId,
		RefreshToken: data["refresh_token"],
		IDToken:      data["id_token"],
		AccessToken:  data["access_token"],
		ExpiresAt:    time.Unix(expiresAtUnix, 0),
	}, nil
}

// DelSession implements [UserSessionDAO].
func (dao *userSessionDaoImpl) DelSession(ctx context.Context, userId int) error {
	ret := redisClient.Del(ctx, getUserSessionKey(userId))
	if ret.Err() != nil {
		return ret.Err()
	}
	log.WithContext(ctx).Infof("Deleted session for user ID %d", userId)
	return nil
}

func getUserSessionKey(userId int) string {
	return fmt.Sprintf("us:%d", userId)
}
