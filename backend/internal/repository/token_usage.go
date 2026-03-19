// Package repository 提供数据访问层
package repository

import (
	"gorm.io/gorm"

	"xiaohongshu/internal/model"
)

// TokenUsageRepository Token使用记录仓库
type TokenUsageRepository struct {
	db *gorm.DB
}

// NewTokenUsageRepository 创建Token使用记录仓库
func NewTokenUsageRepository() *TokenUsageRepository {
	return &TokenUsageRepository{
		db: DB,
	}
}

// Create 创建Token使用记录
func (r *TokenUsageRepository) Create(usage *model.TokenUsage) error {
	return r.db.Create(usage).Error
}

// GetByUserID 根据用户ID获取Token使用记录
func (r *TokenUsageRepository) GetByUserID(userID uint, limit int) ([]model.TokenUsage, error) {
	var usages []model.TokenUsage
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&usages).Error
	return usages, err
}

// GetStatsByUserID 获取用户Token使用统计
func (r *TokenUsageRepository) GetStatsByUserID(userID uint) (*model.TokenUsageStats, error) {
	var stats model.TokenUsageStats

	err := r.db.Model(&model.TokenUsage{}).
		Where("user_id = ?", userID).
		Where("response_status = ?", "success").
		Select("COUNT(*) as total_requests, COALESCE(SUM(input_tokens), 0) as total_prompt_tokens, COALESCE(SUM(output_tokens), 0) as total_completion_tokens, COALESCE(SUM(total_tokens), 0) as total_tokens, COALESCE(SUM(cost), 0) as total_cost").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	stats.SuccessRequests = stats.TotalRequests

	// 获取失败请求数
	var failedCount int64
	r.db.Model(&model.TokenUsage{}).
		Where("user_id = ?", userID).
		Where("response_status = ?", "failed").
		Count(&failedCount)
	stats.FailedRequests = failedCount

	if stats.TotalRequests > 0 {
		stats.AverageTokens = float64(stats.TotalTokens) / float64(stats.TotalRequests)
	}

	return &stats, nil
}

// GetDailyStatsByUserID 获取用户每日Token使用统计
func (r *TokenUsageRepository) GetDailyStatsByUserID(userID uint, days int) ([]model.UserTokenUsage, error) {
	var stats []model.UserTokenUsage

	err := r.db.Model(&model.TokenUsage{}).
		Select("DATE(created_at) as date, SUM(total_tokens) as total_tokens, SUM(cost) as total_cost, COUNT(*) as request_count").
		Where("user_id = ?", userID).
		Group("DATE(created_at)").
		Order("date DESC").
		Limit(days).
		Scan(&stats).Error

	return stats, err
}

// GetStatsByModel 获取按模型统计的使用情况
func (r *TokenUsageRepository) GetStatsByModel(userID uint) ([]model.TokenUsageByModel, error) {
	var stats []model.TokenUsageByModel

	err := r.db.Model(&model.TokenUsage{}).
		Select("model, SUM(total_tokens) as total_tokens, SUM(cost) as total_cost, COUNT(*) as request_count").
		Where("user_id = ?", userID).
		Group("model").
		Order("total_tokens DESC").
		Scan(&stats).Error

	return stats, err
}

// GetGlobalStats 获取全局Token使用统计（仅管理员）
func (r *TokenUsageRepository) GetGlobalStats() (*model.TokenUsageStats, error) {
	var stats model.TokenUsageStats

	err := r.db.Model(&model.TokenUsage{}).
		Where("response_status = ?", "success").
		Select("COUNT(*) as total_requests, COALESCE(SUM(input_tokens), 0) as total_prompt_tokens, COALESCE(SUM(output_tokens), 0) as total_completion_tokens, COALESCE(SUM(total_tokens), 0) as total_tokens, COALESCE(SUM(cost), 0) as total_cost").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	stats.SuccessRequests = stats.TotalRequests

	var failedCount int64
	r.db.Model(&model.TokenUsage{}).
		Where("response_status = ?", "failed").
		Count(&failedCount)
	stats.FailedRequests = failedCount

	if stats.TotalRequests > 0 {
		stats.AverageTokens = float64(stats.TotalTokens) / float64(stats.TotalRequests)
	}

	return &stats, nil
}

// GetGlobalDailyStats 获取全局每日Token使用统计
func (r *TokenUsageRepository) GetGlobalDailyStats(days int) ([]model.UserTokenUsage, error) {
	var stats []model.UserTokenUsage

	err := r.db.Model(&model.TokenUsage{}).
		Select("DATE(created_at) as date, SUM(total_tokens) as total_tokens, SUM(cost) as total_cost, COUNT(*) as request_count").
		Group("DATE(created_at)").
		Order("date DESC").
		Limit(days).
		Scan(&stats).Error

	return stats, err
}
