// Package service 提供业务逻辑层
package service

import (
	"encoding/json"
	"fmt"
	"xiaohongshu/internal/config"
	"xiaohongshu/internal/model"
)

// GenerationService 生成服务
type GenerationService struct {
	aiService         *AIService
	userConfigService *UserConfigService
	tokenUsageSvc     *TokenUsageService
}

// NewGenerationService 创建生成服务实例
func NewGenerationService() *GenerationService {
	return &GenerationService{
		aiService:         NewAIService(),
		userConfigService: NewUserConfigService(),
		tokenUsageSvc:     NewTokenUsageService(),
	}
}

// GenerateContent 生成内容
func (s *GenerationService) GenerateContent(userID uint, req *model.GenerationRequest) (*model.GenerationResponse, error) {
	// 获取用户配置，优先使用用户配置
	apiKey, baseURL, modelName := s.userConfigService.GetLLMConfig(userID)

	// 如果用户没有配置，使用系统默认配置
	if apiKey == "" && config.AppConfig != nil {
		apiKey = config.AppConfig.DeepSeek.APIKey
	}
	if baseURL == "" && config.AppConfig != nil {
		baseURL = config.AppConfig.DeepSeek.BaseURL
	}
	if modelName == "" && config.AppConfig != nil {
		modelName = config.AppConfig.DeepSeek.Model
	}

	if apiKey == "" {
		return s.generateMockContent(req)
	}

	lengthDesc := s.getLengthDescription(req.Length)
	styleDesc := s.getStyleDescription(req.StylePreference)

	prompt := fmt.Sprintf(`你是一个专业的小红书文案写手。请根据以下信息生成小红书风格的文案。

主题内容：%s
内容风格：%s
目标受众：%s
内容长度：%s

要求：
1. 生成吸引人的标题（不超过20个字符）
2. 生成正文内容，使用emoji，分段清晰，语气亲切
3. 生成相关的标签（不超过20个，每个标签不超过10个字符）

请以JSON格式返回，格式如下：
{
  "generated_title": "标题",
  "generated_content": "正文内容",
  "generated_tags": ["标签1", "标签2"]
}

只返回JSON，不要有其他文字说明。`, req.Keywords, styleDesc, req.TargetAudience, lengthDesc)

	messages := []DeepSeekMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := s.aiService.callDeepSeekAPI(messages, apiKey, baseURL, modelName)
	if err != nil {
		// 记录失败的请求
		go s.recordTokenUsage(userID, modelName, "generate_content", prompt, "failed", err.Error())
		return s.generateMockContent(req)
	}

	if len(response.Choices) == 0 {
		// 记录空响应
		go s.recordTokenUsage(userID, modelName, "generate_content", prompt, "failed", "empty response")
		return s.generateMockContent(req)
	}

	content := response.Choices[0].Message.Content

	// 记录成功的请求
	go s.recordTokenUsage(userID, modelName, "generate_content", prompt, "success", "")

	var result model.GenerationResponse
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		// 记录解析失败
		go s.recordTokenUsage(userID, modelName, "generate_content", prompt, "failed", "json unmarshal failed")
		return s.generateMockContent(req)
	}

	return &result, nil
}

// recordTokenUsage 记录Token使用情况
func (s *GenerationService) recordTokenUsage(userID uint, model, requestType, requestContent, responseStatus, errorMessage string) {
	s.tokenUsageSvc.RecordTokenUsage(
		userID,
		model,
		"deepseek",
		requestType,
		requestContent,
		responseStatus,
		errorMessage,
		"",
		"",
		0,
		0,
	)
}

// RewriteContent 改写内容
func (s *GenerationService) RewriteContent(userID uint, req *model.RewriteRequest) (*model.GenerationResponse, error) {
	// 获取用户配置，优先使用用户配置
	apiKey, baseURL, modelName := s.userConfigService.GetLLMConfig(userID)

	// 如果用户没有配置，使用系统默认配置
	if apiKey == "" && config.AppConfig != nil {
		apiKey = config.AppConfig.DeepSeek.APIKey
	}
	if baseURL == "" && config.AppConfig != nil {
		baseURL = config.AppConfig.DeepSeek.BaseURL
	}
	if modelName == "" && config.AppConfig != nil {
		modelName = config.AppConfig.DeepSeek.Model
	}

	if apiKey == "" {
		return s.rewriteMockContent(req)
	}

	lengthDesc := s.getLengthDescription(req.Length)
	styleDesc := s.getStyleDescription(req.StylePreference)

	prompt := fmt.Sprintf(`请帮我改写以下小红书文案。

原文案：
%s

要求：
1. 保持原文案风格：%s
2. 内容长度：%s
3. %s
4. 生成吸引人的标题（不超过20个字符）
5. 生成正文内容，使用emoji，分段清晰，语气亲切
6. 生成相关的标签（不超过20个，每个标签不超过10个字符）

请以JSON格式返回，格式如下：
{
  "generated_title": "标题",
  "generated_content": "正文内容",
  "generated_tags": ["标签1", "标签2"]
}

只返回JSON，不要有其他文字说明。`, req.Content, styleDesc, lengthDesc, func() string {
		if req.PreserveKeyInfo {
			return "请保留原文案的关键信息"
		}
		return "可以自由发挥"
	}())

	messages := []DeepSeekMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := s.aiService.callDeepSeekAPI(messages, apiKey, baseURL, modelName)
	if err != nil {
		// 记录失败的请求
		go s.recordTokenUsage(userID, modelName, "rewrite_content", prompt, "failed", err.Error())
		return s.rewriteMockContent(req)
	}

	if len(response.Choices) == 0 {
		// 记录空响应
		go s.recordTokenUsage(userID, modelName, "rewrite_content", prompt, "failed", "empty response")
		return s.rewriteMockContent(req)
	}

	content := response.Choices[0].Message.Content

	// 记录成功的请求
	go s.recordTokenUsage(userID, modelName, "rewrite_content", prompt, "success", "")

	var result model.GenerationResponse
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		// 记录解析失败
		go s.recordTokenUsage(userID, modelName, "rewrite_content", prompt, "failed", "json unmarshal failed")
		return s.rewriteMockContent(req)
	}

	return &result, nil
}

// generateMockContent 生成模拟内容（当API不可用时）
func (s *GenerationService) generateMockContent(req *model.GenerationRequest) (*model.GenerationResponse, error) {
	return &model.GenerationResponse{
		GeneratedContent: fmt.Sprintf(`✨ 超实用%s分享！

今天给大家带来超棒的%s心得～

📝 首先要注意以下几点：
1️⃣ 第一点真的很重要
2️⃣ 第二点也不能忽略
3️⃣ 第三点是关键

💡 最后给大家一些小建议：
- 建议1
- 建议2
- 建议3

希望对大家有帮助哦～
有问题可以在评论区留言！

喜欢记得点赞收藏关注三连～`, req.Keywords, req.Keywords),
		GeneratedTitle: fmt.Sprintf("%s超全攻略！", req.Keywords),
		GeneratedTags:  []string{req.Keywords, "干货分享", "生活小技巧", "必备", "推荐", "实用"},
	}, nil
}

// rewriteMockContent 改写模拟内容
func (s *GenerationService) rewriteMockContent(req *model.RewriteRequest) (*model.GenerationResponse, error) {
	return &model.GenerationResponse{
		GeneratedContent: req.Content + "\n\n✨ 改写版本来啦！",
		GeneratedTitle:   "改写后的标题",
		GeneratedTags:    []string{"改写", "文案", "小红书"},
	}, nil
}

// getLengthDescription 获取长度描述
func (s *GenerationService) getLengthDescription(length int) string {
	if length <= 100 {
		return "简短精炼（约100字）"
	} else if length <= 300 {
		return "中等长度（约300字）"
	} else if length <= 500 {
		return "详细完整（约500字）"
	}
	return "详细完整（约800字）"
}

// getStyleDescription 获取风格描述
func (s *GenerationService) getStyleDescription(style string) string {
	styleMap := map[string]string{
		"cute":         "活泼可爱",
		"professional": "专业严谨",
		"artistic":     "文艺清新",
		"humorous":     "幽默风趣",
		"informative":  "干货分享",
	}
	if desc, ok := styleMap[style]; ok {
		return desc
	}
	return style
}
