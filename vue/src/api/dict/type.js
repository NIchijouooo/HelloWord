import axios from '@/utils/axios'

// 查询字典类型列表
const listType = (query) =>  {
  return axios.request({
    url: '/dictType/getDictTypeList',
    method: 'post',
    data: query
  })
}

// 查询字典类型详细
const getType = (query)  => {
  return axios.request({
    url: '/dictType/getDictTypeByID',
    method: 'post',
    data: query
  })
}

// 新增字典类型
const addType = (data)  => {
  return axios.request({
    url: '/dictType/createDictType',
    method: 'post',
    data: data
  })
}

// 修改字典类型
const updateType = (data) =>  {
  return axios({
    url: '/dictType/updateDictType',
    method: 'post',
    data: data
  })
}

// 删除字典类型
const delType = (data) =>  {
  return axios.request({
    url: '/dictType/deleteDictType?dictId=' + data.dictIds,
    method: 'post',
    data: data
  })
}

// 清理参数缓存
// const clearCache() {
//   return axios.request({
//     url: '/qianhai/system/dict/type/clearCache',
//     method: 'delete'
//   })
// }

// 导出字典类型
// const exportType(query) {
//   return axios.request({
//     url: '/qianhai/system/dict/type/export',
//     method: 'get',
//     params: query
//   })
// }

// 获取字典选择框列表
const optionselect = () => {
  return axios.request({
    url: '/dictType/getDictTypeList',
    method: 'get'
  })
}
// 回显数据字典
const selectDictLabel = (datas, value)  => {
  var actions = [];
  Object.keys(datas).some((key) => {
    if (datas[key].dictValue == ('' + value)) {
      actions.push(datas[key].dictLabel);
      return true;
    }
  })
  return actions.join('');
}

export default {
  listType,
  getType,
  addType,
  updateType,
  delType,
  optionselect,
  selectDictLabel
}
