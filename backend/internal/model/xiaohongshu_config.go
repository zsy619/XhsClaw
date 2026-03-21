// Package model 数据模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// XHSConfig 小红书配置模型
type XHSConfig struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"index;not null;comment:用户ID"`
	Name        string         `json:"name" gorm:"size:100;not null;comment:配置名称"`
	Cookie      string         `json:"-" gorm:"type:text;comment:Cookie"`
	Token       string         `json:"-" gorm:"type:text;comment:Token"`
	DeviceID    string         `json:"device_id" gorm:"size:100;comment:设备ID"`
	XHSUserID   string         `json:"xhs_user_id" gorm:"size:50;comment:小红书用户ID"`
	IsDefault   bool           `json:"is_default" gorm:"default:false;comment:是否默认配置"`
	IsEnabled   bool           `json:"is_enabled" gorm:"default:true;comment:是否启用"`
	Status      string         `json:"status" gorm:"size:20;default:'pending';comment:状态:pending/active/expired/error"`
	Description string         `json:"description" gorm:"size:255;comment:描述"`
	SortOrder   int            `json:"sort_order" gorm:"default:0;comment:排序"`
	LastLoginAt *time.Time     `json:"last_login_at,omitempty" gorm:"comment:最后登录时间"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (XHSConfig) TableName() string {
	return "xhs_configs"
}

// XHSConfigRequest 创建/更新小红书配置的请求
type XHSConfigRequest struct {
	Name        string `json:"name" binding:"required"`
	Cookie      string `json:"cookie"`
	Token       string `json:"token"`
	DeviceID    string `json:"device_id"`
	XHSUserID   string `json:"xhs_user_id"`
	IsDefault   bool   `json:"is_default"`
	IsEnabled   bool   `json:"is_enabled"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// XHSVerifyResponse 验证响应
type XHSVerifyResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
}
