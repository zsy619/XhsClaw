// Package model 定义数据模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// LLMProvider 大模型服务商配置表
type LLMProvider struct {
	ID          uint           `json:"id" gorm:"primaryKey;comment:配置ID"`
	UserID      uint           `json:"user_id" gorm:"index;not null;comment:用户ID"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name        string         `json:"name" gorm:"size:100;not null;comment:配置名称"`
	Provider    string         `json:"provider" gorm:"size:50;not null;comment:服务商类型 openai/deepseek/自定义"`
	APIKey      string         `json:"api_key" gorm:"size:500;comment:API密钥"`
	BaseURL     string         `json:"base_url" gorm:"size:500;comment:API地址"`
	ModelName   string         `json:"model_name" gorm:"size:100;comment:模型名称"`
	IsDefault   bool           `json:"is_default" gorm:"default:false;comment:是否默认"`
	IsEnabled   bool           `json:"is_enabled" gorm:"default:true;comment:是否启用"`
	Timeout     int            `json:"timeout" gorm:"default:60;comment:超时时间(秒)"`
	RetryCount  int            `json:"retry_count" gorm:"default:3;comment:重试次数"`
	Extra       string         `json:"extra" gorm:"type:text;comment:扩展配置JSON"`
	Description string         `json:"description" gorm:"size:255;comment:描述"`
	SortOrder   int            `json:"sort_order" gorm:"default:0;comment:排序"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (LLMProvider) TableName() string {
	return "llm_providers"
}

// ProviderType 提供商类型常量
const (
	// OpenAI 系列
	ProviderOpenAI = "openai" // OpenAI GPT系列模型
	ProviderAzure  = "azure"  // Azure OpenAI Service

	// Anthropic 系列
	ProviderClaude = "claude" // Claude 系列模型 (Anthropic)

	// Google 系列
	ProviderGemini = "gemini" // Google Gemini 系列模型

	// DeepSeek 系列
	ProviderDeepSeek = "deepseek" // DeepSeek 系列模型

	// 国内主流服务商
	ProviderQianfan  = "qianfan"  // 百度文心一言
	ProviderQwen     = "qwen"     // 阿里通义千问
	ProviderSpark    = "spark"    // 讯飞星火大模型
	ProviderGLM      = "glm"      // 智谱 GLM 系列
	ProviderHunyuan  = "hunyuan"  // 腾讯混元大模型
	ProviderDoubao   = "doubao"   // 字节豆包大模型
	ProviderBaichuan = "baichuan" // 百川大模型
	ProviderMiniMax  = "minimax"  // MiniMax 大模型

	// 其他国际服务商
	ProviderMistral    = "mistral"    // Mistral AI
	ProviderCohere     = "cohere"     // Cohere
	ProviderGroq       = "groq"       // Groq (高速推理)
	ProviderReplicate  = "replicate"  // Replicate 云服务
	ProviderPerplexity = "perplexity" // Perplexity AI

	// 开源/自托管
	ProviderOllama   = "ollama"   // Ollama 本地模型
	ProviderLMStudio = "lmstudio" // LM Studio 本地服务
	ProviderLocalAI  = "localai"  // LocalAI 自托管
	ProvidervLLM     = "vllm"     // vLLM 自托管推理服务

	// 自定义/通用
	ProviderCustom = "custom" // 自定义 API (兼容任意 OpenAI 格式接口)
)

// ProviderDisplayNames 提供商显示名称映射
var ProviderDisplayNames = map[string]string{
	ProviderOpenAI:     "OpenAI",
	ProviderAzure:      "Azure OpenAI",
	ProviderClaude:     "Claude (Anthropic)",
	ProviderGemini:     "Google Gemini",
	ProviderDeepSeek:   "DeepSeek",
	ProviderQianfan:    "文心一言 (百度)",
	ProviderQwen:       "通义千问 (阿里)",
	ProviderSpark:      "讯飞星火",
	ProviderGLM:        "智谱 GLM",
	ProviderHunyuan:    "腾讯混元",
	ProviderDoubao:     "字节豆包",
	ProviderBaichuan:   "百川大模型",
	ProviderMiniMax:    "MiniMax",
	ProviderMistral:    "Mistral AI",
	ProviderCohere:     "Cohere",
	ProviderGroq:       "Groq",
	ProviderReplicate:  "Replicate",
	ProviderPerplexity: "Perplexity",
	ProviderOllama:     "Ollama (本地)",
	ProviderLMStudio:   "LM Studio (本地)",
	ProviderLocalAI:    "LocalAI (自托管)",
	ProvidervLLM:       "vLLM (自托管)",
	ProviderCustom:     "自定义",
}

// ProviderBaseURLs 提供商默认 API 地址映射
var ProviderBaseURLs = map[string]string{
	ProviderOpenAI:     "https://api.openai.com/v1",
	ProviderAzure:      "https://{your-resource-name}.openai.azure.com",
	ProviderClaude:     "https://api.anthropic.com/v1",
	ProviderGemini:     "https://generativelanguage.googleapis.com/v1beta",
	ProviderDeepSeek:   "https://api.deepseek.com/v1",
	ProviderQianfan:    "https://qianfan.baidubce.com/v2",
	ProviderQwen:       "https://dashscope.aliyuncs.com/compatible-mode/v1",
	ProviderSpark:      "https://spark-api.xf-yun.com/v3.5",
	ProviderGLM:        "https://open.bigmodel.cn/api/paas/v4",
	ProviderHunyuan:    "https://hunyuan.cloud.tencent.com/v1",
	ProviderDoubao:     "https://ark.cn-beijing.volces.com/api/v3",
	ProviderBaichuan:   "https://api.baichuan-ai.com/v1",
	ProviderMiniMax:    "https://api.minimax.chat/v1",
	ProviderMistral:    "https://api.mistral.ai/v1",
	ProviderCohere:     "https://api.cohere.ai/v2",
	ProviderGroq:       "https://api.groq.com/openai/v1",
	ProviderReplicate:  "https://api.replicate.com/v1",
	ProviderPerplexity: "https://api.perplexity.ai",
	ProviderOllama:     "http://localhost:11434/v1",
	ProviderLMStudio:   "http://localhost:1234/v1",
	ProviderLocalAI:    "http://localhost:8080/v1",
	ProvidervLLM:       "http://localhost:8000/v1",
}

// GetProviderDisplayName 获取提供商的显示名称
func GetProviderDisplayName(provider string) string {
	if name, ok := ProviderDisplayNames[provider]; ok {
		return name
	}
	return provider
}

// GetProviderBaseURL 获取提供商的默认 API 地址
func GetProviderBaseURL(provider string) string {
	if url, ok := ProviderBaseURLs[provider]; ok {
		return url
	}
	return ""
}

// IsLocalProvider 判断是否为本地/自托管提供商
func IsLocalProvider(provider string) bool {
	switch provider {
	case ProviderOllama, ProviderLMStudio, ProviderLocalAI, ProvidervLLM:
		return true
	}
	return false
}

// LLMProviderRequest 服务商配置请求
type LLMProviderRequest struct {
	Name        string `json:"name" binding:"required"`
	Provider    string `json:"provider" binding:"required"`
	APIKey      string `json:"api_key"`
	BaseURL     string `json:"base_url"`
	ModelName   string `json:"model_name"`
	IsDefault   bool   `json:"is_default"`
	IsEnabled   bool   `json:"is_enabled"`
	Timeout     int    `json:"timeout"`
	RetryCount  int    `json:"retry_count"`
	Extra       string `json:"extra"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// LLMProviderResponse 服务商配置响应
type LLMProviderResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	BaseURL     string `json:"base_url"`
	ModelName   string `json:"model_name"`
	IsDefault   bool   `json:"is_default"`
	IsEnabled   bool   `json:"is_enabled"`
	Timeout     int    `json:"timeout"`
	RetryCount  int    `json:"retry_count"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
}
