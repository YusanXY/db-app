package handler

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/errors"
	"dbapp/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// GetArticleList 获取文章列表
// @Summary 获取文章列表
// @Tags 文章
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query string false "状态"
// @Param category_id query int false "分类ID"
// @Success 200 {object} response.ArticleListResponse
// @Router /api/v1/articles [get]
func (h *ArticleHandler) GetArticleList(c *gin.Context) {
	var req request.ListArticleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误"))
		return
	}

	var userID uint64
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint64)
	}

	result, err := h.articleService.List(&req, userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": result,
	})
}

// GetArticleDetail 获取文章详情
// @Summary 获取文章详情
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.ArticleResponse
// @Router /api/v1/articles/{id} [get]
func (h *ArticleHandler) GetArticleDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的文章ID"))
		return
	}

	var userID uint64
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint64)
	}

	article, err := h.articleService.GetByID(id, userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": article,
	})
}

// CreateArticle 创建文章
// @Summary 创建文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param article body request.CreateArticleRequest true "文章信息"
// @Success 201 {object} response.ArticleResponse
// @Router /api/v1/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req request.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误"))
		return
	}

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint64)

	article, err := h.articleService.Create(&req, userIDUint)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(201, gin.H{
		"code":    201,
		"message": "创建成功",
		"data":    article,
	})
}

// UpdateArticle 更新文章
// @Summary 更新文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param article body request.UpdateArticleRequest true "文章信息"
// @Success 200 {object} response.ArticleResponse
// @Router /api/v1/articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的文章ID"))
		return
	}

	var req request.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误"))
		return
	}

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint64)

	article, err := h.articleService.Update(id, &req, userIDUint)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    article,
	})
}

// DeleteArticle 删除文章
// @Summary 删除文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的文章ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint64)

	if err := h.articleService.Delete(id, userIDUint); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

