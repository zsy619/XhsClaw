// Package repository 提供数据访问层
package repository

import (
	"fmt"
	"xiaohongshu/internal/model"

	"gorm.io/gorm"
)

// LLMProviderRepository 大模型配置仓储
type LLMProviderRepository struct {
	db *gorm.DB
}

// NewLLMProviderRepository 创建大模型配置仓储实例
func NewLLMProviderRepository() *LLMProviderRepository {
	return &LLMProviderRepository{db: DB}
}

// Create 创建大模型配置
func (r *LLMProviderRepository) Create(provider *model.LLMProvider) error {
	return r.db.Create(provider).Error
}

// FindByID 根据ID查找
func (r *LLMProviderRepository) FindByID(id uint) (*model.LLMProvider, error) {
	var provider model.LLMProvider
	err := r.db.First(&provider, id).Error
	if err != nil {
		return nil, fmt.Errorf("查询大模型配置失败：%w", err)
	}
	return &provider, nil
}

// ListByUserID 获取用户的所有配置
func (r *LLMProviderRepository) ListByUserID(userID uint) ([]model.LLMProvider, error) {
	var providers []model.LLMProvider
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&providers).Error
	if err != nil {
		return nil, fmt.Errorf("查询用户大模型配置失败：%w", err)
	}
	return providers, nil
}

// GetDefaultByUserID 获取用户默认配置
func (r *LLMProviderRepository) GetDefaultByUserID(userID uint) (*model.LLMProvider, error) {
	var provider model.LLMProvider
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&provider).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("未找到默认配置")
		}
		return nil, fmt.Errorf("查询默认配置失败：%w", err)
	}
	return &provider, nil
}

// Update 更新配置
func (r *LLMProviderRepository) Update(provider *model.LLMProvider) error {
	return r.db.Save(provider).Error
}

// Delete 删除配置
func (r *LLMProviderRepository) Delete(id uint) error {
	return r.db.Delete(&model.LLMProvider{}, id).Error
}

// SetDefault 设置为默认配置
func (r *LLMProviderRepository) SetDefault(userID, providerID uint) error {
	tx := r.db.Begin()
	
	// 取消该用户的其他默认配置
	if err := tx.Model(&model.LLMProvider{}).
		Where("user_id = ? AND is_default = ?", userID, true).
		Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("取消默认配置失败：%w", err)
	}
	
	// 设置指定配置为默认
	if err := tx.Model(&model.LLMProvider{}).
		Where("id = ?", providerID).
		Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("设置默认配置失败：%w", err)
	}
	
	return tx.Commit().Error
}

// GetOrCreate 获取或创建默认配置
func (r *LLMProviderRepository) GetOrCreate(userID uint) (*model.LLMProvider, error) {
	provider, err := r.GetDefaultByUserID(userID)
	if err == nil {
		return provider, nil
	}
	
	// 创建新的默认配置
	provider = &model.LLMProvider{
		UserID:    userID,
		Name:      "默认配置",
		Provider:  model.ProviderCustom,
		IsDefault: true,
		IsEnabled: true,
		Timeout:   60,
		RetryCount: 3,
	}
	
	if err := r.db.Create(provider).Error; err != nil {
		return nil, fmt.Errorf("创建默认配置失败：%w", err)
	}
	
	return provider, nil
}
