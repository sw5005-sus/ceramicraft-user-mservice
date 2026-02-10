package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/dao/mocks"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
)

func TestGetUserAddressService(t *testing.T) {
	t.Run("Singleton Instance", func(t *testing.T) {
		service1 := GetUserAddressService()
		service2 := GetUserAddressService()
		if service1 != service2 {
			t.Errorf("Expected the same instance, got different instances")
		}
	})

	t.Run("UserAddressDao Not Nil", func(t *testing.T) {
		service := GetUserAddressService()
		if service == nil {
			t.Errorf("Expected service to be initialized, got nil")
		}
		if _, ok := service.(*UserAddressServiceImpl); !ok {
			t.Errorf("Expected service to be of type UserAddressServiceImpl")
		}
	})
}

func TestUserAddressService_GetUserAddress(t *testing.T) {
	initEnv()
	userAddressData := []*model.UserAddress{
		{ID: 1, UserID: 1, DefaultMarkTime: 4},
		{ID: 2, UserID: 1, DefaultMarkTime: 2},
		{ID: 3, UserID: 1, DefaultMarkTime: 10},
	}
	expAddrOrdinals := []int{3, 1, 2}
	t.Run("GetUserAddresses Success", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		userID := 1
		userDao.On("GetUserAddresses", ctx, userID).Return(userAddressData, nil)

		addresses, err := service.GetUserAddresses(ctx, userID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if addresses == nil {
			t.Errorf("Expected addresses to be returned, got nil")
		}
		if addresses[0].IsDefault != true {
			t.Errorf("Expected first address to be default, got %v", addresses[0].IsDefault)
		}
		for i, addr := range addresses {
			if addr.ID != expAddrOrdinals[i] {
				t.Errorf("Expected address ID %d at position %d, got %d", expAddrOrdinals[i], i, addr.ID)
			}
		}
	})

	t.Run("GetUserAddresses No Addresses", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}

		ctx := context.Background()
		userID := 999 // Assuming this user ID has no addresses
		userDao.On("GetUserAddresses", ctx, userID).Return([]*model.UserAddress{}, nil)

		addresses, err := service.GetUserAddresses(ctx, userID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(addresses) != 0 {
			t.Errorf("Expected no addresses, got %d", len(addresses))
		}
	})

	t.Run("GetUserAddresses Error", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		userID := -1 // Assuming this user ID will cause an error
		userDao.On("GetUserAddresses", ctx, userID).Return(nil, assert.AnError)
		_, err := service.GetUserAddresses(ctx, userID)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}
func TestUserAddressService_CreateUserAddress(t *testing.T) {
	initEnv()
	t.Run("CreateUserAddress Success", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		address := &data.UserAddressVO{
			UserID:       1,
			ZipCode:      "123456",
			Country:      "Country",
			Province:     "Province",
			City:         "City",
			Detail:       "Detail",
			FirstName:    "First",
			LastName:     "Last",
			ContactPhone: "1234567890",
			IsDefault:    true,
		}
		userDao.On("CreateUserAddress", ctx, mock.MatchedBy(func(userAddress *model.UserAddress) bool {
			return userAddress.DefaultMarkTime > 0
		})).Return(1, nil)

		createdAddress, err := service.CreateUserAddress(ctx, address)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if createdAddress.ID != 1 {
			t.Errorf("Expected address ID 1, got %d", createdAddress.ID)
		}
	})

	t.Run("CreateUserAddress Error", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		address := &data.UserAddressVO{
			UserID:       1,
			ZipCode:      "123456",
			Country:      "Country",
			Province:     "Province",
			City:         "City",
			Detail:       "Detail",
			FirstName:    "First",
			LastName:     "Last",
			ContactPhone: "1234567890",
			IsDefault:    true,
		}
		userDao.On("CreateUserAddress", ctx, mock.Anything).Return(0, assert.AnError)

		_, err := service.CreateUserAddress(ctx, address)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}
func TestUserAddressService_UpdateUserAddress(t *testing.T) {
	initEnv()
	t.Run("UpdateUserAddress Success", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		address := &data.UserAddressVO{
			ID:           1,
			UserID:       1,
			ZipCode:      "123456",
			Country:      "Country",
			Province:     "Province",
			City:         "City",
			Detail:       "Detail",
			FirstName:    "First",
			LastName:     "Last",
			ContactPhone: "1234567890",
			IsDefault:    true,
		}
		userDao.On("UpdateUserAddress", ctx, mock.MatchedBy(func(userAddress *model.UserAddress) bool {
			return userAddress.ID == 1 && userAddress.DefaultMarkTime > 0
		})).Return(1, nil)

		err := service.UpdateUserAddress(ctx, address)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		userDao.AssertCalled(t, "UpdateUserAddress", ctx, mock.Anything)
	})

	t.Run("UpdateUserAddress No Rows Updated", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		address := &data.UserAddressVO{
			ID:           1,
			UserID:       1,
			ZipCode:      "123456",
			Country:      "Country",
			Province:     "Province",
			City:         "City",
			Detail:       "Detail",
			FirstName:    "First",
			LastName:     "Last",
			ContactPhone: "1234567890",
			IsDefault:    true,
		}
		userDao.On("UpdateUserAddress", ctx, mock.Anything).Return(0, nil)

		err := service.UpdateUserAddress(ctx, address)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("UpdateUserAddress Error", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		address := &data.UserAddressVO{
			ID:           1,
			UserID:       1,
			ZipCode:      "123456",
			Country:      "Country",
			Province:     "Province",
			City:         "City",
			Detail:       "Detail",
			FirstName:    "First",
			LastName:     "Last",
			ContactPhone: "1234567890",
			IsDefault:    true,
		}
		userDao.On("UpdateUserAddress", ctx, mock.Anything).Return(0, assert.AnError)

		err := service.UpdateUserAddress(ctx, address)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}
func TestUserAddressService_DeleteUserAddress(t *testing.T) {
	initEnv()
	t.Run("DeleteUserAddress Success", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		addressID := 1
		userID := 1

		userDao.On("UpdateUserAddress", ctx, mock.MatchedBy(func(userAddress *model.UserAddress) bool {
			return userAddress.ID == addressID && userAddress.UserID == userID && userAddress.DeletedAt.Valid
		})).Return(1, nil)

		err := service.DeleteUserAddress(ctx, addressID, userID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		userDao.AssertCalled(t, "UpdateUserAddress", ctx, mock.Anything)
	})

	t.Run("DeleteUserAddress No Rows Updated", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		addressID := 1
		userID := 1

		userDao.On("UpdateUserAddress", ctx, mock.Anything).Return(0, nil)

		err := service.DeleteUserAddress(ctx, addressID, userID)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("DeleteUserAddress Error", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		addressID := 1
		userID := 1

		userDao.On("UpdateUserAddress", ctx, mock.Anything).Return(0, assert.AnError)

		err := service.DeleteUserAddress(ctx, addressID, userID)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}
func TestUserAddressService_GetDefaultAddress(t *testing.T) {
	initEnv()
	t.Run("GetDefaultAddress Success", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		userID := 1
		expectedAddress := &model.UserAddress{
			ID:              1,
			UserID:          userID,
			ZipCode:         "123456",
			Country:         "Country",
			Province:        "Province",
			City:            "City",
			Detail:          "Detail",
			FirstName:       "First",
			LastName:        "Last",
			ContactPhone:    "1234567890",
			DefaultMarkTime: 4,
		}
		userDao.On("GetDefaultAddress", ctx, userID).Return(expectedAddress, nil)

		address, err := service.GetDefaultAddress(ctx, userID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if address == nil || address.ID != expectedAddress.ID {
			t.Errorf("Expected address ID %d, got %v", expectedAddress.ID, address)
		}
	})

	t.Run("GetDefaultAddress Not Found", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		userID := 999 // Assuming this user ID has no default address
		userDao.On("GetDefaultAddress", ctx, userID).Return(nil, nil)

		address, err := service.GetDefaultAddress(ctx, userID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if address != nil {
			t.Errorf("Expected no address, got %v", address)
		}
	})

	t.Run("GetDefaultAddress Error", func(t *testing.T) {
		userDao := new(mocks.UserAddressDao)
		service := UserAddressServiceImpl{
			userAddressDao: userDao,
		}
		ctx := context.Background()
		userID := 1
		userDao.On("GetDefaultAddress", ctx, userID).Return(nil, assert.AnError)

		address, err := service.GetDefaultAddress(ctx, userID)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if address != nil {
			t.Errorf("Expected no address, got %v", address)
		}
	})
}
