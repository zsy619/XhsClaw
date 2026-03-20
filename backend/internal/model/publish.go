package model

import (
	"time"

	"gorm.io/gorm"
)

// 发布记录状态常量
const (
	// PublishStatusPending 待发布
	PublishStatusPending = 0
	// PublishStatusPublishing 发布中
	PublishStatusPublishing = 1
	// PublishStatusSuccess 发布成功
	PublishStatusSuccess = 2
	// PublishStatusFailed 发布失败
	PublishStatusFailed = 3
	// PublishStatusCancelled 已取消
	PublishStatusCancelled = 4
)

// 内容状态常量
const (
	// ContentStatusDraft 草稿
	ContentStatusDraft = 0
	// ContentStatusPending 待发布
	ContentStatusPending = 1
	// ContentStatusPublishing 发布中
	ContentStatusPublishing = 2
	// ContentStatusPublished 已发布
	ContentStatusPublished = 3
	// ContentStatusPublishFailed 发布失败
	ContentStatusPublishFailed = 4
)

// PublishRecord 发布记录模型
type PublishRecord struct {
	ID          uint           `json:"id" gorm:"primaryKey;comment:发布记录ID"`
	UserID      uint           `json:"user_id" gorm:"index;not null;comment:用户ID"`
	ContentID   uint           `json:"content_id" gorm:"index;not null;comment:内容ID"`
	Status      int            `json:"status" gorm:"default:0;comment:状态 0:待发布 1:发布中 2:成功 3:失败 4:已取消"`
	ErrorMsg    string         `json:"error_msg" gorm:"type:text;comment:错误信息"`
	ScheduledAt time.Time      `json:"scheduled_at" gorm:"comment:计划发布时间"`
	PublishedAt *time.Time    `json:"published_at" gorm:"comment:实际发布时间"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Content     Content        `json:"content,omitempty" gorm:"foreignKey:ContentID"`
}

// TableName 指定表名
func (PublishRecord) TableName() string {
	return "publish_records"
}

// SchedulePublishRequest 定时发布请求
type SchedulePublishRequest struct {
	ContentID   uint   `json:"content_id" binding:"required"`
	PublishTime string `json:"publish_time" binding:"required"`
	Frequency   string `json:"frequency"`
}

// PublishRequest 立即发布请求
type PublishRequest struct {
	ContentID uint `json:"content_id" binding:"required"`
}
