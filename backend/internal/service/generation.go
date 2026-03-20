// Package service 提供业务逻辑层
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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

// GetGenerationInstructions 获取小红书文案生成指令（共用函数）
func GetGenerationInstructions() string {
	return `生成顺序要求：
1. 第一步：生成封面建议
	目标：提供直观、有吸引力的封面设计方案，在海量笔记中脱颖而出。
	要求：
	- 提供 1-2个具体的封面方向（如：真人出镜对比图、产品大图+醒目文字、高颜值场景摆拍、干货信息汇总图）。
	- 如果涉及文字，请提供 3-5个封面文案关键词或短句。文案需突出痛点、卖点或引发好奇心，字体要大而醒目。
	示例：
	- 方向：真人出镜，展示使用前后的惊人对比。
	- 文案关键词："塌鼻星人有救了！"、"3分钟！get立体侧颜"、"手残党必备"
2. 第二步：生成标题
	目标：创作一个在20个字符内，能瞬间抓住用户眼球、激发点击欲望的爆款标题，巧用emoji表情。
	要求：
	- 字符数：严格控制在 20个字符以内（包括标点符号），确保在首页推荐流中完整显示。
	- 吸引力法则：灵活运用悬念、痛点、数据、干货、情绪、反常识、人群标签等技巧。
	- 关键词植入：巧妙融入核心关键词，便于系统推荐和用户搜索。
	示例：
	- 数据型：3天瘦腿2cm，我只做了这件事！
	- 悬念型：千万别买这双鞋，因为实在太显腿长了！
	- 痛点型：毛孔粗大星人，这瓶"磨皮水"给我锁死！
3. 第三步：生成正文内容
	目标：用亲切、真诚、有价值的语言留住用户，引发互动（点赞、收藏、评论）。
	要求：
	- 语气亲切： 像朋友间分享一样自然，多用"你"、"姐妹"、"真的绝了"、"谁懂啊"等口语化表达，拉近距离。
	- 巧用Emoji： 在段落开头、重点内容、产品/步骤前合理使用Emoji进行点缀，增加可读性和生动性。切勿滥用，保持排版清爽。
	- 结构清晰：
	- 开头（钩子）： 用1-2句话承接标题，简单引入主题，或再次强化痛点/吸引力。
	- 中间（干货/分享）：
	分点阐述，逻辑清晰。可以使用数字序号或Emoji小标题（如：🍱今日饮食 | 🏃‍♀️运动打卡 | 💡小贴士）。
	- 如果是教程，步骤要详细、可操作。
	- 如果是测评，结论要客观、优缺点分明。
	- 如果是Vlog，故事要有趣或有共鸣点。
	- 结尾（互动/引导）： 总结感受，并引导用户互动。例如："姐妹们快冲！"、"你们还有什么私藏好物？评论区分享给我吧～"、"如果对你有用，别忘了点个❤️哦！"
	内容价值： 确保笔记有核心价值，无论是情绪价值（让人开心、感动）、实用价值（学到东西）还是审美价值（看着舒服）。
4. 第四步：生成标签
	目标：通过精准的标签，让笔记被更多目标用户看到，进入流量池。
	要求：
	- 数量：建议 5-8个，最多不超过10个。
	- 结构（金字塔型）：
		- 1-2个核心标签： 与笔记内容最相关的、热度较高的大类标签（如：#护肤、#穿搭、#美食教程）。
		- 3-5个细分标签： 更具体的长尾标签，精准定位人群（如：#干皮护肤、#小个子穿搭、#空气炸锅食谱）。
		- 1-2个氛围/场景标签： 拓展曝光渠道（如：#好物分享、#我的日常、#独居生活图鉴、#氛围感）。
	- 格式：每个标签以"#"开头，标签内容中不加空格，控制在 15个字符以内，简洁明了。
5. 第五步：生成emoji表情
	目标：为笔记添加一个有吸引力的emoji表情，增加点击率。
	要求：
	- 选择一个与笔记内容相关、符合目标受众的emoji表情。
	- 避免使用通用的emoji表情，如"👍"、"❤️"等。
	- 考虑用户的情感反应，选择一个能触发情感共鸣的emoji。`
}

// GetGenerationResponseFormat 获取生成响应的JSON格式要求（共用函数）
func GetGenerationResponseFormat() string {
	return `请以JSON格式返回，格式如下：
{
  "generated_emoji": "emoji表情",
  "generated_title": "标题",
  "generated_content": "正文内容",
  "generated_tags": ["标签1", "标签2"],
  "cover_suggestion": "封面文案"
}

只返回JSON，不要有其他文字说明。`
}

// GenerateContent 生成内容
func (s *GenerationService) GenerateContent(userID uint, req *model.GenerationRequest, ipAddress, userAgent string) (*model.GenerationResponse, error) {
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
		fmt.Println("apiKey:-->", apiKey)
		return s.generateMockContent(req, nil)
	}

	lengthDesc := s.getLengthDescription(req.Length)
	styleDesc := s.getStyleDescription(req.StylePreference)

	prompt := fmt.Sprintf(`你是一个专业的小红书文案写手。请根据以下信息，按照指定顺序生成小红书风格的文案。

主题内容：%s
内容风格：%s
目标受众：%s
内容长度：%s

%s

%s`, req.Keywords, styleDesc, req.TargetAudience, lengthDesc, GetGenerationInstructions(), GetGenerationResponseFormat())

	return s.callGenerationAPI(userID, modelName, "generate_content", prompt, apiKey, baseURL, req, ipAddress, userAgent)
}

// RewriteContent 改写内容
func (s *GenerationService) RewriteContent(userID uint, req *model.RewriteRequest, ipAddress, userAgent string) (*model.GenerationResponse, error) {
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
	_ = s.getStyleDescription(req.StylePreference)

	prompt := fmt.Sprintf(`请帮我改写以下小红书文案，按照指定顺序生成。

原文案：
%s

%s

内容长度：%s
%s

%s`, req.Content, GetGenerationInstructions(), lengthDesc, func() string {
		if req.PreserveKeyInfo {
			return "请保留原文案的关键信息"
		}
		return "可以自由发挥，但保持核心主题"
	}(), GetGenerationResponseFormat())

	messages := []DeepSeekMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := s.aiService.callDeepSeekAPI(messages, apiKey, baseURL, modelName)
	if err != nil {
		// 记录失败的请求（使用默认值）
		go s.recordTokenUsage(context.Background(), userID, modelName, "rewrite_content", prompt, "failed", err.Error(), ipAddress, userAgent, 0, 0)
		return s.rewriteMockContent(req, err)
	}

	if len(response.Choices) == 0 {
		// 记录空响应
		go s.recordTokenUsage(context.Background(), userID, modelName, "rewrite_content", prompt, "failed", "empty response", ipAddress, userAgent, response.Usage.PromptTokens, response.Usage.CompletionTokens)
		return s.rewriteMockContent(req)
	}

	content := response.Choices[0].Message.Content

	// 记录成功的请求
	go s.recordTokenUsage(context.Background(), userID, modelName, "rewrite_content", prompt, "success", "", ipAddress, userAgent, response.Usage.PromptTokens, response.Usage.CompletionTokens)

	var result model.GenerationResponse
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		// 记录解析失败
		go s.recordTokenUsage(context.Background(), userID, modelName, "rewrite_content", prompt, "failed", "json unmarshal failed", ipAddress, userAgent, response.Usage.PromptTokens, response.Usage.CompletionTokens)
		return s.rewriteMockContent(req)
	}

	return &result, nil
}

// callGenerationAPI 调用生成API的共用方法
func (s *GenerationService) callGenerationAPI(userID uint, modelName, requestType, prompt, apiKey, baseURL string, req *model.GenerationRequest, ipAddress, userAgent string) (*model.GenerationResponse, error) {
	messages := []DeepSeekMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := s.aiService.callDeepSeekAPI(messages, apiKey, baseURL, modelName)
	if err != nil {
		// 记录失败的请求（使用默认值）
		go s.recordTokenUsage(context.Background(), userID, modelName, requestType, prompt, "failed", err.Error(), ipAddress, userAgent, 0, 0)
		return s.generateMockContent(req, err)
	}

	if len(response.Choices) == 0 {
		// 记录空响应
		go s.recordTokenUsage(context.Background(), userID, modelName, requestType, prompt, "failed", "empty response", ipAddress, userAgent, response.Usage.PromptTokens, response.Usage.CompletionTokens)
		return s.generateMockContent(req, fmt.Errorf("记录空响应"))
	}

	content := response.Choices[0].Message.Content

	// 记录成功的请求
	go s.recordTokenUsage(context.Background(), userID, modelName, requestType, prompt, "success", "", ipAddress, userAgent, response.Usage.PromptTokens, response.Usage.CompletionTokens)

	var result model.GenerationResponse
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		// 记录解析失败
		go s.recordTokenUsage(context.Background(), userID, modelName, requestType, prompt, "failed", "json unmarshal failed", ipAddress, userAgent, response.Usage.PromptTokens, response.Usage.CompletionTokens)
		return s.generateMockContent(req, fmt.Errorf("记录解析失败"))
	}

	return &result, nil
}

// recordTokenUsage 记录Token使用情况（支持token参数）
func (s *GenerationService) recordTokenUsage(ctx context.Context, userID uint, model, requestType, requestContent, responseStatus, errorMessage, ipAddress, userAgent string, promptTokens, completionTokens int) {
	s.tokenUsageSvc.RecordTokenUsage(
		ctx,
		userID,
		model,
		"deepseek",
		requestType,
		requestContent,
		responseStatus,
		errorMessage,
		ipAddress,
		userAgent,
		promptTokens,
		completionTokens,
	)
}

// generateMockContent 生成模拟内容（当API不可用时）
func (s *GenerationService) generateMockContent(req *model.GenerationRequest, err error) (*model.GenerationResponse, error) {
	if err != nil {
		return &model.GenerationResponse{
			GeneratedContent: err.Error(),
			GeneratedTitle:   "生成失败",
			GeneratedTags:    []string{"生成失败"},
			CoverSuggestion:  "生成失败 | 请联系管理员",
		}, nil
	}
	generatedTags := []string{req.Keywords, "实用技巧", "干货分享", "建议收藏", "关注不亏"}

	return &model.GenerationResponse{
		GeneratedContent: fmt.Sprintf(`⚠️ 系统提示：未配置大模型参数

您还未配置大模型API参数，当前返回的是模拟内容。

请前往系统设置页面配置以下参数：
1. API Key
2. Base URL
3. 模型名称

配置完成后，系统将使用真实的AI模型生成更优质的内容。

✨ 模拟内容：

## Title
%s超全攻略！建议收藏💯

## Tags
%s

---

超实用%s分享！

今天给大家带来超棒的%s心得～

📝 首先要注意以下几点：
1️⃣ 第一点真的很重要
2️⃣ 第二点也不能忽略
3️⃣ 第三点是关键

💡 最后给大家一些小建议：
- 建议1
- 建议2
- 建议3

希望对大家有帮助哦～`, req.Keywords, strings.Join(generatedTags, " | "), req.Keywords, req.Keywords),
		GeneratedTitle:  fmt.Sprintf("%s超全攻略！", req.Keywords),
		GeneratedTags:   generatedTags,
		CoverSuggestion: fmt.Sprintf("%s | 实用分享", req.Keywords),
	}, nil
}

// rewriteMockContent 改写模拟内容
func (s *GenerationService) rewriteMockContent(req *model.RewriteRequest, err ...error) (*model.GenerationResponse, error) {
	generatedTags := []string{"改写文案", "内容优化", "小红书技巧", "文案创作", "干货分享"}

	if len(err) > 0 {
		return &model.GenerationResponse{
			GeneratedContent: err[0].Error(),
			GeneratedTitle:   "改写后的标题",
			GeneratedTags:    generatedTags,
			CoverSuggestion:  "内容改写 | 精彩呈现",
		}, nil
	}

	return &model.GenerationResponse{
		GeneratedContent: fmt.Sprintf(`⚠️ 系统提示：未配置大模型参数

您还未配置大模型API参数，当前返回的是模拟内容。

请前往系统设置页面配置以下参数：
1. API Key
2. Base URL
3. 模型名称

配置完成后，系统将使用真实的AI模型生成更优质的内容。

✨ 模拟改写内容：

## Title
改写后的精彩标题✨

## Tags
%s

---

✨ 改写版本来啦！

%s

希望这个改写版本对您有帮助～`, strings.Join(generatedTags, " | "), req.Content),
		GeneratedTitle:  "改写后的标题",
		GeneratedTags:   generatedTags,
		CoverSuggestion: "内容改写 | 精彩呈现",
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
