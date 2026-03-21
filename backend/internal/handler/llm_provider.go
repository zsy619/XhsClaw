// Package handler 提供请求处理
package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/repository"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"
)

// LLMProviderHandler 大模型配置处理器
type LLMProviderHandler struct {
	llmRepo *repository.LLMProviderRepository
}

// NewLLMProviderHandler 创建大模型配置处理器
func NewLLMProviderHandler() *LLMProviderHandler {
	return &LLMProviderHandler{
		llmRepo: repository.NewLLMProviderRepository(),
	}
}

// List 获取大模型配置列表
func (h *LLMProviderHandler) List(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	providers, err := h.llmRepo.ListByUserID(userID)
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	// 计算分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(providers) {
		providers = []model.LLMProvider{}
	} else if end > len(providers) {
		providers = providers[start:]
	} else {
		providers = providers[start:end]
	}

	response.Success(ctx, map[string]interface{}{
		"items": providers,
		"total":  len(providers),
		"page":   page,
	})
}

// Get 获取单个大模型配置
func (h *LLMProviderHandler) Get(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ParamError(ctx, "无效的ID")
		return
	}

	provider, err := h.llmRepo.FindByID(uint(id))
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	// 检查权限
	if provider.UserID != userID {
		response.Error(ctx, errno.Forbidden)
		return
	}

	response.Success(ctx, provider)
}

// Create 创建大模型配置
func (h *LLMProviderHandler) Create(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	var req model.LLMProviderRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	provider := &model.LLMProvider{
		UserID:     userID,
		Name:      req.Name,
		Provider:  req.Provider,
		APIKey:    req.APIKey,
		BaseURL:   req.BaseURL,
		ModelName: req.ModelName,
		IsDefault: req.IsDefault,
		IsEnabled: req.IsEnabled,
		Timeout:   req.Timeout,
		RetryCount: req.RetryCount,
		Description: req.Description,
		SortOrder: req.SortOrder,
	}

	if err := h.llmRepo.Create(provider); err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, provider)
}

// Update 更新大模型配置
func (h *LLMProviderHandler) Update(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ParamError(ctx, "无效的ID")
		return
	}

	provider, err := h.llmRepo.FindByID(uint(id))
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	// 检查权限
	if provider.UserID != userID {
		response.Error(ctx, errno.Forbidden)
		return
	}

	var req model.LLMProviderRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	provider.Name = req.Name
	provider.Provider = req.Provider
	provider.APIKey = req.APIKey
	provider.BaseURL = req.BaseURL
	provider.ModelName = req.ModelName
	provider.IsDefault = req.IsDefault
	provider.IsEnabled = req.IsEnabled
	provider.Timeout = req.Timeout
	provider.RetryCount = req.RetryCount
	provider.Description = req.Description
	provider.SortOrder = req.SortOrder

	if err := h.llmRepo.Update(provider); err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, provider)
}

// Delete 删除大模型配置
func (h *LLMProviderHandler) Delete(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ParamError(ctx, "无效的ID")
		return
	}

	provider, err := h.llmRepo.FindByID(uint(id))
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	// 检查权限
	if provider.UserID != userID {
		response.Error(ctx, errno.Forbidden)
		return
	}

	if err := h.llmRepo.Delete(uint(id)); err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, nil)
}

// GetActive 获取用户当前激活的大模型配置
func (h *LLMProviderHandler) GetActive(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	provider, err := h.llmRepo.GetDefaultByUserID(userID)
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, provider)
}
