package model

import (
	"time"
	"gorm.io/gorm"
)

type Comment struct {
	ID         uint64         `gorm:"primaryKey" json:"id"`
	ArticleID  uint64         `gorm:"not null;index" json:"article_id"`
	UserID     uint64         `gorm:"not null;index" json:"user_id"`
	ParentID   *uint64        `gorm:"index" json:"parent_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	ContentHTML string        `gorm:"type:text" json:"content_html"`
	LikeCount  int            `gorm:"default:0" json:"like_count"`
	ReplyCount int            `gorm:"default:0" json:"reply_count"`
	Status     string         `gorm:"size:20;default:'published'" json:"status"`
	IPAddress  string         `gorm:"size:45" json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User     User       `gorm:"foreignKey:UserID" json:"user"`
	Article  Article    `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	Parent   *Comment   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies  []Comment  `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

