import { defineConfig, type PluginOption } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import { resolve } from 'path'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import viteCompression from 'vite-plugin-compression'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
    // SVG 图标插件
    createSvgIconsPlugin({
      iconDirs: [resolve(process.cwd(), 'src/assets/icons')],
      symbolId: 'icon-[dir]-[name]',
    }),
    // Gzip 压缩
    viteCompression({
      verbose: true,
      disable: false,
      threshold: 10240,
      algorithm: 'gzip',
      ext: '.gz',
    }),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    open: false,
    proxy: {
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'static',
    sourcemap: false,
    // Vite 8 中，我们可以使用默认的 rolldown 压缩器，它比 terser 更快
    // 但如果需要更高级的压缩选项，可以继续使用 terser
    minify: 'terser',
    terserOptions: {
      compress: {
        // productionTypes 在 Vite 8 中可能已被移除，我们保留其他选项
        drop_console: true,
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor-vue': ['vue', 'vue-router', 'pinia'],
          'vendor-utils': ['axios', 'dayjs', 'lodash-es'],
          'vendor-charts': ['echarts'],
          'vendor-editor': ['md-editor-v3', 'marked'],
        },
      },
    },
  },
  css: {
    preprocessorOptions: {
      scss: {
        // 注意：移除全局注入，避免路径问题
        // 变量文件在 src/assets/styles/variables.scss
        // 需要使用时手动导入：@import './variables.scss';
      },
    },
  },
})
