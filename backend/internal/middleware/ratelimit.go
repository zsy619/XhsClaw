// Package middleware 提供请求限流中间件
package middleware

import (
	"context"
	"sync"
	"time"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
)

// RateLimiter 请求限流器
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int           // 时间窗口内最大请求数
	window   time.Duration // 时间窗口
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// 定期清理过期的请求记录
	go rl.cleanup()

	return rl
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// 清理过期的请求记录
	var validRequests []time.Time
	for _, t := range rl.requests[key] {
		if now.Sub(t) < rl.window {
			validRequests = append(validRequests, t)
		}
	}
	rl.requests[key] = validRequests

	// 检查是否超过限制
	if len(validRequests) >= rl.limit {
		return false
	}

	// 记录新请求
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

// cleanup 定期清理过期的请求记录
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, requests := range rl.requests {
			var validRequests []time.Time
			for _, t := range requests {
				if now.Sub(t) < rl.window {
					validRequests = append(validRequests, t)
				}
			}
			if len(validRequests) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validRequests
			}
		}
		rl.mu.Unlock()
	}
}

// 全局限流器实例
var (
	apiLimiter     = NewRateLimiter(100, 1*time.Minute)    // API接口：每分钟100次
	loginLimiter   = NewRateLimiter(5, 5*time.Minute)      // 登录：5分钟5次
	generateLimiter = NewRateLimiter(20, 1*time.Minute)    // 生成内容：每分钟20次
)

// RateLimitMiddleware 通用限流中间件
func RateLimitMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 获取客户端IP
		clientIP := ctx.ClientIP()

		if !apiLimiter.Allow(clientIP) {
			response.Error(ctx, errno.NewErrNo(10007, "请求过于频繁，请稍后再试"))
			ctx.Abort()
			return
		}

		ctx.Next(c)
	}
}

// LoginRateLimitMiddleware 登录限流中间件
func LoginRateLimitMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		clientIP := ctx.ClientIP()

		if !loginLimiter.Allow(clientIP) {
			response.Error(ctx, errno.NewErrNo(10008, "登录尝试过于频繁，请5分钟后再试"))
			ctx.Abort()
			return
		}

		ctx.Next(c)
	}
}

// GenerateRateLimitMiddleware 生成内容限流中间件
func GenerateRateLimitMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		userID := GetUserID(c)
		key := string(rune(userID))
		if key == "" {
			key = ctx.ClientIP()
		}

		if !generateLimiter.Allow(key) {
			response.Error(ctx, errno.NewErrNo(10009, "生成内容过于频繁，请稍后再试"))
			ctx.Abort()
			return
		}

		ctx.Next(c)
	}
}
