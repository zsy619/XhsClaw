// Package repository 提供数据访问层
package repository

import (
	"xiaohongshu/internal/model"

	"gorm.io/gorm"
)

// ContentRepository 内容仓储
type ContentRepository struct {
	db *gorm.DB
}

// NewContentRepository 创建内容仓储实例
func NewContentRepository() *ContentRepository {
	return &ContentRepository{db: DB}
}

// Create 创建内容
func (r *ContentRepository) Create(content *model.Content) error {
	return r.db.Create(content).Error
}

// FindByID 根据ID查找内容
func (r *ContentRepository) FindByID(id uint) (*model.Content, error) {
	var content model.Content
	err := r.db.First(&content, id).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// FindByUserIDAndID 根据用户ID和内容ID查找
func (r *ContentRepository) FindByUserIDAndID(userID, id uint) (*model.Content, error) {
	var content model.Content
	err := r.db.Where("user_id = ? AND id = ?", userID, id).First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// ListByUserID 获取用户的内容列表
func (r *ContentRepository) ListByUserID(userID uint, offset, limit int, status *int) ([]model.Content, int64, error) {
	var contents []model.Content
	var total int64
	
	query := r.db.Model(&model.Content{}).Where("user_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	
	query.Count(&total)
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&contents).Error
	
	return contents, total, err
}

// Update 更新内容
func (r *ContentRepository) Update(content *model.Content) error {
	return r.db.Save(content).Error
}

// Delete 删除内容
func (r *ContentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Content{}, id).Error
}
