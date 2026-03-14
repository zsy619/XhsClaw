#!/bin/bash

# 小红书自动化管理系统 - 启动脚本

echo "========================================="
echo "  小红书自动化管理系统"
echo "========================================="

# 检查是否已安装 Go
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.21+"
    exit 1
fi

# 检查是否已安装 Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js 未安装，请先安装 Node.js 18+"
    exit 1
fi

# 检查是否已安装 npm
if ! command -v npm &> /dev/null; then
    echo "❌ npm 未安装"
    exit 1
fi

# 创建日志目录
mkdir -p logs

echo ""
echo "📦 正在安装后端依赖..."
cd backend
go mod tidy
if [ $? -ne 0 ]; then
    echo "❌ 后端依赖安装失败"
    exit 1
fi
echo "✅ 后端依赖安装成功"

echo ""
echo "📦 正在安装前端依赖..."
cd ../frontend
npm install
if [ $? -ne 0 ]; then
    echo "❌ 前端依赖安装失败"
    exit 1
fi
echo "✅ 前端依赖安装成功"

echo ""
echo "🚀 启动服务..."
echo ""
echo "提示："
echo "  - 后端服务将在 http://localhost:8000 启动"
echo "  - 前端服务将在 http://localhost:5173 启动"
echo "  - 请确保 MySQL 数据库已启动，并配置好数据库连接"
echo ""
echo "按 Ctrl+C 停止服务"
echo ""

# 回到项目根目录
cd ..

# 在后台启动后端
echo "启动后端服务..."
cd backend
nohup go run cmd/server/main.go > ../logs/backend.log 2>&1 &
BACKEND_PID=$!
echo "后端服务 PID: $BACKEND_PID"
cd ..

# 等待一下
sleep 3

# 启动前端
echo "启动前端服务..."
cd frontend
npm run dev
