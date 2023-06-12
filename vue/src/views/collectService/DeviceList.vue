<template>
  <div class="main-container">
    <div v-if="ctxData.idFlag" class="main" style="overflow:hidden;">
      <div class="search-bar">
        <el-form :model="ctxData.screenForm"  :inline="true" ref="searchFormRef" status-icon label-width="90px">
          <el-form-item label="设备名称" prop="name">
            <el-input type="text" v-model="ctxData.screenForm.name" autocomplete="off" placeholder="请输入设备名称">
            </el-input>
          </el-form-item>
          <el-form-item label="设备标签" prop="label">
            <el-input type="text" v-model="ctxData.screenForm.label" autocomplete="off" placeholder="请输入设备标签">
            </el-input>
          </el-form-item>
          <el-form-item label="设备地址" prop="addr">
            <el-input type="text" v-model="ctxData.screenForm.addr" autocomplete="off" placeholder="请输入设备地址">
            </el-input>
          </el-form-item>
          <el-form-item label="设备模型" prop="tsl">
            <el-select v-model="ctxData.screenForm.tsl" clearable style="width: 100%" placeholder="请选择设备模型">
              <el-option
                v-for="(item, index) of ctxData.deviceModelList"
                :key="'dm_' + index"
                :label="item.name"
                :value="item.name"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="采集接口" prop="collInterfaceName">
            <el-select
              v-model="ctxData.screenForm.collInterfaceName"
              clearable
              style="width: 100%"
              placeholder="请选择采集接口"
              multiple
              collapse-tags
              collapse-tags-tooltip
            >
              <el-option
                v-for="(item, index) of ctxData.interfaceList"
                :key="'cin_' + index"
                :label="item.collInterfaceName"
                :value="item.collInterfaceName"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="通信状态">
            <el-radio-group v-model="ctxData.screenForm.commStatus">
              <el-radio label="Line">全部</el-radio>
              <el-radio label="onLine">在线</el-radio>
              <el-radio label="offLine">离线</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item style="margin-left: 20px;">
            <el-button type="primary" @click="doScreen()">搜索</el-button>
          </el-form-item>
          <el-form-item>
            <el-button @click="cancelScreen()">重置</el-button>
          </el-form-item>
          <el-form-item>
            <el-button style="color: #fff" color="#2EA554" @click="refresh()">
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
          :data="filterTableData"
          :cell-style="ctxData.cellStyle"
          :header-cell-style="ctxData.headerCellStyle"
          style="width: 100%"
          :max-height="ctxData.tableMaxHeight"
          stripe
        >
          <el-table-column prop="name" label="设备名称" width="auto" min-width="160" align="center" fixed="left">
          </el-table-column>
          <el-table-column prop="label" label="设备标签" width="auto" min-width="160" align="center"> </el-table-column>
          <el-table-column prop="tsl" label="设备模型" width="auto" min-width="160" align="center">
            <template #default="scope">
              <el-button type="primary" link @click="showBlockParams(scope.row)">{{scope.row.tsl}}</el-button>
            </template>
          </el-table-column>
          <el-table-column prop="collInterfaceName" label="采集接口" width="auto" min-width="160" align="center">
          </el-table-column>
          <el-table-column prop="addr" label="通讯地址" width="auto" min-width="100" align="center"> </el-table-column>
          <el-table-column label="当前通信状态" width="auto" min-width="150" align="center">
            <template #default="scope">
              <el-tag v-if="scope.row.commStatus === 'onLine'" type="success">在线</el-tag>
              <el-tag v-else type="danger">离线</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="lastCommRTC" label="最后通信时间" width="auto" min-width="200" align="center">
          </el-table-column>
          <el-table-column prop="commTotalCnt" label="通信总次数" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column prop="commSuccessCnt" label="通信成功次数" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column label="操作" width="auto" min-width="200" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="showDeviceProperty(scope.row)" text type="success">查看变量</el-button>
            </template>
          </el-table-column>
          <template #empty>
            <div>无数据</div>
          </template>
        </el-table>
        <div class="pagination" style="bottom: 24px;">
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
    <div v-else class="main">
      <DeviceProperty :curDevice="ctxData.curDevice" :pageInfo="'dashboard'" @changeIdFlag="changeIdFlag" style="width: 100%; height: 96%;overflow:hidden;"></DeviceProperty>
    </div>
  </div>
</template>
<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import variables from 'styles/variables.module.scss'
import DeviceModelApi from 'api/deviceModel.js'
import InterfaceApi from 'api/interface.js'
import { userStore } from 'stores/user'
import { useRouter } from 'vue-router'
import DeviceProperty from './deviceManage/Device-property.vue'
const users = userStore()

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
  deviceTableData: [],
  screenData: [],
  deviceTotal: 0,
  deviceOnline: 0,
  dFlag: false,
  dTitle: '添加设备',
  deviceInfo: '',
  deviceModelList: [],
  selectedDevices: [],
  idFlag: true, //设备-属性切换显示表示
  uFlag: false,
  commStatusOptions: [
    {
      label: '全部',
      value: 'Line',
    },
    {
      label: '在线',
      value: 'onLine',
    },
    {
      label: '离线',
      value: 'offLine',
    },
  ],
  sFlag: false,
  screenForm: {
    name: '',
    label: '',
    tsl: '',
    commStatus: 'Line',
    collInterfaceName: [],
    addr: '',
  },
  interfaceList: [],
  interfaces: [],
  checkedInterfaceList: [],
  collDeviceObj: {}, //采集接口与设备对应对象

})

const changeIdFlag = () => {
  ctxData.idFlag = true
  getCollDevices()
}

// 获取采集接口列表
const getInterfaceList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  InterfaceApi.getInterfaceList(pData).then((res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.interfaceList = res.data
      ctxData.interfaces = []
      ctxData.interfaceList.forEach((item) => {
        ctxData.interfaces.push(item.collInterfaceName)
        ctxData.checkedInterfaceList.push(item.collInterfaceName)
      })
      getCollDevices()
    } else {
      showOneResMsg(res)
    }
  })
}
getInterfaceList()

const contentRef = ref(null)
// 获取采集接口下的设备列表
const getCollDevices = (flag) => {
  const pData = {
    token: users.token,
    data: {
      names: ctxData.checkedInterfaceList,
    },
  }
  InterfaceApi.getAllCollDevices(pData).then(async (res) => {
    console.log('getCollDevices -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.deviceTableData = res.data.node == null ? [] : res.data.node
      // ctxData.deviceTableData.forEach((item) => {
      //   item['collInterfaceName'] = 'modbusRCP'
      // })
      ctxData.screenData = ctxData.deviceTableData
      ctxData.deviceTotal = res.data.commTotalCnt
      ctxData.deviceOnline = res.data.commSuccessCnt
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
      ctxData.tableMaxHeight = contentRef.value.clientHeight - 34 - 36 - 22
    })
  })
}
// 获取物模型列表
const getDeviceModelList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  DeviceModelApi.getDeviceModelList(pData).then((res) => {
    console.log('getDeviceModelList -> res', res)
    if (res.code === '0') {
      ctxData.deviceModelList = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
getDeviceModelList()
// 处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
// 过滤表格数据
const filterTableData = computed(() => {
  return ctxData.screenData.slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  return ctxData.screenData
})
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
const refresh = () => {
  getCollDevices(1)
}
const showDeviceProperty = (row) => {
  ctxData.idFlag = false
  ctxData.curDevice = row
}
//处理过滤
const doScreen = () => {
  console.log('doScreen ctxData.deviceTableData = ', ctxData.deviceTableData)
  console.log('doScreen ctxData.screenForm = ', ctxData.screenForm)
  ctxData.screenData = []
  ctxData.deviceTotal = 0
  ctxData.deviceOnline = 0
  ctxData.deviceTableData.forEach((device) => {
    var flag1 = device.name.includes(ctxData.screenForm.name)
    var flag2 = device.label.includes(ctxData.screenForm.label)
    var flag3 = device.addr.includes(ctxData.screenForm.addr)
    var flag4 = ctxData.screenForm.tsl === '' || device.tsl == ctxData.screenForm.tsl
    var flag5 =
      ctxData.screenForm.collInterfaceName.length === 0 ||
      ctxData.screenForm.collInterfaceName.includes(device.collInterfaceName)
    var flag6 = ctxData.screenForm.commStatus === 'Line' || device.commStatus.includes(ctxData.screenForm.commStatus)
    if (flag1 && flag2 && flag3 && flag4 && flag5 && flag6) {
      ctxData.screenData.push(device)
      if (device.commStatus == 'onLine') {
        ctxData.deviceOnline++
      }
      ctxData.deviceTotal++
    }
  })
}
//取消过滤
const cancelScreen = () => {
  initScreenForm()
}
//初始化过滤表单
const initScreenForm = () => {
  ctxData.screenForm = {
    name: '',
    tsl: '',
    label: '',
    commStatus: 'Line',
    collInterfaceName: [],
    addr: '',
  }
}

//显示单个res结果，code不等于 '0' 的message
const showOneResMsg = (res) => {
  ElMessage({
    type: 'error',
    message: res.message,
  })
}
// 查看命令详情
const router = useRouter()
const showBlockParams = (row) => {
  router.push({ path: 'deviceModelRtu', query: { tsl: row.tsl }})
}
</script>
<style lang="scss" scoped>
@use 'styles/custom-scoped.scss' as *;
.title-count {
  display: flex;
  justify-content: space-between;
}
</style>
