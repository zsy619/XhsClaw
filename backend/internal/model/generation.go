// Package model 定义数据模型
package model

// PaginationMode 分页模式
type PaginationMode string

const (
	// PaginationModeSeparator 按 --- 分隔手动分页
	PaginationModeSeparator PaginationMode = "separator"
	// PaginationModeAutoFit 固定尺寸，自动整体缩放内容，避免溢出/大面积留白
	PaginationModeAutoFit PaginationMode = "auto-fit"
	// PaginationModeAutoSplit 根据渲染后高度自动拆分为多张卡片
	PaginationModeAutoSplit PaginationMode = "auto-split"
	// PaginationModeDynamic 根据内容动态调整图片高度
	PaginationModeDynamic PaginationMode = "dynamic"
)

// GenerationRequest 主题生成请求
type GenerationRequest struct {
	Keywords        string        `json:"keywords" binding:"required"`
	StylePreference string        `json:"style_preference"`
	TargetAudience  string        `json:"target_audience"`
	Length          int           `json:"length"`
	// 分页相关参数
	EnablePagination bool           `json:"enable_pagination"` // 是否启用分页
	PaginationMode   PaginationMode `json:"pagination_mode"`  // 分页模式：separator, auto-fit, auto-split, dynamic
	CardWidth       int            `json:"card_width"`      // 卡片宽度（像素）
	CardHeight      int            `json:"card_height"`    // 卡片高度（像素）
	Theme           string         `json:"theme"`          // 主题名称
}

// GenerationResponse 主题生成响应
type GenerationResponse struct {
	GeneratedContent  string   `json:"generated_content"`
	GeneratedTitle     string   `json:"generated_title"`
	GeneratedTags      []string `json:"generated_tags"`
	CoverSuggestion    string   `json:"cover_suggestion"` // 封面建议文案
}

// RewriteRequest 改写请求
type RewriteRequest struct {
	Content          string `json:"content" binding:"required"`
	StylePreference  string `json:"style_preference"`
	PreserveKeyInfo  bool   `json:"preserve_key_info"`
	Length           int    `json:"length"`
}
