// Package handler 提供请求处理
package handler

import (
	"context"
	"strconv"
	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/service"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
)

// ContentHandler 内容处理器
type ContentHandler struct {
	contentService *service.ContentService
}

// NewContentHandler 创建内容处理器实例
func NewContentHandler() *ContentHandler {
	return &ContentHandler{
		contentService: service.NewContentService(),
	}
}

// GenerateContent 生成内容
func (h *ContentHandler) GenerateContent(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	var req model.GenerateContentRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	if req.Count <= 0 {
		req.Count = 1
	}
	if req.Count > 10 {
		req.Count = 10
	}

	resp, err := h.contentService.GenerateContent(userID, &req)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, resp)
}

// SaveContent 保存内容
func (h *ContentHandler) SaveContent(c context.Context, ctx *app.RequestContext) {
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

	content, err := h.contentService.SaveContent(userID, &req)
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

// GetContent 获取内容详情
func (h *ContentHandler) GetContent(c context.Context, ctx *app.RequestContext) {
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

	content, err := h.contentService.GetContent(userID, uint(id))
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

// ListContents 获取内容列表
func (h *ContentHandler) ListContents(c context.Context, ctx *app.RequestContext) {
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
	
	var status *int
	if statusStr := ctx.Query("status"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			status = &s
		}
	}

	contents, total, err := h.contentService.ListContents(userID, page, pageSize, status)
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, map[string]interface{}{
		"list":      contents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateContent 更新内容
func (h *ContentHandler) UpdateContent(c context.Context, ctx *app.RequestContext) {
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

	var req model.UpdateContentRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	content, err := h.contentService.UpdateContent(userID, uint(id), &req)
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

// DeleteContent 删除内容
func (h *ContentHandler) DeleteContent(c context.Context, ctx *app.RequestContext) {
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

	err = h.contentService.DeleteContent(userID, uint(id))
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

// ListContentHistories 获取历史记录列表
func (h *ContentHandler) ListContentHistories(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	contentIDStr := ctx.Query("content_id")
	
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
	
	var contentID *uint
	if contentIDStr != "" {
		if cid, err := strconv.ParseUint(contentIDStr, 10, 32); err == nil {
			cidUint := uint(cid)
			contentID = &cidUint
		}
	}

	histories, total, err := h.contentService.ListContentHistories(userID, page, pageSize, contentID)
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, map[string]interface{}{
		"list":      histories,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetContentHistory 获取历史记录详情
func (h *ContentHandler) GetContentHistory(c context.Context, ctx *app.RequestContext) {
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

	history, err := h.contentService.GetContentHistory(userID, uint(id))
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, history)
}

// RestoreContentHistory 恢复到历史版本
func (h *ContentHandler) RestoreContentHistory(c context.Context, ctx *app.RequestContext) {
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

	content, err := h.contentService.RestoreContentHistory(userID, uint(id))
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
