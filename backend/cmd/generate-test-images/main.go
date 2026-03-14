// Package main 测试图片生成程序
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// 颜色配置
var themeColors = map[string]struct {
	bg   color.Color
	text color.Color
}{
	"playful-geometric": {
		bg:   color.RGBA{R: 255, G: 250, B: 251, A: 255},
		text: color.RGBA{R: 51, G: 51, B: 51, A: 255},
	},
	"retro": {
		bg:   color.RGBA{R: 250, G: 248, B: 255, A: 255},
		text: color.RGBA{R: 51, G: 51, B: 51, A: 255},
	},
	"Sketch": {
		bg:   color.RGBA{R: 240, G: 253, B: 244, A: 255},
		text: color.RGBA{R: 51, G: 51, B: 51, A: 255},
	},
	"terminal": {
		bg:   color.RGBA{R: 239, G: 246, B: 255, A: 255},
		text: color.RGBA{R: 51, G: 51, B: 51, A: 255},
	},
	"auto-fit": {
		bg:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
		text: color.RGBA{R: 51, G: 51, B: 51, A: 255},
	},
}

// drawRect 绘制矩形
func drawRect(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			img.Set(x, y, c)
		}
	}
}

// drawTextLine 绘制简单的文本行（使用方块模拟文字）
func drawTextLine(img *image.RGBA, x, y, width int, c color.Color) {
	for i := 0; i < width; i++ {
		if i%8 < 6 {
			img.Set(x+i, y, c)
			img.Set(x+i, y+1, c)
			img.Set(x+i, y+2, c)
			img.Set(x+i, y+3, c)
			img.Set(x+i, y+4, c)
		}
		if i%8 == 7 {
			x += 2
		}
	}
}

// generateTestImage 生成测试图片
func generateTestImage(filename string, theme string, cardNum int) error {
	width := 1080
	height := 1440

	// 创建图片
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 获取主题颜色
	colors, ok := themeColors[theme]
	if !ok {
		colors = themeColors["playful-geometric"]
	}

	// 填充背景
	drawRect(img, 0, 0, width, height, colors.bg)

	// 绘制顶部装饰条
	decorColor := color.RGBA{R: 255, G: 66, B: 99, A: 255}
	drawRect(img, 0, 0, width, 80, decorColor)

	// 绘制标题区域
	titleBg := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	drawRect(img, 40, 100, width-40, 200, titleBg)

	// 绘制标题文本（模拟）
	for i := 0; i < 3; i++ {
		drawTextLine(img, 60, 120+i*20, 300, colors.text)
	}

	// 绘制内容区域
	contentBg := color.RGBA{R: 255, G: 255, B: 255, A: 200}
	drawRect(img, 40, 220, width-40, 1200, contentBg)

	// 绘制内容行
	for i := 0; i < 20; i++ {
		lineWidth := 200 + (i%5)*50
		drawTextLine(img, 60, 250+i*40, lineWidth, colors.text)
	}

	// 添加卡片编号装饰
	drawRect(img, width-150, height-100, width-50, height-50, decorColor)

	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 保存图片
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func main() {
	fmt.Println("开始生成测试图片...")

	themes := []string{"playful-geometric", "retro", "Sketch", "terminal", "auto-fit"}
	basePath := "./public/images/all_themes"

	for _, theme := range themes {
		fmt.Printf("正在生成主题: %s\n", theme)

		// 生成5张卡片
		for i := 1; i <= 5; i++ {
			filename := filepath.Join(basePath, theme, fmt.Sprintf("card_%d.png", i))
			if err := generateTestImage(filename, theme, i); err != nil {
				fmt.Printf("  生成卡片 %d 失败: %v\n", i, err)
			} else {
				fmt.Printf("  ✓ 卡片 %d 生成成功\n", i)
			}

			// 有些主题只需要1张卡片
			if theme == "auto-fit" && i == 1 {
				break
			}
		}

		// 生成封面
		coverFilename := filepath.Join(basePath, theme, "cover.png")
		if err := generateTestImage(coverFilename, theme, 0); err != nil {
			fmt.Printf("  生成封面失败: %v\n", err)
		} else {
			fmt.Println("  ✓ 封面生成成功")
		}
	}

	fmt.Println("✅ 所有测试图片生成完成！")
}
