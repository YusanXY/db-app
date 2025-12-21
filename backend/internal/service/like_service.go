package service

import (
	"dbapp/internal/errors"
	"dbapp/internal/model"
	"dbapp/internal/repository"
)

type LikeService struct {
	likeRepo    *repository.LikeRepository
	articleRepo *repository.ArticleRepository
	commentRepo *repository.CommentRepository
}

func NewLikeService(
	likeRepo *repository.LikeRepository,
	articleRepo *repository.ArticleRepository,
	commentRepo *repository.CommentRepository,
) *LikeService {
	return &LikeService{
		likeRepo:    likeRepo,
		articleRepo: articleRepo,
		commentRepo: commentRepo,
	}
}

func (s *LikeService) ToggleLike(userID uint64, targetType string, targetID uint64) (bool, error) {
	// 检查是否已点赞
	existing, err := s.likeRepo.GetByUserAndTarget(userID, targetType, targetID)
	if err != nil {
		return false, errors.NewInternalError("查询点赞状态失败")
	}

	if existing != nil {
		// 已点赞，取消点赞
		if err := s.likeRepo.Delete(userID, targetType, targetID); err != nil {
			return false, errors.NewInternalError("取消点赞失败")
		}

		// 减少计数
		if targetType == "article" {
			go s.articleRepo.DecrementLikeCount(targetID)
		} else if targetType == "comment" {
			go s.commentRepo.DecrementLikeCount(targetID)
		}

		return false, nil
	} else {
		// 未点赞，添加点赞
		like := &model.Like{
			UserID:     userID,
			TargetType: targetType,
			TargetID:   targetID,
		}

		if err := s.likeRepo.Create(like); err != nil {
			return false, errors.NewInternalError("点赞失败")
		}

		// 增加计数
		if targetType == "article" {
			go s.articleRepo.IncrementLikeCount(targetID)
		} else if targetType == "comment" {
			go s.commentRepo.IncrementLikeCount(targetID)
		}

		return true, nil
	}
}

func (s *LikeService) IsLiked(userID uint64, targetType string, targetID uint64) (bool, error) {
	return s.likeRepo.IsLikedByUser(userID, targetType, targetID)
}

