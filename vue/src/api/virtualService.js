import axios from '@/utils/axios'
// 获取虚拟设备列表
const getDeviceList = (params) => {
  return axios.request({
    url: '/virtual/devices',
    method: 'get',
    data: params,
  })
}
// 添加虚拟设备
const addDevice = (params) => {
  return axios.request({
    url: '/virtual/device',
    method: 'post',
    data: params,
  })
}
// 编辑虚拟设备
const editDevice = (params) => {
  return axios.request({
    url: '/virtual/device',
    method: 'put',
    data: params,
  })
}
// 删除虚拟设备
const deleteDevice = (params) => {
  return axios.request({
    url: '/virtual/devices',
    method: 'delete',
    data: params,
  })
}
// 获取虚拟设备属性列表
const getPropertyList = (params) => {
  return axios.request({
    url: '/virtual/properties',
    method: 'get',
    data: params,
  })
}
// 添加虚拟设备属性
const addProperty = (params) => {
  return axios.request({
    url: '/virtual/property',
    method: 'post',
    data: params,
  })
}
// 编辑虚拟设备属性
const editProperty = (params) => {
  return axios.request({
    url: '/virtual/property',
    method: 'put',
    data: params,
  })
}
// 批量删除虚拟设备属性
const deletePorperties = (params) => {
  return axios.request({
    url: '/virtual/properties',
    method: 'delete',
    data: params,
  })
}
export default {
  getDeviceList,
  addDevice,
  editDevice,
  deleteDevice,
  getPropertyList,
  addProperty,
  editProperty,
  deletePorperties,
}
