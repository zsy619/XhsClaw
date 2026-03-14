// Package service 提供业务逻辑层
package service

import (
	"errors"
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

// SchedulePublish 定时发布
func (s *PublishService) SchedulePublish(userID uint, req *model.SchedulePublishRequest) (*model.PublishRecord, error) {
	content, err := s.contentRepo.FindByUserIDAndID(userID, req.ContentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ContentNotFound
		}
		return nil, errno.InternalError
	}

	publishTime, err := time.Parse(time.RFC3339, req.PublishTime)
	if err != nil {
		return nil, errno.InvalidParams.WithMessage("发布时间格式错误")
	}

	record := &model.PublishRecord{
		UserID:      userID,
		ContentID:   req.ContentID,
		Status:      0,
		ScheduledAt: publishTime,
	}

	if err := s.publishRepo.Create(record); err != nil {
		return nil, errno.InternalError
	}

	// 更新内容状态
	content.Status = 1
	content.PublishTime = &publishTime
	if err := s.contentRepo.Update(content); err != nil {
		return nil, errno.InternalError
	}

	return record, nil
}

// PublishNow 立即发布
func (s *PublishService) PublishNow(userID uint, req *model.PublishRequest) (*model.PublishRecord, error) {
	_, err := s.contentRepo.FindByUserIDAndID(userID, req.ContentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ContentNotFound
		}
		return nil, errno.InternalError
	}

	now := time.Now()
	record := &model.PublishRecord{
		UserID:      userID,
		ContentID:   req.ContentID,
		Status:      1,
		ScheduledAt: now,
	}

	if err := s.publishRepo.Create(record); err != nil {
		return nil, errno.InternalError
	}

	// 模拟发布（实际项目中需要调用小红书发布API）
	go s.doPublish(record.ID)

	return record, nil
}

// doPublish 执行发布操作
func (s *PublishService) doPublish(recordID uint) {
	record, err := s.publishRepo.FindByID(recordID)
	if err != nil {
		return
	}

	// 这里应该调用小红书的发布API
	// 模拟发布过程
	time.Sleep(2 * time.Second)

	now := time.Now()
	record.Status = 2
	record.PublishedAt = &now

	// 更新内容状态
	content, _ := s.contentRepo.FindByID(record.ContentID)
	if content != nil {
		content.Status = 2
		s.contentRepo.Update(content)
	}

	s.publishRepo.Update(record)
}

// GetPublishRecord 获取发布记录
func (s *PublishService) GetPublishRecord(userID, id uint) (*model.PublishRecord, error) {
	record, err := s.publishRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NotFound
		}
		return nil, errno.InternalError
	}

	if record.UserID != userID {
		return nil, errno.Forbidden
	}

	return record, nil
}

// ListPublishRecords 获取发布记录列表
func (s *PublishService) ListPublishRecords(userID uint, page, pageSize int) ([]model.PublishRecord, int64, error) {
	offset := (page - 1) * pageSize
	return s.publishRepo.ListByUserID(userID, offset, pageSize)
}

// CancelPublish 取消发布
func (s *PublishService) CancelPublish(userID, recordID uint) error {
	record, err := s.publishRepo.FindByID(recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.NotFound
		}
		return errno.InternalError
	}

	if record.UserID != userID {
		return errno.Forbidden
	}

	if record.Status != 0 {
		return errno.InvalidParams.WithMessage("只能取消待发布的记录")
	}

	record.Status = 3
	record.ErrorMsg = "用户取消发布"
	if err := s.publishRepo.Update(record); err != nil {
		return errno.InternalError
	}

	// 更新内容状态
	content, _ := s.contentRepo.FindByID(record.ContentID)
	if content != nil {
		content.Status = 0
		s.contentRepo.Update(content)
	}

	return nil
}

// RetryPublish 重试发布
func (s *PublishService) RetryPublish(userID, recordID uint) (*model.PublishRecord, error) {
	record, err := s.publishRepo.FindByID(recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NotFound
		}
		return nil, errno.InternalError
	}

	if record.UserID != userID {
		return nil, errno.Forbidden
	}

	if record.Status != 3 {
		return nil, errno.InvalidParams.WithMessage("只能重失败的记录")
	}

	// 重置记录状态
	now := time.Now()
	record.Status = 1
	record.ErrorMsg = ""
	record.UpdatedAt = now
	if err := s.publishRepo.Update(record); err != nil {
		return nil, errno.InternalError
	}

	// 更新内容状态
	content, _ := s.contentRepo.FindByID(record.ContentID)
	if content != nil {
		content.Status = 1
		s.contentRepo.Update(content)
	}

	// 重新执行发布
	go s.doPublish(record.ID)

	return record, nil
}
