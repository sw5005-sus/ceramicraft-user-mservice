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

type UserAddressDao interface {
	GetUserAddresses(ctx context.Context, userID int) ([]*model.UserAddress, error)
	CreateUserAddress(ctx context.Context, address *model.UserAddress) (int, error)
	UpdateUserAddress(ctx context.Context, address *model.UserAddress) (int, error)
	GetDefaultAddress(ctx context.Context, userID int) (*model.UserAddress, error)
}

type UserAddressDaoImpl struct {
	db *gorm.DB
}

var (
	userAddressOnce sync.Once
	userAddressDao  UserAddressDao
)

func GetUserAddressDao() UserAddressDao {
	userAddressOnce.Do(func() {
		if userAddressDao == nil {
			userAddressDao = &UserAddressDaoImpl{db: repository.DB}
		}
	})
	return userAddressDao
}

func (dao *UserAddressDaoImpl) GetUserAddresses(ctx context.Context, userID int) ([]*model.UserAddress, error) {
	var addresses []*model.UserAddress
	ret := dao.db.WithContext(ctx).Where("user_id = ? and deleted_at is null", userID).Find(&addresses)
	if ret.Error != nil {
		log.Logger.Errorf("Failed to get user addresses: %v", ret.Error)
		return nil, ret.Error
	}
	return addresses, nil
}

func (dao *UserAddressDaoImpl) CreateUserAddress(ctx context.Context, address *model.UserAddress) (int, error) {
	ret := dao.db.WithContext(ctx).Create(address)
	if ret.Error != nil {
		log.Logger.Errorf("Failed to create user address: %v", ret.Error)
		return 0, ret.Error
	}
	log.Logger.Infof("User address created with ID: %d", address.ID)
	return address.ID, nil
}

func (dao *UserAddressDaoImpl) UpdateUserAddress(ctx context.Context, address *model.UserAddress) (int, error) {
	ret := dao.db.WithContext(ctx).Model(&model.UserAddress{}).
		Where("id = ? and user_id=?", address.ID, address.UserID).
		Updates(address)
	if ret.Error != nil {
		log.Logger.Errorf("Failed to update user address: %v", ret.Error)
		return 0, ret.Error
	}
	log.Logger.Infof("User address updated with ID: %d", address.ID)

	return int(ret.RowsAffected), ret.Error
}

func (dao *UserAddressDaoImpl) GetDefaultAddress(ctx context.Context, userID int) (*model.UserAddress, error) {
	var address model.UserAddress
	ret := dao.db.WithContext(ctx).Where("user_id = ? and deleted_at is null", userID).Order("default_mark_time desc, id asc").First(&address)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Errorf("Failed to get default user address: %v", ret.Error)
		return nil, ret.Error
	}
	return &address, nil
}
