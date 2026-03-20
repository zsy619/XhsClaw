// Package service 提供业务逻辑层
package service

import (
	"fmt"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/repository"
	"xiaohongshu/pkg/errno"
)

// UserConfigService 用户配置服务
type UserConfigService struct {
	userConfigRepo *repository.UserConfigRepository
}

// NewUserConfigService 创建用户配置服务实例
func NewUserConfigService() *UserConfigService {
	var repo *repository.UserConfigRepository
	if repository.DB != nil {
		repo = repository.NewUserConfigRepository()
	}
	return &UserConfigService{
		userConfigRepo: repo,
	}
}

// GetUserConfig 获取用户配置
func (s *UserConfigService) GetUserConfig(userID uint) (*model.UserConfig, error) {
	if s.userConfigRepo == nil {
		return &model.UserConfig{UserID: userID}, nil
	}
	config, err := s.userConfigRepo.GetOrCreate(userID)
	if err != nil {
		return nil, errno.InternalError
	}
	return config, nil
}

// UpdateUserConfig 更新用户配置
func (s *UserConfigService) UpdateUserConfig(userID uint, req *model.UserConfigRequest) (*model.UserConfig, error) {
	if s.userConfigRepo == nil {
		return &model.UserConfig{UserID: userID}, nil
	}
	config, err := s.userConfigRepo.GetOrCreate(userID)
	if err != nil {
		return nil, errno.InternalError
	}

	if req.LLMAPIKey != "" {
		config.LLMAPIKey = req.LLMAPIKey
	}
	if req.LLMBaseURL != "" {
		config.LLMBaseURL = req.LLMBaseURL
	}
	if req.LLMModel != "" {
		config.LLMModel = req.LLMModel
	}
	if req.XiaohongshuCookie != "" {
		config.XiaohongshuCookie = req.XiaohongshuCookie
	}
	if req.XiaohongshuUserId != "" {
		config.XiaohongshuUserId = req.XiaohongshuUserId
	}
	if req.XiaohongshuToken != "" {
		config.XiaohongshuToken = req.XiaohongshuToken
	}
	if req.DefaultPublishTime != "" {
		config.DefaultPublishTime = req.DefaultPublishTime
	}

	config.AutoPublishEnabled = req.AutoPublishEnabled

	if err := s.userConfigRepo.Update(config); err != nil {
		return nil, errno.InternalError
	}

	return config, nil
}

// GetLLMConfig 获取大模型配置（优先使用用户配置，否则使用系统默认配置）
func (s *UserConfigService) GetLLMConfig(userID uint) (apiKey, baseURL, model string) {
	if s.userConfigRepo == nil {
		return "", "", ""
	}
	config, err := s.userConfigRepo.FindByUserID(userID)
	if err == nil {
		if config.LLMAPIKey != "" {
			apiKey = config.LLMAPIKey
		}
		if config.LLMBaseURL != "" {
			baseURL = config.LLMBaseURL
		}
		if config.LLMModel != "" {
			model = config.LLMModel
		}
	} else {
		fmt.Println("查询用户配置失败:", err)
	}
	return apiKey, baseURL, model
}
