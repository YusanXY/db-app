package service

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/dto/response"
	"dbapp/internal/errors"
	"dbapp/internal/model"
	"dbapp/internal/repository"
	"html"
	"strings"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
	articleRepo *repository.ArticleRepository
	likeRepo    *repository.LikeRepository
}

func NewCommentService(
	commentRepo *repository.CommentRepository,
	articleRepo *repository.ArticleRepository,
	likeRepo *repository.LikeRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
		likeRepo:    likeRepo,
	}
}

func (s *CommentService) Create(req *request.CreateCommentRequest, userID uint64) (*response.CommentResponse, error) {
	// 验证文章是否存在
	article, err := s.articleRepo.GetByID(req.ArticleID)
	if err != nil || article == nil {
		return nil, errors.NewNotFoundError("文章不存在")
	}

	// 如果是指定父评论的回复，验证父评论是否存在
	if req.ParentID != nil && *req.ParentID > 0 {
		parent, err := s.commentRepo.GetByID(*req.ParentID)
		if err != nil || parent == nil {
			return nil, errors.NewNotFoundError("父评论不存在")
		}
		// 确保父评论属于同一篇文章
		if parent.ArticleID != req.ArticleID {
			return nil, errors.NewBadRequestError("父评论不属于该文章")
		}
	}

	// 转义HTML，防止XSS
	contentHTML := html.EscapeString(req.Content)
	// 将换行转换为<br>
	contentHTML = strings.ReplaceAll(contentHTML, "\n", "<br>")

	comment := &model.Comment{
		ArticleID:   req.ArticleID,
		UserID:      userID,
		ParentID:    req.ParentID,
		Content:     req.Content,
		ContentHTML: contentHTML,
		Status:      "published",
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, errors.NewInternalError("创建评论失败")
	}

	// 如果是回复，增加父评论的回复计数
	if req.ParentID != nil && *req.ParentID > 0 {
		go s.commentRepo.IncrementReplyCount(*req.ParentID)
	}

	// 增加文章的评论计数
	go s.articleRepo.IncrementCommentCount(req.ArticleID)

	// 重新加载评论以获取关联数据
	comment, _ = s.commentRepo.GetByID(comment.ID)
	return s.toResponse(comment, userID), nil
}

func (s *CommentService) GetByID(id uint64, userID uint64) (*response.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("评论不存在")
	}
	return s.toResponse(comment, userID), nil
}

func (s *CommentService) ListByArticle(articleID uint64, req *request.ListCommentRequest, userID uint64) (*response.CommentListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	comments, total, err := s.commentRepo.ListByArticle(articleID, req.Page, req.PageSize)
	if err != nil {
		return nil, errors.NewInternalError("查询评论列表失败")
	}

	items := make([]response.CommentResponse, len(comments))
	for i, comment := range comments {
		items[i] = *s.toResponse(&comment, userID)
	}

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &response.CommentListResponse{
		Items: items,
		Pagination: response.Pagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *CommentService) Update(id uint64, req *request.UpdateCommentRequest, userID uint64) (*response.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("评论不存在")
	}

	// 检查权限：只能修改自己的评论
	if comment.UserID != userID {
		return nil, errors.NewForbiddenError("无权限修改此评论")
	}

	// 转义HTML
	contentHTML := html.EscapeString(req.Content)
	contentHTML = strings.ReplaceAll(contentHTML, "\n", "<br>")

	comment.Content = req.Content
	comment.ContentHTML = contentHTML

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, errors.NewInternalError("更新评论失败")
	}

	comment, _ = s.commentRepo.GetByID(id)
	return s.toResponse(comment, userID), nil
}

func (s *CommentService) Delete(id uint64, userID uint64) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return errors.NewNotFoundError("评论不存在")
	}

	// 检查权限：只能删除自己的评论（管理员可以删除任何评论，这里简化处理）
	if comment.UserID != userID {
		return errors.NewForbiddenError("无权限删除此评论")
	}

	// 如果是回复，减少父评论的回复计数
	if comment.ParentID != nil && *comment.ParentID > 0 {
		go s.commentRepo.DecrementReplyCount(*comment.ParentID)
	}

	// 减少文章的评论计数
	go s.articleRepo.DecrementCommentCount(comment.ArticleID)

	if err := s.commentRepo.Delete(id); err != nil {
		return errors.NewInternalError("删除评论失败")
	}

	return nil
}

func (s *CommentService) toResponse(comment *model.Comment, userID uint64) *response.CommentResponse {
	resp := &response.CommentResponse{
		ID:          comment.ID,
		ArticleID:   comment.ArticleID,
		Content:     comment.Content,
		ContentHTML: comment.ContentHTML,
		User:        *toUserResponse(&comment.User),
		ParentID:    comment.ParentID,
		LikeCount:   comment.LikeCount,
		ReplyCount:  comment.ReplyCount,
		Status:      comment.Status,
		CreatedAt:   comment.CreatedAt,
		UpdatedAt:   comment.UpdatedAt,
	}

	// 检查当前用户是否点赞
	if userID > 0 {
		isLiked, _ := s.likeRepo.IsLikedByUser(userID, "comment", comment.ID)
		resp.IsLiked = isLiked
	}

	// 处理父评论
	if comment.Parent != nil {
		parentResp := s.toResponse(comment.Parent, userID)
		resp.Parent = parentResp
	}

	// 处理回复列表
	if len(comment.Replies) > 0 {
		replies := make([]response.CommentResponse, len(comment.Replies))
		for i, reply := range comment.Replies {
			replies[i] = *s.toResponse(&reply, userID)
		}
		resp.Replies = replies
	}

	return resp
}

func toUserResponse(user *model.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

