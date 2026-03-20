// Package service 提供业务逻辑层
package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"xiaohongshu/internal/model"
	"xiaohongshu/internal/repository"
	"xiaohongshu/internal/utils"
	"xiaohongshu/pkg/errno"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.LoginResponse, error) {
	// 检查用户名是否已存在
	_, err := s.userRepo.FindByUsername(req.Username)
	if err == nil {
		log.Printf("[UserService] Register failed: username %s already exists", req.Username)
		return nil, errno.UserAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[UserService] Register error: failed to check username: %v", err)
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}

	// 检查邮箱是否已存在
	_, err = s.userRepo.FindByEmail(req.Email)
	if err == nil {
		log.Printf("[UserService] Register failed: email %s already exists", req.Email)
		return nil, errno.UserAlreadyExists.WithMessage("邮箱已被注册")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[UserService] Register error: failed to check email: %v", err)
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("[UserService] Register error: failed to hash password: %v", err)
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		RoleID:   3, // 默认普通用户角色ID为3
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		log.Printf("[UserService] Register error: failed to create user: %v", err)
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	log.Printf("[UserService] Register success: user %d created", user.ID)

	// 生成token（默认用户角色代码为"user"）
	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		log.Printf("[UserService] Register error: failed to generate token for user %d: %v", user.ID, err)
		return nil, fmt.Errorf("生成认证令牌失败: %w", err)
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[UserService] Login failed: user %s not found", req.Username)
			return nil, errno.UserNotFound
		}
		log.Printf("[UserService] Login error: failed to find user: %v", err)
		return nil, fmt.Errorf("查找用户失败: %w", err)
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		log.Printf("[UserService] Login failed: wrong password for user %s", req.Username)
		return nil, errno.WrongPassword
	}

	// 检查用户状态
	if user.Status != 1 {
		log.Printf("[UserService] Login failed: user %s is disabled", req.Username)
		return nil, errno.UserDisabled
	}

	// 确定角色代码
	roleCode := "user"
	if user.Role != nil {
		roleCode = user.Role.Code
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, roleCode)
	if err != nil {
		log.Printf("[UserService] Login error: failed to generate token for user %d: %v", user.ID, err)
		return nil, fmt.Errorf("生成认证令牌失败: %w", err)
	}

	log.Printf("[UserService] Login success: user %d logged in", user.ID)
	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[UserService] GetUserInfo failed: user %d not found", userID)
			return nil, errno.UserNotFound
		}
		log.Printf("[UserService] GetUserInfo error: failed to find user %d: %v", userID, err)
		return nil, fmt.Errorf("查找用户失败: %w", err)
	}
	return user, nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}
