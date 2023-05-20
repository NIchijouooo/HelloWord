import axios from '@/utils/axios'
// 获取所有网卡偏好设置
const getNetworkList = (params) => {
  return axios.request({
    url: '/network/params',
    method: 'get',
    data: params,
  })
}

// 添加网卡偏好设置
const addNetwork = (params) => {
  return axios.request({
    url: '/network/param',
    method: 'post',
    data: params,
  })
}
// 编辑网卡偏好设置
const editNetwork = (params) => {
  return axios.request({
    url: '/network/param',
    method: 'put',
    data: params,
  })
}
// 删除网卡偏好设置
const deleteNetwork = (params) => {
  return axios.request({
    url: '/network/param',
    method: 'delete',
    data: params,
  })
}
// 获取所有通信接口
const getNetworkNames = (params) => {
  return axios.request({
    url: '/network/names',
    method: 'get',
    data: params,
  })
}
export default {
  getNetworkList,
  addNetwork,
  editNetwork,
  deleteNetwork,
  getNetworkNames,
}
