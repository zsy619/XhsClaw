// Package handler 提供HTTP请求处理
package handler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"xiaohongshu/internal/repository"
	"xiaohongshu/pkg/errno"
	"xiaohongshu/pkg/response"
)

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// HealthCheck 健康检查接口
func HealthCheck(c context.Context, ctx *app.RequestContext) {
	healthStatus := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  make(map[string]string),
	}

	// 检查数据库连接
	db := repository.GetDB()
	sqlDB, err := db.DB()
	if err != nil {
		healthStatus.Status = "unhealthy"
		healthStatus.Services["database"] = "error: " + err.Error()
	} else {
		// 测试数据库连接
		pingCtx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(pingCtx); err != nil {
			healthStatus.Status = "unhealthy"
			healthStatus.Services["database"] = "unreachable: " + err.Error()
		} else {
			healthStatus.Services["database"] = "healthy"
			// 获取数据库连接池统计信息
			stats := sqlDB.Stats()
			healthStatus.Services["database_open_connections"] = fmt.Sprintf("%d", stats.OpenConnections)
			healthStatus.Services["database_in_use"] = fmt.Sprintf("%d", stats.InUse)
			healthStatus.Services["database_idle"] = fmt.Sprintf("%d", stats.Idle)
		}
	}

	response.Success(ctx, healthStatus)
}

// ReadyCheck 就绪检查接口
func ReadyCheck(c context.Context, ctx *app.RequestContext) {
	// 检查数据库是否可用
	db := repository.GetDB()
	sqlDB, err := db.DB()
	if err != nil {
		response.ErrorWithMessage(ctx, errno.ServiceUnavailable, "database unavailable")
		return
	}

	// 测试数据库连接
	ctxTimeout, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctxTimeout); err != nil {
		response.ErrorWithMessage(ctx, errno.ServiceUnavailable, "database not ready")
		return
	}

	response.Success(ctx, map[string]string{
		"status": "ready",
	})
}

// LivenessCheck 存活检查接口
func LivenessCheck(c context.Context, ctx *app.RequestContext) {
	response.Success(ctx, map[string]string{
		"status": "alive",
	})
}

// GetDBStats 获取数据库统计信息
func GetDBStats() sql.DBStats {
	db := repository.GetDB()
	sqlDB, _ := db.DB()
	return sqlDB.Stats()
}
