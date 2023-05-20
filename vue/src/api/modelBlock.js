import axios from '@/utils/axios'

// 获取模型命令
const getDeviceModelBlock = (params) => {
  return axios.request({
    url: '/tsl/modbus/cmd',
    method: 'get',
    data: params,
  })
}
// 添加模型命令
const addDeviceModelBlock = (params) => {
  return axios.request({
    url: '/tsl/modbus/cmd',
    method: 'post',
    data: params,
  })
}
// 修改模型命令
const editDeviceModelBlock = (params) => {
  return axios.request({
    url: '/tsl/modbus/cmd',
    method: 'put',
    data: params,
  })
}

// 删除模型命令
const deleteDeviceModelBlock = (params) => {
  return axios.request({
    url: '/tsl/modbus/cmd',
    method: 'delete',
    data: params,
  })
}
// 批量导入模型命令
const importDeviceModelBlock = (params) => {
  return axios.request({
    url: '/tsl/modbus/block/xlsx',
    method: 'post',
    data: params,
  })
}
// 批量导出模型命令
const exportDeviceModelBlock = (params) => {
  return axios.request({
    url: '/tsl/modbus/block/xlsx',
    method: 'get',
    data: params,
  })
}
// 获取模型命令参数
const getDeviceModelBlockProperty = (params) => {
  return axios.request({
    url: '/tsl/modbus/cmd/params',
    method: 'get',
    data: params,
  })
}
// 添加模型命令参数
const addDeviceModelBlockProperty = (params) => {
  return axios.request({
    url: '/tsl/modbus/cmd/param',
    method: 'post',
    data: params,
  })
}
// 修改模型命令参数
const editDeviceModelBlockProperty = (params) => {
  return axios.request({
    url: '/tsl/modbus/cmd/param',
    method: 'put',
    data: params,
  })
}

// 删除模型命令参数
const deleteDeviceModelBlockProperty = (params) => {
  return axios.request({
    url: '/tsl/modbus/cmd/params',
    method: 'delete',
    data: params,
  })
}

// 批量导入模型命令参数
const importDeviceModelBlockProperty = (params) => {
  return axios.request({
    url: '/tsl/modbus/cmd/param/xlsx',
    method: 'post',
    data: params,
  })
}
// 批量导出模型命令参数
const exportDeviceModelBlockProperty = (params) => {
  return axios.request({
    url: '/tsl/modbus/cmd/param/xlsx',
    method: 'get',
    data: params,
  })
}
// 获取模型所有属性
const getDeviceModelPropertyList = (params) => {
  return axios.request({
    url: '/tsl/modbus/properties',
    method: 'get',
    data: params,
  })
}

//
export default {
  getDeviceModelBlock,
  addDeviceModelBlock,
  editDeviceModelBlock,
  deleteDeviceModelBlock,
  importDeviceModelBlock,
  exportDeviceModelBlock,

  getDeviceModelBlockProperty,
  addDeviceModelBlockProperty,
  editDeviceModelBlockProperty,
  deleteDeviceModelBlockProperty,
  importDeviceModelBlockProperty,
  exportDeviceModelBlockProperty,

  getDeviceModelPropertyList,
}
