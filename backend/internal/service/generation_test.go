// Package service 提供业务逻辑层测试
package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"xiaohongshu/internal/model"
)

// TestNewGenerationService 测试创建生成服务实例
func TestNewGenerationService(t *testing.T) {
	service := NewGenerationService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.aiService)
	assert.NotNil(t, service.userConfigService)
}

// TestGenerateContent 测试生成内容
func TestGenerateContent(t *testing.T) {
	service := NewGenerationService()

	// 创建测试请求
	req := &model.GenerationRequest{
		Keywords:        "小红书运营",
		StylePreference: "informative",
		TargetAudience:  "新手博主",
		Length:          300,
	}

	// 测试生成内容（应该返回模拟内容）
	result, err := service.GenerateContent(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.GeneratedContent)
	assert.NotEmpty(t, result.GeneratedTitle)
	assert.NotEmpty(t, result.GeneratedTags)
	assert.Greater(t, len(result.GeneratedTags), 0)
}

// TestRewriteContent 测试改写内容
func TestRewriteContent(t *testing.T) {
	service := NewGenerationService()

	// 创建测试请求
	req := &model.RewriteRequest{
		Content:         "原有的小红书文案内容",
		StylePreference: "cute",
		PreserveKeyInfo: true,
		Length:          200,
	}

	// 测试改写内容（应该返回模拟内容）
	result, err := service.RewriteContent(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.GeneratedContent)
	assert.NotEmpty(t, result.GeneratedTitle)
	assert.NotEmpty(t, result.GeneratedTags)
}

// TestGetLengthDescription 测试获取长度描述
func TestGetLengthDescription(t *testing.T) {
	service := NewGenerationService()

	testCases := []struct {
		name     string
		length   int
		expected string
	}{
		{
			name:     "简短精炼",
			length:   50,
			expected: "简短精炼（约100字）",
		},
		{
			name:     "中等长度",
			length:   200,
			expected: "中等长度（约300字）",
		},
		{
			name:     "详细完整",
			length:   400,
			expected: "详细完整（约500字）",
		},
		{
			name:     "很长",
			length:   600,
			expected: "详细完整（约800字）",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.getLengthDescription(tc.length)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestGetStyleDescription 测试获取风格描述
func TestGetStyleDescription(t *testing.T) {
	service := NewGenerationService()

	testCases := []struct {
		name     string
		style    string
		expected string
	}{
		{
			name:     "活泼可爱",
			style:    "cute",
			expected: "活泼可爱",
		},
		{
			name:     "专业严谨",
			style:    "professional",
			expected: "专业严谨",
		},
		{
			name:     "文艺清新",
			style:    "artistic",
			expected: "文艺清新",
		},
		{
			name:     "幽默风趣",
			style:    "humorous",
			expected: "幽默风趣",
		},
		{
			name:     "干货分享",
			style:    "informative",
			expected: "干货分享",
		},
		{
			name:     "未知风格",
			style:    "unknown",
			expected: "unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.getStyleDescription(tc.style)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestGenerateMockContent 测试生成模拟内容
func TestGenerateMockContent(t *testing.T) {
	service := NewGenerationService()

	req := &model.GenerationRequest{
		Keywords:        "测试主题",
		StylePreference: "cute",
		TargetAudience:  "测试用户",
		Length:          300,
	}

	result, err := service.generateMockContent(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.GeneratedContent, "测试主题")
	assert.Contains(t, result.GeneratedTitle, "测试主题")
	assert.Greater(t, len(result.GeneratedTags), 0)
}

// TestRewriteMockContent 测试改写模拟内容
func TestRewriteMockContent(t *testing.T) {
	service := NewGenerationService()

	req := &model.RewriteRequest{
		Content:         "原文案内容",
		StylePreference: "cute",
		PreserveKeyInfo: true,
		Length:          200,
	}

	result, err := service.rewriteMockContent(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.GeneratedContent, "原文案内容")
	assert.NotEmpty(t, result.GeneratedTitle)
	assert.Greater(t, len(result.GeneratedTags), 0)
}
