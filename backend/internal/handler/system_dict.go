// Package handler 提供请求处理
package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"xiaohongshu/internal/repository"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"
)

// SystemDictHandler 系统字典处理器
type SystemDictHandler struct {
	dictRepo *repository.SystemDictRepository
}

// NewSystemDictHandler 创建系统字典处理器
func NewSystemDictHandler() *SystemDictHandler {
	return &SystemDictHandler{
		dictRepo: repository.NewSystemDictRepository(),
	}
}

// GetByCategory 获取指定类别的字典数据
func (h *SystemDictHandler) GetByCategory(c context.Context, ctx *app.RequestContext) {
	category := ctx.Param("category")
	if category == "" {
		response.Error(ctx, &errno.ErrNo{Code: 400, Message: "类别不能为空"})
		return
	}

	dicts, err := h.dictRepo.GetByCategory(category)
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, dicts)
}

// List 获取所有字典数据
func (h *SystemDictHandler) List(c context.Context, ctx *app.RequestContext) {
	dicts, err := h.dictRepo.List()
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, dicts)
}

// GetCategories 获取所有类别
func (h *SystemDictHandler) GetCategories(c context.Context, ctx *app.RequestContext) {
	categories, err := h.dictRepo.GetCategories()
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, categories)
}
