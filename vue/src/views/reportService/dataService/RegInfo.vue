<template>
  <div class="main-container">
    <div class="main">
      <div class="title" style="justify-content: space-between">
        <div class="title-left">
          <el-button type="primary" plain @click="toDataService()" style="margin-right: 20px">
            <el-icon class="el-input__icon"><back /></el-icon>
            返回上报服务
          </el-button>
          <el-input style="width: 200px" placeholder="请输入 地址/名称 过滤" clearable v-model="ctxData.regInfo">
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
        </div>
        <div>
          <el-button type="primary" plain class="right-btn" @click="importReg()">
            <el-icon class="el-input__icon"><download /></el-icon>
            导入寄存器
          </el-button>
          <el-button type="primary" plain class="right-btn" @click="exportReg">
            <el-icon class="el-input__icon"><upload /></el-icon>
            导出寄存器
          </el-button>
          <el-button type="primary" bg class="right-btn" @click="reportNode()">
            <el-icon class="btn-icon">
              <Icon name="local-report" size="14px" color="#ffffff" />
            </el-icon>
            主动上报
          </el-button>
          <el-button type="primary" bg class="right-btn" @click="addReg()">
            <el-icon class="btn-icon">
              <Icon name="local-add" size="14px" color="#ffffff" />
            </el-icon>
            添加
          </el-button>
          <el-button type="danger" bg class="right-btn" @click="deleteReg()">
            <el-icon class="btn-icon">
              <Icon name="local-delete" size="14px" color="#ffffff" />
            </el-icon>
            删除
          </el-button>
          <el-button style="color: #fff" color="#2EA554" class="right-btn" @click="refresh()">
            <el-icon class="btn-icon">
              <Icon name="local-refresh" size="14px" color="#ffffff" />
            </el-icon>
            刷新
          </el-button>
        </div>
      </div>
      <div class="title" style="top: 60px; height: 76px; padding: 20px 0; justify-content: flex-start">
        <div class="tName">{{ props.curGateway.serviceName }}</div>
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
          @row-dblclick="editReg"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="regName" label="寄存器名称" width="auto" min-width="180" align="center">
          </el-table-column>
          <el-table-column prop="label" label="寄存器标签" width="auto" min-width="180" align="center">
          </el-table-column>
          <el-table-column prop="propertyType" label="属性类型" width="auto" min-width="120" align="center">
          </el-table-column>
          <el-table-column prop="collName" label="采集接口名称" width="auto" min-width="180" align="center">
          </el-table-column>
          <el-table-column prop="nodeName" label="设备名称" width="auto" min-width="180" align="center">
          </el-table-column>
          <el-table-column prop="propertyName" label="设备属性名称" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column prop="regType" label="寄存器类型" width="auto" min-width="100" align="center">
          </el-table-column>
          <el-table-column prop="regAddr" label="寄存器地址" width="auto" min-width="100" align="center">
          </el-table-column>
          <el-table-column prop="regCnt" label="寄存器数量" width="auto" min-width="100" align="center">
          </el-table-column>
          <el-table-column prop="rule" label="解析规则" width="auto" min-width="100" align="center"> </el-table-column>
          <el-table-column label="操作" width="auto" min-width="120" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="editReg(scope.row)" text type="primary">编辑</el-button>
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
    <!-- 添加寄存器 -->
    <el-dialog v-model="ctxData.rFlag" :title="ctxData.rTitle" width="600px">
      <div class="dialog-content">
        <el-form
          :model="ctxData.regForm"
          :rules="ctxData.regRules"
          ref="regFormRef"
          status-icon
          label-position="right"
          label-width="120px"
        >
          <el-form-item label="采集接口名称" prop="collName">
            <el-select
              v-model="ctxData.regForm.collName"
              style="width: 100%"
              placeholder="请选择采集接口名称"
              @change="selectCollName"
              clearable
              filterable
            >
              <el-option v-for="item in ctxData.collOptions" :key="item.name" :label="item.label" :value="item.name">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="设备名称" prop="nodeName">
            <el-select
              v-model="ctxData.regForm.nodeName"
              style="width: 100%"
              placeholder="请选择设备名称"
              @change="selectNodeName"
              clearable
              filterable
            >
              <el-option v-for="item in ctxData.nodeOptions" :key="item.name" :label="item.label" :value="item.name">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="设备属性名称" prop="propertyName">
            <el-select
              v-model="ctxData.regForm.propertyName"
              style="width: 100%"
              placeholder="请选择设备属性名称"
              clearable
              filterable
            >
              <el-option
                v-for="item in ctxData.propertyOptions"
                :key="item.name"
                :label="item.label"
                :value="item.name"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="寄存器名称" prop="regName">
            <el-input type="text" v-model="ctxData.regForm.regName" autocomplete="off" placeholder="请输入寄存器名称">
            </el-input>
          </el-form-item>
          <el-form-item label="寄存器标签" prop="label">
            <el-input type="text" v-model="ctxData.regForm.label" autocomplete="off" placeholder="请输入寄存器标签">
            </el-input>
          </el-form-item>

          <el-form-item label="属性类型" prop="propertyType">
            <el-select
              v-model="ctxData.regForm.propertyType"
              style="width: 100%"
              placeholder="请选择属性类型"
              clearable
              filterable
            >
              <el-option v-for="item in ctxData.typeOptions" :key="item.name" :label="item.label" :value="item.name">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="寄存器类型" prop="regType">
            <el-select
              v-model="ctxData.regForm.regType"
              style="width: 100%"
              placeholder="请选择寄存器类型"
              clearable
              filterable
            >
              <el-option v-for="item in ctxData.regTypeOptions" :key="item.name" :label="item.label" :value="item.name">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="寄存器地址" prop="regAddr">
            <el-input
              type="text"
              v-model.number="ctxData.regForm.regAddr"
              autocomplete="off"
              placeholder="请输入寄存器地址"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="寄存器数量" prop="regCnt">
            <el-select
              v-model.number="ctxData.regForm.regCnt"
              style="width: 100%"
              @change="changeRegCnt()"
              placeholder="请选择寄存器数量"
            >
              <el-option
                v-for="item in ctxData.regCntOptions"
                :key="'regCount_' + item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="解析规则" prop="rule">
            <el-select v-model.number="ctxData.regForm.rule" style="width: 100%" placeholder="请选择解析规则">
              <el-option
                v-for="item in ctxData.ruleTypeOptions['rt' + ctxData.regForm.regCnt]"
                :key="'type_' + item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelSubmit()">取消</el-button>
          <el-button type="primary" @click="submitReg()">保存</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- 导入寄存器 -->
    <el-dialog
      v-model="ctxData.uFlag"
      title="上传寄存器xlsx文件"
      width="600px"
      :before-close="beforeCloseUploadReg"
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
          <el-button @click="cancelUploadReg">取消</el-button>
          <el-button type="primary" @click="submitUploadReg">上传</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { Search, Back, Download, Upload } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import variables from 'styles/variables.module.scss'
import InterfaceApi from 'api/interface.js'
import DeviceModelApi from 'api/deviceModel.js'
import ServiceApi from 'api/service.js'
import TransferModelApi from 'api/transferModel.js'
import { userStore } from 'stores/user'
const users = userStore()
const props = defineProps({
  curGateway: {
    type: Object,
    default: {},
  },
})
console.log('id -> props', props)

const emit = defineEmits(['changeDnFlag'])
const toDataService = () => {
  emit('changeDnFlag')
}

const regCnt = /^[0-9]*[1-9][0-9]*$/
const validateRegAddr = (rule, value, callback) => {
  if (value !== 0 && !regCnt.test(value)) {
    callback(new Error('只能输入自然数！'))
  } else {
    callback()
  }
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
  regTableData: [],
  rFlag: false,
  rTitle: '添加寄存器',
  regInfo: '',
  regForm: {
    serviceName: '', // 服务名称
    regName: '', // 寄存器名称
    label: '', // 寄存器标签
    propertyType: 0, // 属性类型
    collName: '', // 采集接口名称
    nodeName: '', // 设备名称
    propertyName: '', //设备属性名称
    regType: 0, // 寄存器类型
    regAddr: 0, // 寄存器地址
    regCnt: 1, //寄存器数量
    rule: 'Int_AB', //解析规则
  },
  regRules: {
    regName: [
      {
        required: true,
        message: '寄存器名称不能为空！',
        trigger: 'blur',
      },
    ],
    label: [
      {
        required: true,
        message: '寄存器标签不能为空！',
        trigger: 'blur',
      },
    ],
    collName: [
      {
        required: true,
        message: '采集接口名称不能为空！',
        trigger: 'blur',
      },
    ],
    nodeName: [
      {
        required: true,
        message: '设备名称不能为空！',
        trigger: 'blur',
      },
    ],
    propertyName: [
      {
        required: true,
        message: '设备属性名称不能为空！',
        trigger: 'blur',
      },
    ],
    regAddr: [
      {
        required: true,
        message: '寄存器地址不能为空！',
        trigger: 'blur',
      },
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    regCnt: [
      {
        required: true,
        message: '寄存器数量不能为空！',
        trigger: 'blur',
      },
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    rule: [
      {
        required: true,
        message: '解析规则不能为空！',
        trigger: 'blur',
      },
    ],
  },
  collOptions: [], //采集接口选项
  nodeOptions: [], //设备选项
  propertyOptions: [], //设备属性选项
  nameTotsl: {}, //设备名称对应设备模板名称
  typeOptions: [
    { name: 0, label: '设备属性' },
    { name: 1, label: '系统属性' },
  ],
  regTypeOptions: [
    { name: 0, label: '线圈状态寄存器' },
    { name: 1, label: '离散输入状态寄存器' },
    { name: 2, label: '保持寄存器' },
    { name: 3, label: '输入寄存器' },
  ],
  regCntOptions: [
    { label: 1, value: 1 },
    { label: 2, value: 2 },
    { label: 4, value: 4 },
  ],
  //数据类型
  ruleTypeOptions: {
    rt1: [
      { label: 'Int_AB', value: 'Int_AB' },
      { label: 'Int_BA', value: 'Int_BA' },
    ],
    rt2: [
      { label: 'Long_ABCD', value: 'Long_ABCD' },
      { label: 'Long_BADC', value: 'Long_BADC' },
      { label: 'Long_DCBA', value: 'Long_DCBA' },
      { label: 'Long_CDAB', value: 'Long_CDAB' },
      { label: 'Float_ABCD', value: 'Float_ABCD' },
      { label: 'Float_BADC', value: 'Float_BADC' },
      { label: 'Float_DCBA', value: 'Float_DCBA' },
      { label: 'Float_CDAB', value: 'Float_CDAB' },
    ],
    rt4: [
      { label: 'Double_ABCDEFGH', value: 'Double_ABCDEFGH' },
      { label: 'Double_GHEFCDAB', value: 'Double_GHEFCDAB' },
      { label: 'Double_BADCFEHG', value: 'Double_BADCFEHG' },
      { label: 'Double_HGFEDCBA', value: 'Double_HGFEDCBA' },
    ],
  },
  selectedNodes: [],
  uFlag: false,
})
// 获取寄存器列表
const getRegList = (flag) => {
  const pData = {
    token: users.token,
    data: {
      serviceName: props.curGateway.serviceName,
    },
  }
  ServiceApi.getRegByServiceIdList(pData).then(async (res) => {
    console.log('getRegByServiceIdList -> res ', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.regTableData = res.data
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
getRegList()
const contentRef = ref(null)
//获取接口列表
const getInterfaceList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  InterfaceApi.getInterfaceList(pData).then((res) => {
    console.log('getInterfaceList -> res', res)
    if (res.code === '0') {
      ctxData.collOptions = []
      res.data.forEach((item) => {
        const temp = {
          label: item.collInterfaceName,
          name: item.collInterfaceName,
        }
        ctxData.collOptions.push(temp)
      })
    } else {
      showOneResMsg(res)
    }
  })
}
getInterfaceList()
// 选择采集接口触发事件
const selectCollName = (val) => {
  getCollDevices(val).then((res) => {
    handleDevice(res)
  })
  ctxData.propertyOptions = []
}
//处理设备
const handleDevice = (res) => {
  if (res.code === '0') {
    ctxData.nodeOptions = []
    ctxData.nameTotsl = {}
    res.data.deviceNodeMap.forEach((item) => {
      const temp = {
        label: item.label,
        name: item.name,
        tsl: item.tsl,
      }
      ctxData.nameTotsl['' + item.name] = item.tsl
      ctxData.nodeOptions.push(temp)
    })
  } else {
    showOneResMsg(res)
  }
}
// 获取采集接口下的设备列表
const getCollDevices = (collName) => {
  const pData = {
    token: users.token,
    data: {
      name: collName,
    },
  }
  return InterfaceApi.getCollDevices(pData)
  // InterfaceApi.getCollDevices(pData).then((res) => {
  //   if (res.code === '0') {
  //     ctxData.nodeOptions = []
  //     ctxData.nameTotsl = {}
  //     res.data.deviceNodeMap.forEach((item) => {
  //       const temp = {
  //         label: item.label,
  //         name: item.name,
  //         tsl: item.tsl,
  //       }
  //       ctxData.nameTotsl['' + item.name] = item.tsl
  //       ctxData.nodeOptions.push(temp)
  //     })
  //   } else {
  //     showOneResMsg(res)
  //   }
  // })
}
// 选择采集接口触发事件
const selectNodeName = (nodeName) => {
  getNodeProperty(nodeName)
}
// 获取设备属性列表
const getNodeProperty = (nodeName) => {
  const pData = {
    token: users.token,
    data: {
      name: ctxData.nameTotsl[nodeName],
    },
  }
  DeviceModelApi.getDeviceModelProperty(pData).then((res) => {
    console.log('getDeviceModelProperty -> res ', res)
    if (res.code === '0') {
      ctxData.propertyOptions = []
      res.data.forEach((item) => {
        const temp = {
          label: item.label,
          name: item.name,
        }
        ctxData.propertyOptions.push(temp)
      })
    } else {
      showOneResMsg(res)
    }
  })
}
// 处理复选框事件
const handleSelectionChange = (val) => {
  ctxData.selectedNodes = val
}
// 处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
// 过滤表格数据
const filterTableData = computed(() => {
  let regInfo = ctxData.regInfo
  return ctxData.regTableData
    .filter(
      (item) =>
        !regInfo ||
        item.deviceName.toLowerCase().includes(regInfo.toLowerCase()) ||
        item.deviceAddr.toLowerCase().includes(regInfo.toLowerCase())
    )
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  let regInfo = ctxData.regInfo
  return ctxData.regTableData.filter(
    (item) =>
      !regInfo ||
      item.deviceName.toLowerCase().includes(regInfo.toLowerCase()) ||
      item.deviceAddr.toLowerCase().includes(regInfo.toLowerCase())
  )
})
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
const refresh = () => {
  getRegList(1)
}

// 添加寄存器
const addReg = () => {
  ctxData.rFlag = true
  ctxData.rTitle = '添加寄存器'
  initRegForm()
}
const editReg = (row) => {
  ctxData.rFlag = true
  ctxData.rTitle = '编辑寄存器'
  ctxData.regForm = {
    regName: row.regName,
    label: row.label,
    propertyType: row.propertyType,
    collName: row.collName,
    nodeName: row.nodeName,
    propertyName: row.propertyName,
    regType: row.regType,
    regAddr: row.regAddr,
    regCnt: row.regCnt,
    rule: row.rule,
  }
  getCollDevices(ctxData.regForm.collName).then((res) => {
    handleDevice(res)
    getNodeProperty(ctxData.regForm.nodeName)
  })
}
// 初始化寄存器表单
const initRegForm = () => {
  ctxData.regForm = {
    regName: '', // 寄存器名称
    label: '', // 寄存器标签
    propertyType: 0, // 属性类型
    collName: '', // 采集接口名称
    nodeName: '', // 设备名称
    propertyName: '', //设备属性名称
    regType: 0, // 寄存器类型
    regAddr: 0, // 寄存器地址
    regCnt: 1, //寄存器数量
    rule: 'Int_AB',
  }
}
//修改寄存器长度
const changeRegCnt = () => {
  const val = ctxData.regForm.regCnt
  if (val === 1) {
    ctxData.regForm.rule = 'Int_AB'
  }
  if (val === 2) {
    ctxData.regForm.rule = 'Long_ABCD'
  }
  if (val === 4) {
    ctxData.regForm.rule = 'Double_ABCDEFGH'
  }
}
const reportNode = () => {
  let dList = []
  if (ctxData.selectedNodes.length === 0) {
    ElMessage.info('请至少选择一个设备！')
    return
  } else {
    ctxData.selectedNodes.forEach((item) => {
      dList.push(item.deviceName)
    })
  }
  const pData = {
    token: users.token,
    data: {
      serviceName: props.curGateway.serviceName,
      deviceNames: dList,
    },
  }
  TransferModelApi.reportNodes(pData).then((res) => {
    if (res.code === '0') {
      ElMessage.success(res.message)
    } else {
      showOneResMsg(res)
    }
  })
}
const regFormRef = ref(null)
//提交添加寄存器
const submitReg = () => {
  regFormRef.value.validate((valid) => {
    if (valid) {
      const pData = {
        token: users.token,
        data: {
          serviceName: props.curGateway.serviceName,
          register: ctxData.regForm,
        },
      }
      if (ctxData.rTitle.includes('添加')) {
        ServiceApi.addReg(pData).then((res) => {
          handleResult(res, getRegList)
          cancelSubmit()
        })
      } else {
        ServiceApi.editReg(pData).then((res) => {
          handleResult(res, getRegList)
          cancelSubmit()
        })
      }
    } else {
      return false
    }
  })
}
//取消添加节点
const cancelSubmit = () => {
  //
  ctxData.rFlag = false
  initRegForm()
}

// 批量删除寄存器
const deleteReg = () => {
  let dList = []
  if (ctxData.selectedNodes.length === 0) {
    ElMessage.info('请至少选择一个寄存器！')
    return
  } else {
    ctxData.selectedNodes.forEach((item) => {
      dList.push(item.regName)
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
          serviceName: props.curGateway.serviceName,
          registers: dList,
        },
      }
      ServiceApi.deleteReg(pData).then((res) => {
        handleResult(res, getRegList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}

// 导出寄存器
const exportReg = () => {
  const pData = {
    token: users.token,
    responseType: 'blob',
    data: {
      serviceName: props.curGateway.serviceName,
    },
  }
  ServiceApi.exportReg(pData).then((res) => {
    console.log('exportReg -> res', res)
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
//导入寄存器
const importReg = () => {
  ctxData.uFlag = true
}
/**
 * 提交上传的插件
 */
const submitUploadReg = () => {
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
  let typeFlag = fileTypeList.includes(file.type)
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
    formData.append('serviceName', props.curGateway.serviceName)
    formData.append('fileName', file)
    const pData = {
      token: users.token,
      contentType: 'multipart/form-data',
      data: formData,
    }
    ServiceApi.importReg(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success(res.message)
        ctxData.uFlag = false
        getRegList()
      } else {
        showOneResMsg(res)
      }
    })
  }
}
//取消上传设备模型文件
const cancelUploadReg = () => {
  ctxData.uFlag = false
  uploadRef.value.clearFiles()
}
// 弹窗右上角关闭事件处理
const beforeCloseUploadReg = () => {
  cancelUploadReg()
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
.dialog-cLeft {
  position: absolute;
  top: 0;
  left: 0;
  width: calc(50% - 12px);
  bottom: 0;
}
.dialog-cRight {
  position: absolute;
  top: 0;
  right: 0;
  left: calc(50% + 12px);
  bottom: 0;
}
.card-box {
  border-radius: 0;
  height: 100%;
  font-size: 14px;
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 100%;
  }
}
.dialog-content {
  overflow: unset;
  padding: 0;
  min-height: 408px;
}
:deep(.el-card__body) {
  position: absolute;
  top: 56px;
  left: 0;
  right: 0;
  bottom: 0;
}
.card-main {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: auto;
  .upload-item {
    display: flex;
    align-items: center;
    height: 40px;
    line-height: 40px;
    width: 100%;
    color: #606266;
    cursor: pointer;
    .checkbox-custom {
      margin: 2px 10px 2px 20px;
    }
    &-text {
      display: flex;
      align-items: center;
      height: 40px;
      line-height: 40px;
      font-size: 14px;
      width: 100%;
    }
  }
  .upload-item:hover {
    background-color: #f5f7fa;
  }
  .isClicked {
    background-color: #eaeefd;
  }
  .upload-item .el-icon {
    margin: 10px;
  }
}
</style>
