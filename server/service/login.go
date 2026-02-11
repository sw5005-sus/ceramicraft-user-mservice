package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/sw5005-sus/ceramicraft-user-mservice/common/bo"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/utils"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/dao"
)

type LoginService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context) error
}

type LoginServiceImpl struct {
	userDao dao.UserDao
}

var (
	loginServiceOnce sync.Once
	loginServiceInst *LoginServiceImpl
)

func GetLoginService() *LoginServiceImpl {
	loginServiceOnce.Do(func() {
		loginServiceInst = &LoginServiceImpl{userDao: dao.GetUserDao()}
	})
	return loginServiceInst
}

func (ls *LoginServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := ls.userDao.GetUserByEmail(ctx, email)
	if err != nil {
		log.Logger.Errorf("Failed to get user by email: %v", err)
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	if VerifyPassword(user.Password, password) != nil {
		log.Logger.Errorf("Failed to verify password")
		return "", fmt.Errorf("invalid password")
	}

	token, err := utils.GenerateJWTToken(&bo.UserBO{ID: user.ID, Email: user.Email})
	if err != nil {
		return "", err
	}

	return token, nil
}
