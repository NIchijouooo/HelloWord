import { createApp } from 'vue'
import App from './App.vue'
// 覆盖element样式
//import 'styles/element-variables.scss'
//import ElementPlus from 'element-plus'
import router from './router'
import { createPinia } from 'pinia'
import axios from 'axios'
import { configStore } from '@/stores/app.js'
import './assets/css/reset.css'
import { registerIcons } from 'utils/common'

const pinia = createPinia()
const app = createApp(App)
app.use(pinia)
const config = configStore()
registerIcons(app)
axios.get('./config.json').then((res) => {
  if (res.status === 200) {
    config.setConfigInfo(res.data)
  } else {
    ElMessage({
      type: 'error',
      message: '缺少config.json配置文件！',
    })
  }
  console.log(res)
  app.use(router)
  //渲染
  app.mount('#app')
})
