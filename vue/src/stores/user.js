export const userStore = defineStore({
  id: 'user',
  state: () => ({
    // 用户信息
    userInfo: null,
    // 路由信息
    routers: [],
    // 密钥
    token: '',
  }),
  actions: {
    setUserInfo(userInfo) {
      this.userInfo = userInfo
    },
    // 登录
    setLoginInfo(userInfo, routers, token) {
      this.userInfo = userInfo
      this.routers = routers
      this.token = token
    },
    // 注销
    logout() {
      this.clearUserInfo()
      localStorage.setItem('loginInfo', '')
      localStorage.setItem('logout', JSON.stringify({ type: 'warning', message: '注销成功！' }))
    },
    // 修改密码
    updatePwd() {
      this.clearUserInfo()
      localStorage.setItem('loginInfo', '')
      localStorage.setItem('logout', JSON.stringify({ type: 'success', message: '密码修改成功，请重新登录！' }))
    },
    clearUserInfo() {
      this.userInfo = null
      this.routers = []
      this.token = ''
    },
  },
})
