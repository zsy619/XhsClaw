# Vite 8 技能文档

## 目录
1. [Vite 8 概述](#vite-8-概述)
2. [核心特性与新功能](#核心特性与新功能)
3. [从旧版本升级](#从旧版本升级)
4. [配置优化方法](#配置优化方法)
5. [性能提升技巧](#性能提升技巧)
6. [常见问题解决方案](#常见问题解决方案)
7. [与之前版本的主要差异对比](#与之前版本的主要差异对比)

---

## Vite 8 概述

Vite 8 是 Vite 构建工具的最新主要版本，于 2026 年 3 月 12 日发布。该版本带来了重大的架构改进、性能优化和新特性，其中最重要的变化是与 Rolldown 的深度集成。

### 主要亮点
- **Rolldown-Vite 合并**：将 Rolldown 作为默认的打包器，提供更快的构建速度
- **更新的浏览器目标**：支持更现代的浏览器特性
- **开发工具集成**：内置 DevTools 支持
- **改进的控制台日志**：将浏览器控制台日志转发到开发服务器终端

---

## 核心特性与新功能

### 1. Rolldown 集成（最重大的变化）

Vite 8 最重要的改进是与 Rolldown 的深度集成。Rolldown 是一个用 Rust 编写的高性能 JavaScript 打包器，作为 Rollup 的替代品。

**优势：**
- 显著更快的构建速度
- 更好的 Tree Shaking 性能
- 更低的内存占用
- 与 Rollup 插件兼容

**相关配置变化：**
```typescript
// vite.config.ts
export default defineConfig({
  build: {
    // 使用默认的 rolldown 压缩器（推荐）
    minify: 'rolldown', 
    // 或者继续使用 terser（如果需要高级选项）
    // minify: 'terser',
  }
})
```

### 2. 更新的默认浏览器目标

Vite 8 更新了默认的浏览器兼容性目标，支持更现代的浏览器特性。

**新的默认目标包括：**
- ES2025 语法支持
- 更新的 iOS 目标
- LightningCSS 的 ES2024/ES2025 构建目标支持

### 3. 浏览器控制台日志转发

Vite 8 新增了将浏览器控制台日志和错误转发到开发服务器终端的功能，方便调试。

### 4. 集成 DevTools

Vite 8 内置了 DevTools 支持，提供更好的开发体验。

### 5. WASM SSR 支持

为 `.wasm?init` 添加了 SSR 支持。

### 6. Manifest 增强

为独立的 CSS 入口点添加了 `assets` 字段。

### 7. 其他改进

- 支持通配符主机上的端口冲突检测
- 快捷键不区分大小写
- 为 `optimizeDeps` 添加了 `ignoreOutdatedRequests` 选项
- 高度实验性的完整打包模式
- 支持环境变量文件挂载（FIFOs）

---

## 从旧版本升级

### 前置要求

- **Node.js 版本**：需要 Node.js 20.19+ 或 22.12+
- **现有项目**：建议从 Vite 6 或 7 升级

### 升级步骤

#### 1. 更新 package.json 依赖

```json
{
  "devDependencies": {
    "vite": "^8.0.0",
    "@vitejs/plugin-vue": "^7.0.0",
    "@vitejs/plugin-vue-jsx": "^5.0.0"
  }
}
```

#### 2. 处理破坏性变更

Vite 8 包含以下破坏性变更：

##### a. 移除 import.meta.hot.accept 解析回退

**影响：** 如果你在使用 `import.meta.hot.accept` 时依赖旧的回退行为，需要更新代码。

**解决方案：** 确保显式指定要接受的模块。

##### b. 更新默认浏览器目标

**影响：** 构建输出可能不再兼容非常旧的浏览器。

**解决方案：** 如果你需要支持旧浏览器，可以在配置中显式设置 `build.target`：

```typescript
export default defineConfig({
  build: {
    target: 'es2020' // 根据需要调整
  }
})
```

##### c. Rolldown-Vite 合并

**影响：** 构建行为可能有细微变化。

**解决方案：** 大多数情况下无需修改，但如果你使用了特定的 Rollup 插件或配置，可能需要测试验证。

#### 3. 更新 vite.config.ts

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  build: {
    // 如果使用了 terser 的 productionTypes 选项，需要移除
    minify: 'terser', // 或者改为 'rolldown' 以获得更好的性能
    terserOptions: {
      compress: {
        // productionTypes 选项已移除
        drop_console: true,
        drop_debugger: true,
      },
    },
  },
})
```

#### 4. 清理并重新安装依赖

```bash
# 删除旧的依赖
rm -rf node_modules package-lock.json

# 重新安装
npm install

# 或使用 yarn
rm -rf node_modules yarn.lock
yarn install
```

#### 5. 测试构建

```bash
# 运行开发服务器
npm run dev

# 构建生产版本
npm run build
```

---

## 配置优化方法

### 1. 利用 Rolldown 进行更快的构建

```typescript
export default defineConfig({
  build: {
    // 使用 rolldown 作为默认压缩器（推荐）
    minify: 'rolldown',
    
    // 配置 rolldown 选项
    rollupOptions: {
      // 你的 rollup 配置仍然有效
      output: {
        manualChunks: {
          'vendor-vue': ['vue', 'vue-router', 'pinia'],
        },
      },
    },
  },
})
```

### 2. 优化依赖预构建

```typescript
export default defineConfig({
  optimizeDeps: {
    // 忽略过时的请求，提高开发体验
    ignoreOutdatedRequests: true,
    
    // 显式指定需要预构建的依赖
    include: ['vue', 'vue-router', 'pinia'],
    
    // 排除不需要预构建的依赖
    exclude: ['some-large-package'],
  },
})
```

### 3. 配置 CSS 优化

```typescript
export default defineConfig({
  css: {
    // 使用 LightningCSS（如果需要）
    lightningcss: {
      targets: {
        chrome: 100,
        safari: 16,
      },
    },
    preprocessorOptions: {
      scss: {
        // SCSS 配置
      },
    },
  },
})
```

### 4. 开发服务器优化

```typescript
export default defineConfig({
  server: {
    // 配置允许的主机（防止 DNS 重绑定攻击）
    allowedHosts: ['localhost', '.your-domain.com'],
    
    // 配置 CORS
    cors: true,
    
    // 预热请求
    warmup: {
      clientFiles: ['./src/main.ts', './src/App.vue'],
    },
  },
})
```

---

## 性能提升技巧

### 1. 使用 Rolldown 替代 Terser

```typescript
export default defineConfig({
  build: {
    minify: 'rolldown', // 比 terser 快得多
  },
})
```

### 2. 合理配置代码分割

```typescript
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          // 将第三方库拆分成独立的 chunk
          'vendor-vue': ['vue', 'vue-router', 'pinia'],
          'vendor-utils': ['axios', 'dayjs', 'lodash-es'],
          'vendor-charts': ['echarts'],
          'vendor-editor': ['md-editor-v3', 'marked'],
        },
      },
    },
  },
})
```

### 3. 优化开发环境

```typescript
export default defineConfig({
  optimizeDeps: {
    // 只在需要时重新优化依赖
    ignoreOutdatedRequests: true,
  },
  server: {
    // 预热常用文件
    warmup: {
      clientFiles: ['./src/main.ts'],
    },
  },
})
```

### 4. 使用 Gzip/Brotli 压缩

```typescript
import viteCompression from 'vite-plugin-compression'

export default defineConfig({
  plugins: [
    viteCompression({
      algorithm: 'gzip',
      threshold: 10240,
    }),
    // 可以同时添加 Brotli 压缩
    viteCompression({
      algorithm: 'brotliCompress',
      ext: '.br',
      threshold: 10240,
    }),
  ],
})
```

### 5. 减少构建时间

```typescript
export default defineConfig({
  build: {
    // 禁用不需要的 sourcemap
    sourcemap: false,
    
    // 配置输出目录
    outDir: 'dist',
    assetsDir: 'static',
  },
})
```

---

## 常见问题解决方案

### 1. 依赖安装失败

**问题：** 升级后 npm install 失败

**解决方案：**
```bash
# 清理 npm 缓存
npm cache clean --force

# 删除 node_modules 和 lock 文件
rm -rf node_modules package-lock.json

# 重新安装
npm install
```

### 2. 构建错误：找不到模块

**问题：** 构建时出现模块找不到的错误

**解决方案：**
```typescript
// 检查 resolve.alias 配置是否正确
export default defineConfig({
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
})
```

### 3. HMR 不工作

**问题：** 热模块替换不工作

**解决方案：**
- 确保没有使用 `import.meta.hot.accept` 的旧回退行为
- 检查 `vite.config.ts` 中的 server 配置
- 尝试清除浏览器缓存

### 4. TypeScript 类型错误

**问题：** 升级后出现 TypeScript 类型错误

**解决方案：**
```bash
# 更新 vue-tsc
npm install vue-tsc@latest --save-dev

# 运行类型检查
npm run build
```

### 5. 插件兼容性问题

**问题：** 某些插件在 Vite 8 中不工作

**解决方案：**
- 检查插件是否有支持 Vite 8 的版本
- 查看插件的 GitHub issues
- 考虑替代插件或暂时降级

### 6. 构建速度变慢

**问题：** 升级后构建速度反而变慢

**解决方案：**
```typescript
export default defineConfig({
  build: {
    // 确保使用 rolldown 压缩器
    minify: 'rolldown',
    
    // 检查是否有不必要的插件
  },
})
```

---

## 与之前版本的主要差异对比

### Vite 8 vs Vite 7

| 特性 | Vite 7 | Vite 8 |
|------|--------|--------|
| 默认打包器 | Rollup | Rolldown |
| 默认压缩器 | Terser | Rolldown |
| Node.js 要求 | 20.19+, 22.12+ | 20.19+, 22.12+ |
| 浏览器目标 | 更新过 | 进一步更新 |
| DevTools 集成 | 部分 | 完整 |
| 控制台日志转发 | 无 | 有 |
| import.meta.hot.accept 回退 | 有 | 已移除 |

### Vite 8 vs Vite 6

| 特性 | Vite 6 | Vite 8 |
|------|--------|--------|
| 默认打包器 | Rollup | Rolldown |
| 默认压缩器 | Terser | Rolldown |
| Node.js 要求 | 18+ | 20.19+, 22.12+ |
| 浏览器目标 | 较旧 | 现代 |
| DevTools 集成 | 无 | 完整 |
| 控制台日志转发 | 无 | 有 |
| server.allowedHosts | 无 | 有 |
| server.cors 默认值 | true | false |

### 架构变化

#### Vite 6/7 架构
```
源代码 → Vite 开发服务器 → ESBuild（预构建）
                    ↓
              Rollup（生产构建）
```

#### Vite 8 架构
```
源代码 → Vite 开发服务器 → ESBuild（预构建）
                    ↓
              Rolldown（生产构建）
```

### 性能对比

根据 Vite 团队的测试，Vite 8 相比 Vite 7：

- **开发服务器启动**：类似或稍快
- **HMR 更新**：类似或稍快
- **生产构建**：快 20-50%（取决于项目大小）
- **内存占用**：降低 10-30%

---

## 总结

Vite 8 是一个重大的版本升级，带来了显著的性能改进和新特性。虽然有一些破坏性变更，但对于大多数项目来说，升级过程相对简单。

**升级建议：**
1. 先在开发环境测试
2. 确保所有依赖都兼容 Vite 8
3. 充分测试构建后的应用
4. 利用 Rolldown 获得更好的性能
5. 阅读官方迁移指南获取最新信息

**参考资源：**
- [Vite 官方文档](https://vitejs.dev/)
- [Vite 8 发布公告](https://github.com/vitejs/vite/releases/tag/v8.0.0)
- [Vite 迁移指南](https://vitejs.dev/guide/migration.html)

