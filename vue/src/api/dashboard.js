import axios from '@/utils/axios'
// 获取系统参数
const getSysParams = (params) => {
  return axios.request({
    url: '/system/params',
    method: 'get',
    data: params,
  })
}
const getCpuList = (params) => {
  return axios.request({
    url: '/system/params/cpu',
    method: 'get',
    data: params,
  })
}
// 获取系统内存运行数据
const getMemoryList = (params) => {
  return axios.request({
    url: '/system/params/memory',
    method: 'get',
    data: params,
  })
}
// 获取系统硬盘运行数据
const getDiskList = (params) => {
  return axios.request({
    url: '/system/params/disk',
    method: 'get',
    data: params,
  })
}
// 获取系统设备丢包数据
const getDevicePacketLossList = (params) => {
  return axios.request({
    url: '/system/params/device/packetLoss',
    method: 'get',
    data: params,
  })
}
// 获取系统设备在线数据
const getDeviceOnlineList = (params) => {
  return axios.request({
    url: '/system/params/device/online',
    method: 'get',
    data: params,
  })
}
export default {
  getSysParams,
  getCpuList,
  getMemoryList,
  getDiskList,
  getDevicePacketLossList,
  getDeviceOnlineList,
}
