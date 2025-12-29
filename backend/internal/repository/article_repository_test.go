package repository

import (
	"dbapp/internal/model"
	"dbapp/internal/test"
	"testing"
)

func TestArticleRepository_Create(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewArticleRepository(db)

	// 创建测试用户
	user := test.CreateTestUser(db, "testuser", "test@example.com")

	article := &model.Article{
		Title:    "测试文章",
		Slug:     "test-article",
		Content:  "这是测试内容",
		Summary:  "测试摘要",
		AuthorID: user.ID,
		Status:   "published",
	}

	err := repo.Create(article)
	if err != nil {
		t.Fatalf("创建文章失败: %v", err)
	}

	if article.ID == 0 {
		t.Error("文章ID应该被设置")
	}
}

func TestArticleRepository_GetByID(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewArticleRepository(db)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	article := test.CreateTestArticle(db, user.ID, "测试文章")

	// 查询文章
	found, err := repo.GetByID(article.ID)
	if err != nil {
		t.Fatalf("查询文章失败: %v", err)
	}

	if found.Title != "测试文章" {
		t.Errorf("期望标题 测试文章, 得到 %s", found.Title)
	}

	if found.Author.ID != user.ID {
		t.Errorf("期望作者ID %d, 得到 %d", user.ID, found.Author.ID)
	}
}

func TestArticleRepository_GetBySlug(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewArticleRepository(db)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	article := &model.Article{
		Title:    "测试文章",
		Slug:     "test-slug",
		Content:  "测试内容",
		AuthorID: user.ID,
		Status:   "published",
	}
	db.Create(article)

	// 查询文章
	found, err := repo.GetBySlug("test-slug")
	if err != nil {
		t.Fatalf("查询文章失败: %v", err)
	}

	if found.Slug != "test-slug" {
		t.Errorf("期望slug test-slug, 得到 %s", found.Slug)
	}
}

func TestArticleRepository_List(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewArticleRepository(db)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	test.CreateTestArticle(db, user.ID, "文章1")
	test.CreateTestArticle(db, user.ID, "文章2")
	test.CreateTestArticle(db, user.ID, "文章3")

	// 查询列表
	conditions := map[string]interface{}{
		"status": "published",
	}
	articles, total, err := repo.List(1, 10, conditions)
	if err != nil {
		t.Fatalf("查询文章列表失败: %v", err)
	}

	if total != 3 {
		t.Errorf("期望总数 3, 得到 %d", total)
	}

	if len(articles) != 3 {
		t.Errorf("期望文章数量 3, 得到 %d", len(articles))
	}
}

func TestArticleRepository_ListWithPagination(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewArticleRepository(db)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	for i := 1; i <= 15; i++ {
		test.CreateTestArticle(db, user.ID, "文章"+string(rune(i+'0')))
	}

	// 第一页
	conditions := map[string]interface{}{
		"status": "published",
	}
	articles, total, err := repo.List(1, 10, conditions)
	if err != nil {
		t.Fatalf("查询文章列表失败: %v", err)
	}

	if total != 15 {
		t.Errorf("期望总数 15, 得到 %d", total)
	}

	if len(articles) != 10 {
		t.Errorf("期望第一页文章数量 10, 得到 %d", len(articles))
	}

	// 第二页
	articles2, _, err := repo.List(2, 10, conditions)
	if err != nil {
		t.Fatalf("查询文章列表失败: %v", err)
	}

	if len(articles2) != 5 {
		t.Errorf("期望第二页文章数量 5, 得到 %d", len(articles2))
	}
}

func TestArticleRepository_ListWithFilter(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewArticleRepository(db)

	// 创建测试数据
	user1 := test.CreateTestUser(db, "user1", "user1@example.com")
	user2 := test.CreateTestUser(db, "user2", "user2@example.com")

	test.CreateTestArticle(db, user1.ID, "用户1的文章")
	test.CreateTestArticle(db, user2.ID, "用户2的文章")

	// 按作者筛选
	conditions := map[string]interface{}{
		"status":    "published",
		"author_id": user1.ID,
	}
	articles, total, err := repo.List(1, 10, conditions)
	if err != nil {
		t.Fatalf("查询文章列表失败: %v", err)
	}

	if total != 1 {
		t.Errorf("期望总数 1, 得到 %d", total)
	}

	if len(articles) != 1 {
		t.Errorf("期望文章数量 1, 得到 %d", len(articles))
	}

	if articles[0].AuthorID != user1.ID {
		t.Errorf("期望作者ID %d, 得到 %d", user1.ID, articles[0].AuthorID)
	}
}

func TestArticleRepository_Update(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewArticleRepository(db)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	article := test.CreateTestArticle(db, user.ID, "测试文章")

	// 更新文章
	article.Title = "更新后的标题"
	article.Content = "更新后的内容"
	err := repo.Update(article)
	if err != nil {
		t.Fatalf("更新文章失败: %v", err)
	}

	// 验证更新
	found, _ := repo.GetByID(article.ID)
	if found.Title != "更新后的标题" {
		t.Errorf("期望标题 更新后的标题, 得到 %s", found.Title)
	}
}

func TestArticleRepository_IncrementViewCount(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewArticleRepository(db)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	article := test.CreateTestArticle(db, user.ID, "测试文章")

	initialCount := article.ViewCount

	// 增加浏览次数
	err := repo.IncrementViewCount(article.ID)
	if err != nil {
		t.Fatalf("增加浏览次数失败: %v", err)
	}

	// 验证
	found, _ := repo.GetByID(article.ID)
	if found.ViewCount != initialCount+1 {
		t.Errorf("期望浏览次数 %d, 得到 %d", initialCount+1, found.ViewCount)
	}
}

func TestArticleRepository_Delete(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewArticleRepository(db)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	article := test.CreateTestArticle(db, user.ID, "测试文章")

	// 删除文章
	err := repo.Delete(article.ID)
	if err != nil {
		t.Fatalf("删除文章失败: %v", err)
	}

	// 验证删除
	_, err = repo.GetByID(article.ID)
	if err == nil {
		t.Error("文章应该已被删除")
	}
}

