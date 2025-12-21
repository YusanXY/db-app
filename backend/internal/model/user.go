package model

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID            uint64         `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email         string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash  string         `gorm:"size:255;not null" json:"-"`
	Nickname      string         `gorm:"size:100" json:"nickname"`
	AvatarURL     string         `gorm:"size:500" json:"avatar_url"`
	Bio           string         `gorm:"type:text" json:"bio"`
	Role          string         `gorm:"size:20;not null;default:'user'" json:"role"`
	Status        string         `gorm:"size:20;not null;default:'active'" json:"status"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

