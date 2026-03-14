// Package model 定义数据模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// UserConfig 用户配置模型
type UserConfig struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	
	// 大模型配置
	LLMAPIKey   string `json:"llm_api_key" gorm:"size:500"`
	LLMBaseURL  string `json:"llm_base_url" gorm:"size:500"`
	LLMModel    string `json:"llm_model" gorm:"size:100"`
	
	// 小红书配置
	XiaohongshuCookie string `json:"xiaohongshu_cookie" gorm:"type:text"`
	XiaohongshuUserId string `json:"xiaohongshu_user_id" gorm:"size:100"`
	XiaohongshuToken  string `json:"xiaohongshu_token" gorm:"size:500"`
	
	// 发布配置
	DefaultPublishTime string `json:"default_publish_time" gorm:"size:20"` // 默认发布时间 HH:mm
	AutoPublishEnabled bool   `json:"auto_publish_enabled" gorm:"default:false"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (UserConfig) TableName() string {
	return "user_configs"
}

// UserConfigRequest 用户配置请求
type UserConfigRequest struct {
	LLMAPIKey           string `json:"llm_api_key"`
	LLMBaseURL          string `json:"llm_base_url"`
	LLMModel            string `json:"llm_model"`
	XiaohongshuCookie   string `json:"xiaohongshu_cookie"`
	XiaohongshuUserId   string `json:"xiaohongshu_user_id"`
	XiaohongshuToken    string `json:"xiaohongshu_token"`
	DefaultPublishTime   string `json:"default_publish_time"`
	AutoPublishEnabled  bool   `json:"auto_publish_enabled"`
}
