// Package model 定义数据模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// XiaohongshuConfig 小红书账号配置表
type XiaohongshuConfig struct {
	ID          uint           `json:"id" gorm:"primaryKey;comment:配置ID"`
	UserID      uint           `json:"user_id" gorm:"index;not null;comment:用户ID"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name        string         `json:"name" gorm:"size:100;not null;comment:配置名称"`
	Cookie      string         `json:"cookie" gorm:"type:text;comment:Cookie"`
	XHSUserID   string         `json:"xhs_user_id" gorm:"column:xhs_user_id;size:100;comment:小红书用户ID"`
	Token       string         `json:"token" gorm:"size:500;comment:Token"`
	DeviceID    string         `json:"device_id" gorm:"size:100;comment:设备ID"`
	IsDefault   bool           `json:"is_default" gorm:"default:false;comment:是否默认"`
	IsEnabled   bool           `json:"is_enabled" gorm:"default:true;comment:是否启用"`
	Status      string         `json:"status" gorm:"size:20;default:pending;comment:状态 pending/active/expired/error"`
	LastLoginAt *time.Time     `json:"last_login_at" gorm:"comment:最后登录时间"`
	ErrorMsg    string         `json:"error_msg" gorm:"type:text;comment:错误信息"`
	Extra       string         `json:"extra" gorm:"type:text;comment:扩展配置JSON"`
	Description string         `json:"description" gorm:"size:255;comment:描述"`
	SortOrder   int            `json:"sort_order" gorm:"default:0;comment:排序"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (XiaohongshuConfig) TableName() string {
	return "xiaohongshu_configs"
}

// ConfigStatus 配置状态常量
const (
	XHSStatusPending = "pending" // 待验证
	XHSStatusActive  = "active"  // 正常
	XHSStatusExpired = "expired" // 已过期
	XHSStatusError   = "error"   // 错误
)

// XiaohongshuConfigRequest 小红书配置请求
type XiaohongshuConfigRequest struct {
	Name        string `json:"name" binding:"required"`
	Cookie      string `json:"cookie"`
	XHSUserID   string `json:"xhs_user_id"`
	Token       string `json:"token"`
	DeviceID    string `json:"device_id"`
	IsDefault   bool   `json:"is_default"`
	IsEnabled   bool   `json:"is_enabled"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// XiaohongshuConfigResponse 小红书配置响应
type XiaohongshuConfigResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	XHSUserID   string  `json:"xhs_user_id"`
	IsDefault   bool    `json:"is_default"`
	IsEnabled   bool    `json:"is_enabled"`
	Status      string  `json:"status"`
	LastLoginAt *string `json:"last_login_at"`
	Description string  `json:"description"`
	SortOrder   int     `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
}

// PublishConfig 发布配置（保留在用户配置中）
type PublishConfig struct {
	ID                   uint      `json:"id" gorm:"primaryKey;comment:配置ID"`
	UserID               uint      `json:"user_id" gorm:"uniqueIndex;not null;comment:用户ID"`
	DefaultPublishTime   string    `json:"default_publish_time" gorm:"size:20;comment:默认发布时间 HH:mm"`
	AutoPublishEnabled   bool      `json:"auto_publish_enabled" gorm:"default:false;comment:是否启用自动发布"`
	DefaultXHSConfigID   uint      `json:"default_xhs_config_id" gorm:"comment:默认小红书配置ID"`
	DefaultLLMProviderID uint      `json:"default_llm_provider_id" gorm:"comment:默认LLM配置ID"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// TableName 指定表名
func (PublishConfig) TableName() string {
	return "publish_configs"
}

// PublishConfigRequest 发布配置请求
type PublishConfigRequest struct {
	DefaultPublishTime   string `json:"default_publish_time"`
	AutoPublishEnabled   bool   `json:"auto_publish_enabled"`
	DefaultXHSConfigID   uint   `json:"default_xhs_config_id"`
	DefaultLLMProviderID uint   `json:"default_llm_provider_id"`
}
