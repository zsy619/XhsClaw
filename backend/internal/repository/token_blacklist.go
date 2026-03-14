// Package repository 提供数据访问层
package repository

import (
	"time"
	"xiaohongshu/internal/model"
)

// TokenBlacklistRepository Token黑名单仓库
type TokenBlacklistRepository struct{}

// NewTokenBlacklistRepository 创建Token黑名单仓库实例
func NewTokenBlacklistRepository() *TokenBlacklistRepository {
	return &TokenBlacklistRepository{}
}

// AddToBlacklist 将Token添加到黑名单
func (r *TokenBlacklistRepository) AddToBlacklist(token string, userID uint, expiresAt time.Time) error {
	blacklist := &model.TokenBlacklist{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return DB.Create(blacklist).Error
}

// IsTokenBlacklisted 检查Token是否在黑名单中
func (r *TokenBlacklistRepository) IsTokenBlacklisted(token string) (bool, error) {
	var count int64
	err := DB.Model(&model.TokenBlacklist{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CleanExpiredTokens 清理已过期的Token
func (r *TokenBlacklistRepository) CleanExpiredTokens() error {
	return DB.Where("expires_at < ?", time.Now()).Delete(&model.TokenBlacklist{}).Error
}
