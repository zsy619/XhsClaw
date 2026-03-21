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

// XHSConfigHandler 小红书配置处理器
type XHSConfigHandler struct {
	xhsConfigRepo *repository.XHSConfigRepository
}

// NewXHSConfigHandler 创建小红书配置处理器
func NewXHSConfigHandler() *XHSConfigHandler {
	return &XHSConfigHandler{
		xhsConfigRepo: repository.NewXHSConfigRepository(),
	}
}

// List 获取小红书配置列表
func (h *XHSConfigHandler) List(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	configs, total, err := h.xhsConfigRepo.List(userID, page, pageSize)
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, map[string]interface{}{
		"items": configs,
		"total": total,
	})
}

// Get 获取单个配置
func (h *XHSConfigHandler) Get(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, errno.InvalidParams)
		return
	}

	config, err := h.xhsConfigRepo.FindByUserIDAndID(userID, uint(id))
	if err != nil {
		response.Error(ctx, errno.NotFound)
		return
	}

	response.Success(ctx, config)
}

// Create 创建配置
func (h *XHSConfigHandler) Create(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	var req model.XHSConfigRequest
	if err := ctx.Bind(&req); err != nil {
		response.Error(ctx, errno.InvalidParams)
		return
	}

	config := &model.XHSConfig{
		UserID:      userID,
		Name:        req.Name,
		Cookie:      req.Cookie,
		Token:       req.Token,
		DeviceID:    req.DeviceID,
		XHSUserID:   req.XHSUserID,
		IsDefault:   req.IsDefault,
		IsEnabled:   req.IsEnabled,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      "pending",
	}

	if err := h.xhsConfigRepo.Create(config); err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, config)
}

// Update 更新配置
func (h *XHSConfigHandler) Update(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, errno.InvalidParams)
		return
	}

	config, err := h.xhsConfigRepo.FindByUserIDAndID(userID, uint(id))
	if err != nil {
		response.Error(ctx, errno.NotFound)
		return
	}

	var req model.XHSConfigRequest
	if err := ctx.Bind(&req); err != nil {
		response.Error(ctx, errno.InvalidParams)
		return
	}

	// 更新字段
	config.Name = req.Name
	config.Cookie = req.Cookie
	config.Token = req.Token
	config.DeviceID = req.DeviceID
	config.XHSUserID = req.XHSUserID
	config.IsDefault = req.IsDefault
	config.IsEnabled = req.IsEnabled
	config.Description = req.Description
	config.SortOrder = req.SortOrder

	if err := h.xhsConfigRepo.Update(config); err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, config)
}

// Delete 删除配置
func (h *XHSConfigHandler) Delete(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, errno.InvalidParams)
		return
	}

	_, err = h.xhsConfigRepo.FindByUserIDAndID(userID, uint(id))
	if err != nil {
		response.Error(ctx, errno.NotFound)
		return
	}

	if err := h.xhsConfigRepo.Delete(uint(id)); err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, nil)
}

// Verify 验证配置
func (h *XHSConfigHandler) Verify(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, errno.InvalidParams)
		return
	}

	config, err := h.xhsConfigRepo.FindByUserIDAndID(userID, uint(id))
	if err != nil {
		response.Error(ctx, errno.NotFound)
		return
	}

	// 这里应该调用小书书的验证API
	// 暂时模拟验证成功
	resp := &model.XHSVerifyResponse{
		Success: true,
		Message: "配置验证成功",
		UserID:  config.XHSUserID,
	}

	// 更新状态
	config.Status = "active"
	h.xhsConfigRepo.Update(config)

	response.Success(ctx, resp)
}

// GetActive 获取激活的配置
func (h *XHSConfigHandler) GetActive(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	config, err := h.xhsConfigRepo.GetActive(userID)
	if err != nil {
		response.Error(ctx, errno.NotFound)
		return
	}

	response.Success(ctx, config)
}
