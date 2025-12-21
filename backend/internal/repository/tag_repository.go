package repository

import (
	"dbapp/internal/model"
	"gorm.io/gorm"
)

type TagRepository struct {
	*BaseRepository
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *TagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepository) GetByID(id uint64) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.First(&tag, id).Error
	return &tag, err
}

func (r *TagRepository) GetBySlug(slug string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tag, err
}

func (r *TagRepository) GetByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tag, err
}

func (r *TagRepository) List(keyword string, sort string, order string, limit int) ([]model.Tag, error) {
	var tags []model.Tag
	query := r.db.Model(&model.Tag{})

	if keyword != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 排序
	if sort == "" {
		sort = "article_count"
	}
	if order == "" {
		order = "DESC"
	}
	query = query.Order(sort + " " + order)

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&tags).Error
	return tags, err
}

func (r *TagRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

func (r *TagRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

func (r *TagRepository) IncrementArticleCount(id uint64) error {
	return r.db.Model(&model.Tag{}).Where("id = ?", id).
		UpdateColumn("article_count", gorm.Expr("article_count + 1")).Error
}

func (r *TagRepository) DecrementArticleCount(id uint64) error {
	return r.db.Model(&model.Tag{}).Where("id = ?", id).
		UpdateColumn("article_count", gorm.Expr("GREATEST(article_count - 1, 0)")).Error
}

