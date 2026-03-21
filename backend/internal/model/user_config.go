// Package model 定义数据模型
package model

// UserConfigRequest 用户配置请求（用于 API 请求）
// 注意：UserConfig 模型已迁移到独立的表中（llm_providers, xiaohongshu_configs, publish_configs）
// 此处保留 Request 类型用于 API 请求兼容
type UserConfigRequest struct {
	// LLM 配置
	LLMAPIKey  string `json:"llm_api_key"`
	LLMBaseURL string `json:"llm_base_url"`
	LLMModel   string `json:"llm_model"`

	// 小红书配置
	XiaohongshuCookie string `json:"xiaohongshu_cookie"`
	XiaohongshuUserId string `json:"xiaohongshu_user_id"`
	XiaohongshuToken  string `json:"xiaohongshu_token"`

	// 发布配置
	DefaultPublishTime string `json:"default_publish_time"`
	AutoPublishEnabled bool   `json:"auto_publish_enabled"`
}
