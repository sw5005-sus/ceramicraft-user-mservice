package dao

import (
	"context"
	"errors"
	"sync"

	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository"

	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
	"gorm.io/gorm"
)

type UserDao interface {
	CreateUser(ctx context.Context, user *model.User) (int, error)
	UpdateUserInTransaction(ctx context.Context, user *model.User, tx *gorm.DB) error
	UpdateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(context.Context, string) (*model.User, error)
	GetUserById(context.Context, int) (*model.User, error)
}

type UserDaoImpl struct {
	db *gorm.DB
}

var (
	userOnce sync.Once
	userDao  *UserDaoImpl
)

func GetUserDao() *UserDaoImpl {
	userOnce.Do(func() {
		if userDao == nil {
			userDao = &UserDaoImpl{db: repository.DB}
		}
	})
	return userDao
}

func (dao *UserDaoImpl) CreateUser(ctx context.Context, user *model.User) (int, error) {
	ret := dao.db.WithContext(ctx).Save(&user)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrDuplicatedKey) {
			return 0, errors.New("user already exists")
		}
		log.Logger.Errorf("Failed to create user: ", ret.Error)
		return 0, ret.Error
	}
	log.Logger.Infof("User created with ID: %d", user.ID)
	return user.ID, nil
}

func (dao *UserDaoImpl) UpdateUserInTransaction(ctx context.Context, user *model.User, tx *gorm.DB) error {
	ret := tx.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(user)
	if ret.Error != nil {
		log.Logger.Errorf("Failed to update user: %v", ret.Error)
		return ret.Error
	}
	return nil
}

func (dao *UserDaoImpl) UpdateUser(ctx context.Context, user *model.User) error {
	ret := dao.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(user)
	if ret.Error != nil {
		log.Logger.Errorf("Failed to update user: %v", ret.Error)
		return ret.Error
	}
	return nil
}

func (dao *UserDaoImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	ret := dao.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Errorf("Failed to get user by email: %v", ret.Error)
		return nil, ret.Error
	}
	return &user, nil
}

func (dao *UserDaoImpl) GetUserById(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	ret := dao.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Errorf("Failed to get user by id: %v", ret.Error)
		return nil, ret.Error
	}
	return &user, nil
}
