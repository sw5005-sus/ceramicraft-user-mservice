package service

import (
	"context"
	"database/sql"
	"sync"

	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/dao"
)

type UserProfileService interface {
	GetUserProfile(ctx context.Context, userID int) (*data.UserProfileVO, error)
	UpdateUserProfile(ctx context.Context, userID int, profile *data.UserProfileVO) error
}

var (
	userProfileServiceInst *UserProfileServiceImpl
	userProfileOnce        sync.Once
)

func GetUserProfileService() *UserProfileServiceImpl {
	userProfileOnce.Do(func() {
		userProfileServiceInst = &UserProfileServiceImpl{
			userDao: dao.GetUserDao(),
		}
	})
	return userProfileServiceInst
}

type UserProfileServiceImpl struct {
	userDao dao.UserDao
}

func (u *UserProfileServiceImpl) GetUserProfile(ctx context.Context, userID int) (*data.UserProfileVO, error) {
	user, err := u.userDao.GetUserById(ctx, userID)
	if err != nil {
		log.Logger.Errorf("Failed to get user by id: %v", err)
		return nil, err
	}
	if user == nil {
		log.Logger.Warnf("User not found with id: %d", userID)
		return nil, nil
	}
	return &data.UserProfileVO{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Avatar: user.AvatarId,
	}, nil
}

func (u *UserProfileServiceImpl) UpdateUserProfile(ctx context.Context, userID int, profile *data.UserProfileVO) error {
	user, err := u.userDao.GetUserById(ctx, userID)
	if err != nil {
		log.Logger.Errorf("Failed to get user by id: %v", err)
		return err
	}
	if user == nil {
		log.Logger.Warnf("User not found with id: %d", userID)
		return sql.ErrNoRows
	}
	user.Name = profile.Name
	user.AvatarId = profile.Avatar
	err = u.userDao.UpdateUser(ctx, user)
	log.Logger.Infof("User profile updated for user id: %d\terr=%v", userID, err)
	return err
}
