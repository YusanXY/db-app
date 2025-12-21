package service

import (
	"dbapp/internal/dto/request"
	"dbapp/internal/dto/response"
	"dbapp/internal/errors"
	"dbapp/internal/model"
	"dbapp/internal/repository"
	"dbapp/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Register(req *request.RegisterRequest) (*response.UserResponse, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.NewBadRequestError("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, _ = s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.NewBadRequestError("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalError("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Nickname:     req.Nickname,
		Role:         "user",
		Status:       "active",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.NewInternalError("创建用户失败")
	}

	return s.toResponse(user), nil
}

func (s *UserService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.NewUnauthorizedError("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.NewUnauthorizedError("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.NewForbiddenError("用户已被禁用")
	}

	// 生成JWT Token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.NewInternalError("生成Token失败")
	}

	// 更新最后登录时间
	go s.userRepo.UpdateLastLogin(user.ID)

	return &response.LoginResponse{
		Token:     token,
		ExpiresIn: 3600,
		User:      s.toResponse(user),
	}, nil
}

func (s *UserService) GetByID(id uint64) (*response.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("用户不存在")
	}
	return s.toResponse(user), nil
}

func (s *UserService) toResponse(user *model.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}

