package handler

import (
	"bytes"
	"dbapp/internal/dto/request"
	"dbapp/internal/repository"
	"dbapp/internal/service"
	"dbapp/internal/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandler_Register(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authHandler := NewAuthHandler(userService)

	router := setupRouter()
	router.POST("/api/v1/auth/register", authHandler.Register)

	reqBody := request.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "测试用户",
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 201, int(response["code"].(float64)))
}

func TestAuthHandler_Register_InvalidData(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authHandler := NewAuthHandler(userService)

	router := setupRouter()
	router.POST("/api/v1/auth/register", authHandler.Register)

	reqBody := map[string]interface{}{
		"username": "", // 无效的用户名
		"email":    "invalid-email",
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestAuthHandler_Login(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authHandler := NewAuthHandler(userService)

	// 先注册用户
	registerReq := &request.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err := userService.Register(registerReq)
	if err != nil {
		t.Fatalf("注册用户失败: %v", err)
	}

	router := setupRouter()
	router.POST("/api/v1/auth/login", authHandler.Login)

	reqBody := request.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, int(response["code"].(float64)))
	assert.NotNil(t, response["data"])
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authHandler := NewAuthHandler(userService)

	router := setupRouter()
	router.POST("/api/v1/auth/login", authHandler.Login)

	reqBody := request.LoginRequest{
		Username: "nonexistent",
		Password: "wrongpassword",
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

