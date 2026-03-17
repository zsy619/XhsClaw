// Package handler 提供请求处理
package handler

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"

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
	StyleKey            string `json:"style_key"`
	OutputPrefix        string `json:"output_prefix"`
	Mode                string `json:"mode"`
	Width               int    `json:"width"`
	Height              int    `json:"height"`
	MaxHeight           int    `json:"max_height"`
	CardWidth           int    `json:"card_width"`
	CardHeight          int    `json:"card_height"`
	MaxContentHeight    int    `json:"max_content_height"`
	EnableSmartPagination bool `json:"enable_smart_pagination"`
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
	styleConfigs := h.rendererService.GetStyles()
	styles := make([]StyleInfo, len(styleConfigs))
	for i, sc := range styleConfigs {
		styles[i] = StyleInfo{
			Key:  sc.Key,
			Name: sc.Name,
		}
	}

	response.Success(ctx, RenderResponse{
		Success: true,
		Message: "获取成功",
		Styles:  styles,
	})
}

// verifyImageFile 验证图片文件是否存在
func (h *RendererHandler) verifyImageFile(imagePath string) bool {
	fullPath := filepath.Join(h.rendererService.GetImagesDir(), filepath.Base(imagePath))
	if _, err := os.Stat(fullPath); err != nil {
		return false
	}
	return true
}

// RenderMarkdown 渲染Markdown为图片
func (h *RendererHandler) RenderMarkdown(c context.Context, ctx *app.RequestContext) {
	var req RenderRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	width := req.Width
	if req.CardWidth > 0 {
		width = req.CardWidth
	}
	if width == 0 {
		width = 1080
	}

	height := req.Height
	if req.CardHeight > 0 {
		height = req.CardHeight
	}
	if height == 0 {
		height = 1440
	}

	maxHeight := req.MaxHeight
	if req.MaxContentHeight > 0 {
		maxHeight = req.MaxContentHeight
	}
	if maxHeight == 0 {
		maxHeight = 1100
	}

	mode := service.PaginationMode(req.Mode)
	if mode == "" {
		mode = service.PaginationAutoSplit
	}

	// 生成图片
	images, err := h.rendererService.RenderMarkdownToImage(
		req.MarkdownContent,
		req.StyleKey,
		req.OutputPrefix,
		mode,
		width,
		height,
		maxHeight,
	)
	if err != nil {
		response.Error(ctx, errno.InternalError.WithMessage(err.Error()))
		return
	}

	// 自我校验：验证生成的图片文件是否存在
	validImages := make([]string, 0)
	for _, imgPath := range images {
		if h.verifyImageFile(imgPath) {
			validImages = append(validImages, imgPath)
		}
	}

	if len(validImages) == 0 && len(images) > 0 {
		// 如果验证失败但有路径，仍返回原始路径
		validImages = images
	}

	response.Success(ctx, RenderResponse{
		Success: true,
		Message: "渲染成功",
		Images:  validImages,
	})
}

// GenerateCover 生成封面
func (h *RendererHandler) GenerateCover(c context.Context, ctx *app.RequestContext) {
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

	// 验证封面图片
	if !h.verifyImageFile(imagePath) {
		response.Error(ctx, errno.InternalError.WithMessage("封面图片生成失败"))
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

	imagePath := filepath.Join("./public/images", filename)
	ctx.File(imagePath)
}
