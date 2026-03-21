// Package handler 提供处理器层
package handler

import (
	"context"
	"fmt"
	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/service"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
)

// TokenUsageHandler Token使用记录处理器
type TokenUsageHandler struct {
	tokenUsageService *service.TokenUsageService
}

// NewTokenUsageHandler 创建Token使用记录处理器
func NewTokenUsageHandler() *TokenUsageHandler {
	return &TokenUsageHandler{
		tokenUsageService: service.NewTokenUsageService(),
	}
}

// GetUserTokenUsage 获取用户Token使用记录
func (h *TokenUsageHandler) GetUserTokenUsage(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	// 从limit参数获取条数，默认50条
	limit := ctx.DefaultQuery("limit", "50")
	var limitInt int
	if _, err := fmt.Sscanf(limit, "%d", &limitInt); err != nil {
		limitInt = 50
	}

	usages, err := h.tokenUsageService.GetUserTokenUsage(c, userID, limitInt)
	if err != nil {
		response.ErrorWithMessage(ctx, errno.InternalError, err.Error())
		return
	}

	response.Success(ctx, usages)
}

// GetUserTokenStats 获取用户Token使用统计
func (h *TokenUsageHandler) GetUserTokenStats(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	stats, err := h.tokenUsageService.GetUserTokenStats(c, userID)
	if err != nil {
		response.ErrorWithMessage(ctx, errno.InternalError, err.Error())
		return
	}

	response.Success(ctx, stats)
}

// GetUserDailyStats 获取用户每日Token使用统计
func (h *TokenUsageHandler) GetUserDailyStats(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	days := ctx.DefaultQuery("days", "30")
	var daysInt int
	if _, err := fmt.Sscanf(days, "%d", &daysInt); err != nil {
		daysInt = 30
	}

	stats, err := h.tokenUsageService.GetUserDailyStats(c, userID, daysInt)
	if err != nil {
		response.ErrorWithMessage(ctx, errno.InternalError, err.Error())
		return
	}

	response.Success(ctx, stats)
}

// GetUserStatsByModel 获取用户按模型统计的使用情况
func (h *TokenUsageHandler) GetUserStatsByModel(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	stats, err := h.tokenUsageService.GetUserStatsByModel(c, userID)
	if err != nil {
		response.ErrorWithMessage(ctx, errno.InternalError, err.Error())
		return
	}

	response.Success(ctx, stats)
}

// GetGlobalTokenStats 获取全局Token使用统计（仅管理员）
func (h *TokenUsageHandler) GetGlobalTokenStats(c context.Context, ctx *app.RequestContext) {
	stats, err := h.tokenUsageService.GetGlobalTokenStats(c)
	if err != nil {
		response.ErrorWithMessage(ctx, errno.InternalError, err.Error())
		return
	}

	response.Success(ctx, stats)
}

// GetGlobalDailyStats 获取全局每日Token使用统计
func (h *TokenUsageHandler) GetGlobalDailyStats(c context.Context, ctx *app.RequestContext) {
	days := ctx.DefaultQuery("days", "30")
	var daysInt int
	if _, err := fmt.Sscanf(days, "%d", &daysInt); err != nil {
		daysInt = 30
	}

	stats, err := h.tokenUsageService.GetGlobalDailyStats(c, daysInt)
	if err != nil {
		response.ErrorWithMessage(ctx, errno.InternalError, err.Error())
		return
	}

	response.Success(ctx, stats)
}
