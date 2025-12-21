package handler

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/errors"
	"dbapp/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetCategoryList 获取分类列表
// @Summary 获取分类列表
// @Tags 分类
// @Accept json
// @Produce json
// @Param parent_id query int false "父分类ID"
// @Param is_active query bool false "是否启用"
// @Param tree query bool false "是否返回树形结构"
// @Success 200 {array} response.CategoryResponse
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetCategoryList(c *gin.Context) {
	var req request.ListCategoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误"))
		return
	}

	categories, err := h.categoryService.List(&req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": categories,
	})
}

// GetCategoryDetail 获取分类详情
// @Summary 获取分类详情
// @Tags 分类
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.CategoryResponse
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategoryDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的分类ID"))
		return
	}

	category, err := h.categoryService.GetByID(id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": category,
	})
}

// GetCategoryBySlug 通过slug获取分类详情
// @Summary 通过slug获取分类详情
// @Tags 分类
// @Accept json
// @Produce json
// @Param slug path string true "分类slug"
// @Success 200 {object} response.CategoryResponse
// @Router /api/v1/categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		errors.HandleError(c, errors.NewBadRequestError("无效的分类slug"))
		return
	}

	category, err := h.categoryService.GetBySlug(slug)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": category,
	})
}

// CreateCategory 创建分类
// @Summary 创建分类
// @Tags 分类
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body request.CreateCategoryRequest true "分类信息"
// @Success 201 {object} response.CategoryResponse
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误: "+err.Error()))
		return
	}

	category, err := h.categoryService.Create(&req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(201, gin.H{
		"code":    201,
		"message": "创建成功",
		"data":    category,
	})
}

// UpdateCategory 更新分类
// @Summary 更新分类
// @Tags 分类
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Param category body request.UpdateCategoryRequest true "分类信息"
// @Success 200 {object} response.CategoryResponse
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的分类ID"))
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误"))
		return
	}

	category, err := h.categoryService.Update(id, &req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    category,
	})
}

// DeleteCategory 删除分类
// @Summary 删除分类
// @Tags 分类
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的分类ID"))
		return
	}

	if err := h.categoryService.Delete(id); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

