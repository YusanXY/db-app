package service

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/repository"
	"dbapp/internal/test"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := NewUserService(userRepo)

	req := &request.RegisterRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		Nickname: "新用户",
	}

	user, err := userService.Register(req)
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}

	if user.Username != "newuser" {
		t.Errorf("期望用户名 newuser, 得到 %s", user.Username)
	}

	if user.Email != "newuser@example.com" {
		t.Errorf("期望邮箱 newuser@example.com, 得到 %s", user.Email)
	}
}

func TestUserService_Register_DuplicateUsername(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := NewUserService(userRepo)

	// 创建已存在的用户
	test.CreateTestUser(db, "existinguser", "existing@example.com")

	req := &request.RegisterRequest{
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "password123",
	}

	_, err := userService.Register(req)
	if err == nil {
		t.Error("应该返回用户名已存在的错误")
	}
}

func TestUserService_Register_DuplicateEmail(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := NewUserService(userRepo)

	// 创建已存在的用户
	test.CreateTestUser(db, "user1", "existing@example.com")

	req := &request.RegisterRequest{
		Username: "newuser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	_, err := userService.Register(req)
	if err == nil {
		t.Error("应该返回邮箱已存在的错误")
	}
}

func TestUserService_Login(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := NewUserService(userRepo)

	// 创建测试用户（密码需要是bcrypt哈希）
	req := &request.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err := userService.Register(req)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 登录
	loginReq := &request.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	result, err := userService.Login(loginReq)
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	if result.Token == "" {
		t.Error("Token应该被生成")
	}

	if result.User.Username != "testuser" {
		t.Errorf("期望用户名 testuser, 得到 %s", result.User.Username)
	}
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := NewUserService(userRepo)

	// 创建测试用户
	req := &request.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err := userService.Register(req)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 错误的密码
	loginReq := &request.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	_, err = userService.Login(loginReq)
	if err == nil {
		t.Error("应该返回登录失败的错误")
	}
}

func TestUserService_GetByID(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := NewUserService(userRepo)

	// 创建测试用户
	user := test.CreateTestUser(db, "testuser", "test@example.com")

	// 查询用户
	found, err := userService.GetByID(user.ID)
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}

	if found.Username != "testuser" {
		t.Errorf("期望用户名 testuser, 得到 %s", found.Username)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := NewUserService(userRepo)

	// 查询不存在的用户
	_, err := userService.GetByID(99999)
	if err == nil {
		t.Error("应该返回用户不存在的错误")
	}
}

