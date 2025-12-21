package test

import (
	"dbapp/internal/model"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

// SetupTestDB 设置测试数据库
func SetupTestDB(t *testing.T) *gorm.DB {
	// 使用SQLite内存数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("无法连接测试数据库: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.Category{},
		&model.Tag{},
		&model.Comment{},
		&model.Like{},
	)
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	TestDB = db
	return db
}

// TeardownTestDB 清理测试数据库
func TeardownTestDB(db *gorm.DB) {
	if db != nil {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}

// InitTestLogger 初始化测试日志
func InitTestLogger() {
	logger.Init("debug")
}

// CreateTestUser 创建测试用户
func CreateTestUser(db *gorm.DB, username, email string) *model.User {
	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: "$2a$10$testhash",
		Nickname:     "测试用户",
		Role:         "user",
		Status:       "active",
	}
	db.Create(user)
	return user
}

// CreateTestArticle 创建测试文章
func CreateTestArticle(db *gorm.DB, authorID uint64, title string) *model.Article {
	article := &model.Article{
		Title:    title,
		Slug:     "test-article",
		Content:  "测试内容",
		Summary:  "测试摘要",
		AuthorID: authorID,
		Status:   "published",
	}
	db.Create(article)
	return article
}

// IsCI 检查是否在CI环境中
func IsCI() bool {
	return os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true"
}

