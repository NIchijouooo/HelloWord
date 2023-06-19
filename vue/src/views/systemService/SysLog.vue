<template>
  <div class="SysLog">
    <div class="title">
      <el-check-tag size="large" :checked="ctxData.logType === 0" @change="onChange(0)">采集日志</el-check-tag>
      <el-check-tag size="large" :checked="ctxData.logType === 1" @change="onChange(1)">上报日志</el-check-tag>
      <el-check-tag size="large" :checked="ctxData.logType === 2" @change="onChange(2)">系统日志</el-check-tag>
    </div>
    <div class="content">
      <div class="cTitle">
        <div class="tName">{{ ctxData.logName[ctxData.logType] }}</div>
        <div class="option">
          <el-select
            v-model="ctxData.collectName"
            v-if="ctxData.logType === 0"
            style="width: 100%"
            placeholder="请选择采集接口"
          >
            <el-option
              v-for="item in ctxData.collectList"
              :key="'collect_' + item.collInterfaceName"
              :label="item.collInterfaceName"
              :value="item.collInterfaceName"
            >
            </el-option>
          </el-select>
          <el-select
            v-model="ctxData.uploadName"
            v-if="ctxData.logType === 1"
            style="width: 100%"
            placeholder="请选择上报接口"
          >
            <el-option
              v-for="item in ctxData.uploadList"
              :key="'upload_' + item.serviceName"
              :label="item.serviceName"
              :value="item.serviceName"
            >
            </el-option>
          </el-select>
          <el-button v-if="ctxData.startStopFlag" @click="changeSSFlag(false)" class="right-btn" type="success" plain>
            <el-icon class="el-input__icon"><caret-left /></el-icon>
            开始
          </el-button>
          <el-button v-else class="right-btn" @click="changeSSFlag(true)" type="danger" plain>
            <el-icon class="el-input__icon"><close-bold /></el-icon>
            停止
          </el-button>
          <el-button class="right-btn" @click="clearDataTable()" type="warning" plain>
            <el-icon class="el-input__icon"><delete /></el-icon>
            清空
          </el-button>
          <el-button class="right-btn" @click="exportCsv()" type="primary" plain>
            <el-icon class="el-input__icon"><download /></el-icon>
            导出
          </el-button>
        </div>
      </div>
      <div class="cTable" ref="contentRef">
        <el-table
          :data="filterTableData"
          :cell-style="ctxData.cellStyle"
          :header-cell-style="ctxData.headerCellStyle"
          :max-height="ctxData.tableMaxHeight"
          style="width: 100%"
          stripe
        >
          <el-table-column prop="timeStamp" label="时间戳" width="auto" min-width="220" align="center">
          </el-table-column>
          <el-table-column label="数据方向" width="auto" min-width="150" align="center">
            <template #default="scope">
              <el-tag :type="scope.row.direction === 0 ? 'success' : 'warning'">{{
                ctxData.dataType[scope.row.direction]
              }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="label" label="数据标识" width="auto" min-width="150" align="center"> </el-table-column>
          <el-table-column prop="content" label="数据内容" width="auto" min-width="1000" align="left">
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
            :total="ctxData.dataTable.length"
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
import { CaretLeft, CloseBold, Delete, Download } from '@element-plus/icons-vue'
import variables from 'styles/variables.module.scss'
import InterfaceApi from 'api/interface.js'
import ServiceApi from 'api/service.js'
import { userStore } from 'stores/user'
import Papa from 'papaparse'
import dayjs from 'dayjs'
const users = userStore()
console.log('location -> ', location)

const contentRef = ref(null)
const ctxData = reactive({
  websocket: null, //websocket
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
  logType: 0, //0,1,2
  logName: ['采集日志', '上报日志', '系统日志'],
  collectName: '', //采集接口名称
  collectList: [], //采集接口列表
  uploadName: '', //上报接口名称
  uploadList: [], //上报接口列表
  startStopFlag: true,
  dataTable: [],
  dataType: ['接收方', '发送方'],
})
//日志类型更改
const onChange = (val) => {
  console.log('onChange -> val', val)
  ctxData.logType = val
  console.log('ctxData.startStopFlag = ', ctxData.startStopFlag)
  if (!ctxData.startStopFlag) {
    ctxData.startStopFlag = true
    closeWs()
  }
}
//开始结束状态更改
const changeSSFlag = (val) => {
  ctxData.startStopFlag = val
  //建立websocket连接
  if (val === false) {
    initWebsocket()
  } else {
    closeWs()
  }
}
//初始化websocket
const initWebsocket = () => {
  let ip = ''
  if (process.env.NODE_ENV === 'production') {
    ip = location.origin.replace('https://', '').replace('http://', '')
  } else {
    ip = 'localhost:8080'
  }
  let param =
    '?type=' +
    ctxData.logType +
    '&name=' +
    (ctxData.logType === 0 ? ctxData.collectName : ctxData.logType === 1 ? ctxData.uploadName : 'systemLog')

  let url = 'ws://' + ip + '/api/ws' + param
  console.log('url = ' + url)
  ctxData.websocket = new WebSocket(url)
  //获取websocket信息
  getMessage()
  //监听连接断开
  onClose()
  //监听连接异常
  onError()
}

// 发送消息
const sendMessage = (msg) => {
  console.log('sendMessage -> ', ctxData.websocket)
  console.log('sendMessage -> msg', msg)
  if (ctxData.websocket && ctxData.websocket.readyState === 1) {
    ctxData.websocket.send(msg)
  }
}
// 获取消息
const getMessage = () => {
  ctxData.websocket.onmessage = function (event) {
    console.log('event.data', event.data)
    if (event.data && event.data != '') {
      ctxData.dataTable.push(JSON.parse(event.data))
    }
  }
}

//关闭websocket连接
const closeWs = () => {
  //关闭连接
  if (ctxData.websocket) {
    console.log('关闭ws命令')
    // const connectionInfo = {
    //   type: 'close',
    //   name: 'close',
    // }
    // sendMessage(JSON.stringify(connectionInfo))
    ctxData.websocket.close()
    console.log('closeWs', ctxData.websocket)
  }
}
// 监听websocket关闭事件
const onClose = () => {
  ctxData.websocket.onclose = (e) => {
    console.log('onClose - > e', e)
    ctxData.startStopFlag = true
    ctxData.websocket = null
  }
}
// 监听websocket连接异常事件
const onError = () => {
  ctxData.websocket.onerror = (e) => {
    console.log('onError - > e', e)
    ctxData.startStopFlag = true
    ctxData.websocket = null
  }
}
// 获取采集接口信息
const getInterfaceList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  InterfaceApi.getInterfaceList(pData).then((res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.collectList = res.data
      if (res.data.length > 0) {
        ctxData.collectName = res.data[0].collInterfaceName
      }
    } else {
      showOneResMsg(res)
    }
  })
}
getInterfaceList()
// 获取上报接口信息
const getUploadList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  ServiceApi.getGatewayList(pData).then((res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.uploadList = res.data
      if (res.data.length > 0) {
        ctxData.uploadName = res.data[0].serviceName
      }
    } else {
      showOneResMsg(res)
    }
  })
}
getUploadList()

nextTick(() => {
  ctxData.tableMaxHeight = contentRef.value.clientHeight - 34 - 36 - 22
})
//处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
//处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
// 过滤表格数据
const filterTableData = computed(() =>
  ctxData.dataTable.slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
)
// 清空数据表
const clearDataTable = () => {
  ctxData.dataTable = []
}
//监听变化
watch(
  () => ctxData.collectName,
  (name, prevName) => {
    ctxData.startStopFlag = true
    closeWs()
  }
)
watch(
  () => ctxData.uploadName,
  (name, prevName) => {
    ctxData.startStopFlag = true
    closeWs()
  }
)

const exportCsv = () => {
  let outPutData = []
  ctxData.dataTable.forEach((item) => {
    outPutData.push({
      timeStamp: item.timeStamp + '\t',
      direction: ctxData.dataType[item.direction],
      label: item.label,
      content: item.content,
    })
  })
  var csv = Papa.unparse(outPutData)
  //定义文件内容，类型必须为Blob 否则createObjectURL会报错
  let content = new Blob([csv])
  //生成url对象
  let urlObject = window.URL || window.webkitURL || window
  let url = urlObject.createObjectURL(content)
  //生成<a></a>DOM元素
  let el = document.createElement('a')
  //链接赋值
  el.href = url
  let fileName = ''
  if (ctxData.logType === 0) {
    fileName = '采集日志'
  } else if (ctxData.logType === 1) {
    fileName = '上报日志'
  } else {
    fileName = '系统日志'
  }
  var nowTime = dayjs(new Date()).format('YYYY-MM-DD HH-mm-ss')
  fileName += nowTime
  el.download = fileName + '.csv'
  //必须点击否则不会下载
  el.click()
  //移除链接释放资源
  urlObject.revokeObjectURL(url)
}
//显示单个res结果，code不等于 '0' 的message
const showOneResMsg = (res) => {
  ElMessage({
    type: 'error',
    message: res.message,
  })
}
onUnmounted(() => {
  console.log('页面关闭了')
  closeWs()
})
</script>
<style lang="scss" scoped>
@use 'styles/custom-scoped.scss' as *;
.SysLog {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  .title {
    position: absolute;
    top: 20px;
    height: 76px;
    left: 20px;
    right: 20px;
    display: flex;
    align-items: center;
    padding: 24px;
    box-sizing: border-box;
    background-color: #f5f8fa;
    border-radius: 4px;
  }
  .content {
    position: absolute;
    top: 124px;
    left: 20px;
    right: 20px;
    bottom: 20px;
    background-color: #f5f8fa;
    border-radius: 4px;
    padding: 24px;
    box-sizing: border-box;
    .cTitle {
      position: relative;
      height: 36px;
      box-sizing: border-box;
      display: flex;
      justify-content: space-between;
      align-items: center;
      .tName {
        line-height: 16px;
        font-size: 18px;
        border-left: 4px solid #3054eb;
        padding-left: 20px;
      }
      .option {
        display: flex;
        align-items: center;
      }
    }

    .cTable {
      position: absolute;
      top: 80px;
      bottom: 24px;
      left: 24px;
      right: 24px;
    }
  }
  .right-btn {
    margin-left: 20px;
  }
}
.pagination {
  position: absolute;
  bottom: 12px;
  left: 0;
  right: 0;
  display: flex;
  justify-content: center;
}
:deep(.el-check-tag) {
  line-height: 22px;
  margin-right: 20px;
  background-color: #e6ecef;
}
:deep(.el-check-tag.is-checked) {
  background-color: #c6e2ff;
}
</style>
