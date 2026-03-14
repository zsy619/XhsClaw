// Package repository 提供数据访问层
package repository

import (
	"xiaohongshu/internal/model"

	"gorm.io/gorm"
)

// PublishRepository 发布记录仓储
type PublishRepository struct {
	db *gorm.DB
}

// NewPublishRepository 创建发布记录仓储实例
func NewPublishRepository() *PublishRepository {
	return &PublishRepository{db: DB}
}

// Create 创建发布记录
func (r *PublishRepository) Create(record *model.PublishRecord) error {
	return r.db.Create(record).Error
}

// FindByID 根据ID查找发布记录
func (r *PublishRepository) FindByID(id uint) (*model.PublishRecord, error) {
	var record model.PublishRecord
	err := r.db.Preload("Content").First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// ListByUserID 获取用户的发布记录列表
func (r *PublishRepository) ListByUserID(userID uint, offset, limit int) ([]model.PublishRecord, int64, error) {
	var records []model.PublishRecord
	var total int64
	
	r.db.Model(&model.PublishRecord{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Preload("Content").Where("user_id = ?", userID).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&records).Error
	
	return records, total, err
}

// Update 更新发布记录
func (r *PublishRepository) Update(record *model.PublishRecord) error {
	return r.db.Save(record).Error
}

// GetPendingRecords 获取待发布的记录
func (r *PublishRepository) GetPendingRecords() ([]model.PublishRecord, error) {
	var records []model.PublishRecord
	err := r.db.Where("status = ? AND scheduled_at <= NOW()", 0).
		Preload("Content").Preload("User").Find(&records).Error
	return records, err
}
