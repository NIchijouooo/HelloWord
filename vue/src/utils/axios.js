import axios from 'axios'
import router from '@/router'
import { ElMessage } from 'element-plus'

const apiUrl = import.meta.env.VITE_REQUEST_API_URL
// 创建 axios
const service = axios.create({
  timeout: 5 * 60 * 100, // 超时时间5分钟
})
// 请求拦截 request
service.interceptors.request.use(
  (config) => {
    // 发请求前做的一些处理，数据转化，配置请求头，设置token,设置loading等
    config.headers['Content-Type'] = 'application/json'
    const pData = config.data
    if (pData.contentType) {
      config.headers['Content-Type'] = pData.contentType
    }
    if (pData.responseType) {
      config.responseType = pData.responseType
    }
    config.headers['token'] = pData.token == null ? '' : pData.token
    // 请求体参数
    if (config.method === 'get') {
      config.params = pData.data
    } else {
      if (pData.contentType) {
        config.data = pData.data
      } else {
        config.data = JSON.stringify(pData.data)
      }
    }

    // console.log('request config', process.env.NODE_ENV)
    if (!pData.mock) {
      config.url = apiUrl + config.url
    } else {
      //参数存在mock表示本地模拟了数据，
      // 判断是否是生产环境
      if (process.env.NODE_ENV !== 'development') {
        config.url = apiUrl + config.url
      }
    }
    return config
  },
  (error) => {
    Promise.reject(error)
  }
)

// 响应拦截 reponse
service.interceptors.response.use(
  (response) => {
    const res = response.data
    // console.log('response res ->', res)
    var fileName = ''
    if (response.headers['content-disposition']) {
      fileName = response.headers['content-disposition'].split(';')[1].split('=')[1]
      return { fileName: decodeURIComponent(fileName), blob: res }
    }
    return res
  },
  (err) => {
    console.log('service.interceptors.response -> err -> ', err)
    if (err.response.status === 400 || err.response.status === 401) {
      ElMessage({
        type: 'error',
        message: err.response.data.message,
      })
      // 跳转到登录页面
      router.push('/')
    } else {
      if (err.response.data.code === 401 || err.response.data.code === 400) {
        ElMessage({
          type: 'warning',
          message: err.response.data.message,
        })
        // 跳转到登录页面
        router.push('/')
      } else {
        if (err && err.message) {
          ElMessage({
            type: 'warning',
            message: err.message,
          })
        }
        return Promise.reject(err)
      }
    }
  }
)

export default service
