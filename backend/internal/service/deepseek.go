// Package service 提供业务逻辑层
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"xiaohongshu/internal/config"
	"xiaohongshu/pkg/errno"
)

// ContentItem 内容项结构
type ContentItem struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// DeepSeekMessage DeepSeek消息结构
type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekRequest DeepSeek请求结构
type DeepSeekRequest struct {
	Model    string           `json:"model"`
	Messages []DeepSeekMessage `json:"messages"`
}

// DeepSeekChoice DeepSeek选择结构
type DeepSeekChoice struct {
	Index        int            `json:"index"`
	Message      DeepSeekMessage `json:"message"`
	FinishReason string         `json:"finish_reason"`
}

// DeepSeekUsage DeepSeek使用量结构
type DeepSeekUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// DeepSeekResponse DeepSeek响应结构
type DeepSeekResponse struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created int64           `json:"created"`
	Model   string          `json:"model"`
	Choices []DeepSeekChoice `json:"choices"`
	Usage   DeepSeekUsage   `json:"usage"`
}

// GeneratedContent 生成的内容结构
type GeneratedContent struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// AIService AI服务
type AIService struct {
	defaultCfg     *config.DeepSeekConfig
	tokenUsageSvc  *TokenUsageService
}

// NewAIService 创建AI服务实例
func NewAIService() *AIService {
	var defaultCfg *config.DeepSeekConfig
	if config.AppConfig != nil {
		defaultCfg = &config.AppConfig.DeepSeek
	} else {
		// 如果配置未初始化，使用默认值
		defaultCfg = &config.DeepSeekConfig{
			Model:   "deepseek-chat",
			BaseURL: "https://api.deepseek.com",
		}
	}
	return &AIService{
		defaultCfg:    defaultCfg,
		tokenUsageSvc: NewTokenUsageService(),
	}
}

// GenerateXiaohongshuContent 生成小红书内容
func (s *AIService) GenerateXiaohongshuContent(skillContent string, count int, length string, userAPIKey, userBaseURL, userModel string) ([]ContentItem, error) {
	apiKey := s.defaultCfg.APIKey
	baseURL := s.defaultCfg.BaseURL
	model := s.defaultCfg.Model
	
	if userAPIKey != "" {
		apiKey = userAPIKey
	}
	if userBaseURL != "" {
		baseURL = userBaseURL
	}
	if userModel != "" {
		model = userModel
	}
	
	if apiKey == "" {
		return nil, errno.ServiceUnavailable.WithMessage("DeepSeek API Key未配置")
	}

	// 构建提示词
	lengthDesc := map[string]string{
		"short":  "简短精炼（约100-200字）",
		"medium": "中等长度（约300-500字）",
		"long":   "详细完整（约600-800字）",
	}[length]
	if lengthDesc == "" {
		lengthDesc = "中等长度（约300-500字）"
	}

	prompt := fmt.Sprintf(`你是一个专业的小红书文案写手。请根据以下技能内容，生成%d个小红书风格的文案。

技能内容：%s

要求：
1. 每个文案包含：标题（不超过20个字符）、正文（%s）、标签（不超过20个，每个标签不超过10个字符）
2. 标题要吸引人，符合小红书风格
3. 正文要使用emoji，分段清晰，语气亲切
4. 标签要相关且热门
5. 请以JSON数组格式返回，格式如下：
[
  {
    "title": "标题",
    "description": "正文内容",
    "tags": ["标签1", "标签2"]
  }
]

只返回JSON，不要有其他文字说明。`, count, skillContent, lengthDesc)

	// 调用DeepSeek API
	messages := []DeepSeekMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := s.callDeepSeekAPI(messages, apiKey, baseURL, model)
	
	// 记录Token使用情况
	var promptTokens, completionTokens int
	if err == nil && response != nil {
		promptTokens = response.Usage.PromptTokens
		completionTokens = response.Usage.CompletionTokens
	}
	
	if err != nil {
		// 记录失败的请求
		go s.RecordUsage(0, model, "deepseek", "generate_content", prompt, "failed", err.Error(), "", "", promptTokens, completionTokens)
		return nil, err
	}

	// 解析响应
	if len(response.Choices) == 0 {
		go s.RecordUsage(0, model, "deepseek", "generate_content", prompt, "failed", "empty response", "", "", promptTokens, completionTokens)
		return nil, errno.GenerateFailed
	}

	content := response.Choices[0].Message.Content
	
	// 记录成功的请求
	go s.RecordUsage(0, model, "deepseek", "generate_content", prompt, "success", "", "", "", promptTokens, completionTokens)
	
	// 解析JSON
	var items []ContentItem
	err = json.Unmarshal([]byte(content), &items)
	if err != nil {
		return nil, errno.GenerateFailed.WithMessage("解析AI响应失败")
	}

	// 验证生成的内容
	for i := range items {
		if len([]rune(items[i].Title)) > 20 {
			items[i].Title = string([]rune(items[i].Title)[:20])
		}
		if len(items[i].Tags) > 20 {
			items[i].Tags = items[i].Tags[:20]
		}
		for j := range items[i].Tags {
			if len([]rune(items[i].Tags[j])) > 10 {
				items[i].Tags[j] = string([]rune(items[i].Tags[j])[:10])
			}
		}
	}

	return items, nil
}

// callDeepSeekAPI 调用DeepSeek API，返回响应和使用量
func (s *AIService) callDeepSeekAPI(messages []DeepSeekMessage, apiKey, baseURL, model string) (*DeepSeekResponse, error) {
	reqBody := DeepSeekRequest{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/chat/completions", baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败: %s", string(body))
	}

	var result DeepSeekResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// RecordUsage 记录Token使用情况
func (s *AIService) RecordUsage(
	userID uint,
	model, provider, requestType, requestContent, responseStatus, errorMessage, ipAddress, userAgent string,
	promptTokens, completionTokens int,
) error {
	if s.tokenUsageSvc == nil {
		s.tokenUsageSvc = NewTokenUsageService()
	}
	return s.tokenUsageSvc.RecordTokenUsage(
		userID, model, provider, requestType, requestContent, responseStatus, errorMessage, ipAddress, userAgent,
		promptTokens, completionTokens,
	)
}
