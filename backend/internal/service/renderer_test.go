// Package service 提供业务逻辑层测试
package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewRendererService 测试创建渲染服务实例
func TestNewRendererService(t *testing.T) {
	service, err := NewRendererService()
	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.NotNil(t, service.themes)
	assert.Greater(t, len(service.themes), 0)
}

// TestGetStyles 测试获取所有样式
func TestGetStyles(t *testing.T) {
	service, err := NewRendererService()
	assert.NoError(t, err)
	styles := service.GetStyles()

	assert.NotNil(t, styles)
	assert.Greater(t, len(styles), 0)

	// 验证默认样式是否存在
	foundDefault := false
	foundTerminal := false
	for _, style := range styles {
		if style.Key == "default" {
			foundDefault = true
			assert.Equal(t, "简约灰", style.Name)
		}
		if style.Key == "terminal" {
			foundTerminal = true
			assert.Equal(t, "终端风格", style.Name)
		}
	}
	assert.True(t, foundDefault, "应该找到默认样式")
	assert.True(t, foundTerminal, "应该找到终端风格样式")
}

// TestGetStyle 测试获取指定样式
func TestGetStyle(t *testing.T) {
	service, err := NewRendererService()
	assert.NoError(t, err)

	// 测试获取存在的样式
	style := service.GetStyle("terminal")
	assert.Equal(t, "terminal", style.Key)
	assert.Equal(t, "终端风格", style.Name)

	// 测试获取不存在的样式（应返回默认样式）
	style = service.GetStyle("nonexistent")
	assert.Equal(t, "default", style.Key)
}

// TestThemeConfigProperties 测试主题配置属性
func TestThemeConfigProperties(t *testing.T) {
	service, err := NewRendererService()
	assert.NoError(t, err)
	styles := service.GetStyles()

	for _, style := range styles {
		assert.NotEmpty(t, style.Key)
		assert.NotEmpty(t, style.Name)
		assert.NotEmpty(t, style.CoverBg)
		assert.NotEmpty(t, style.CardBg)
		assert.NotEmpty(t, style.TitleGradient)
		assert.NotEmpty(t, style.AccentColor)
	}
}

// TestCalculateTitleSize 测试标题字号计算
func TestCalculateTitleSize(t *testing.T) {
	width := 1080

	// 测试不同长度的标题
	assert.Equal(t, int(float64(width)*0.14), calculateTitleSize("短标题", width))          // <=6字
	assert.Equal(t, int(float64(width)*0.12), calculateTitleSize("这是一个中等标题", width))     // 7-10字
	assert.Equal(t, int(float64(width)*0.09), calculateTitleSize("这是一个比较长的标题内容", width)) // 11-18字
	// 注意：以下测试需要确保字符数正确
	longTitle := "这是一个非常非常长的标题内容测试" // 11-18字
	assert.Equal(t, int(float64(width)*0.09), calculateTitleSize(longTitle, width))
	veryLongTitle := "这是一个超级超级超级超级超级长的标题内容测试" // 19-30字
	assert.Equal(t, int(float64(width)*0.07), calculateTitleSize(veryLongTitle, width))
	// 测试真正的超长标题 (>30字)
	veryVeryLongTitle := "这是一个超级超级超级超级超级超级超级超级超级超级超级超级长的标题内容测试"
	assert.Equal(t, int(float64(width)*0.055), calculateTitleSize(veryVeryLongTitle, width))
}

// TestIsTagLine 测试标签行判断
func TestIsTagLine(t *testing.T) {
	// 测试纯标签行
	assert.True(t, isTagLine("#标签1"))
	assert.True(t, isTagLine("#标签1 #标签2"))
	assert.True(t, isTagLine("#小红书 #内容创作"))

	// 测试非标签行
	assert.False(t, isTagLine("# 标题")) // 标题行（有空格）
	assert.False(t, isTagLine("普通文本")) // 普通文本
	assert.False(t, isTagLine(""))     // 空行
}

// TestSimpleMarkdownToHTML 测试Markdown转HTML
func TestSimpleMarkdownToHTML(t *testing.T) {
	service, err := NewRendererService()
	assert.NoError(t, err)

	// 测试标题转换
	html := service.simpleMarkdownToHTML("# 一级标题")
	assert.Contains(t, html, "<h1>一级标题</h1>")

	// 测试二级标题
	html = service.simpleMarkdownToHTML("## 二级标题")
	assert.Contains(t, html, "<h2>二级标题</h2>")

	// 测试粗体
	html = service.simpleMarkdownToHTML("**粗体文本**")
	assert.Contains(t, html, "<strong>粗体文本</strong>")

	// 测试斜体
	html = service.simpleMarkdownToHTML("*斜体文本*")
	assert.Contains(t, html, "<em>斜体文本</em>")

	// 测试行内代码
	html = service.simpleMarkdownToHTML("`code`")
	assert.Contains(t, html, "<code>code</code>")
}

// TestEscapeHTML 测试HTML转义
func TestEscapeHTML(t *testing.T) {
	service, err := NewRendererService()
	assert.NoError(t, err)

	// 测试特殊字符转义
	assert.Equal(t, "&lt;div&gt;", service.escapeHTML("<div>"))
	assert.Equal(t, "&amp;", service.escapeHTML("&"))
	assert.Equal(t, "&quot;", service.escapeHTML("\""))
	assert.Equal(t, "&#39;", service.escapeHTML("'"))
}
