#!/bin/bash

# 系统验证脚本
# 验证所有功能模块是否正常

echo "🔍 小红书内容生成与管理系统 - 功能验证"
echo "========================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 计数器
PASSED=0
FAILED=0

# 检查函数
check_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ $2${NC}"
        ((FAILED++))
    fi
}

# 1. 检查 Python 环境
echo "📦 检查 Python 环境..."
if command -v python3 &> /dev/null; then
    PYTHON_VERSION=$(python3 --version)
    check_result 0 "Python 环境：$PYTHON_VERSION"
else
    check_result 1 "Python 环境：未找到 Python3"
fi

# 2. 检查 Node.js 环境
echo ""
echo "📦 检查 Node.js 环境..."
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    check_result 0 "Node.js 环境：$NODE_VERSION"
else
    check_result 1 "Node.js 环境：未找到 Node.js"
fi

# 3. 检查后端依赖
echo ""
echo "📦 检查后端依赖..."
cd backend
if [ -f "requirements.txt" ]; then
    check_result 0 "requirements.txt 存在"
else
    check_result 1 "requirements.txt 不存在"
fi

# 4. 检查前端依赖
echo ""
echo "📦 检查前端依赖..."
cd ../frontend
if [ -f "package.json" ]; then
    check_result 0 "package.json 存在"
else
    check_result 1 "package.json 不存在"
fi

if [ -d "node_modules" ]; then
    check_result 0 "前端依赖已安装"
else
    check_result 1 "前端依赖未安装"
fi

cd ..

# 5. 检查核心文件
echo ""
echo "📄 检查核心文件..."

FILES=(
    "backend/init_database.py"
    "backend/test_system.py"
    "backend/app/main.py"
    "backend/app/domain/services/deepseek_service.py"
    "backend/app/domain/services/icon_generator.py"
    "backend/app/domain/services/xiaohongshu_publisher.py"
    "backend/app/api/routes/generation.py"
    "backend/app/api/routes/publish.py"
    "backend/app/api/routes/stats.py"
    "backend/app/api/routes/health.py"
    "frontend/src/components/MarkdownEditor.vue"
    "frontend/src/components/charts/StatsCharts.vue"
    "frontend/src/views/Dashboard.vue"
    "spec.md"
    "tasks.md"
    "checklist.md"
    "README.md"
    "DEPLOYMENT.md"
    "COMPLETION_REPORT.md"
    "start.sh"
)

for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        check_result 0 "文件存在：$file"
    else
        check_result 1 "文件缺失：$file"
    fi
done

# 6. 检查数据库配置
echo ""
echo "🗄️  检查数据库配置..."
if [ -f "backend/.env" ]; then
    check_result 0 ".env 配置文件存在"
    
    if grep -q "MYSQL_HOST" backend/.env; then
        check_result 0 "MYSQL_HOST 已配置"
    else
        check_result 1 "MYSQL_HOST 未配置"
    fi
    
    if grep -q "DEEPSEEK_API_KEY" backend/.env; then
        check_result 0 "DEEPSEEK_API_KEY 已配置"
    else
        check_result 1 "DEEPSEEK_API_KEY 未配置"
    fi
else
    check_result 1 ".env 配置文件不存在"
fi

# 7. 运行后端测试
echo ""
echo "🧪 运行后端功能测试..."
cd backend
if [ -f "test_system.py" ]; then
    python3 test_system.py
    if [ $? -eq 0 ]; then
        check_result 0 "后端功能测试通过"
    else
        check_result 1 "后端功能测试失败"
    fi
else
    check_result 1 "测试脚本不存在"
fi
cd ..

# 8. 检查启动脚本
echo ""
echo "🚀 检查启动脚本..."
if [ -f "start.sh" ]; then
    if [ -x "start.sh" ]; then
        check_result 0 "启动脚本可执行"
    else
        check_result 1 "启动脚本不可执行"
    fi
else
    check_result 1 "启动脚本不存在"
fi

# 总结
echo ""
echo "========================================"
echo "📊 验证结果汇总"
echo "========================================"
echo -e "${GREEN}✅ 通过：$PASSED${NC}"
echo -e "${RED}❌ 失败：$FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 所有验证通过！系统已准备就绪！${NC}"
    echo ""
    echo "下一步:"
    echo "1. 配置 backend/.env 文件"
    echo "2. 运行 ./start.sh 启动系统"
    echo "3. 访问 http://localhost:5173"
    exit 0
else
    echo -e "${RED}⚠️  部分验证失败，请检查上述错误${NC}"
    exit 1
fi
