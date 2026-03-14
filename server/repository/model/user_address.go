package model

import (
	"database/sql"
	"time"

	ceramicraftsecure "github.com/sw5005-sus/ceramicraft-secure"
)

type UserAddress struct {
	ID              int          `gorm:"primaryKey"`
	UserID          int          `gorm:"type:int;not null"`
	ZipCode         string       `gorm:"type:varchar(128);not null"`
	Country         string       `gorm:"type:varchar(64);not null"`
	Province        string       `gorm:"type:varchar(64);not null"`
	City            string       `gorm:"type:varchar(64);not null"`
	Detail          string       `gorm:"type:varchar(255);not null"`
	FirstName       string       `gorm:"type:varchar(64);not null"`
	LastName        string       `gorm:"type:varchar(64);not null"`
	ContactPhone    string       `gorm:"type:varchar(128);not null"`
	DefaultMarkTime int64        `gorm:"column:default_mark_time;not null;default:0"`
	CreatedAt       time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time    `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       sql.NullTime `gorm:"column:deleted_at;"`
}

func (UserAddress) TableName() string {
	return "user_addresses"
}

func (u *UserAddress) Encrypt() error {
	if u.ZipCode != "" {
		encryptZipCode, err := ceramicraftsecure.AesEncrypt(u.ZipCode)
		if err != nil {
			return err
		}
		u.ZipCode = encryptZipCode
	}

	if u.Detail != "" {
		encryptedDetail, err := ceramicraftsecure.AesEncrypt(u.Detail)
		if err != nil {
			return err
		}
		u.Detail = encryptedDetail
	}

	if u.ContactPhone != "" {
		encryptedContactPhone, err := ceramicraftsecure.AesEncrypt(u.ContactPhone)
		if err != nil {
			return err
		}
		u.ContactPhone = encryptedContactPhone
	}

	return nil
}

func (u *UserAddress) Decrypt() error {
	if u.ZipCode != "" {
		decryptedZipCode, err := ceramicraftsecure.AesDecrypt(u.ZipCode)
		if err != nil {
			return err
		}
		u.ZipCode = decryptedZipCode
	}

	if u.Detail != "" {
		decryptedDetail, err := ceramicraftsecure.AesDecrypt(u.Detail)
		if err != nil {
			return err
		}
		u.Detail = decryptedDetail
	}

	if u.ContactPhone != "" {
		decryptedContactPhone, err := ceramicraftsecure.AesDecrypt(u.ContactPhone)
		if err != nil {
			return err
		}
		u.ContactPhone = decryptedContactPhone
	}
	return nil

}
