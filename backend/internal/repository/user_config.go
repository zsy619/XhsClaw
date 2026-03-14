// Package repository 提供数据访问层
package repository

import (
	"xiaohongshu/internal/model"

	"gorm.io/gorm"
)

// UserConfigRepository 用户配置仓储
type UserConfigRepository struct {
	db *gorm.DB
}

// NewUserConfigRepository 创建用户配置仓储实例
func NewUserConfigRepository() *UserConfigRepository {
	return &UserConfigRepository{
		db: DB,
	}
}

// FindByUserID 根据用户ID查找配置
func (r *UserConfigRepository) FindByUserID(userID uint) (*model.UserConfig, error) {
	var config model.UserConfig
	err := r.db.Where("user_id = ?", userID).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Create 创建用户配置
func (r *UserConfigRepository) Create(config *model.UserConfig) error {
	return r.db.Create(config).Error
}

// Update 更新用户配置
func (r *UserConfigRepository) Update(config *model.UserConfig) error {
	return r.db.Save(config).Error
}

// Upsert 更新或创建用户配置
func (r *UserConfigRepository) Upsert(config *model.UserConfig) error {
	var existing model.UserConfig
	err := r.db.Where("user_id = ?", config.UserID).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.db.Create(config).Error
		}
		return err
	}
	
	config.ID = existing.ID
	config.CreatedAt = existing.CreatedAt
	return r.db.Save(config).Error
}

// GetOrCreate 获取或创建用户配置
func (r *UserConfigRepository) GetOrCreate(userID uint) (*model.UserConfig, error) {
	var config model.UserConfig
	err := r.db.Where("user_id = ?", userID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			config = model.UserConfig{
				UserID: userID,
			}
			if err := r.db.Create(&config).Error; err != nil {
				return nil, err
			}
			return &config, nil
		}
		return nil, err
	}
	return &config, nil
}
