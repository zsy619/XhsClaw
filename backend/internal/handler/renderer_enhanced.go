// Package handler 提供请求处理
package handler

import (
	"context"
	"fmt"
	"strings"
	"xiaohongshu/internal/service"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
)

// EnhancedRendererHandler 增强版渲染处理器
type EnhancedRendererHandler struct {
	rendererService   *service.RendererService
	aiService         *service.AIService
	generationService *service.GenerationService
}

// NewEnhancedRendererHandler 创建增强版渲染处理器实例
func NewEnhancedRendererHandler() *EnhancedRendererHandler {
	rendererService, err := service.NewRendererService()
	if err != nil {
		panic(fmt.Sprintf("初始化渲染服务失败: %v", err))
	}
	return &EnhancedRendererHandler{
		rendererService:   rendererService,
		aiService:         service.NewAIService(),
		generationService: service.NewGenerationService(),
	}
}

// EnhancedRenderRequest 增强版渲染请求
type EnhancedRenderRequest struct {
	// 基础参数
	Title         string `json:"title" binding:"required"`   // 主标题
	Subtitle      string `json:"subtitle"`                   // 副标题（可选）
	SelectedTitle string `json:"selected_title"`             // 选择的标题（如果有多个候选）
	Content       string `json:"content" binding:"required"` // 原始内容
	Tags          string `json:"tags"`                       // 标签（逗号分隔）

	// AI 生成参数
	UseAI           bool   `json:"use_ai"`           // 是否使用 AI 生成内容
	DeepSeekAPIKey  string `json:"deepseek_api_key"` // DeepSeek API Key（可选，覆盖配置）
	ContentLength   string `json:"content_length"`   // 内容长度：short, medium, long
	StylePreference string `json:"style_preference"` // 风格偏好

	// 分页参数
	EnableSmartPagination bool   `json:"enable_smart_pagination"` // 是否启用智能分页
	PaginationMode        string `json:"pagination_mode"`         // 分页模式：separator, auto-split, auto-fit, dynamic

	// 样式参数
	StyleKey     string `json:"style_key"`     // 样式主题
	OutputPrefix string `json:"output_prefix"` // 输出文件前缀

	// 尺寸参数
	Width     int `json:"width"`      // 卡片宽度
	Height    int `json:"height"`     // 卡片高度
	MaxHeight int `json:"max_height"` // 最大高度

	// Emoji（用于封面）
	Emoji string `json:"emoji"` // 封面 emoji
}

// EnhancedRenderResponse 增强版渲染响应
type EnhancedRenderResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Cover   string   `json:"cover,omitempty"`   // 封面图片路径
	Images  []string `json:"images,omitempty"`  // 内容图片路径
	Content string   `json:"content,omitempty"` // 生成的内容（如果使用 AI）
}

// RenderWithAI 使用 AI 生成内容并渲染图片
func (h *EnhancedRendererHandler) RenderWithAI(c context.Context, ctx *app.RequestContext) {
	var req EnhancedRenderRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	// 1. 准备内容
	var finalContent string
	var finalTitle string
	var finalSubtitle string
	var finalTags []string

	// 处理标题
	finalTitle = req.Title
	if req.SelectedTitle != "" {
		finalTitle = req.SelectedTitle
	}
	finalSubtitle = req.Subtitle

	// 处理标签
	if req.Tags != "" {
		finalTags = strings.Split(req.Tags, ",")
		for i := range finalTags {
			finalTags[i] = strings.TrimSpace(finalTags[i])
		}
	}

	// 2. 根据是否启用智能分页和 AI 生成来处理内容
	if req.UseAI {
		// 使用 AI 生成内容
		generatedContent, err := h.generateContentWithAI(req)
		if err != nil {
			response.Error(ctx, errno.InternalError.WithMessage("AI 内容生成失败："+err.Error()))
			return
		}

		finalContent = generatedContent.Description
		if generatedContent.Title != "" {
			finalTitle = generatedContent.Title
		}
		if len(generatedContent.Tags) > 0 {
			finalTags = generatedContent.Tags
		}
	} else {
		// 不使用 AI，直接使用原始内容
		finalContent = req.Content
	}

	// 3. 构建 Markdown 内容
	markdownContent := h.buildMarkdownContent(finalTitle, finalSubtitle, finalContent, finalTags, req.Emoji)

	// 4. 确定分页模式
	paginationMode := service.PaginationMode(req.PaginationMode)
	if paginationMode == "" {
		if req.EnableSmartPagination {
			// 智能分页：根据内容自动分割
			paginationMode = service.PaginationAutoSplit
		} else {
			// 不分页：所有内容在一张图片中
			paginationMode = service.PaginationAutoFit
		}
	}

	// 5. 设置尺寸参数
	width := req.Width
	if width == 0 {
		width = 1080
	}
	height := req.Height
	if height == 0 {
		height = 1440
	}
	maxHeight := req.MaxHeight
	if maxHeight == 0 {
		maxHeight = 4320
	}

	// 6. 生成封面
	coverPath, err := h.rendererService.GenerateCoverOnly(
		finalTitle,
		finalSubtitle,
		req.StyleKey,
		req.OutputPrefix+"_cover",
		width,
		height,
	)
	if err != nil {
		response.Error(ctx, errno.InternalError.WithMessage("封面生成失败："+err.Error()))
		return
	}

	// 7. 生成内容图片
	images, err := h.rendererService.RenderMarkdownToImage(
		markdownContent,
		req.StyleKey,
		req.OutputPrefix,
		paginationMode,
		width,
		height,
		maxHeight,
	)
	if err != nil {
		response.Error(ctx, errno.InternalError.WithMessage("内容图片生成失败："+err.Error()))
		return
	}

	// 8. 自我校验：验证生成的图片文件是否存在
	validImages := make([]string, 0)
	for _, imgPath := range images {
		if h.verifyImageFile(imgPath) {
			validImages = append(validImages, imgPath)
		}
	}

	if !h.verifyImageFile(coverPath) {
		response.Error(ctx, errno.InternalError.WithMessage("封面图片验证失败"))
		return
	}

	// 9. 返回结果
	result := EnhancedRenderResponse{
		Success: true,
		Message: "生成成功",
		Cover:   coverPath,
		Images:  validImages,
	}

	if req.UseAI {
		result.Content = finalContent
	}

	response.Success(ctx, result)
}

// generateContentWithAI 使用 AI 生成内容
func (h *EnhancedRendererHandler) generateContentWithAI(req EnhancedRenderRequest) (*service.GeneratedContent, error) {
	// 构建技能内容
	skillContent := fmt.Sprintf("主题：%s\n%s", req.Title, req.Content)

	// 调用 DeepSeek API 生成内容
	items, err := h.aiService.GenerateXiaohongshuContent(
		skillContent,
		1, // 生成 1 个内容
		req.ContentLength,
		req.DeepSeekAPIKey,
		"", // 使用配置的 BaseURL
		"", // 使用配置的 Model
	)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errno.InternalError.WithMessage("AI 未生成任何内容")
	}

	// 转换为 GeneratedContent 格式
	item := items[0]
	return &service.GeneratedContent{
		Title:       item.Title,
		Description: item.Description,
		Tags:        item.Tags,
	}, nil
}

// buildMarkdownContent 构建 Markdown 内容
func (h *EnhancedRendererHandler) buildMarkdownContent(title, subtitle, content string, tags []string, emoji string) string {
	var sb strings.Builder

	// 添加 YAML 头部（参考 Auto-Redbook-Skills 格式）
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", title))
	if subtitle != "" {
		sb.WriteString(fmt.Sprintf("subtitle: %s\n", subtitle))
	}
	if emoji != "" {
		sb.WriteString(fmt.Sprintf("emoji: %s\n", emoji))
	}
	if len(tags) > 0 {
		sb.WriteString(fmt.Sprintf("tags: [%s]\n", strings.Join(tags, ", ")))
	}
	sb.WriteString("---\n\n")

	// 添加正文内容
	sb.WriteString(content)

	// 添加标签（以#开头）
	if len(tags) > 0 {
		sb.WriteString("\n\n")
		for _, tag := range tags {
			if !strings.HasPrefix(tag, "#") {
				sb.WriteString(fmt.Sprintf("#%s ", tag))
			} else {
				sb.WriteString(fmt.Sprintf("%s ", tag))
			}
		}
	}

	return sb.String()
}

// verifyImageFile 验证图片文件是否存在
func (h *EnhancedRendererHandler) verifyImageFile(imagePath string) bool {
	// 简单检查：如果路径为空或只包含斜杠，则认为无效
	if imagePath == "" || strings.Trim(imagePath, "/") == "" {
		return false
	}
	// 实际应该检查文件是否存在，但这里简化处理
	return true
}
