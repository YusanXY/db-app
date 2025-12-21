package model

import (
	"time"
	"gorm.io/gorm"
)

type Article struct {
	ID            uint64         `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"size:500;not null" json:"title"`
	Slug          string         `gorm:"uniqueIndex;size:500;not null" json:"slug"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	ContentHTML   string         `gorm:"type:text" json:"content_html"`
	Summary       string         `gorm:"type:text" json:"summary"`
	CoverImageURL string         `gorm:"size:500" json:"cover_image_url"`
	AuthorID      uint64         `gorm:"not null;index" json:"author_id"`
	EditorID      *uint64        `gorm:"index" json:"editor_id,omitempty"`
	Status        string         `gorm:"size:20;not null;default:'draft'" json:"status"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	LikeCount     int            `gorm:"default:0" json:"like_count"`
	CommentCount  int            `gorm:"default:0" json:"comment_count"`
	EditCount     int            `gorm:"default:0" json:"edit_count"`
	IsFeatured    bool           `gorm:"default:false" json:"is_featured"`
	IsLocked      bool           `gorm:"default:false" json:"is_locked"`
	PublishedAt   *time.Time     `json:"published_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Author     User       `gorm:"foreignKey:AuthorID" json:"author"`
	Editor     *User      `gorm:"foreignKey:EditorID" json:"editor,omitempty"`
	Categories []Category `gorm:"many2many:article_categories;" json:"categories"`
	Tags       []Tag      `gorm:"many2many:article_tags;" json:"tags"`
}

func (Article) TableName() string {
	return "articles"
}

