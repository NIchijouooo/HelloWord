import axios from '@/utils/axios'
// 登录
const login = (params) => {
  return axios.request({
    url: '/account/login',
    method: 'post',
    data: params,
  })
}
// 修改密码
const updatePassword = (params) => {
  return axios.request({
    url: '/account/role/password',
    method: 'put',
    data: params,
  })
}

export default {
  login,
  updatePassword,
}
