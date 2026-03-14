// Package model 定义数据模型
package model

// GenerationRequest 主题生成请求
type GenerationRequest struct {
	Keywords        string `json:"keywords" binding:"required"`
	StylePreference string `json:"style_preference"`
	TargetAudience  string `json:"target_audience"`
	Length          int    `json:"length"`
}

// GenerationResponse 主题生成响应
type GenerationResponse struct {
	GeneratedContent string   `json:"generated_content"`
	GeneratedTitle   string   `json:"generated_title"`
	GeneratedTags    []string `json:"generated_tags"`
}

// RewriteRequest 改写请求
type RewriteRequest struct {
	Content          string `json:"content" binding:"required"`
	StylePreference  string `json:"style_preference"`
	PreserveKeyInfo  bool   `json:"preserve_key_info"`
	Length           int    `json:"length"`
}
