package main

import (
	"os"
	"dbapp/internal/config"
	"dbapp/internal/handler"
	"dbapp/internal/middleware"
	"dbapp/internal/model"
	"dbapp/internal/repository"
	"dbapp/internal/service"
	"dbapp/pkg/database"
	"dbapp/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("加载配置失败: " + err.Error())
	}

	// 初始化日志
	logger.Init(cfg.Server.Mode)

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	db, err := database.InitPostgreSQL(cfg.Database)
	if err != nil {
		logger.Fatal("初始化数据库失败", zap.String("error", err.Error()))
	}

	// 自动迁移数据库表结构
	// 可以通过环境变量 AUTO_MIGRATE=false 来禁用自动迁移
	autoMigrate := os.Getenv("AUTO_MIGRATE")
	if autoMigrate != "false" {
		logger.Info("开始执行数据库迁移...")
		if err := db.AutoMigrate(
			&model.User{},
			&model.Article{},
			&model.Category{},
			&model.Tag{},
			&model.Comment{},
			&model.Like{},
		); err != nil {
			logger.Fatal("数据库迁移失败", zap.String("error", err.Error()))
		} else {
			logger.Info("数据库迁移成功")
		}
	}

	// 初始化Repository
	userRepo := repository.NewUserRepository(db)

	// 创建默认管理员用户
	_, err = userRepo.GetByUsername("root")
	if err != nil {
		// 管理员用户不存在，创建它
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("qweasdzxc"), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("创建管理员用户失败：密码加密失败", zap.String("error", err.Error()))
		} else {
			adminUser := &model.User{
				Username:     "root",
				Email:        "admin@dbapp.local",
				PasswordHash: string(hashedPassword),
				Nickname:     "管理员",
				Role:         "admin",
				Status:       "active",
			}
			if err := userRepo.Create(adminUser); err != nil {
				logger.Error("创建管理员用户失败", zap.String("error", err.Error()))
			} else {
				logger.Info("管理员用户创建成功", zap.String("username", "root"))
			}
		}
	} else {
		logger.Info("管理员用户已存在", zap.String("username", "root"))
	}
	articleRepo := repository.NewArticleRepository(db)

	// 初始化Service
	userService := service.NewUserService(userRepo)
	articleService := service.NewArticleService(articleRepo, userRepo)

	// 初始化Handler
	authHandler := handler.NewAuthHandler(userService)
	articleHandler := handler.NewArticleHandler(articleService)

	// 初始化路由
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API路由
	api := router.Group("/api/v1")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetMe)
		}

		// 文章路由
		articles := api.Group("/articles")
		{
			articles.GET("", articleHandler.GetArticleList)
			articles.GET("/:id", articleHandler.GetArticleDetail)
			articles.POST("", middleware.AuthMiddleware(), articleHandler.CreateArticle)
			articles.PUT("/:id", middleware.AuthMiddleware(), articleHandler.UpdateArticle)
			articles.DELETE("/:id", middleware.AuthMiddleware(), articleHandler.DeleteArticle)
		}
	}

	// 启动服务
	addr := ":" + cfg.Server.Port
	logger.Info("服务器启动", zap.String("address", addr))
	if err := router.Run(addr); err != nil {
		logger.Fatal("启动服务器失败", zap.String("error", err.Error()))
	}
}

