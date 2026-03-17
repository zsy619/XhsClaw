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

// RendererService 渲染服务
type RendererService struct {
	themes    map[string]ThemeConfig
	templates map[string]*template.Template
	imagesDir string
	assetsDir string
}

// NewRendererService 创建渲染服务实例
func NewRendererService() *RendererService {
	// 获取项目根目录
	projectRoot := getProjectRoot()
	assetsDir := filepath.Join(projectRoot, "assets")
	imagesDir := filepath.Join(projectRoot, "public", "images")

	// 确保目录存在
	os.MkdirAll(imagesDir, 0755)

	s := &RendererService{
		imagesDir: imagesDir,
		assetsDir: assetsDir,
		themes:    make(map[string]ThemeConfig),
		templates: make(map[string]*template.Template),
	}

	// 初始化主题配置
	s.initThemes()

	// 加载模板
	s.loadTemplates()

	return s
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
	}
}

// loadTemplates 加载HTML模板
func (s *RendererService) loadTemplates() {
	templatesDir := filepath.Join(s.assetsDir, "templates")

	// 加载封面模板
	coverPath := filepath.Join(templatesDir, "cover.html")
	if content, err := os.ReadFile(coverPath); err == nil {
		s.templates["cover"] = template.Must(template.New("cover").Parse(string(content)))
	}

	// 加载卡片模板
	cardPath := filepath.Join(templatesDir, "card.html")
	if content, err := os.ReadFile(cardPath); err == nil {
		s.templates["card"] = template.Must(template.New("card").Parse(string(content)))
	}
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
	Title             string
	TitleSize         int
	TitleGradient     string
	Subtitle          string
	SubtitleSize      int
}

// CardData 卡片模板数据
type CardData struct {
	Width      int
	Height     int
	Background string
	Content    string
	ThemeCSS   string
	PageNumber string
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
		Title:             title,
		TitleSize:         calculateTitleSize(title, width),
		TitleGradient:     theme.TitleGradient,
		Subtitle:          subtitle,
		SubtitleSize:      int(float64(width) * 0.067),
	}

	// 渲染HTML
	var htmlBuf bytes.Buffer
	if tmpl, ok := s.templates["cover"]; ok {
		if err := tmpl.Execute(&htmlBuf, data); err != nil {
			return "", fmt.Errorf("渲染封面HTML失败: %v", err)
		}
	} else {
		return "", fmt.Errorf("封面模板未加载")
	}

	// 使用Chromium渲染图片
	imagePath, err := s.renderHTMLToImage(htmlBuf.String(), outputPrefix, "cover", width, height)
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
		}

		// 渲染HTML
		var htmlBuf bytes.Buffer
		if tmpl, ok := s.templates["card"]; ok {
			if err := tmpl.Execute(&htmlBuf, data); err != nil {
				return nil, fmt.Errorf("渲染卡片HTML失败: %v", err)
			}
		} else {
			return nil, fmt.Errorf("卡片模板未加载")
		}

		// 使用Chromium渲染图片
		imagePath, err := s.renderHTMLToImage(htmlBuf.String(), outputPrefix, fmt.Sprintf("card_%d", i+1), width, height)
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
	// 创建chromedp上下文
	ctx, cancel := chromedp.NewContext(context.Background())
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
