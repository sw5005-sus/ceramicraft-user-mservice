package service

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/utils"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/dao/mocks"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
)

func initEnv() {
	config.Config = &config.Conf{
		LogConfig: &config.LogConfig{
			Level:    "debug",
			FilePath: "",
		},
		EmailConfig: &config.EmailConfig{},
		KafkaConfig: &config.KafkaConfig{
			UserActivatedTopic: "user_activated",
		},
	}
	log.InitLogger()
	err := os.Setenv("JWT_SECRET", "TEST_SECRET_KEY")
	if err != nil {
		panic(err)
	}
	utils.InitJwtSecret()
}

func TestLogin(t *testing.T) {
	initEnv()
	ctx := context.Background()
	mockDao := new(mocks.UserDao)
	loginService := &LoginServiceImpl{
		userDao: mockDao,
	}
	hashedPwd, _ := HashPassword("correctpassword")
	existUser := &model.User{ID: 1, Email: "test@example.com", Password: hashedPwd}
	nonExistEmail := "nonexistent@example.com"
	mockDao.On("GetUserByEmail", mock.Anything, existUser.Email).Return(existUser, nil)
	mockDao.On("GetUserByEmail", mock.Anything, nonExistEmail).Return(nil, nil)

	tests := []struct {
		email    string
		password string
		expected string
		hasError bool
	}{
		{existUser.Email, "correctpassword", "token", false},
		{nonExistEmail, "anyPassword", "", true},
		{existUser.Email, "wrongpassword", "", true},
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			token, err := loginService.Login(ctx, test.email, test.password)
			if (err != nil) != test.hasError {
				t.Errorf("expected error: %v, got: %v", test.hasError, err)
			}
			if test.expected == "token" && token == "" {
				t.Errorf("expected token: %s, got: %s", test.expected, token)
			}
		})
	}
}
func TestGetLoginService(t *testing.T) {
	initEnv()

	// Test singleton behavior
	service1 := GetLoginService()
	service2 := GetLoginService()

	if service1 != service2 {
		t.Errorf("Expected the same instance, got different instances")
	}

	// Test if userDao is initialized
	if service1.userDao == nil {
		t.Error("Expected userDao to be initialized, got nil")
	}
}
