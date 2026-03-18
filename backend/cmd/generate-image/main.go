// Package main 图片生成测试程序
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

	// 示例 Markdown 内容
	markdownContent := `📖 小红书运营技巧分享

如何让你的笔记获得更多曝光？

- 标题要吸引人，使用表情符号增加点击率
- 内容要有价值，解决用户实际问题
- 配图要精美，符合小红书风格
- 标签要精准，提高搜索曝光率

1. 每天坚持发布，保持活跃度
2. 互动评论区，增加账号权重
3. 分析数据，优化内容策略

记住：内容为王，质量第一！

#小红书运营 #内容创作 #涨粉技巧`

	fmt.Println("正在生成小红书笔记图片...")

	// 生成图片
	images, err := renderer.RenderMarkdownToImage(
		markdownContent,
		"default",                         // 主题样式
		"test_note",                       // 输出前缀
		service.PaginationSeparator,       // 分页模式
		1080,                              // 宽度
		1440,                              // 高度
		4320,                              // 最大高度
	)
	if err != nil {
		log.Fatalf("生成图片失败: %v", err)
	}

	fmt.Println("✅ 图片生成成功！")
	fmt.Printf("共生成 %d 张图片:\n", len(images))
	for i, img := range images {
		fmt.Printf("  %d. %s\n", i+1, img)
	}
	fmt.Println("\n图片保存在 ./public/images/ 目录下")
}
