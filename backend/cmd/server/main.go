// Package main 程序入口
package main

import (
	"fmt"
	"xiaohongshu/internal/app"
	"xiaohongshu/internal/config"
	"xiaohongshu/internal/repository"
	"xiaohongshu/internal/utils"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// 初始化日志系统
	logger := utils.GetLogger()
	defer logger.Close()

	logger.Info("========================================")
	logger.Info("小红书内容生成系统启动中...")
	logger.Info("========================================")

	// 加载配置
	logger.Info("正在加载配置文件...")
	cfg, err := config.LoadConfig("")
	if err != nil {
		logger.Fatal("Failed to load config: %v", err)
	}
	logger.Info("配置文件加载成功")

	// 初始化数据库
	logger.Info("正在初始化数据库...")
	if err := repository.InitDatabase(&cfg.Database); err != nil {
		logger.Fatal("Failed to init database: %v", err)
	}
	logger.Info("数据库初始化成功")

	// 创建Hertz服务器
	logger.Info("正在创建HTTP服务器...")
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%d", cfg.Server.Port)),
	)

	// 设置路由
	logger.Info("正在配置路由...")
	app.SetupRouter(h)

	logger.Info("========================================")
	logger.Info("服务器启动成功，监听端口: %d", cfg.Server.Port)
	logger.Info("========================================")

	h.Spin()
}
