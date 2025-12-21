package repository

import (
	"dbapp/internal/model"
	"gorm.io/gorm"
)

type LikeRepository struct {
	*BaseRepository
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *LikeRepository) Create(like *model.Like) error {
	return r.db.Create(like).Error
}

func (r *LikeRepository) GetByUserAndTarget(userID uint64, targetType string, targetID uint64) (*model.Like, error) {
	var like model.Like
	err := r.db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		First(&like).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &like, err
}

func (r *LikeRepository) Delete(userID uint64, targetType string, targetID uint64) error {
	return r.db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Delete(&model.Like{}).Error
}

func (r *LikeRepository) CountByTarget(targetType string, targetID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Like{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error
	return count, err
}

func (r *LikeRepository) IsLikedByUser(userID uint64, targetType string, targetID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Like{}).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Count(&count).Error
	return count > 0, err
}

