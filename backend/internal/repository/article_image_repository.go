package repository

import (
	"dbapp/internal/model"
	"gorm.io/gorm"
)

type ArticleImageRepository struct {
	*BaseRepository
}

func NewArticleImageRepository(db *gorm.DB) *ArticleImageRepository {
	return &ArticleImageRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建文章图片记录
func (r *ArticleImageRepository) Create(image *model.ArticleImage) error {
	return r.db.Create(image).Error
}

// GetByArticleID 获取文章的所有图片
func (r *ArticleImageRepository) GetByArticleID(articleID uint64) ([]*model.ArticleImage, error) {
	var images []*model.ArticleImage
	err := r.db.Where("article_id = ?", articleID).Find(&images).Error
	return images, err
}

// GetByURL 根据URL获取图片记录
func (r *ArticleImageRepository) GetByURL(url string) (*model.ArticleImage, error) {
	var image model.ArticleImage
	err := r.db.Where("image_url = ?", url).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// DeleteByArticleID 删除文章的所有图片记录（软删除）
func (r *ArticleImageRepository) DeleteByArticleID(articleID uint64) error {
	return r.db.Where("article_id = ?", articleID).Delete(&model.ArticleImage{}).Error
}

// DeleteByID 删除指定图片记录
func (r *ArticleImageRepository) DeleteByID(id uint64) error {
	return r.db.Delete(&model.ArticleImage{}, id).Error
}

// CountByArticleID 统计文章使用的图片数量
func (r *ArticleImageRepository) CountByArticleID(articleID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.ArticleImage{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

