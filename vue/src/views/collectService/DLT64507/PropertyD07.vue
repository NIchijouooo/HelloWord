<template>
<div class="main-container">
  <div class="main">
    <div class="search-bar">
      <el-form :inline="true" ref="searchFormRef" status-icon label-width="90px">
        <el-form-item style="margin-left: 20px;">
          <el-button type="primary" plain @click="toDeviceModel()">
            <el-icon class="el-input__icon"><back /></el-icon>
            返回设备模型
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" plain class="right-btn" @click="exportDPS()">
            <el-icon class="el-input__icon"><upload /></el-icon>
            导出模型属性
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="search-bar" style="display: flex;">
      <div class="title" style="position: relative;margin-right: 40px;justify-content: flex-start;padding: 0px 0px;height: 40px;">
        <div class="tName">{{ props.curDeviceModel.label }}：属性列表</div>
      </div>
      <el-form :inline="true" ref="searchFormRef2" status-icon label-width="90px">
        <el-form-item label="">
          <el-input style="width: 200px" placeholder="请输入属性名称" v-model="ctxData.deviceModelProperty">
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button style="color: #fff" color="#2EA554" class="right-btn" @click="refresh('deviceModelProperty')">
            <el-icon class="btn-icon">
              <Icon name="local-refresh" size="14px" color="#ffffff" />
            </el-icon>
            刷新
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="content" ref="contentRef">
      <el-table
        :data="filterDMPTableData"
        :cell-style="ctxData.cellStyle"
        :header-cell-style="ctxData.headerCellStyle"
        style="width: 100%"
        :max-height="ctxData.tableMaxHeight"
        stripe
      >
        <el-table-column type="index" width="60">
          <template #header> 序号 </template>
        </el-table-column>
        <el-table-column sortable prop="name" label="属性名称" width="auto" min-width="150" align="center"> </el-table-column>
        <el-table-column sortable prop="label" label="属性标签" width="auto" min-width="150" align="center"> </el-table-column>

        <el-table-column sortable prop="rulerId" label="数据标识" width="auto" min-width="150" align="center"> </el-table-column>
        <el-table-column sortable prop="format" label="数据格式" width="auto" min-width="150" align="center"> </el-table-column>
        <el-table-column sortable prop="len" label="数据长度" width="auto" min-width="150" align="center"> </el-table-column>
        <el-table-column sortable label="读写属性" width="auto" min-width="80" align="center">
          <template #default="scope">
            {{ ctxData.accessModeNames['am' + scope.row.accessMode] }}
          </template>
        </el-table-column>
        <el-table-column sortable prop="type" label="数据类型" width="auto" min-width="100" align="center">
          <template #default="scope">
            {{ ctxData.typeNames['t' + scope.row.type] }}
          </template>
        </el-table-column>
        <el-table-column sortable prop="unit" label="单位" width="auto" min-width="80" align="center" ></el-table-column>
        <el-table-column sortable prop="blockAddOffset" label="块偏移地址" width="auto" min-width="150" align="center"> </el-table-column>
        <el-table-column sortable prop="rulerAddOffset" label="标识偏移地址" width="auto" min-width="150" align="center"> </el-table-column>
        <el-table-column sortable prop="step" label="步长" width="auto" min-width="120" align="center" />
        <template #empty>
          <div>无数据</div>
        </template>
      </el-table>
      <div class="pagination" style="z-index: 9;">
        <el-pagination
          :current-page="ctxData.currentPage"
          :page-size="ctxData.pagesize"
          :page-sizes="[20, 50, 200, 500]"
          :total="filterTableDataPage.length"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
          background
          layout="total, sizes, prev, pager, next, jumper"
          style="margin-top: 46px"
        ></el-pagination>
      </div>
    </div>
  </div>
</div>
</template>
<script setup>
import { Search, Back, Upload } from '@element-plus/icons-vue'
import variables from 'styles/variables.module.scss'
import DeviceModelApi from 'api/deviceModel.js'
import { userStore } from 'stores/user'
const users = userStore()

const props = defineProps({
  curDeviceModel: {
    type: Object,
    default: {},
  },
})
console.log('id -> props', props)

// 返回设备模型
const emit = defineEmits(['changeShowFlag'])
const toDeviceModel = () => {
  console.log('toDeviceModel')
  emit('changeShowFlag')
}

const contentRef = ref(null)
const ctxData = reactive({
  headerCellStyle: {
    background: variables.primaryColor,
    color: variables.fontWhiteColor,
    height: '54px',
  },
  cellStyle: {
    height: '48px',
  },
  tableMaxHeight: 0,
  currentPage: 1, // 默认当前页是第一页
  pagesize: 20, // 每页数据个数
  propertyList: [],
  deviceModelProperty: '',

  typeOptions: [
    { label: 'uint32', value: 0 },
    { label: 'int32', value: 1 },
    { label: 'double', value: 2 },
    { label: 'string', value: 3 },
  ],

  typeNames: {
    t0: 'uint32',
    t1: 'int32',
    t2: 'double',
    t3: 'string',
  },

  accessModeNames: {
    am0: '只读',
    am1: '只写',
    am2: '读写',
  },

  pFlag: false, //属性对话框标识
  pTitle: '添加属性', //属性对话框titleName
  propertyForm: {
    name: '', // 属性名称，只能是字母+数字的组合，不可以是中文
    label: '', // 属性标签
    rulerId:'',//数据标识
    format:'',//数据格式
    len:0,//数据长度
    accessMode: 0, // 读写属性
    type: 0,
    unit: '', // 单位
    blockAddOffset:0,
    rulerAddOffset:0,
  },
  psFlag: false,
  selectedProperties: [],
})
// 获取设备模型属性
const getDeviceModelPropertyList = (flag) => {
  const pData = {
    token: users.token,
    data: {
      name: props.curDeviceModel.name,
    },
  }
  DeviceModelApi.getDeviceModelProperty(pData).then(async (res) => {
    console.log('getDeviceModelProperty -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.propertyList = res.data
      if (flag === 1) {
        ElMessage({
          type: 'success',
          message: '刷新成功！',
        })
      }
    } else {
      showOneResMsg(res)
    }
    await nextTick(() => {
      ctxData.tableMaxHeight = contentRef.value.clientHeight - 34 - 36 - 132
    })
  })
}

getDeviceModelPropertyList()
const refresh = () => {
  getDeviceModelPropertyList(1)
}
const filterDMPTableData = computed(() => {
  console.log('ctxData.propertyList ->', ctxData.propertyList)
  return ctxData.propertyList
    .filter((item) => {
      var a = !ctxData.deviceModelProperty
      var b = item.name.toLowerCase().includes(ctxData.deviceModelProperty.toLowerCase())

      return a || b
    })
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  return ctxData.propertyList.filter((item) => {
    var a = !ctxData.deviceModelProperty
    var b = item.name.toLowerCase().includes(ctxData.deviceModelProperty.toLowerCase())

    return a || b
  })
})
// 处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
// 导出设备属性和服务
const exportDPS = () => {
  console.log('exportDPS')
  const pData = {
    token: users.token,
    responseType: 'blob',
    data: {
      name: props.curDeviceModel.name,
    },
  }
  DeviceModelApi.exportDeviceModelProptyAndService(pData).then((res) => {
    if (res && res.code === '1') {
      showOneResMsg(res)
      return
    }
    const blob = new Blob([res.blob])
    const fileName = res.fileName

    const elink = document.createElement('a')
    elink.download = fileName
    elink.style.display = 'none'
    elink.href = URL.createObjectURL(blob)
    document.body.appendChild(elink)
    elink.click()
    URL.revokeObjectURL(elink.href) // 释放URL 对象
    document.body.removeChild(elink)
  })
}

//显示单个res结果，code不等于 '0' 的message
const showOneResMsg = (res) => {
  ElMessage({
    type: 'error',
    message: res.message,
  })
}
/**
 * 处理返回的结果
 * @param {结果} res
 * @param {要调用的方法} doFunction
 */
// eslint-disable-next-line no-unused-vars
const handleResult = (res, doFunction) => {
  ElMessage({
    type: res.code === '0' ? 'success' : 'error',
    message: res.message,
  })
  if (res.code === '0' && doFunction) {
    doFunction()
  }
}
</script>
<style lang="scss" scoped>
@use 'styles/custom-scoped.scss' as *;
.form-title {
  margin-bottom: 12px;
}
.tName {
  line-height: 14px;
  font-size: 14px;
  border-left: 3px solid #3054eb;
  padding-left: 15px;
}
</style>
