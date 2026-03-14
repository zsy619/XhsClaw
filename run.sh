#!/bin/bash
"""
小红书自动发布系统启动脚本

功能：
  - 启动后端服务（FastAPI）
  - 启动前端服务（HTTP 服务器）
  - 检查服务状态
  - 停止服务

使用方法：
  ./run.sh start     # 启动所有服务
  ./run.sh status    # 检查服务状态
  ./run.sh stop      # 停止所有服务
  ./run.sh restart   # 重启所有服务
"""

# 项目根目录
PROJECT_ROOT="$(dirname "$(realpath "$0")")"
BACKEND_DIR="${PROJECT_ROOT}/backend"
FRONTEND_DIR="${PROJECT_ROOT}/frontend"

# 端口配置
BACKEND_PORT=8000
FRONTEND_PORT=5173

# 日志文件
LOG_DIR="${PROJECT_ROOT}/logs"
BACKEND_LOG="${LOG_DIR}/backend.log"
FRONTEND_LOG="${LOG_DIR}/frontend.log"

# 进程文件
PID_DIR="${PROJECT_ROOT}/pids"
BACKEND_PID="${PID_DIR}/backend.pid"
FRONTEND_PID="${PID_DIR}/frontend.pid"

# 颜色定义
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

# 创建必要的目录
mkdir -p "${LOG_DIR}"
mkdir -p "${PID_DIR}"

# 显示日志
show_log() {
    local log_file="$1"
    local service_name="$2"
    
    echo -e "${BLUE}=== ${service_name} 日志 ===${NC}"
    if [ -f "${log_file}" ]; then
        tail -n 20 "${log_file}"
    else
        echo -e "${YELLOW}日志文件不存在${NC}"
    fi
    echo
}

# 检查端口是否被占用
check_port() {
    local port="$1"
    local service_name="$2"
    
    if lsof -i :${port} > /dev/null 2>&1; then
        echo -e "${GREEN}${service_name} 服务已在端口 ${port} 运行${NC}"
        return 0
    else
        echo -e "${YELLOW}${service_name} 服务未运行${NC}"
        return 1
    fi
}

# 启动后端服务
start_backend() {
    echo -e "${BLUE}启动后端服务...${NC}"
    
    # 检查后端是否已运行
    if check_port ${BACKEND_PORT} "后端"; then
        return
    fi
    
    # 检查虚拟环境
    if [ ! -d "${BACKEND_DIR}/venv" ]; then
        echo -e "${RED}错误: 后端虚拟环境不存在，请先创建${NC}"
        echo -e "${BLUE}提示: 运行 ${BACKEND_DIR}/create_venv.sh 创建虚拟环境${NC}"
        return 1
    fi
    
    # 启动后端服务
    cd "${BACKEND_DIR}" && source venv/bin/activate && uvicorn app.main:app --host 0.0.0.0 --port ${BACKEND_PORT} > "${BACKEND_LOG}" 2>&1 &
    
    # 保存进程 ID
    echo $! > "${BACKEND_PID}"
    
    # 等待服务启动
    echo -e "${BLUE}等待后端服务启动...${NC}"
    sleep 3
    
    if check_port ${BACKEND_PORT} "后端"; then
        echo -e "${GREEN}后端服务启动成功！${NC}"
        echo -e "${GREEN}访问地址: http://localhost:${BACKEND_PORT}${NC}"
        echo -e "${GREEN}API 文档: http://localhost:${BACKEND_PORT}/docs${NC}"
    else
        echo -e "${RED}后端服务启动失败${NC}"
        show_log "${BACKEND_LOG}" "后端"
    fi
}

# 启动前端服务
start_frontend() {
    echo -e "${BLUE}启动前端服务...${NC}"
    
    # 检查前端是否已运行
    if check_port ${FRONTEND_PORT} "前端"; then
        return
    fi
    
    # 检查 node_modules 是否存在
    if [ ! -d "${FRONTEND_DIR}/node_modules" ]; then
        echo -e "${YELLOW}提示: 前端依赖未安装，正在执行 npm install...${NC}"
        cd "${FRONTEND_DIR}" && npm install
    fi
    
    # 启动前端服务（使用 Vite 开发服务器）
    cd "${FRONTEND_DIR}" && npm run dev > "${FRONTEND_LOG}" 2>&1 &
    
    # 保存进程 ID
    echo $! > "${FRONTEND_PID}"
    
    # 等待服务启动（Vite 可能需要更长时间）
    echo -e "${BLUE}等待前端服务启动...${NC}"
    sleep 5
    
    if check_port ${FRONTEND_PORT} "前端"; then
        echo -e "${GREEN}前端服务启动成功！${NC}"
        echo -e "${GREEN}访问地址: http://localhost:${FRONTEND_PORT}${NC}"
    else
        echo -e "${RED}前端服务启动失败${NC}"
        show_log "${FRONTEND_LOG}" "前端"
    fi
}

# 停止后端服务
stop_backend() {
    echo -e "${BLUE}停止后端服务...${NC}"
    
    if [ -f "${BACKEND_PID}" ]; then
        local pid=$(cat "${BACKEND_PID}")
        if kill -0 ${pid} 2>/dev/null; then
            kill ${pid}
            echo -e "${GREEN}后端服务已停止${NC}"
            rm "${BACKEND_PID}"
        else
            echo -e "${YELLOW}后端服务进程不存在${NC}"
            rm "${BACKEND_PID}"
        fi
    else
        echo -e "${YELLOW}后端服务未运行${NC}"
    fi
}

# 停止前端服务
stop_frontend() {
    echo -e "${BLUE}停止前端服务...${NC}"
    
    if [ -f "${FRONTEND_PID}" ]; then
        local pid=$(cat "${FRONTEND_PID}")
        if kill -0 ${pid} 2>/dev/null; then
            kill ${pid}
            echo -e "${GREEN}前端服务已停止${NC}"
            rm "${FRONTEND_PID}"
        else
            echo -e "${YELLOW}前端服务进程不存在${NC}"
            rm "${FRONTEND_PID}"
        fi
    else
        echo -e "${YELLOW}前端服务未运行${NC}"
    fi
}

# 检查服务状态
check_status() {
    echo -e "${BLUE}=== 服务状态检查 ===${NC}"
    
    # 检查后端状态
    if check_port ${BACKEND_PORT} "后端"; then
        echo -e "${GREEN}后端服务: 运行中${NC}"
        echo -e "${GREEN}后端地址: http://localhost:${BACKEND_PORT}${NC}"
    else
        echo -e "${RED}后端服务: 未运行${NC}"
    fi
    
    # 检查前端状态
    if check_port ${FRONTEND_PORT} "前端"; then
        echo -e "${GREEN}前端服务: 运行中${NC}"
        echo -e "${GREEN}前端地址: http://localhost:${FRONTEND_PORT}${NC}"
    else
        echo -e "${RED}前端服务: 未运行${NC}"
    fi
    
    echo
}

# 主函数
main() {
    case "$1" in
        start)
            echo -e "${BLUE}=== 启动所有服务 ===${NC}"
            start_backend
            echo
            start_frontend
            echo
            check_status
            ;;
        stop)
            echo -e "${BLUE}=== 停止所有服务 ===${NC}"
            stop_backend
            echo
            stop_frontend
            echo
            check_status
            ;;
        restart)
            echo -e "${BLUE}=== 重启所有服务 ===${NC}"
            stop_backend
            stop_frontend
            echo
            start_backend
            echo
            start_frontend
            echo
            check_status
            ;;
        status)
            check_status
            ;;
        logs)
            echo -e "${BLUE}=== 查看服务日志 ===${NC}"
            show_log "${BACKEND_LOG}" "后端"
            show_log "${FRONTEND_LOG}" "前端"
            ;;
        *)
            echo -e "${BLUE}使用方法:${NC}"
            echo -e "  ${GREEN}./run.sh start${NC}     # 启动所有服务"
            echo -e "  ${GREEN}./run.sh stop${NC}      # 停止所有服务"
            echo -e "  ${GREEN}./run.sh restart${NC}   # 重启所有服务"
            echo -e "  ${GREEN}./run.sh status${NC}    # 检查服务状态"
            echo -e "  ${GREEN}./run.sh logs${NC}      # 查看服务日志"
            ;;
    esac
}

# 执行主函数
main "$1"
