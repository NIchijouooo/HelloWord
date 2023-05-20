import axios from '@/utils/axios'
// 备份系统配置文件
const exportSystemConfig = (params) => {
  return axios.request({
    url: '/system/backup',
    method: 'post',
    data: params,
  })
}
// 还原系统备份文件
const importSystemConfig = (params) => {
  return axios.request({
    url: '/system/recover',
    method: 'post',
    data: params,
  })
}
// 升级系统
const updateSystem = (params) => {
  return axios.request({
    url: '/system/update',
    method: 'post',
    data: params,
  })
}
// 重启系统
const rebootSystem = (params) => {
  return axios.request({
    url: '/system/reboot/system',
    method: 'post',
    data: params,
  })
}
// 重启服务
const rebootServer = (params) => {
  return axios.request({
    url: '/system/reboot/service',
    method: 'post',
    data: params,
  })
}
// 备份系统配置文件到锐泰云平台
const backupRemote = (params) => {
  return axios.request({
    url: '/system/backup/remote/rt',
    method: 'get',
    data: params,
  })
}
// 获取网关在锐泰云平台中备份文件列表
const getRemoteFileInfo = (params) => {
  return axios.request({
    url: '/system/recover/remote/rt/files/info',
    method: 'get',
    data: params,
  })
}
// 获取网关在锐泰云平台中指定配置文件
const getRemoteFile = (params) => {
  return axios.request({
    url: '/system/recover/remote/rt/file',
    method: 'get',
    data: params,
  })
}
export default {
  exportSystemConfig,
  importSystemConfig,
  updateSystem,
  rebootSystem,
  rebootServer,
  backupRemote,
  getRemoteFileInfo,
  getRemoteFile,
}
