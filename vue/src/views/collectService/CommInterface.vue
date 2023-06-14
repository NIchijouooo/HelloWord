<template>
  <div class="main-container">
    <div class="main">
      <div class="search-bar">
        <el-form :inline="true" ref="searchFormRef" status-icon label-width="90px">
          <el-form-item label="接口名称">
            <el-input style="width: 200px" placeholder="请输入接口名称" clearable v-model="ctxData.interfaceName">
              <template #prefix>
                <el-icon class="el-input__icon"><search /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button style="color: #fff; margin-left: 20px" color="#2EA554" class="right-btn" @click="refresh()">
              <el-icon class="btn-icon">
                <Icon name="local-refresh" size="14px" color="#ffffff" />
              </el-icon>
              刷新
            </el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="tool-bar">
        <el-button type="primary" bg class="right-btn" @click="addInterface()">
          <el-icon class="btn-icon">
            <Icon name="local-add" size="14px" color="#ffffff" />
          </el-icon>
          添加
        </el-button>
      </div>
      <div class="content-up">
        <div class="cu-main">
          <el-collapse v-model="ctxData.activeNames" @change="handleChange">
            <el-collapse-item v-for="(comm, cKey, cIndex) in filterTableDataPage" :key="cIndex" :name="cKey">
              <template #title>
                <div class="tName">
                  {{ ctxData.typeName[cKey] }}
                </div>
              </template>
              <div class="cu-items">
                <div v-for="(item, key, index) in comm.data" :key="'interface' + index" class="cu-item">
                  <div class="cui-content">
                    <el-card class="box-card" shadow="hover">
                      <template #header>
                        <div class="card-header">
                          <div class="cui-header">
                            <span style="font-weight: 600">{{ item.name }}</span>
                            <el-tag size="small" style="margin-left: 8px" :type="item.status ? 'success' : 'danger'">{{
                              item.status === 1 ? '连接' : '未连接'
                            }}</el-tag>
                          </div>
                          <div>
                            <el-button class="head-tag" text type="danger" @click="deleteInterface(item)" size="small"
                              >删除</el-button
                            >
                            <el-button class="head-tag" text type="primary" @click="editInterface(item)" size="small"
                              >编辑</el-button
                            >
                          </div>
                        </div>
                      </template>
                      <div class="card-content">
                        <div class="cc-body">
                          <div v-for="(param, key1, index) in item.param" class="ccb-item" :key="index">
                            <div class="">{{ ctxData.paramName[key1] }}：</div>
                            <div class="head-name" v-if="key1 == 'ledModuleEnable'">{{ param ? '是' : '否' }}</div>
                            <div class="head-name" v-else>{{ param }}</div>
                          </div>
                        </div>
                      </div>
                    </el-card>
                  </div>
                </div>
              </div>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>
    </div>
    <!-- 添加编辑通讯接口 -->
    <el-dialog
      v-model="ctxData.tFlag"
      :title="ctxData.tTitle"
      width="600px"
      :before-close="handleClose"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.interfaceForm"
          :rules="ctxData.interfaceRules"
          ref="interfaceFormRef"
          status-icon
          label-position="right"
          label-width="120px"
        >
          <el-form-item label="接口名称" prop="name">
            <el-input
              :disabled="ctxData.tTitle.includes('编辑')"
              type="text"
              v-model="ctxData.interfaceForm.name"
              autocomplete="off"
              placeholder="请输入接口名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="接口类型" prop="type">
            <el-select v-model="ctxData.interfaceForm.type" style="width: 100%" placeholder="请选择接口类型">
              <el-option
                v-for="item in ctxData.ciOptions"
                :key="'type' + item.name"
                :label="item.label"
                :value="item.name"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="ctxData.serialList.includes(ctxData.interfaceForm.type)" label="串口名称" prop="pName">
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.pName"
              autocomplete="off"
              placeholder="请输入串口名称，例如：/dev/ttyS1或者COM1"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.serialList.includes(ctxData.interfaceForm.type)" label="波特率" prop="pBaudRate">
            <el-select v-model="ctxData.interfaceForm.pBaudRate" style="width: 100%" placeholder="请选择波特率">
              <el-option
                v-for="item in ctxData.baudRateOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="ctxData.serialList.includes(ctxData.interfaceForm.type)" label="数据位" prop="pDataBits">
            <el-select v-model="ctxData.interfaceForm.pDataBits" style="width: 100%" placeholder="请选择数据位">
              <el-option
                v-for="item in ctxData.dataBitsOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="ctxData.serialList.includes(ctxData.interfaceForm.type)" label="停止位" prop="pStopBits">
            <el-select v-model="ctxData.interfaceForm.pStopBits" style="width: 100%" placeholder="请选择停止位">
              <el-option
                v-for="item in ctxData.stopBitsOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="ctxData.serialList.includes(ctxData.interfaceForm.type)" label="校验位" prop="pParity">
            <el-select v-model="ctxData.interfaceForm.pParity" style="width: 100%" placeholder="请选择校验位">
              <el-option
                v-for="item in ctxData.parityOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="ctxData.tcpUdpList.includes(ctxData.interfaceForm.type)" label="IP地址" prop="ip">
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.ip"
              autocomplete="off"
              :placeholder="'请输入' + (ctxData.interfaceForm.type.includes('Client') ? '远端' : '本地') + 'IP地址'"
            >
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="
              ctxData.tcpUdpList.includes(ctxData.interfaceForm.type) ||
              ctxData.interfaceForm.type === 'sac009' ||
              ctxData.interfaceForm.type === 'dr504' ||
              ctxData.interfaceForm.type === 'hdkj'
            "
            label="端口"
            prop="port"
          >
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.port"
              autocomplete="off"
              :placeholder="
                '请输入' +
                (ctxData.portList.includes(ctxData.interfaceForm.type)
                  ? '目标'
                  : ctxData.interfaceForm.type.includes('Client')
                  ? '远端'
                  : '本地') +
                '端口'
              "
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.ioList.includes(ctxData.interfaceForm.type)" label="IO端口名称" prop="ioName">
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.ioName"
              autocomplete="off"
              placeholder="请输入IO端口名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.interfaceForm.type === 'hdkj'" label="间隔时间" prop="reportInterval">
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.reportInterval"
              autocomplete="off"
              placeholder="请输入间隔时间"
            >
              <template #append>单位分钟</template>
            </el-input>
          </el-form-item>

          <el-form-item v-if="ctxData.interfaceForm.type === 'hdkj'" label="离线判断次数" prop="offlinePeriod">
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.offlinePeriod"
              autocomplete="off"
              placeholder="请输入离线判断次数"
            >
              <template #append>单位次</template>
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="
              ctxData.interfaceForm.type !== '' &&
              !ctxData.ioList.includes(ctxData.interfaceForm.type) &&
              ctxData.interfaceForm.type !== 'hdkj'
            "
            label="超时时间"
            prop="timeout"
          >
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.timeout"
              autocomplete="off"
              placeholder="请输入超时时间"
            >
              <template #append>单位毫秒</template>
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="ctxData.intervalList.includes(ctxData.interfaceForm.type)"
            label="间隔时间"
            prop="interval"
          >
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.interval"
              autocomplete="off"
              placeholder="请输入间隔时间"
            >
              <template #append>单位毫秒</template>
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.interfaceForm.type === 's7'" label="连接时间" prop="connect">
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.connect"
              autocomplete="off"
              placeholder="请输入连接时间"
            >
              <template #append>单位毫秒</template>
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="ctxData.serialList.includes(ctxData.interfaceForm.type)"
            label="启用led模块"
            prop="ledModuleEnable"
          >
            <el-switch v-model="ctxData.interfaceForm.ledModuleEnable" />
          </el-form-item>
          <el-form-item
            v-if="
              ctxData.serialList.includes(ctxData.interfaceForm.type) && ctxData.interfaceForm.ledModuleEnable == true
            "
            label="LED管脚编号"
            prop="ledGPIO"
          >
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.ledGPIO"
              autocomplete="off"
              placeholder="请输入LED管脚编号"
            >
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="
              ctxData.serialList.includes(ctxData.interfaceForm.type) && ctxData.interfaceForm.ledModuleEnable == true
            "
            label="LED点亮值"
            prop="ledOn"
          >
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.ledOn"
              autocomplete="off"
              placeholder="请输入LED点亮值"
            >
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="
              ctxData.serialList.includes(ctxData.interfaceForm.type) && ctxData.interfaceForm.ledModuleEnable == true
            "
            label="LED熄灭值"
            prop="ledOff"
          >
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.ledOff"
              autocomplete="off"
              placeholder="请输入LED熄灭值"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.interfaceForm.type == 'httpSmartNode'" label="应用编号" prop="appCode">
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.appCode"
              autocomplete="off"
              placeholder="请输入应用编号"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.interfaceForm.type == 'httpSmartNode'" label="应用秘钥" prop="appKey">
            <el-input
              type="text"
              v-model="ctxData.interfaceForm.appKey"
              autocomplete="off"
              placeholder="请输入应用秘钥"
            >
            </el-input>
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelSubmit()">取消</el-button>
          <el-button type="primary" @click="submitInterfaceForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { Search } from '@element-plus/icons-vue'
import CommInterfaceApi from 'api/commInterface.js'
import { userStore } from 'stores/user'
const users = userStore()

// 验证方法start
const regCnt = /^[0-9]*[1-9][0-9]*$/
const validateRegCnt = (rule, value, callback) => {
  if (!regCnt.test(value)) {
    callback(new Error('只能输入正整数数字！'))
  } else {
    callback()
  }
}
const validateNumber = (rule, value, callback) => {
  if (isNaN(Number(value))) {
    callback(new Error('只能输入数字！'))
  } else {
    callback()
  }
}
const refExpIP = /^((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))$/
const validateIP = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('IP地址不能为空！'))
  } else {
    if (!refExpIP.test(value)) {
      callback(new Error('IP地址格式错误！'))
    } else {
      callback()
    }
  }
}
// 验证方法end
const ctxData = reactive({
  currentPage: 1, // 默认当前页是第一页
  pagesize: 20, // 每页数据个数
  interfaceName: '',
  commInterfaceData: [],
  tFlag: false,
  tTitle: '添加通讯接口',
  activeNames: [],
  interfaceForm: {
    name: '',
    type: '',
    //串口
    pName: '', //串口名称，比如"/dev/ttyS1"或者"COM1"
    pBaudRate: '', //波特率
    pDataBits: '', //数据位，5，6，7，8
    pStopBits: '', //停止位，1，1.5，2
    pParity: '', //校验位，N（无校验），O（奇校验），E（偶校验）
    timeout: '', //超时时间（毫秒）
    interval: '', //间隔时间（毫秒）
    ledModuleEnable: false, //启用led模块
    ledGPIO: '', //LED管脚编号
    ledOn: '', //LED点亮
    ledOff: '', //LED熄灭
    //tcp/udp
    ip: '', //本地（远端）IP地址
    port: '', //本地（远端）端口
    //io
    ioName: '', //IO端口名称
    //s7
    connect: '', //连接时间
    reportInterval: '',
    offlinePeriod: '',
    //SmartNode
    appCode: '', //应用编码
    appKey: '', //应用密钥
  },
  commListObj: {},
  interfaceFormUpload: {
    name: '',
    type: '',
    param: {},
  },
  interfaceRules: {
    name: [
      {
        required: true,
        message: '接口名称不能为空',
        trigger: 'blur',
      },
    ],
    type: [
      {
        required: true,
        message: '通讯接口类型不能为空',
        trigger: 'blur',
      },
    ],
    pName: [
      {
        required: true,
        message: '串口名称不能为空',
        trigger: 'blur',
      },
    ],
    pBaudRate: [
      {
        required: true,
        message: '波特率不能为空',
        trigger: 'blur',
      },
    ],
    pDataBits: [
      {
        required: true,
        message: '数据位不能为空',
        trigger: 'blur',
      },
    ],
    pStopBits: [
      {
        required: true,
        message: '停止位不能为空',
        trigger: 'blur',
      },
    ],
    pParity: [
      {
        required: true,
        message: '校验位不能为空',
        trigger: 'blur',
      },
    ],
    ip: [
      {
        required: true,
        validator: validateIP,
        trigger: 'blur',
      },
    ],
    port: [
      {
        required: true,
        message: '端口不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateNumber,
      },
    ],
    ioName: [
      {
        required: true,
        message: 'IO端口名称不能为空',
        trigger: 'blur',
      },
    ],
    ledGPIO: [
      {
        required: true,
        message: 'LED管脚编号不能为空',
        trigger: 'blur',
      },
    ],
    ledOn: [
      {
        required: true,
        message: 'LED点亮值不能为空',
        trigger: 'blur',
      },
    ],
    ledOff: [
      {
        required: true,
        message: 'LED熄灭值不能为空',
        trigger: 'blur',
      },
    ],
    timeout: [
      {
        required: true,
        message: '超时时间不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateNumber,
      },
    ],
    interval: [
      {
        required: true,
        message: '间隔时间不能为空',
        trigger: 'blur',
      },
      {
        validator: validateNumber,
        trigger: 'blur',
      },
    ],
    reportInterval: [
      {
        required: true,
        message: '间隔时间不能为空',
        trigger: 'blur',
      },
      {
        validator: validateNumber,
        trigger: 'blur',
      },
    ],
    connect: [
      {
        required: true,
        message: '连接时间不能为空',
        trigger: 'blur',
      },
      {
        validator: validateNumber,
        trigger: 'blur',
      },
    ],
    offlinePeriod: [
      {
        required: true,
        message: '离线判断次数',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateRegCnt,
      },
    ],
    appCode: [
      {
        required: true,
        message: '应用编码不能为空',
        trigger: 'blur',
      },
    ],
    appKey: [
      {
        required: true,
        message: '应用秘钥不能为空',
        trigger: 'blur',
      },
    ],
  },
  typeName: {
    localSerial: '本地串口',
    tcpClient: 'TCP客户端',
    tcpServer: 'TCP服务端',
    udpClient: 'UPD客户端',
    udpServer: 'UDP服务端',
    ioOut: '开关量输出',
    ioIn: '开关量输入',
    httpTianGang: '天罡NBIOT水表',
    s7: '西门子PLC协议',
    sac009: '天睿信诚分体空调',
    dr504: '济南有人DTU',
    mbTCPClient: 'ModbusTCP客户端',
    mbRTUClient: 'ModbusRTU客户端',
    hdkj: '和达科技水表',
    httpSmartNode: 'SmartNodeAPI协议',
  },
  paramName: {
    name: '串口名称',
    baudRate: '波特率',
    dataBits: '数据位',
    stopBits: '停止位',
    parity: '校验位',
    timeout: '超时时间（单位毫秒）',
    interval: '间隔时间（单位毫秒）',
    ip: 'IP地址',
    port: '端口',
    connect: '连接时间',
    reportInterval: '间隔时间（单位分钟）',
    offlinePeriod: '离线判断次数',
    ledModuleEnable: '启用LED模块',
    ledGPIO: 'LED管脚编号',
    ledOn: 'LED点亮',
    ledOff: 'LED熄灭',
    appCode: '应用编码',
    appKey: '应用秘钥',
  },
  ciOptions: [
    {
      name: 'localSerial',
      label: '本地串口',
    },
    {
      name: 'tcpClient',
      label: 'TCP客户端',
    },
    {
      name: 'tcpServer',
      label: 'TCP服务端',
    },
    {
      name: 'udpClient',
      label: 'UDP客户端',
    },
    {
      name: 'udpServer',
      label: 'UDP服务端',
    },
    {
      name: 'ioOut',
      label: '开关量输出',
    },
    {
      name: 'ioIn',
      label: '开关量输入',
    },
    {
      name: 'httpTianGang',
      label: '天罡NBIOT水表',
    },
    {
      name: 's7',
      label: '西门子PLC协议',
    },
    {
      name: 'sac009',
      label: '天睿信诚分体空调',
    },
    {
      name: 'dr504',
      label: '济南有人DTU',
    },
    {
      name: 'mbTCPClient',
      label: 'ModbusTCP客户端',
    },
    {
      name: 'mbRTUClient',
      label: 'ModbusRTU客户端',
    },
    {
      value: 'hdkj',
      label: '和达科技水表',
    },
    {
      value: 'httpSmartNode',
      label: 'SmartNodeAPI协议',
    },
  ],
  baudRateOptions: [
    { value: '1200', label: '1200' },
    { value: '2400', label: '2400' },
    { value: '4800', label: '4800' },
    { value: '9600', label: '9600' },
    { value: '14400', label: '14400' },
    { value: '19200', label: '19200' },
    { value: '115200', label: '115200' },
  ],
  dataBitsOptions: [
    { value: '5', label: '5' },
    { value: '6', label: '6' },
    { value: '7', label: '7' },
    { value: '8', label: '8' },
  ],
  stopBitsOptions: [
    { value: '1', label: '1' },
    { value: '1.5', label: '1.5' },
    { value: '2', label: '2' },
  ],
  parityOptions: [
    { value: 'N', label: '无校验' },
    { value: 'O', label: '奇校验' },
    { value: 'E', label: '偶校验' },
  ],
  serialList: ['localSerial', 'mbRTUClient'],
  tcpUdpList: ['tcpClient', 'tcpServer', 'udpClient', 'udpServer', 's7', 'mbTCPClient', 'httpSmartNode'],
  ioList: ['ioOut', 'ioIn'],
  intervalList: [
    'localSerial',
    'tcpClient',
    'tcpServer',
    'udpClient',
    'udpServer',
    'mbTCPClient',
    'mbRTUClient',
    'httpSmartNode',
  ],
  portList: ['s7', 'sac009', 'dr504'],
  curRow: null,
  sFlag: false,
})
//获取当前网关通讯接口的通讯协议
const getCommProtocolList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  CommInterfaceApi.getCommProtocolList(pData).then((res) => {
    if (!res) return
    console.log('getCommProtocolList -> res', res)
    if (res.code === '0') {
      ctxData.ciOptions = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
getCommProtocolList()
// 获取接口列表
const getInterfaceList = (flag) => {
  const pData = {
    token: users.token,
    data: {},
  }
  CommInterfaceApi.getCommInterfaceList(pData).then((res) => {
    if (!res) return
    console.log('getCommInterfaceList -> res', res)
    if (res.code === '0') {
      ctxData.commInterfaceData = res.data
      ctxData.commInterfaceData.forEach((item) => {
        if (ctxData.activeNames.indexOf(item.type) == -1) {
          ctxData.activeNames.push(item.type)
        }
      })
      if (flag === 1) {
        ElMessage({
          type: 'success',
          message: '刷新成功！',
        })
      }
    } else {
      showOneResMsg(res)
    }
  })
}
getInterfaceList()
// 刷新
const refresh = () => {
  getInterfaceList(1)
}
//
const handleChange = (val) => {
  console.log(val)
}
// 过滤表格数据
const filterTableData = computed(() =>
  ctxData.commInterfaceData
    .filter((item) => !ctxData.interfaceName || item.name.toLowerCase().includes(ctxData.interfaceName.toLowerCase()))
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
)

const filterTableDataPage = computed(() => {
  ctxData.commListObj = {}
  ctxData.commInterfaceData
    .filter((item) => !ctxData.interfaceName || item.name.toLowerCase().includes(ctxData.interfaceName.toLowerCase()))
    .forEach((item) => {
      if (!ctxData.commListObj[item.type]) {
        ctxData.commListObj[item.type] = {
          opened: true,
          data: [],
        }
        ctxData.commListObj[item.type].data.push(item)
      } else {
        ctxData.commListObj[item.type].data.push(item)
      }
    })
  return ctxData.commListObj
})

// 添加接口
const addInterface = () => {
  //
  ctxData.tFlag = true
  ctxData.tTitle = '添加通讯接口'
}
// 编辑接口
const editInterface = (row) => {
  //
  ctxData.tFlag = true
  ctxData.tTitle = '编辑通讯接口'
  ctxData.interfaceForm.name = row.name
  ctxData.interfaceForm.type = row.type
  if (ctxData.serialList.includes(row.type)) {
    ctxData.interfaceForm.pName = row.param.name
    ctxData.interfaceForm.pBaudRate = row.param.baudRate
    ctxData.interfaceForm.pDataBits = row.param.dataBits
    ctxData.interfaceForm.pStopBits = row.param.stopBits
    ctxData.interfaceForm.pParity = row.param.parity
    ctxData.interfaceForm.timeout = row.param.timeout
    ctxData.interfaceForm.interval = row.param.interval
    ctxData.interfaceForm.ledModuleEnable = row.param.ledModuleEnable
    ctxData.interfaceForm.ledGPIO = row.param.ledGPIO
    ctxData.interfaceForm.ledOn = row.param.ledOn
    ctxData.interfaceForm.ledOff = row.param.ledOff
  }
  if (ctxData.tcpUdpList.includes(row.type)) {
    ctxData.interfaceForm.ip = row.param.ip
    ctxData.interfaceForm.port = row.param.port
    ctxData.interfaceForm.timeout = row.param.timeout
    ctxData.interfaceForm.interval = row.param.interval
  }
  if (ctxData.ioList.includes(row.type)) {
    ctxData.interfaceForm.ioName = row.param.name
  }
  if (row.type === 'httpTianGang') {
    // NBIOT天罡水表
    ctxData.interfaceForm.timeout = row.param.timeout
  }
  if (row.type === 's7') {
    // s7
    ctxData.interfaceForm.ip = row.param.ip
    ctxData.interfaceForm.port = row.param.port
    ctxData.interfaceForm.timeout = row.param.timeout
    ctxData.interfaceForm.connect = row.param.connect
  }
  if (row.type === 'sac009' || row.type === 'dr504') {
    ctxData.interfaceForm.port = row.param.port
    ctxData.interfaceForm.timeout = row.param.timeout
  }
  if (row.type === 'hdkj') {
    ctxData.interfaceForm.port = row.param.port
    ctxData.interfaceForm.reportInterval = row.param.reportInterval
    ctxData.interfaceForm.offlinePeriod = row.param.offlinePeriod
  }
  if (row.type === 'httpSmartNode') {
    ctxData.interfaceForm.ip = row.param.ip
    ctxData.interfaceForm.port = row.param.port
    ctxData.interfaceForm.appCode = row.param.appCode
    ctxData.interfaceForm.appKey = row.param.appKey
    ctxData.interfaceForm.timeout = row.param.timeout
    ctxData.interfaceForm.interval = row.param.interval
  }
}
const interfaceFormRef = ref()
// 提交接口表单
const submitInterfaceForm = () => {
  interfaceFormRef.value.validate((valid) => {
    if (valid) {
      let param = {}

      console.log('submitInterfaceForm -> ctxData.interfaceForm.type', ctxData.interfaceForm.type)
      if (ctxData.serialList.includes(ctxData.interfaceForm.type)) {
        param['name'] = ctxData.interfaceForm.pName
        param['baudRate'] = ctxData.interfaceForm.pBaudRate
        param['dataBits'] = ctxData.interfaceForm.pDataBits
        param['stopBits'] = ctxData.interfaceForm.pStopBits
        param['parity'] = ctxData.interfaceForm.pParity
        param['timeout'] = ctxData.interfaceForm.timeout
        param['interval'] = ctxData.interfaceForm.interval
        param['ledModuleEnable'] = ctxData.interfaceForm.ledModuleEnable
        param['ledGPIO'] = ctxData.interfaceForm.ledGPIO
        param['ledOn'] = ctxData.interfaceForm.ledOn
        param['ledOff'] = ctxData.interfaceForm.ledOff
      } else if (ctxData.interfaceForm.type === 'httpTianGang') {
        //NBIOT天罡水表
        param['timeout'] = ctxData.interfaceForm.timeout
      } else if (ctxData.interfaceForm.type === 's7') {
        param['ip'] = ctxData.interfaceForm.ip
        param['port'] = ctxData.interfaceForm.port
        param['timeout'] = ctxData.interfaceForm.timeout
        param['connect'] = ctxData.interfaceForm.connect
      } else if (ctxData.interfaceForm.type === 'sac009' || ctxData.interfaceForm.type === 'dr504') {
        param['port'] = ctxData.interfaceForm.port
        param['timeout'] = ctxData.interfaceForm.timeout
      } else if (ctxData.interfaceForm.type === 'hdkj') {
        param['port'] = ctxData.interfaceForm.port
        param['reportInterval'] = ctxData.interfaceForm.reportInterval
        param['offlinePeriod'] = ctxData.interfaceForm.offlinePeriod
      } else if (ctxData.interfaceForm.type === 'httpSmartNode') {
        param['ip'] = ctxData.interfaceForm.ip
        param['port'] = ctxData.interfaceForm.port
        param['appCode'] = ctxData.interfaceForm.appCode
        param['appKey'] = ctxData.interfaceForm.appKey
        param['timeout'] = ctxData.interfaceForm.timeout
        param['interval'] = ctxData.interfaceForm.interval
      } else {
        if (ctxData.tcpUdpList.includes(ctxData.interfaceForm.type)) {
          param['ip'] = ctxData.interfaceForm.ip
          param['port'] = ctxData.interfaceForm.port
          param['timeout'] = ctxData.interfaceForm.timeout
          if (ctxData.interfaceForm.type !== 's7') {
            param['interval'] = ctxData.interfaceForm.interval
          }
        }
        if (ctxData.ioList.includes(ctxData.interfaceForm.type)) {
          param['name'] = ctxData.interfaceForm.ioName
        }
      }
      console.log('submitInterfaceForm -> param', param)
      const pData = {
        token: users.token,
        data: {
          name: ctxData.interfaceForm.name,
          type: ctxData.interfaceForm.type,
          param,
        },
      }
      if (ctxData.tTitle.includes('添加')) {
        CommInterfaceApi.addCommInterface(pData).then((res) => {
          handleResult(res, getInterfaceList)
          cancelSubmit()
        })
      } else {
        CommInterfaceApi.editCommInterface(pData).then((res) => {
          handleResult(res, getInterfaceList)
          cancelSubmit()
        })
      }
    } else {
      return false
    }
  })
}
// 处理弹窗右上角关闭icon
const handleClose = (done) => {
  //
  cancelSubmit()
}
// 取消提交表单
const cancelSubmit = () => {
  //
  ctxData.tFlag = false
  interfaceFormRef.value.resetFields()
  initInterfaceForm()
}
//删除接口
const deleteInterface = (row) => {
  ElMessageBox.confirm('确定要删除这个通讯接口吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          name: row.name,
        },
      }
      CommInterfaceApi.deleteCommInterface(pData).then((res) => {
        handleResult(res, getInterfaceList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}
//初始化接口表单
const initInterfaceForm = () => {
  ctxData.interfaceForm = {
    name: '',
    type: '',
    //串口
    pName: '', //串口名称，比如"/dev/ttyS1"或者"COM1"
    pBaudRate: '', //波特率
    pDataBits: '', //数据位，5，6，7，8
    pStopBits: '', //停止位，1，1.5，2
    pParity: '', //校验位，N（无校验），O（奇校验），E（偶校验）
    timeout: '', //超时时间（毫秒）
    interval: '', //间隔时间（毫秒）
    ledModuleEnable: false, //启用led模块
    ledGPIO: '', //LED管脚编号
    ledOn: '', //LED点亮
    ledOff: '', //LED熄灭
    //tcp/udp
    pIp: '', //本地（远端）IP地址
    pPort: '', //本地（远端）端口
    //io
    ioName: '', //IO端口名称
    connect: '',
    reportInterval: '',
    offlinePeriod: '',
  }
}
//显示参数切换name
const showParam = (type) => {
  ctxData.paramName.name = ctxData.ioList.includes(type) ? 'IO端口名称' : '串口名称'
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
.main-container .main {
  min-width: 748px;
  overflow: auto;
}
.content-up {
  overflow: auto;
  min-width: 700px;
  padding: 0 8px;
}
.cu-main {
  position: relative;
  height: 100%;
  width: 100%;
  box-sizing: border-box;
}
.tName {
  display: flex;
  line-height: 16px;
  font-size: 18px;
  border-left: 4px solid #3054eb;
  padding-left: 20px;
  align-items: center;
}
.btn-icon-arrow {
  margin-left: 8px;
  cursor: pointer;
}
.btn-icon-arrow .el-icon {
  margin-right: 0;
}
.arrow-opened {
  transition: all 0.5s;
}
.arrow-closed {
  transition: all 0.5s;
  transform: rotate(90deg);
}
.cu-items {
  position: relative;
  left: -8px;
  padding-left: 8px;
  right: -8px;
  top: 0px;
  bottom: 0;
  box-sizing: border-box;
  display: flex;
  flex-wrap: wrap;
  overflow: auto;
}
.cui-opened {
  transition: all 1s ease;
}
.cui-closed {
  overflow-y: hidden;
  transition: all 1s ease;
}
.cu-item {
  position: relative;
  width: 20%;
  min-width: 322px;
  box-sizing: border-box;
  margin-bottom: 18px;
}
.cui-content {
  position: relative;
  height: 100%;
  margin: 8px 8px 16px 8px;
}
.cui-header {
  display: flex;
  align-items: center;
  width: 190px;
  white-space: nowrap;
  overflow: auto;
}
.head-tag {
  float: right;
  border: 0;
  margin-left: 4px;
  height: 22px;
  padding: 0 2px;
}
.card-content {
  position: relative;
  font-size: 14px;
  .cc-body {
    position: relative;
    top: 10px;
    left: 18px;
    right: 18px;
    bottom: 8px;
    overflow: auto;
  }
  .ccb-item {
    float: left;
    height: 36px;
    width: calc(100% - 44px);
    padding: 0 4px;
    line-height: 36px;
    display: flex;
    justify-content: space-between;
  }
}
:deep(.card-header) {
  display: flex;
  justify-content: space-between;
}
:deep(.el-card) {
  height: 100%;
}
:deep(.el-card__header) {
  padding: 12px 20px;
}
:deep(.el-card__body) {
  height: calc(100% - 48px);
  position: relative;
  box-sizing: border-box;
  padding: 0;
}
:deep(.el-collapse-item__header.is-active) {
  background-color: #f5f5f5 !important;
}
:deep(.el-collapse-item__content) {
  padding-bottom: 15px;
  background-color: #f5f5f5 !important;
}
</style>
<style lang="scss">
.el-popover.el-popper {
  padding: 20px;
}
</style>
