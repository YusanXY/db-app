package handler

import (
	"bytes"
	"dbapp/internal/config"
	"dbapp/internal/dto/request"
	"dbapp/internal/middleware"
	"dbapp/internal/repository"
	"dbapp/internal/service"
	"dbapp/internal/test"
	"dbapp/pkg/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	// 设置测试用的 JWT 配置
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret-key-for-testing",
			ExpiresIn: 3600,
		},
	}
}

func TestArticleHandler_GetArticleList(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := service.NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)
	articleHandler := NewArticleHandler(articleService)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	test.CreateTestArticle(db, user.ID, "文章1")
	test.CreateTestArticle(db, user.ID, "文章2")

	router := setupRouter()
	router.GET("/api/v1/articles", articleHandler.GetArticleList)

	req, _ := http.NewRequest("GET", "/api/v1/articles?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, int(response["code"].(float64)))
}

func TestArticleHandler_GetArticleDetail(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := service.NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)
	articleHandler := NewArticleHandler(articleService)

	// 创建测试数据
	user := test.CreateTestUser(db, "testuser", "test@example.com")
	_ = test.CreateTestArticle(db, user.ID, "测试文章")

	router := setupRouter()
	router.GET("/api/v1/articles/:id", articleHandler.GetArticleDetail)

	req, _ := http.NewRequest("GET", "/api/v1/articles/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestArticleHandler_CreateArticle(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := service.NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)
	articleHandler := NewArticleHandler(articleService)

	// 创建测试用户
	user := test.CreateTestUser(db, "testuser", "test@example.com")

	// 生成token
	token, _ := utils.GenerateJWT(user.ID, user.Username, user.Role)

	router := setupRouter()
	router.POST("/api/v1/articles", middleware.AuthMiddleware(), articleHandler.CreateArticle)

	reqBody := request.CreateArticleRequest{
		Title:   "新文章",
		Content: "文章内容",
		Summary: "文章摘要",
		Status:  "published",
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/articles", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestArticleHandler_CreateArticle_Unauthorized(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(db)

	articleRepo := repository.NewArticleRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	articleImageRepo := repository.NewArticleImageRepository(db)
	articleService := service.NewArticleService(articleRepo, userRepo, likeRepo, articleImageRepo)
	articleHandler := NewArticleHandler(articleService)

	router := setupRouter()
	router.POST("/api/v1/articles", middleware.AuthMiddleware(), articleHandler.CreateArticle)

	reqBody := request.CreateArticleRequest{
		Title:   "新文章",
		Content: "文章内容",
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/articles", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	// 不设置Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

