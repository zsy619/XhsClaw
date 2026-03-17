// Package model 定义数据模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// ContentAttributes 内容属性
type ContentAttributes struct {
	ContentStyle   string   `json:"content_style"`
	CustomStyle    string   `json:"custom_style"`
	TargetAudience []string `json:"target_audience"`
}

// RenderAttributes 渲染属性
type RenderAttributes struct {
	ImageStyleTheme     string `json:"image_style_theme"`
	EnableSmartPagination bool `json:"enable_smart_pagination"`
	CardWidth         int    `json:"card_width"`
	CardHeight        int    `json:"card_height"`
}

// Content 内容模型（生成的小红书笔记）
type Content struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	UserID            uint           `json:"user_id" gorm:"index;not null"`
	Title             string         `json:"title" gorm:"size:50;not null"`
	TitleOptions      string         `json:"title_options" gorm:"type:text"` // JSON格式存储备选标题数组
	SelectedTitleIndex int            `json:"selected_title_index" gorm:"default:0"` // 选中的备选标题索引
	Description       string         `json:"description" gorm:"type:text;not null"`
	Tags              string         `json:"tags" gorm:"type:text"` // JSON格式存储标签数组
	Images            string         `json:"images" gorm:"type:text"` // JSON格式存储图片URL数组
	ContentAttributes string         `json:"content_attributes" gorm:"type:text"` // JSON格式存储内容属性
	RenderAttributes  string         `json:"render_attributes" gorm:"type:text"` // JSON格式存储渲染属性
	Status            int            `json:"status" gorm:"default:0"` // 0:草稿, 1:待发布, 2:已发布, 3:发布失败
	PublishTime       *time.Time     `json:"publish_time"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
	
	User              User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (Content) TableName() string {
	return "contents"
}

// ContentHistory 内容历史记录模型
type ContentHistory struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	ContentID         uint           `json:"content_id" gorm:"index;not null"`
	UserID            uint           `json:"user_id" gorm:"index;not null"`
	Type              string         `json:"type" gorm:"size:20;not null"` // 'create' | 'edit' | 'delete' | 'publish'
	Title             string         `json:"title" gorm:"size:50;not null"`
	TitleOptions      string         `json:"title_options" gorm:"type:text"` // JSON格式存储备选标题数组
	SelectedTitleIndex int           `json:"selected_title_index" gorm:"default:0"`
	Description       string         `json:"description" gorm:"type:text;not null"`
	Tags              string         `json:"tags" gorm:"type:text"` // JSON格式存储标签数组
	Images            string         `json:"images" gorm:"type:text"` // JSON格式存储图片URL数组
	ContentAttributes string         `json:"content_attributes" gorm:"type:text"` // JSON格式存储内容属性
	RenderAttributes  string         `json:"render_attributes" gorm:"type:text"` // JSON格式存储渲染属性
	ChangeReason      string         `json:"change_reason" gorm:"size:255"`
	CreatedAt         time.Time      `json:"created_at"`

	Content           Content        `json:"content,omitempty" gorm:"foreignKey:ContentID"`
	User              User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (ContentHistory) TableName() string {
	return "content_histories"
}

// GenerateContentRequest 生成内容请求
type GenerateContentRequest struct {
	SkillContent string `json:"skill_content" binding:"required"`
	Count        int    `json:"count" binding:"min=1,max=10"` // 生成数量
	Length       string `json:"length"` // short, medium, long
}

// GenerateContentResponse 生成内容响应
type GenerateContentResponse struct {
	Contents []ContentItem `json:"contents"`
}

// ContentItem 单个生成的内容项
type ContentItem struct {
	Title               string           `json:"title"`
	Description         string           `json:"description"`
	Tags                []string         `json:"tags"`
	ContentAttributes   ContentAttributes `json:"content_attributes"`
	RenderAttributes    RenderAttributes  `json:"render_attributes"`
}

// ContentSaveRequest 保存内容请求
type ContentSaveRequest struct {
	Title               string           `json:"title" binding:"required"`
	TitleOptions        []string         `json:"title_options"` // 备选标题数组
	SelectedTitleIndex  int              `json:"selected_title_index"` // 选中的备选标题索引
	Description         string           `json:"description" binding:"required"`
	Tags                []string         `json:"tags"`
	Images              []string         `json:"images"` // 生成的图片路径数组
	ContentAttributes   ContentAttributes `json:"content_attributes"`
	RenderAttributes    RenderAttributes  `json:"render_attributes"`
}

// UpdateContentRequest 更新内容请求
type UpdateContentRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Images      []string `json:"images"`
	Status      *int     `json:"status"`
	PublishTime *string  `json:"publish_time"` // RFC3339格式
}
