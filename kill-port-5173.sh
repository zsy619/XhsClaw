#!/bin/bash

# 终止5173端口的进程

echo "正在查找5173端口的进程..."

# 使用lsof查找5173端口的进程
PORT=5173
PID=$(lsof -ti :$PORT)

if [ -n "$PID" ]; then
    echo "找到进程: $PID"
    echo "正在终止进程..."
    kill -9 $PID
    echo "进程已终止"
else
    echo "未找到5173端口的进程"
fi

echo "操作完成"
