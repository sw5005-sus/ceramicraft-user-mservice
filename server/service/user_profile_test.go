package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/dao/mocks"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
)

func TestGetUserProfileService(t *testing.T) {
	service1 := GetUserProfileService()
	service2 := GetUserProfileService()
	assert.Equal(t, service1, service2)
}

func TestGetUserProfile(t *testing.T) {
	mockDao := new(mocks.UserDao)
	service := &UserProfileServiceImpl{userDao: mockDao}
	userID := 1

	mockDao.On("GetUserById", context.Background(), userID).Return(&model.User{ID: userID, Email: "test@example.com", Name: "Test User", AvatarId: "avatar123"}, nil)

	profile, err := service.GetUserProfile(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, userID, profile.ID)
	assert.Equal(t, "test@example.com", profile.Email)
	assert.Equal(t, "Test User", profile.Name)
	assert.Equal(t, "avatar123", profile.Avatar)

	mockDao.AssertExpectations(t)
}

func TestGetUserProfile_UserNotFound(t *testing.T) {
	initEnv()
	mockDao := new(mocks.UserDao)
	userProfileService := &UserProfileServiceImpl{userDao: mockDao}
	userID := 2

	mockDao.On("GetUserById", context.Background(), userID).Return(nil, nil)

	profile, err := userProfileService.GetUserProfile(context.Background(), userID)
	assert.NoError(t, err)
	assert.Nil(t, profile)

	mockDao.AssertExpectations(t)
}

func TestUpdateUserProfile(t *testing.T) {
	initEnv()
	mockDao := new(mocks.UserDao)
	service := &UserProfileServiceImpl{userDao: mockDao}
	userID := 1
	profile := &data.UserProfileVO{Name: "Updated User", Avatar: "newAvatar123"}

	mockDao.On("GetUserById", context.Background(), userID).Return(&model.User{ID: userID, Email: "test@example.com", Name: "Test User", AvatarId: "avatar123"}, nil)
	mockDao.On("UpdateUser", context.Background(), mock.Anything).Return(nil)

	err := service.UpdateUserProfile(context.Background(), userID, profile)
	assert.NoError(t, err)

	mockDao.AssertExpectations(t)
}

func TestUpdateUserProfile_UserNotFound(t *testing.T) {
	initEnv()
	mockDao := new(mocks.UserDao)
	service := &UserProfileServiceImpl{userDao: mockDao}
	userID := 2
	profile := &data.UserProfileVO{Name: "Updated User", Avatar: "newAvatar123"}

	mockDao.On("GetUserById", context.Background(), userID).Return(nil, nil)

	err := service.UpdateUserProfile(context.Background(), userID, profile)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	mockDao.AssertExpectations(t)
}
