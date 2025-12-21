package handler

import (
	"dbapp/internal/config"
	"dbapp/pkg/response"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FileHandler struct {
	cfg *config.Config
}

func NewFileHandler(cfg *config.Config) *FileHandler {
	return &FileHandler{
		cfg: cfg,
	}
}

// UploadFile 上传文件
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, "请选择要上传的文件")
		return
	}

	// 检查文件大小
	if file.Size > h.cfg.File.MaxSize {
		response.Error(c, 400, fmt.Sprintf("文件大小不能超过 %d MB", h.cfg.File.MaxSize/1024/1024))
		return
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	ext = strings.TrimPrefix(ext, ".")
	allowed := false
	for _, allowedExt := range h.cfg.File.AllowedExt {
		if ext == strings.ToLower(allowedExt) {
			allowed = true
			break
		}
	}
	if !allowed {
		response.Error(c, 400, fmt.Sprintf("不支持的文件类型，允许的类型: %v", h.cfg.File.AllowedExt))
		return
	}

	// 创建上传目录
	uploadPath := h.cfg.File.UploadPath
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		zap.L().Error("创建上传目录失败", zap.String("error", err.Error()))
		response.Error(c, 500, "创建上传目录失败")
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(uploadPath, filename)

	// 保存文件
	src, err := file.Open()
	if err != nil {
		zap.L().Error("打开上传文件失败", zap.String("error", err.Error()))
		response.Error(c, 500, "打开上传文件失败")
		return
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		zap.L().Error("创建文件失败", zap.String("error", err.Error()))
		response.Error(c, 500, "创建文件失败")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		zap.L().Error("保存文件失败", zap.String("error", err.Error()))
		response.Error(c, 500, "保存文件失败")
		return
	}

	// 返回文件URL
	fileURL := fmt.Sprintf("/uploads/%s", filename)
	response.Success(c, gin.H{
		"url":  fileURL,
		"name": file.Filename,
		"size": file.Size,
	})
}

