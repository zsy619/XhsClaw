// Package model 定义数据模型
package model

import (
	"time"
)

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	// 内容统计
	TotalContents    int64 `json:"total_contents"`    // 总内容数
	PublishedCount   int64 `json:"published_count"`   // 已发布数
	DraftCount       int64 `json:"draft_count"`       // 草稿数
	PendingCount     int64 `json:"pending_count"`     // 待发布数
	FailedCount      int64 `json:"failed_count"`      // 失败数

	// 今日统计
	TodayContents    int64  `json:"today_contents"`    // 今日生成
	TodayPublished    int64  `json:"today_published"`   // 今日发布
	TodayViews        int64  `json:"today_views"`       // 今日浏览

	// 趋势数据（最近7天）
	WeeklyTrend []DailyStats `json:"weekly_trend"`

	// 发布状态分布
	StatusDistribution []StatusCount `json:"status_distribution"`

	// 创作效率
	AvgGenerationTime float64 `json:"avg_generation_time"` // 平均生成时间(秒)
	SuccessRate       float64 `json:"success_rate"`         // 成功率(%)

	// Token使用
	TotalTokens    int64 `json:"total_tokens"`    // 总Token数
	TodayTokens    int64 `json:"today_tokens"`    // 今日Token数
	TotalCost      float64 `json:"total_cost"`    // 总费用(元)

	// 最近活动时间
	LastActivityTime string `json:"last_activity_time"`

	UpdatedAt time.Time `json:"updated_at"`
}

// DailyStats 每日统计数据
type DailyStats struct {
	Date         string `json:"date"`          // 日期 YYYY-MM-DD
	Contents     int64  `json:"contents"`      // 生成内容数
	Published    int64  `json:"published"`     // 发布数
	Views        int64  `json:"views"`          // 浏览数
	Tokens       int64  `json:"tokens"`         // Token数
}

// StatusCount 状态数量
type StatusCount struct {
	Status int    `json:"status"` // 状态值: 0草稿, 1待发布, 2已发布, 3失败
	Label  string `json:"label"`  // 状态标签
	Count  int64  `json:"count"`  // 数量
	Color  string `json:"color"`  // 颜色
}

// UserActivity 用户活动
type UserActivity struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Type      string    `json:"type"`      // generate, save, publish, edit
	Title     string    `json:"title"`     // 内容标题
	Status    string    `json:"status"`    // 状态
	Time      time.Time `json:"time"`      // 时间
	TimeAgo   string    `json:"time_ago"`  // 相对时间
}

// ContentTrend 内容趋势
type ContentTrend struct {
	Date     string `json:"date"`
	Generate int64  `json:"generate"` // 生成数
	Publish  int64  `json:"publish"`  // 发布数
}

// DashboardResponse 仪表盘响应
type DashboardResponse struct {
	Stats     DashboardStats    `json:"stats"`
	Activities []UserActivity   `json:"activities"`
	Trends    []ContentTrend   `json:"trends"`
}