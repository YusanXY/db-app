package handler

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/errors"
	"dbapp/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type TagHandler struct {
	tagService *service.TagService
}

func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// GetTagList 获取标签列表
// @Summary 获取标签列表
// @Tags 标签
// @Accept json
// @Produce json
// @Param keyword query string false "关键词"
// @Param sort query string false "排序字段" default(article_count)
// @Param order query string false "排序方向" default(desc)
// @Param limit query int false "限制数量"
// @Success 200 {array} response.TagResponse
// @Router /api/v1/tags [get]
func (h *TagHandler) GetTagList(c *gin.Context) {
	var req request.ListTagRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误"))
		return
	}

	tags, err := h.tagService.List(&req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": tags,
	})
}

// GetTagDetail 获取标签详情
// @Summary 获取标签详情
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} response.TagResponse
// @Router /api/v1/tags/{id} [get]
func (h *TagHandler) GetTagDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的标签ID"))
		return
	}

	tag, err := h.tagService.GetByID(id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": tag,
	})
}

// GetTagBySlug 通过slug获取标签详情
// @Summary 通过slug获取标签详情
// @Tags 标签
// @Accept json
// @Produce json
// @Param slug path string true "标签slug"
// @Success 200 {object} response.TagResponse
// @Router /api/v1/tags/slug/{slug} [get]
func (h *TagHandler) GetTagBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		errors.HandleError(c, errors.NewBadRequestError("无效的标签slug"))
		return
	}

	tag, err := h.tagService.GetBySlug(slug)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": tag,
	})
}

// CreateTag 创建标签
// @Summary 创建标签
// @Tags 标签
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tag body request.CreateTagRequest true "标签信息"
// @Success 201 {object} response.TagResponse
// @Router /api/v1/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req request.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误: "+err.Error()))
		return
	}

	tag, err := h.tagService.Create(&req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(201, gin.H{
		"code":    201,
		"message": "创建成功",
		"data":    tag,
	})
}

// UpdateTag 更新标签
// @Summary 更新标签
// @Tags 标签
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "标签ID"
// @Param tag body request.UpdateTagRequest true "标签信息"
// @Success 200 {object} response.TagResponse
// @Router /api/v1/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的标签ID"))
		return
	}

	var req request.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误"))
		return
	}

	tag, err := h.tagService.Update(id, &req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    tag,
	})
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Tags 标签
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "标签ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("无效的标签ID"))
		return
	}

	if err := h.tagService.Delete(id); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

