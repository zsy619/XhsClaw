// Package service 提供业务逻辑层 - 小红书图片渲染服务
package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	DefaultWidth     = 1080
	DefaultHeight    = 1440
	PaddingLeft      = 60
	PaddingRight     = 60
	PaddingTop       = 120
	PaddingBottom    = 120
	LineHeight       = 55
	TitleSize        = 64
	BodySize         = 42
	ListIndent       = 70
	ParagraphSpacing = 40 // 段落间距
	TitleSpacing     = 80 // 标题后的间距
	ListItemSpacing  = 20 // 列表项间距
)

type PaginationMode string

const (
	PaginationAutoSplit = "auto-split"
	PaginationSeparator = "separator"
	PaginationAutoFit   = "auto-fit"
	PaginationDynamic   = "dynamic"
)

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
}

type RendererService struct {
	styles    map[string]StyleConfig
	font      font.Face
	fontBold  font.Face
	fontData  *truetype.Font
	emojiFont *truetype.Font
	imagesDir string
}

// GetImagesDir 获取图片存储目录
func (s *RendererService) GetImagesDir() string {
	return s.imagesDir
}

type TextBlock struct {
	Text     string
	IsTitle  bool
	IsList   bool
	ListChar string
}

// 系统字体路径 - 按优先级排序
var systemFontPaths = []string{
	"/System/Library/Fonts/PingFang.ttc",
	"/System/Library/Fonts/PingFangSC-Regular.otf",
	"/System/Library/Fonts/Hiragino Sans GB.ttc",
	"/System/Library/Fonts/STHeiti Light.ttc",
	"/System/Library/Fonts/STHeiti Medium.ttc",
	"/Library/Fonts/STHeiti Light.ttc",
	"/Library/Fonts/STHeiti Medium.ttc",
	"C:/Windows/Fonts/msyh.ttc",
	"C:/Windows/Fonts/simhei.ttf",
	"C:/Windows/Fonts/simsun.ttc",
	"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
}

// Emoji 字体路径
var emojiFontPaths = []string{
	"/System/Library/Fonts/Apple Color Emoji.ttc",
	"/System/Library/Fonts/PingFang.ttc",
	"C:/Windows/Fonts/seguiemj.ttf",
}

func loadSystemFont() (*truetype.Font, font.Face, font.Face) {
	// 尝试加载系统字体
	for _, fontPath := range systemFontPaths {
		fontData, err := ioutil.ReadFile(fontPath)
		if err != nil {
			continue
		}

		f, err := freetype.ParseFont(fontData)
		if err != nil {
			fmt.Printf("解析字体失败 %s: %v\n", fontPath, err)
			continue
		}

		// 创建常规字体
		regularFace := truetype.NewFace(f, &truetype.Options{
			Size: float64(BodySize),
			DPI:  72,
		})

		// 创建粗体字体 (复用 regularFace，因为 truetype.Options 没有 Weight)
		boldFace := regularFace

		fmt.Printf("成功加载系统字体: %s\n", fontPath)
		return f, regularFace, boldFace
	}

	// 降级使用基本字体
	fmt.Println("无法加载系统字体，使用 basicfont")
	return nil, basicfont.Face7x13, basicfont.Face7x13
}

func loadEmojiFont() *truetype.Font {
	for _, fontPath := range emojiFontPaths {
		fontData, err := ioutil.ReadFile(fontPath)
		if err != nil {
			continue
		}

		f, err := freetype.ParseFont(fontData)
		if err != nil {
			fmt.Printf("解析Emoji字体失败 %s: %v\n", fontPath, err)
			continue
		}

		fmt.Printf("成功加载Emoji字体: %s\n", fontPath)
		return f
	}

	fmt.Println("无法加载Emoji字体")
	return nil
}

func getProjectRoot() string {
	return "/Volumes/E/JYW/创意项目/XhsClaw/backend"
}

func NewRendererService() *RendererService {
	imagesDir := filepath.Join(getProjectRoot(), "public", "images")
	os.MkdirAll(imagesDir, 0755)

	fontData, fontFace, fontBold := loadSystemFont()
	emojiFontData := loadEmojiFont()

	return &RendererService{
		font:      fontFace,
		fontBold:  fontBold,
		fontData:  fontData,
		emojiFont: emojiFontData,
		imagesDir: imagesDir,
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
			},
			"playful-geometric": {
				Key:           "playful-geometric",
				Name:          "活泼几何",
				Primary:       color.RGBA{R: 59, G: 130, B: 246, A: 255},
				Secondary:     color.RGBA{R: 96, G: 165, B: 250, A: 255},
				Background:    color.RGBA{R: 239, G: 246, B: 255, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 102, G: 102, B: 102, A: 255},
				Accent:        color.RGBA{R: 59, G: 130, B: 246, A: 255},
			},
			"neo-brutalism": {
				Key:           "neo-brutalism",
				Name:          "新野兽派",
				Primary:       color.RGBA{R: 251, G: 146, B: 60, A: 255},
				Secondary:     color.RGBA{R: 255, G: 183, B: 77, A: 255},
				Background:    color.RGBA{R: 255, G: 250, B: 240, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 102, G: 102, B: 102, A: 255},
				Accent:        color.RGBA{R: 251, G: 146, B: 60, A: 255},
			},
			"botanical": {
				Key:           "botanical",
				Name:          "植物系",
				Primary:       color.RGBA{R: 72, G: 187, B: 120, A: 255},
				Secondary:     color.RGBA{R: 129, G: 230, B: 176, A: 255},
				Background:    color.RGBA{R: 240, G: 253, B: 244, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 102, G: 102, B: 102, A: 255},
				Accent:        color.RGBA{R: 72, G: 187, B: 120, A: 255},
			},
			"professional": {
				Key:           "professional",
				Name:          "专业商务",
				Primary:       color.RGBA{R: 107, G: 114, B: 128, A: 255},
				Secondary:     color.RGBA{R: 148, G: 163, B: 184, A: 255},
				Background:    color.RGBA{R: 245, G: 247, B: 250, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 31, G: 41, B: 55, A: 255},
				TextSecondary: color.RGBA{R: 107, G: 114, B: 128, A: 255},
				Accent:        color.RGBA{R: 107, G: 114, B: 128, A: 255},
			},
			"retro": {
				Key:           "retro",
				Name:          "复古风格",
				Primary:       color.RGBA{R: 147, G: 112, B: 219, A: 255},
				Secondary:     color.RGBA{R: 187, G: 154, B: 247, A: 255},
				Background:    color.RGBA{R: 250, G: 248, B: 255, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 102, G: 102, B: 102, A: 255},
				Accent:        color.RGBA{R: 147, G: 112, B: 219, A: 255},
			},
			"terminal": {
				Key:           "terminal",
				Name:          "终端风格",
				Primary:       color.RGBA{R: 249, G: 115, B: 22, A: 255},
				Secondary:     color.RGBA{R: 251, G: 146, B: 60, A: 255},
				Background:    color.RGBA{R: 17, G: 24, B: 39, A: 255},
				CardInner:     color.RGBA{R: 31, G: 41, B: 55, A: 255},
				TextPrimary:   color.RGBA{R: 249, G: 250, B: 251, A: 255},
				TextSecondary: color.RGBA{R: 148, G: 163, B: 184, A: 255},
				Accent:        color.RGBA{R: 249, G: 115, B: 22, A: 255},
			},
			"sketch": {
				Key:           "sketch",
				Name:          "手绘风格",
				Primary:       color.RGBA{R: 236, G: 72, B: 153, A: 255},
				Secondary:     color.RGBA{R: 244, G: 114, B: 182, A: 255},
				Background:    color.RGBA{R: 255, G: 245, B: 247, A: 255},
				CardInner:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
				TextPrimary:   color.RGBA{R: 51, G: 51, B: 51, A: 255},
				TextSecondary: color.RGBA{R: 102, G: 102, B: 102, A: 255},
				Accent:        color.RGBA{R: 236, G: 72, B: 153, A: 255},
			},
		},
	}
}

func (s *RendererService) GetStyles() []StyleConfig {
	styles := make([]StyleConfig, 0, len(s.styles))
	for _, style := range s.styles {
		styles = append(styles, style)
	}
	return styles
}

func (s *RendererService) GetStyle(key string) StyleConfig {
	style, exists := s.styles[key]
	if !exists {
		style = s.styles["default"]
	}
	return style
}

func (s *RendererService) RenderMarkdownToImage(markdown, styleKey, outputPrefix string, mode PaginationMode, width, height, maxHeight int) ([]string, error) {
	if width <= 0 {
		width = DefaultWidth
	}
	if height <= 0 {
		height = DefaultHeight
	}
	if maxHeight <= 0 {
		maxHeight = height
	}

	style := s.GetStyle(styleKey)
	blocks := parseMarkdown(markdown)

	switch mode {
	case PaginationSeparator:
		return s.renderWithSeparator(blocks, style, outputPrefix, width, height)
	case PaginationAutoFit:
		return s.renderAutoFit(blocks, style, outputPrefix, width, height)
	case PaginationDynamic:
		return s.renderDynamic(blocks, style, outputPrefix, width, height, maxHeight)
	default:
		return s.renderAutoSplit(blocks, style, outputPrefix, width, height, maxHeight)
	}
}

func parseMarkdown(markdown string) []TextBlock {
	lines := strings.Split(markdown, "\n")
	var blocks []TextBlock

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "---") {
			blocks = append(blocks, TextBlock{Text: "---", IsList: false})
			continue
		}

		if strings.HasPrefix(line, "#") {
			title := strings.TrimPrefix(line, "#")
			title = strings.TrimSpace(title)
			blocks = append(blocks, TextBlock{Text: title, IsTitle: true})
			continue
		}

		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "• ") {
			text := strings.TrimPrefix(line, "- ")
			text = strings.TrimPrefix(text, "• ")
			blocks = append(blocks, TextBlock{Text: text, IsList: true, ListChar: "•"})
			continue
		}

		if len(line) > 0 && (line[0] >= '0' && line[0] <= '9') && strings.Contains(line, ".") {
			parts := strings.SplitN(line, ".", 2)
			if len(parts) >= 2 {
				blocks = append(blocks, TextBlock{Text: strings.TrimSpace(parts[1]), IsList: true, ListChar: parts[0] + "."})
				continue
			}
		}

		blocks = append(blocks, TextBlock{Text: line, IsTitle: false, IsList: false})
	}

	return blocks
}

func (s *RendererService) renderWithSeparator(blocks []TextBlock, style StyleConfig, outputPrefix string, width, height int) ([]string, error) {
	var images []string
	var currentPageBlocks []TextBlock
	pageIndex := 0

	addPage := func() error {
		if len(currentPageBlocks) == 0 {
			return nil
		}
		img, err := s.renderPage(currentPageBlocks, style, width, height)
		if err != nil {
			return err
		}
		filename := s.saveImage(img, outputPrefix, pageIndex)
		images = append(images, filename)
		pageIndex++
		currentPageBlocks = nil
		return nil
	}

	for _, block := range blocks {
		if block.Text == "---" {
			if err := addPage(); err != nil {
				return nil, err
			}
			continue
		}
		currentPageBlocks = append(currentPageBlocks, block)
	}

	if err := addPage(); err != nil {
		return nil, err
	}

	if len(images) == 0 {
		img, err := s.renderPage(blocks, style, width, height)
		if err != nil {
			return nil, err
		}
		filename := s.saveImage(img, outputPrefix, 0)
		images = append(images, filename)
	}

	return images, nil
}

func (s *RendererService) renderAutoSplit(blocks []TextBlock, style StyleConfig, outputPrefix string, width, height, maxContentHeight int) ([]string, error) {
	var images []string
	var currentPageBlocks []TextBlock
	currentHeight := 0
	pageIndex := 0

	for _, block := range blocks {
		blockHeight := estimateBlockHeight(block, width)

		if currentHeight+blockHeight > maxContentHeight && len(currentPageBlocks) > 0 {
			img, err := s.renderPage(currentPageBlocks, style, width, height)
			if err != nil {
				return nil, err
			}
			filename := s.saveImage(img, outputPrefix, pageIndex)
			images = append(images, filename)
			pageIndex++
			currentPageBlocks = nil
			currentHeight = 0
		}

		currentPageBlocks = append(currentPageBlocks, block)
		currentHeight += blockHeight
	}

	if len(currentPageBlocks) > 0 {
		img, err := s.renderPage(currentPageBlocks, style, width, height)
		if err != nil {
			return nil, err
		}
		filename := s.saveImage(img, outputPrefix, pageIndex)
		images = append(images, filename)
	}

	if len(images) == 0 {
		img, err := s.renderPage(blocks, style, width, height)
		if err != nil {
			return nil, err
		}
		filename := s.saveImage(img, outputPrefix, 0)
		images = append(images, filename)
	}

	return images, nil
}

func (s *RendererService) renderAutoFit(blocks []TextBlock, style StyleConfig, outputPrefix string, width, height int) ([]string, error) {
	img, err := s.renderPage(blocks, style, width, height)
	if err != nil {
		return nil, err
	}
	filename := s.saveImage(img, outputPrefix, 0)
	return []string{filename}, nil
}

func (s *RendererService) renderDynamic(blocks []TextBlock, style StyleConfig, outputPrefix string, width, height, maxHeight int) ([]string, error) {
	contentHeight := calculateContentHeight(blocks, width)
	if contentHeight > maxHeight {
		return s.renderAutoSplit(blocks, style, outputPrefix, width, height, maxHeight-200)
	}

	img := s.createImage(width, contentHeight+PaddingTop+PaddingBottom)
	draw.Draw(img, img.Bounds(), &image.Uniform{style.Background}, image.Point{}, draw.Src)

	s.drawCard(img, 0, 0, width, contentHeight+PaddingTop+PaddingBottom, style)
	s.drawTextBlocks(img, blocks, style, width, contentHeight+PaddingTop+PaddingBottom)

	filename := s.saveImage(img, outputPrefix, 0)
	return []string{filename}, nil
}

func estimateBlockHeight(block TextBlock, width int) int {
	maxTextWidth := width - PaddingLeft - PaddingRight
	if block.IsList {
		maxTextWidth -= ListIndent
	}
	charsPerLine := maxTextWidth / (BodySize / 2)
	if charsPerLine == 0 {
		charsPerLine = 20
	}
	lines := (len([]rune(block.Text)) + charsPerLine - 1) / charsPerLine
	if lines == 0 {
		lines = 1
	}
	if block.IsTitle {
		return TitleSize + LineHeight*2
	}
	return (BodySize + LineHeight) * lines
}

func calculateContentHeight(blocks []TextBlock, width int) int {
	height := PaddingTop
	for _, block := range blocks {
		if block.IsTitle {
			height += TitleSize + LineHeight*2
		} else {
			h := estimateBlockHeight(block, width)
			height += h
		}
	}
	height += PaddingBottom
	return height
}

func (s *RendererService) renderPage(blocks []TextBlock, style StyleConfig, width, height int) (*image.RGBA, error) {
	img := s.createImage(width, height)
	draw.Draw(img, img.Bounds(), &image.Uniform{style.Background}, image.Point{}, draw.Src)

	cardY := 50
	cardHeight := height - 100
	s.drawCard(img, 30, cardY, width-60, cardHeight, style)

	s.drawTextBlocks(img, blocks, style, width, cardHeight)

	return img, nil
}

func (s *RendererService) createImage(width, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

func (s *RendererService) drawCard(img *image.RGBA, x, y, width, height int, style StyleConfig) {
	cardRect := image.Rect(x, y, x+width, y+height)
	draw.Draw(img, cardRect, &image.Uniform{style.Background}, image.Point{}, draw.Src)

	innerPadding := 15
	innerRect := image.Rect(x+innerPadding, y+innerPadding, x+width-innerPadding, y+height-innerPadding)
	draw.Draw(img, innerRect, &image.Uniform{style.CardInner}, image.Point{}, draw.Src)
}

func (s *RendererService) drawTextBlocks(img *image.RGBA, blocks []TextBlock, style StyleConfig, width, height int) {
	currentY := PaddingTop + 40
	prevWasParagraph := false

	for _, block := range blocks {
		if block.IsTitle {
			s.drawCenteredText(img, block.Text, width/2, currentY+40, TitleSize, style.TextPrimary)
			currentY += TitleSize + TitleSpacing
		} else if block.IsList {
			// 列表项之间添加间距
			if prevWasParagraph {
				currentY += ParagraphSpacing / 2
			}
			x := PaddingLeft + ListIndent
			s.drawText(img, block.ListChar+" ", x-25, currentY, BodySize, style.Accent)
			s.drawText(img, block.Text, x, currentY, BodySize, style.TextPrimary)
			currentY += BodySize + ListItemSpacing
		} else {
			// 段落之间添加更大的间距
			if prevWasParagraph {
				currentY += ParagraphSpacing
			}
			lines := s.wrapText(block.Text, width-PaddingLeft-PaddingRight)
			for _, line := range lines {
				s.drawText(img, line, PaddingLeft, currentY, BodySize, style.TextPrimary)
				currentY += BodySize + LineHeight/2
			}
		}
		prevWasParagraph = !block.IsTitle && !block.IsList
	}
}

func (s *RendererService) drawText(img *image.RGBA, text string, x, y, size int, col color.Color) {
	// 检测文本是否包含 emoji
	hasEmoji := s.containsEmoji(text)

	if hasEmoji {
		// 分段绘制，分别使用不同的字体
		s.drawTextWithEmoji(img, text, x, y, size, col)
		return
	}

	// 使用 freetype 的 DrawString 来绘制文本
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(s.fontData)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(&image.Uniform{col})
	c.SetFontSize(float64(size))

	point := freetype.Pt(x, y+size)
	c.DrawString(text, point)
}

func (s *RendererService) containsEmoji(text string) bool {
	if s.emojiFont == nil {
		return false
	}
	for _, r := range []rune(text) {
		if r >= 0x1F300 && r <= 0x1F9FF { // Emoji range
			return true
		}
	}
	return false
}

func (s *RendererService) drawTextWithEmoji(img *image.RGBA, text string, x, y, size int, col color.Color) {
	if s.emojiFont == nil {
		// 没有 emoji 字体，使用普通字体绘制
		s.drawText(img, text, x, y, size, col)
		return
	}

	runes := []rune(text)
	currentX := x
	emojiFont := s.emojiFont

	for _, r := range runes {
		var fontToUse *truetype.Font
		if r >= 0x1F300 && r <= 0x1F9FF {
			fontToUse = emojiFont
		} else {
			fontToUse = s.fontData
		}

		c := freetype.NewContext()
		c.SetDPI(72)
		c.SetFont(fontToUse)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(&image.Uniform{col})
		c.SetFontSize(float64(size))

		point := freetype.Pt(currentX, y+size)
		_, err := c.DrawString(string(r), point)
		if err == nil {
			// 计算这个字符的宽度
			bounds, advance, _ := s.font.GlyphBounds(r)
			if bounds != (fixed.Rectangle26_6{}) {
				currentX += int(advance >> 6)
			}
		}
	}
}

func (s *RendererService) drawCenteredText(img *image.RGBA, text string, x, y, size int, col color.Color) {
	// 使用 freetype 的 DrawString 来绘制文本
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(s.fontData)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(&image.Uniform{col})
	c.SetFontSize(float64(size))

	// 使用 font.Face 来计算文本宽度
	runes := []rune(text)
	var totalWidth fixed.Int26_6
	for _, r := range runes {
		_, advance, _ := s.font.GlyphBounds(r)
		totalWidth += advance
	}
	startX := x - int(totalWidth)/2

	point := freetype.Pt(startX, y+size)
	c.DrawString(text, point)
}

func (s *RendererService) wrapText(text string, maxWidth int) []string {
	var lines []string
	runes := []rune(text)

	// 估算每个字符的宽度（根据字体大小）
	// 中文字符宽度约等于字体大小，英文约等于字体大小的一半
	charWidth := BodySize / 2
	if charWidth == 0 {
		charWidth = 21
	}
	maxChars := maxWidth / charWidth
	if maxChars == 0 {
		maxChars = 20
	}

	var currentLine string
	for _, r := range runes {
		// 检测是否是中文或表情符号（宽字符）
		isWideChar := r >= 0x4E00 && r <= 0x9FFF || // 中文
			r >= 0x3000 && r <= 0x303F || // 中文标点
			r >= 0xFF00 && r <= 0xFFEF || // 全角字符
			r >= 0x1F300 && r <= 0x1F9FF // Emoji

		// 宽字符占2个位置
		charLen := 1
		if isWideChar {
			charLen = 2
		}

		// 检查是否需要换行
		currentRunes := []rune(currentLine)
		currentLen := 0
		for _, cr := range currentRunes {
			isWide := cr >= 0x4E00 && cr <= 0x9FFF || cr >= 0x3000 && cr <= 0x303F || cr >= 0xFF00 && cr <= 0xFFEF || cr >= 0x1F300 && cr <= 0x1F9FF
			if isWide {
				currentLen += 2
			} else {
				currentLen += 1
			}
		}

		if currentLen+charLen > maxChars {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = string(r)
		} else {
			currentLine += string(r)
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func (s *RendererService) saveImage(img *image.RGBA, prefix string, index int) string {
	if prefix == "" {
		prefix = "note"
	}

	timestamp := time.Now().Format("20060102150405")
	var filename string
	if index == 0 {
		filename = fmt.Sprintf("%s_%s.png", prefix, timestamp)
	} else {
		filename = fmt.Sprintf("%s_%s_%d.png", prefix, timestamp, index)
	}

	fullPath := filepath.Join(s.imagesDir, filename)

	f, err := os.Create(fullPath)
	if err != nil {
		randomName := generateRandomName(prefix, "png")
		fullPath = filepath.Join(s.imagesDir, randomName)
		f, err = os.Create(fullPath)
	}
	if err != nil {
		return ""
	}
	defer f.Close()

	png.Encode(f, img)

	return "/xiaohongshu-renderer/image/" + filename
}

func generateRandomName(prefix, ext string) string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%s_%s.%s", prefix, hex.EncodeToString(bytes), ext)
}

func (s *RendererService) GenerateCoverOnly(title, subtitle, styleKey, outputPrefix string, width, height int) (string, error) {
	if width <= 0 {
		width = DefaultWidth
	}
	if height <= 0 {
		height = DefaultHeight
	}

	style := s.GetStyle(styleKey)

	img := s.createImage(width, height)
	draw.Draw(img, img.Bounds(), &image.Uniform{style.Background}, image.Point{}, draw.Src)

	cardPadding := 50
	cardRect := image.Rect(cardPadding, cardPadding, width-cardPadding, height-cardPadding)
	draw.Draw(img, cardRect, &image.Uniform{style.Primary}, image.Point{}, draw.Src)

	innerPadding := 25
	innerRect := image.Rect(
		cardPadding+innerPadding,
		cardPadding+innerPadding,
		width-cardPadding-innerPadding,
		height-cardPadding-innerPadding,
	)
	draw.Draw(img, innerRect, &image.Uniform{style.CardInner}, image.Point{}, draw.Src)

	titleY := height / 2
	if subtitle != "" {
		titleY = height/2 - 50
		s.drawCenteredText(img, subtitle, width/2, titleY+80, BodySize, style.TextSecondary)
	}

	s.drawCenteredText(img, title, width/2, titleY, TitleSize+12, style.Primary)

	filename := s.saveCoverImage(img, outputPrefix)
	return filename, nil
}

func (s *RendererService) saveCoverImage(img *image.RGBA, prefix string) string {
	if prefix == "" {
		prefix = "cover"
	}

	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.png", prefix, timestamp)

	fullPath := filepath.Join(s.imagesDir, filename)

	f, err := os.Create(fullPath)
	if err != nil {
		randomName := generateRandomName(prefix, "png")
		fullPath = filepath.Join(s.imagesDir, randomName)
		f, err = os.Create(fullPath)
	}
	if err != nil {
		return ""
	}
	defer f.Close()

	png.Encode(f, img)

	return "/xiaohongshu-renderer/image/" + filename
}
