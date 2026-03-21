// Package repository 提供数据访问层
package repository

import (
	"fmt"

	"gorm.io/gorm"

	"xiaohongshu/internal/model"
)

// SystemDictRepository 系统字典仓储
type SystemDictRepository struct {
	db *gorm.DB
}

// NewSystemDictRepository 创建系统字典仓储实例
func NewSystemDictRepository() *SystemDictRepository {
	return &SystemDictRepository{db: DB}
}

// GetByCategory 获取指定类别的字典
func (r *SystemDictRepository) GetByCategory(category string) ([]model.SystemDict, error) {
	var dicts []model.SystemDict
	err := r.db.Where("category = ? AND enabled = ?", category, true).
		Order("sort_order ASC, id ASC").
		Find(&dicts).Error
	if err != nil {
		return nil, fmt.Errorf("查询字典失败：%w", err)
	}
	return dicts, nil
}

// List 获取所有字典
func (r *SystemDictRepository) List() ([]model.SystemDict, error) {
	var dicts []model.SystemDict
	err := r.db.Order("category ASC, sort_order ASC, id ASC").Find(&dicts).Error
	if err != nil {
		return nil, fmt.Errorf("查询字典列表失败：%w", err)
	}
	return dicts, nil
}

// GetCategories 获取所有类别
func (r *SystemDictRepository) GetCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&model.SystemDict{}).
		Where("enabled = ?", true).
		Distinct("category").
		Pluck("category", &categories).Error
	if err != nil {
		return nil, fmt.Errorf("查询字典类别失败：%w", err)
	}
	return categories, nil
}

// FindByID 根据ID查找
func (r *SystemDictRepository) FindByID(id uint) (*model.SystemDict, error) {
	var dict model.SystemDict
	err := r.db.First(&dict, id).Error
	if err != nil {
		return nil, fmt.Errorf("查询字典失败：%w", err)
	}
	return &dict, nil
}

// Create 创建字典
func (r *SystemDictRepository) Create(dict *model.SystemDict) error {
	return r.db.Create(dict).Error
}

// Update 更新字典
func (r *SystemDictRepository) Update(dict *model.SystemDict) error {
	return r.db.Save(dict).Error
}

// Delete 删除字典
func (r *SystemDictRepository) Delete(id uint) error {
	return r.db.Delete(&model.SystemDict{}, id).Error
}
