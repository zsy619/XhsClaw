// Package model 数据模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// PublishConfig 发布配置模型
type PublishConfig struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	UserID             uint           `json:"user_id" gorm:"uniqueIndex;not null;comment:用户ID"`
	DefaultPublishTime string         `json:"default_publish_time" gorm:"size:10;comment:默认发布时间"`
	AutoPublishEnabled bool           `json:"auto_publish_enabled" gorm:"default:false;comment:是否启用自动发布"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (PublishConfig) TableName() string {
	return "publish_configs"
}
