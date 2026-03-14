// Package app 应用程序核心
package app

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"xiaohongshu/internal/handler"
	"xiaohongshu/internal/middleware"
)

// SetupRouter 设置路由
func SetupRouter(h *server.Hertz) {
	// 安全响应头中间件
	h.Use(middleware.SecurityHeadersMiddleware())

	// 请求ID中间件
	h.Use(middleware.RequestIDMiddleware())

	// 通用限流中间件
	h.Use(middleware.RateLimitMiddleware())

	// CORS中间件
	h.Use(middleware.CORSMiddleware())

	// 创建处理器实例
	userHandler := handler.NewUserHandler()
	userConfigHandler := handler.NewUserConfigHandler()
	contentHandler := handler.NewContentHandler()
	publishHandler := handler.NewPublishHandler()
	generationHandler := handler.NewGenerationHandler()
	rendererHandler := handler.NewRendererHandler()

	// 公共路由
	api := h.Group("/api")
	{
		// 用户认证 - /api/auth
		auth := api.Group("/auth")
		{
			// 登录限流中间件
			auth.Use(middleware.LoginRateLimitMiddleware())
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// 公共静态资源 - 图片访问不需要认证（img标签无法携带token）
		api.Static("/xiaohongshu-renderer/image", "./public/images")

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			authorized.GET("/auth/me", userHandler.GetUserInfo)
			authorized.POST("/auth/logout", userHandler.Logout)
			authorized.GET("/auth/profile", userHandler.GetProfile)
			authorized.GET("/user/info", userHandler.GetUserInfo)
			authorized.GET("/users", userHandler.ListUsers)

			// 用户配置
			config := authorized.Group("/user/config")
			{
				config.GET("", userConfigHandler.GetConfig)
				config.PUT("", userConfigHandler.UpdateConfig)
			}

			// 内容生成 - /api/generation
			generation := authorized.Group("/generation")
			{
				// 生成内容限流中间件
				generation.Use(middleware.GenerateRateLimitMiddleware())
				generation.POST("/theme", generationHandler.GenerateContent)
				generation.POST("/rewrite", generationHandler.RewriteContent)
			}

			// 内容管理 - /api/content
			content := authorized.Group("/content")
			{
				content.POST("/generate", contentHandler.GenerateContent)
				content.POST("/save", generationHandler.SaveContent)
				content.GET("/list", contentHandler.ListContents)
				content.GET("/:id", contentHandler.GetContent)
				content.PUT("/:id", contentHandler.UpdateContent)
				content.DELETE("/:id", contentHandler.DeleteContent)
			}

			// 发布管理 - /api/publish
			publish := authorized.Group("/publish")
			{
				publish.POST("/schedule", publishHandler.SchedulePublish)
				publish.POST("/now", publishHandler.PublishNow)
				publish.GET("/list", publishHandler.ListPublishRecords)
				publish.GET("/:id", publishHandler.GetPublishRecord)
				publish.DELETE("/:id/cancel", publishHandler.CancelPublish)
			}

			// 小红书渲染器 - /api/xiaohongshu-renderer
			renderer := authorized.Group("/xiaohongshu-renderer")
			{
				renderer.GET("/styles", rendererHandler.GetRendererStyles)
				renderer.POST("/render", rendererHandler.RenderMarkdown)
				renderer.POST("/cover", rendererHandler.GenerateCover)
			}
		}

		// 旧版本兼容 - /api/v1
		v1 := api.Group("/v1")
		{
			// v1 版本公共静态资源
			v1.Static("/xiaohongshu-renderer/image", "./public/images")

			// 用户认证
			v1Auth := v1.Group("/auth")
			{
				v1Auth.POST("/register", userHandler.Register)
				v1Auth.POST("/login", userHandler.Login)
			}

			// 需要认证的路由
			v1Authorized := v1.Group("")
			v1Authorized.Use(middleware.AuthMiddleware())
			{
				v1Authorized.GET("/user/info", userHandler.GetUserInfo)
				v1Authorized.GET("/users", userHandler.ListUsers)

				// 用户配置
				v1Config := v1Authorized.Group("/user/config")
				{
					v1Config.GET("", userConfigHandler.GetConfig)
					v1Config.PUT("", userConfigHandler.UpdateConfig)
				}

				v1Content := v1Authorized.Group("/content")
				{
					v1Content.POST("/generate", contentHandler.GenerateContent)
					v1Content.POST("/save", contentHandler.SaveContent)
					v1Content.GET("/list", contentHandler.ListContents)
					v1Content.GET("/:id", contentHandler.GetContent)
					v1Content.PUT("/:id", contentHandler.UpdateContent)
					v1Content.DELETE("/:id", contentHandler.DeleteContent)
				}

				v1Publish := v1Authorized.Group("/publish")
				{
					v1Publish.POST("/schedule", publishHandler.SchedulePublish)
					v1Publish.POST("/now", publishHandler.PublishNow)
					v1Publish.GET("/list", publishHandler.ListPublishRecords)
					v1Publish.GET("/:id", publishHandler.GetPublishRecord)
					v1Publish.DELETE("/:id/cancel", publishHandler.CancelPublish)
				}

				// 内容生成 - /api/v1/generation
				v1Generation := v1Authorized.Group("/generation")
				{
					v1Generation.POST("/theme", generationHandler.GenerateContent)
					v1Generation.POST("/rewrite", generationHandler.RewriteContent)
				}

				// 小红书渲染器 - /api/v1/xiaohongshu-renderer
				v1Renderer := v1Authorized.Group("/xiaohongshu-renderer")
				{
					v1Renderer.GET("/styles", rendererHandler.GetRendererStyles)
					v1Renderer.POST("/render", rendererHandler.RenderMarkdown)
					v1Renderer.POST("/cover", rendererHandler.GenerateCover)
				}
			}
		}
	}
}
