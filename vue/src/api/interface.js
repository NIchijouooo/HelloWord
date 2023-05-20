import axios from '@/utils/axios'
// 添加采集接口
const addInterface = (params) => {
  return axios.request({
    url: '/interface/collect',
    method: 'post',
    data: params,
  })
}
// 修改采集接口
const editInterface = (params) => {
  return axios.request({
    url: '/interface/collect',
    method: 'put',
    data: params,
  })
}
// 删除采集接口
const deleteInterface = (params) => {
  return axios.request({
    url: '/interface/collect',
    method: 'delete',
    data: params,
  })
}
// 获取采集接口列表
const getInterfaceList = (params) => {
  return axios.request({
    url: '/interface/collects',
    method: 'get',
    data: params,
  })
}

// 采集接口设备管理
// 获取单个采集接口下所有设备
const getCollDevices = (params) => {
  return axios.request({
    url: '/interface/collect',
    method: 'get',
    data: params,
  })
}
// 获取采集接口下的设备
const getAllCollDevices = (params) => {
  return axios.request({
    url: '/interface/collect/devices',
    method: 'post',
    data: params,
  })
}
// 添加采集接口下设备
const addCollDevice = (params) => {
  return axios.request({
    url: '/interface/collect/device',
    method: 'post',
    data: params,
  })
}
// 从CSV文件中批量增加节点
const addDeviceFromCSV = (params) => {
  return axios.request({
    url: '/interface/collect/devices/xlsx/import',
    method: 'post',
    data: params,
  })
}
// 批量导出设备
const exportDevice = (params) => {
  return axios.request({
    url: '/interface/collect/devices/xlsx/export',
    method: 'post',
    data: params,
  })
}
// 修改采集接口下设备
const editCollDevice = (params) => {
  return axios.request({
    url: '/interface/collect/device',
    method: 'put',
    data: params,
  })
}
// 批量删除采集接口下设备
const deleteCollDevices = (params) => {
  return axios.request({
    url: '/interface/collect/devices',
    method: 'delete',
    data: params,
  })
}
// 设备数据
// 获取某采集接口，某个设备的缓存数据
const getDeviceDataCache = (params) => {
  return axios.request({
    url: '/interface/collect/variable/cache',
    method: 'get',
    data: params,
  })
}
// 获取某采集接口，某个设备的实时数据
const getDeviceDataReal = (params) => {
  return axios.request({
    url: '/interface/collect/variable/real',
    method: 'get',
    data: params,
  })
}
// 获取某采集接口，某个设备的历史数据
// const getDeviceDataHis = (params) => {
//   return axios.request({
//     url: '/interface/collect/device/variable/real',
//     method: 'get',
//     data: params,
//   })
// }
// 调用设备服务接口
const invokeDeviceService = (params) => {
  return axios.request({
    url: '/interface/collect/node/service',
    method: 'post',
    data: params,
  })
}
export default {
  addInterface,
  editInterface,
  deleteInterface,
  getInterfaceList,
  getCollDevices,
  getAllCollDevices,
  addCollDevice,
  editCollDevice,
  deleteCollDevices,
  getDeviceDataCache,
  getDeviceDataReal,
  addDeviceFromCSV,
  exportDevice,
  invokeDeviceService,
}
