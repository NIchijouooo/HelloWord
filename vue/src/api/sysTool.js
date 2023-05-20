import axios from '@/utils/axios'
// 修改NTP服务
const updateNTP = (params) => {
  return axios.request({
    url: '/ntp/param',
    method: 'put',
    data: params,
  })
}
// 获取NTP服务
const getNTPInfo = (params) => {
  return axios.request({
    url: '/ntp/param',
    method: 'get',
    data: params,
  })
}
// 立即校时
const sendTimingCmd = (params) => {
  return axios.request({
    url: '/ntp/cmd',
    method: 'post',
    data: params,
  })
}
// 发送ping命令
const sendPingCmd = (params) => {
  return axios.request({
    url: '/system/ping',
    method: 'post',
    data: params,
  })
}
// 发送通讯报文
const sendMessage = (params) => {
  return axios.request({
    url: '/system/commTool',
    method: 'post',
    data: params,
  })
}
export default {
  updateNTP,
  getNTPInfo,
  sendTimingCmd,
  sendPingCmd,
  sendMessage,
}
