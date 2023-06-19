<template>
  <div class="main-container">
    <div class="main">
      <div class="search-bar">
        <el-form :inline="true" ref="searchFormRef" status-icon label-width="120px">
          <el-form-item style="margin-left: 20px;">
            <el-button type="primary" plain @click="toDataService()" style="margin-right: 20px">
              <el-icon class="el-input__icon"><back /></el-icon>
              返回上报服务
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" plain class="right-btn" @click="importDevice()">
              <el-icon class="el-input__icon"><download /></el-icon>
              导入设备
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" plain class="right-btn" @click="exportDevice">
              <el-icon class="el-input__icon"><upload /></el-icon>
              导出设备
            </el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="search-bar" style="display: flex;">
        <div class="title" style="position: relative;margin-right: 40px;height: 40px; padding: 0px 0; justify-content: flex-start">
          <div class="tName">{{ props.curGateway.serviceName }}</div>
        </div>
        <el-form :inline="true" ref="searchFormRef2" status-icon label-width="90px">
          <el-form-item label="">
            <el-input style="width: 200px" placeholder="请输入 地址/名称 过滤" clearable v-model="ctxData.nodeInfo">
              <template #prefix>
                <el-icon class="el-input__icon"><search /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" bg class="right-btn" @click="reportNode()">
              <el-icon class="btn-icon">
                <Icon name="local-report" size="14px" color="#ffffff" />
              </el-icon>
              主动上报
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" bg class="right-btn" @click="addNode()">
              <el-icon class="btn-icon">
                <Icon name="local-add" size="14px" color="#ffffff" />
              </el-icon>
              添加
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-button type="danger" bg class="right-btn" @click="deleteNode()">
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
          @row-dblclick="editNode"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column type="expand">
            <template #default="scope">
              <div class="param-content">
                <div class="pc-title">
                  <div class="pct-info">
                    <b> {{ scope.row.deviceName }} </b>
                    上报参数
                  </div>
                </div>
                <div class="pc-content">
                  <div class="param-item" v-for="(item, key, index) of scope.row.param" :key="index">
                    <div class="param-value">{{ ctxData.paramName[key] }}：</div>
                    <div class="param-name">{{ item || '-' }}</div>
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column sortable prop="deviceName" label="设备名称" width="auto" min-width="180" align="center">
          </el-table-column>
          <el-table-column sortable prop="deviceLabel" label="设备标签" width="auto" min-width="180" align="center">
          </el-table-column>
          <el-table-column sortable prop="deviceAddr" label="通讯地址" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column sortable prop="uploadModel" label="上报模型" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column sortable prop="collInterfaceName" label="采集接口名称" width="auto" min-width="150" align="center">
          </el-table-column>
          <!-- el-table-column sortable label="通信状态" width="auto" min-width="100" align="center">
            <template #default="scope">
              {{ scope.row.commStatus === 'onLine' ? '在线' : '离线' }}
            </template>
          </el-table-column -->

          <el-table-column sortable label="上报状态" width="auto" min-width="100" align="center">
            <template #default="scope">
              {{ scope.row.reportStatus === 'onLine' ? '在线' : '离线' }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="auto" min-width="120" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="editNode(scope.row)" text type="primary">编辑</el-button>
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
    <!-- 添加设备 -->
    <el-dialog v-model="ctxData.nFlag" title="添加设备" width="1000px">
      <div class="dialog-content">
        <div class="dialog-cLeft">
          <el-card class="card-box" shadow="hover">
            <template #header>
              <div class="card-header">
                <span>请选择采集设备</span>
              </div>
            </template>
            <div class="card-main">
              <div><el-input v-model="ctxData.filterText" placeholder="节点名称过滤" /></div>
              <el-tree
                ref="nodeTreeRef"
                :props="ctxData.treeProps"
                :load="loadNode"
                :filter-node-method="filterNode"
                lazy
                show-checkbox
              />
            </div>
          </el-card>
        </div>
        <div class="dialog-cRight">
          <el-card class="card-box" shadow="hover">
            <template #header>
              <div class="card-header">
                <span>请选择一个上报模型</span>
              </div>
            </template>
            <div class="card-main">
              <el-radio-group v-model="ctxData.curModelName" style="width: 100%">
                <div
                  class="upload-item"
                  v-for="item in ctxData.transferModelList"
                  :key="'tml' + item.name"
                  :class="{ isClicked: item.isClicked }"
                >
                  <div class="upload-item-text">
                    <el-radio :label="item.name">{{ item.label + ' (' + item.name + ')' }}</el-radio>
                  </div>
                </div>
              </el-radio-group>
            </div>
          </el-card>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelAddNode()">取消</el-button>
          <el-button type="primary" @click="submitAddNode()">添加</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- 编辑设备 -->
    <el-dialog
      v-model="ctxData.dFlag"
      title="编辑设备"
      width="600px"
      :before-close="handleClose"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.nodeForm"
          :rules="ctxData.nodeRules"
          ref="nodeFormRef"
          status-icon
          label-position="right"
          label-width="120px"
        >
          <el-form-item label="服务名称" prop="serviceName">
            <el-input
              disabled
              type="text"
              v-model="ctxData.nodeForm.serviceName"
              autocomplete="off"
              placeholder="请输入服务名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="采集接口名称" prop="collInterfaceName">
            <el-input
              disabled
              type="text"
              v-model="ctxData.nodeForm.collInterfaceName"
              autocomplete="off"
              placeholder="请输入采集接口名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="设备名称" prop="deviceName">
            <el-input
              disabled
              type="text"
              v-model="ctxData.nodeForm.deviceName"
              autocomplete="off"
              placeholder="请输入设备名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="设备地址" prop="deviceAddr">
            <el-input
              disabled
              type="text"
              v-model="ctxData.nodeForm.deviceAddr"
              autocomplete="off"
              placeholder="请输入设备地址"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="协议名称" prop="protocol">
            <el-input
              disabled
              type="text"
              v-model="ctxData.nodeForm.protocol"
              autocomplete="off"
              placeholder="请输入协议名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="上报模型" prop="uploadModel">
            <el-select v-model="ctxData.nodeForm.uploadModel" style="width: 100%" placeholder="请选择上报模型">
              <el-option
                v-for="item in ctxData.transferModelList"
                :key="'model_' + item.name"
                :label="item.label"
                :value="item.name"
              >
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item v-if="ctxData.nodeForm.protocol != 'FSJY.MQTT' && ctxData.nodeForm.protocol != 'ZXJS.MQTT'" label="通讯编码" prop="deviceCode">
            <el-input type="text" v-model="ctxData.nodeForm.deviceCode" autocomplete="off" placeholder="请输入通讯编码">
            </el-input>
          </el-form-item>

          <el-form-item v-if="ctxData.nodeForm.protocol == 'FSJY.MQTT'" label="产品密钥" prop="productKey">
            <el-input type="text" v-model="ctxData.nodeForm.productKey" autocomplete="off" placeholder="请输入网关产品密钥">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.nodeForm.protocol == 'FSJY.MQTT'" label="通讯地址" prop="deviceID">
            <el-input type="text" v-model="ctxData.nodeForm.deviceID" autocomplete="off" placeholder="请输入网关通讯地址">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.nodeForm.protocol == 'FSJY.MQTT'" label="设备密钥" prop="deviceSecret">
            <el-input type="text" v-model="ctxData.nodeForm.deviceSecret" autocomplete="off" placeholder="请输入网关设备密钥">
            </el-input>
          </el-form-item>

          <el-form-item v-if="ctxData.nodeForm.protocol == 'ZXJS.MQTT'" label="产品序列号" prop="productSn">
            <el-input type="text" v-model="ctxData.nodeForm.productSn" autocomplete="off" placeholder="请输入产品序列号">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.nodeForm.protocol == 'ZXJS.MQTT'" label="设备序列号" prop="deviceSn">
            <el-input type="text" v-model="ctxData.nodeForm.deviceSn" autocomplete="off" placeholder="请输入设备序列号">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.nodeForm.protocol == 'ZXJS.MQTT'" label="设备密码" prop="devicePwd">
            <el-input type="text" v-model="ctxData.nodeForm.devicePwd" autocomplete="off" placeholder="请输入设备密码">
            </el-input>
          </el-form-item>



        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelSubmit()">取消</el-button>
          <el-button type="primary" @click="submitNodeForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- 批量添加设备 -->
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
import InterfaceApi from 'api/interface.js'
import ServiceApi from 'api/service.js'
import TransferModelApi from 'api/transferModel.js'
import { userStore } from 'stores/user'
const users = userStore()
const props = defineProps({
  curGateway: {
    type: Object,
    default: {},
  },
  pageInfo: String,
})
console.log('id -> props', props)

const emit = defineEmits(['changeDnFlag'])
const toDataService = () => {
  //lp update 2023-06-12 首页上报信息查看详情，对点击返回按钮进行操作标识，防止死循环跳转显示详情页面
  emit('changeDnFlag', 'goBack')
}
const ctxData = reactive({
  headerCellStyle: {
    background: variables.primaryColor,
    color: variables.fontWhiteColor,
    height: '54px',
  },
  paramName: {
    ProductKey: '产品密钥',
    DeviceID: '通讯地址',
    DeviceSecret: '设备密钥',
    ProductSn: '产品序列号',
    DeviceSn: '设备序列号',
    DevicePwd: '设备密码',
  },
  cellStyle: {
    height: '48px',
  },
  tableMaxHeight: 0,
  currentPage: 1, // 默认当前页是第一页
  pagesize: 20, // 每页数据个数
  nodeTableData: [],
  dFlag: false,
  dTitle: '添加设备',
  nodeInfo: '',
  nodeForm: {
    serviceName: '', // 服务名称
    collInterfaceName: '', // 采集接口名称
    deviceName: '', // 设备名称
    uploadModel: '', //上报模型
    deviceAddr: '', // 设备地址
    protocol: '', // 协议名称
    deviceCode: '', //通讯编码
    //gwai add 2023-04-05
    productKey: '',
    deviceID: '',
    deviceSecret: '',
    productSn: '',
    deviceSn: '',
    devicePwd: '',
  },
  nodeRules: {
    serviceName: [
      {
        required: true,
        trigger: 'blur',
      },
    ],
    collInterfaceName: [
      {
        required: true,
        trigger: 'blur',
      },
    ],
    deviceName: [
      {
        required: true,
        trigger: 'blur',
      },
    ],
    deviceAddr: [
      {
        required: true,
        trigger: 'blur',
      },
    ],
    protocol: [
      {
        required: true,
        trigger: 'blur',
      },
    ],
    uploadModel: [
      {
        required: true,
        message: '上报模型不能为空',
        trigger: 'blur',
      },
    ],
    deviceCode: [
      {
        required: true,
        message: '通讯编码不能为空',
        trigger: 'blur',
      },
    ],
    //gwai add 2023-04-05
    productKey: [
      {
        required: true,
        message: '产品密钥不能为空',
        trigger: 'blur',
      },
    ],
    deviceID: [
      {
        required: true,
        message: '产品密钥不能为空',
        trigger: 'blur',
      },
    ],
    deviceSecret: [
      {
        required: true,
        message: '产品密钥不能为空',
        trigger: 'blur',
      },
    ],
    productSn: [
      {
        required: true,
        message: '产品序列号不能为空',
        trigger: 'blur',
      },
    ],
    deviceSn: [
      {
        required: true,
        message: '设备序列号不能为空',
        trigger: 'blur',
      },
    ],
    devicePwd: [
      {
        required: true,
        message: '设备密码',
        trigger: 'blur',
      },
    ],

  },
  nodeModelList: [],
  selectedNodes: [],
  dpFlag: true, //设备-属性切换显示表示
  interfaceList: [], //采集接口列表
  nFlag: false, //添加弹框标识
  treeProps: {
    children: 'children',
    label: 'label',
    isLeaf: 'leaf',
  },
  deviceList: [], //某个采集接口下的设备列表
  interfaceNames: [],
  uFlag: false,
  filterText: '',
  transferModelList: [],
  curModelName: '',
})

const getNodeList = (flag) => {
  const pData = {
    token: users.token,
    data: {
      serviceName: props.curGateway.serviceName,
    },
  }
  ServiceApi.getDeviceByServiceIdList(pData).then(async (res) => {
    console.log('getDeviceByServiceIdList -> res ', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.nodeTableData = res.data
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
getNodeList()
const contentRef = ref(null)
const getInterfaceList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  InterfaceApi.getInterfaceList(pData).then((res) => {
    console.log('getInterfaceList -> res', res)
    if (res.code === '0') {
      ctxData.interfaceList = []
      ctxData.interfaceNames = []
      res.data.forEach((item) => {
        const temp = {
          label: item.collInterfaceName,
          name: item.collInterfaceName,
          hasChildren: true,
        }
        ctxData.interfaceList.push(temp)
        ctxData.interfaceNames.push(item.collInterfaceName)
      })
    } else {
      showOneResMsg(res)
    }
  })
}
getInterfaceList()

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
  let nodeInfo = ctxData.nodeInfo
  return ctxData.nodeTableData
    .filter(
      (item) =>
        !nodeInfo ||
        item.deviceName.toLowerCase().includes(nodeInfo.toLowerCase()) ||
        item.deviceAddr.toLowerCase().includes(nodeInfo.toLowerCase())
    )
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  let nodeInfo = ctxData.nodeInfo
  return ctxData.nodeTableData.filter(
    (item) =>
      !nodeInfo ||
      item.deviceName.toLowerCase().includes(nodeInfo.toLowerCase()) ||
      item.deviceAddr.toLowerCase().includes(nodeInfo.toLowerCase())
  )
})
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
const refresh = () => {
  getNodeList(1)
}
// 加载node tree
const loadNode = (node, resolve) => {
  //
  console.log('loadNode-> node === ', node)
  if (node.level === 0) {
    return resolve(ctxData.interfaceList)
  } else if (node.level === 1) {
    console.log('node->', node)
    getCollDevices(node.data.name).then(async (res) => {
      console.log('getCollDevices -> res', res)
      if (res.code === '0') {
        ctxData.deviceList = []
        if (res.data.deviceNodeMap) {
          res.data.deviceNodeMap.forEach((item) => {
            const temp = {
              label: item.label,
              name: item.name,
              addr: item.addr,
              collName: node.data.name,
              hasChildren: false,
            }
            ctxData.deviceList.push(temp)
          })
        }
        console.log('node->ctxData.deviceList', ctxData.deviceList)
        return resolve(ctxData.deviceList)
      } else {
        showOneResMsg(res)
      }
      return resolve([])
    })
  } else {
    return resolve([])
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
}

// 获取上报模型列表
const getTransferList = (flag) => {
  const pData = {
    token: users.token,
    data: {},
  }
  TransferModelApi.getModelList(pData).then((res) => {
    console.log('getModelList -> res', res)
    if (res.code === '0') {
      ctxData.transferModelList = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
getTransferList()

// 添加设备
const addNode = () => {
  ctxData.nFlag = true
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

//提交添加节点
const nodeTreeRef = ref(null)
const submitAddNode = () => {
  //
  var checkedNodes = nodeTreeRef.value.getCheckedNodes()
  checkedNodes = checkedNodes.filter((item) => {
    return !item.hasChildren
  })
  console.log('checkedNodes => ', checkedNodes)
  console.log('props.curGateway => ', props.curGateway)
  if (checkedNodes.length === 0) {
    ElMessage.info('请先展开节点，获取节点下的子设备数据！')
    return
  } else {
    let count = 0
    let success = 0
    checkedNodes.forEach((item) => {
      if (item.collName) {
        const pData = {
          token: users.token,
          data: {
            serviceName: props.curGateway.serviceName,
            collInterfaceName: item.collName,
            deviceName: item.name,
            deviceLabel: item.label,
            deviceAddr: item.addr,
            protocol: props.curGateway.protocol,
            uploadModel: ctxData.curModelName,
            param: {},
          },
        }
        console.log('pData item -> ', item)
        console.log('pData', pData)
        ServiceApi.addDevice(pData).then((res) => {
          count++
          console.log('addDevice -> res ', res)
          if (res.code !== '0') {
            showOneResMsg(res)
          } else {
            success++
          }
          if (count === checkedNodes.length) {
            ElMessage.success('添加' + count + '个设备，成功' + success + '个，失败' + (count - success) + '个。')
            getNodeList()
          }
        })
      }
    })

    cancelAddNode()
  }
}

watch(
  () => ctxData.filterText,
  (val) => {
    nodeTreeRef.value && nodeTreeRef.value.filter(val)
  }
)

const filterNode = (value, data) => {
  console.log('filterNode value', value)
  if (!value) return true
  return data.label.includes(value)
}
//取消添加节点
const cancelAddNode = () => {
  //
  ctxData.nFlag = false
}
// 编辑设备
const editNode = (row) => {
  ctxData.dFlag = true
  ctxData.dTitle = '编辑设备'
  ctxData.nodeForm = {
    serviceName: row.serviceName, // 服务名称
    collInterfaceName: row.collInterfaceName, // 采集接口名称
    deviceName: row.deviceName, // 设备名称
    deviceLabel: row.deviceLabel,
    uploadModel: row.uploadModel,
    deviceAddr: row.deviceAddr, // 设备地址
    protocol: row.protocol, // 协议名称
    deviceCode: row.param.deviceCode === '' ? row.deviceAddr : row.param.deviceCode, //通讯编码

    //gwai add 2023-04-05
    //productKey: row.productKey,
    //deviceID: row.deviceID,
    //deviceSecret: row.deviceSecret,
  }
  
  ctxData.nodeForm['productKey'] = row.param.ProductKey
  ctxData.nodeForm['deviceID'] = row.param.DeviceID
  ctxData.nodeForm['deviceSecret'] = row.param.DeviceSecret

  ctxData.nodeForm['productSn'] = row.param.ProductSn
  ctxData.nodeForm['deviceSn'] = row.param.DeviceSn
  ctxData.nodeForm['devicePwd'] = row.param.DevicePwd

}
const nodeFormRef = ref(null)
const submitNodeForm = () => {
  nodeFormRef.value.validate((valid) => {
    console.log('ctxData.nodeForm', ctxData.nodeForm)
    if (valid) {
      const nForm = {
        serviceName: ctxData.nodeForm.serviceName, // 服务名称
        collInterfaceName: ctxData.nodeForm.collInterfaceName, // 采集接口名称
        uploadModel: ctxData.nodeForm.uploadModel, // 设备名称
        deviceName: ctxData.nodeForm.deviceName, // 设备名称
        deviceLabel: ctxData.nodeForm.deviceLabel, // 设备标签
        deviceAddr: ctxData.nodeForm.deviceAddr, // 设备地址
        protocol: ctxData.nodeForm.protocol, // 协议名称
        param: {
          deviceCode: ctxData.nodeForm.deviceCode, //通讯编码
          //gwai add 2023-04-05
          productKey: ctxData.nodeForm.productKey,
          deviceID: ctxData.nodeForm.deviceID,
          deviceSecret: ctxData.nodeForm.deviceSecret,
          productSn: ctxData.nodeForm.productSn,
          deviceSn: ctxData.nodeForm.deviceSn,
          devicePwd: ctxData.nodeForm.devicePwd,
        },
      }
      const pData = {
        token: users.token,
        data: nForm,
      }
      ServiceApi.editDevice(pData).then((res) => {
        handleResult(res, getNodeList)
        cancelSubmit()
      })
    } else {
      return false
    }
  })
}
// 取消提交
const cancelSubmit = () => {
  ctxData.dFlag = false
  nodeFormRef.value.resetFields()
}
//处理弹出框右上角关闭图标
const handleClose = (done) => {
  cancelSubmit()
}

// 批量删除设备
const deleteNode = () => {
  let dList = []
  if (ctxData.selectedNodes.length === 0) {
    ElMessage.info('请至少选择一个设备！')
    return
  } else {
    ctxData.selectedNodes.forEach((item) => {
      dList.push(item.deviceName)
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
          deviceNames: dList,
        },
      }
      ServiceApi.deleteDevice(pData).then((res) => {
        handleResult(res, getNodeList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}

// 批量导入设备
const importDevice = (row) => {
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
    formData.append('file', file)
    const pData = {
      token: users.token,
      contentType: 'multipart/form-data',
      data: formData,
    }
    ServiceApi.addDeviceFromCSV(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success(res.message)
        ctxData.uFlag = false
        getNodeList()
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
// 导出设备
const exportDevice = () => {
  const pData = {
    token: users.token,
    responseType: 'blob',
    data: {
      serviceName: props.curGateway.serviceName,
    },
  }
  ServiceApi.exportDevice(pData).then((res) => {
    console.log('exportDevice -> res', res)
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
