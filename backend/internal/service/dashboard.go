// Package service 提供业务逻辑层
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"xiaohongshu/internal/model"
	"xiaohongshu/internal/repository"
)

// DashboardService 仪表盘服务
type DashboardService struct {
	contentRepo        *repository.ContentRepository
	contentHistoryRepo *repository.ContentHistoryRepository
	tokenUsageRepo     *repository.TokenUsageRepository
}

// NewDashboardService 创建仪表盘服务实例
func NewDashboardService() *DashboardService {
	return &DashboardService{
		contentRepo:        repository.NewContentRepository(),
		contentHistoryRepo: repository.NewContentHistoryRepository(),
		tokenUsageRepo:     repository.NewTokenUsageRepository(),
	}
}

// GetUserDashboardStats 获取用户仪表盘统计数据
func (s *DashboardService) GetUserDashboardStats(ctx context.Context, userID uint) (*model.DashboardStats, error) {
	stats := &model.DashboardStats{}

	// 获取内容统计
	contentStats, err := s.getContentStats(ctx, userID)
	if err != nil {
		log.Printf("[DashboardService] GetUserDashboardStats error: failed to get content stats: %v", err)
		return nil, fmt.Errorf("获取内容统计失败: %w", err)
	}
	stats.TotalContents = contentStats.Total
	stats.PublishedCount = contentStats.Published
	stats.DraftCount = contentStats.Draft
	stats.PendingCount = contentStats.Pending
	stats.FailedCount = contentStats.Failed

	// 获取今日统计
	todayStats, err := s.getTodayStats(ctx, userID)
	if err != nil {
		log.Printf("[DashboardService] GetUserDashboardStats error: failed to get today stats: %v", err)
		return nil, fmt.Errorf("获取今日统计失败: %w", err)
	}
	stats.TodayContents = todayStats.Contents
	stats.TodayPublished = todayStats.Published
	stats.TodayViews = todayStats.Views

	// 获取每周趋势
	weeklyTrend, err := s.getWeeklyTrend(ctx, userID)
	if err != nil {
		log.Printf("[DashboardService] GetUserDashboardStats error: failed to get weekly trend: %v", err)
		return nil, fmt.Errorf("获取每周趋势失败: %w", err)
	}
	stats.WeeklyTrend = weeklyTrend

	// 获取发布状态分布
	statusDist, err := s.getStatusDistribution(ctx, userID)
	if err != nil {
		log.Printf("[DashboardService] GetUserDashboardStats error: failed to get status distribution: %v", err)
		return nil, fmt.Errorf("获取状态分布失败: %w", err)
	}
	stats.StatusDistribution = statusDist

	// 获取创作效率
	avgTime, successRate, err := s.getGenerationEfficiency(ctx, userID)
	if err != nil {
		log.Printf("[DashboardService] GetUserDashboardStats error: failed to get generation efficiency: %v", err)
		return nil, fmt.Errorf("获取创作效率失败: %w", err)
	}
	stats.AvgGenerationTime = avgTime
	stats.SuccessRate = successRate

	// 获取Token使用统计
	tokenStats, err := s.getTokenStats(ctx, userID)
	if err != nil {
		log.Printf("[DashboardService] GetUserDashboardStats error: failed to get token stats: %v", err)
		return nil, fmt.Errorf("获取Token统计失败: %w", err)
	}
	stats.TotalTokens = tokenStats.Total
	stats.TodayTokens = tokenStats.Today
	stats.TotalCost = tokenStats.Cost

	// 获取最近活动时间
	lastActivity, err := s.getLastActivityTime(ctx, userID)
	if err == nil && !lastActivity.IsZero() {
		stats.LastActivityTime = lastActivity.Format("2006-01-02 15:04:05")
	}

	stats.UpdatedAt = time.Now()

	log.Printf("[DashboardService] GetUserDashboardStats success: user %d, total contents %d", userID, stats.TotalContents)
	return stats, nil
}

// contentStatsResult 内容统计结果
type contentStatsResult struct {
	Total     int64
	Published int64
	Draft     int64
	Pending   int64
	Failed    int64
}

// getContentStats 获取内容统计
func (s *DashboardService) getContentStats(ctx context.Context, userID uint) (*contentStatsResult, error) {
	result := &contentStatsResult{}

	// 统计总数
	if err := repository.DB.Model(&model.Content{}).Where("user_id = ?", userID).Count(&result.Total).Error; err != nil {
		return nil, err
	}

	// 统计各状态数量
	if err := repository.DB.Model(&model.Content{}).Where("user_id = ? AND status = ?", userID, 0).Count(&result.Draft).Error; err != nil {
		return nil, err
	}
	if err := repository.DB.Model(&model.Content{}).Where("user_id = ? AND status = ?", userID, 1).Count(&result.Pending).Error; err != nil {
		return nil, err
	}
	if err := repository.DB.Model(&model.Content{}).Where("user_id = ? AND status = ?", userID, 2).Count(&result.Published).Error; err != nil {
		return nil, err
	}
	if err := repository.DB.Model(&model.Content{}).Where("user_id = ? AND status = ?", userID, 3).Count(&result.Failed).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// todayStatsResult 今日统计结果
type todayStatsResult struct {
	Contents  int64
	Published int64
	Views     int64
}

// getTodayStats 获取今日统计
func (s *DashboardService) getTodayStats(ctx context.Context, userID uint) (*todayStatsResult, error) {
	result := &todayStatsResult{}
	now := time.Now()
	today := now.Format("2006-01-02")

	// 今日生成的内容数
	if err := repository.DB.Model(&model.Content{}).
		Where("user_id = ? AND DATE(created_at) = ?", userID, today).
		Count(&result.Contents).Error; err != nil {
		return nil, err
	}

	// 今日发布的内容数
	if err := repository.DB.Model(&model.Content{}).
		Where("user_id = ? AND status = ? AND DATE(publish_time) = ?", userID, 2, today).
		Count(&result.Published).Error; err != nil {
		return nil, err
	}

	// 今日浏览数（暂无实现，暂时设为0）
	result.Views = 0

	return result, nil
}

// getWeeklyTrend 获取每周趋势
func (s *DashboardService) getWeeklyTrend(ctx context.Context, userID uint) ([]model.DailyStats, error) {
	stats := make([]model.DailyStats, 7)
	now := time.Now()

	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		dailyStat := model.DailyStats{
			Date: dateStr,
		}

		// 统计当日生成的内容数
		if err := repository.DB.Model(&model.Content{}).
			Where("user_id = ? AND DATE(created_at) = ?", userID, dateStr).
			Count(&dailyStat.Contents).Error; err != nil {
			log.Printf("[DashboardService] getWeeklyTrend error: failed to count contents for date %s: %v", dateStr, err)
		}

		// 统计当日发布的内容数
		if err := repository.DB.Model(&model.Content{}).
			Where("user_id = ? AND status = ? AND DATE(publish_time) = ?", userID, 2, dateStr).
			Count(&dailyStat.Published).Error; err != nil {
			log.Printf("[DashboardService] getWeeklyTrend error: failed to count published for date %s: %v", dateStr, err)
		}

		// 浏览数（暂无实现，暂时设为0）
		dailyStat.Views = 0

		// 获取当日Token使用
		var tokens int64
		if err := repository.DB.Model(&model.TokenUsage{}).
			Where("user_id = ? AND DATE(created_at) = ?", userID, dateStr).
			Select("COALESCE(SUM(input_tokens + output_tokens), 0)").
			Scan(&tokens).Error; err != nil {
			log.Printf("[DashboardService] getWeeklyTrend error: failed to get tokens for date %s: %v", dateStr, err)
		}
		dailyStat.Tokens = tokens

		stats[6-i] = dailyStat
	}

	return stats, nil
}

// getStatusDistribution 获取发布状态分布
func (s *DashboardService) getStatusDistribution(ctx context.Context, userID uint) ([]model.StatusCount, error) {
	statusLabels := map[int]string{
		0: "草稿",
		1: "待发布",
		2: "已发布",
		3: "发布失败",
	}

	statusColors := map[int]string{
		0: "#667eea",
		1: "#4facfe",
		2: "#43e97b",
		3: "#999999",
	}

	distribution := make([]model.StatusCount, 0, 4)

	for status := 0; status <= 3; status++ {
		var count int64
		if err := repository.DB.Model(&model.Content{}).
			Where("user_id = ? AND status = ?", userID, status).
			Count(&count).Error; err != nil {
			log.Printf("[DashboardService] getStatusDistribution error: failed to count status %d: %v", status, err)
			continue
		}

		distribution = append(distribution, model.StatusCount{
			Status: status,
			Label:  statusLabels[status],
			Count:  count,
			Color:  statusColors[status],
		})
	}

	return distribution, nil
}

// getGenerationEfficiency 获取创作效率
func (s *DashboardService) getGenerationEfficiency(ctx context.Context, userID uint) (avgTime float64, successRate float64, err error) {
	// 计算成功率
	var total, success int64
	if err := repository.DB.Model(&model.Content{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		log.Printf("[DashboardService] getGenerationEfficiency error: failed to count total: %v", err)
		return 5.2, 0, err
	}

	if err := repository.DB.Model(&model.Content{}).
		Where("user_id = ? AND status != ?", userID, 3).
		Count(&success).Error; err != nil {
		log.Printf("[DashboardService] getGenerationEfficiency error: failed to count success: %v", err)
		return 5.2, 0, err
	}

	if total > 0 {
		successRate = float64(success) / float64(total) * 100
	}

	// 平均生成时间（从TokenUsage表获取，暂无实现，设为模拟值）
	avgTime = 5.2

	return avgTime, successRate, nil
}

// tokenStatsResult Token统计结果
type tokenStatsResult struct {
	Total int64
	Today int64
	Cost  float64
}

// getTokenStats 获取Token统计
func (s *DashboardService) getTokenStats(ctx context.Context, userID uint) (*tokenStatsResult, error) {
	result := &tokenStatsResult{}
	now := time.Now()
	today := now.Format("2006-01-02")

	// 统计总Token数
	if err := repository.DB.Model(&model.TokenUsage{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(input_tokens + output_tokens), 0)").
		Scan(&result.Total).Error; err != nil {
		log.Printf("[DashboardService] getTokenStats error: failed to get total tokens: %v", err)
		return nil, err
	}

	// 统计今日Token数
	if err := repository.DB.Model(&model.TokenUsage{}).
		Where("user_id = ? AND DATE(created_at) = ?", userID, today).
		Select("COALESCE(SUM(input_tokens + output_tokens), 0)").
		Scan(&result.Today).Error; err != nil {
		log.Printf("[DashboardService] getTokenStats error: failed to get today tokens: %v", err)
		return nil, err
	}

	// 计算费用（假设每百万Token 1元）
	result.Cost = float64(result.Total) / 1000000 * 1.0

	return result, nil
}

// getLastActivityTime 获取最近活动时间
func (s *DashboardService) getLastActivityTime(ctx context.Context, userID uint) (time.Time, error) {
	var lastTime time.Time

	// 从内容表获取最近创建时间
	var contentTime time.Time
	err := repository.DB.Model(&model.Content{}).
		Where("user_id = ?", userID).
		Select("COALESCE(MAX(created_at), '1970-01-01')").
		Scan(&contentTime).Error
	if err == nil && !contentTime.IsZero() {
		lastTime = contentTime
	}

	// 从发布记录表获取最近发布时间
	var publishTime time.Time
	err = repository.DB.Model(&model.PublishRecord{}).
		Where("user_id = ?", userID).
		Select("COALESCE(MAX(created_at), '1970-01-01')").
		Scan(&publishTime).Error
	if err == nil && !publishTime.IsZero() && publishTime.After(lastTime) {
		lastTime = publishTime
	}

	return lastTime, nil
}

// GetUserActivities 获取用户最近活动
func (s *DashboardService) GetUserActivities(ctx context.Context, userID uint, limit int) ([]model.UserActivity, error) {
	if limit <= 0 {
		limit = 5 // 默认显示 5 个活动
	}

	activities := make([]model.UserActivity, 0)

	// 获取最近的内容创建记录
	var contents []model.Content
	if err := repository.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&contents).Error; err != nil {
		log.Printf("[DashboardService] GetUserActivities error: failed to get activities: %v", err)
		return nil, err
	}

	for _, content := range contents {
		var statusText string
		switch content.Status {
		case 0:
			statusText = "草稿"
		case 1:
			statusText = "待发布"
		case 2:
			statusText = "已发布"
		case 3:
			statusText = "发布失败"
		default:
			statusText = "未知"
		}

		activities = append(activities, model.UserActivity{
			ID:      content.ID,
			UserID:  content.UserID,
			Type:    "generate",
			Title:   content.Title,
			Status:  statusText,
			Time:    content.CreatedAt,
			TimeAgo: formatTimeAgo(content.CreatedAt),
		})
	}

	return activities, nil
}

// GetContentTrends 获取内容趋势
func (s *DashboardService) GetContentTrends(ctx context.Context, userID uint, days int) ([]model.ContentTrend, error) {
	if days <= 0 {
		days = 7
	}

	trends := make([]model.ContentTrend, days)
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		trend := model.ContentTrend{
			Date: dateStr,
		}

		// 统计当日生成的内容数
		if err := repository.DB.Model(&model.Content{}).
			Where("user_id = ? AND DATE(created_at) = ?", userID, dateStr).
			Count(&trend.Generate).Error; err != nil {
			log.Printf("[DashboardService] GetContentTrends error: failed to count generate: %v", err)
		}

		// 统计当日发布的内容数
		if err := repository.DB.Model(&model.Content{}).
			Where("user_id = ? AND status = ? AND DATE(publish_time) = ?", userID, 2, dateStr).
			Count(&trend.Publish).Error; err != nil {
			log.Printf("[DashboardService] GetContentTrends error: failed to count publish: %v", err)
		}

		trends[days-1-i] = trend
	}

	return trends, nil
}

// formatTimeAgo 格式化相对时间
func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "刚刚"
	} else if diff < time.Hour {
		return fmt.Sprintf("%d 分钟前", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%d 小时前", int(diff.Hours()))
	} else if diff < 7*24*time.Hour {
		return fmt.Sprintf("%d 天前", int(diff.Hours()/24))
	} else {
		return t.Format("01-02 15:04")
	}
}

// GetDashboardData 获取完整的仪表盘数据
func (s *DashboardService) GetDashboardData(ctx context.Context, userID uint) (*model.DashboardResponse, error) {
	stats, err := s.GetUserDashboardStats(ctx, userID)
	if err != nil {
		log.Printf("[DashboardService] GetDashboardData error: failed to get stats: %v", err)
		return nil, fmt.Errorf("获取仪表盘统计失败: %w", err)
	}

	activities, err := s.GetUserActivities(ctx, userID, 5) // 只显示前 5 个活动
	if err != nil {
		log.Printf("[DashboardService] GetDashboardData error: failed to get activities: %v", err)
		return nil, fmt.Errorf("获取用户活动失败: %w", err)
	}

	trends, err := s.GetContentTrends(ctx, userID, 7)
	if err != nil {
		log.Printf("[DashboardService] GetDashboardData error: failed to get trends: %v", err)
		return nil, fmt.Errorf("获取内容趋势失败: %w", err)
	}

	return &model.DashboardResponse{
		Stats:      *stats,
		Activities: activities,
		Trends:     trends,
	}, nil
}

// Content 属性解析辅助函数
func parseContentTags(content *model.Content) []string {
	var tags []string
	if content.Tags != "" {
		json.Unmarshal([]byte(content.Tags), &tags)
	}
	return tags
}

func parseContentImages(content *model.Content) []string {
	var images []string
	if content.Images != "" {
		json.Unmarshal([]byte(content.Images), &images)
	}
	return images
}
