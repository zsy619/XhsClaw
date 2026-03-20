// Package repository 提供数据访问层
package repository

import (
	"xiaohongshu/internal/model"

	"gorm.io/gorm"
)

// ContentHistoryRepository 内容历史记录仓库
type ContentHistoryRepository struct {
	db *gorm.DB
}

// NewContentHistoryRepository 创建内容历史记录仓库实例
func NewContentHistoryRepository() *ContentHistoryRepository {
	return &ContentHistoryRepository{
		db: GetDB(),
	}
}

// Create 创建历史记录
func (r *ContentHistoryRepository) Create(history *model.ContentHistory) error {
	return r.db.Create(history).Error
}

// FindByID 根据ID查找历史记录
func (r *ContentHistoryRepository) FindByID(id uint) (*model.ContentHistory, error) {
	var history model.ContentHistory
	err := r.db.Where("id = ?", id).First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// ListByUserID 根据用户ID获取历史记录列表
func (r *ContentHistoryRepository) ListByUserID(userID uint, offset, pageSize int, contentID *uint) ([]model.ContentHistory, int64, error) {
	var histories []model.ContentHistory
	var total int64

	query := r.db.Model(&model.ContentHistory{}).Where("user_id = ?", userID)
	if contentID != nil {
		query = query.Where("content_id = ?", *contentID)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// CreateFromContent 从内容创建历史记录
func (r *ContentHistoryRepository) CreateFromContent(content *model.Content, historyType, changeReason string) error {
	history := &model.ContentHistory{
		ContentID:          content.ID,
		UserID:             content.UserID,
		Type:               historyType,
		Title:              content.Title,
		TitleOptions:       content.TitleOptions,
		SelectedTitleIndex: content.SelectedTitleIndex,
		Description:       content.Description,
		Tags:               content.Tags,
		Images:             content.Images,
		CoverSuggestion:    content.CoverSuggestion,
		ContentAttributes:  content.ContentAttributes,
		RenderAttributes:   content.RenderAttributes,
		ChangeReason:      changeReason,
	}
	return r.Create(history)
}
