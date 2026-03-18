// Package app 应用程序核心
package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"xiaohongshu/internal/handler"
	"xiaohongshu/internal/middleware"
	"xiaohongshu/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// getProjectRoot 获取项目根目录的绝对路径
func getProjectRoot() string {
	return "/Volumes/E/JYW/创意项目/XhsClaw/backend"
}

// serveImage 自定义静态图片文件处理器
func serveImage(imagesDir string) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		filename := ctx.Param("filepath")
		logger := utils.GetLogger()

		logger.Info("=== serveImage 被调用 ===")
		logger.Info("请求 filepath 参数：%s", filename)
		logger.Info("请求 URL: %s", string(ctx.Request.URI().Path()))

		if filename == "" {
			logger.Warn("文件名为空")
			ctx.String(400, "文件名不能为空")
			return
		}

		logger.Info("请求图片：%s", filename)

		// 安全地构建文件路径
		imagePath := filepath.Join(imagesDir, filename)

		// 检查路径是否安全（防止路径遍历攻击）
		absImagesDir, err := filepath.Abs(imagesDir)
		if err != nil {
			logger.Error("获取绝对路径失败: %v", err)
			ctx.String(500, "服务器错误")
			return
		}

		absImagePath, err := filepath.Abs(imagePath)
		if err != nil {
			logger.Error("获取图片绝对路径失败: %v", err)
			ctx.String(500, "服务器错误")
			return
		}

		// 确保图片路径在 imagesDir 目录下
		if !strings.HasPrefix(absImagePath, absImagesDir) {
			logger.Warn("非法的图片路径请求：%s", filename)
			ctx.String(403, "访问被拒绝")
			return
		}

		// 检查文件是否存在
		if _, err := os.Stat(absImagePath); os.IsNotExist(err) {
			logger.Warn("图片文件不存在：%s", absImagePath)
			ctx.String(404, "图片不存在")
			return
		}

		logger.Info("提供图片：%s", absImagePath)

		// 读取文件内容
		fileContent, err := os.ReadFile(absImagePath)
		if err != nil {
			logger.Error("读取图片文件失败：%v", err)
			ctx.String(500, "读取文件失败")
			return
		}

		// 设置 Content-Type
		ctx.Header("Content-Type", "image/png")
		ctx.Header("Content-Length", fmt.Sprintf("%d", len(fileContent)))

		// 直接写入响应
		ctx.Write(fileContent)
	}
}

// SetupRouter 设置路由
func SetupRouter(h *server.Hertz) {
	// 获取项目根目录的绝对路径
	projectRoot := getProjectRoot()
	imagesDir := filepath.Join(projectRoot, "public", "images")

	logger := utils.GetLogger()
	logger.Info("静态文件目录: %s", imagesDir)

	// CORS中间件
	h.Use(middleware.CORSMiddleware())

	// 创建处理器实例
	userHandler := handler.NewUserHandler()
	userConfigHandler := handler.NewUserConfigHandler()
	contentHandler := handler.NewContentHandler()
	publishHandler := handler.NewPublishHandler()
	generationHandler := handler.NewGenerationHandler()
	rendererHandler := handler.NewRendererHandler()
	enhancedRendererHandler := handler.NewEnhancedRendererHandler()
	roleHandler := handler.NewRoleHandler()

	// 健康检查 - 不需要认证
	h.GET("/health", handler.HealthCheck)
	h.GET("/ready", handler.ReadyCheck)
	h.GET("/live", handler.LivenessCheck)

	// 公共路由
	api := h.Group("/api")
	{
		// 健康检查（API路径）
		api.GET("/health", handler.HealthCheck)
		api.GET("/ready", handler.ReadyCheck)
		api.GET("/live", handler.LivenessCheck)
		// 用户认证 - /api/auth
		auth := api.Group("/auth")
		{
			// 登录限流中间件
			auth.Use(middleware.LoginRateLimitMiddleware())
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// ======== 重要：/api/v1 路由组放在认证路由之前 ========
		v1 := api.Group("/v1")
		{
			// v1 版本小红书渲染器 - 不需要认证
			v1.GET("/xiaohongshu-renderer/styles", rendererHandler.GetRendererStyles)
			v1.POST("/xiaohongshu-renderer/render", rendererHandler.RenderMarkdown)
			v1.POST("/xiaohongshu-renderer/cover", rendererHandler.GenerateCover)
			// 增强版渲染器 - 支持 AI 生成和智能分页
			v1.POST("/xiaohongshu-renderer/render-with-ai", enhancedRendererHandler.RenderWithAI)

			// v1 用户认证
			v1Auth := v1.Group("/auth")
			{
				v1Auth.POST("/register", userHandler.Register)
				v1Auth.POST("/login", userHandler.Login)
			}

			// Token使用记录路由
			tokenUsageHandler := handler.NewTokenUsageHandler()
			
			// v1 需要认证的路由
			v1Authorized := v1.Group("")
			v1Authorized.Use(middleware.AuthMiddleware())
			{
				// Token使用记录
				v1Authorized.GET("/token-usage", tokenUsageHandler.GetUserTokenUsage)
				v1Authorized.GET("/token-usage/stats", tokenUsageHandler.GetUserTokenStats)
				v1Authorized.GET("/token-usage/daily", tokenUsageHandler.GetUserDailyStats)
				v1Authorized.GET("/token-usage/by-model", tokenUsageHandler.GetUserStatsByModel)
				
				// 全局Token使用统计（仅管理员）
				v1Authorized.GET("/admin/token-usage/global", tokenUsageHandler.GetGlobalTokenStats)
				v1Authorized.GET("/admin/token-usage/global/daily", tokenUsageHandler.GetGlobalDailyStats)

				v1Authorized.GET("/user/info", userHandler.GetUserInfo)
				v1Authorized.GET("/users", userHandler.ListUsers)

				// 角色和权限管理
				v1Authorized.GET("/roles", roleHandler.ListRoles)
				v1Authorized.GET("/roles/all", roleHandler.ListAllRoles)
				v1Authorized.GET("/roles/:id", roleHandler.GetRole)
				v1Authorized.POST("/roles", roleHandler.CreateRole)
				v1Authorized.PUT("/roles/:id", roleHandler.UpdateRole)
				v1Authorized.DELETE("/roles/:id", roleHandler.DeleteRole)
				v1Authorized.GET("/permissions", roleHandler.ListPermissions)
				v1Authorized.PUT("/users/:id/role", roleHandler.UpdateUserRole)
				v1Authorized.PUT("/users/:id/status", roleHandler.UpdateUserStatus)

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
					v1Content.GET("/histories/list", contentHandler.ListContentHistories)
					v1Content.GET("/histories/:id", contentHandler.GetContentHistory)
					v1Content.POST("/histories/:id/restore", contentHandler.RestoreContentHistory)
				}

				v1Publish := v1Authorized.Group("/publish")
				{
					v1Publish.POST("/schedule", publishHandler.SchedulePublish)
					v1Publish.POST("/now", publishHandler.PublishNow)
					v1Publish.GET("/list", publishHandler.ListPublishRecords)
					v1Publish.GET("/:id", publishHandler.GetPublishRecord)
					v1Publish.POST("/:id/cancel", publishHandler.CancelPublish)
					v1Publish.POST("/:id/retry", publishHandler.RetryPublish)
				}

				v1Generation := v1Authorized.Group("/generation")
				{
					v1Generation.POST("/theme", generationHandler.GenerateContent)
					v1Generation.POST("/rewrite", generationHandler.RewriteContent)
				}
			}
		}

		// 图片路由 - 在 v1 组外单独注册，避免通配符路由问题
		api.GET("/v1/xiaohongshu-renderer/image/*filepath", serveImage(imagesDir))

		// ======== 需要认证的路由 ========
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			authorized.GET("/auth/me", userHandler.GetUserInfo)
			authorized.POST("/auth/logout", userHandler.Logout)
			authorized.GET("/auth/profile", userHandler.GetProfile)
			authorized.GET("/user/info", userHandler.GetUserInfo)
			authorized.GET("/users", userHandler.ListUsers)

			// 角色和权限管理
			authorized.GET("/roles", roleHandler.ListRoles)
			authorized.GET("/roles/all", roleHandler.ListAllRoles)
			authorized.GET("/roles/:id", roleHandler.GetRole)
			authorized.POST("/roles", roleHandler.CreateRole)
			authorized.PUT("/roles/:id", roleHandler.UpdateRole)
			authorized.DELETE("/roles/:id", roleHandler.DeleteRole)
			authorized.GET("/permissions", roleHandler.ListPermissions)
			authorized.PUT("/users/:id/role", roleHandler.UpdateUserRole)
			authorized.PUT("/users/:id/status", roleHandler.UpdateUserStatus)

			// 用户配置
			config := authorized.Group("/user/config")
			{
				config.GET("", userConfigHandler.GetConfig)
				config.PUT("", userConfigHandler.UpdateConfig)
			}

			// 内容生成 - /api/generation
			generation := authorized.Group("/generation")
			{
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
				content.GET("/histories/list", contentHandler.ListContentHistories)
				content.GET("/histories/:id", contentHandler.GetContentHistory)
				content.POST("/histories/:id/restore", contentHandler.RestoreContentHistory)
			}

			// 发布管理 - /api/publish
			publish := authorized.Group("/publish")
			{
				publish.POST("/schedule", publishHandler.SchedulePublish)
				publish.POST("/now", publishHandler.PublishNow)
				publish.GET("/list", publishHandler.ListPublishRecords)
				publish.GET("/:id", publishHandler.GetPublishRecord)
				publish.POST("/:id/cancel", publishHandler.CancelPublish)
				publish.POST("/:id/retry", publishHandler.RetryPublish)
			}

			// 小红书渲染器 - /api/xiaohongshu-renderer
			renderer := authorized.Group("/xiaohongshu-renderer")
			{
				renderer.GET("/styles", rendererHandler.GetRendererStyles)
				renderer.POST("/render", rendererHandler.RenderMarkdown)
				renderer.POST("/cover", rendererHandler.GenerateCover)
			}
		}
	}
}
