// Package service 提供业务逻辑层
package service

import (
	"encoding/json"
	"errors"
	"time"
	"xiaohongshu/internal/model"
	"xiaohongshu/internal/repository"
	"xiaohongshu/pkg/errno"

	"gorm.io/gorm"
)

// ContentService 内容服务
type ContentService struct {
	contentRepo        *repository.ContentRepository
	contentHistoryRepo *repository.ContentHistoryRepository
	aiService          *AIService
}

// NewContentService 创建内容服务实例
func NewContentService() *ContentService {
	return &ContentService{
		contentRepo:        repository.NewContentRepository(),
		contentHistoryRepo: repository.NewContentHistoryRepository(),
		aiService:          NewAIService(),
	}
}

// GenerateContent 生成内容
func (s *ContentService) GenerateContent(userID uint, req *model.GenerateContentRequest) (*model.GenerateContentResponse, error) {
	items, err := s.aiService.GenerateXiaohongshuContent(req.SkillContent, req.Count, req.Length, "", "", "")
	if err != nil {
		return nil, err
	}

	modelItems := make([]model.ContentItem, len(items))
	for i, item := range items {
		modelItems[i] = model.ContentItem{
			Title:       item.Title,
			Description: item.Description,
			Tags:        item.Tags,
		}
	}

	return &model.GenerateContentResponse{
		Contents: modelItems,
	}, nil
}

// SaveContent 保存内容（支持备选标题）
func (s *ContentService) SaveContent(userID uint, req *model.ContentSaveRequest) (*model.Content, error) {
	tagsJSON, _ := json.Marshal(req.Tags)
	imagesJSON, _ := json.Marshal(req.Images)
	contentAttrsJSON, _ := json.Marshal(req.ContentAttributes)
	renderAttrsJSON, _ := json.Marshal(req.RenderAttributes)
	
	// 处理备选标题
	titleOptionsJSON := []byte("[]")
	if len(req.TitleOptions) > 0 {
		titleOptionsJSON, _ = json.Marshal(req.TitleOptions)
	}
	
	// 使用选中的标题作为主标题
	title := req.Title
	if len(req.TitleOptions) > 0 && req.SelectedTitleIndex >= 0 && req.SelectedTitleIndex < len(req.TitleOptions) {
		title = req.TitleOptions[req.SelectedTitleIndex]
	}
	
	content := &model.Content{
		UserID:            userID,
		Title:             title,
		TitleOptions:      string(titleOptionsJSON),
		SelectedTitleIndex: req.SelectedTitleIndex,
		Description:       req.Description,
		Tags:              string(tagsJSON),
		Images:            string(imagesJSON),
		ContentAttributes: string(contentAttrsJSON),
		RenderAttributes:  string(renderAttrsJSON),
		Status:            0,
	}

	if err := s.contentRepo.Create(content); err != nil {
		return nil, errno.InternalError
	}

	// 创建历史记录
	if err := s.contentHistoryRepo.CreateFromContent(content, "create", "创建新内容"); err != nil {
		// 历史记录创建失败不影响主流程，只记录日志
	}

	return content, nil
}

// GetContent 获取内容详情
func (s *ContentService) GetContent(userID, id uint) (*model.Content, error) {
	content, err := s.contentRepo.FindByUserIDAndID(userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ContentNotFound
		}
		return nil, errno.InternalError
	}
	return content, nil
}

// ListContents 获取内容列表
func (s *ContentService) ListContents(userID uint, page, pageSize int, status *int) ([]model.Content, int64, error) {
	offset := (page - 1) * pageSize
	return s.contentRepo.ListByUserID(userID, offset, pageSize, status)
}

// UpdateContent 更新内容
func (s *ContentService) UpdateContent(userID, id uint, req *model.UpdateContentRequest) (*model.Content, error) {
	content, err := s.contentRepo.FindByUserIDAndID(userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ContentNotFound
		}
		return nil, errno.InternalError
	}

	// 更新前先保存历史记录
	if err := s.contentHistoryRepo.CreateFromContent(content, "edit", "更新内容"); err != nil {
		// 历史记录创建失败不影响主流程
	}

	if req.Title != "" {
		content.Title = req.Title
	}
	if req.Description != "" {
		content.Description = req.Description
	}
	if req.Tags != nil {
		tagsJSON, _ := json.Marshal(req.Tags)
		content.Tags = string(tagsJSON)
	}
	if req.Images != nil {
		imagesJSON, _ := json.Marshal(req.Images)
		content.Images = string(imagesJSON)
	}
	if req.Status != nil {
		content.Status = *req.Status
	}
	if req.PublishTime != nil {
		t, err := time.Parse(time.RFC3339, *req.PublishTime)
		if err == nil {
			content.PublishTime = &t
		}
	}

	if err := s.contentRepo.Update(content); err != nil {
		return nil, errno.InternalError
	}

	return content, nil
}

// DeleteContent 删除内容
func (s *ContentService) DeleteContent(userID, id uint) error {
	content, err := s.contentRepo.FindByUserIDAndID(userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ContentNotFound
		}
		return errno.InternalError
	}

	// 删除前先保存历史记录
	if err := s.contentHistoryRepo.CreateFromContent(content, "delete", "删除内容"); err != nil {
		// 历史记录创建失败不影响主流程
	}

	if err := s.contentRepo.Delete(id); err != nil {
		return errno.InternalError
	}

	return nil
}

// ListContentHistories 获取历史记录列表
func (s *ContentService) ListContentHistories(userID uint, page, pageSize int, contentID *uint) ([]model.ContentHistory, int64, error) {
	offset := (page - 1) * pageSize
	return s.contentHistoryRepo.ListByUserID(userID, offset, pageSize, contentID)
}

// GetContentHistory 获取历史记录详情
func (s *ContentService) GetContentHistory(userID, id uint) (*model.ContentHistory, error) {
	history, err := s.contentHistoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ContentNotFound
		}
		return nil, errno.InternalError
	}
	
	// 验证权限
	if history.UserID != userID {
		return nil, errno.ContentNotFound
	}
	
	return history, nil
}

// RestoreContentHistory 恢复到历史版本
func (s *ContentService) RestoreContentHistory(userID, historyID uint) (*model.Content, error) {
	history, err := s.GetContentHistory(userID, historyID)
	if err != nil {
		return nil, err
	}
	
	// 查找内容
	content, err := s.contentRepo.FindByUserIDAndID(userID, history.ContentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果内容已被删除，重新创建
			content = &model.Content{
				UserID:            history.UserID,
				Title:             history.Title,
				TitleOptions:      history.TitleOptions,
				SelectedTitleIndex: history.SelectedTitleIndex,
				Description:       history.Description,
				Tags:              history.Tags,
				Images:            history.Images,
				ContentAttributes: history.ContentAttributes,
				RenderAttributes:  history.RenderAttributes,
				Status:            0,
			}
			
			if err := s.contentRepo.Create(content); err != nil {
				return nil, errno.InternalError
			}
			
			// 创建历史记录
			if err := s.contentHistoryRepo.CreateFromContent(content, "create", "从历史记录恢复"); err != nil {
				// 历史记录创建失败不影响主流程
			}
			
			return content, nil
		}
		return nil, errno.InternalError
	}
	
	// 更新前保存历史记录
	if err := s.contentHistoryRepo.CreateFromContent(content, "edit", "恢复到历史版本"); err != nil {
		// 历史记录创建失败不影响主流程
	}
	
	// 更新内容
	content.Title = history.Title
	content.TitleOptions = history.TitleOptions
	content.SelectedTitleIndex = history.SelectedTitleIndex
	content.Description = history.Description
	content.Tags = history.Tags
	content.Images = history.Images
	content.ContentAttributes = history.ContentAttributes
	content.RenderAttributes = history.RenderAttributes
	
	if err := s.contentRepo.Update(content); err != nil {
		return nil, errno.InternalError
	}
	
	return content, nil
}
