#!/bin/bash

# 小红书自动化管理系统 - 停止脚本

echo "正在停止服务..."

# 停止后端（查找占用8000端口的进程）
BACKEND_PID=$(lsof -ti :8000)
if [ -n "$BACKEND_PID" ]; then
    kill -9 $BACKEND_PID
    echo "✅ 后端服务已停止"
else
    echo "ℹ️  后端服务未运行"
fi

# 停止前端（查找占用5173端口的进程）
FRONTEND_PID=$(lsof -ti :5173)
if [ -n "$FRONTEND_PID" ]; then
    kill -9 $FRONTEND_PID
    echo "✅ 前端服务已停止"
else
    echo "ℹ️  前端服务未运行"
fi

echo "所有服务已停止"
