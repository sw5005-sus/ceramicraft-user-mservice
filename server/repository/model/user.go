package model

import (
	"time"

	ceramicraftsecure "github.com/sw5005-sus/ceramicraft-secure"
)

const (
	UserStatusInactive = -1
	UserStatusActive   = 1
)

type User struct {
	ID           int        `gorm:"primaryKey"`
	Email        string     `gorm:"type:varchar(128);unique;not null"`
	EmailHash    string     `gorm:"type:varchar(128);unique;not null"`
	ZitadelSub   string     `gorm:"type:varchar(128);default:''"`
	Password     string     `gorm:"type:varchar(255);not null"`
	Status       int        `gorm:"type:int;not null"`
	Name         string     `gorm:"type:varchar(64)"`
	AvatarId     string     `gorm:"type:varchar(64)"`
	ActivateTime *time.Time `gorm:"column:activate_time"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName sets the insert table name for this struct type
func (User) TableName() string {
	return "users"
}

func (u *User) EncryptAndSetHash() error {
	if u.Email == "" {
		return nil
	}
	if u.EmailHash == "" {
		emailHash, err := GetEmailHash(u.Email)
		if err != nil {
			return err
		}
		u.EmailHash = emailHash
	}
	encryptedEmail, err := ceramicraftsecure.AesEncrypt(u.Email)
	if err != nil {
		return err
	}
	u.Email = encryptedEmail
	return nil
}

func (u *User) Decrpt() error {
	if u.EmailHash == "" {
		return nil
	}
	decryptedEmail, err := ceramicraftsecure.AesDecrypt(u.Email)
	if err != nil {
		return err
	}
	u.Email = decryptedEmail
	return nil
}

func GetEmailHash(plainEmial string) (string, error) {
	hashRet, err := ceramicraftsecure.GenHmacSha256(plainEmial)
	if err != nil {
		return "", err
	}
	return hashRet, nil
}
