package service

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/dto/response"
	"dbapp/internal/errors"
	"dbapp/internal/model"
	"dbapp/internal/repository"
	"github.com/gosimple/slug"
)

type TagService struct {
	tagRepo *repository.TagRepository
}

func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

func (s *TagService) Create(req *request.CreateTagRequest) (*response.TagResponse, error) {
	// 如果未提供slug，自动生成
	tagSlug := req.Slug
	if tagSlug == "" {
		tagSlug = slug.Make(req.Name)
	}

	// 检查slug是否已存在
	existing, _ := s.tagRepo.GetBySlug(tagSlug)
	if existing != nil {
		return nil, errors.NewBadRequestError("标签标识已存在")
	}

	// 检查名称是否已存在
	existing, _ = s.tagRepo.GetByName(req.Name)
	if existing != nil {
		return nil, errors.NewBadRequestError("标签名称已存在")
	}

	tag := &model.Tag{
		Name:        req.Name,
		Slug:        tagSlug,
		Description: req.Description,
		Color:       req.Color,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, errors.NewInternalError("创建标签失败")
	}

	tag, _ = s.tagRepo.GetByID(tag.ID)
	return s.toResponse(tag), nil
}

func (s *TagService) GetByID(id uint64) (*response.TagResponse, error) {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("标签不存在")
	}
	return s.toResponse(tag), nil
}

func (s *TagService) GetBySlug(slug string) (*response.TagResponse, error) {
	tag, err := s.tagRepo.GetBySlug(slug)
	if err != nil || tag == nil {
		return nil, errors.NewNotFoundError("标签不存在")
	}
	return s.toResponse(tag), nil
}

func (s *TagService) List(req *request.ListTagRequest) ([]response.TagResponse, error) {
	tags, err := s.tagRepo.List(req.Keyword, req.Sort, req.Order, req.Limit)
	if err != nil {
		return nil, errors.NewInternalError("查询标签列表失败")
	}

	items := make([]response.TagResponse, len(tags))
	for i, tag := range tags {
		items[i] = *s.toResponse(&tag)
	}

	return items, nil
}

func (s *TagService) Update(id uint64, req *request.UpdateTagRequest) (*response.TagResponse, error) {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("标签不存在")
	}

	// 更新字段
	if req.Name != "" {
		// 检查新名称是否与其他标签冲突
		existing, _ := s.tagRepo.GetByName(req.Name)
		if existing != nil && existing.ID != id {
			return nil, errors.NewBadRequestError("标签名称已存在")
		}
		tag.Name = req.Name
	}
	if req.Description != "" {
		tag.Description = req.Description
	}
	if req.Color != "" {
		tag.Color = req.Color
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, errors.NewInternalError("更新标签失败")
	}

	tag, _ = s.tagRepo.GetByID(id)
	return s.toResponse(tag), nil
}

func (s *TagService) Delete(id uint64) error {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return errors.NewNotFoundError("标签不存在")
	}

	// 检查是否有关联文章
	if tag.ArticleCount > 0 {
		return errors.NewBadRequestError("标签下存在文章，无法删除")
	}

	if err := s.tagRepo.Delete(id); err != nil {
		return errors.NewInternalError("删除标签失败")
	}

	return nil
}

func (s *TagService) toResponse(tag *model.Tag) *response.TagResponse {
	return &response.TagResponse{
		ID:           tag.ID,
		Name:         tag.Name,
		Slug:         tag.Slug,
		Description:  tag.Description,
		Color:        tag.Color,
		ArticleCount: tag.ArticleCount,
		CreatedAt:    tag.CreatedAt,
	}
}

