import axios from '@/utils/axios'
// 获取所有上报服务列表
const getGatewayList = (params) => {
  return axios.request({
    url: '/service/report/gateways',
    method: 'get',
    data: params,
  })
}

// 添加上报服务
const addGateway = (params) => {
  return axios.request({
    url: '/service/report/gateway',
    method: 'post',
    data: params,
  })
}
// 编辑上报服务
const editGateway = (params) => {
  return axios.request({
    url: '/service/report/gateway',
    method: 'put',
    data: params,
  })
}
// 删除上报服务
const deleteGateway = (params) => {
  return axios.request({
    url: '/service/report/gateway',
    method: 'delete',
    data: params,
  })
}

// 获取某个上报服务的所有设备列表
const getDeviceByServiceIdList = (params) => {
  return axios.request({
    url: '/service/report/devices',
    method: 'get',
    data: params,
  })
}

// 添加上报设备
const addDevice = (params) => {
  return axios.request({
    url: '/service/report/device',
    method: 'post',
    data: params,
  })
}
// 编辑上报设备
const editDevice = (params) => {
  return axios.request({
    url: '/service/report/device',
    method: 'put',
    data: params,
  })
}
// 删除上报设备
const deleteDevice = (params) => {
  return axios.request({
    url: '/service/report/devices',
    method: 'delete',
    data: params,
  })
}
// 批量添加上报设备
const addDeviceFromCSV = (params) => {
  return axios.request({
    url: '/service/report/devices/xlsx',
    method: 'post',
    data: params,
  })
}
//导出设备
const exportDevice = (params) => {
  return axios.request({
    url: '/service/report/devices/xlsx',
    method: 'get',
    data: params,
  })
}

// 获取某个上报服务的所有寄存器
const getRegByServiceIdList = (params) => {
  return axios.request({
    url: '/service/report/mbTCP/registers',
    method: 'get',
    data: params,
  })
}

// 添加寄存器
const addReg = (params) => {
  return axios.request({
    url: '/service/report/mbTCP/register',
    method: 'post',
    data: params,
  })
}
// 编辑寄存器
const editReg = (params) => {
  return axios.request({
    url: '/service/report/mbTCP/register',
    method: 'put',
    data: params,
  })
}
// 删除寄存器
const deleteReg = (params) => {
  return axios.request({
    url: '/service/report/mbTCP/registers',
    method: 'delete',
    data: params,
  })
}

// 导入寄存器
const importReg = (params) => {
  return axios.request({
    url: '/service/report/mbTCP/register/xlsx',
    method: 'post',
    data: params,
  })
}
// 导出寄存器
const exportReg = (params) => {
  return axios.request({
    url: '/service/report/mbTCP/register/xlsx',
    method: 'get',
    data: params,
  })
}
// 获取当前网关上报协议
const getReportProtocolList = (params) => {
  return axios.request({
    url: '/service/report/protocol',
    method: 'get',
    data: params,
  })
}
//
export default {
  getGatewayList,
  addGateway,
  editGateway,
  deleteGateway,
  getDeviceByServiceIdList,
  addDevice,
  editDevice,
  deleteDevice,
  addDeviceFromCSV,
  exportDevice,
  getRegByServiceIdList,
  addReg,
  editReg,
  deleteReg,
  importReg,
  exportReg,
  getReportProtocolList,
}
