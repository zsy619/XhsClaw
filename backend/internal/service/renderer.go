// Package service 提供业务逻辑层 - 小红书图片渲染服务
// 使用 Chromium (chromedp) 渲染 HTML 生成图片
package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/chromedp/chromedp"
)

// 默认尺寸配置 (3:4 比例)
const (
	DefaultWidth  = 1080
	DefaultHeight = 1440
	MaxHeight     = 4320 // dynamic 模式最大高度
)

// PaginationMode 分页模式类型
type PaginationMode string

const (
	PaginationAutoSplit PaginationMode = "auto-split"
	PaginationSeparator PaginationMode = "separator"
	PaginationAutoFit   PaginationMode = "auto-fit"
	PaginationDynamic   PaginationMode = "dynamic"
)

// ThemeConfig 主题配置
type ThemeConfig struct {
	Key           string
	Name          string
	CoverBg       string // 封面背景渐变
	CardBg        string // 卡片背景渐变
	TitleGradient string // 标题文字渐变
	AccentColor   string // 强调色
}

// templateCache 模板缓存
type templateCache struct {
	templates    map[string]*template.Template
	lastModified map[string]time.Time
}

// RendererService 渲染服务
type RendererService struct {
	themes    map[string]ThemeConfig
	templates *templateCache
	imagesDir string
	assetsDir string
}

// NewRendererService 创建渲染服务实例
func NewRendererService() (*RendererService, error) {
	// 获取项目根目录
	projectRoot := getProjectRoot()
	assetsDir := filepath.Join(projectRoot, "assets")
	imagesDir := filepath.Join(projectRoot, "public", "images")

	// 确保目录存在
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return nil, fmt.Errorf("创建图片目录失败: %v", err)
	}

	s := &RendererService{
		imagesDir: imagesDir,
		assetsDir: assetsDir,
		themes:    make(map[string]ThemeConfig),
		templates: &templateCache{
			templates:    make(map[string]*template.Template),
			lastModified: make(map[string]time.Time),
		},
	}

	// 初始化主题配置
	s.initThemes()

	// 加载模板
	if err := s.loadTemplates(); err != nil {
		return nil, fmt.Errorf("加载模板失败: %v", err)
	}

	return s, nil
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 尝试从环境变量获取
	if root := os.Getenv("PROJECT_ROOT"); root != "" {
		return root
	}
	// 默认路径
	return "/Volumes/E/JYW/创意项目/XhsClaw/backend"
}

// initThemes 初始化主题配置
func (s *RendererService) initThemes() {
	s.themes = map[string]ThemeConfig{
		"default": {
			Key:           "default",
			Name:          "简约灰",
			CoverBg:       "linear-gradient(180deg, #f3f3f3 0%, #f9f9f9 100%)",
			CardBg:        "linear-gradient(180deg, #f3f3f3 0%, #f9f9f9 100%)",
			TitleGradient: "linear-gradient(180deg, #111827 0%, #4B5563 100%)",
			AccentColor:   "#6366f1",
		},
		"xiaohongshu": {
			Key:           "xiaohongshu",
			Name:          "小红书红",
			CoverBg:       "linear-gradient(180deg, #ff4757 0%, #ff6b81 100%)",
			CardBg:        "linear-gradient(135deg, #ff4757 0%, #ff6b81 100%)",
			TitleGradient: "linear-gradient(180deg, #ff2442 0%, #ff6b81 100%)",
			AccentColor:   "#ff2442",
		},
		"purple": {
			Key:           "purple",
			Name:          "紫韵",
			CoverBg:       "linear-gradient(180deg, #8b5cf6 0%, #a78bfa 100%)",
			CardBg:        "linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)",
			TitleGradient: "linear-gradient(180deg, #7c3aed 0%, #a78bfa 100%)",
			AccentColor:   "#8b5cf6",
		},
		"mint": {
			Key:           "mint",
			Name:          "清新薄荷",
			CoverBg:       "linear-gradient(180deg, #22c55e 0%, #4ade80 100%)",
			CardBg:        "linear-gradient(135deg, #22c55e 0%, #4ade80 100%)",
			TitleGradient: "linear-gradient(180deg, #16a34a 0%, #4ade80 100%)",
			AccentColor:   "#22c55e",
		},
		"sunset": {
			Key:           "sunset",
			Name:          "日落橙",
			CoverBg:       "linear-gradient(180deg, #f97316 0%, #fb923c 100%)",
			CardBg:        "linear-gradient(135deg, #f97316 0%, #fb923c 100%)",
			TitleGradient: "linear-gradient(180deg, #ea580c 0%, #fb923c 100%)",
			AccentColor:   "#f97316",
		},
		"ocean": {
			Key:           "ocean",
			Name:          "深海蓝",
			CoverBg:       "linear-gradient(180deg, #0284c7 0%, #38bdf8 100%)",
			CardBg:        "linear-gradient(135deg, #0284c7 0%, #38bdf8 100%)",
			TitleGradient: "linear-gradient(180deg, #0369a1 0%, #38bdf8 100%)",
			AccentColor:   "#0284c7",
		},
		"elegant": {
			Key:           "elegant",
			Name:          "优雅白",
			CoverBg:       "linear-gradient(180deg, #e5e5e5 0%, #f5f5f5 100%)",
			CardBg:        "linear-gradient(135deg, #e5e5e5 0%, #fafafa 100%)",
			TitleGradient: "linear-gradient(180deg, #171717 0%, #404040 100%)",
			AccentColor:   "#171717",
		},
		"dark": {
			Key:           "dark",
			Name:          "暗黑模式",
			CoverBg:       "linear-gradient(180deg, #171717 0%, #262626 100%)",
			CardBg:        "linear-gradient(135deg, #0a0a0a 0%, #171717 100%)",
			TitleGradient: "linear-gradient(180deg, #ffffff 0%, #a3a3a3 100%)",
			AccentColor:   "#3b82f6",
		},
		"playful-geometric": {
			Key:           "playful-geometric",
			Name:          "活泼几何",
			CoverBg:       "linear-gradient(180deg, #8B5CF6 0%, #F472B6 100%)",
			CardBg:        "linear-gradient(135deg, #8B5CF6 0%, #F472B6 100%)",
			TitleGradient: "linear-gradient(180deg, #7C3AED 0%, #F472B6 100%)",
			AccentColor:   "#8B5CF6",
		},
		"neo-brutalism": {
			Key:           "neo-brutalism",
			Name:          "新野兽派",
			CoverBg:       "linear-gradient(180deg, #FF4757 0%, #FECA57 100%)",
			CardBg:        "linear-gradient(135deg, #FF4757 0%, #FECA57 100%)",
			TitleGradient: "linear-gradient(180deg, #000000 0%, #FF4757 100%)",
			AccentColor:   "#FF4757",
		},
		"botanical": {
			Key:           "botanical",
			Name:          "植物系",
			CoverBg:       "linear-gradient(180deg, #4A7C59 0%, #8FBC8F 100%)",
			CardBg:        "linear-gradient(135deg, #4A7C59 0%, #8FBC8F 100%)",
			TitleGradient: "linear-gradient(180deg, #1F2937 0%, #4A7C59 100%)",
			AccentColor:   "#22c55e",
		},
		"professional": {
			Key:           "professional",
			Name:          "专业商务",
			CoverBg:       "linear-gradient(180deg, #2563EB 0%, #3B82F6 100%)",
			CardBg:        "linear-gradient(135deg, #2563EB 0%, #3B82F6 100%)",
			TitleGradient: "linear-gradient(180deg, #1E3A8A 0%, #2563EB 100%)",
			AccentColor:   "#2563eb",
		},
		"retro": {
			Key:           "retro",
			Name:          "复古风格",
			CoverBg:       "linear-gradient(180deg, #D35400 0%, #F39C12 100%)",
			CardBg:        "linear-gradient(135deg, #D35400 0%, #F39C12 100%)",
			TitleGradient: "linear-gradient(180deg, #8B4513 0%, #D35400 100%)",
			AccentColor:   "#ea580c",
		},
		"terminal": {
			Key:           "terminal",
			Name:          "终端风格",
			CoverBg:       "linear-gradient(180deg, #0D1117 0%, #21262D 100%)",
			CardBg:        "linear-gradient(135deg, #0D1117 0%, #161B22 100%)",
			TitleGradient: "linear-gradient(180deg, #39D353 0%, #58A6FF 100%)",
			AccentColor:   "#39d353",
		},
		"sketch": {
			Key:           "sketch",
			Name:          "手绘风格",
			CoverBg:       "linear-gradient(180deg, #555555 0%, #999999 100%)",
			CardBg:        "linear-gradient(135deg, #555555 0%, #888888 100%)",
			TitleGradient: "linear-gradient(180deg, #111827 0%, #6B7280 100%)",
			AccentColor:   "#ec4899",
		},
		"pink-cream": {
			Key:           "pink-cream",
			Name:          "粉色奶油",
			CoverBg:       "linear-gradient(180deg, #f472b6 0%, #ec4899 100%)",
			CardBg:        "linear-gradient(135deg, #f472b6 0%, #ec4899 100%)",
			TitleGradient: "linear-gradient(180deg, #db2777 0%, #f472b6 100%)",
			AccentColor:   "#ec4899",
		},
		"coral": {
			Key:           "coral",
			Name:          "珊瑚粉",
			CoverBg:       "linear-gradient(180deg, #fb7185 0%, #f43f5e 100%)",
			CardBg:        "linear-gradient(135deg, #fb7185 0%, #f43f5e 100%)",
			TitleGradient: "linear-gradient(180deg, #e11d48 0%, #fb7185 100%)",
			AccentColor:   "#f43f5e",
		},
		"lavender": {
			Key:           "lavender",
			Name:          "薰衣草紫",
			CoverBg:       "linear-gradient(180deg, #a78bfa 0%, #8b5cf6 100%)",
			CardBg:        "linear-gradient(135deg, #a78bfa 0%, #8b5cf6 100%)",
			TitleGradient: "linear-gradient(180deg, #7c3aed 0%, #a78bfa 100%)",
			AccentColor:   "#8b5cf6",
		},
		"cream": {
			Key:           "cream",
			Name:          "奶黄包",
			CoverBg:       "linear-gradient(180deg, #fbbf24 0%, #f59e0b 100%)",
			CardBg:        "linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%)",
			TitleGradient: "linear-gradient(180deg, #d97706 0%, #fbbf24 100%)",
			AccentColor:   "#f59e0b",
		},
		"nordic": {
			Key:           "nordic",
			Name:          "北欧风格",
			CoverBg:       "linear-gradient(180deg, #e2e8f0 0%, #f1f5f9 100%)",
			CardBg:        "linear-gradient(135deg, #e2e8f0 0%, #f8fafc 100%)",
			TitleGradient: "linear-gradient(180deg, #1e293b 0%, #475569 100%)",
			AccentColor:   "#334155",
		},
		"peach": {
			Key:           "peach",
			Name:          "蜜桃粉",
			CoverBg:       "linear-gradient(180deg, #fda4af 0%, #f43f5e 100%)",
			CardBg:        "linear-gradient(135deg, #fda4af 0%, #f43f5e 100%)",
			TitleGradient: "linear-gradient(180deg, #e11d48 0%, #fb7185 100%)",
			AccentColor:   "#f43f5e",
		},
		// 新增 20 个小红书风格样式
		"cream-custard": {
			Key:           "cream-custard",
			Name:          "奶油布丁",
			CoverBg:       "linear-gradient(180deg, #fffbf0 0%, #fff5d6 100%)",
			CardBg:        "linear-gradient(135deg, #fff9e6 0%, #fff5d6 100%)",
			TitleGradient: "linear-gradient(135deg, #f5e6c8 0%, #f0d5a8 100%)",
			AccentColor:   "#d4a574",
		},
		"sakura-pink": {
			Key:           "sakura-pink",
			Name:          "樱花粉",
			CoverBg:       "linear-gradient(180deg, #fff5f7 0%, #ffeef2 100%)",
			CardBg:        "linear-gradient(135deg, #fff0f5 0%, #ffeef2 100%)",
			TitleGradient: "linear-gradient(135deg, #ffb7c5 0%, #ff9eb5 100%)",
			AccentColor:   "#ff9eb5",
		},
		"matcha-latte": {
			Key:           "matcha-latte",
			Name:          "抹茶拿铁",
			CoverBg:       "linear-gradient(180deg, #f7f9f4 0%, #e8ede3 100%)",
			CardBg:        "linear-gradient(135deg, #f7f9f4 0%, #e8ede3 100%)",
			TitleGradient: "linear-gradient(135deg, #6b8c5e 0%, #8ba874 100%)",
			AccentColor:   "#6b8c5e",
		},
		"blueberry-cheese": {
			Key:           "blueberry-cheese",
			Name:          "蓝莓芝士",
			CoverBg:       "linear-gradient(180deg, #f5f6fa 0%, #e8eaf6 100%)",
			CardBg:        "linear-gradient(135deg, #f5f6fa 0%, #e8eaf6 100%)",
			TitleGradient: "linear-gradient(135deg, #5c6bc0 0%, #7986cb 100%)",
			AccentColor:   "#5c6bc0",
		},
		"caramel-macchiato": {
			Key:           "caramel-macchiato",
			Name:          "焦糖玛奇朵",
			CoverBg:       "linear-gradient(180deg, #faf6f1 0%, #f0e9df 100%)",
			CardBg:        "linear-gradient(135deg, #faf6f1 0%, #f0e9df 100%)",
			TitleGradient: "linear-gradient(135deg, #8b6f47 0%, #a08060 100%)",
			AccentColor:   "#8b6f47",
		},
		"honey-peach": {
			Key:           "honey-peach",
			Name:          "蜜桃蜂蜜",
			CoverBg:       "linear-gradient(180deg, #fff9f5 0%, #fff0e8 100%)",
			CardBg:        "linear-gradient(135deg, #fff9f5 0%, #fff0e8 100%)",
			TitleGradient: "linear-gradient(135deg, #e88d67 0%, #ff997a 100%)",
			AccentColor:   "#e88d67",
		},
		"vanilla-milk": {
			Key:           "vanilla-milk",
			Name:          "香草牛奶",
			CoverBg:       "linear-gradient(180deg, #fafaf9 0%, #f5f5f0 100%)",
			CardBg:        "linear-gradient(135deg, #fafaf9 0%, #f5f5f0 100%)",
			TitleGradient: "linear-gradient(135deg, #6d6d5a 0%, #8c8c7a 100%)",
			AccentColor:   "#6d6d5a",
		},
		"chocolate-mint": {
			Key:           "chocolate-mint",
			Name:          "巧克力薄荷",
			CoverBg:       "linear-gradient(180deg, #f5f9f7 0%, #e8f0eb 100%)",
			CardBg:        "linear-gradient(135deg, #f5f9f7 0%, #e8f0eb 100%)",
			TitleGradient: "linear-gradient(135deg, #4a6755 0%, #5d8a6a 100%)",
			AccentColor:   "#4a6755",
		},
		"strawberry-milk": {
			Key:           "strawberry-milk",
			Name:          "草莓牛奶",
			CoverBg:       "linear-gradient(180deg, #fff5f8 0%, #ffeef2 100%)",
			CardBg:        "linear-gradient(135deg, #fff5f8 0%, #ffeef2 100%)",
			TitleGradient: "linear-gradient(135deg, #e86b8a 0%, #ff859e 100%)",
			AccentColor:   "#e86b8a",
		},
		"mango-pudding": {
			Key:           "mango-pudding",
			Name:          "芒果布丁",
			CoverBg:       "linear-gradient(180deg, #fffdf0 0%, #fff9e0 100%)",
			CardBg:        "linear-gradient(135deg, #fffdf0 0%, #fff9e0 100%)",
			TitleGradient: "linear-gradient(135deg, #e8a838 0%, #ffb84a 100%)",
			AccentColor:   "#e8a838",
		},
		"taro-milktea": {
			Key:           "taro-milktea",
			Name:          "芋头奶茶",
			CoverBg:       "linear-gradient(180deg, #f9f6fa 0%, #f0e8f0 100%)",
			CardBg:        "linear-gradient(135deg, #f9f6fa 0%, #f0e8f0 100%)",
			TitleGradient: "linear-gradient(135deg, #8b6b8c 0%, #a080a0 100%)",
			AccentColor:   "#8b6b8c",
		},
		"coconut-cream": {
			Key:           "coconut-cream",
			Name:          "椰子奶油",
			CoverBg:       "linear-gradient(180deg, #fafaf9 0%, #f5f5f4 100%)",
			CardBg:        "linear-gradient(135deg, #fafaf9 0%, #f5f5f4 100%)",
			TitleGradient: "linear-gradient(135deg, #8b8574 0%, #a09a8a 100%)",
			AccentColor:   "#8b8574",
		},
		"red-velvet": {
			Key:           "red-velvet",
			Name:          "红丝绒",
			CoverBg:       "linear-gradient(180deg, #fef6f6 0%, #f5e8e8 100%)",
			CardBg:        "linear-gradient(135deg, #fef6f6 0%, #f5e8e8 100%)",
			TitleGradient: "linear-gradient(135deg, #c44569 0%, #e05676 100%)",
			AccentColor:   "#c44569",
		},
		"pistachio-green": {
			Key:           "pistachio-green",
			Name:          "开心果绿",
			CoverBg:       "linear-gradient(180deg, #f7f9f7 0%, #e8efe8 100%)",
			CardBg:        "linear-gradient(135deg, #f7f9f7 0%, #e8efe8 100%)",
			TitleGradient: "linear-gradient(135deg, #5e8c61 0%, #7ab380 100%)",
			AccentColor:   "#5e8c61",
		},
		"bubblegum-pink": {
			Key:           "bubblegum-pink",
			Name:          "泡泡糖粉",
			CoverBg:       "linear-gradient(180deg, #fff0f5 0%, #ffe0eb 100%)",
			CardBg:        "linear-gradient(135deg, #fff0f5 0%, #ffe0eb 100%)",
			TitleGradient: "linear-gradient(135deg, #e86b9e 0%, #ff85b3 100%)",
			AccentColor:   "#e86b9e",
		},
		"lemon-meringue": {
			Key:           "lemon-meringue",
			Name:          "柠檬蛋白",
			CoverBg:       "linear-gradient(180deg, #fffdf5 0%, #fff9e6 100%)",
			CardBg:        "linear-gradient(135deg, #fffdf5 0%, #fff9e6 100%)",
			TitleGradient: "linear-gradient(135deg, #e8c838 0%, #ffd84a 100%)",
			AccentColor:   "#e8c838",
		},
		"blackberry-sage": {
			Key:           "blackberry-sage",
			Name:          "黑莓鼠尾草",
			CoverBg:       "linear-gradient(180deg, #f6f6f9 0%, #e8e8f0 100%)",
			CardBg:        "linear-gradient(135deg, #f6f6f9 0%, #e8e8f0 100%)",
			TitleGradient: "linear-gradient(135deg, #5c5c8c 0%, #7575a8 100%)",
			AccentColor:   "#5c5c8c",
		},
		"peaches-cream": {
			Key:           "peaches-cream",
			Name:          "蜜桃奶油",
			CoverBg:       "linear-gradient(180deg, #fff9f7 0%, #fff0ec 100%)",
			CardBg:        "linear-gradient(135deg, #fff9f7 0%, #fff0ec 100%)",
			TitleGradient: "linear-gradient(135deg, #e87d67 0%, #ff9585 100%)",
			AccentColor:   "#e87d67",
		},
		"earl-grey": {
			Key:           "earl-grey",
			Name:          "伯爵茶",
			CoverBg:       "linear-gradient(180deg, #f7f7f9 0%, #eeeff4 100%)",
			CardBg:        "linear-gradient(135deg, #f7f7f9 0%, #eeeff4 100%)",
			TitleGradient: "linear-gradient(135deg, #6b6b8c 0%, #8585a8 100%)",
			AccentColor:   "#6b6b8c",
		},
		"tiramisu": {
			Key:           "tiramisu",
			Name:          "提拉米苏",
			CoverBg:       "linear-gradient(180deg, #faf8f5 0%, #f0ebe3 100%)",
			CardBg:        "linear-gradient(135deg, #faf8f5 0%, #f0ebe3 100%)",
			TitleGradient: "linear-gradient(135deg, #8b7355 0%, #a0856a 100%)",
			AccentColor:   "#8b7355",
		},
		"pomegranate": {
			Key:           "pomegranate",
			Name:          "石榴",
			CoverBg:       "linear-gradient(180deg, #fef6f9 0%, #f5e8ec 100%)",
			CardBg:        "linear-gradient(135deg, #fef6f9 0%, #f5e8ec 100%)",
			TitleGradient: "linear-gradient(135deg, #c44569 0%, #d8567a 100%)",
			AccentColor:   "#c44569",
		},
		"sage-green": {
			Key:           "sage-green",
			Name:          "鼠尾草绿",
			CoverBg:       "linear-gradient(180deg, #f7f9f8 0%, #e8efe9 100%)",
			CardBg:        "linear-gradient(135deg, #f7f9f8 0%, #e8efe9 100%)",
			TitleGradient: "linear-gradient(135deg, #5e7a6b 0%, #7a9a85 100%)",
			AccentColor:   "#5e7a6b",
		},
		"honey-ginger": {
			Key:           "honey-ginger",
			Name:          "蜂蜜姜茶",
			CoverBg:       "linear-gradient(180deg, #fff8f0 0%, #fff0e0 100%)",
			CardBg:        "linear-gradient(135deg, #fff5e6 0%, #ffe8d0 100%)",
			TitleGradient: "linear-gradient(135deg, #d4a574 0%, #e8c49a 100%)",
			AccentColor:   "#d4a574",
		},
		"rose-milk": {
			Key:           "rose-milk",
			Name:          "玫瑰奶茶",
			CoverBg:       "linear-gradient(180deg, #fff5f7 0%, #ffeef2 100%)",
			CardBg:        "linear-gradient(135deg, #fff0f5 0%, #ffe4eb 100%)",
			TitleGradient: "linear-gradient(135deg, #c98b95 0%, #e0a5b0 100%)",
			AccentColor:   "#c98b95",
		},
		"lavender-honey": {
			Key:           "lavender-honey",
			Name:          "薰衣草蜂蜜",
			CoverBg:       "linear-gradient(180deg, #f5f3ff 0%, #ede9fe 100%)",
			CardBg:        "linear-gradient(135deg, #f0edff 0%, #e2dff5 100%)",
			TitleGradient: "linear-gradient(135deg, #8b7ec8 0%, #a795d5 100%)",
			AccentColor:   "#8b7ec8",
		},
		"blue-lagoon": {
			Key:           "blue-lagoon",
			Name:          "蓝色泻湖",
			CoverBg:       "linear-gradient(180deg, #f0f9ff 0%, #e0f2fe 100%)",
			CardBg:        "linear-gradient(135deg, #e0f7fa 0%, #b2ebf2 100%)",
			TitleGradient: "linear-gradient(135deg, #0288d1 0%, #29b6f6 100%)",
			AccentColor:   "#0288d1",
		},
		"aurora-green": {
			Key:           "aurora-green",
			Name:          "极光绿",
			CoverBg:       "linear-gradient(180deg, #f0fdf4 0%, #dcfce7 100%)",
			CardBg:        "linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%)",
			TitleGradient: "linear-gradient(135deg, #059669 0%, #10b981 100%)",
			AccentColor:   "#059669",
		},
		"pearl-white": {
			Key:           "pearl-white",
			Name:          "珍珠白",
			CoverBg:       "linear-gradient(180deg, #fafafa 0%, #f5f5f5 100%)",
			CardBg:        "linear-gradient(135deg, #ffffff 0%, #f8f8f8 100%)",
			TitleGradient: "linear-gradient(135deg, #404040 0%, #737373 100%)",
			AccentColor:   "#404040",
		},
		"blush-pink": {
			Key:           "blush-pink",
			Name:          "腮红粉",
			CoverBg:       "linear-gradient(180deg, #fff1f2 0%, #ffe4e6 100%)",
			CardBg:        "linear-gradient(135deg, #ffe4e6 0%, #fecdd3 100%)",
			TitleGradient: "linear-gradient(135deg, #e11d48 0%, #fb7185 100%)",
			AccentColor:   "#e11d48",
		},
		"ocean-mist": {
			Key:           "ocean-mist",
			Name:          "海雾",
			CoverBg:       "linear-gradient(180deg, #f0f8ff 0%, #e6f2ff 100%)",
			CardBg:        "linear-gradient(135deg, #f5f9ff 0%, #e0efff 100%)",
			TitleGradient: "linear-gradient(135deg, #1e40af 0%, #3b82f6 100%)",
			AccentColor:   "#1e40af",
		},
		"lily-white": {
			Key:           "lily-white",
			Name:          "百合白",
			CoverBg:       "linear-gradient(180deg, #fefefe 0%, #f8f8f8 100%)",
			CardBg:        "linear-gradient(135deg, #fafafa 0%, #f0f0f0 100%)",
			TitleGradient: "linear-gradient(135deg, #2d3748 0%, #4a5568 100%)",
			AccentColor:   "#2d3748",
		},
		"sun-kissed": {
			Key:           "sun-kissed",
			Name:          "阳光亲吻",
			CoverBg:       "linear-gradient(180deg, #fffbeb 0%, #fef3c7 100%)",
			CardBg:        "linear-gradient(135deg, #fffbeb 0%, #fde68a 100%)",
			TitleGradient: "linear-gradient(135deg, #d97706 0%, #f59e0b 100%)",
			AccentColor:   "#d97706",
		},
		"berry-smoothie": {
			Key:           "berry-smoothie",
			Name:          "莓果奶昔",
			CoverBg:       "linear-gradient(180deg, #fdf2f8 0%, #fce7f3 100%)",
			CardBg:        "linear-gradient(135deg, #fdf4ff 0%, #fae8ff 100%)",
			TitleGradient: "linear-gradient(135deg, #c026d3 0%, #db2777 100%)",
			AccentColor:   "#c026d3",
		},
		"winter-sky": {
			Key:           "winter-sky",
			Name:          "冬日天空",
			CoverBg:       "linear-gradient(180deg, #f0f9ff 0%, #e0f2fe 100%)",
			CardBg:        "linear-gradient(135deg, #f5faff 0%, #e0f2fe 100%)",
			TitleGradient: "linear-gradient(135deg, #0369a1 0%, #0ea5e9 100%)",
			AccentColor:   "#0369a1",
		},
		"ivory-cream": {
			Key:           "ivory-cream",
			Name:          "象牙奶油",
			CoverBg:       "linear-gradient(180deg, #fffff0 0%, #fffaf0 100%)",
			CardBg:        "linear-gradient(135deg, #fffef5 0%, #fff8e7 100%)",
			TitleGradient: "linear-gradient(135deg, #78716c 0%, #a8a29e 100%)",
			AccentColor:   "#78716c",
		},
		"floral-pink": {
			Key:           "floral-pink",
			Name:          "花漾粉",
			CoverBg:       "linear-gradient(180deg, #fdf4ff 0%, #fae8ff 100%)",
			CardBg:        "linear-gradient(135deg, #faf5ff 0%, #f3d8fa 100%)",
			TitleGradient: "linear-gradient(135deg, #a855f7 0%, #c084fc 100%)",
			AccentColor:   "#a855f7",
		},
		"mint-chocolate": {
			Key:           "mint-chocolate",
			Name:          "薄荷巧克力",
			CoverBg:       "linear-gradient(180deg, #f0fdf4 0%, #dcfce7 100%)",
			CardBg:        "linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%)",
			TitleGradient: "linear-gradient(135deg, #047857 0%, #059669 100%)",
			AccentColor:   "#047857",
		},
		"autumn-leaves": {
			Key:           "autumn-leaves",
			Name:          "秋日落叶",
			CoverBg:       "linear-gradient(180deg, #fff7ed 0%, #ffedd5 100%)",
			CardBg:        "linear-gradient(135deg, #fff5f2 0%, #ffe4d6 100%)",
			TitleGradient: "linear-gradient(135deg, #c2410c 0%, #ea580c 100%)",
			AccentColor:   "#c2410c",
		},
		"rainbow-sorbet": {
			Key:           "rainbow-sorbet",
			Name:          "彩虹冰糕",
			CoverBg:       "linear-gradient(180deg, #fefce8 0%, #fef9c3 100%)",
			CardBg:        "linear-gradient(135deg, #fefce8 0%, #fef08a 100%)",
			TitleGradient: "linear-gradient(135deg, #a16207 0%, #ca8a04 100%)",
			AccentColor:   "#a16207",
		},
		"cherry-blush": {
			Key:           "cherry-blush",
			Name:          "Cherry腮红",
			CoverBg:       "linear-gradient(180deg, #fff1f2 0%, #ffe4e6 100%)",
			CardBg:        "linear-gradient(135deg, #fff0f3 0%, #ffe3e8 100%)",
			TitleGradient: "linear-gradient(135deg, #be123c 0%, #e11d48 100%)",
			AccentColor:   "#be123c",
		},
		"sea-glass": {
			Key:           "sea-glass",
			Name:          "海玻璃",
			CoverBg:       "linear-gradient(180deg, #f0fdfa 0%, #ccfbf1 100%)",
			CardBg:        "linear-gradient(135deg, #f0fdfa 0%, #99f6e4 100%)",
			TitleGradient: "linear-gradient(135deg, #0d9488 0%, #14b8a6 100%)",
			AccentColor:   "#0d9488",
		},
		"cotton-candy": {
			Key:           "cotton-candy",
			Name:          "棉花糖",
			CoverBg:       "linear-gradient(180deg, #fdf4ff 0%, #fce7f3 100%)",
			CardBg:        "linear-gradient(135deg, #fdf2f8 0%, #fbcfe8 100%)",
			TitleGradient: "linear-gradient(135deg, #db2777 0%, #ec4899 100%)",
			AccentColor:   "#db2777",
		},
	}
}

// loadTemplates 加载HTML模板
func (s *RendererService) loadTemplates() error {
	templatesDir := filepath.Join(s.assetsDir, "templates")

	// 加载封面模板
	coverPath := filepath.Join(templatesDir, "cover.html")
	if err := s.loadTemplate("cover", coverPath); err != nil {
		return err
	}

	// 加载卡片模板
	cardPath := filepath.Join(templatesDir, "card.html")
	if err := s.loadTemplate("card", cardPath); err != nil {
		return err
	}

	return nil
}

// loadTemplate 加载单个模板
func (s *RendererService) loadTemplate(name, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取%s模板失败: %v", name, err)
	}

	tmpl := template.Must(template.New(name).Parse(string(content)))
	s.templates.templates[name] = tmpl

	// 更新最后修改时间
	if info, err := os.Stat(path); err == nil {
		s.templates.lastModified[name] = info.ModTime()
	}

	return nil
}

// getTemplate 获取模板（带缓存检查）
func (s *RendererService) getTemplate(name string) (*template.Template, error) {
	templatesDir := filepath.Join(s.assetsDir, "templates")
	path := filepath.Join(templatesDir, name+".html")

	// 检查文件是否存在
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("模板文件不存在: %v", err)
	}

	// 检查是否需要重新加载
	lastModified, exists := s.templates.lastModified[name]
	if !exists || info.ModTime().After(lastModified) {
		// 重新加载模板
		if err := s.loadTemplate(name, path); err != nil {
			return nil, err
		}
	}

	tmpl, exists := s.templates.templates[name]
	if !exists {
		return nil, fmt.Errorf("模板未加载: %s", name)
	}

	return tmpl, nil
}

// GetImagesDir 获取图片存储目录
func (s *RendererService) GetImagesDir() string {
	return s.imagesDir
}

// GetStyles 获取所有样式配置
func (s *RendererService) GetStyles() []ThemeConfig {
	styles := make([]ThemeConfig, 0, len(s.themes))
	for _, style := range s.themes {
		styles = append(styles, style)
	}
	return styles
}

// GetStyle 获取指定样式配置
func (s *RendererService) GetStyle(key string) ThemeConfig {
	style, exists := s.themes[key]
	if !exists {
		style = s.themes["default"]
	}
	return style
}

// loadThemeCSS 加载主题CSS
func (s *RendererService) loadThemeCSS(theme string) string {
	cssPath := filepath.Join(s.assetsDir, "themes", theme+".css")
	if content, err := os.ReadFile(cssPath); err == nil {
		return fmt.Sprintf("<style>%s</style>", string(content))
	}
	// 返回默认主题
	defaultPath := filepath.Join(s.assetsDir, "themes", "default.css")
	if content, err := os.ReadFile(defaultPath); err == nil {
		return fmt.Sprintf("<style>%s</style>", string(content))
	}
	return ""
}

// renderTemplate 统一渲染模板
func (s *RendererService) renderTemplate(templateName string, data interface{}) (string, error) {
	// 验证模板数据
	if err := s.validateTemplateData(data); err != nil {
		return "", err
	}
	
	// 获取模板（带缓存检查）
	tmpl, err := s.getTemplate(templateName)
	if err != nil {
		return "", err
	}
	
	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, data); err != nil {
		return "", fmt.Errorf("渲染%s模板失败: %v", templateName, err)
	}
	return htmlBuf.String(), nil
}

// validateTemplateData 验证模板数据
func (s *RendererService) validateTemplateData(data interface{}) error {
	switch d := data.(type) {
	case CardData:
		if d.Width <= 0 {
			return fmt.Errorf("无效的卡片宽度: %d", d.Width)
		}
		if d.Height <= 0 {
			return fmt.Errorf("无效的卡片高度: %d", d.Height)
		}
		if d.Content == "" {
			return fmt.Errorf("卡片内容不能为空")
		}
		if d.FontSize < 0 {
			return fmt.Errorf("无效的字体大小: %d", d.FontSize)
		}
		if d.LineHeight < 0 {
			return fmt.Errorf("无效的行高: %f", d.LineHeight)
		}
		if d.Padding < 0 {
			return fmt.Errorf("无效的内边距: %d", d.Padding)
		}
		if d.BorderRadius < 0 {
			return fmt.Errorf("无效的边框圆角: %d", d.BorderRadius)
		}
	case CoverData:
		if d.Width <= 0 {
			return fmt.Errorf("无效的封面宽度: %d", d.Width)
		}
		if d.Height <= 0 {
			return fmt.Errorf("无效的封面高度: %d", d.Height)
		}
		if d.Title == "" {
			return fmt.Errorf("封面标题不能为空")
		}
	default:
		return fmt.Errorf("未知的模板数据类型")
	}
	return nil
}

// CoverData 封面模板数据
type CoverData struct {
	Width             int
	Height            int
	Background        string
	InnerWidth        int
	InnerHeight       int
	InnerLeft         int
	InnerTop          int
	PaddingTop        int
	PaddingRight      int
	Emoji             string
	EmojiSize         int
	EmojiMarginBottom int
	EmojiColor        string
	Title             string
	TitleSize         int
	TitleGradient     string
	TitleColor        string
	Subtitle          string
	SubtitleSize      int
	SubtitleColor     string
}

// CardData 卡片模板数据
type CardData struct {
	Width      int
	Height     int
	Background string
	Content    string
	ThemeCSS   string
	PageNumber string
	FontSize   int
	LineHeight float64
	Padding    int
	BorderRadius int
}

// calculateTitleSize 根据标题长度计算字体大小
func calculateTitleSize(title string, width int) int {
	titleLen := len([]rune(title))
	switch {
	case titleLen <= 6:
		return int(float64(width) * 0.14) // 极大
	case titleLen <= 10:
		return int(float64(width) * 0.12) // 大
	case titleLen <= 18:
		return int(float64(width) * 0.09) // 中
	case titleLen <= 30:
		return int(float64(width) * 0.07) // 小
	default:
		return int(float64(width) * 0.055) // 极小
	}
}

// GenerateCoverOnly 生成封面图片
func (s *RendererService) GenerateCoverOnly(title, subtitle, styleKey, outputPrefix string, width, height int) (string, error) {
	if width <= 0 {
		width = DefaultWidth
	}
	if height <= 0 {
		height = DefaultHeight
	}

	theme := s.GetStyle(styleKey)

	// 准备模板数据
	data := CoverData{
		Width:             width,
		Height:            height,
		Background:        theme.CoverBg,
		InnerWidth:        int(float64(width) * 0.88),
		InnerHeight:       int(float64(height) * 0.91),
		InnerLeft:         int(float64(width) * 0.06),
		InnerTop:          int(float64(height) * 0.045),
		PaddingTop:        int(float64(width) * 0.074),
		PaddingRight:      int(float64(width) * 0.079),
		Emoji:             "📝",
		EmojiSize:         int(float64(width) * 0.167),
		EmojiMarginBottom: int(float64(height) * 0.035),
		EmojiColor:        "",
		Title:             title,
		TitleSize:         calculateTitleSize(title, width),
		TitleGradient:     theme.TitleGradient,
		TitleColor:        "",
		Subtitle:          subtitle,
		SubtitleSize:      int(float64(width) * 0.067),
		SubtitleColor:     "",
	}

	// 渲染HTML
	htmlContent, err := s.renderTemplate("cover", data)
	if err != nil {
		return "", err
	}

	// 使用Chromium渲染图片
	imagePath, err := s.renderHTMLToImage(htmlContent, outputPrefix, "cover", width, height)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}

// RenderMarkdownToImage 将Markdown渲染为图片
func (s *RendererService) RenderMarkdownToImage(markdown, styleKey, outputPrefix string, mode PaginationMode, width, height, maxHeight int) ([]string, error) {
	if width <= 0 {
		width = DefaultWidth
	}
	if height <= 0 {
		height = DefaultHeight
	}
	if maxHeight <= 0 {
		maxHeight = MaxHeight
	}

	theme := s.GetStyle(styleKey)

	// 解析Markdown内容
	contentParts := s.parseMarkdownContent(markdown, mode)

	var images []string

	for i, part := range contentParts {
		pageNumber := ""
		if len(contentParts) > 1 {
			pageNumber = fmt.Sprintf("%d/%d", i+1, len(contentParts))
		}

		// 转换Markdown为HTML
		htmlContent := s.markdownToHTML(part)

		// 准备模板数据
		data := CardData{
			Width:      width,
			Height:     height,
			Background: theme.CardBg,
			Content:    htmlContent,
			ThemeCSS:   s.loadThemeCSS(styleKey),
			PageNumber: pageNumber,
			FontSize:   42,
			LineHeight: 1.7,
			Padding:    60,
			BorderRadius: 20,
		}

		// 渲染HTML
		htmlContent, err := s.renderTemplate("card", data)
		if err != nil {
			return nil, err
		}

		// 使用Chromium渲染图片
		imagePath, err := s.renderHTMLToImage(htmlContent, outputPrefix, fmt.Sprintf("card_%d", i+1), width, height)
		if err != nil {
			return nil, err
		}

		images = append(images, imagePath)
	}

	return images, nil
}

// parseMarkdownContent 根据模式解析Markdown内容
func (s *RendererService) parseMarkdownContent(markdown string, mode PaginationMode) []string {
	switch mode {
	case PaginationSeparator:
		// 按 --- 分隔符分割
		parts := regexp.MustCompile(`\n---+\n`).Split(markdown, -1)
		var result []string
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
		if len(result) == 0 {
			return []string{markdown}
		}
		return result
	default:
		// 默认返回整个内容
		return []string{strings.TrimSpace(markdown)}
	}
}

// markdownToHTML 将Markdown转换为HTML
func (s *RendererService) markdownToHTML(md string) string {
	// 处理标签 (#开头的内容)
	lines := strings.Split(md, "\n")
	var contentLines []string
	var tags []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			// 检查是否是标签行（只有标签的行）
			if isTagLine(line) {
				// 提取标签
				tagMatches := regexp.MustCompile(`#([\p{Han}\w]+)`).FindAllStringSubmatch(line, -1)
				for _, match := range tagMatches {
					if len(match) > 1 {
						tags = append(tags, match[1])
					}
				}
				continue
			}
		}
		contentLines = append(contentLines, line)
	}

	content := strings.Join(contentLines, "\n")

	// 简单的Markdown转HTML
	html := s.simpleMarkdownToHTML(content)

	// 添加标签
	if len(tags) > 0 {
		tagsHTML := `<div class="tags-container">`
		for _, tag := range tags {
			tagsHTML += fmt.Sprintf(`<span class="tag">#%s</span>`, tag)
		}
		tagsHTML += `</div>`
		html += tagsHTML
	}

	return html
}

// isTagLine 判断是否是纯标签行
func isTagLine(line string) bool {
	// 匹配只有标签的行，如 "#标签1 #标签2"
	// 使用 \p{Han} 匹配中文字符，\w 匹配字母数字下划线
	tagPattern := regexp.MustCompile(`^(#[\p{Han}\w]+\s*)+$`)
	return tagPattern.MatchString(line)
}

// simpleMarkdownToHTML 简单的Markdown转HTML实现
func (s *RendererService) simpleMarkdownToHTML(md string) string {
	// 处理代码块
	md = regexp.MustCompile("```(\\w*)\n([\\s\\S]*?)```").ReplaceAllStringFunc(md, func(match string) string {
		parts := regexp.MustCompile("```(\\w*)\n([\\s\\S]*?)```").FindStringSubmatch(match)
		if len(parts) >= 3 {
			return fmt.Sprintf("<pre><code>%s</code></pre>", s.escapeHTML(parts[2]))
		}
		return match
	})

	// 处理行内代码
	md = regexp.MustCompile("`([^`]+)`").ReplaceAllStringFunc(md, func(match string) string {
		parts := regexp.MustCompile("`([^`]+)`").FindStringSubmatch(match)
		if len(parts) >= 2 {
			return fmt.Sprintf("<code>%s</code>", s.escapeHTML(parts[1]))
		}
		return match
	})

	// 处理标题
	md = regexp.MustCompile(`(?m)^### (.+)$`).ReplaceAllString(md, "<h3>$1</h3>")
	md = regexp.MustCompile(`(?m)^## (.+)$`).ReplaceAllString(md, "<h2>$1</h2>")
	md = regexp.MustCompile(`(?m)^# (.+)$`).ReplaceAllString(md, "<h1>$1</h1>")

	// 处理粗体和斜体
	md = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(md, "<strong>$1</strong>")
	md = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(md, "<em>$1</em>")

	// 处理引用块
	md = regexp.MustCompile(`(?m)^> (.+)$`).ReplaceAllString(md, "<blockquote><p>$1</p></blockquote>")

	// 处理无序列表
	md = regexp.MustCompile(`(?m)^[-*] (.+)$`).ReplaceAllString(md, "<li>$1</li>")
	md = regexp.MustCompile(`(<li>.*</li>\n?)+`).ReplaceAllStringFunc(md, func(match string) string {
		return "<ul>" + match + "</ul>"
	})

	// 处理有序列表
	md = regexp.MustCompile(`(?m)^\d+\. (.+)$`).ReplaceAllString(md, "<li>$1</li>")

	// 处理分割线
	md = regexp.MustCompile(`(?m)^---+$`).ReplaceAllString(md, "<hr>")

	// 处理段落
	lines := strings.Split(md, "\n")
	var result []string
	var currentParagraph strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if currentParagraph.Len() > 0 {
				para := currentParagraph.String()
				if !strings.HasPrefix(para, "<h") && !strings.HasPrefix(para, "<ul") &&
					!strings.HasPrefix(para, "<ol") && !strings.HasPrefix(para, "<blockquote") &&
					!strings.HasPrefix(para, "<pre") && !strings.HasPrefix(para, "<hr") {
					para = fmt.Sprintf("<p>%s</p>", para)
				}
				result = append(result, para)
				currentParagraph.Reset()
			}
		} else {
			if currentParagraph.Len() > 0 {
				currentParagraph.WriteString(" ")
			}
			currentParagraph.WriteString(line)
		}
	}

	if currentParagraph.Len() > 0 {
		para := currentParagraph.String()
		if !strings.HasPrefix(para, "<h") && !strings.HasPrefix(para, "<ul") &&
			!strings.HasPrefix(para, "<ol") && !strings.HasPrefix(para, "<blockquote") &&
			!strings.HasPrefix(para, "<pre") && !strings.HasPrefix(para, "<hr") {
			para = fmt.Sprintf("<p>%s</p>", para)
		}
		result = append(result, para)
	}

	return strings.Join(result, "\n")
}

// escapeHTML 转义HTML特殊字符
func (s *RendererService) escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&#39;")
	return text
}

// renderHTMLToImage 使用Chromium渲染HTML为图片
func (s *RendererService) renderHTMLToImage(htmlContent, outputPrefix, suffix string, width, height int) (string, error) {
	// 指定 Chrome 浏览器路径（macOS）
	chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	
	// 检查 Chrome 是否存在
	if _, err := os.Stat(chromePath); os.IsNotExist(err) {
		// 尝试其他可能的路径
		chromePath = "/Applications/Chromium.app/Contents/MacOS/Chromium"
	}
	
	// 创建 chromedp 执行分配器，指定 Chrome 路径
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(chromePath),
			chromedp.Headless,
			chromedp.NoSandbox,
			chromedp.DisableGPU,
			chromedp.WindowSize(width, height),
		)...,
	)
	defer cancelAlloc()

	// 创建 chromedp 上下文
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 设置超时
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte

	// 创建临时HTML文件
	tempFile := filepath.Join(s.imagesDir, ".temp_render.html")
	if err := os.WriteFile(tempFile, []byte(htmlContent), 0644); err != nil {
		return "", fmt.Errorf("创建临时HTML文件失败: %v", err)
	}
	defer os.Remove(tempFile)

	// 执行渲染
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(int64(width), int64(height), chromedp.EmulateScale(2)),
		chromedp.Navigate("file://"+tempFile),
		chromedp.WaitReady("body"),
		chromedp.Sleep(500*time.Millisecond), // 等待字体加载
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		return "", fmt.Errorf("Chromium渲染失败: %v", err)
	}

	// 保存图片
	filename := s.generateFilename(outputPrefix, suffix)
	fullPath := filepath.Join(s.imagesDir, filename)

	if err := os.WriteFile(fullPath, buf, 0644); err != nil {
		return "", fmt.Errorf("保存图片失败: %v", err)
	}

	return "/xiaohongshu-renderer/image/" + filename, nil
}

// generateFilename 生成文件名
func (s *RendererService) generateFilename(prefix, suffix string) string {
	if prefix == "" {
		prefix = "note"
	}

	timestamp := time.Now().Format("20060102150405")

	if suffix != "" {
		return fmt.Sprintf("%s_%s_%s.png", prefix, timestamp, suffix)
	}
	return fmt.Sprintf("%s_%s.png", prefix, timestamp)
}

// generateRandomName 生成随机文件名
func generateRandomName(prefix, ext string) string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%s_%s.%s", prefix, hex.EncodeToString(bytes), ext)
}
