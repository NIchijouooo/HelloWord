import axios from '@/utils/axios'
// 添加物模型
const addDeviceModel = (params) => {
  return axios.request({
    url: '/tsl',
    method: 'post',
    data: params,
  })
}
// 编辑物模型
const editDeviceModel = (params) => {
  return axios.request({
    url: '/tsl',
    method: 'put',
    data: params,
  })
}
// 删除物模型插件
const deleteDeviceModel = (params) => {
  return axios.request({
    url: '/tsl',
    method: 'delete',
    data: params,
  })
}
// 获取所有物模型
const getDeviceModelList = (params) => {
  return axios.request({
    url: '/tsls',
    method: 'get',
    data: params,
  })
}
// 获取单个物模型内容
const getDeviceModelContent = (params) => {
  return axios.request({
    url: '/tsl/contents',
    method: 'get',
    data: params,
  })
}
// 批量导入单个物模型的属性和服务
const importDeviceModelProptyAndService = (params) => {
  return axios.request({
    url: '/tsl/contents/xlsx',
    method: 'post',
    data: params,
  })
}
// 批量导出单个物模型的属性和服务
const exportDeviceModelProptyAndService = (params) => {
  return axios.request({
    url: '/tsl/contents/xlsx',
    method: 'get',
    data: params,
  })
}
// 同步物模型属性
const syncDeviceModelProperty = (params) => {
  return axios.request({
    url: '/tsl/contents/plugin/xlsx',
    method: 'get',
    data: params,
  })
}
// 添加物模型属性
const addDeviceModelProperty = (params) => {
  return axios.request({
    url: '/tsl/content/property',
    method: 'post',
    data: params,
  })
}
// 编辑物模型属性
const editDeviceModelProperty = (params) => {
  return axios.request({
    url: '/tsl/content/property',
    method: 'put',
    data: params,
  })
} // 删除物模型属性
const deleteDeviceModelProperty = (params) => {
  return axios.request({
    url: '/tsl/content/properties',
    method: 'delete',
    data: params,
  })
}
// 获取物模型属性
const getDeviceModelProperty = (params) => {
  return axios.request({
    url: '/tsl/content/properties',
    method: 'get',
    data: params,
  })
}

// 添加物模型服务
const addDeviceModelService = (params) => {
  return axios.request({
    url: '/tsl/content/service',
    method: 'post',
    data: params,
  })
}
// 编辑物模型服务
const editDeviceModelService = (params) => {
  return axios.request({
    url: '/tsl/content/service',
    method: 'put',
    data: params,
  })
} // 删除物模型服务
const deleteDeviceModelService = (params) => {
  return axios.request({
    url: '/tsl/content/service',
    method: 'delete',
    data: params,
  })
}
// 获取物模型服务
const getDeviceModelService = (params) => {
  return axios.request({
    url: '/tsl/content/service',
    method: 'get',
    data: params,
  })
}

//
export default {
  addDeviceModel,
  editDeviceModel,
  deleteDeviceModel,
  getDeviceModelList,
  getDeviceModelContent,
  importDeviceModelProptyAndService,
  exportDeviceModelProptyAndService,
  //属性
  syncDeviceModelProperty,
  addDeviceModelProperty,
  editDeviceModelProperty,
  deleteDeviceModelProperty,
  getDeviceModelProperty,
  //服务
  addDeviceModelService,
  editDeviceModelService,
  deleteDeviceModelService,
  getDeviceModelService,
}
