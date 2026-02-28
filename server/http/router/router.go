package router

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/sw5005-sus/ceramicraft-user-mservice/common/middleware"
	_ "github.com/sw5005-sus/ceramicraft-user-mservice/server/docs"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/http/api"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

const (
	serviceURIPrefix = "/user-ms/v1"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("password", passwordStrengthValidator)
		if err != nil {
			panic(err)
		}
	}

	basicGroup := r.Group(serviceURIPrefix)
	{
		basicGroup.GET("/swagger/*any", gs.WrapHandler(
			swaggerFiles.Handler,
			gs.URL("/user-ms/v1/swagger/doc.json"),
		))
		basicGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	v1UnAuthed := basicGroup.Group("")
	{
		v1UnAuthed.POST("/customer/login", api.UserLogin)
		v1UnAuthed.POST("/customer/users", api.Register)
		v1UnAuthed.PUT("/customer/users/activate", api.Validate)
		v1UnAuthed.POST("/customer/oauth-callback", api.OAuthCallback)

		v1UnAuthed.POST("/merchant/login", api.UserLogin)
	}
	v1Authed := basicGroup.Group("")
	{
		v1Authed.Use(middleware.AuthMiddleware())
		v1Authed.POST("/customer/logout", api.UserLogout)
		v1Authed.GET("/customer/users/self", api.GetUserProfile)
		v1Authed.PUT("/customer/users/self", api.UpdateUserProfile)
		v1Authed.GET("/customer/users/self/addresses", api.ListUserAddresses)
		v1Authed.POST("/customer/users/self/addresses", api.AddUserAddress)
		v1Authed.PUT("/customer/users/self/addresses/:address_id", api.UpdateUserAddress)
		v1Authed.DELETE("/customer/users/self/addresses/:address_id", api.DeleteUserAddress)

		v1Authed.POST("/merchant/logout", api.UserLogout)
	}
	return r
}

// Custom password validation rules
var passwordStrengthValidator validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// atleast one letter
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
	// atleast one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	// min length 8
	isValidLength := len(password) >= 8

	return hasLetter && hasDigit && isValidLength
}
