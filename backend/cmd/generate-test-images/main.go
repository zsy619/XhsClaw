// Package main 完整图片生成测试程序
package main

import (
	"fmt"
	"log"

	"xiaohongshu/internal/service"
)

func main() {
	// 创建渲染服务
	renderer, err := service.NewRendererService()
	if err != nil {
		log.Fatalf("初始化渲染服务失败: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("小红书图片生成测试")
	fmt.Println("========================================")

	// 测试1: 生成封面
	fmt.Println("\n📝 测试1: 生成封面图片...")
	coverPath, err := renderer.GenerateCoverOnly(
		"小红书运营秘籍",
		"让你快速涨粉的10个技巧",
		"playful-geometric",
		"test_cover",
		1080,
		1440,
	)
	if err != nil {
		log.Fatalf("生成封面失败: %v", err)
	}
	fmt.Printf("✅ 封面生成成功: %s\n", coverPath)

	// 测试2: 测试不同主题
	themes := []string{"default", "playful-geometric", "neo-brutalism", "terminal", "sketch"}
	
	for _, theme := range themes {
		fmt.Printf("\n📝 测试主题: %s\n", theme)
		
		// 测试Markdown内容
		markdownContent := fmt.Sprintf(`# %s主题测试

这是一个使用 **%s** 主题生成的测试卡片。

## 功能特点

- 支持 Markdown 语法
- 支持多种主题样式
- 支持自动分页

> 这是一条引用文字

代码示例: ` + "`print('Hello')`" + `

#测试 #%s主题`, theme, theme, theme)

		images, err := renderer.RenderMarkdownToImage(
			markdownContent,
			theme,
			fmt.Sprintf("test_%s", theme),
			service.PaginationSeparator,
			1080,
			1440,
			4320,
		)
		if err != nil {
			log.Printf("❌ 主题 %s 生成失败: %v", theme, err)
			continue
		}
		fmt.Printf("✅ 主题 %s 生成成功，共 %d 张图片\n", theme, len(images))
		for i, img := range images {
			fmt.Printf("   %d. %s\n", i+1, img)
		}
	}

	// 测试3: 测试分页功能
	fmt.Println("\n📝 测试3: 分页功能测试...")
	multiPageContent := `# 第一页内容

这是第一页的内容。

---

# 第二页内容

这是第二页的内容。

---

# 第三页内容

这是第三页的内容。

#分页测试`

	images, err := renderer.RenderMarkdownToImage(
		multiPageContent,
		"default",
		"test_multipage",
		service.PaginationSeparator,
		1080,
		1440,
		4320,
	)
	if err != nil {
		log.Fatalf("分页测试失败: %v", err)
	}
	fmt.Printf("✅ 分页测试成功，共生成 %d 张图片\n", len(images))
	for i, img := range images {
		fmt.Printf("   %d. %s\n", i+1, img)
	}

	// 测试4: 测试标题长度自适应
	fmt.Println("\n📝 测试4: 标题长度自适应测试...")
	titleTests := []struct {
		title    string
		expected string
	}{
		{"短标题", "极大字号"},
		{"这是一个中等标题", "大字号"},
		{"这是一个比较长的标题内容", "中字号"},
		{"这是一个非常非常长的标题内容测试", "小字号"},
		{"这是一个超级超级超级超级超级长的标题内容测试", "极小字号"},
	}

	for _, tt := range titleTests {
		coverPath, err := renderer.GenerateCoverOnly(
			tt.title,
			tt.expected,
			"default",
			fmt.Sprintf("test_title_%d", len(tt.title)),
			1080,
			1440,
		)
		if err != nil {
			log.Printf("❌ 标题测试失败: %v", err)
			continue
		}
		fmt.Printf("✅ 标题 '%s' (%s): %s\n", tt.title, tt.expected, coverPath)
	}

	fmt.Println("\n========================================")
	fmt.Println("所有测试完成！")
	fmt.Println("========================================")
}
