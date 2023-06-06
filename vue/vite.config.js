import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import ElementPlus from 'unplugin-element-plus/vite'
import path from 'path'
import { viteMockServe } from 'vite-plugin-mock'
import { svgBuilder } from './src/components/icon/svg/index'

const { ElementPlusResolver } = require('unplugin-vue-components/resolvers')
// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  let prodMock = false
  return {
    plugins: [
      vue(),
      AutoImport({
        imports: ['pinia', 'vue', 'vue-router'],
        resolvers: [ElementPlusResolver()],
      }),
      ElementPlus({
        importStyle: 'sass',
        useSource: true,
      }),
      Components({
        resolvers: [ElementPlusResolver()],
      }),
      viteMockServe({
        mockPath: './mock',
        supportTs: false,
        localEnabled: command === 'serve',
        prodEnabled: command !== 'serve' && prodMock,
        watchFiles: true,
        injectCode: `
          import { setupProdMockServer } from './mockProdServer';
          setupProdMockServer();
        `,
      }),
      svgBuilder('./src/assets/icons/'),
    ],
    base: '/',
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src'),
        comps: path.resolve(__dirname, 'src/components'),
        styles: path.resolve(__dirname, 'src/styles'),
        views: path.resolve(__dirname, 'src/views'),
        layout: path.resolve(__dirname, 'src/layout'),
        utils: path.resolve(__dirname, 'src/utils'),
        api: path.resolve(__dirname, 'src/api'),
        assets: path.resolve(__dirname, 'src/assets'),
        stores: path.resolve(__dirname, 'src/stores'),
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: `@use "@/styles/element-variables.scss" as *;`,
        },
      },
    },
    // 打包配置
    build: {
      target: 'modules',
      outDir: '../webroot', // 指定输出路径
      emptyOutDir: true,
      assetsDir: 'assets', // 指定生成静态资源的存放路径
      minify: 'terser', // 混淆器，terser构建后文件体积更小
      chunkSizeWarningLimit: 1000, // 包限制提醒
      rollupOptions: {
        // rollup打包配置
        output: {
          manualChunks(id) {
            if (id.includes('node_modules')) {
              return id.toString().split('node_modules/')[1].split('/')[0].toString()
            }
          },
        },
      },
    },
    // 本地运行配置及代理配置
    server: {
      host: 'localhost',
      port: 8101,
      open: false, // 在服务器启动时自动在浏览器中打开应用程序
      // 代理配置
      proxy: {
        '/api': {
          changeOrigin: true,
          //target: 'http://192.168.4.162:80', // 代理地址 测试环境
          //target: ' http://192.168.4.196:8080', // 代理地址 测试环境
          // target: 'http://127.0.0.1:7070', // 代理地址 本地环境
          target: 'http://120.78.181.76', // 代理地址 本地环境
          rewrite: (path) => path.replace(/^\/api/, '/api'),
        },
      },
    },
  }
})
