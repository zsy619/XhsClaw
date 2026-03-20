// Package model 定义数据模型
package model

import (
	"time"
)

// SystemDict 系统字典表
type SystemDict struct {
	ID          uint      `json:"id" gorm:"primaryKey;comment:字典ID"`
	Category    string    `json:"category" gorm:"size:50;index;not null;comment:字典分类"`
	Code        string    `json:"code" gorm:"size:100;not null;comment:字典编码"`
	Name        string    `json:"name" gorm:"size:100;not null;comment:字典名称"`
	Value       string    `json:"value" gorm:"type:text;comment:字典值"`
	Description string    `json:"description" gorm:"size:255;comment:字典描述"`
	SortOrder   int       `json:"sort_order" gorm:"default:0;comment:排序"`
	Enabled     bool      `json:"enabled" gorm:"default:true;comment:是否启用"`
	Extra       string    `json:"extra" gorm:"type:text;comment:扩展信息JSON"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (SystemDict) TableName() string {
	return "system_dicts"
}

// DictCategory 常量定义
const (
	DictCategoryLLMProvider = "llm_provider" // 大模型服务商
	DictCategoryLLMModel    = "llm_model"    // 大模型
)

// SystemDictRequest 字典请求
type SystemDictRequest struct {
	Category    string `json:"category" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Value       string `json:"value"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Enabled     *bool  `json:"enabled"`
	Extra       string `json:"extra"`
}
