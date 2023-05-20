import axios from '@/utils/axios'
// 导入单个物模型插件
const importOnePlugin = (params) => {
  return axios.request({
    url: '/tsl/plugin/file',
    method: 'post',
    data: params,
  })
}
// 导出单个物模型插件
const exportOnePlugin = (params) => {
  return axios.request({
    url: '/tsl/plugin/file',
    method: 'get',
    data: params,
  })
}
// 删除物模型插件
const deletePlugin = (params) => {
  return axios.request({
    url: '/tsl/plugin/file',
    method: 'delete',
    data: params,
  })
}
// 获取所有物模型插件
const getPlugin = (params) => {
  return axios.request({
    url: '/tsl/plugin/param',
    method: 'get',
    data: params,
  })
}
// 获取mock中的数据
const getTestList = (params) => {
  return axios.request({
    url: '/api/test01',
    method: 'get',
    data: params,
  })
}
export default {
  importOnePlugin,
  exportOnePlugin,
  deletePlugin,
  getPlugin,
  getTestList,
}
