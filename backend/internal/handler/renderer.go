// Package handler 提供请求处理
package handler

import (
	"context"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"

	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/service"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"
)

// RendererHandler 渲染处理器
type RendererHandler struct {
	rendererService *service.RendererService
}

// NewRendererHandler 创建渲染处理器实例
func NewRendererHandler() *RendererHandler {
	return &RendererHandler{
		rendererService: service.NewRendererService(),
	}
}

// StyleInfo 样式信息
type StyleInfo struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// RenderRequest 渲染请求
type RenderRequest struct {
	MarkdownContent      string `json:"markdown_content" binding:"required"`
	StyleKey           string `json:"style_key"`
	OutputPrefix       string `json:"output_prefix"`
	Mode               string `json:"mode"` // separator, auto-fit, auto-split, dynamic
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	MaxHeight          int    `json:"max_height"`
	CardWidth          int    `json:"card_width"`           // 兼容前端字段
	CardHeight         int    `json:"card_height"`          // 兼容前端字段
	MaxContentHeight   int    `json:"max_content_height"`   // 兼容前端字段
	EnableSmartPagination bool `json:"enable_smart_pagination"` // 兼容前端字段
}

// CoverRequest 封面生成请求
type CoverRequest struct {
	Title        string `json:"title" binding:"required"`
	Subtitle     string `json:"subtitle"`
	StyleKey     string `json:"style_key"`
	OutputPrefix string `json:"output_prefix"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

// RenderResponse 渲染响应
type RenderResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Images  []string    `json:"images"`
	Styles  []StyleInfo `json:"styles,omitempty"`
}

// CoverResponse 封面响应
type CoverResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Image   string `json:"image"`
}

// GetRendererStyles 获取渲染样式
func (h *RendererHandler) GetRendererStyles(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	styleConfigs := h.rendererService.GetStyles()
	styles := make([]StyleInfo, len(styleConfigs))
	for i, sc := range styleConfigs {
		styles[i] = StyleInfo{
			Key: func() string {
				switch sc.Name {
				case "小红书红":
					return "xiaohongshu"
				case "紫韵":
					return "purple"
				case "清新薄荷":
					return "mint"
				case "日落橙":
					return "sunset"
				case "深海蓝":
					return "ocean"
				case "优雅白":
					return "elegant"
				case "暗黑模式":
					return "dark"
				default:
					return "xiaohongshu"
				}
			}(),
			Name: sc.Name,
		}
	}

	response.Success(ctx, RenderResponse{
		Success: true,
		Message: "获取成功",
		Styles:  styles,
	})
}

// RenderMarkdown 渲染Markdown为图片
func (h *RendererHandler) RenderMarkdown(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	var req RenderRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	// 兼容前端字段 - 优先使用 card_width 等字段
	width := req.Width
	if req.CardWidth > 0 {
		width = req.CardWidth
	}
	if width == 0 {
		width = 1080 // 默认宽度
	}

	height := req.Height
	if req.CardHeight > 0 {
		height = req.CardHeight
	}
	if height == 0 {
		height = 1440 // 默认高度
	}

	maxHeight := req.MaxHeight
	if req.MaxContentHeight > 0 {
		maxHeight = req.MaxContentHeight
	}
	if maxHeight == 0 {
		maxHeight = height - 340 // 默认最大内容高度
	}

	images, err := h.rendererService.RenderMarkdownToImage(
		req.MarkdownContent,
		req.StyleKey,
		req.OutputPrefix,
		width,
		height,
		maxHeight,
	)
	if err != nil {
		response.Error(ctx, errno.InternalError.WithMessage(err.Error()))
		return
	}

	response.Success(ctx, RenderResponse{
		Success: true,
		Message: "渲染成功",
		Images:  images,
	})
}

// GenerateCover 生成封面
func (h *RendererHandler) GenerateCover(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	var req CoverRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	imagePath, err := h.rendererService.GenerateCoverOnly(
		req.Title,
		req.Subtitle,
		req.StyleKey,
		req.OutputPrefix,
		req.Width,
		req.Height,
	)
	if err != nil {
		response.Error(ctx, errno.InternalError.WithMessage(err.Error()))
		return
	}

	response.Success(ctx, CoverResponse{
		Success: true,
		Message: "封面生成成功",
		Image:   imagePath,
	})
}

// GetRenderedImage 获取渲染后的图片
func GetRenderedImage(c context.Context, ctx *app.RequestContext) {
	filename := ctx.Param("filename")
	if filename == "" {
		response.ParamError(ctx, "文件名不能为空")
		return
	}

	// 安全地构建文件路径
	imagePath := filepath.Join("./public/images", filename)

	// 返回图片文件
	ctx.File(imagePath)
}
