// Package handler 提供请求处理
package handler

import (
	"context"
	"errors"
	"strconv"
	"time"
	"xiaohongshu/internal/config"
	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/repository"
	"xiaohongshu/internal/service"
	"xiaohongshu/internal/utils"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

// UserHandler 用户处理器
type UserHandler struct {
	userRepo *repository.UserRepository
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userRepo: repository.NewUserRepository(),
	}
}

// toAuthUserInfo 转换为认证用户信息格式
func toAuthUserInfo(user *model.User) *model.AuthUserInfo {
	// 确定角色代码
	roleCode := "user"
	if user.Role != nil {
		roleCode = user.Role.Code
	}
	
	return &model.AuthUserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		IsActive:  user.Status == 1,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Role:      roleCode,
	}
}

// Register 用户注册
func (h *UserHandler) Register(c context.Context, ctx *app.RequestContext) {
	var req model.RegisterRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		response.ParamError(ctx, err.Error())
		return
	}

	userService := service.NewUserService()
	resp, err := userService.Register(c, &req)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	// 兼容前端格式
	authResp := map[string]interface{}{
		"access_token": resp.Token,
		"token_type":   "Bearer",
		"expires_in":   config.AppConfig.JWT.Expire * 3600,
		"user":         toAuthUserInfo(&resp.User),
	}

	response.Success(ctx, authResp)
}

// Login 用户登录
func (h *UserHandler) Login(c context.Context, ctx *app.RequestContext) {
	var req model.LoginRequest
	
	// 尝试从表单绑定（支持URLSearchParams格式）
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	
	if username != "" && password != "" {
		req.Username = username
		req.Password = password
	} else {
		// 从JSON绑定
		if err := ctx.BindAndValidate(&req); err != nil {
			response.ParamError(ctx, err.Error())
			return
		}
	}

	userService := service.NewUserService()
	resp, err := userService.Login(c, &req)
	if err != nil {
		if e, ok := err.(*errno.ErrNo); ok {
			response.Error(ctx, e)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	// 兼容前端格式
	authResp := map[string]interface{}{
		"access_token": resp.Token,
		"token_type":   "Bearer",
		"expires_in":   config.AppConfig.JWT.Expire * 3600,
		"user":         toAuthUserInfo(&resp.User),
	}

	response.Success(ctx, authResp)
}

// GetUserInfo 获取当前用户信息（兼容前端格式 /auth/me）
func (h *UserHandler) GetUserInfo(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, errno.UserNotFound)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, toAuthUserInfo(user))
}

// Logout 用户登出
func (h *UserHandler) Logout(c context.Context, ctx *app.RequestContext) {
	// 获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	// 获取 Token
	authHeader := string(ctx.GetHeader("Authorization"))
	token := ""
	if authHeader != "" {
		token = authHeader[len("Bearer "):]
	}

	// 解析 Token 获取过期时间
	claims, err := utils.ParseToken(token)
	if err != nil {
		// 即使 Token 解析失败，也返回登出成功（前端会清除本地 Token）
		response.Success(ctx, map[string]interface{}{
			"success": true,
			"message": "登出成功",
		})
		return
	}

	// 将 Token 添加到黑名单
	tokenBlacklistRepo := repository.NewTokenBlacklistRepository()
	err = tokenBlacklistRepo.AddToBlacklist(token, userID, claims.ExpiresAt.Time)
	if err != nil {
		// 即使添加黑名单失败，也返回登出成功
		response.Success(ctx, map[string]interface{}{
			"success": true,
			"message": "登出成功",
		})
		return
	}

	response.Success(ctx, map[string]interface{}{
		"success": true,
		"message": "登出成功",
	})
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(ctx, errno.Unauthorized)
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, errno.UserNotFound)
		} else {
			response.Error(ctx, errno.InternalError)
		}
		return
	}

	response.Success(ctx, user)
}

// ListUsers 获取用户列表（管理员）
func (h *UserHandler) ListUsers(c context.Context, ctx *app.RequestContext) {
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

	userService := service.NewUserService()
	users, total, err := userService.ListUsers(c, page, pageSize)
	if err != nil {
		response.Error(ctx, errno.InternalError)
		return
	}

	response.Success(ctx, map[string]interface{}{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
