// Package repository 提供数据访问层
package repository

import (
	"fmt"

	"gorm.io/gorm"

	"xiaohongshu/internal/model"
)

// XHSConfigRepository 小红书配置仓储
type XHSConfigRepository struct {
	db *gorm.DB
}

// NewXHSConfigRepository 创建小红书配置仓储实例
func NewXHSConfigRepository() *XHSConfigRepository {
	return &XHSConfigRepository{db: DB}
}

// List 获取用户的小红书配置列表（支持分页）
func (r *XHSConfigRepository) List(userID uint, page, pageSize int) ([]model.XHSConfig, int64, error) {
	var configs []model.XHSConfig
	var total int64

	query := r.db.Model(&model.XHSConfig{}).Where("user_id = ?", userID)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计小红书配置数量失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&configs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询小红书配置列表失败: %w", err)
	}

	return configs, total, nil
}

// FindByID 根据ID查找
func (r *XHSConfigRepository) FindByID(id uint) (*model.XHSConfig, error) {
	var config model.XHSConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		return nil, fmt.Errorf("查询小红书配置失败: %w", err)
	}
	return &config, nil
}

// FindByUserIDAndID 根据用户ID和配置ID查找
func (r *XHSConfigRepository) FindByUserIDAndID(userID, id uint) (*model.XHSConfig, error) {
	var config model.XHSConfig
	err := r.db.Where("user_id = ? AND id = ?", userID, id).First(&config).Error
	if err != nil {
		return nil, fmt.Errorf("查询小红书配置失败: %w", err)
	}
	return &config, nil
}

// GetActive 获取用户激活的配置
func (r *XHSConfigRepository) GetActive(userID uint) (*model.XHSConfig, error) {
	var config model.XHSConfig
	err := r.db.Where("user_id = ? AND is_enabled = ? AND is_default = ?", userID, true, true).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有默认配置，返回第一个启用的配置
			err = r.db.Where("user_id = ? AND is_enabled = ?", userID, true).First(&config).Error
			if err != nil {
				return nil, fmt.Errorf("没有找到激活的小红书配置: %w", err)
			}
			return &config, nil
		}
		return nil, fmt.Errorf("查询激活配置失败: %w", err)
	}
	return &config, nil
}

// Create 创建配置
func (r *XHSConfigRepository) Create(config *model.XHSConfig) error {
	// 如果设置为默认，取消其他默认配置
	if config.IsDefault {
		if err := r.clearDefaultFlags(config.UserID); err != nil {
			return err
		}
	}

	return r.db.Create(config).Error
}

// Update 更新配置
func (r *XHSConfigRepository) Update(config *model.XHSConfig) error {
	// 如果设置为默认，取消其他默认配置
	if config.IsDefault {
		if err := r.clearDefaultFlags(config.UserID); err != nil {
			return err
		}
	}

	return r.db.Save(config).Error
}

// Delete 删除配置
func (r *XHSConfigRepository) Delete(id uint) error {
	return r.db.Delete(&model.XHSConfig{}, id).Error
}

// clearDefaultFlags 清除用户的所有默认配置标志
func (r *XHSConfigRepository) clearDefaultFlags(userID uint) error {
	return r.db.Model(&model.XHSConfig{}).
		Where("user_id = ? AND is_default = ?", userID, true).
		Update("is_default", false).Error
}
