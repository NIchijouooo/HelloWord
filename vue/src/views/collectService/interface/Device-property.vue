<template>
  <div class="main">
    <div class="search-bar">
      <el-form :inline="true" ref="searchFormRef" status-icon label-width="90px">
        <el-form-item style="margin-left: 20px;">
          <el-button type="primary" plain @click="toDevice()" style="margin-right: 20px">
            <el-icon class="el-input__icon"><back /></el-icon>
            返回采集设备
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="search-bar" style="display: flex;">
      <div class="title" style="position: relative;margin-right: 40px;justify-content: flex-start;padding: 0px 0px;height: 40px;">
        <div class="tName">{{ props.curDevice.name }}({{ props.curDevice.label }})</div>
      </div>
      <el-form :inline="true" ref="searchFormRef2" status-icon label-width="90px">
        <el-form-item label="">
          <el-input style="width: 200px" placeholder="请输入 名称/标签 过滤" clearable v-model="ctxData.propertyInfo">
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" plain class="right-btn" @click="showProperties()">
            <el-icon class="el-input__icon"><edit /></el-icon>
            写属性
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button
            style="color: #fff"
            color="#2EA554"
            :loading="ctxData.isLoading"
            class="right-btn"
            @click="refresh()"
          >
            <el-icon class="btn-icon">
              <Icon name="local-refresh" size="14px" color="#ffffff" />
            </el-icon>
            刷新
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="content" ref="contentRef" style="top: 136px">
      <el-table
        :data="filterTableData"
        :cell-style="ctxData.cellStyle"
        :header-cell-style="ctxData.headerCellStyle"
        height="660"
        style="width: 100%"
        stripe
      >
        <el-table-column prop="index" label="序号" width="55" />
        <el-table-column prop="name" label="变量名称" width="auto" min-width="180" align="center"> </el-table-column>
        <el-table-column prop="label" label="变量标签" width="auto" min-width="180" align="center"> </el-table-column>
        <el-table-column prop="type" label="变量类型" width="auto" min-width="120" align="center"> </el-table-column>
        <el-table-column prop="value" label="变量值" width="auto" min-width="200" align="center"> </el-table-column>
        <el-table-column prop="unit" label="单位" width="auto" min-width="200" align="center"> </el-table-column>
        <el-table-column prop="explain" label="说明" width="auto" min-width="200" align="center"> </el-table-column>
        <el-table-column prop="timestamp" label="实测时间" width="auto" min-width="220" align="center">
        </el-table-column>
        <el-table-column label="操作" width="auto" min-width="200" align="center" fixed="right">
          <template #default="scope">
            <el-button @click="showHisData(scope.row)" text type="success">查看历史数据</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <div>无数据</div>
        </template>
      </el-table>
      <div class="pagination">
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
    <!-- 写变量弹出页 -->
    <el-dialog
      v-model="ctxData.pFlag"
      title="写属性"
      width="1000px"
      :before-close="handleCloseProperty"
      :close-on-click-modal="false"
    >
      <div class="dialog-content" style="min-height: 408px; overflow: unset; padding: 0">
        <div class="dialog-content-head">
          <el-input
            placeholder="请输入属性名称或标签"
            style="width: 200px"
            v-model="ctxData.pInfo"
            @change="changePInfo()"
          >
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" class="right-btn" @click="sendCmd('Write')">
            <el-icon class="el-input__icon"><edit /></el-icon>
            写属性
          </el-button>
        </div>
        <div class="dialog-content-content">
          <!--

            
          -->
          <el-table
            :data="ctxData.pTableData"
            @selection-change="handleSelectionChange"
            :header-cell-style="ctxData.userHeadStyle"
            :cell-style="ctxData.userCellStyle"
            style="border: 1px solid #c0c4cc"
            max-height="290"
            stripe
          >
            <el-table-column type="selection" width="55"> </el-table-column>
            <el-table-column prop="name" label="属性名称" min-width="120" align="center"> </el-table-column>
            <el-table-column prop="label" label="属性标签" min-width="140" align="center"> </el-table-column>
            <el-table-column label="写入值" min-width="100" align="center">
              <template #default="scope">
                <el-input placeholder="请输入写入值" v-model="scope.row.sendValue"> </el-input>
              </template>
            </el-table-column>
            <el-table-column label="返回结果" min-width="250" align="center">
              <template #default="scope">
                <el-tag v-show="scope.row.result.Code === '0'" :type="'success'"> 操作成功 </el-tag>
                <el-tag v-show="scope.row.result.Code === '1'" :type="'danger'"> 操作失败 </el-tag>
                <span v-show="scope.row.result.Code === ''">-</span>
                <span v-show="scope.row.result.Code === '1'">
                  {{ scope.row.result.Message }}
                </span>
              </template>
            </el-table-column>
          </el-table>
          <div class="pagination dialog-pagination">
            <el-pagination
              :current-page="ctxData.pCurrentPage"
              :page-size="ctxData.pPagesize"
              :page-sizes="[5, 10, 20, 50]"
              :total="ctxData.pTableData1.length"
              @current-change="handlePCurrentChange"
              @size-change="handlePSizeChange"
              background
              layout="total, sizes, prev, pager, next, jumper"
            ></el-pagination>
          </div>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelPorperty()">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { Search, Back, Edit } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import variables from 'styles/variables.module.scss'
import InterfaceApi from 'api/interface.js'
import { userStore } from 'stores/user'
import DeviceModelApi from 'api/deviceModel.js'
const users = userStore()
const props = defineProps({
  curDevice: {
    type: Object,
    default: {},
  },
  collInterfaceName: {
    type: String,
    default: '',
  },
})
console.log('id -> props', props)

const emit = defineEmits(['changeDpFlag'])
const toDevice = () => {
  emit('changeDpFlag')
}
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
  propertyTableData: [],
  propertyInfo: '',
  pFlag: false,
  pInfo: '',
  pCurrentPage: 1, // 默认当前页是第一页
  pPagesize: 5, // 每页数据个数
  userHeadStyle: {
    color: variables.headerTextColor,
    borderColor: '#C0C4CC',
    height: '48px',
  },
  userCellStyle: {
    borderColor: '#C0C4CC',
    height: '48px',
  },
  pTableData: [],
  pTableData1: [],
  selectPorperties: [],
  isLoading: false,
})
const contentRef = ref(null)
// 获取采集接口下的设备属性
const getDeviceDataReal = (flag) => {
  const pData = {
    token: users.token,
    data: {
      collInterfaceName: props.collInterfaceName,
      deviceName: props.curDevice.name,
    },
  }
  ctxData.isLoading = true
  InterfaceApi.getDeviceDataReal(pData).then(async (res) => {
    console.log('getDeviceDataReal -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.propertyTableData = res.data
      if (flag === 1) {
        ElMessage.success('刷新成功！')
      }
    } else {
      showOneResMsg(res)
    }

    ctxData.isLoading = false
    console.log('getDeviceDataReal -> ctxData.propertyTableData', ctxData.propertyTableData)
    await nextTick(() => {
      ctxData.tableMaxHeight = contentRef.value.clientHeight - 34 - 36 - 22
    })
  })
}
getDeviceDataReal()
const refresh = () => {
  getDeviceDataReal(1)
}
// 过滤表格数据
const filterTableData = computed(() => {
  let propertyInfo = ctxData.propertyInfo
  return ctxData.propertyTableData
    .filter(
      (item) =>
        !propertyInfo ||
        item.name.toLowerCase().includes(propertyInfo.toLowerCase()) ||
        item.label.toLowerCase().includes(propertyInfo.toLowerCase())
    )
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  let propertyInfo = ctxData.propertyInfo
  return ctxData.propertyTableData.filter(
    (item) =>
      !propertyInfo ||
      item.name.toLowerCase().includes(propertyInfo.toLowerCase()) ||
      item.label.toLowerCase().includes(propertyInfo.toLowerCase())
  )
})
// 处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}

const showHisData = (row) => {
  //
  ElMessage.info('功能完善中...')
}
//展示写变量页面
const showProperties = () => {
  const pData = {
    token: users.token,
    data: {
      name: props.curDevice.tsl,
    },
  }
  DeviceModelApi.getDeviceModelProperty(pData).then((res) => {
    if (res.code === '0') {
      ctxData.pTableData = []
      ctxData.pTableData1 = []
      const resData = res.data.filter((item) => item.accessMode !== 0)
      resData.forEach((item) => {
        const temp = {
          name: item.name,
          label: item.label,
          sendValue: '',
          result: {
            Code: '',
            Message: '',
          },
        }
        ctxData.pTableData.push(temp)
        ctxData.pTableData1.push(temp)
      })
      console.log('ctxData.pTableData -> ', ctxData.pTableData)
      ctxData.pFlag = true
    } else {
      showOneResMsg(res)
    }
  })
}
const cancelPorperty = () => {
  ctxData.pFlag = false
}
const handleCloseProperty = (done) => {
  cancelPorperty()
}

const changePInfo = () => {
  ctxData.pTableData = ctxData.pTableData1.filter(
    (item) => !ctxData.pInfo || item.name.includes(ctxData.pInfo) || item.label.includes(ctxData.pInfo)
  )
}
// 处理当前页变化
const handlePCurrentChange = (val) => {
  ctxData.pCurrentPage = val
  ctxData.pTableData = ctxData.pTableData1.slice(
    (ctxData.pCurrentPage - 1) * ctxData.pPagesize,
    ctxData.pCurrentPage * ctxData.pPagesize
  )
}
// 处理每页大小变化
const handlePSizeChange = (val) => {
  ctxData.pPagesize = val
  ctxData.pTableData = ctxData.pTableData1.slice(
    (ctxData.pCurrentPage - 1) * ctxData.pPagesize,
    ctxData.pCurrentPage * ctxData.pPagesize
  )
}
const handleSelectionChange = (val) => {
  console.log('handleSelectionChange val = ', val)
  console.log('handleSelectionChange pTableData = ', ctxData.pTableData)
  ctxData.selectPorperties = val
}

const sendCmd = (cmd) => {
  console.log('sendCmd', ctxData.selectPorperties)
  if (ctxData.selectPorperties.length < 1) {
    ElMessage.warning('请至少选择一个属性操作')
    return
  }
  const serviceParam = {}
  ctxData.selectPorperties.forEach((item) => {
    serviceParam[item.name] = item.sendValue
  })
  //serviceParam[ctxData.selectPorperties[0].name] = ctxData.selectPorperties[0].sendValue
  console.log(serviceParam)
  const pData = {
    token: users.token,
    data: {
      collInterfaceName: props.collInterfaceName,
      deviceName: props.curDevice.name,
      serviceName: 'SetVariables',
      serviceParam: serviceParam,
    },
  }
  if (cmd === 'Read') {
  } else {
    InterfaceApi.invokeDeviceService(pData).then((res) => {
      ElMessage.success('写属性命令下发成功！')
      ctxData.pTableData.forEach((item) => {
        ctxData.selectPorperties.forEach((property) => {
          if (item.name === property.name) {
            item.result = res
          }
        })
      })
    })
  }
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
.dialog-content-head {
  position: relative;
  width: 100%;
  height: 60px;
  display: flex;
  justify-content: space-between;
  box-sizing: border-box;
  padding-bottom: 24px;
}
.dialog-content-content {
  position: relative;
  width: 100%;
  text-align: center;
  height: 348px;
}
</style>
