// Package handler 提供请求处理
package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/service"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"
)

// UserConfigHandler 用户配置处理器
type UserConfigHandler struct {
	userConfigService *service.UserConfigService
}

// NewUserConfigHandler 创建用户配置处理器实例
func NewUserConfigHandler() *UserConfigHandler {
	return &UserConfigHandler{
		userConfigService: service.NewUserConfigService(),
	}
}

// GetConfig 获取用户配置
func (h *UserConfigHandler) GetConfig(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	config, err := h.userConfigService.GetUserConfig(userID)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, config)
}

// UpdateConfig 更新用户配置
func (h *UserConfigHandler) UpdateConfig(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	var req model.UserConfigRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	config, err := h.userConfigService.UpdateUserConfig(userID, &req)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, config)
}
