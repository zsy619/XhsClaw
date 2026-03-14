// Package service 提供业务逻辑层
package service

import (
	"errors"
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
func (s *UserService) Register(req *model.RegisterRequest) (*model.LoginResponse, error) {
	// 检查用户名是否已存在
	_, err := s.userRepo.FindByUsername(req.Username)
	if err == nil {
		return nil, errno.UserAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 检查邮箱是否已存在
	_, err = s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errno.UserAlreadyExists.WithMessage("邮箱已被注册")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errno.InternalError
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Role:     "user",
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errno.InternalError
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errno.InternalError
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.UserNotFound
		}
		return nil, errno.InternalError
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errno.WrongPassword
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errno.UserDisabled
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errno.InternalError
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.UserNotFound
		}
		return nil, errno.InternalError
	}
	return user, nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}
