// Package handler 提供请求处理
package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/service"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"
)

// DashboardHandler 仪表盘处理器
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler 创建仪表盘处理器实例
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		dashboardService: service.NewDashboardService(),
	}
}

// GetDashboardStats 获取仪表盘统计数据
func (h *DashboardHandler) GetDashboardStats(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	stats, err := h.dashboardService.GetUserDashboardStats(c, userID)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, stats)
}

// GetDashboardData 获取完整仪表盘数据
func (h *DashboardHandler) GetDashboardData(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	data, err := h.dashboardService.GetDashboardData(c, userID)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, data)
}

// GetUserActivities 获取用户最近活动
func (h *DashboardHandler) GetUserActivities(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	activities, err := h.dashboardService.GetUserActivities(c, userID, 10)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, activities)
}

// GetContentTrends 获取内容趋势
func (h *DashboardHandler) GetContentTrends(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	trends, err := h.dashboardService.GetContentTrends(c, userID, 7)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, trends)
}
