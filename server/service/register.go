package service

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/mq"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/proxy"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/dao"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/repository/model"
	"gorm.io/gorm"
)

type RegisterService interface {
	Register(ctx context.Context, email, password string) error
	VerifyAndActivate(ctx context.Context, activationCode string) error
}

type RegisterImpl struct {
	userDao        dao.UserDao
	userActivation dao.UserActivationDao
	emailService   proxy.EmailService
	txBeginner     repository.TxBeginner
	kafkaProducer  mq.KafkaProducer
}

var (
	registerServiceInst *RegisterImpl
	registerOnce        sync.Once
)

const activationExpiryDuration = time.Minute * 5

func GetRegisterService() *RegisterImpl {
	registerOnce.Do(func() {
		if registerServiceInst == nil {
			registerServiceInst = &RegisterImpl{
				userDao:        dao.GetUserDao(),
				userActivation: dao.GetUserActivationDao(),
				emailService:   proxy.GetEmailInstance(),
				txBeginner:     repository.DB,
				kafkaProducer:  mq.GetKafkaProducer(),
			}
		}
	})
	return registerServiceInst
}

func (rs *RegisterImpl) Register(ctx context.Context, email, password string) error {
	user, err := rs.userDao.GetUserByEmail(ctx, email)
	if err != nil {
		log.Logger.Errorf("Failed to get user by email: %v", err)
		return err
	}
	if user != nil && user.Status == model.UserStatusActive {
		log.Logger.Errorf("User already exists with email: %s", email)
		return errors.New("user already exists")
	}
	if user == nil {
		hashedPassword, err := HashPassword(password)
		if err != nil {
			log.Logger.Errorf("Failed to hash password: %v", err)
			return err
		}
		user = &model.User{
			Email:     email,
			Password:  hashedPassword,
			Status:    model.UserStatusInactive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err = rs.userDao.CreateUser(ctx, user)
		if err != nil {
			log.Logger.Errorf("Failed to create user: %v", err)
			return err
		}
	}
	code, err := generateVerificationCode()
	if err != nil {
		log.Logger.Errorf("Failed to generate verification code: %v", err)
		return err
	}
	err = rs.userActivation.Replace(ctx, &model.UserActivation{
		UserID:    user.ID,
		Code:      code,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(activationExpiryDuration),
	})
	if err != nil {
		log.Logger.Errorf("Failed to create user activation: %v", err)
		return err
	}
	err = rs.emailService.Send("Your activation code is: "+code, email, "CermiCraft Activation Code")
	if err != nil {
		log.Logger.Errorf("Failed to send activation email: %v", err)
		return err
	}
	log.Logger.Infof("Activation email sent for user: %d", user.ID)
	return nil
}

func (rs *RegisterImpl) VerifyAndActivate(ctx context.Context, activationCode string) error {
	userActivation, err := rs.userActivation.GetByCode(ctx, activationCode)
	if err != nil {
		log.Logger.Errorf("Failed to get user activation by code: %v", err)
		return err
	}
	if userActivation == nil || userActivation.ExpiresAt.Before(time.Now()) {
		log.Logger.Warnf("Invalid or expired activation code: %s", activationCode)
		return errors.New("invalid or expired activation code")
	}
	err = rs.txBeginner.Transaction(func(tx *gorm.DB) error {
		curTime := time.Now()
		err = rs.userDao.UpdateUserInTransaction(ctx, &model.User{ID: userActivation.UserID, Status: model.UserStatusActive, ActivateTime: &curTime, UpdatedAt: curTime}, tx)
		if err != nil {
			log.Logger.Errorf("Failed to update user status: %v", err)
			return err
		}
		log.Logger.Infof("User %d activated successfully", userActivation.UserID)
		err = rs.userActivation.DeleteByUserId(ctx, userActivation.UserID, tx)
		if err != nil {
			log.Logger.Warnf("Failed to delete user activation after activation: %v", err)
			return err
		}
		log.Logger.Infof("User activation %d marked as used", userActivation.ID)
		eventMsg := &mq.UserActivatedEvent{UserID: userActivation.UserID, ActivateTime: curTime.Unix()}
		err = rs.kafkaProducer.Produce(ctx, config.Config.KafkaConfig.UserActivatedTopic, fmt.Sprintf("%d", userActivation.UserID), eventMsg.ToBytes())
		if err != nil {
			log.Logger.Errorf("Failed to produce user activated event: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Logger.Errorf("Failed to start transaction: %v", err)
		return err
	}
	return nil
}

func generateVerificationCode() (string, error) {
	var num uint32
	err := binary.Read(rand.Reader, binary.BigEndian, &num)
	if err != nil {
		return "", err
	}
	codeInt := int(num % 1000000) // Convert to int
	return fmt.Sprintf("%06d", codeInt), nil
}
