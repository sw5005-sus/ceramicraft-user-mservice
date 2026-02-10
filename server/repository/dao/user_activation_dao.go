package dao

import (
	"context"
	"errors"
	"sync"

	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
	"gorm.io/gorm"
)

type UserActivationDao interface {
	DeleteByUserId(ctx context.Context, userId int, tx *gorm.DB) error
	Create(ctx context.Context, activation *model.UserActivation, tx *gorm.DB) error
	GetByCode(ctx context.Context, code string) (*model.UserActivation, error)
	Update(ctx context.Context, activation *model.UserActivation, tx *gorm.DB) error
	Replace(ctx context.Context, activation *model.UserActivation) error
}

type UserActivationDaoImpl struct {
	db *gorm.DB
}

var (
	userActiveSyncOnce sync.Once
	userActivationDao  *UserActivationDaoImpl
)

func GetUserActivationDao() *UserActivationDaoImpl {
	userActiveSyncOnce.Do(func() {
		if userActivationDao == nil {
			userActivationDao = &UserActivationDaoImpl{db: repository.DB}
		}
	})
	return userActivationDao
}

func (dao *UserActivationDaoImpl) Replace(ctx context.Context, activation *model.UserActivation) error {
	ret := dao.db.WithContext(ctx).Save(activation)
	return ret.Error
}

func (dao *UserActivationDaoImpl) DeleteByUserId(ctx context.Context, userId int, tx *gorm.DB) error {
	ret := tx.WithContext(ctx).Where("user_id = ?", userId).Delete(&model.UserActivation{})
	return ret.Error
}

func (dao *UserActivationDaoImpl) Create(ctx context.Context, activation *model.UserActivation, tx *gorm.DB) error {
	ret := tx.WithContext(ctx).Create(activation)
	return ret.Error
}

func (dao *UserActivationDaoImpl) GetByCode(ctx context.Context, code string) (*model.UserActivation, error) {
	var activation model.UserActivation
	ret := dao.db.WithContext(ctx).Where("code = ?", code).First(&activation)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ret.Error
	}
	return &activation, nil
}

func (dao *UserActivationDaoImpl) Update(ctx context.Context, activation *model.UserActivation, tx *gorm.DB) error {
	ret := tx.WithContext(ctx).Save(activation)
	return ret.Error
}
