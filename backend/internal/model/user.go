package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;comment:用户ID"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100;comment:邮箱"`
	Password  string         `json:"-" gorm:"size:255;not null;comment:密码"`
	Nickname  string         `json:"nickname" gorm:"size:50;not null;comment:昵称"`
	Avatar    string         `json:"avatar" gorm:"size:255;comment:头像URL"`
	RoleID    uint           `json:"role_id" gorm:"index;default:3;comment:角色ID"`
	Role      *Role          `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Status    int            `json:"status" gorm:"default:1;comment:状态 1:正常 0:禁用"`
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
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// AuthUserInfo 认证用户信息
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
	ID        uint      `json:"id" gorm:"primaryKey;comment:记录ID"`
	Token     string    `json:"token" gorm:"uniqueIndex;size:500;not null;comment:Token值"`
	UserID    uint      `json:"user_id" gorm:"index;not null;comment:用户ID"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index;not null;comment:过期时间"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (TokenBlacklist) TableName() string {
	return "token_blacklists"
}

// Role 角色模型
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey;comment:角色ID"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null;comment:角色名称"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null;comment:角色代码"`
	Description string         `json:"description" gorm:"size:255;comment:角色描述"`
	Permissions string         `json:"permissions" gorm:"type:text;comment:权限列表JSON"`
	IsSystem    bool           `json:"is_system" gorm:"default:false;comment:是否系统角色"`
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
	ID          uint      `json:"id" gorm:"primaryKey;comment:权限ID"`
	Name        string    `json:"name" gorm:"size:50;not null;comment:权限名称"`
	Code        string    `json:"code" gorm:"uniqueIndex;size:100;not null;comment:权限代码"`
	Module      string    `json:"module" gorm:"size:50;not null;comment:所属模块"`
	Description string    `json:"description" gorm:"size:255;comment:权限描述"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}
