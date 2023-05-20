import { createRouter, createWebHashHistory } from 'vue-router'
import { routes } from './router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { configStore } from '@/stores/app.js'
NProgress.configure({ showSpinner: false })
const router = createRouter({
  history: createWebHashHistory(), // hash模式
  // history: createWebHistory(),// HTML5模式
  routes: [...routes],
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  },
})
// 导航守卫
router.beforeEach((to, from, next) => {
  console.log('to', to)
  console.log('to', to.path)
  const config = configStore()
  NProgress.start()
  document.title = to.meta && to.meta.title ? to.meta.title + ' - ' + config.configInfo.name : config.configInfo.name
  const loginInfo = localStorage.getItem('loginInfo')
  if (to.path !== '/login') {
    if (!loginInfo || loginInfo === '') {
      localStorage.setItem('logout', JSON.stringify({ type: 'warning', message: '没有权限，请先登录！' }))
      next('/')
    } else {
      next()
    }
  } else {
    if (to.path === '/error') {
      next('/')
    } else {
      next()
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})

export default router
