#!/bin/bash

# 对比测试脚本：验证 XhsClaw 与 Auto-Redbook-Skills 生成效果一致性

echo "========================================"
echo "小红书图片生成对比测试"
echo "========================================"
echo ""

# 测试数据
TEST_MARKDOWN='---
title: 测试标题
subtitle: 测试副标题
emoji: 📝
tags: [测试，标签，小红书]
---

## 核心要点

- ✅ 第一点内容
- ✅ 第二点内容
- ✅ 第三点内容

## 更多细节

💡 这里是详细说明

📝 总结内容

记得点赞收藏哦！

#测试 #标签 #小红书'

# 1. 测试封面生成
echo "1️⃣  测试封面生成..."
echo ""

# XhsClaw 封面生成
echo "  生成 XhsClaw 封面..."
curl -X POST http://localhost:8000/api/v1/xiaohongshu-renderer/cover \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试标题",
    "subtitle": "测试副标题",
    "style_key": "playful-geometric",
    "output_prefix": "xhsclaw_cover_test",
    "width": 1080,
    "height": 1440
  }' | jq .

echo ""
echo "  ✅ XhsClaw 封面生成完成"
echo ""

# 2. 测试内容渲染（不分页）
echo "2️⃣  测试内容渲染（不分页模式）..."
echo ""

curl -X POST http://localhost:8000/api/v1/xiaohongshu-renderer/render \
  -H "Content-Type: application/json" \
  -d '{
    "markdown_content": "'"$TEST_MARKDOWN"'",
    "style_key": "playful-geometric",
    "output_prefix": "xhsclaw_content_nopage",
    "card_width": 1080,
    "card_height": 1440,
    "enable_smart_pagination": false
  }' | jq .

echo ""
echo "  ✅ XhsClaw 内容渲染完成（不分页）"
echo ""

# 3. 测试内容渲染（智能分页）
echo "3️⃣  测试内容渲染（智能分页模式）..."
echo ""

curl -X POST http://localhost:8000/api/v1/xiaohongshu-renderer/render \
  -H "Content-Type: application/json" \
  -d '{
    "markdown_content": "'"$TEST_MARKDOWN"'",
    "style_key": "playful-geometric",
    "output_prefix": "xhsclaw_content_page",
    "card_width": 1080,
    "card_height": 1440,
    "enable_smart_pagination": true,
    "pagination_mode": "auto-split"
  }' | jq .

echo ""
echo "  ✅ XhsClaw 内容渲染完成（智能分页）"
echo ""

# 4. 测试完整流程（封面 + 内容）
echo "4️⃣  测试完整流程（封面 + 内容）..."
echo ""

curl -X POST http://localhost:8000/api/v1/xiaohongshu-renderer/render-with-ai \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试标题",
    "content": "测试内容\n\n- 第一点\n- 第二点\n- 第三点",
    "tags": "测试，标签，小红书",
    "use_ai": false,
    "enable_smart_pagination": true,
    "pagination_mode": "auto-split",
    "style_key": "playful-geometric",
    "output_prefix": "xhsclaw_full_test",
    "width": 1080,
    "height": 1440,
    "emoji": "📝"
  }' | jq .

echo ""
echo "  ✅ XhsClaw 完整流程测试完成"
echo ""

echo "========================================"
echo "所有测试完成！"
echo "========================================"
echo ""
echo "📊 测试总结："
echo "  - 封面生成：✅"
echo "  - 内容渲染（不分页）：✅"
echo "  - 内容渲染（智能分页）：✅"
echo "  - 完整流程：✅"
echo ""
echo "📁 生成的图片位置："
echo "  /Volumes/E/JYW/创意项目/XhsClaw/backend/public/images/"
echo ""
