package handler

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/errors"
	"dbapp/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "注册信息"
// @Success 201 {object} response.UserResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误: "+err.Error()))
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(201, gin.H{
		"code":    201,
		"message": "注册成功",
		"data":    user,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "登录信息"
// @Success 200 {object} response.LoginResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("参数错误"))
		return
	}

	result, err := h.userService.Login(&req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": result,
	})
}

// GetMe 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.UserResponse
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint64)

	user, err := h.userService.GetByID(userIDUint)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": user,
	})
}

