package service

import (
	"context"
	"database/sql"
	"sort"
	"sync"
	"time"

	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/dao"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
)

type UserAddressService interface {
	GetUserAddresses(ctx context.Context, userID int) ([]*data.UserAddressVO, error)
	GetDefaultAddress(ctx context.Context, userID int) (*data.UserAddressVO, error)
	CreateUserAddress(ctx context.Context, address *data.UserAddressVO) (*data.UserAddressVO, error)
	UpdateUserAddress(ctx context.Context, address *data.UserAddressVO) error
	DeleteUserAddress(ctx context.Context, addressID int, userId int) error
}

type UserAddressServiceImpl struct {
	userAddressDao dao.UserAddressDao
}

var (
	userAddressServiceInst UserAddressService
	userAddressOnce        sync.Once
)

func GetUserAddressService() UserAddressService {
	userAddressOnce.Do(func() {
		userAddressServiceInst = &UserAddressServiceImpl{
			userAddressDao: dao.GetUserAddressDao(),
		}
	})
	return userAddressServiceInst
}

func (u *UserAddressServiceImpl) GetUserAddresses(ctx context.Context, userID int) ([]*data.UserAddressVO, error) {
	addresses, err := u.userAddressDao.GetUserAddresses(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(addresses) == 0 {
		log.Logger.Infof("No addresses found for user ID: %d", userID)
		return []*data.UserAddressVO{}, nil
	}
	var defaultAddr *data.UserAddressVO
	maxDefaultMarkTime := int64(0)
	var addressVOs []*data.UserAddressVO
	for _, addr := range addresses {
		addressVOs = append(addressVOs, &data.UserAddressVO{
			ID:           addr.ID,
			UserID:       addr.UserID,
			ZipCode:      addr.ZipCode,
			Country:      addr.Country,
			Province:     addr.Province,
			City:         addr.City,
			Detail:       addr.Detail,
			FirstName:    addr.FirstName,
			LastName:     addr.LastName,
			ContactPhone: addr.ContactPhone,
			IsDefault:    false,
		})
		// Find the most recently marked default address
		if maxDefaultMarkTime < addr.DefaultMarkTime {
			maxDefaultMarkTime = addr.DefaultMarkTime
			defaultAddr = addressVOs[len(addressVOs)-1]
		}
	}
	if defaultAddr != nil {
		defaultAddr.IsDefault = true
	}
	sort.SliceStable(addressVOs, func(i, j int) bool {
		if addressVOs[i].IsDefault != addressVOs[j].IsDefault {
			return addressVOs[i].IsDefault
		}
		return addressVOs[i].ID < addressVOs[j].ID
	})
	return addressVOs, nil
}

func (u *UserAddressServiceImpl) GetDefaultAddress(ctx context.Context, userID int) (*data.UserAddressVO, error) {
	addr, err := u.userAddressDao.GetDefaultAddress(ctx, userID)
	if err != nil {
		log.Logger.Errorf("Failed to get default address for user ID %d: %v", userID, err)
		return nil, err
	}
	if addr == nil {
		return nil, nil
	}
	return &data.UserAddressVO{
		ID:           addr.ID,
		UserID:       addr.UserID,
		ZipCode:      addr.ZipCode,
		Country:      addr.Country,
		Province:     addr.Province,
		City:         addr.City,
		Detail:       addr.Detail,
		FirstName:    addr.FirstName,
		LastName:     addr.LastName,
		ContactPhone: addr.ContactPhone,
		IsDefault:    true,
	}, nil
}

func (u *UserAddressServiceImpl) CreateUserAddress(ctx context.Context, address *data.UserAddressVO) (*data.UserAddressVO, error) {
	addrModel := &model.UserAddress{
		UserID:       address.UserID,
		ZipCode:      address.ZipCode,
		Country:      address.Country,
		Province:     address.Province,
		City:         address.City,
		Detail:       address.Detail,
		FirstName:    address.FirstName,
		LastName:     address.LastName,
		ContactPhone: address.ContactPhone,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if address.IsDefault {
		addrModel.DefaultMarkTime = time.Now().Unix()
	}
	id, err := u.userAddressDao.CreateUserAddress(ctx, addrModel)
	if err != nil {
		return nil, err
	}
	address.ID = id
	return address, nil
}

func (u *UserAddressServiceImpl) UpdateUserAddress(ctx context.Context, address *data.UserAddressVO) error {
	addrModel := &model.UserAddress{
		ID:           address.ID,
		UserID:       address.UserID,
		ZipCode:      address.ZipCode,
		Country:      address.Country,
		Province:     address.Province,
		City:         address.City,
		Detail:       address.Detail,
		FirstName:    address.FirstName,
		LastName:     address.LastName,
		ContactPhone: address.ContactPhone,
		UpdatedAt:    time.Now(),
	}
	if address.IsDefault {
		addrModel.DefaultMarkTime = time.Now().Unix()
	}
	ret, err := u.userAddressDao.UpdateUserAddress(ctx, addrModel)
	if err != nil {
		log.Logger.Errorf("Failed to update user address: %v", err)
		return err
	}
	if ret == 0 {
		log.Logger.Warnf("No user address updated for ID: %d", address.ID)
		return sql.ErrNoRows
	}
	return nil
}

func (u *UserAddressServiceImpl) DeleteUserAddress(ctx context.Context, addressID int, userId int) error {
	addrModel := &model.UserAddress{
		ID:        addressID,
		UserID:    userId,
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
	ret, err := u.userAddressDao.UpdateUserAddress(ctx, addrModel)
	if err != nil {
		log.Logger.Errorf("Failed to update user address: %v", err)
		return err
	}
	if ret == 0 {
		log.Logger.Warnf("No user address updated for ID: %d", addressID)
		return sql.ErrNoRows
	}
	return nil
}
