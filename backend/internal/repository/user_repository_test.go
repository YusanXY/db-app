package repository

import (
	"dbapp/internal/model"
	"dbapp/internal/test"
	"testing"
	"time"
)

func TestUserRepository_Create(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewUserRepository(db)

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$testhash",
		Nickname:     "测试用户",
		Role:         "user",
		Status:       "active",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}

	if user.ID == 0 {
		t.Error("用户ID应该被设置")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewUserRepository(db)

	// 创建测试用户
	user := test.CreateTestUser(db, "testuser", "test@example.com")

	// 查询用户
	found, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}

	if found.Username != "testuser" {
		t.Errorf("期望用户名 testuser, 得到 %s", found.Username)
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewUserRepository(db)

	// 创建测试用户
	test.CreateTestUser(db, "testuser", "test@example.com")

	// 查询用户
	found, err := repo.GetByUsername("testuser")
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}

	if found.Username != "testuser" {
		t.Errorf("期望用户名 testuser, 得到 %s", found.Username)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewUserRepository(db)

	// 创建测试用户
	test.CreateTestUser(db, "testuser", "test@example.com")

	// 查询用户
	found, err := repo.GetByEmail("test@example.com")
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}

	if found.Email != "test@example.com" {
		t.Errorf("期望邮箱 test@example.com, 得到 %s", found.Email)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewUserRepository(db)

	// 创建测试用户
	user := test.CreateTestUser(db, "testuser", "test@example.com")

	// 更新用户
	user.Nickname = "新昵称"
	err := repo.Update(user)
	if err != nil {
		t.Fatalf("更新用户失败: %v", err)
	}

	// 验证更新
	found, _ := repo.GetByID(user.ID)
	if found.Nickname != "新昵称" {
		t.Errorf("期望昵称 新昵称, 得到 %s", found.Nickname)
	}
}

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	repo := NewUserRepository(db)

	// 创建测试用户
	user := test.CreateTestUser(db, "testuser", "test@example.com")

	// 更新最后登录时间
	err := repo.UpdateLastLogin(user.ID)
	if err != nil {
		t.Fatalf("更新最后登录时间失败: %v", err)
	}

	// 验证更新
	found, _ := repo.GetByID(user.ID)
	if found.LastLoginAt == nil {
		t.Error("最后登录时间应该被设置")
	}

	// 验证时间在合理范围内（最近1分钟内）
	if time.Since(*found.LastLoginAt) > time.Minute {
		t.Error("最后登录时间应该在最近1分钟内")
	}
}

