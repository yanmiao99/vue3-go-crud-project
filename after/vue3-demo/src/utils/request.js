/*
*
* axios 二次封装
*
* 全局配置
*
* 响应拦截
*
* request 请求的方法
*
* 封装成 对象调用的方式
*
* request.get('/api',{xxx:xxx})
*
*
* */

// 引入文件
import axios from "axios"
import {ElMessage} from "element-plus";

// 全局配置
const service = axios.create({
    baseURL: '/user', // 根路径
    timeout: 8000  // 请求超时时间
})

// 响应拦截
service.interceptors.response.use(res => {
    const {code, data, msg} = res.data

    if (code === 200) {
        // 请求成功
        ElMessage.success(msg)
        return data
    } else if (code === 400) {
        // 请求失败
        ElMessage.error(msg)
    }
})

// request 方法
function request(options) {
    options.method = options.method || 'get'

    if (options.method.toLowerCase() === 'get') options.params = options.data

    return service(options)
}

// * request.post('/api',{xxx:xxx})

['get', 'post', 'put', 'delete'].forEach(item => {
    request[item] = (url, data) => {
        return request({
            url,
            data,
            method: item
        })
    }
})

export default request



