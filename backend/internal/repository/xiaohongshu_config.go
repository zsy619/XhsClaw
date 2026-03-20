// Package repository 提供数据访问层
package repository

import (
	"fmt"
	"xiaohongshu/internal/model"

	"gorm.io/gorm"
)

// XiaohongshuConfigRepository 小红书配置仓储
type XiaohongshuConfigRepository struct {
	db *gorm.DB
}

// NewXiaohongshuConfigRepository 创建小红书配置仓储实例
func NewXiaohongshuConfigRepository() *XiaohongshuConfigRepository {
	return &XiaohongshuConfigRepository{db: DB}
}

// Create 创建小红书配置
func (r *XiaohongshuConfigRepository) Create(config *model.XiaohongshuConfig) error {
	return r.db.Create(config).Error
}

// FindByID 根据ID查找
func (r *XiaohongshuConfigRepository) FindByID(id uint) (*model.XiaohongshuConfig, error) {
	var config model.XiaohongshuConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		return nil, fmt.Errorf("查询小红书配置失败：%w", err)
	}
	return &config, nil
}

// ListByUserID 获取用户的所有配置
func (r *XiaohongshuConfigRepository) ListByUserID(userID uint) ([]model.XiaohongshuConfig, error) {
	var configs []model.XiaohongshuConfig
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&configs).Error
	if err != nil {
		return nil, fmt.Errorf("查询用户小红书配置失败：%w", err)
	}
	return configs, nil
}

// GetDefaultByUserID 获取用户默认配置
func (r *XiaohongshuConfigRepository) GetDefaultByUserID(userID uint) (*model.XiaohongshuConfig, error) {
	var config model.XiaohongshuConfig
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("未找到默认配置")
		}
		return nil, fmt.Errorf("查询默认配置失败：%w", err)
	}
	return &config, nil
}

// Update 更新配置
func (r *XiaohongshuConfigRepository) Update(config *model.XiaohongshuConfig) error {
	return r.db.Save(config).Error
}

// Delete 删除配置
func (r *XiaohongshuConfigRepository) Delete(id uint) error {
	return r.db.Delete(&model.XiaohongshuConfig{}, id).Error
}

// SetDefault 设置为默认配置
func (r *XiaohongshuConfigRepository) SetDefault(userID, configID uint) error {
	tx := r.db.Begin()
	
	// 取消该用户的其他默认配置
	if err := tx.Model(&model.XiaohongshuConfig{}).
		Where("user_id = ? AND is_default = ?", userID, true).
		Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("取消默认配置失败：%w", err)
	}
	
	// 设置指定配置为默认
	if err := tx.Model(&model.XiaohongshuConfig{}).
		Where("id = ?", configID).
		Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("设置默认配置失败：%w", err)
	}
	
	return tx.Commit().Error
}

// GetOrCreate 获取或创建默认配置
func (r *XiaohongshuConfigRepository) GetOrCreate(userID uint) (*model.XiaohongshuConfig, error) {
	config, err := r.GetDefaultByUserID(userID)
	if err == nil {
		return config, nil
	}
	
	// 创建新的默认配置
	config = &model.XiaohongshuConfig{
		UserID:    userID,
		Name:      "默认账号",
		IsDefault: true,
		IsEnabled: true,
		Status:    model.XHSStatusPending,
	}
	
	if err := r.db.Create(config).Error; err != nil {
		return nil, fmt.Errorf("创建默认配置失败：%w", err)
	}
	
	return config, nil
}
