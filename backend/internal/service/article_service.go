package service

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/dto/response"
	"dbapp/internal/errors"
	"dbapp/internal/model"
	"dbapp/internal/repository"
	"github.com/gosimple/slug"
	"regexp"
	"strings"
	"time"
)

type ArticleService struct {
	articleRepo     *repository.ArticleRepository
	userRepo        *repository.UserRepository
	likeRepo        *repository.LikeRepository
	articleImageRepo *repository.ArticleImageRepository
}

func NewArticleService(
	articleRepo *repository.ArticleRepository,
	userRepo *repository.UserRepository,
	likeRepo *repository.LikeRepository,
	articleImageRepo *repository.ArticleImageRepository,
) *ArticleService {
	return &ArticleService{
		articleRepo:      articleRepo,
		userRepo:         userRepo,
		likeRepo:         likeRepo,
		articleImageRepo: articleImageRepo,
	}
}

func (s *ArticleService) Create(req *request.CreateArticleRequest, userID uint64) (*response.ArticleResponse, error) {
	// 生成slug
	articleSlug := slug.Make(req.Title)

	// 检查slug是否已存在
	existing, _ := s.articleRepo.GetBySlug(articleSlug)
	if existing != nil {
		// 如果slug已存在，添加时间戳
		articleSlug = articleSlug + "-" + time.Now().Format("20060102150405")
	}

	now := time.Now()
	article := &model.Article{
		Title:   req.Title,
		Slug:    articleSlug,
		Content: req.Content,
		Summary: req.Summary,
		CoverImageURL: req.CoverImageURL,
		AuthorID: userID,
		Status:  req.Status,
	}

	if req.Status == "published" {
		article.PublishedAt = &now
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, errors.NewInternalError("创建文章失败")
	}

	// 更新分类关联
	if len(req.CategoryIDs) > 0 {
		if err := s.articleRepo.UpdateCategories(article.ID, req.CategoryIDs); err != nil {
			// 记录错误但不影响文章创建
		}
	}

	// 更新标签关联
	if len(req.TagIDs) > 0 {
		if err := s.articleRepo.UpdateTags(article.ID, req.TagIDs); err != nil {
			// 记录错误但不影响文章创建
		}
	}

	// 提取文章内容中的图片并保存到数据库
	s.extractAndSaveImages(article.ID, article.Content)

	// 加载关联数据
	article, _ = s.articleRepo.GetByID(article.ID)
	return s.toResponse(article, userID), nil
}

func (s *ArticleService) GetByID(id uint64, userID uint64) (*response.ArticleResponse, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("文章不存在")
	}

	// 增加浏览次数
	go s.articleRepo.IncrementViewCount(id)

	resp := s.toResponse(article, userID)
	
	// 检查当前用户是否点赞
	if userID > 0 && s.likeRepo != nil {
		isLiked, _ := s.likeRepo.IsLikedByUser(userID, "article", id)
		resp.IsLiked = isLiked
	}

	return resp, nil
}

func (s *ArticleService) List(req *request.ListArticleRequest, userID uint64) (*response.ArticleListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	conditions := make(map[string]interface{})
	if req.Status != "" {
		conditions["status"] = req.Status
	} else {
		// 如果用户已登录，显示所有状态的文章；否则只显示已发布的
		if userID > 0 {
			// 已登录用户可以看到所有文章（包括自己的草稿）
			// 不设置status条件，显示所有状态
		} else {
			conditions["status"] = "published" // 未登录用户只显示已发布的
		}
	}
	if req.CategoryID > 0 {
		conditions["category_id"] = req.CategoryID
	}
	if req.TagID > 0 {
		conditions["tag_id"] = req.TagID
	}
	if req.AuthorID > 0 {
		conditions["author_id"] = req.AuthorID
	}
	if req.Keyword != "" {
		conditions["keyword"] = req.Keyword
	}
	if req.Sort != "" {
		conditions["sort"] = req.Sort
	} else {
		conditions["sort"] = "created_at"
	}
	if req.Order != "" {
		conditions["order"] = req.Order
	} else {
		conditions["order"] = "DESC"
	}

	articles, total, err := s.articleRepo.List(req.Page, req.PageSize, conditions)
	if err != nil {
		return nil, errors.NewInternalError("查询文章列表失败")
	}

	items := make([]*response.ArticleResponse, len(articles))
	for i, article := range articles {
		resp := s.toResponse(&article, userID)
		// 检查当前用户是否点赞
		if userID > 0 && s.likeRepo != nil {
			isLiked, _ := s.likeRepo.IsLikedByUser(userID, "article", article.ID)
			resp.IsLiked = isLiked
		}
		items[i] = resp
	}

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &response.ArticleListResponse{
		Items: items,
		Pagination: response.Pagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *ArticleService) Update(id uint64, req *request.UpdateArticleRequest, userID uint64) (*response.ArticleResponse, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("文章不存在")
	}

	// 检查权限
	if article.AuthorID != userID {
		return nil, errors.NewForbiddenError("无权限修改此文章")
	}

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
		article.Slug = slug.Make(req.Title)
	}
	if req.Content != "" {
		article.Content = req.Content
		// 更新文章内容时，重新提取图片
		s.extractAndSaveImages(article.ID, article.Content)
	}
	if req.Summary != "" {
		article.Summary = req.Summary
	}
	if req.CoverImageURL != "" {
		article.CoverImageURL = req.CoverImageURL
	}
	if req.Status != "" {
		article.Status = req.Status
		if req.Status == "published" && article.PublishedAt == nil {
			now := time.Now()
			article.PublishedAt = &now
		}
	}

	article.EditCount++
	editorID := userID
	article.EditorID = &editorID

	if err := s.articleRepo.Update(article); err != nil {
		return nil, errors.NewInternalError("更新文章失败")
	}

	// 更新分类关联（即使是空数组也要更新，表示清除所有分类）
	if req.CategoryIDs != nil {
		if err := s.articleRepo.UpdateCategories(article.ID, req.CategoryIDs); err != nil {
			// 记录错误但不影响文章更新
		}
	}

	// 更新标签关联（即使是空数组也要更新，表示清除所有标签）
	if req.TagIDs != nil {
		if err := s.articleRepo.UpdateTags(article.ID, req.TagIDs); err != nil {
			// 记录错误但不影响文章更新
		}
	}

	article, _ = s.articleRepo.GetByID(id)
	return s.toResponse(article, userID), nil
}

func (s *ArticleService) Delete(id uint64, userID uint64) error {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return errors.NewNotFoundError("文章不存在")
	}

	// 检查权限
	if article.AuthorID != userID {
		return errors.NewForbiddenError("无权限删除此文章")
	}

	if err := s.articleRepo.Delete(id); err != nil {
		return errors.NewInternalError("删除文章失败")
	}

	return nil
}

func (s *ArticleService) toResponse(article *model.Article, userID uint64) *response.ArticleResponse {
	author := &response.UserResponse{
		ID:        article.Author.ID,
		Username:  article.Author.Username,
		Nickname:  article.Author.Nickname,
		AvatarURL: article.Author.AvatarURL,
	}

	var editor *response.UserResponse
	if article.Editor != nil {
		editor = &response.UserResponse{
			ID:        article.Editor.ID,
			Username:  article.Editor.Username,
			Nickname:  article.Editor.Nickname,
			AvatarURL: article.Editor.AvatarURL,
		}
	}

	categories := make([]response.CategoryResponse, len(article.Categories))
	for i, cat := range article.Categories {
		categories[i] = response.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
			Slug: cat.Slug,
		}
	}

	tags := make([]response.TagResponse, len(article.Tags))
	for i, tag := range article.Tags {
		tags[i] = response.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		}
	}

	return &response.ArticleResponse{
		ID:            article.ID,
		Title:         article.Title,
		Slug:          article.Slug,
		Content:       article.Content,
		ContentHTML:   article.ContentHTML,
		Summary:       article.Summary,
		CoverImageURL: article.CoverImageURL,
		Author:        author,
		Editor:        editor,
		Categories:    categories,
		Tags:          tags,
		ViewCount:     article.ViewCount,
		LikeCount:     article.LikeCount,
		CommentCount:  article.CommentCount,
		IsFeatured:    article.IsFeatured,
		Status:        article.Status,
		PublishedAt:   article.PublishedAt,
		CreatedAt:     article.CreatedAt,
		UpdatedAt:     article.UpdatedAt,
	}
}

// extractAndSaveImages 从Markdown内容中提取图片URL并保存到数据库
func (s *ArticleService) extractAndSaveImages(articleID uint64, content string) {
	if s.articleImageRepo == nil || content == "" {
		return
	}

	// 匹配Markdown图片语法: ![alt](url) 或 <img src="url" alt="alt" />
	// 也匹配HTML img标签
	imgRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)|<img[^>]+src=["']([^"']+)["'][^>]*>`)
	matches := imgRegex.FindAllStringSubmatch(content, -1)

	// 收集所有图片URL
	imageURLs := make(map[string]bool)
	for _, match := range matches {
		var url string
		if len(match) > 2 && match[2] != "" {
			// Markdown格式: ![alt](url)
			url = strings.TrimSpace(match[2])
		} else if len(match) > 3 && match[3] != "" {
			// HTML格式: <img src="url" />
			url = strings.TrimSpace(match[3])
		}
		
		if url != "" && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			// 只处理相对路径的图片（上传到服务器的图片）
			imageURLs[url] = true
		}
	}

	// 获取现有的图片记录
	existingImages, _ := s.articleImageRepo.GetByArticleID(articleID)
	existingURLs := make(map[string]bool)
	for _, img := range existingImages {
		existingURLs[img.ImageURL] = true
	}

	// 添加新图片
	for url := range imageURLs {
		if !existingURLs[url] {
			// 检查是否已存在（可能其他文章也在使用）
			existing, _ := s.articleImageRepo.GetByURL(url)
			if existing == nil {
				// 创建新记录
				articleImage := &model.ArticleImage{
					ArticleID: articleID,
					ImageURL:  url,
				}
				s.articleImageRepo.Create(articleImage)
			}
		}
	}

	// 删除不再使用的图片记录（软删除）
	for _, img := range existingImages {
		if !imageURLs[img.ImageURL] {
			s.articleImageRepo.DeleteByID(img.ID)
		}
	}
}

