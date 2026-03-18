// Package service 提供业务逻辑层
package service

import (
	"strings"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/repository"
)

// TokenUsageService Token使用记录服务
type TokenUsageService struct {
	repo *repository.TokenUsageRepository
}

// NewTokenUsageService 创建Token使用记录服务
func NewTokenUsageService() *TokenUsageService {
	return &TokenUsageService{
		repo: repository.NewTokenUsageRepository(),
	}
}

// TokenPricePerMillion 每百万tokens的价格（美元）
var TokenPricePerMillion = map[string]map[string]float64{
	"deepseek": {
		"prompt":      0.27,      // 输入: $0.27/百万tokens
		"completion":  1.10,     // 输出: $1.10/百万tokens
	},
	"openai": {
		"prompt":      2.50,      // gpt-4 输入
		"completion":  10.00,     // gpt-4 输出
	},
	"anthropic": {
		"prompt":      3.00,      // claude 输入
		"completion":  15.00,     // claude 输出
	},
	"default": {
		"prompt":      1.00,      // 默认输入价格
		"completion":  3.00,      // 默认输出价格
	},
}

// RecordTokenUsage 记录Token使用情况
func (s *TokenUsageService) RecordTokenUsage(
	userID uint,
	modelName, provider, requestType, requestContent, responseStatus, errorMessage, ipAddress, userAgent string,
	promptTokens, completionTokens int,
) error {
	// 计算费用
	cost := s.calculateCost(provider, promptTokens, completionTokens)
	
	// 截断请求内容（保留前500字符）
	if len(requestContent) > 500 {
		requestContent = requestContent[:500]
	}
	
	usage := &model.TokenUsage{
		UserID:           userID,
		Model:            modelName,
		Provider:         provider,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
		Cost:             cost,
		RequestType:      requestType,
		RequestContent:   requestContent,
		ResponseStatus:   responseStatus,
		ErrorMessage:    errorMessage,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
	}
	
	return s.repo.Create(usage)
}

// calculateCost 计算费用
func (s *TokenUsageService) calculateCost(provider string, promptTokens, completionTokens int) float64 {
	provider = strings.ToLower(provider)
	prices, ok := TokenPricePerMillion[provider]
	if !ok {
		prices = TokenPricePerMillion["default"]
	}
	
	promptCost := float64(promptTokens) / 1_000_000 * prices["prompt"]
	completionCost := float64(completionTokens) / 1_000_000 * prices["completion"]
	
	return promptCost + completionCost
}

// GetUserTokenUsage 获取用户Token使用记录
func (s *TokenUsageService) GetUserTokenUsage(userID uint, limit int) ([]model.TokenUsage, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetByUserID(userID, limit)
}

// GetUserTokenStats 获取用户Token使用统计
func (s *TokenUsageService) GetUserTokenStats(userID uint) (*model.TokenUsageStats, error) {
	return s.repo.GetStatsByUserID(userID)
}

// GetUserDailyStats 获取用户每日Token使用统计
func (s *TokenUsageService) GetUserDailyStats(userID uint, days int) ([]model.UserTokenUsage, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetDailyStatsByUserID(userID, days)
}

// GetUserStatsByModel 获取用户按模型统计的使用情况
func (s *TokenUsageService) GetUserStatsByModel(userID uint) ([]model.TokenUsageByModel, error) {
	return s.repo.GetStatsByModel(userID)
}

// GetGlobalTokenStats 获取全局Token使用统计（仅管理员）
func (s *TokenUsageService) GetGlobalTokenStats() (*model.TokenUsageStats, error) {
	return s.repo.GetGlobalStats()
}

// GetGlobalDailyStats 获取全局每日Token使用统计
func (s *TokenUsageService) GetGlobalDailyStats(days int) ([]model.UserTokenUsage, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetGlobalDailyStats(days)
}
