// Package repository 提供数据访问层
package repository

import (
	"xiaohongshu/internal/model"

	"gorm.io/gorm"
)

// RoleRepository 角色仓库
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓库实例
func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		db: GetDB(),
	}
}

// Create 创建角色
func (r *RoleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

// FindByID 根据ID查找角色
func (r *RoleRepository) FindByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByCode 根据代码查找角色
func (r *RoleRepository) FindByCode(code string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (r *RoleRepository) List(offset, pageSize int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	err := r.db.Model(&model.Role{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Order("id ASC").Offset(offset).Limit(pageSize).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// ListAll 获取所有角色
func (r *RoleRepository) ListAll() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Order("id ASC").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// PermissionRepository 权限仓库
type PermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓库实例
func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{
		db: GetDB(),
	}
}

// Create 创建权限
func (r *PermissionRepository) Create(permission *model.Permission) error {
	return r.db.Create(permission).Error
}

// ListAll 获取所有权限
func (r *PermissionRepository) ListAll() ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Order("module ASC, id ASC").Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// ListByModule 根据模块获取权限
func (r *PermissionRepository) ListByModule(module string) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Where("module = ?", module).Order("id ASC").Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
