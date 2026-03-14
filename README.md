# 小红书自动化管理系统

基于 Vue3 + TailwindCSS + Go + Gorm + Hertz + Eino 构建的小红书内容生成与发布管理系统。

## 功能特性

1. **AI内容生成**：调用 DeepSeek API 自动生成符合小红书风格的技能描述、标题和标签
2. **内容管理**：支持内容的保存、编辑、删除和查看
3. **定时发布**：支持设置发布时间和发布频率
4. **发布记录**：记录所有发布历史，方便管理和查看
5. **用户系统**：支持用户注册、登录和权限管理
6. **可视化界面**：基于 Vue3 + TailwindCSS 的现代化管理界面

## 技术栈

### 前端
- Vue 3
- TypeScript
- TailwindCSS
- Element Plus
- Pinia
- Vue Router
- Axios

### 后端
- Go 1.21+
- Hertz (Web框架)
- GORM (ORM)
- MySQL (数据库)
- JWT (认证)
- Viper (配置管理)

## 项目结构

```
xiaohongshu/
├── frontend/              # 前端项目
│   ├── src/
│   │   ├── api/          # API接口
│   │   ├── assets/       # 静态资源
│   │   ├── components/   # 组件
│   │   ├── composables/  # 组合式函数
│   │   ├── router/       # 路由
│   │   ├── stores/       # 状态管理
│   │   ├── views/        # 页面
│   │   ├── App.vue
│   │   └── main.ts
│   ├── index.html
│   ├── package.json
│   ├── tailwind.config.js
│   └── vite.config.ts
├── backend/              # 后端项目
│   ├── cmd/
│   │   └── server/      # 服务入口
│   ├── internal/
│   │   ├── app/         # 应用核心
│   │   ├── config/      # 配置
│   │   ├── handler/     # 处理器
│   │   ├── middleware/  # 中间件
│   │   ├── model/       # 数据模型
│   │   ├── repository/  # 数据访问
│   │   ├── service/     # 业务逻辑
│   │   └── utils/       # 工具函数
│   ├── pkg/
│   │   ├── errno/       # 错误码
│   │   └── response/    # 响应处理
│   ├── config.yaml      # 配置文件
│   ├── .env.example     # 环境变量示例
│   └── go.mod
└── docs/
    └── requirement.md   # 需求文档
```

## 快速开始

### 前置要求

- Node.js 18+
- Go 1.21+
- MySQL 8.0+

### 1. 数据库准备

创建MySQL数据库：

```sql
CREATE DATABASE xiaohongshu CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 后端配置

进入后端目录：

```bash
cd backend
```

复制配置文件并修改：

```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库连接和 DeepSeek API Key
```

或者修改 `config.yaml` 文件。

安装依赖：

```bash
go mod tidy
```

启动后端服务：

**方式一：使用 Air 热重载（推荐开发时使用）**

```bash
# 首先安装 Air（如果未安装）
go install github.com/cosmtrek/air@latest

# 使用 Air 启动（支持代码变更自动重启）
air
```

**方式二：直接运行**

```bash
go run cmd/server/main.go
```

后端服务将在 http://localhost:8000 启动

### 3. 前端配置

进入前端目录：

```bash
cd frontend
```

安装依赖：

```bash
npm install
```

启动开发服务器：

```bash
npm run dev
```

前端服务将在 http://localhost:5173 启动

## API 文档

### 认证接口

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

### 用户接口

- `GET /api/v1/user/info` - 获取当前用户信息
- `GET /api/v1/users` - 获取用户列表（管理员）

### 内容接口

- `POST /api/v1/content/generate` - 生成内容
- `POST /api/v1/content/save` - 保存内容
- `GET /api/v1/content/list` - 获取内容列表
- `GET /api/v1/content/:id` - 获取内容详情
- `PUT /api/v1/content/:id` - 更新内容
- `DELETE /api/v1/content/:id` - 删除内容

### 发布接口

- `POST /api/v1/publish/schedule` - 定时发布
- `POST /api/v1/publish/now` - 立即发布
- `GET /api/v1/publish/list` - 获取发布记录列表
- `GET /api/v1/publish/:id` - 获取发布记录详情
- `DELETE /api/v1/publish/:id/cancel` - 取消发布

## 配置说明

### DeepSeek API

获取 DeepSeek API Key：https://platform.deepseek.com/

在配置文件中设置：

```yaml
deepseek:
  api_key: "your_api_key_here"
  model: deepseek-chat
  base_url: https://api.deepseek.com
```

### JWT 配置

生产环境请务必修改 JWT 密钥：

```yaml
jwt:
  secret: "your-secret-key-here"
  expire: 24
```

## 开发说明

### 前端开发

- 前端使用 Vite 作为构建工具
- 使用 TailwindCSS 进行样式开发
- 使用 Element Plus 作为 UI 组件库
- 使用 Pinia 进行状态管理

### 后端开发

- 后端使用 Hertz 作为 Web 框架
- 使用 GORM 进行数据库操作
- 遵循 MVC 架构模式
- 代码结构清晰，易于扩展
- 使用 Air 进行开发时的热重载

#### Air 热重载配置

项目已配置 `.air.toml，默认监听端口 8000。Air 会自动监控代码变更并重启服务。

## 许可证

MIT License
