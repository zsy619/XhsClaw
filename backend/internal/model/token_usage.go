package model

import (
	"time"
)

// TokenUsage 大模型Token使用记录
type TokenUsage struct {
	ID              uint      `gorm:"primarykey;comment:记录ID" json:"id"`
	UserID          uint      `gorm:"index;comment:用户ID" json:"user_id"`
	Username        string    `gorm:"-" json:"username"`
	Model           string    `gorm:"size:100;comment:使用的模型" json:"model"`
	Provider        string    `gorm:"size:50;comment:提供商" json:"provider"`
	InputTokens     int       `gorm:"column:input_tokens;comment:输入tokens" json:"input_tokens"`
	OutputTokens    int       `gorm:"column:output_tokens;comment:输出tokens" json:"output_tokens"`
	TotalTokens     int       `gorm:"comment:总tokens" json:"total_tokens"`
	Cost            float64   `gorm:"type:decimal(10,6);comment:费用(美元)" json:"cost"`
	RequestType     string    `gorm:"size:50;comment:请求类型" json:"request_type"`
	RequestContent  string    `gorm:"type:text;comment:请求内容摘要" json:"request_content"`
	ResponseStatus  string    `gorm:"size:20;comment:响应状态" json:"response_status"`
	ErrorMessage    string    `gorm:"type:text;comment:错误信息" json:"error_message,omitempty"`
	IPAddress       string    `gorm:"size:50;comment:IP地址" json:"ip_address,omitempty"`
	UserAgent       string    `gorm:"size:500;comment:用户代理" json:"user_agent,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName 指定表名
func (TokenUsage) TableName() string {
	return "token_usage"
}

// TokenUsageStats Token使用统计
type TokenUsageStats struct {
	TotalRequests     int64   `json:"total_requests"`
	SuccessRequests   int64   `json:"success_requests"`
	FailedRequests    int64   `json:"failed_requests"`
	TotalPromptTokens int64   `json:"total_prompt_tokens"`
	TotalCompletionTokens int64 `json:"total_completion_tokens"`
	TotalTokens       int64   `json:"total_tokens"`
	TotalCost         float64 `json:"total_cost"`
	AverageTokens     float64 `json:"average_tokens"`
}

// UserTokenUsage 用户Token使用统计
type UserTokenUsage struct {
	Date         string  `json:"date"`
	TotalTokens  int64   `json:"total_tokens"`
	TotalCost    float64 `json:"total_cost"`
	RequestCount int64   `json:"request_count"`
}

// TokenUsageByModel 按模型统计
type TokenUsageByModel struct {
	Model        string  `json:"model"`
	TotalTokens  int64   `json:"total_tokens"`
	TotalCost    float64 `json:"total_cost"`
	RequestCount int64   `json:"request_count"`
}
