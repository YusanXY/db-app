package handler

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/errors"
	"dbapp/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// GetCommentList 获取文章评论列表
// @Summary 获取文章评论列表
// @Tags 评论
// @Accept json
// @Produce json
// @Param article_id path int true "文章ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.CommentListResponse
// @Router /api/v1/articles/{article_id}/comments [get]
func (h *CommentHandler) GetCommentList(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的文章ID"))
		return
	}

	var req request.ListCommentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误"))
		return
	}

	var userID uint64
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint64)
	}

	result, err := h.commentService.ListByArticle(articleID, &req, userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": result,
	})
}

// CreateComment 创建评论
// @Summary 创建评论
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param article_id path int true "文章ID"
// @Param comment body request.CreateCommentRequest true "评论信息"
// @Success 201 {object} response.CommentResponse
// @Router /api/v1/articles/{article_id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的文章ID"))
		return
	}

	var req request.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误: "+err.Error()))
		return
	}

	// 从路径参数设置article_id（前端不需要在请求体中发送）
	req.ArticleID = articleID

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint64)

	comment, err := h.commentService.Create(&req, userIDUint)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(201, gin.H{
		"code":    201,
		"message": "评论成功",
		"data":    comment,
	})
}

// UpdateComment 更新评论
// @Summary 更新评论
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Param comment body request.UpdateCommentRequest true "评论信息"
// @Success 200 {object} response.CommentResponse
// @Router /api/v1/comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的评论ID"))
		return
	}

	var req request.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误"))
		return
	}

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint64)

	comment, err := h.commentService.Update(id, &req, userIDUint)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    comment,
	})
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的评论ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint64)

	if err := h.commentService.Delete(id, userIDUint); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

