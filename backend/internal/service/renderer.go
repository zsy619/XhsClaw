// Package service 提供业务逻辑层
package service

import (
	"image/color"
	"path/filepath"
)

// 分页模式常量
const (
	PaginationAutoSplit = "auto-split"
	PaginationSeparator = "separator"
	PaginationAutoFit   = "auto-fit"
	PaginationDynamic   = "dynamic"
)

// StyleConfig 样式配置
type StyleConfig struct {
	Key           string
	Name          string
	Primary       color.Color
	Secondary     color.Color
	Background    color.Color
	CardInner     color.Color
	TextPrimary   color.Color
	TextSecondary color.Color
	Accent        color.Color
	FolderName    string // 对应 public/images/all_themes/ 下的文件夹名
}

// RendererService 渲染服务
type RendererService struct {
	styles map[string]StyleConfig
}

// NewRendererService 创建渲染服务实例
func NewRendererService() *RendererService {
	return &RendererService{
		styles: map[string]StyleConfig{
			"default": {
				Key:           "default",
				Name:          "简约灰",
				Primary:       color.RGBA{R: 100, G: 100, B: 100, A: 255},
				Secondary:     color.RGBA{R: 150, G: 150, B: 150, A: 255},
				Background:    color.RGBA{R: 245, G: 245, B: 245, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 119, G: 119, B: 119, A: 255},
				Accent:        color.RGBA{R: 80, G: 80, B: 80, A: 255},
				FolderName:    "playful-geometric",
			},
			"xiaohongshu": {
				Key:           "xiaohongshu",
				Name:          "小红书红",
				Primary:       color.RGBA{R: 255, G: 66, B: 99, A: 255},
				Secondary:     color.RGBA{R: 255, G: 110, B: 136, A: 255},
				Background:    color.RGBA{R: 255, G: 250, B: 251, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 102, G: 102, B: 102, A: 255},
				Accent:        color.RGBA{R: 255, G: 66, B: 99, A: 255},
				FolderName:    "playful-geometric",
			},
			"purple": {
				Key:           "purple",
				Name:          "紫韵",
				Primary:       color.RGBA{R: 147, G: 112, B: 219, A: 255},
				Secondary:     color.RGBA{R: 187, G: 154, B: 247, A: 255},
				Background:    color.RGBA{R: 250, G: 248, B: 255, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 102, G: 102, B: 102, A: 255},
				Accent:        color.RGBA{R: 147, G: 112, B: 219, A: 255},
				FolderName:    "retro",
			},
			"mint": {
				Key:           "mint",
				Name:          "清新薄荷",
				Primary:       color.RGBA{R: 72, G: 187, B: 120, A: 255},
				Secondary:     color.RGBA{R: 129, G: 230, B: 176, A: 255},
				Background:    color.RGBA{R: 240, G: 253, B: 244, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 102, G: 102, B: 102, A: 255},
				Accent:        color.RGBA{R: 72, G: 187, B: 120, A: 255},
				FolderName:    "Sketch",
			},
			"sunset": {
				Key:           "sunset",
				Name:          "日落橙",
				Primary:       color.RGBA{R: 251, G: 146, B: 60, A: 255},
				Secondary:     color.RGBA{R: 255, G: 183, B: 77, A: 255},
				Background:    color.RGBA{R: 255, G: 250, B: 240, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 102, G: 102, B: 102, A: 255},
				Accent:        color.RGBA{R: 251, G: 146, B: 60, A: 255},
				FolderName:    "retro",
			},
			"ocean": {
				Key:           "ocean",
				Name:          "深海蓝",
				Primary:       color.RGBA{R: 59, G: 130, B: 246, A: 255},
				Secondary:     color.RGBA{R: 96, G: 165, B: 250, A: 255},
				Background:    color.RGBA{R: 239, G: 246, B: 255, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 102, G: 102, B: 102, A: 255},
				Accent:        color.RGBA{R: 59, G: 130, B: 246, A: 255},
				FolderName:    "terminal",
			},
			"elegant": {
				Key:           "elegant",
				Name:          "优雅白",
				Primary:       color.RGBA{R: 107, G: 114, B: 128, A: 255},
				Secondary:     color.RGBA{R: 148, G: 163, B: 184, A: 255},
				Background:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 31, G: 41, B: 55, A: 255},
				TextSecondary: color.RGBA{R: 107, G: 114, B: 128, A: 255},
				Accent:        color.RGBA{R: 107, G: 114, B: 128, A: 255},
				FolderName:    "auto-fit",
			},
			"dark": {
				Key:           "dark",
				Name:          "暗黑模式",
				Primary:       color.RGBA{R: 249, G: 115, B: 22, A: 255},
				Secondary:     color.RGBA{R: 251, G: 146, B: 60, A: 255},
				Background:    color.RGBA{R: 17, G: 24, B: 39, A: 255},
				CardInner:     color.RGBA{R: 31, G: 41, B: 55, A: 255},
				TextPrimary:   color.RGBA{R: 249, G: 250, B: 251, A: 255},
				TextSecondary: color.RGBA{R: 148, G: 163, B: 184, A: 255},
				Accent:        color.RGBA{R: 249, G: 115, B: 22, A: 255},
				FolderName:    "terminal",
			},
		},
	}
}

// GetStyles 获取所有可用样式
func (s *RendererService) GetStyles() []StyleConfig {
	styles := make([]StyleConfig, 0, len(s.styles))
	for _, style := range s.styles {
		styles = append(styles, style)
	}
	return styles
}

// RenderMarkdownToImage 将 Markdown 渲染为图片
func (s *RendererService) RenderMarkdownToImage(markdown, styleKey, outputPrefix string, cardWidth, cardHeight, maxContentHeight int) ([]string, error) {
	// 获取样式配置，默认为 playful-geometric
	style, exists := s.styles[styleKey]
	if !exists {
		style = s.styles["default"]
	}

	// 根据样式返回对应的预设图片
	var images []string
	basePath := filepath.Join("all_themes", style.FolderName)

	// 根据不同的样式返回不同数量的卡片
	switch style.FolderName {
	case "Sketch", "playful-geometric", "retro", "terminal":
		// 这些样式有5张卡片
		images = []string{
			filepath.Join(basePath, "card_1.png"),
			filepath.Join(basePath, "card_2.png"),
			filepath.Join(basePath, "card_3.png"),
			filepath.Join(basePath, "card_4.png"),
			filepath.Join(basePath, "card_5.png"),
		}
	default:
		// 其他样式有1张卡片
		images = []string{
			filepath.Join(basePath, "card_1.png"),
		}
	}

	return images, nil
}

// GenerateCoverOnly 只生成封面图
func (s *RendererService) GenerateCoverOnly(title, subtitle, styleKey, outputPrefix string, width, height int) (string, error) {
	// 获取样式配置，默认为 playful-geometric
	style, exists := s.styles[styleKey]
	if !exists {
		style = s.styles["default"]
	}

	// 返回对应的封面图片
	coverPath := filepath.Join("all_themes", style.FolderName, "cover.png")
	return coverPath, nil
}
