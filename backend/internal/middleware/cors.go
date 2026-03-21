// Package middleware 提供中间件
package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// CORSMiddleware CORS中间件
// 处理跨域资源共享(CORS)的预检请求和响应头
func CORSMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 允许所有来源（生产环境应该配置具体的域名）
		ctx.Header("Access-Control-Allow-Origin", "*")

		// 允许所有常用的HTTP方法
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS")

		// 允许常用的请求头
		ctx.Header("Access-Control-Allow-Headers",
			"Origin, Content-Type, Accept, Authorization, X-Requested-With, Cache-Control")

		// 暴露响应头（让客户端可以访问这些响应头）
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		// 允许携带凭证（cookies、authorization headers等）
		ctx.Header("Access-Control-Allow-Credentials", "true")

		// 预检请求的缓存时间（秒）
		ctx.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求（OPTIONS方法）
		if string(ctx.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next(c)
	}
}
