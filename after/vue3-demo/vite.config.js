import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        /* 是否开启 $ref , vue3.2 的语法糖 */
        vue({
            refTransform: true,
            reactivityTransform: true
        }),
    ],
    server: {
        // port: 8080 //指定端口号

        proxy: {
            '/user': {
                target: "http://127.0.0.1:3000/"
            }
        }
    },

    base: './', //打包相对路径
})
