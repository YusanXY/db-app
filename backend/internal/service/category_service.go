package service

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/dto/response"
	"dbapp/internal/errors"
	"dbapp/internal/model"
	"dbapp/internal/repository"
	"github.com/gosimple/slug"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) Create(req *request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	// 如果未提供slug，自动生成
	categorySlug := req.Slug
	if categorySlug == "" {
		categorySlug = slug.Make(req.Name)
	}

	// 检查slug是否已存在
	existing, _ := s.categoryRepo.GetBySlug(categorySlug)
	if existing != nil {
		return nil, errors.NewBadRequestError("分类标识已存在")
	}

	// 如果指定了父分类，验证父分类是否存在
	if req.ParentID != nil && *req.ParentID > 0 {
		parent, err := s.categoryRepo.GetByID(*req.ParentID)
		if err != nil || parent == nil {
			return nil, errors.NewNotFoundError("父分类不存在")
		}
	}

	category := &model.Category{
		Name:        req.Name,
		Slug:        categorySlug,
		Description: req.Description,
		ParentID:    req.ParentID,
		IconURL:     req.IconURL,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, errors.NewInternalError("创建分类失败")
	}

	category, _ = s.categoryRepo.GetByID(category.ID)
	return s.toResponse(category), nil
}

func (s *CategoryService) GetByID(id uint64) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("分类不存在")
	}
	return s.toResponse(category), nil
}

func (s *CategoryService) GetBySlug(slug string) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.GetBySlug(slug)
	if err != nil || category == nil {
		return nil, errors.NewNotFoundError("分类不存在")
	}
	return s.toResponse(category), nil
}

func (s *CategoryService) List(req *request.ListCategoryRequest) ([]response.CategoryResponse, error) {
	var categories []model.Category
	var err error

	if req.Tree {
		categories, err = s.categoryRepo.GetTree()
	} else {
		categories, err = s.categoryRepo.List(req.ParentID, req.IsActive)
	}

	if err != nil {
		return nil, errors.NewInternalError("查询分类列表失败")
	}

	items := make([]response.CategoryResponse, len(categories))
	for i, category := range categories {
		items[i] = *s.toResponse(&category)
	}

	return items, nil
}

func (s *CategoryService) Update(id uint64, req *request.UpdateCategoryRequest) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("分类不存在")
	}

	// 更新字段
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil {
		// 验证父分类是否存在（如果指定了父分类）
		if *req.ParentID > 0 {
			parent, err := s.categoryRepo.GetByID(*req.ParentID)
			if err != nil || parent == nil {
				return nil, errors.NewNotFoundError("父分类不存在")
			}
			// 防止循环引用
			if *req.ParentID == id {
				return nil, errors.NewBadRequestError("不能将自己设为父分类")
			}
		}
		category.ParentID = req.ParentID
	}
	if req.IconURL != "" {
		category.IconURL = req.IconURL
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, errors.NewInternalError("更新分类失败")
	}

	category, _ = s.categoryRepo.GetByID(id)
	return s.toResponse(category), nil
}

func (s *CategoryService) Delete(id uint64) error {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return errors.NewNotFoundError("分类不存在")
	}

	// 检查是否有子分类
	childCount, err := s.categoryRepo.CountChildren(id)
	if err != nil {
		return errors.NewInternalError("查询子分类失败")
	}
	if childCount > 0 {
		return errors.NewBadRequestError("存在子分类，无法删除")
	}

	// 检查是否有关联文章
	if category.ArticleCount > 0 {
		return errors.NewBadRequestError("分类下存在文章，无法删除")
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return errors.NewInternalError("删除分类失败")
	}

	return nil
}

func (s *CategoryService) toResponse(category *model.Category) *response.CategoryResponse {
	resp := &response.CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		Description:  category.Description,
		ParentID:     category.ParentID,
		IconURL:      category.IconURL,
		SortOrder:    category.SortOrder,
		ArticleCount: category.ArticleCount,
		IsActive:     category.IsActive,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}

	if category.Parent != nil {
		parentResp := s.toResponse(category.Parent)
		resp.Parent = parentResp
	}

	if len(category.Children) > 0 {
		children := make([]response.CategoryResponse, len(category.Children))
		for i, child := range category.Children {
			children[i] = *s.toResponse(&child)
		}
		resp.Children = children
	}

	return resp
}

