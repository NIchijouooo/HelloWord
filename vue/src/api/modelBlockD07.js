import axios from '@/utils/axios'

// 获取模型命令
const getDeviceModelBlockD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/cmd',
    method: 'get',
    data: params,
  })
}
// 添加模型命令
const addDeviceModelBlockD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/cmd',
    method: 'post',
    data: params,
  })
}
// 修改模型命令
const editDeviceModelBlockD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/cmd',
    method: 'put',
    data: params,
  })
}

// 删除模型命令
const deleteDeviceModelBlockD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/cmd',
    method: 'delete',
    data: params,
  })
}
// 批量导入模型命令
const importDeviceModelBlockD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/block/xlsx',
    method: 'post',
    data: params,
  })
}
// 批量导出模型命令
const exportDeviceModelBlockD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/block/xlsx',
    method: 'get',
    data: params,
  })
}
// 获取模型命令参数
const getDeviceModelBlockPropertyD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/cmd/params',
    method: 'get',
    data: params,
  })
}
// 添加模型命令参数
const addDeviceModelBlockPropertyD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/cmd/param',
    method: 'post',
    data: params,
  })
}
// 修改模型命令参数
const editDeviceModelBlockPropertyD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/cmd/param',
    method: 'put',
    data: params,
  })
}

// 删除模型命令参数
const deleteDeviceModelBlockPropertyD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/cmd/params',
    method: 'delete',
    data: params,
  })
}

// 批量导入模型命令参数
const importDeviceModelBlockPropertyD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/cmd/param/xlsx',
    method: 'post',
    data: params,
  })
}
// 批量导出模型命令参数
const exportDeviceModelBlockPropertyD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/cmd/param/xlsx',
    method: 'get',
    data: params,
  })
}
// 获取模型所有属性
const getDeviceModelPropertyListD07 = (params) => {
  return axios.request({
    url: '/tsl/dlt64507/properties',
    method: 'get',
    data: params,
  })
}

//
export default {
  getDeviceModelBlockD07,
  addDeviceModelBlockD07,
  editDeviceModelBlockD07,
  deleteDeviceModelBlockD07,
  importDeviceModelBlockD07,
  exportDeviceModelBlockD07,

  getDeviceModelBlockPropertyD07,
  addDeviceModelBlockPropertyD07,
  editDeviceModelBlockPropertyD07,
  deleteDeviceModelBlockPropertyD07,
  importDeviceModelBlockPropertyD07,
  exportDeviceModelBlockPropertyD07,

  getDeviceModelPropertyListD07,
}
