// Package app 应用程序核心
// 提供路由配置和初始化功能
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

// GetProjectRoot 获取项目根目录的绝对路径
func GetProjectRoot() string {
	return "/Volumes/E/JYW/创意项目/XhsClaw/backend"
}

// GetImagesDir 获取图片目录路径
func GetImagesDir(projectRoot string) string {
	return filepath.Join(projectRoot, "public", "images")
}

// GetLogger 获取日志实例
func GetLogger() *utils.Logger {
	return utils.GetLogger()
}

// SetupRouter 配置所有路由
// 这是路由配置的入口函数，由 main.go 调用
func SetupRouter(h *server.Hertz) {
	logger := GetLogger()
	logger.Info("开始配置路由...")

	// 注册公开路由（无需认证）
	registerPublicRoutes(h)

	// 注册认证路由（登录、注册等）
	registerAuthRoutes(h)

	// 注册需要认证的路由
	setupAuthorizedRoutes(h)

	logger.Info("路由配置完成")
}

// registerPublicRoutes 注册公开路由（无需认证）
// 包含健康检查、渲染器等公开接口
func registerPublicRoutes(h *server.Hertz) {
	logger := GetLogger()
	logger.Info("注册公开路由")

	// 根路径健康检查
	h.GET("/health", handler.HealthCheck)
	h.GET("/ready", handler.ReadyCheck)
	h.GET("/live", handler.LivenessCheck)

	// API 前缀健康检查
	api := h.Group("/api")
	api.GET("/health", handler.HealthCheck)
	api.GET("/ready", handler.ReadyCheck)
	api.GET("/live", handler.LivenessCheck)

	// v1 版本公开路由 - 小红书渲染器
	v1 := h.Group("/api/v1")
	v1.GET("/xhsclaw/styles", handler.NewRendererHandler().GetRendererStyles)
	v1.POST("/xhsclaw/render", handler.NewRendererHandler().RenderMarkdown)
	v1.POST("/xhsclaw/render-with-ai", handler.NewEnhancedRendererHandler().RenderWithAI)

	// 图片访问路由
	projectRoot := GetProjectRoot()
	imagesDir := GetImagesDir(projectRoot)
	h.GET("/api/v1/xhsclaw/image/*filepath", ServeImageHandler(imagesDir))
}

// registerAuthRoutes 注册认证相关路由
// 包含登录、注册等无需认证的路由
func registerAuthRoutes(h *server.Hertz) {
	logger := GetLogger()
	logger.Info("注册认证路由")

	// 公开路由组
	auth := h.Group("/api/auth")
	auth.Use(middleware.LoginRateLimitMiddleware())
	{
		auth.POST("/register", handler.NewUserHandler().Register)
		auth.POST("/login", handler.NewUserHandler().Login)
	}

	// v1 版本认证路由
	v1Auth := h.Group("/api/v1/auth")
	v1Auth.Use(middleware.LoginRateLimitMiddleware())
	{
		v1Auth.POST("/register", handler.NewUserHandler().Register)
		v1Auth.POST("/login", handler.NewUserHandler().Login)
	}
}

// setupAuthorizedRoutes 配置需要认证的路由
// 将认证中间件应用到路由组
func setupAuthorizedRoutes(h *server.Hertz) {
	logger := GetLogger()
	logger.Info("配置需要认证的路由")

	// API v1 版本授权路由组
	v1Authorized := h.Group("/api/v1")
	v1Authorized.Use(middleware.AuthMiddleware())

	// 用户信息
	userHandler := handler.NewUserHandler()
	v1Authorized.GET("/user/info", userHandler.GetUserInfo)
	v1Authorized.GET("/users", userHandler.ListUsers)

	// 用户配置
	userConfigHandler := handler.NewUserConfigHandler()
	config := v1Authorized.Group("/user/config")
	{
		config.GET("", userConfigHandler.GetConfig)
		config.PUT("", userConfigHandler.UpdateConfig)
	}

	// 仪表盘统计
	dashboardHandler := handler.NewDashboardHandler()
	v1Authorized.GET("/dashboard/stats", dashboardHandler.GetDashboardStats)
	v1Authorized.GET("/dashboard/data", dashboardHandler.GetDashboardData)
	v1Authorized.GET("/dashboard/activities", dashboardHandler.GetUserActivities)
	v1Authorized.GET("/dashboard/trends", dashboardHandler.GetContentTrends)

	// 角色和权限管理
	roleHandler := handler.NewRoleHandler()
	v1Authorized.GET("/roles", roleHandler.ListRoles)
	v1Authorized.GET("/roles/all", roleHandler.ListAllRoles)
	v1Authorized.GET("/roles/:id", roleHandler.GetRole)
	v1Authorized.POST("/roles", roleHandler.CreateRole)
	v1Authorized.PUT("/roles/:id", roleHandler.UpdateRole)
	v1Authorized.DELETE("/roles/:id", roleHandler.DeleteRole)
	v1Authorized.GET("/permissions", roleHandler.ListPermissions)
	v1Authorized.PUT("/users/:id/role", roleHandler.UpdateUserRole)
	v1Authorized.PUT("/users/:id/status", roleHandler.UpdateUserStatus)

	// 内容管理
	contentHandler := handler.NewContentHandler()
	content := v1Authorized.Group("/content")
	{
		content.POST("/generate", contentHandler.GenerateContent)
		content.POST("/save", contentHandler.SaveContent)
		content.GET("/list", contentHandler.ListContents)
		content.GET("/:id", contentHandler.GetContent)
		content.PUT("/:id", contentHandler.UpdateContent)
		content.DELETE("/:id", contentHandler.DeleteContent)
		content.GET("/histories/list", contentHandler.ListContentHistories)
		content.GET("/histories/:id", contentHandler.GetContentHistory)
		content.POST("/histories/:id/restore", contentHandler.RestoreContentHistory)
	}

	// 发布管理
	publishHandler := handler.NewPublishHandler()
	publish := v1Authorized.Group("/publish")
	{
		publish.POST("/schedule", publishHandler.SchedulePublish)
		publish.POST("/now", publishHandler.PublishNow)
		publish.GET("/list", publishHandler.ListPublishRecords)
		publish.GET("/:id", publishHandler.GetPublishRecord)
		publish.POST("/:id/cancel", publishHandler.CancelPublish)
		publish.POST("/:id/retry", publishHandler.RetryPublish)
	}

	// 内容生成
	generationHandler := handler.NewGenerationHandler()
	rendererHandler := handler.NewRendererHandler()
	generation := v1Authorized.Group("/generation")
	generation.Use(middleware.GenerateRateLimitMiddleware())
	{
		generation.POST("/theme", generationHandler.GenerateContent)
		generation.POST("/rewrite", generationHandler.RewriteContent)
		generation.GET("/styles", rendererHandler.GetRendererStyles)
		generation.POST("/render", rendererHandler.RenderMarkdown)
		generation.POST("/cover", rendererHandler.GenerateCover)
	}

	// 大模型配置管理
	llmHandler := handler.NewLLMProviderHandler()
	llm := v1Authorized.Group("/llm")
	{
		llm.GET("/providers", llmHandler.List)
		llm.GET("/providers/:id", llmHandler.Get)
		llm.POST("/providers", llmHandler.Create)
		llm.PUT("/providers/:id", llmHandler.Update)
		llm.DELETE("/providers/:id", llmHandler.Delete)
		llm.GET("/active", llmHandler.GetActive)
	}

	// 系统字典
	dictHandler := handler.NewSystemDictHandler()
	dict := v1Authorized.Group("/dict")
	{
		dict.GET("/category/:category", dictHandler.GetByCategory)
		dict.GET("/list", dictHandler.List)
		dict.GET("/categories", dictHandler.GetCategories)
	}

	// 小红书配置管理
	xhsConfigHandler := handler.NewXHSConfigHandler()
	xhs := v1Authorized.Group("/xhs")
	{
		xhs.GET("/configs", xhsConfigHandler.List)
		xhs.GET("/configs/:id", xhsConfigHandler.Get)
		xhs.POST("/configs", xhsConfigHandler.Create)
		xhs.PUT("/configs/:id", xhsConfigHandler.Update)
		xhs.DELETE("/configs/:id", xhsConfigHandler.Delete)
		xhs.POST("/configs/:id/verify", xhsConfigHandler.Verify)
		xhs.GET("/active", xhsConfigHandler.GetActive)
	}

	// Token 使用统计
	tokenHandler := handler.NewTokenUsageHandler()
	v1Authorized.GET("/token-usage", tokenHandler.GetUserTokenUsage)
	v1Authorized.GET("/token-usage/stats", tokenHandler.GetUserTokenStats)
	v1Authorized.GET("/token-usage/daily", tokenHandler.GetUserDailyStats)
	v1Authorized.GET("/token-usage/by-model", tokenHandler.GetUserStatsByModel)
	v1Authorized.GET("/admin/token-usage/global", tokenHandler.GetGlobalTokenStats)
	v1Authorized.GET("/admin/token-usage/global/daily", tokenHandler.GetGlobalDailyStats)
}

// ServeImageHandler 创建图片服务处理器
// 用于安全地提供静态图片文件
func ServeImageHandler(imagesDir string) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		filename := ctx.Param("filepath")
		logger := GetLogger()

		logger.Debug("serveImage 被调用，filepath: %s", filename)

		if filename == "" {
			logger.Warn("文件名为空")
			ctx.String(400, "文件名不能为空")
			return
		}

		// 安全地构建文件路径
		imagePath := filepath.Join(imagesDir, strings.TrimPrefix(filename, "/"))

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

		logger.Debug("提供图片：%s", absImagePath)

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
