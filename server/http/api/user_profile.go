package api

import (
	"database/sql"
	"net/http"

	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/service"

	"github.com/gin-gonic/gin"
)

// Get Current UserProfile.
// @Summary Get User Profile
// @Description This endpoint allows current login user fetch his/her profile in JSON format.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} data.BaseResponse "data is UserProfileVO"
// @Failure 500 {object} data.BaseResponse
// @Router /user-ms/v1/customer/users/self [get]
func GetUserProfile(c *gin.Context) {
	userId, exist := c.Get("userID")
	if !exist || userId.(int) <= 0 {
		c.JSON(http.StatusUnauthorized, data.BaseResponse{Code: http.StatusUnauthorized, ErrMsg: "Unauthorized"})
		return
	}
	userProfile, err := service.GetUserProfileService().GetUserProfile(c.Request.Context(), userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, data.BaseResponse{Code: http.StatusInternalServerError, ErrMsg: err.Error()})
		return
	}
	if userProfile == nil {
		c.JSON(http.StatusNotFound, data.BaseResponse{Code: http.StatusNotFound, ErrMsg: "User not found"})
		return
	}
	userDefaultAddress, err := service.GetUserAddressService().GetDefaultAddress(c, userId.(int))
	if err != nil {
		log.Logger.Errorf("Failed to get default address for user ID %d: %v", userId.(int), err)
	}
	userProfile.DefaultAddress = userDefaultAddress
	c.JSON(http.StatusOK, data.BaseResponse{Code: http.StatusOK, Data: userProfile})
}

// Update UserProfile.
// @Summary Update User Profile
// @Description This endpoint allows current login user update his/her profile in JSON format.
// @Tags User
// @Accept json
// @Produce json
// @Param user body data.UserProfileVO true "User profile to update"
// @Success 200 {object} data.BaseResponse "data is UserProfileVO"
// @Failure 404 {object} data.BaseResponse
// @Failure 500 {object} data.BaseResponse
// @Router /user-ms/v1/customer/users/self [put]
func UpdateUserProfile(c *gin.Context) {
	var userProfile data.UserProfileVO
	if err := c.ShouldBindJSON(&userProfile); err != nil {
		c.JSON(http.StatusBadRequest, data.BaseResponse{Code: http.StatusBadRequest, ErrMsg: "Invalid input"})
		return
	}
	userId, exist := c.Get("userID")
	if !exist || userId.(int) <= 0 {
		c.JSON(http.StatusUnauthorized, data.BaseResponse{Code: http.StatusUnauthorized, ErrMsg: "Unauthorized"})
		return
	}
	if userProfile.ID != userId.(int) {
		log.Logger.Warnf("User ID mismatch: token ID %d, payload ID %d", userId.(int), userProfile.ID)
		c.JSON(http.StatusBadRequest, data.BaseResponse{Code: http.StatusBadRequest, ErrMsg: "User ID does not match current user"})
		return
	}
	userProfile.Email = "" // Email should not be updated here
	err := service.GetUserProfileService().UpdateUserProfile(c.Request.Context(), userProfile.ID, &userProfile)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, data.BaseResponse{Code: http.StatusNotFound, ErrMsg: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, data.BaseResponse{Code: http.StatusInternalServerError, ErrMsg: err.Error()})
		return
	}
	updatedUserProfile, err := service.GetUserProfileService().GetUserProfile(c.Request.Context(), userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, data.BaseResponse{Code: http.StatusInternalServerError, ErrMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.BaseResponse{Code: http.StatusOK, Data: updatedUserProfile})
}
