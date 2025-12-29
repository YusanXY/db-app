package repository

import (
	"dbapp/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	*BaseRepository
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) GetByID(id uint64) (*model.Category, error) {
	var category model.Category
	err := r.db.Preload("Parent").Preload("Children").First(&category, id).Error
	if err != nil {
		return &category, err
	}
	// 动态计算文章数
	var count int64
	r.db.Table("article_categories").Where("category_id = ?", category.ID).Count(&count)
	category.ArticleCount = int(count)
	return &category, nil
}

func (r *CategoryRepository) GetBySlug(slug string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("slug = ?", slug).
		Preload("Parent").Preload("Children").
		First(&category).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}

func (r *CategoryRepository) List(parentID *uint64, isActive *bool) ([]model.Category, error) {
	var categories []model.Category
	query := r.db.Model(&model.Category{})

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Preload("Parent").
		Order("sort_order ASC, created_at ASC").
		Find(&categories).Error

	if err != nil {
		return categories, err
	}

	// 动态计算每个分类的文章数
	for i := range categories {
		var count int64
		r.db.Table("article_categories").Where("category_id = ?", categories[i].ID).Count(&count)
		categories[i].ArticleCount = int(count)
	}

	return categories, nil
}

func (r *CategoryRepository) GetTree() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("parent_id IS NULL AND is_active = ?", true).
		Preload("Children", "is_active = ?", true).
		Order("sort_order ASC, created_at ASC").
		Find(&categories).Error

	if err != nil {
		return categories, err
	}

	// 动态计算每个分类及其子分类的文章数
	for i := range categories {
		var count int64
		r.db.Table("article_categories").Where("category_id = ?", categories[i].ID).Count(&count)
		categories[i].ArticleCount = int(count)
		
		for j := range categories[i].Children {
			var childCount int64
			r.db.Table("article_categories").Where("category_id = ?", categories[i].Children[j].ID).Count(&childCount)
			categories[i].Children[j].ArticleCount = int(childCount)
		}
	}

	return categories, nil
}

func (r *CategoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *CategoryRepository) IncrementArticleCount(id uint64) error {
	return r.db.Model(&model.Category{}).Where("id = ?", id).
		UpdateColumn("article_count", gorm.Expr("article_count + 1")).Error
}

func (r *CategoryRepository) DecrementArticleCount(id uint64) error {
	return r.db.Model(&model.Category{}).Where("id = ?", id).
		UpdateColumn("article_count", gorm.Expr("GREATEST(article_count - 1, 0)")).Error
}

func (r *CategoryRepository) CountChildren(parentID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}

