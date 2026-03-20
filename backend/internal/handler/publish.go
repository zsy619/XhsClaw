// Package handler 提供请求处理
package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/service"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"
)

// PublishHandler 发布处理器
type PublishHandler struct {
	publishService *service.PublishService
}

// NewPublishHandler 创建发布处理器实例
func NewPublishHandler() *PublishHandler {
	return &PublishHandler{
		publishService: service.NewPublishService(),
	}
}

// SchedulePublish 定时发布
func (h *PublishHandler) SchedulePublish(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	var req model.SchedulePublishRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	record, err := h.publishService.SchedulePublish(c, userID, &req)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, record)
}

// PublishNow 立即发布
func (h *PublishHandler) PublishNow(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	var req model.PublishRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	record, err := h.publishService.PublishNow(c, userID, &req)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, record)
}

// GetPublishRecord 获取发布记录
func (h *PublishHandler) GetPublishRecord(c context.Context, ctx *app.RequestContext) {
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

	record, err := h.publishService.GetPublishRecord(c, userID, uint(id))
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, record)
}

// ListPublishRecords 获取发布记录列表
func (h *PublishHandler) ListPublishRecords(c context.Context, ctx *app.RequestContext) {
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

	records, total, err := h.publishService.ListPublishRecords(c, userID, page, pageSize)
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, map[string]interface{}{
		"list":      records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CancelPublish 取消发布
func (h *PublishHandler) CancelPublish(c context.Context, ctx *app.RequestContext) {
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

	err = h.publishService.CancelPublish(c, userID, uint(id))
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, nil)
}

// RetryPublish 重试发布
func (h *PublishHandler) RetryPublish(c context.Context, ctx *app.RequestContext) {
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

	record, err := h.publishService.RetryPublish(c, userID, uint(id))
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, record)
}
