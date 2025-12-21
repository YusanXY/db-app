package model

import (
	"time"
)

type Tag struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Slug        string    `gorm:"uniqueIndex;size:50;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Color       string    `gorm:"size:20" json:"color"`
	ArticleCount int      `gorm:"default:0" json:"article_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Tag) TableName() string {
	return "tags"
}

