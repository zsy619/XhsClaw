// Package service 提供业务逻辑层测试
package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewRendererService 测试创建渲染服务实例
func TestNewRendererService(t *testing.T) {
	service := NewRendererService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.styles)
	assert.Greater(t, len(service.styles), 0)
}

// TestGetStyles 测试获取所有样式
func TestGetStyles(t *testing.T) {
	service := NewRendererService()
	styles := service.GetStyles()

	assert.NotNil(t, styles)
	assert.Greater(t, len(styles), 0)

	// 验证默认样式是否存在
	foundDefault := false
	foundXiaohongshu := false
	for _, style := range styles {
		if style.Key == "default" {
			foundDefault = true
			assert.Equal(t, "简约灰", style.Name)
		}
		if style.Key == "xiaohongshu" {
			foundXiaohongshu = true
			assert.Equal(t, "小红书红", style.Name)
		}
	}
	assert.True(t, foundDefault, "应该找到默认样式")
	assert.True(t, foundXiaohongshu, "应该找到小红书红样式")
}

// TestRenderMarkdownToImage 测试渲染Markdown为图片
func TestRenderMarkdownToImage(t *testing.T) {
	service := NewRendererService()

	// 测试基本功能
	images, err := service.RenderMarkdownToImage(
		"# 测试内容\n这是测试内容",
		"default",
		"test",
		1080,
		1440,
		2160,
	)

	assert.NoError(t, err)
	assert.NotNil(t, images)
}

// TestGenerateCoverOnly 测试只生成封面
func TestGenerateCoverOnly(t *testing.T) {
	service := NewRendererService()

	// 测试基本功能
	imagePath, err := service.GenerateCoverOnly(
		"测试标题",
		"测试副标题",
		"default",
		"test_cover",
		1080,
		1440,
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, imagePath)
}

// TestStyleConfigProperties 测试样式配置属性
func TestStyleConfigProperties(t *testing.T) {
	service := NewRendererService()
	styles := service.GetStyles()

	for _, style := range styles {
		assert.NotEmpty(t, style.Key)
		assert.NotEmpty(t, style.Name)
		assert.NotNil(t, style.Primary)
		assert.NotNil(t, style.Secondary)
		assert.NotNil(t, style.Background)
		assert.NotNil(t, style.CardInner)
		assert.NotNil(t, style.TextPrimary)
		assert.NotNil(t, style.TextSecondary)
		assert.NotNil(t, style.Accent)
	}
}
