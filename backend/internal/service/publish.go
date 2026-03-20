// Package service 提供业务逻辑层
package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"xiaohongshu/internal/model"
	"xiaohongshu/internal/repository"
	"xiaohongshu/pkg/errno"

	"gorm.io/gorm"
)

// PublishService 发布服务
type PublishService struct {
	publishRepo *repository.PublishRepository
	contentRepo *repository.ContentRepository
}

// NewPublishService 创建发布服务实例
func NewPublishService() *PublishService {
	return &PublishService{
		publishRepo: repository.NewPublishRepository(),
		contentRepo: repository.NewContentRepository(),
	}
}

// SchedulePublish 定时发布内容
func (s *PublishService) SchedulePublish(ctx context.Context, userID uint, req *model.SchedulePublishRequest) (*model.PublishRecord, error) {
	// 验证内容存在且属于当前用户
	content, err := s.contentRepo.FindByUserIDAndID(userID, req.ContentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[PublishService] SchedulePublish failed: content %d not found for user %d", req.ContentID, userID)
			return nil, errno.ContentNotFound
		}
		log.Printf("[PublishService] SchedulePublish error: %v", err)
		return nil, fmt.Errorf("查找内容失败: %w", err)
	}

	// 解析发布时间
	publishTime, err := time.Parse(time.RFC3339, req.PublishTime)
	if err != nil {
		log.Printf("[PublishService] SchedulePublish failed: invalid time format %s for user %d", req.PublishTime, userID)
		return nil, errno.InvalidParams.WithMessage("发布时间格式错误，应使用 RFC3339 格式")
	}

	// 验证发布时间在当前时间之后
	if publishTime.Before(time.Now()) {
		return nil, errno.InvalidParams.WithMessage("发布时间不能早于当前时间")
	}

	// 创建发布记录
	record := &model.PublishRecord{
		UserID:      userID,
		ContentID:   req.ContentID,
		Status:      model.PublishStatusPending,
		ScheduledAt: publishTime,
	}

	if err := s.publishRepo.Create(record); err != nil {
		log.Printf("[PublishService] SchedulePublish failed to create record: %v", err)
		return nil, fmt.Errorf("创建发布记录失败: %w", err)
	}

	// 更新内容状态为待发布
	content.Status = model.ContentStatusPending
	content.PublishTime = &publishTime
	if err := s.contentRepo.Update(content); err != nil {
		log.Printf("[PublishService] SchedulePublish failed to update content status: %v", err)
		// 不影响主流程，仅记录日志
	}

	log.Printf("[PublishService] SchedulePublish success: record %d created for content %d by user %d", record.ID, req.ContentID, userID)
	return record, nil
}

// PublishNow 立即发布内容
func (s *PublishService) PublishNow(ctx context.Context, userID uint, req *model.PublishRequest) (*model.PublishRecord, error) {
	// 验证内容存在且属于当前用户
	content, err := s.contentRepo.FindByUserIDAndID(userID, req.ContentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[PublishService] PublishNow failed: content %d not found for user %d", req.ContentID, userID)
			return nil, errno.ContentNotFound
		}
		log.Printf("[PublishService] PublishNow error: %v", err)
		return nil, fmt.Errorf("查找内容失败: %w", err)
	}

	now := time.Now()
	record := &model.PublishRecord{
		UserID:      userID,
		ContentID:   req.ContentID,
		Status:      model.PublishStatusPublishing,
		ScheduledAt: now,
	}

	if err := s.publishRepo.Create(record); err != nil {
		log.Printf("[PublishService] PublishNow failed to create record: %v", err)
		return nil, fmt.Errorf("创建发布记录失败: %w", err)
	}

	// 启动异步发布任务
	go s.doPublishWithContext(record.ID, content.ID)

	log.Printf("[PublishService] PublishNow success: record %d created, publishing started for content %d", record.ID, req.ContentID)
	return record, nil
}

// doPublishWithContext 使用context执行发布操作，支持取消和超时
func (s *PublishService) doPublishWithContext(recordID, contentID uint) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 创建一个带超时的channel，用于监听发布完成信号
	done := make(chan error, 1)

	go func() {
		done <- s.doPublish(ctx, recordID, contentID)
	}()

	select {
	case <-ctx.Done():
		// 超时或取消
		err := ctx.Err()
		log.Printf("[PublishService] doPublish context done: %v for record %d", err, recordID)
		s.handlePublishFailure(recordID, contentID, fmt.Sprintf("发布超时或被取消: %v", err))
	case err := <-done:
		if err != nil {
			log.Printf("[PublishService] doPublish failed: %v for record %d", err, recordID)
			s.handlePublishFailure(recordID, contentID, err.Error())
		}
	}
}

// doPublish 执行实际的发布操作
func (s *PublishService) doPublish(ctx context.Context, recordID, contentID uint) error {
	// 检查context是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	record, err := s.publishRepo.FindByID(recordID)
	if err != nil {
		return fmt.Errorf("查找发布记录失败: %w", err)
	}

	// 检查记录状态，确保未被取消
	if record.Status == model.PublishStatusCancelled {
		return errors.New("发布已取消")
	}

	log.Printf("[PublishService] doPublish started: record %d, content %d", recordID, contentID)

	// 模拟小红书发布API调用（实际项目中需要调用真实API）
	// 添加重试机制
	maxRetries := 3
	var publishErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		log.Printf("[PublishService] doPublish attempt %d/%d for record %d", attempt, maxRetries, recordID)

		publishErr = s.callPublishAPI(ctx, contentID)
		if publishErr == nil {
			break
		}

		log.Printf("[PublishService] doPublish attempt %d failed: %v", attempt, publishErr)

		if attempt < maxRetries {
			// 指数退避重试
			backoff := time.Duration(attempt*attempt) * time.Second
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}
	}

	if publishErr != nil {
		return fmt.Errorf("小红书API发布失败(已重试%d次): %w", maxRetries, publishErr)
	}

	// 发布成功，更新状态
	now := time.Now()
	record.Status = model.PublishStatusSuccess
	record.PublishedAt = &now
	record.ErrorMsg = ""

	if err := s.publishRepo.Update(record); err != nil {
		log.Printf("[PublishService] doPublish failed to update record status: %v", err)
	}

	// 更新内容状态
	content, _ := s.contentRepo.FindByID(contentID)
	if content != nil {
		content.Status = model.ContentStatusPublished
		s.contentRepo.Update(content)
		log.Printf("[PublishService] doPublish success: content %d status updated to published", contentID)
	}

	log.Printf("[PublishService] doPublish success: record %d published at %v", recordID, now)
	return nil
}

// callPublishAPI 调用小红书发布API（模拟实现）
func (s *PublishService) callPublishAPI(ctx context.Context, contentID uint) error {
	// 检查context是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 模拟API调用延迟
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(2 * time.Second):
	}

	// 实际项目中，这里应该调用小红书的发布API
	// 例如：
	// resp, err := http.PostForm(ctx, "https://api.xiaohongshu.com/publish", data)
	// if err != nil {
	//     return err
	// }

	log.Printf("[PublishService] callPublishAPI success for content %d", contentID)
	return nil
}

// handlePublishFailure 处理发布失败
func (s *PublishService) handlePublishFailure(recordID, contentID uint, errMsg string) {
	record, err := s.publishRepo.FindByID(recordID)
	if err != nil {
		log.Printf("[PublishService] handlePublishFailure: failed to find record %d: %v", recordID, err)
		return
	}

	record.Status = model.PublishStatusFailed
	record.ErrorMsg = errMsg

	if err := s.publishRepo.Update(record); err != nil {
		log.Printf("[PublishService] handlePublishFailure: failed to update record %d: %v", recordID, err)
	}

	// 更新内容状态为待发布（可重新发布）
	content, _ := s.contentRepo.FindByID(contentID)
	if content != nil {
		content.Status = model.ContentStatusPending
		s.contentRepo.Update(content)
	}

	log.Printf("[PublishService] handlePublishFailure: record %d failed with error: %s", recordID, errMsg)
}

// GetPublishRecord 获取发布记录详情
func (s *PublishService) GetPublishRecord(ctx context.Context, userID, id uint) (*model.PublishRecord, error) {
	record, err := s.publishRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NotFound
		}
		return nil, fmt.Errorf("查找发布记录失败: %w", err)
	}

	// 验证用户权限
	if record.UserID != userID {
		log.Printf("[PublishService] GetPublishRecord forbidden: user %d accessing record %d owned by user %d", userID, id, record.UserID)
		return nil, errno.Forbidden
	}

	return record, nil
}

// ListPublishRecords 获取用户的发布记录列表
func (s *PublishService) ListPublishRecords(ctx context.Context, userID uint, page, pageSize int) ([]model.PublishRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.publishRepo.ListByUserID(userID, offset, pageSize)
}

// CancelPublish 取消待发布的记录
func (s *PublishService) CancelPublish(ctx context.Context, userID, recordID uint) error {
	record, err := s.publishRepo.FindByID(recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.NotFound
		}
		return fmt.Errorf("查找发布记录失败: %w", err)
	}

	// 验证用户权限
	if record.UserID != userID {
		log.Printf("[PublishService] CancelPublish forbidden: user %d cancelling record %d owned by user %d", userID, recordID, record.UserID)
		return errno.Forbidden
	}

	// 只能取消待发布的记录
	if record.Status != model.PublishStatusPending {
		return errno.InvalidParams.WithMessage("只能取消待发布的记录")
	}

	record.Status = model.PublishStatusCancelled
	record.ErrorMsg = "用户主动取消发布"

	if err := s.publishRepo.Update(record); err != nil {
		return fmt.Errorf("更新发布记录失败: %w", err)
	}

	// 重置内容状态为草稿
	content, _ := s.contentRepo.FindByID(record.ContentID)
	if content != nil {
		content.Status = model.ContentStatusDraft
		s.contentRepo.Update(content)
	}

	log.Printf("[PublishService] CancelPublish success: record %d cancelled by user %d", recordID, userID)
	return nil
}

// RetryPublish 重试失败的发布
func (s *PublishService) RetryPublish(ctx context.Context, userID, recordID uint) (*model.PublishRecord, error) {
	record, err := s.publishRepo.FindByID(recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NotFound
		}
		return nil, fmt.Errorf("查找发布记录失败: %w", err)
	}

	// 验证用户权限
	if record.UserID != userID {
		log.Printf("[PublishService] RetryPublish forbidden: user %d retrying record %d owned by user %d", userID, recordID, record.UserID)
		return nil, errno.Forbidden
	}

	// 只能重试失败的记录
	if record.Status != model.PublishStatusFailed {
		return nil, errno.InvalidParams.WithMessage("只能重试失败的记录")
	}

	// 重置记录状态
	now := time.Now()
	record.Status = model.PublishStatusPublishing
	record.ErrorMsg = ""
	record.UpdatedAt = now

	if err := s.publishRepo.Update(record); err != nil {
		return nil, fmt.Errorf("更新发布记录失败: %w", err)
	}

	// 更新内容状态为发布中
	content, _ := s.contentRepo.FindByID(record.ContentID)
	if content != nil {
		content.Status = model.ContentStatusPublishing
		s.contentRepo.Update(content)
	}

	// 重新执行发布
	go s.doPublishWithContext(record.ID, record.ContentID)

	log.Printf("[PublishService] RetryPublish success: record %d retrying for user %d", recordID, userID)
	return record, nil
}
