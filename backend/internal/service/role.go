// Package service 提供业务逻辑层
package service

import (
	"encoding/json"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/repository"
	"xiaohongshu/pkg/errno"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo       *repository.RoleRepository
	permissionRepo *repository.PermissionRepository
	userRepo       *repository.UserRepository
}

// NewRoleService 创建角色服务实例
func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo:       repository.NewRoleRepository(),
		permissionRepo: repository.NewPermissionRepository(),
		userRepo:       repository.NewUserRepository(),
	}
}

// ListRoles 获取角色列表
func (s *RoleService) ListRoles(page, pageSize int) ([]model.Role, int64, error) {
	return s.roleRepo.List(page, pageSize)
}

// ListAllRoles 获取所有角色
func (s *RoleService) ListAllRoles() ([]model.Role, error) {
	return s.roleRepo.ListAll()
}

// GetRole 获取角色详情
func (s *RoleService) GetRole(id uint) (*model.Role, error) {
	return s.roleRepo.FindByID(id)
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(name, code, description string, permissions []string) (*model.Role, error) {
	// 检查角色代码是否已存在
	_, err := s.roleRepo.FindByCode(code)
	if err == nil {
		return nil, errno.RoleAlreadyExists
	}

	permissionsJSON, _ := json.Marshal(permissions)

	role := &model.Role{
		Name:        name,
		Code:        code,
		Description: description,
		Permissions: string(permissionsJSON),
		IsSystem:    false,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, errno.InternalError
	}

	return role, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(id uint, name, description string, permissions []string) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, errno.RoleNotFound
	}

	// 系统角色不可修改名称和代码
	if !role.IsSystem {
		if name != "" {
			role.Name = name
		}
	}

	if description != "" {
		role.Description = description
	}

	if permissions != nil {
		permissionsJSON, _ := json.Marshal(permissions)
		role.Permissions = string(permissionsJSON)
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, errno.InternalError
	}

	return role, nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return errno.RoleNotFound
	}

	// 系统角色不可删除
	if role.IsSystem {
		return errno.CannotDeleteSystemRole
	}

	// 检查是否有用户使用该角色
	users, _, err := s.userRepo.List(0, 100)
	if err != nil {
		return errno.InternalError
	}

	for _, user := range users {
		if user.RoleID == id {
			return errno.RoleInUse
		}
	}

	return s.roleRepo.Delete(id)
}

// ListPermissions 获取所有权限
func (s *RoleService) ListPermissions() ([]model.Permission, error) {
	return s.permissionRepo.ListAll()
}

// ListPermissionsByModule 根据模块获取权限
func (s *RoleService) ListPermissionsByModule(module string) ([]model.Permission, error) {
	return s.permissionRepo.ListByModule(module)
}

// UpdateUserRole 更新用户角色
func (s *RoleService) UpdateUserRole(userID, roleID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errno.UserNotFound
	}

	// 验证角色是否存在
	_, err = s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, errno.RoleNotFound
	}

	user.RoleID = roleID

	if err := s.userRepo.Update(user); err != nil {
		return nil, errno.InternalError
	}

	return user, nil
}

// UpdateUserStatus 更新用户状态
func (s *RoleService) UpdateUserStatus(userID uint, status int) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errno.UserNotFound
	}

	user.Status = status

	if err := s.userRepo.Update(user); err != nil {
		return nil, errno.InternalError
	}

	return user, nil
}
