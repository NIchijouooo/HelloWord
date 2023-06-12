import axios from '@/utils/axios'
// 添加通讯接口
const addCommInterface = (params) => {
  return axios.request({
    url: '/interface/communication',
    method: 'post',
    data: params,
  })
}
// 编辑通讯接口
const editCommInterface = (params) => {
  return axios.request({
    url: '/interface/communication',
    method: 'put',
    data: params,
  })
}
// 删除通讯接口
const deleteCommInterface = (params) => {
  return axios.request({
    url: '/interface/communication',
    method: 'delete',
    data: params,
  })
}
// 获取所有通信接口
const getCommInterfaceList = (params) => {
  return axios.request({
    url: '/interface/communication',
    method: 'get',
    data: params,
  })
}
// 获取当前网关通讯协议
const getCommProtocolList = (params) => {
  return axios.request({
    url: '/interface/communication/protocol',
    method: 'get',
    data: params,
  })
}
export default {
  addCommInterface,
  editCommInterface,
  deleteCommInterface,
  getCommInterfaceList,
  getCommProtocolList,
}
