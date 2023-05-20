<template>
  <div class="main-container">
    <div class="main">
      <div class="title" style="justify-content: space-between">
        <el-input style="width: 200px" placeholder="请输入模块名称" clearable v-model="ctxData.moduleName">
          <template #prefix>
            <el-icon class="el-input__icon"><search /></el-icon>
          </template>
        </el-input>
        <div>
          <el-button v-show="ctxData.showAddBtn" type="primary" bg class="right-btn" @click="addModule()">
            <el-icon class="btn-icon">
              <Icon name="local-add" size="14px" color="#ffffff" />
            </el-icon>
            添加
          </el-button>
          <el-button style="color: #fff" color="#2EA554" class="right-btn" @click="refresh()">
            <el-icon class="btn-icon">
              <Icon name="local-refresh" size="14px" color="#ffffff" />
            </el-icon>
            刷新
          </el-button>
        </div>
      </div>
      <div class="content" ref="contentRef">
        <el-table
          :data="filterTableData"
          :cell-style="ctxData.cellStyle"
          :header-cell-style="ctxData.headerCellStyle"
          :max-height="ctxData.tableMaxHeight"
          style="width: 100%"
          stripe
          @row-dblclick="editModule"
        >
          <el-table-column prop="name" label="模块名称" width="auto" min-width="120" align="center"> </el-table-column>
          <el-table-column prop="model" label="模块型号" width="auto" min-width="120" align="center"> </el-table-column>
          <el-table-column label="sim卡号" width="auto" min-width="180" align="center">
            <template #default="scope">
              <el-tooltip
                v-if="scope.row.runParam.iccid !== ''"
                class="box-item"
                effect="dark"
                content="单击显示SIM卡条形码"
                placement="top-start"
              >
                <el-tag style="cursor: pointer" @click="showSimNum(scope.row.runParam.iccid)">{{
                  scope.row.runParam.iccid
                }}</el-tag>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column label="imei" width="auto" min-width="150" align="center">
            <template #default="scope">
              {{ scope.row.runParam.imei }}
            </template>
          </el-table-column>
          <el-table-column label="信号强度" width="auto" min-width="80" align="center">
            <template #default="scope">
              {{ scope.row.runParam.csq }}
            </template>
          </el-table-column>
          <el-table-column label="流量（KByte）" width="auto" min-width="150" align="center">
            <template #default="scope">
              {{ scope.row.runParam.flow }}
            </template>
          </el-table-column>
          <el-table-column label="基站定位Lac" width="auto" min-width="110" align="center">
            <template #default="scope">
              {{ scope.row.runParam.lac }}
            </template>
          </el-table-column>
          <el-table-column label="基站定位Ci" width="auto" min-width="100" align="center">
            <template #default="scope">
              {{ scope.row.runParam.ci }}
            </template></el-table-column
          >
          <el-table-column label="sim卡插入状态" width="auto" min-width="120" align="center">
            <template #default="scope">
              <el-tag :type="scope.row.runParam.simInsert ? 'success' : 'danger'">{{
                scope.row.runParam.simInsert ? '是' : '否'
              }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="网络注册状态" width="auto" min-width="120" align="center">
            <template #default="scope">
              <el-tag :type="scope.row.runParam.netRegister ? 'success' : 'danger'">
                {{ scope.row.runParam.netRegister ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="auto" min-width="150" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="editModule(scope.row)" text type="primary">编辑</el-button>
              <el-button @click="deleteModule(scope.row)" text type="danger">删除</el-button>
            </template>
          </el-table-column>
          <template #empty>
            <div>无数据</div>
          </template>
        </el-table>
        <div class="mobile-tips">
          <el-tag class="ml-2" type="danger">注：移动网络配置完成后，必须重启网关，才能生效！</el-tag>
        </div>
      </div>
    </div>
    <el-dialog
      v-model="ctxData.mFlag"
      :title="ctxData.mTitle"
      width="900px"
      :before-close="handleClose"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.moduleForm"
          :rules="ctxData.moduleRules"
          ref="moduleFormRef"
          status-icon
          label-position="right"
          label-width="120px"
          inline="true"
        >
          <el-form-item label="模块名称" prop="name">
            <el-input
              :disabled="ctxData.mTitle.includes('编辑')"
              type="text"
              style="width: 270px"
              v-model="ctxData.moduleForm.name"
              autocomplete="off"
              placeholder="请输入模块名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="模块型号" prop="model">
            <el-input
              type="text"
              style="width: 270px"
              v-model="ctxData.moduleForm.model"
              autocomplete="off"
              placeholder="请输入模块型号"
            >
            </el-input>
          </el-form-item>
          <div class="form-title"><div class="tName">配置参数</div></div>
          <el-form-item label="启用流量报警" prop="flowAlarm">
            <el-switch
              style="width: 270px"
              v-model="ctxData.moduleForm.flowAlarm"
              inline-prompt
              active-text="是"
              inactive-text="否"
            />
          </el-form-item>
          <el-form-item label="流量报警值" :prop="ctxData.moduleForm.flowAlarm ? 'flowAlarmValue' : ''">
            <el-input
              :disabled="!ctxData.moduleForm.flowAlarm"
              type="text"
              style="width: 270px"
              v-model="ctxData.moduleForm.flowAlarmValue"
              autocomplete="off"
              placeholder="请输入流量报警值"
            >
            </el-input>
          </el-form-item>
          <div class="form-title"><div class="tName">通信参数</div></div>
          <el-form-item label="串口号" prop="serialNumber">
            <el-input
              type="text"
              style="width: 270px"
              v-model="ctxData.moduleForm.serialNumber"
              autocomplete="off"
              placeholder="请输入串口名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="波特率" prop="baudRate">
            <el-select v-model="ctxData.moduleForm.baudRate" style="width: 270px" placeholder="请选择波特率">
              <el-option
                v-for="item in ctxData.baudRateOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="数据位" prop="dataBits">
            <el-select v-model="ctxData.moduleForm.dataBits" style="width: 270px" placeholder="请选择数据位">
              <el-option
                v-for="item in ctxData.dataBitsOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="停止位" prop="stopBits">
            <el-select v-model="ctxData.moduleForm.stopBits" style="width: 270px" placeholder="请选择停止位">
              <el-option
                v-for="item in ctxData.stopBitsOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="校验位" prop="parity">
            <el-select v-model="ctxData.moduleForm.parity" style="width: 270px" placeholder="请选择校验位">
              <el-option
                v-for="item in ctxData.parityOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="超时时间" prop="timeout">
            <el-input
              type="text"
              style="width: 270px"
              v-model="ctxData.moduleForm.timeout"
              autocomplete="off"
              placeholder="请输入超时时间"
            >
              <template #append>单位毫秒</template>
            </el-input>
          </el-form-item>
          <el-form-item label="间隔时间" prop="interval">
            <el-input
              type="text"
              style="width: 270px"
              v-model="ctxData.moduleForm.interval"
              autocomplete="off"
              placeholder="请输入间隔时间"
            >
              <template #append>单位毫秒</template>
            </el-input>
          </el-form-item>
          <el-form-item label="轮询时间" prop="polling">
            <el-input
              type="text"
              style="width: 270px"
              v-model="ctxData.moduleForm.polling"
              autocomplete="off"
              placeholder="请输入轮询时间"
            >
              <template #append>单位秒</template>
            </el-input>
          </el-form-item>
          <div class="form-title"><div class="tName">保活参数</div></div>
          <el-form-item label="保活检测主IP" prop="ipMaster">
            <el-input
              type="text"
              style="width: 270px"
              v-model="ctxData.moduleForm.ipMaster"
              autocomplete="off"
              :placeholder="'请输入保活检测主IP'"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="保活检测备用IP" prop="ipSlave">
            <el-input
              type="text"
              style="width: 270px"
              v-model="ctxData.moduleForm.ipSlave"
              autocomplete="off"
              :placeholder="'请输入保活检测备用IP'"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="保活检测周期" prop="pollPeriod">
            <el-input
              type="text"
              style="width: 270px"
              v-model.number="ctxData.moduleForm.pollPeriod"
              autocomplete="off"
              placeholder="请输入保活检测周期"
            >
              <template #append>单位秒</template>
            </el-input>
          </el-form-item>
          <el-form-item label="离线判断次数" prop="offlineCnt">
            <el-input
              type="text"
              style="width: 270px"
              v-model.number="ctxData.moduleForm.offlineCnt"
              autocomplete="off"
              placeholder="请输入离线判断次数"
            >
              <template #append>单位次</template>
            </el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelSubmit()">取消</el-button>
          <el-button type="primary" @click="submitModuleForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
    <!--sim卡sn-->
    <el-dialog v-model="ctxData.sFlag" title="sim卡号条形码" width="30%">
      <div style="display: flex; width: 100%; justify-content: center; text-align: center">
        <vue3-barcode :value="ctxData.barcodeVale" :height="60" />
      </div>
    </el-dialog>
  </div>
</template>
<script setup>
import variables from 'styles/variables.module.scss'
import ModuleApi from 'api/module.js'
import { Search } from '@element-plus/icons-vue'
import { userStore } from 'stores/user'
import Vue3Barcode from 'vue3-barcode'
const users = userStore()

const refExpIP =
  /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/
const refExpYM =
  /^(?=^.{3,255}$)(http(s)?:\/\/)?(www\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:\d+)*(\/\w+\.\w+)*$/
const validateIP = (rule, value, callback) => {
  console.log('validateIP')
  if (refExpIP.test(value) || refExpYM.test(value)) {
    callback()
  } else {
    callback(new Error('IP格式错误！'))
  }
}
const validateNumber = (rule, value, callback) => {
  if (isNaN(Number(value))) {
    callback(new Error('只能输入数字！'))
  } else {
    callback()
  }
}
const contentRef = ref(null)
const ctxData = reactive({
  moduleName: '',
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
  moduleTableData: [],
  showAddBtn: true,
  mFlag: false,
  mTitle: '添加移动网络',
  moduleForm: {
    name: '',
    model: '',
    flowAlarm: false, //流量报警
    flowAlarmValue: '', //流量报警值
    serialNumber: '', //串口号
    baudRate: '', //波特率
    dataBits: '', //数据位
    stopBits: '', //停止位
    parity: '', //校验位
    timeout: '', //超时时间（毫秒）
    interval: '', //间隔时间（毫秒）
    polling: '', //轮询时间（秒）
    ipMaster: '', //保活检测主IP
    ipSlave: '', //保活检测备用IP
    pollPeriod: '', //保活检测周期
    offlineCnt: '', //离线判断次数
  },
  moduleRules: {
    name: [
      {
        required: true,
        message: '模块名称不能为空',
        trigger: 'blur',
      },
    ],
    model: [
      {
        required: true,
        message: '模块型号不能为空',
        trigger: 'blur',
      },
    ],
    serialNumber: [
      {
        required: true,
        message: '串口号不能为空',
        trigger: 'blur',
      },
    ],
    baudRate: [
      {
        required: true,
        message: '波特率不能为空',
        trigger: 'blur',
      },
    ],
    dataBits: [
      {
        required: true,
        message: '数据位不能为空',
        trigger: 'blur',
      },
    ],
    stopBits: [
      {
        required: true,
        message: '停止位不能为空',
        trigger: 'blur',
      },
    ],
    parity: [
      {
        required: true,
        message: '校验位不能为空',
        trigger: 'blur',
      },
    ],
    flowAlarmValue: [
      {
        required: true,
        message: '流量报警值不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateNumber,
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
        trigger: 'blur',
        validator: validateNumber,
      },
    ],
    polling: [
      {
        required: true,
        message: '轮询时间不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateNumber,
      },
    ],
    ipMaster: [
      {
        required: true,
        validator: validateIP,
        trigger: 'blur',
      },
    ],
    ipSlave: [
      {
        validator: validateIP,
        trigger: 'blur',
      },
    ],
    pollPeriod: [
      {
        required: true,
        message: '保活检测周期不能为空',
        trigger: 'blur',
      },
      {
        type: 'number',
        message: '保活检测周期只能输入数字',
      },
    ],
    offlineCnt: [
      {
        required: true,
        message: '离线判断次数不能为空',
        trigger: 'blur',
      },
      {
        type: 'number',
        message: '离线判断次数只能输入数字',
      },
    ],
  },
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
  sFlag: false,
  barcodeVale: 0,
})
const getModuleList = (flag) => {
  const pData = {
    token: users.token,
    data: {},
  }
  ModuleApi.getModuleList(pData).then(async (res) => {
    console.log('getModuleList -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.moduleTableData = res.data
      ctxData.showAddBtn = res.data.length === 0
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
getModuleList()
// 刷新
const refresh = () => {
  getModuleList(1)
}
//
const showSimNum = (sim) => {
  ctxData.sFlag = true
  ctxData.barcodeVale = sim
}

// 过滤表格数据
const filterTableData = computed(() =>
  ctxData.moduleTableData
    .filter((item) => !ctxData.moduleName || item.name.toLowerCase().includes(ctxData.moduleName.toLowerCase()))
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
)
const filterTableDataPage = computed(() =>
  ctxData.moduleTableData.filter(
    (item) => !ctxData.moduleName || item.name.toLowerCase().includes(ctxData.moduleName.toLowerCase())
  )
)
// 添加以移动网络
const addModule = () => {
  ctxData.mFlag = true
}
// 编辑以移动网络
const editModule = (row) => {
  ctxData.mFlag = true
  ctxData.mTitle = '编辑移动网络'
  ctxData.moduleForm.name = row.name
  ctxData.moduleForm.model = row.model
  ctxData.moduleForm.flowAlarm = row.configParam.flowAlarm
  ctxData.moduleForm.flowAlarmValue = row.configParam.flowAlarmValue
  ctxData.moduleForm.serialNumber = row.commParam.name
  ctxData.moduleForm.baudRate = row.commParam.baudRate
  ctxData.moduleForm.dataBits = row.commParam.dataBits
  ctxData.moduleForm.stopBits = row.commParam.stopBits
  ctxData.moduleForm.parity = row.commParam.parity
  ctxData.moduleForm.timeout = row.commParam.timeout
  ctxData.moduleForm.interval = row.commParam.interval
  ctxData.moduleForm.polling = row.commParam.pollPeriod

  ctxData.moduleForm.ipMaster = row.keepAliveParam.ipMaster
  ctxData.moduleForm.ipSlave = row.keepAliveParam.ipSlave
  ctxData.moduleForm.pollPeriod = row.keepAliveParam.pollPeriod
  ctxData.moduleForm.offlineCnt = row.keepAliveParam.offlineCnt
}
// 提交表单
const moduleFormRef = ref(null)
const submitModuleForm = () => {
  moduleFormRef.value.validate((valid) => {
    const data = {
      name: ctxData.moduleForm.name,
      model: ctxData.moduleForm.model,
      commParam: {},
      configParam: {},
    }
    const commParam = {
      name: ctxData.moduleForm.serialNumber,
      baudRate: ctxData.moduleForm.baudRate,
      dataBits: ctxData.moduleForm.dataBits,
      stopBits: ctxData.moduleForm.stopBits,
      parity: ctxData.moduleForm.parity,
      timeout: ctxData.moduleForm.timeout,
      interval: ctxData.moduleForm.interval,
      pollPeriod: ctxData.moduleForm.polling,
    }
    const configParam = {
      flowAlarm: ctxData.moduleForm.flowAlarm,
      flowAlarmValue: ctxData.moduleForm.flowAlarmValue,
    }
    const keepAliveParam = {
      ipMaster: ctxData.moduleForm.ipMaster,
      ipSlave: ctxData.moduleForm.ipSlave,
      pollPeriod: ctxData.moduleForm.pollPeriod,
      offlineCnt: ctxData.moduleForm.offlineCnt,
    }
    data.commParam = commParam
    data.configParam = configParam
    data.keepAliveParam = keepAliveParam
    if (valid) {
      const pData = {
        token: users.token,
        data: data,
      }
      if (ctxData.mTitle.includes('添加')) {
        ModuleApi.addModule(pData).then((res) => {
          handleResult(res, getModuleList)
          cancelSubmit()
        })
      } else {
        ModuleApi.editModule(pData).then((res) => {
          handleResult(res, getModuleList)
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
const deleteModule = (row) => {
  ElMessageBox.confirm('确定要删除这个移动网络吗?', '警告', {
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
      ModuleApi.deleteModule(pData).then((res) => {
        handleResult(res, getModuleList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}
// 取消提交表单
const cancelSubmit = () => {
  //
  ctxData.mFlag = false
  moduleFormRef.value.resetFields()
  initmoduleForm()
}
const initmoduleForm = () => {
  ctxData.moduleForm = {
    name: '',
    model: '',
    flowAlarm: false, //流量报警
    flowAlarmValue: '', //流量报警值
    serialNumber: '', //串口号
    baudRate: '', //波特率
    dataBits: '', //数据位
    stopBits: '', //停止位
    parity: '', //校验位
    timeout: '', //超时时间（毫秒）
    interval: '', //间隔时间（毫秒）
    polling: '', //轮询时间（秒）
    ipMaster: '', //保活检测主IP
    ipSlave: '', //保活检测备用IP
    pollPeriod: '', //保活检测周期
    offlineCnt: '', //离线判断次数
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
.mobile-tips {
  position: absolute;
  bottom: 20px;
}
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
