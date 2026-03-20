// Package repository 提供数据访问层
package repository

import (
	"fmt"
	"xiaohongshu/internal/model"

	"gorm.io/gorm"
)

// PublishConfigRepository 发布配置仓储
type PublishConfigRepository struct {
	db *gorm.DB
}

// NewPublishConfigRepository 创建发布配置仓储实例
func NewPublishConfigRepository() *PublishConfigRepository {
	return &PublishConfigRepository{db: DB}
}

// Create 创建发布配置
func (r *PublishConfigRepository) Create(config *model.PublishConfig) error {
	return r.db.Create(config).Error
}

// FindByID 根据ID查找
func (r *PublishConfigRepository) FindByID(id uint) (*model.PublishConfig, error) {
	var config model.PublishConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		return nil, fmt.Errorf("查询发布配置失败：%w", err)
	}
	return &config, nil
}

// GetByUserID 根据用户ID获取配置
func (r *PublishConfigRepository) GetByUserID(userID uint) (*model.PublishConfig, error) {
	var config model.PublishConfig
	err := r.db.Where("user_id = ?", userID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回空的配置对象
			return &model.PublishConfig{
				UserID: userID,
			}, nil
		}
		return nil, fmt.Errorf("查询发布配置失败：%w", err)
	}
	return &config, nil
}

// Update 更新配置
func (r *PublishConfigRepository) Update(config *model.PublishConfig) error {
	return r.db.Save(config).Error
}

// GetOrCreate 获取或创建发布配置
func (r *PublishConfigRepository) GetOrCreate(userID uint) (*model.PublishConfig, error) {
	config, err := r.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	
	// 如果不存在，则创建
	if config.ID == 0 {
		config = &model.PublishConfig{
			UserID: userID,
		}
		if err := r.db.Create(config).Error; err != nil {
			return nil, fmt.Errorf("创建发布配置失败：%w", err)
		}
	}
	
	return config, nil
}
