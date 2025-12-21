package handler

import (
	"dbapp/internal/errors"
	"dbapp/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type LikeHandler struct {
	likeService *service.LikeService
}

func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
}

// ToggleArticleLike 点赞/取消点赞文章
// @Summary 点赞/取消点赞文章
// @Tags 点赞
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/articles/{id}/like [post]
func (h *LikeHandler) ToggleArticleLike(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的文章ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint64)

	isLiked, err := h.likeService.ToggleLike(userIDUint, "article", id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	action := "取消点赞"
	if isLiked {
		action = "点赞"
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": action + "成功",
		"data": gin.H{
			"is_liked": isLiked,
		},
	})
}

// ToggleCommentLike 点赞/取消点赞评论
// @Summary 点赞/取消点赞评论
// @Tags 点赞
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/comments/{id}/like [post]
func (h *LikeHandler) ToggleCommentLike(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的评论ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint64)

	isLiked, err := h.likeService.ToggleLike(userIDUint, "comment", id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	action := "取消点赞"
	if isLiked {
		action = "点赞"
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": action + "成功",
		"data": gin.H{
			"is_liked": isLiked,
		},
	})
}

