package model

import (
	"time"
)

type Like struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `gorm:"not null;index" json:"user_id"`
	TargetType string    `gorm:"size:20;not null;index" json:"target_type"` // article, comment
	TargetID   uint64    `gorm:"not null;index" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user"`
}

func (Like) TableName() string {
	return "likes"
}

