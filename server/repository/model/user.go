package model

import (
	"time"
)

const (
	UserStatusInactive = -1
	UserStatusActive   = 1
)

type User struct {
	ID           int        `gorm:"primaryKey"`
	Email        string     `gorm:"type:varchar(128);unique;not null"`
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
