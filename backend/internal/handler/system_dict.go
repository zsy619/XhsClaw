// Package handler 提供请求处理
package handler

import (
	"context"
	"strconv"

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

// List 获取所有字典数据（支持分页）
func (h *SystemDictHandler) List(c context.Context, ctx *app.RequestContext) {
	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	dicts, total, err := h.dictRepo.List(page, pageSize)
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	// 返回分页结果
	response.Success(ctx, map[string]interface{}{
		"items": dicts,
		"total": total,
	})
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
