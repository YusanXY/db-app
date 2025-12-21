package model

import (
	"time"
	"gorm.io/gorm"
)

// ArticleImage 文章图片表，用于追踪文章内容中使用的图片
// 便于后续的图片管理和清理
type ArticleImage struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	ArticleID uint64         `gorm:"not null;index" json:"article_id"`
	ImageURL  string         `gorm:"size:500;not null" json:"image_url"`
	ImagePath string         `gorm:"size:500" json:"image_path"` // 服务器上的文件路径
	FileSize  int64          `gorm:"default:0" json:"file_size"`  // 文件大小（字节）
	MimeType  string         `gorm:"size:100" json:"mime_type"`   // 文件类型
	Width     int            `json:"width"`                        // 图片宽度
	Height    int            `json:"height"`                       // 图片高度
	Alt       string         `gorm:"size:500" json:"alt"`          // 图片描述
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Article Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
}

func (ArticleImage) TableName() string {
	return "article_images"
}

