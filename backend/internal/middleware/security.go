// Package middleware 提供安全相关中间件
package middleware

import (
	"context"
	"math/rand"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}

// SecurityHeadersMiddleware 安全响应头中间件
func SecurityHeadersMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// X-Content-Type-Options: 防止MIME类型嗅探
		ctx.Header("X-Content-Type-Options", "nosniff")

		// X-Frame-Options: 防止点击劫持
		ctx.Header("X-Frame-Options", "DENY")

		// X-XSS-Protection: 启用XSS过滤器
		ctx.Header("X-XSS-Protection", "1; mode=block")

		// Content-Security-Policy: 内容安全策略
		ctx.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https:; frame-ancestors 'none';")

		// Referrer-Policy: 控制Referer头
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Strict-Transport-Security: 强制使用HTTPS（生产环境）
		// ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Permissions-Policy: 权限策略
		ctx.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		ctx.Next(c)
	}
}

// RequestIDMiddleware 为每个请求添加唯一ID
func RequestIDMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从请求头获取或生成新的请求ID
		requestIDBytes := ctx.GetHeader("X-Request-ID")
		requestID := string(requestIDBytes)
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 设置响应头
		ctx.Header("X-Request-ID", requestID)

		// 存入上下文
		c = context.WithValue(c, "request_id", requestID)

		ctx.Next(c)
	}
}

// generateRequestID 生成简单的请求ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
