#!/bin/bash

# 测试新增 20 个样式的渲染效果
# 使用简单的 Markdown 内容测试每个样式

BASE_URL="http://localhost:8000/api/v1/xiaohongshu-renderer"
OUTPUT_DIR="/Volumes/E/JYW/创意项目/XhsClaw/test-styles"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 测试内容
TEST_MARKDOWN="---
title: 测试标题
tags: [测试，样式，小红书]
---

# 这是一个测试标题

## 小标题 1

这里是测试内容，用于验证样式的渲染效果。

### 要点列表

- 第一个要点
- 第二个要点
- 第三个要点

**粗体文字** 和 *斜体文字*

> 这是一个引用块

---

页脚内容"

# 新增的 20 个样式列表
STYLES=(
    "cream-custard"
    "sakura-pink"
    "matcha-latte"
    "blueberry-cheese"
    "caramel-macchiato"
    "honey-peach"
    "vanilla-milk"
    "chocolate-mint"
    "strawberry-milk"
    "mango-pudding"
    "taro-milktea"
    "coconut-cream"
    "red-velvet"
    "pistachio-green"
    "bubblegum-pink"
    "lemon-meringue"
    "blackberry-sage"
    "peaches-cream"
    "earl-grey"
    "tiramisu"
    "pomegranate"
    "sage-green"
)

echo "开始测试新增的 20 个样式渲染效果..."
echo "========================================"

for style in "${STYLES[@]}"; do
    echo "测试样式：$style"
    
    # 发送渲染请求
    curl -X POST "$BASE_URL/render" \
        -H "Content-Type: application/json" \
        -d "{
            \"markdown_content\": $(echo "$TEST_MARKDOWN" | jq -Rs .),
            \"style_key\": \"$style\",
            \"output_prefix\": \"test_${style}\",
            \"enable_smart_pagination\": false,
            \"card_width\": 1080,
            \"card_height\": 1440
        }" \
        -o "$OUTPUT_DIR/${style}_response.json" 2>/dev/null
    
    if [ $? -eq 0 ]; then
        echo "  ✅ $style 渲染成功"
    else
        echo "  ❌ $style 渲染失败"
    fi
done

echo "========================================"
echo "所有样式测试完成！"
echo "结果保存在：$OUTPUT_DIR"
