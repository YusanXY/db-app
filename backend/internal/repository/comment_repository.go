package repository

import (
	"dbapp/internal/model"
	"gorm.io/gorm"
)

type CommentRepository struct {
	*BaseRepository
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) GetByID(id uint64) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("User").Preload("Parent").Preload("Replies.User").
		First(&comment, id).Error
	return &comment, err
}

func (r *CommentRepository) ListByArticle(articleID uint64, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	// 只查询顶级评论（parent_id为NULL）
	query := r.db.Model(&model.Comment{}).
		Where("article_id = ? AND parent_id IS NULL AND status = ?", articleID, "published")

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "published").Preload("User").Order("created_at ASC")
		}).
		Scopes(r.Paginate(page, pageSize)).
		Order("created_at DESC").
		Find(&comments).Error

	return comments, total, err
}

func (r *CommentRepository) ListByUser(userID uint64, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := r.db.Model(&model.Comment{}).
		Where("user_id = ? AND status = ?", userID, "published")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Article").
		Scopes(r.Paginate(page, pageSize)).
		Order("created_at DESC").
		Find(&comments).Error

	return comments, total, err
}

func (r *CommentRepository) Update(comment *model.Comment) error {
	return r.db.Save(comment).Error
}

func (r *CommentRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Comment{}, id).Error
}

func (r *CommentRepository) IncrementLikeCount(id uint64) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *CommentRepository) DecrementLikeCount(id uint64) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

func (r *CommentRepository) IncrementReplyCount(id uint64) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).
		UpdateColumn("reply_count", gorm.Expr("reply_count + 1")).Error
}

func (r *CommentRepository) DecrementReplyCount(id uint64) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).
		UpdateColumn("reply_count", gorm.Expr("GREATEST(reply_count - 1, 0)")).Error
}

func (r *CommentRepository) CountByArticle(articleID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).
		Where("article_id = ? AND status = ?", articleID, "published").
		Count(&count).Error
	return count, err
}

