// Package handler 提供请求处理
package handler

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/service"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"
)

// getClientIP 获取客户端真实 IP 地址
func getClientIP(ctx *app.RequestContext) string {
	// 优先从 X-Forwarded-For 获取
	xff := string(ctx.GetHeader("X-Forwarded-For"))
	if xff != "" {
		// X-Forwarded-For 可能包含多个 IP，取第一个
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// 其次从 X-Real-IP 获取
	xri := string(ctx.GetHeader("X-Real-IP"))
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// 最后使用 RemoteAddr
	return ctx.RemoteAddr().String()
}

// GenerationHandler 生成处理器
type GenerationHandler struct {
	generationService *service.GenerationService
}

// NewGenerationHandler 创建生成处理器实例
func NewGenerationHandler() *GenerationHandler {
	return &GenerationHandler{
		generationService: service.NewGenerationService(),
	}
}

// GenerateContent 生成内容
func (h *GenerationHandler) GenerateContent(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	// 获取真实 IP 地址
	ipAddress := getClientIP(ctx)

	// 获取 User-Agent
	userAgent := string(ctx.GetHeader("User-Agent"))

	var req model.GenerationRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	resp, err := h.generationService.GenerateContent(userID, &req, ipAddress, userAgent)
	if err != nil {
		response.Error(ctx, errno.GenerateFailed)
		return
	}

	response.Success(ctx, resp)
}

// RewriteContent 改写内容
func (h *GenerationHandler) RewriteContent(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	// 获取真实 IP 地址
	ipAddress := getClientIP(ctx)

	// 获取 User-Agent
	userAgent := string(ctx.GetHeader("User-Agent"))

	var req model.RewriteRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	resp, err := h.generationService.RewriteContent(userID, &req, ipAddress, userAgent)
	if err != nil {
		response.Error(ctx, errno.GenerateFailed)
		return
	}

	response.Success(ctx, resp)
}

// SaveContent 保存内容
func (h *GenerationHandler) SaveContent(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	var req model.ContentSaveRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	contentService := service.NewContentService()
	content, err := contentService.SaveContent(userID, &req)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, content)
}
