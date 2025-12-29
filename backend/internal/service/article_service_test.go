package service

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/repository"
	"dbapp/internal/test"
	"testing"
)

func TestArticleService_Create(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)

	// 创建测试用户
	user := test.CreateTestUser(db, "testuser", "test@example.com")

	req := &request.CreateArticleRequest{
		Title:   "测试文章",
		Content: "这是测试内容",
		Summary: "测试摘要",
		Status:  "published",
	}

	article, err := articleService.Create(req, user.ID)
	if err != nil {
		t.Fatalf("创建文章失败: %v", err)
	}

	if article.Title != "测试文章" {
		t.Errorf("期望标题 测试文章, 得到 %s", article.Title)
	}

	if article.Author.ID != user.ID {
		t.Errorf("期望作者ID %d, 得到 %d", user.ID, article.Author.ID)
	}
}

func TestArticleService_GetByID(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	article := test.CreateTestArticle(db, user.ID, "测试文章")

	// 查询文章
	found, err := articleService.GetByID(article.ID, user.ID)
	if err != nil {
		t.Fatalf("查询文章失败: %v", err)
	}

	if found.Title != "测试文章" {
		t.Errorf("期望标题 测试文章, 得到 %s", found.Title)
	}
}

func TestArticleService_GetByID_NotFound(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)

	// 查询不存在的文章
	_, err := articleService.GetByID(99999, 1)
	if err == nil {
		t.Error("应该返回文章不存在的错误")
	}
}

func TestArticleService_List(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	test.CreateTestArticle(db, user.ID, "文章1")
	test.CreateTestArticle(db, user.ID, "文章2")
	test.CreateTestArticle(db, user.ID, "文章3")

	req := &request.ListArticleRequest{
		Page:     1,
		PageSize: 10,
	}

	result, err := articleService.List(req, user.ID)
	if err != nil {
		t.Fatalf("查询文章列表失败: %v", err)
	}

	if result.Pagination.Total != 3 {
		t.Errorf("期望总数 3, 得到 %d", result.Pagination.Total)
	}

	if len(result.Items) != 3 {
		t.Errorf("期望文章数量 3, 得到 %d", len(result.Items))
	}
}

func TestArticleService_ListWithPagination(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	for i := 1; i <= 15; i++ {
		test.CreateTestArticle(db, user.ID, "文章")
	}

	req := &request.ListArticleRequest{
		Page:     1,
		PageSize: 10,
	}

	result, err := articleService.List(req, user.ID)
	if err != nil {
		t.Fatalf("查询文章列表失败: %v", err)
	}

	if result.Pagination.Total != 15 {
		t.Errorf("期望总数 15, 得到 %d", result.Pagination.Total)
	}

	if len(result.Items) != 10 {
		t.Errorf("期望第一页文章数量 10, 得到 %d", len(result.Items))
	}

	if result.Pagination.TotalPages != 2 {
		t.Errorf("期望总页数 2, 得到 %d", result.Pagination.TotalPages)
	}
}

func TestArticleService_Update(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	article := test.CreateTestArticle(db, user.ID, "原始标题")

	req := &request.UpdateArticleRequest{
		Title:   "更新后的标题",
		Content: "更新后的内容",
	}

	updated, err := articleService.Update(article.ID, req, user.ID)
	if err != nil {
		t.Fatalf("更新文章失败: %v", err)
	}

	if updated.Title != "更新后的标题" {
		t.Errorf("期望标题 更新后的标题, 得到 %s", updated.Title)
	}
}

func TestArticleService_Update_Unauthorized(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)

	// 创建测试数据
	author := test.CreateTestUser(db, "author", "author@example.com")
	otherUser := test.CreateTestUser(db, "other", "other@example.com")
	article := test.CreateTestArticle(db, author.ID, "测试文章")

	req := &request.UpdateArticleRequest{
		Title: "尝试修改",
	}

	_, err := articleService.Update(article.ID, req, otherUser.ID)
	if err == nil {
		t.Error("应该返回无权限的错误")
	}
}

func TestArticleService_Delete(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	article := test.CreateTestArticle(db, user.ID, "测试文章")

	err := articleService.Delete(article.ID, user.ID)
	if err != nil {
		t.Fatalf("删除文章失败: %v", err)
	}

	// 验证删除
	_, err = articleService.GetByID(article.ID, user.ID)
	if err == nil {
		t.Error("文章应该已被删除")
	}
}

func TestArticleService_Delete_Unauthorized(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)

	// 创建测试数据
	author := test.CreateTestUser(db, "author", "author@example.com")
	otherUser := test.CreateTestUser(db, "other", "other@example.com")
	article := test.CreateTestArticle(db, author.ID, "测试文章")

	err := articleService.Delete(article.ID, otherUser.ID)
	if err == nil {
		t.Error("应该返回无权限的错误")
	}
}

