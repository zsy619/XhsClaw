// Package service 提供业务逻辑层
package service

import (
	"fmt"

	"xiaohongshu/internal/model"
	"xiaohongshu/internal/repository"
	"xiaohongshu/pkg/errno"
)

// UserConfigService 用户配置服务（整合 LLM、小红书和发布配置）
type UserConfigService struct {
	llmRepo     *repository.LLMProviderRepository
	xhsRepo     *repository.XHSConfigRepository
	publishRepo *repository.XHSConfigRepository
}

// NewUserConfigService 创建用户配置服务实例
func NewUserConfigService() *UserConfigService {
	return &UserConfigService{
		llmRepo:     repository.NewLLMProviderRepository(),
		xhsRepo:     repository.NewXHSConfigRepository(),
		publishRepo: repository.NewXHSConfigRepository(),
	}
}

// UserConfigResponse 整合后的用户配置响应
type UserConfigResponse struct {
	UserID               uint   `json:"user_id"`
	DefaultLLMProviderID uint   `json:"default_llm_provider_id"`
	DefaultXHSConfigID   uint   `json:"default_xhs_config_id"`
	DefaultPublishTime   string `json:"default_publish_time"`
	AutoPublishEnabled   bool   `json:"auto_publish_enabled"`
}

// GetUserConfig 获取整合后的用户配置
func (s *UserConfigService) GetUserConfig(userID uint) (*UserConfigResponse, error) {
	// 获取默认的 LLM Provider
	var defaultLLMID uint
	llmProvider, err := s.llmRepo.GetDefaultByUserID(userID)
	if err == nil && llmProvider != nil {
		defaultLLMID = llmProvider.ID
	}

	// 获取默认的小红书配置
	var defaultXHSID uint
	xhsConfig, err := s.xhsRepo.GetActive(userID)
	if err == nil && xhsConfig != nil {
		defaultXHSID = xhsConfig.ID
	}

	// 获取发布配置（暂时忽略）
	_, _ = s.publishRepo.GetActive(userID)

	return &UserConfigResponse{
		UserID:               userID,
		DefaultLLMProviderID: defaultLLMID,
		DefaultXHSConfigID:   defaultXHSID,
		DefaultPublishTime:   "",
		AutoPublishEnabled:   false,
	}, nil
}

// UpdateUserConfig 更新用户配置
func (s *UserConfigService) UpdateUserConfig(userID uint, req *model.UserConfigRequest) (*UserConfigResponse, error) {
	// 更新 LLM 配置
	if req.LLMAPIKey != "" || req.LLMBaseURL != "" || req.LLMModel != "" {
		if err := s.updateLLMConfig(userID, req); err != nil {
			return nil, errno.InternalError
		}
	}

	// 更新小红书配置
	if req.XiaohongshuCookie != "" || req.XiaohongshuUserId != "" || req.XiaohongshuToken != "" {
		if err := s.updateXHSConfig(userID, req); err != nil {
			return nil, errno.InternalError
		}
	}

	// 更新发布配置
	if err := s.updatePublishConfig(userID, req); err != nil {
		return nil, errno.InternalError
	}

	// 返回更新后的配置
	return s.GetUserConfig(userID)
}

// updateLLMConfig 更新 LLM 配置
func (s *UserConfigService) updateLLMConfig(userID uint, req *model.UserConfigRequest) error {
	// 尝试获取现有的默认 LLM Provider
	existingProvider, err := s.llmRepo.GetDefaultByUserID(userID)
	if err != nil {
		// 如果不存在，创建一个新的
		provider := &model.LLMProvider{
			UserID:     userID,
			Name:       "默认配置",
			Provider:   model.ProviderCustom,
			APIKey:     req.LLMAPIKey,
			BaseURL:    req.LLMBaseURL,
			ModelName:  req.LLMModel,
			IsDefault:  true,
			IsEnabled:  true,
			Timeout:    60,
			RetryCount: 3,
		}
		return s.llmRepo.Create(provider)
	}

	// 更新现有配置
	if req.LLMAPIKey != "" {
		existingProvider.APIKey = req.LLMAPIKey
	}
	if req.LLMBaseURL != "" {
		existingProvider.BaseURL = req.LLMBaseURL
	}
	if req.LLMModel != "" {
		existingProvider.ModelName = req.LLMModel
	}

	return s.llmRepo.Update(existingProvider)
}

// updateXHSConfig 更新小红书配置
func (s *UserConfigService) updateXHSConfig(userID uint, req *model.UserConfigRequest) error {
	// 尝试获取现有的默认小红书配置
	existingConfig, err := s.xhsRepo.GetActive(userID)
	if err != nil {
		// 如果不存在，创建一个新的
		config := &model.XHSConfig{
			UserID:      userID,
			Name:        "默认账号",
			Cookie:      req.XiaohongshuCookie,
			XHSUserID:   req.XiaohongshuUserId,
			Token:       req.XiaohongshuToken,
			IsDefault:   true,
			IsEnabled:   true,
			Status:      "pending",
			Description: "通过用户配置更新",
		}
		return s.xhsRepo.Create(config)
	}

	// 更新现有配置
	if req.XiaohongshuCookie != "" {
		existingConfig.Cookie = req.XiaohongshuCookie
	}
	if req.XiaohongshuUserId != "" {
		existingConfig.XHSUserID = req.XiaohongshuUserId
	}
	if req.XiaohongshuToken != "" {
		existingConfig.Token = req.XiaohongshuToken
	}

	return s.xhsRepo.Update(existingConfig)
}

// updatePublishConfig 更新发布配置（暂时禁用）
func (s *UserConfigService) updatePublishConfig(userID uint, req *model.UserConfigRequest) error {
	// 暂时忽略发布配置更新
	return nil
}

// GetLLMConfig 获取大模型配置（优先使用用户配置，否则使用系统默认配置）
func (s *UserConfigService) GetLLMConfig(userID uint) (apiKey, baseURL, model string) {
	provider, err := s.llmRepo.GetDefaultByUserID(userID)
	if err == nil && provider != nil {
		if provider.APIKey != "" {
			apiKey = provider.APIKey
		}
		if provider.BaseURL != "" {
			baseURL = provider.BaseURL
		}
		if provider.ModelName != "" {
			model = provider.ModelName
		}
	} else {
		fmt.Println("查询用户 LLM 配置失败:", err)
	}
	return apiKey, baseURL, model
}
