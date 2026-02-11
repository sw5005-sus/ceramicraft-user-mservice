package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/service"
)

// Register handles the user registration process.
// @Summary Register a new user
// @Description This endpoint allows a new user to register by providing their details in JSON format.
// @Tags Register
// @Accept json
// @Produce json
// @Param user body data.UserLoginVO true "User registration details"
// @Param client path string true "Client identifier" Enums(customer, merchant)
// @Success 200
// @Failure 500 {object} data.BaseResponse
// @Router /user-ms/v1/{client}/users [post]
func Register(c *gin.Context) {
	user := &data.UserLoginVO{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.GetRegisterService().Register(c.Request.Context(), user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful, please check your email to activate your account"})
}

// Activate handles the user registration activation process.
// @Summary Activate a new user
// @Description This endpoint allows a new user to activate by providing their verification code in JSON format.
// @Tags Register
// @Accept json
// @Produce json
// @Param user body data.UserActivateReq true "User activate request"
// @Param client path string true "Client identifier" Enums(customer, merchant)
// @Success 200
// @Failure 500 {object} data.BaseResponse
// @Router /user-ms/v1/{client}/users/activate [put]
func Validate(c *gin.Context) {
	req := &data.UserActivateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.GetRegisterService().VerifyAndActivate(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Activation successful, you can now log in"})
}
