package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/service"
)

// Create UserAddress.
// @Summary Create New User Address
// @Description This endpoint allows user create new address in JSON format.
// @Tags UserAddress
// @Accept json
// @Produce json
// @Param address body data.UserAddressVO true "user address to create"
// @Success 201 {object} data.BaseResponse "data is UserAddressVO"
// @Failure 500 {object} data.BaseResponse
// @Router /user-ms/v1/customer/users/self/addresses [post]
func AddUserAddress(c *gin.Context) {
	userAddress := &data.UserAddressVO{}
	if err := c.ShouldBindJSON(userAddress); err != nil {
		c.JSON(http.StatusBadRequest, data.BaseResponse{Code: http.StatusBadRequest, ErrMsg: err.Error()})
		return
	}
	userId, exist := c.Get("userID")
	if !exist || userId.(int) <= 0 {
		c.JSON(http.StatusUnauthorized, data.BaseResponse{Code: http.StatusUnauthorized, ErrMsg: "Unauthorized"})
		return
	}
	userAddress.UserID = userId.(int)
	userAddress, err := service.GetUserAddressService().CreateUserAddress(c.Request.Context(), userAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, data.BaseResponse{Code: http.StatusInternalServerError, ErrMsg: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, data.BaseResponse{Code: http.StatusCreated, Data: userAddress})
}

// Update UserAddress.
// @Summary Update existing User Address
// @Description This endpoint allows user update their old address in JSON format.
// @Tags UserAddress
// @Accept json
// @Produce json
// @Param address body data.UserAddressVO true "user address to update"
// @Success 200 {object} data.BaseResponse "data is UserAddressVO"
// @Failure 500 {object} data.BaseResponse
// @Router /user-ms/v1/customer/users/self/addresses/{address_id} [put]
func UpdateUserAddress(c *gin.Context) {
	addressID := c.Param("address_id")
	if addressID == "" {
		c.JSON(http.StatusBadRequest, data.BaseResponse{Code: http.StatusBadRequest, ErrMsg: "Address ID is required"})
		return
	}
	userAddress := &data.UserAddressVO{}
	if err := c.ShouldBindJSON(userAddress); err != nil {
		c.JSON(http.StatusBadRequest, data.BaseResponse{Code: http.StatusBadRequest, ErrMsg: err.Error()})
		return
	}
	if id, err := strconv.Atoi(addressID); err != nil || id != userAddress.ID {
		c.JSON(http.StatusBadRequest, data.BaseResponse{Code: http.StatusBadRequest, ErrMsg: "Address ID is invalid"})
		return
	}
	userId, exist := c.Get("userID")
	if !exist || userId.(int) <= 0 {
		c.JSON(http.StatusUnauthorized, data.BaseResponse{Code: http.StatusUnauthorized, ErrMsg: "Unauthorized"})
		return
	}
	userAddress.UserID = userId.(int)
	err := service.GetUserAddressService().UpdateUserAddress(c.Request.Context(), userAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, data.BaseResponse{Code: http.StatusNotFound, ErrMsg: "Address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, data.BaseResponse{Code: http.StatusInternalServerError, ErrMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.BaseResponse{Code: http.StatusOK, Data: userAddress})
}

// List UserAddress.
// @Summary List User Addresses
// @Description This endpoint allows user list all their addresses.
// @Tags UserAddress
// @Accept json
// @Produce json
// @Success 200 {object} data.BaseResponse "data is []UserAddressVO"
// @Failure 500 {object} data.BaseResponse
// @Router /user-ms/v1/customer/users/self/addresses [get]
func ListUserAddresses(c *gin.Context) {
	userId, exist := c.Get("userID")
	if !exist || userId.(int) <= 0 {
		c.JSON(http.StatusUnauthorized, data.BaseResponse{Code: http.StatusUnauthorized, ErrMsg: "Unauthorized"})
		return
	}
	userAddresses, err := service.GetUserAddressService().GetUserAddresses(c.Request.Context(), userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, data.BaseResponse{Code: http.StatusInternalServerError, ErrMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.BaseResponse{Code: http.StatusOK, Data: userAddresses})
}

// Delete UserAddress.
// @Summary Delete User Address
// @Description This endpoint allows user delete their address.
// @Tags UserAddress
// @Accept json
// @Produce json
// @Param address_id path int true "Address ID"
// @Success 200 {object} data.BaseResponse
// @Failure 500 {object} data.BaseResponse
// @Router /user-ms/v1/customer/users/self/addresses/{address_id} [delete]
func DeleteUserAddress(c *gin.Context) {
	addressID := c.Param("address_id")
	if addressID == "" {
		c.JSON(http.StatusBadRequest, data.BaseResponse{Code: http.StatusBadRequest, ErrMsg: "Address ID is required"})
		return
	}
	id, err := strconv.Atoi(addressID)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, data.BaseResponse{Code: http.StatusBadRequest, ErrMsg: "Address ID is invalid"})
		return
	}
	userId, exist := c.Get("userID")
	if !exist || userId.(int) <= 0 {
		c.JSON(http.StatusUnauthorized, data.BaseResponse{Code: http.StatusUnauthorized, ErrMsg: "Unauthorized"})
		return
	}
	err = service.GetUserAddressService().DeleteUserAddress(c.Request.Context(), id, userId.(int))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, data.BaseResponse{Code: http.StatusNotFound, ErrMsg: "Address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, data.BaseResponse{Code: http.StatusInternalServerError, ErrMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.BaseResponse{Code: http.StatusOK, Data: "Address deleted successfully"})
}
