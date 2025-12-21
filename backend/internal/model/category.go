package model

import (
	"time"
)

type Category struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Slug        string     `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Description string     `gorm:"type:text" json:"description"`
	ParentID    *uint64    `gorm:"index" json:"parent_id"`
	IconURL     string     `gorm:"size:500" json:"icon_url"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	ArticleCount int       `gorm:"default:0" json:"article_count"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联
	Parent   *Category   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category   `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}

