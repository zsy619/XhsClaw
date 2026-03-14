// Package model 定义数据模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100"`
	Password  string         `json:"-" gorm:"size:255;not null"` // 不返回密码
	Nickname  string         `json:"nickname" gorm:"size:50;not null"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	RoleID    uint           `json:"role_id" gorm:"index;default:3"` // 关联角色ID
	Role      *Role          `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Status    int            `json:"status" gorm:"default:1"` // 1:正常, 0:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required,min=3,max=50"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token      string `json:"token"`
	User       User   `json:"user"`
	// 兼容前端格式
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// AuthUserInfo 认证用户信息（兼容前端格式）
type AuthUserInfo struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	LastLogin string `json:"last_login,omitempty"`
	Role      string `json:"role,omitempty"`
}

// TokenBlacklist Token黑名单模型
type TokenBlacklist struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"uniqueIndex;size:500;not null"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (TokenBlacklist) TableName() string {
	return "token_blacklists"
}

// Role 角色模型
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Description string         `json:"description" gorm:"size:255"`
	Permissions string         `json:"permissions" gorm:"type:text"` // JSON格式存储权限数组
	IsSystem    bool           `json:"is_system" gorm:"default:false"` // 是否系统角色（不可删除）
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:100;not null"`
	Module      string         `json:"module" gorm:"size:50;not null"` // 所属模块
	Description string         `json:"description" gorm:"size:255"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}
