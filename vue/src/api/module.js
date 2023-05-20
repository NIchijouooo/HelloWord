import axios from '@/utils/axios'
// 添加模块
const addModule = (params) => {
  return axios.request({
    url: '/network/mobile',
    method: 'post',
    data: params,
  })
}
// 编辑模块
const editModule = (params) => {
  return axios.request({
    url: '/network/mobile',
    method: 'put',
    data: params,
  })
}
// 删除模块
const deleteModule = (params) => {
  return axios.request({
    url: '/network/mobile',
    method: 'delete',
    data: params,
  })
}
// 获取所有模块
const getModuleList = (params) => {
  return axios.request({
    url: '/network/mobiles',
    method: 'get',
    data: params,
  })
}
export default {
  addModule,
  editModule,
  deleteModule,
  getModuleList,
}
