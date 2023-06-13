<template>
<div class="main-container">
  <div v-if="ctxData.dpFlag" class="main">
    <div class="search-bar">
      <el-form :inline="true" ref="searchFormRef" status-icon label-width="90px">
        <el-form-item style="margin-left: 20px;">
          <el-button type="primary" plain @click="toInterface()" style="margin-right: 20px">
            <el-icon class="el-input__icon"><back /></el-icon>
            返回采集接口
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" plain class="right-btn" @click="importDevice()">
            <el-icon class="el-input__icon"><download /></el-icon>
            批量导入设备
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" plain class="right-btn" @click="exportDevice()">
            <el-icon class="el-input__icon"><upload /></el-icon>
            批量导出设备
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="search-bar" style="display: flex;">
      <div class="title" style="position: relative;margin-right: 40px;justify-content: flex-start;padding: 0px 0px;height: 40px;">
        <div class="tName">{{ props.curInterface.collInterfaceName }}</div>
      </div>
      <el-form :inline="true" ref="searchFormRef2" status-icon label-width="90px">
        <el-form-item label="">
          <el-input
            style="width: 230px; margin-right: 20px"
            placeholder="请输入 名称/标签/地址/ 过滤"
            clearable
            v-model="ctxData.deviceInfo"
          >
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="通信状态">
          <el-select v-model="ctxData.commStatus" style="width: 100px" placeholder="请选择通信状态">
            <el-option
              v-for="(item, index) of ctxData.commStatusOptions"
              :key="'cs_' + index"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" bg class="right-btn" @click="addDevice()">
            <el-icon class="btn-icon">
              <Icon name="local-add" size="14px" color="#ffffff" />
            </el-icon>
            添加
          </el-button>
        </el-form-item>
        <el-form-item>
          <!-- <el-button type="primary" bg class="right-btn" @click="editDevices()">
            <el-icon class="btn-icon tianjia"></el-icon>
            批量修改
          </el-button> -->

          <el-button type="primary" bg class="right-btn" @click="allCollect()">
            <el-icon class="btn-icon">
              <Icon name="local-refresh" size="14px" color="#ffffff" />
            </el-icon>
            批量采集
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="danger" bg class="right-btn" @click="deleteDevice()">
            <el-icon class="btn-icon">
              <Icon name="local-delete" size="14px" color="#ffffff" />
            </el-icon>
            删除
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button style="color: #fff" color="#2EA554" class="right-btn" @click="refresh()">
            <el-icon class="btn-icon">
              <Icon name="local-refresh" size="14px" color="#ffffff" />
            </el-icon>
            刷新
          </el-button>
        </el-form-item>
        <el-form-item>
          <span style="margin-right: 10px">设备总数：{{ ctxData.deviceTotal }}</span>
          <span>在线：{{ ctxData.deviceOnline }}</span>
        </el-form-item>
      </el-form>
    </div>
    <div class="content" ref="contentRef" style="top: 136px">
      <el-table
        :data="filterTableData"
        :cell-style="ctxData.cellStyle"
        :header-cell-style="ctxData.headerCellStyle"
        :max-height="ctxData.tableMaxHeight"
        style="width: 100%"
        stripe
        @selection-change="handleSelectionChange"
        @row-dblclick="editDevice"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column sortable prop="name" label="设备名称" width="auto" min-width="180" align="center"> </el-table-column>
        <el-table-column sortable prop="label" label="设备标签" width="auto" min-width="180" align="center"> </el-table-column>
        <el-table-column sortable prop="tsl" label="设备模型" width="auto" min-width="180" align="center"> </el-table-column>
        <el-table-column sortable prop="addr" label="通讯地址" width="auto" min-width="100" align="center"> </el-table-column>
        <el-table-column sortable label="当前通信状态" width="auto" min-width="150" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.commStatus === 'onLine'" type="success">在线</el-tag>
            <el-tag v-else type="danger">离线</el-tag>
          </template>
        </el-table-column>
        <el-table-column sortable prop="lastCommRTC" label="最后通信时间" width="auto" min-width="200" align="center">
        </el-table-column>
        <el-table-column sortable prop="commTotalCnt" label="通信总次数" width="auto" min-width="150" align="center">
        </el-table-column>
        <el-table-column sortable prop="commSuccessCnt" label="通信成功次数" width="auto" min-width="150" align="center">
        </el-table-column>
        <el-table-column label="操作" width="auto" min-width="200" align="center" fixed="right">
          <template #default="scope">
            <el-button @click="showDeviceProperty(scope.row)" text type="success">查看变量</el-button>
            <el-button @click="editDevice(scope.row)" text type="primary">编辑</el-button>
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
  </div>
  <DeviceProperty
    v-else
    :collInterfaceName="props.curInterface.collInterfaceName"
    :curDevice="ctxData.curDevice"
    @changeDpFlag="changeDpFlag"
  ></DeviceProperty>
  <!-- 添加编辑设备 -->
  <el-dialog
    v-model="ctxData.dFlag"
    :title="ctxData.dTitle"
    width="600px"
    :before-close="handleClose"
    :close-on-click-modal="false"
  >
    <div class="dialog-content">
      <el-form
        :model="ctxData.deviceForm"
        :rules="ctxData.deviceRules"
        ref="deviceFormRef"
        status-icon
        label-position="right"
        label-width="100px"
      >
        <el-form-item label="设备名称" prop="name">
          <el-input
            :disabled="ctxData.dTitle.includes('编辑')"
            type="text"
            v-model="ctxData.deviceForm.name"
            autocomplete="off"
            placeholder="请输入设备名称"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="设备标签" prop="label">
          <el-input type="text" v-model="ctxData.deviceForm.label" autocomplete="off" placeholder="请输入设备标签">
          </el-input>
        </el-form-item>
        <el-form-item label="设备模型" prop="tsl">
          <el-select v-model="ctxData.deviceForm.tsl" style="width: 100%" placeholder="请选择设备模型">
            <el-option
              v-for="(item, index) of ctxData.deviceModelList"
              :key="'dm_' + index"
              :label="item.name"
              :value="item.name"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="采集接口" prop="interfaceName">
          <el-input
            disabled
            type="text"
            v-model="ctxData.deviceForm.interfaceName"
            autocomplete="off"
            placeholder="请输入采集周期"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="通讯地址" prop="addr">
          <el-input type="text" v-model="ctxData.deviceForm.addr" autocomplete="off" placeholder="请输入通讯地址">
          </el-input>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="cancelSubmit()">取消</el-button>
        <el-button type="primary" @click="submitDeviceForm()">保存</el-button>
      </span>
    </template>
  </el-dialog>
  <!-- 批量导入设备 -->
  <el-dialog
    v-model="ctxData.uFlag"
    title="上传设备xlsx文件"
    width="600px"
    :before-close="beforeCloseUploadDevice"
    :close-on-click-modal="false"
  >
    <el-upload
      ref="uploadRef"
      action=""
      :auto-upload="false"
      :http-request="myRequest"
      :limit="1"
      :on-exceed="handleExceed"
      :before-upload="beforeUpload"
    >
      <el-button type="primary">选择文件</el-button>
      <template #tip>
        <div class="el-upload__tip">只能上传一个文件，只支持xlsx格式文件！</div>
      </template>
    </el-upload>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="cancelUploadDevice">取消</el-button>
        <el-button type="primary" @click="submitUploadDevice">上传</el-button>
      </span>
    </template>
  </el-dialog>
</div>
</template>
<script setup>
import { Search, Back, Download, Upload } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import variables from 'styles/variables.module.scss'
import DeviceModelApi from 'api/deviceModel.js'
import InterfaceApi from 'api/interface.js'
import { userStore } from 'stores/user'
import DeviceProperty from './Device-property.vue'
const users = userStore()
const props = defineProps({
  curInterface: {
    type: Object,
    default: {},
  },
})
console.log('id -> props', props)

const emit = defineEmits(['changeIdFlag'])
const toInterface = () => {
  emit('changeIdFlag')
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
  deviceTableData: [],
  deviceTotal: 0,
  deviceOnline: 0,
  dFlag: false,
  dTitle: '添加设备',
  deviceInfo: '',
  deviceForm: {
    interfaceName: props.curInterface.collInterfaceName, // 采集接口名称，只能是字母+数字的组合，不可以是中文
    name: '', // 设备名称
    label: '', // 设备标签
    addr: '', // 设备通信地址
    tsl: '', // 物模型
  },
  deviceRules: {
    interfaceName: [
      {
        required: true,
        validator: '采集接口不能为空',
        trigger: 'blur',
      },
    ],
    name: [
      {
        required: true,
        message: '设备名称不能为空',
        trigger: 'blur',
      },
    ],
    label: [
      {
        required: true,
        message: '设备标签不能为空',
        trigger: 'blur',
      },
    ],
    addr: [
      {
        required: true,
        message: '通信地址不能为空',
        trigger: 'blur',
      },
    ],
    tsl: [
      {
        required: true,
        message: '物模型不能为空',
        trigger: 'blur',
      },
    ],
  },
  deviceModelList: [],
  selectedDevices: [],
  dpFlag: true, //设备-属性切换显示表示
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
  commStatus: 'Line',
})

const changeDpFlag = () => {
  ctxData.dpFlag = true
  getCollDevices()
}

const contentRef = ref(null)
// 获取采集接口下的设备列表
const getCollDevices = (flag) => {
  const pData = {
    token: users.token,
    data: {
      name: props.curInterface.collInterfaceName,
    },
  }
  InterfaceApi.getCollDevices(pData).then(async (res) => {
    console.log('getCollDevices -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.deviceTableData = res.data.deviceNodeMap == null ? [] : res.data.deviceNodeMap
      ctxData.deviceTotal = res.data.deviceNodeCnt
      ctxData.deviceOnline = res.data.deviceNodeOnlineCnt
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
getCollDevices()
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
// 处理复选框事件
const handleSelectionChange = (val) => {
  ctxData.selectedDevices = val
}
// 处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
// 过滤表格数据
const filterTableData = computed(() => {
  let deviceInfo = ctxData.deviceInfo
  return ctxData.deviceTableData
    .filter(
      (item) =>
        !deviceInfo ||
        item.name.toLowerCase().includes(deviceInfo.toLowerCase()) ||
        item.label.toLowerCase().includes(deviceInfo.toLowerCase()) ||
        item.addr.toLowerCase().includes(deviceInfo.toLowerCase())
    )
    .filter((item) => item.commStatus.includes(ctxData.commStatus))
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  let deviceInfo = ctxData.deviceInfo
  return ctxData.deviceTableData
    .filter(
      (item) =>
        !deviceInfo ||
        item.name.toLowerCase().includes(deviceInfo.toLowerCase()) ||
        item.label.toLowerCase().includes(deviceInfo.toLowerCase()) ||
        item.addr.toLowerCase().includes(deviceInfo.toLowerCase())
    )
    .filter((item) => item.commStatus.includes(ctxData.commStatus))
})
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
const refresh = () => {
  getCollDevices(1)
}
// 添加设备
const addDevice = () => {
  ctxData.dFlag = true
  ctxData.dTitle = '添加设备'
}
// 编辑设备
const editDevice = (row) => {
  ctxData.dFlag = true
  ctxData.dTitle = '编辑设备'
  ctxData.deviceForm = {
    interfaceName: props.curInterface.collInterfaceName, // 采集接口名称，只能是字母+数字的组合，不可以是中文
    name: row.name, // 设备名称
    label: row.label, // 设备标签
    addr: row.addr, // 设备通信地址
    tsl: row.tsl, // 物模型
  }
}
const deviceFormRef = ref(null)
const submitDeviceForm = () => {
  deviceFormRef.value.validate((valid) => {
    if (valid) {
      const pData = {
        token: users.token,
        data: ctxData.deviceForm,
      }
      if (ctxData.dTitle.includes('添加')) {
        InterfaceApi.addCollDevice(pData).then((res) => {
          handleResult(res, getCollDevices)
          cancelSubmit()
        })
      } else {
        InterfaceApi.editCollDevice(pData).then((res) => {
          handleResult(res, getCollDevices)
          cancelSubmit()
        })
      }
    } else {
      return false
    }
  })
}
// 取消提交
const cancelSubmit = () => {
  ctxData.dFlag = false
  deviceFormRef.value.resetFields()
  initDeviceForm()
}
//处理弹出框右上角关闭图标
const handleClose = (done) => {
  cancelSubmit()
}
const initDeviceForm = () => {
  ctxData.deviceForm = {
    interfaceName: props.curInterface.collInterfaceName, // 采集接口名称，只能是字母+数字的组合，不可以是中文
    name: '', // 设备名称
    label: '', // 设备标签
    addr: '', // 设备通信地址
    tsl: '', // 物模型
  }
}

const showDeviceProperty = (row) => {
  ctxData.dpFlag = false
  ctxData.curDevice = row
}
// 批量修改设备
const editDevices = () => {
  ElMessage.info('功能完善中...')
}
// 批量采集设备
const allCollect = () => {
  if (ctxData.selectedDevices.length === 0) {
    ElMessage.info('请至少选择一个设备！')
    return
  } else {
    let count = 0
    ctxData.selectedDevices.forEach((item) => {
      const pData = {
        token: users.token,
        data: {
          collInterfaceName: props.curInterface.collInterfaceName,
          deviceName: item.name,
        },
      }
      InterfaceApi.getDeviceDataReal(pData).then((res) => {
        count++
        if (res.code === '0') {
          ElMessage.success(item.label + '：采集成功！')
        } else {
          showOneResMsg(res)
        }
        if (count === ctxData.selectedDevices.length) {
          getCollDevices(1)
        }
      })
    })
  }
}
// 批量删除设备
const deleteDevice = () => {
  let dList = []
  if (ctxData.selectedDevices.length === 0) {
    ElMessage.info('请至少选择一个设备！')
    return
  } else {
    ctxData.selectedDevices.forEach((item) => {
      dList.push(item.name)
    })
  }
  console.log('dList', dList)
  ElMessageBox.confirm('确定要删除这些吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          interfaceName: props.curInterface.collInterfaceName,
          deviceNames: dList,
        },
      }
      InterfaceApi.deleteCollDevices(pData).then((res) => {
        handleResult(res, getCollDevices)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}
// 批量导出设备
const exportDevice = () => {
  const pData = {
    token: users.token,
    responseType: 'blob',
    data: {
      name: props.curInterface.collInterfaceName,
    },
  }
  InterfaceApi.exportDevice(pData).then((res) => {
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
// 批量导入设备
const importDevice = () => {
  ctxData.uFlag = true
}
/**
 * 提交上传的插件
 */
const submitUploadDevice = () => {
  uploadRef.value && uploadRef.value.submit()
}

// 取消插件自带的xhr
const myRequest = (obj) => {}
const uploadRef = ref(null)
/**
 * 上传文件大于limit时事件
 * @param {要上传的文件} filesimportPlugin
 */
const handleExceed = (files) => {
  uploadRef.value.clearFiles()
  //超过limit取第一个
  const file = files[0]
  uploadRef.value.handleStart(file)
}
/**
 * 文件上传前事件
 * @param {要上传的文件} file
 */
const beforeUpload = (file) => {
  console.log('beforeUpload -> file', file)
  const fileTypeList = ['application/vnd.openxmlformats-officedocument.spreadsheetml.sheet']
  let typeFlag = fileTypeList.includes(file.type) || (file.name != '' && file.name.indexOf('xlsx') > -1)
  if (!typeFlag) {
    ElMessage({
      type: 'error',
      message: '文件格式不正确,必须是xlsx文件！',
    })
    return
  }
  // 调上传接口
  if (typeFlag) {
    let formData = new FormData()
    formData.append('name', props.curInterface.collInterfaceName)
    formData.append('fileName', file)
    const pData = {
      token: users.token,
      contentType: 'multipart/form-data',
      data: formData,
    }
    InterfaceApi.addDeviceFromCSV(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success(res.message)
        ctxData.uFlag = false
        getCollDevices()
      } else {
        showOneResMsg(res)
      }
    })
  }
}
//取消上传设备模型文件
const cancelUploadDevice = () => {
  ctxData.uFlag = false
  uploadRef.value.clearFiles()
}
// 弹窗右上角关闭事件处理
const beforeCloseUploadDevice = () => {
  cancelUploadDevice()
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
</style>
