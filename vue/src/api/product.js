import axios from '@/utils/axios'
// 获取网关SN号
const getGatewaySN = (params) => {
  return axios.request({
    url: '/product/sn',
    method: 'get',
    data: params,
  })
}
// 设置网关SN号
const setGatewaySN = (params) => {
  return axios.request({
    url: '/product/sn',
    method: 'post',
    data: params,
  })
}
//
export default {
  getGatewaySN,
  setGatewaySN,
}
