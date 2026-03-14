// Package handler 提供请求处理
package handler

import (
	"context"
	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/service"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
)

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

	var req model.GenerationRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	resp, err := h.generationService.GenerateContent(userID, &req)
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

	var req model.RewriteRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	resp, err := h.generationService.RewriteContent(userID, &req)
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
