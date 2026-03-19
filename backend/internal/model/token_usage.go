// Package model 提供数据模型定义
package model

import (
	"time"
)

// TokenUsageTokenUsage 大模型Token使用记录
type TokenUsage struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	UserID          uint      `gorm:"index" json:"user_id"`                    // 用户ID
	Username        string    `gorm:"-" json:"username"`                         // 用户名（不存储，用于展示）
	Model           string    `gorm:"size:100" json:"model"`                    // 使用的模型
	Provider        string    `gorm:"size:50" json:"provider"`                   // 提供商 (deepseek/openai/anthropic等)
	InputTokens     int       `gorm:"column:input_tokens" json:"input_tokens"`   // 输入tokens
	OutputTokens    int       `gorm:"column:output_tokens" json:"output_tokens"`  // 输出tokens
	TotalTokens     int       `json:"total_tokens"`                             // 总tokens
	Cost            float64   `gorm:"type:decimal(10,6)" json:"cost"`           // 费用（美元）
	RequestType     string    `gorm:"size:50" json:"request_type"`             // 请求类型 (generate_content/render_image等)
	RequestContent  string    `gorm:"type:text" json:"request_content"`         // 请求内容摘要
	ResponseStatus  string    `gorm:"size:20" json:"response_status"`           // 响应状态 (success/failed)
	ErrorMessage    string    `gorm:"type:text" json:"error_message,omitempty"` // 错误信息
	IPAddress       string    `gorm:"size:50" json:"ip_address,omitempty"`     // IP地址
	UserAgent       string    `gorm:"size:500" json:"user_agent,omitempty"`     // 用户代理
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName 指定表名
func (TokenUsage) TableName() string {
	return "token_usage"
}

// TokenUsageStats Token使用统计
type TokenUsageStats struct {
	TotalRequests    int64   `json:"total_requests"`     // 总请求数
	SuccessRequests  int64   `json:"success_requests"`   // 成功请求数
	FailedRequests   int64   `json:"failed_requests"`    // 失败请求数
	TotalPromptTokens    int64   `json:"total_prompt_tokens"`    // 总输入tokens
	TotalCompletionTokens int64   `json:"total_completion_tokens"` // 总输出tokens
	TotalTokens      int64   `json:"total_tokens"`       // 总tokens
	TotalCost        float64 `json:"total_cost"`        // 总费用
	AverageTokens    float64 `json:"average_tokens"`     // 平均每次请求tokens
}

// UserTokenUsage 用户Token使用统计
type UserTokenUsage struct {
	Date            string  `json:"date"`              // 日期
	TotalTokens     int64   `json:"total_tokens"`      // 当日总tokens
	TotalCost       float64 `json:"total_cost"`        // 当日总费用
	RequestCount    int64   `json:"request_count"`    // 当日请求数
}

// TokenUsageByModel 按模型统计
type TokenUsageByModel struct {
	Model       string  `json:"model"`        // 模型名称
	TotalTokens int64   `json:"total_tokens"` // 总tokens
	TotalCost   float64 `json:"total_cost"`   // 总费用
	RequestCount int64  `json:"request_count"` // 请求次数
}
