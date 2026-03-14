// Package middleware 提供中间件
package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// CORSMiddleware CORS中间件
func CORSMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Header("Access-Control-Max-Age", "86400")

		if string(ctx.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next(c)
	}
}
