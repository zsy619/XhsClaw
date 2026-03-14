// Package middleware 提供中间件
package middleware

import (
	"context"
	"strings"
	"xiaohongshu/internal/repository"
	"xiaohongshu/internal/utils"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	// UserIDKey 用户ID上下文键
	UserIDKey ContextKey = "user_id"
	// UsernameKey 用户名上下文键
	UsernameKey ContextKey = "username"
	// RoleKey 用户角色上下文键
	RoleKey ContextKey = "role"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		authHeader := string(ctx.GetHeader("Authorization"))
		if authHeader == "" {
			response.Error(ctx, errno.Unauthorized)
			ctx.Abort()
			return
		}

		// 去掉 "Bearer " 前缀
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			response.Error(ctx, errno.Unauthorized)
			ctx.Abort()
			return
		}

		// 检查 Token 是否在黑名单中
		tokenBlacklistRepo := repository.NewTokenBlacklistRepository()
		isBlacklisted, err := tokenBlacklistRepo.IsTokenBlacklisted(token)
		if err != nil {
			response.Error(ctx, errno.InternalError)
			ctx.Abort()
			return
		}
		if isBlacklisted {
			response.Error(ctx, errno.Unauthorized)
			ctx.Abort()
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			response.Error(ctx, errno.Unauthorized)
			ctx.Abort()
			return
		}

		// 将用户信息存入上下文
		c = context.WithValue(c, UserIDKey, claims.UserID)
		c = context.WithValue(c, UsernameKey, claims.Username)
		c = context.WithValue(c, RoleKey, claims.Role)

		ctx.Next(c)
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c context.Context) uint {
	if userID, ok := c.Value(UserIDKey).(uint); ok {
		return userID
	}
	return 0
}

// GetUsername 从上下文中获取用户名
func GetUsername(c context.Context) string {
	if username, ok := c.Value(UsernameKey).(string); ok {
		return username
	}
	return ""
}

// GetRole 从上下文中获取用户角色
func GetRole(c context.Context) string {
	if role, ok := c.Value(RoleKey).(string); ok {
		return role
	}
	return ""
}
