import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [
    vue({
      // 确保正确转换编译宏
      script: {
        defineModel: false,
        propsDestructure: false
      },
      // 确保所有 .vue 文件都被处理
      include: [/\.vue$/]
    })
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
    chunkSizeWarningLimit: 2000,
    // 暂时禁用代码分割，确保 Vue 编译宏正确工作
    rollupOptions: {
      output: {
        manualChunks: undefined
      }
    }
  }
})

