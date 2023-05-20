import axios from '@/utils/axios'
// 获取所有上报模型列表
const getModelList = (params) => {
  return axios.request({
    url: '/service/report/models',
    method: 'get',
    data: params,
  })
}

// 添加上报模型
const addModel = (params) => {
  return axios.request({
    url: '/service/report/model',
    method: 'post',
    data: params,
  })
}
// 编辑上报模型
const editModel = (params) => {
  return axios.request({
    url: '/service/report/model',
    method: 'put',
    data: params,
  })
}
// 删除上报模型
const deleteModel = (params) => {
  return axios.request({
    url: '/service/report/model',
    method: 'delete',
    data: params,
  })
}

// 获取某个上报模型的所有属性列表
const getPropertiesByModelIdList = (params) => {
  return axios.request({
    url: '/service/report/model/properties',
    method: 'get',
    data: params,
  })
}

// 添加上报模型属性
const addProperty = (params) => {
  return axios.request({
    url: '/service/report/model/property',
    method: 'post',
    data: params,
  })
}
// 编辑上报模型属性
const editProperty = (params) => {
  return axios.request({
    url: '/service/report/model/property',
    method: 'put',
    data: params,
  })
}
// 删除上报模型属性
const deleteProperty = (params) => {
  return axios.request({
    url: '/service/report/model/properties',
    method: 'delete',
    data: params,
  })
}
// 批量添加上报模型属性
const addPropertyFromCSV = (params) => {
  return axios.request({
    url: '/service/report/model/properties/xlsx',
    method: 'post',
    data: params,
  })
}
//导出上报模型属性
const exportProperty = (params) => {
  return axios.request({
    url: '/service/report/model/properties',
    method: 'get',
    data: params,
  })
}

const reportNodes = (params) => {
  return axios.request({
    url: '/service/report/device/cmd/report',
    method: 'post',
    data: params,
  })
}
//
export default {
  getModelList,
  addModel,
  editModel,
  deleteModel,
  getPropertiesByModelIdList,
  addProperty,
  editProperty,
  deleteProperty,
  addPropertyFromCSV,
  exportProperty,
  reportNodes,
}
