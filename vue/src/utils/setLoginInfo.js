import router from '@/router'
const setLoginInfoToStore = (users) => {
  const loginInfo = localStorage.getItem('loginInfo')
  if (loginInfo && loginInfo != '') {
    const { username, permissions, token } = JSON.parse(loginInfo)
    const userInfo = {
      userName: username,
    }
    users.setLoginInfo(userInfo, permissions, token)
  } else {
    localStorage.setItem('logout', JSON.stringify({ type: 'warning', message: '没有权限，请先登录！' }))
    router.push('/')
  }
}

export default setLoginInfoToStore
