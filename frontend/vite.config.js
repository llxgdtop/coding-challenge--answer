import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    // 配置代理解决开发环境跨域问题
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 后端服务地址
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, '') // 如果后端没有 /api 前缀，取消注释
      }
    }
  }
})
