package repository

import (
	"dbapp/internal/model"
	"gorm.io/gorm"
)

type ArticleRepository struct {
	*BaseRepository
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *ArticleRepository) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

func (r *ArticleRepository) GetByID(id uint64) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("Author").Preload("Editor").
		Preload("Categories").Preload("Tags").
		First(&article, id).Error
	return &article, err
}

func (r *ArticleRepository) GetBySlug(slug string) (*model.Article, error) {
	var article model.Article
	err := r.db.Where("slug = ?", slug).
		Preload("Author").Preload("Editor").
		Preload("Categories").Preload("Tags").
		First(&article).Error
	return &article, err
}

func (r *ArticleRepository) List(page, pageSize int, conditions map[string]interface{}) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.Model(&model.Article{})

	// 应用条件
	if status, ok := conditions["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if categoryID, ok := conditions["category_id"]; ok && categoryID.(uint64) > 0 {
		query = query.Joins("JOIN article_categories ON articles.id = article_categories.article_id").
			Where("article_categories.category_id = ?", categoryID)
	}
	if tagID, ok := conditions["tag_id"]; ok && tagID.(uint64) > 0 {
		query = query.Joins("JOIN article_tags ON articles.id = article_tags.article_id").
			Where("article_tags.tag_id = ?", tagID)
	}
	if authorID, ok := conditions["author_id"]; ok && authorID.(uint64) > 0 {
		query = query.Where("author_id = ?", authorID)
	}
	if keyword, ok := conditions["keyword"]; ok && keyword != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+keyword.(string)+"%", "%"+keyword.(string)+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	sort := "created_at"
	order := "DESC"
	if s, ok := conditions["sort"]; ok && s != "" {
		sort = s.(string)
	}
	if o, ok := conditions["order"]; ok && o != "" {
		order = o.(string)
	}

	// 分页查询
	err := query.Preload("Author").Preload("Categories").Preload("Tags").
		Scopes(r.Paginate(page, pageSize)).
		Order(sort + " " + order).
		Find(&articles).Error

	return articles, total, err
}

func (r *ArticleRepository) Update(article *model.Article) error {
	return r.db.Save(article).Error
}

func (r *ArticleRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Article{}, id).Error
}

func (r *ArticleRepository) IncrementViewCount(id uint64) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *ArticleRepository) IncrementLikeCount(id uint64) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *ArticleRepository) DecrementLikeCount(id uint64) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

func (r *ArticleRepository) IncrementCommentCount(id uint64) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
}

func (r *ArticleRepository) DecrementCommentCount(id uint64) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)")).Error
}

// UpdateCategories 更新文章的分类关联
func (r *ArticleRepository) UpdateCategories(articleID uint64, categoryIDs []uint64) error {
	article := &model.Article{ID: articleID}
	
	// 先清除现有的分类关联
	if err := r.db.Model(article).Association("Categories").Clear(); err != nil {
		return err
	}
	
	// 如果有新的分类ID，添加关联
	if len(categoryIDs) > 0 {
		var categories []model.Category
		if err := r.db.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
			return err
		}
		if err := r.db.Model(article).Association("Categories").Replace(categories); err != nil {
			return err
		}
	}
	
	return nil
}

// UpdateTags 更新文章的标签关联
func (r *ArticleRepository) UpdateTags(articleID uint64, tagIDs []uint64) error {
	article := &model.Article{ID: articleID}
	
	// 先清除现有的标签关联
	if err := r.db.Model(article).Association("Tags").Clear(); err != nil {
		return err
	}
	
	// 如果有新的标签ID，添加关联
	if len(tagIDs) > 0 {
		var tags []model.Tag
		if err := r.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return err
		}
		if err := r.db.Model(article).Association("Tags").Replace(tags); err != nil {
			return err
		}
	}
	
	return nil
}

