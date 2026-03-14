// Package model 定义数据模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// PublishRecord 发布记录模型
type PublishRecord struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"index;not null"`
	ContentID  uint           `json:"content_id" gorm:"index;not null"`
	Status     int            `json:"status" gorm:"default:0"` // 0:待发布, 1:发布中, 2:成功, 3:失败
	ErrorMsg   string         `json:"error_msg" gorm:"type:text"`
	ScheduledAt time.Time     `json:"scheduled_at"`
	PublishedAt *time.Time    `json:"published_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Content    Content        `json:"content,omitempty" gorm:"foreignKey:ContentID"`
}

// TableName 指定表名
func (PublishRecord) TableName() string {
	return "publish_records"
}

// SchedulePublishRequest 定时发布请求
type SchedulePublishRequest struct {
	ContentID   uint   `json:"content_id" binding:"required"`
	PublishTime string `json:"publish_time" binding:"required"` // RFC3339格式
	Frequency   string `json:"frequency"` // once, daily, weekly
}

// PublishRequest 立即发布请求
type PublishRequest struct {
	ContentID uint `json:"content_id" binding:"required"`
}
