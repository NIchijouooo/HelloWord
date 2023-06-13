import axios from '@/utils/axios'

// 查询字典数据列表
const listData = (data) =>  {
  return axios.request({
    url: '/dictData/getDictDataList',
    method: 'post',
    data: data
  })
}

// 查询字典数据详细
const getData = (data) => {
  return axios.request({
    url: '/dictData/getDictDataByID',
    method: 'post',
    data: data,
  })
}

// 根据字典类型查询字典数据信息
const getDicts = (data) => {
  return axios.request({
    url: '/dictData/getDictTypeListByDictTypeId',
    method: 'post',
    data: data,
  })
}

// 新增字典数据
const addData = (data) => {
  return axios.request({
    url: '/dictData/createDictData',
    method: 'post',
    data: data
  })
}

// 修改字典数据
const updateData = (data) => {
  return axios.request({
    url: '/dictData/updateDictData',
    method: 'post',
    data: data
  })
}

// 删除字典数据
const delData = (data) => {
  return axios.request({
    url: '/dictData/deleteDictData',
    method: 'post',
    data: data
  })
}

// 导出字典数据
// const exportData(query) {
//   return axios.request({
//     url: '/qianhai/system/dict/data/export',
//     method: 'get',
//     params: query
//   })
// }

export default {
  listData,
  getData,
  getDicts,
  addData,
  updateData,
  delData,
}
