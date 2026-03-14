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
	contentRepo *repository.ContentRepository
	aiService   *AIService
}

// NewContentService 创建内容服务实例
func NewContentService() *ContentService {
	return &ContentService{
		contentRepo: repository.NewContentRepository(),
		aiService:   NewAIService(),
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
		ContentAttributes: string(contentAttrsJSON),
		RenderAttributes:  string(renderAttrsJSON),
		Status:            0,
	}

	if err := s.contentRepo.Create(content); err != nil {
		return nil, errno.InternalError
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
	_, err := s.contentRepo.FindByUserIDAndID(userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ContentNotFound
		}
		return errno.InternalError
	}

	if err := s.contentRepo.Delete(id); err != nil {
		return errno.InternalError
	}

	return nil
}
