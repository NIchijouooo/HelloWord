<template>
  <div class="main-container">
    <div class="main">
      <div class="title" style="justify-content: space-between">
        <div class="tName">属性列表</div>
        <div style="display: flex; align-items: center">
          <el-input
            style="width: 200px; margin-right: 12px; height: 32px"
            placeholder="请输入属性名称或者标签"
            clearable
            v-model="ctxData.propertyInfo"
          >
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" bg class="right-btn" @click="addProperty()">
            <el-icon class="btn-icon">
              <Icon name="local-add" size="14px" color="#ffffff" />
            </el-icon>
            添加
          </el-button>
          <el-button type="danger" bg class="right-btn" @click="deleteProperty()">
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
      <div class="content" ref="contentRef">
        <el-table
          :data="filterTableData"
          :cell-style="ctxData.cellStyle"
          :header-cell-style="ctxData.headerCellStyle"
          :max-height="ctxData.tableMaxHeight"
          style="width: 100%"
          stripe
          @selection-change="handleSelectionChange"
          @row-dblclick="editProperty"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column type="expand">
            <template #default="scope">
              <div class="param-content">
                <div class="pc-title">
                  <div class="pct-info">
                    <b> {{ scope.row.name }} </b>
                    参数详情
                  </div>
                </div>
                <div class="pc-content">
                  <div class="param-item" v-for="(item, key, index) of scope.row.params" :key="index">
                    <div class="param-value">{{ ctxData.paramName[key] }}：</div>
                    <div class="param-name">{{ typeof item === 'boolean' ? (item ? '是' : '否') : item }}</div>
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column sortable prop="name" label="属性名称" width="auto" min-width="180" align="center"> </el-table-column>
          <el-table-column sortable prop="label" label="属性标签" width="auto" min-width="180" align="center"> </el-table-column>
          <el-table-column sortable prop="type" label="数据类型" width="auto" min-width="100" align="center">
            <template #default="scope">
              {{ ctxData.typeNames[['t' + scope.row.type]] }}
            </template>
          </el-table-column>
          <el-table-column sortable prop="decimals" label="小数位" width="auto" min-width="80" align="center"> </el-table-column>
          <el-table-column sortable prop="unit" label="单位" width="auto" min-width="80" align="center"> </el-table-column>
          <el-table-column label="操作" width="auto" min-width="120" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="editProperty(scope.row)" text type="primary">编辑</el-button>
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
    <!-- 添加/编辑属性 -->
    <el-dialog v-model="ctxData.pFlag" :title="ctxData.pTitle" width="600px">
      <div class="dialog-content">
        <el-form
          :model="ctxData.propertyForm"
          :rules="ctxData.propertyRules"
          ref="propertyFormRef"
          status-icon
          label-position="right"
          label-width="120px"
        >
          <el-form-item label="采集接口名称" prop="collName">
            <el-select
              v-model="ctxData.propertyForm.collName"
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
          <el-form-item label="设备名称" prop="deviceName">
            <el-select
              v-model="ctxData.propertyForm.deviceName"
              style="width: 100%"
              placeholder="请选择设备名称"
              @change="selectDeviceName"
              clearable
              filterable
            >
              <el-option v-for="item in ctxData.nodeOptions" :key="item.name" :label="item.label" :value="item.name">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="设备属性名称" prop="propertyName">
            <el-select
              v-model="ctxData.propertyForm.propertyName"
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
          <el-form-item label="属性名称" prop="name">
            <el-input type="text" v-model="ctxData.propertyForm.name" autocomplete="off" placeholder="请输入属性名称">
            </el-input>
          </el-form-item>
          <el-form-item label="属性标签" prop="label">
            <el-input type="text" v-model="ctxData.propertyForm.label" autocomplete="off" placeholder="请输入属性标签">
            </el-input>
          </el-form-item>

          <el-form-item label="数据类型" prop="type">
            <el-select
              v-model="ctxData.propertyForm.type"
              style="width: 100%"
              placeholder="请选择数据类型"
              clearable
              filterable
            >
              <el-option v-for="item in ctxData.typeOptions" :key="item.value" :label="item.label" :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type === 2" label="小数位数" prop="decimals">
            <el-input
              type="text"
              v-model="ctxData.propertyForm.decimals"
              autocomplete="off"
              placeholder="请输入小数位数"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="单位" prop="unit">
            <el-input type="text" v-model="ctxData.propertyForm.unit" autocomplete="off" placeholder="请输入单位">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type !== 3" label="范围报警" prop="minMaxAlarm">
            <el-switch v-model="ctxData.propertyForm.minMaxAlarm" inline-prompt active-text="是" inactive-text="否" />
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.minMaxAlarm"
            label="最小值"
            prop="min"
          >
            <el-input type="text" v-model="ctxData.propertyForm.min" autocomplete="off" placeholder="请输入最小值">
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.minMaxAlarm"
            label="最大值"
            prop="max"
          >
            <el-input type="text" v-model="ctxData.propertyForm.max" autocomplete="off" placeholder="请输入最大值">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type !== 3" label="步长报警" prop="stepAlarm">
            <el-switch v-model="ctxData.propertyForm.stepAlarm" inline-prompt active-text="是" inactive-text="否" />
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.stepAlarm"
            label="步长"
            prop="step"
          >
            <el-input type="text" v-model="ctxData.propertyForm.step" autocomplete="off" placeholder="请输入步长">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type === 3" label="字符串长度报警" prop="dataLengthAlarm">
            <el-switch
              v-model="ctxData.propertyForm.dataLengthAlarm"
              inline-prompt
              active-text="是"
              inactive-text="否"
            />
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type === 3 && ctxData.propertyForm.dataLengthAlarm"
            label="字符串长度"
            prop="dataLength"
          >
            <el-input
              type="text"
              v-model="ctxData.propertyForm.dataLength"
              autocomplete="off"
              placeholder="请输入字符串长度"
            >
            </el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelSubmit()">取消</el-button>
          <el-button type="primary" @click="submitForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { Search } from '@element-plus/icons-vue'
import variables from 'styles/variables.module.scss'
import InterfaceApi from 'api/interface.js'
import DeviceModelApi from 'api/deviceModel.js'
import VirtualDeviceApi from 'api/virtualDevice.js'
import { userStore } from 'stores/user'
const users = userStore()
const regCnt = /^[0-9]*[1-9][0-9]*$/
const validateNum = (rule, value, callback) => {
  console.log('validateNum  value = ', value)
  console.log('validateNum  value.typeof = ', typeof value)
  if (value !== '0' && !regCnt.test(parseInt(value))) {
    callback(new Error('只能输入自然数！'))
  } else {
    callback()
  }
}
const ctxData = reactive({
  propertyInfo: '',
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
  pFlag: false,
  pTitle: '添加属性',
  paramName: {
    collName: '采集接口名称',
    deviceName: '设备名称',
    propertyName: '设备属性名称',
    min: '最小值',
    max: '最大值',
    minMaxAlarm: '范围报警',
    step: '步长',
    stepAlarm: '步长报警',
    dataLength: '字符串长度',
    dataLengthAlarm: '字符串长度报警',
  },
  propertyForm: {
    name: '',
    label: '',
    type: 0,
    decimals: '0',
    unit: '',
    collName: '',
    deviceName: '',
    propertyName: '',
    min: '', // 属性最小值，只有uint32，int32，double有效
    max: '', // 属性最大值，只有uint32，int32，double有效
    minMaxAlarm: false, // 范围报警，只有uint32，int32，double有效
    step: '', // 步长，只有uint32，int32，double有效
    stepAlarm: false, // 步长报警，只有uint32，int32，double有效
    dataLength: '', // 字符串长度，只有string有效
    dataLengthAlarm: false, // 字符串长度报警，只有string有效
  },
  propertyRules: {
    name: [
      {
        required: true,
        message: '属性名称不能为空！',
        trigger: 'blur',
      },
    ],
    label: [
      {
        required: true,
        message: '属性标签不能为空！',
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
    deviceName: [
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
    type: [
      {
        required: true,
        message: '数据类型不能为空！',
        trigger: 'blur',
      },
    ],
    decimals: [
      {
        trigger: 'blur',
        validator: validateNum,
      },
    ],
  },
  typeOptions: [
    //数据类型
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
  collOptions: [], //采集接口选项
  nodeOptions: [], //设备选项
  propertyOptions: [], //设备属性选项
  nameTotsl: {}, //设备名称对应设备模板名称
  selectedNodes: [],
})

const contentRef = ref(null)
const getPropertyList = (flag) => {
  const pData = {
    token: users.token,
    data: {},
  }
  VirtualDeviceApi.getPropertyList(pData).then(async (res) => {
    console.log('getPropertyList -> res ', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.propertyTableData = res.data
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
getPropertyList()
//获取采集接口
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
  getCollDevices(val)
  ctxData.propertyOptions = []
}
// 获取采集接口下的设备列表
const getCollDevices = (collName) => {
  const pData = {
    token: users.token,
    data: {
      name: collName,
    },
  }
  InterfaceApi.getCollDevices(pData).then((res) => {
    console.log('getCollDevices res => ', res)
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
  })
}
// 选择采集接口触发事件
const selectDeviceName = (nodeName) => {
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
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
const refresh = () => {
  getPropertyList(1)
}

const addProperty = () => {
  ctxData.pFlag = true
  ctxData.pTitle = '添加属性'
  initPropertyForm()
}

const editProperty = (row) => {
  ctxData.pFlag = true
  ctxData.pTitle = '编辑属性'
  ctxData.propertyForm.name = row.name
  ctxData.propertyForm.label = row.label
  ctxData.propertyForm.type = row.type
  ctxData.propertyForm.decimals = row.decimals
  ctxData.propertyForm.unit = row.unit
  ctxData.propertyForm.collName = row.params.collName
  ctxData.propertyForm.deviceName = row.params.deviceName
  ctxData.propertyForm.propertyName = row.params.propertyName
  ctxData.propertyForm.minMaxAlarm = row.alarmParams.minMaxAlarm
  ctxData.propertyForm.min = row.alarmParams.min
  ctxData.propertyForm.max = row.alarmParams.max
  ctxData.propertyForm.step = row.alarmParams.step
  ctxData.propertyForm.stepAlarm = row.alarmParams.stepAlarm
  ctxData.propertyForm.dataLengthAlarm = row.alarmParams.dataLengthAlarm
  ctxData.propertyForm.dataLength = row.alarmParams.dataLength
}

const propertyFormRef = ref(null)
//保存
const submitForm = () => {
  propertyFormRef.value.validate((valid) => {
    if (valid) {
      const postForm = {
        name: ctxData.propertyForm.name,
        label: ctxData.propertyForm.label,
        type: ctxData.propertyForm.type,
        decimals: ctxData.propertyForm.decimals,
        unit: ctxData.propertyForm.unit,
        params: {
          collName: ctxData.propertyForm.collName,
          deviceName: ctxData.propertyForm.deviceName,
          propertyName: ctxData.propertyForm.propertyName,
        },
        alarmParams: {
          min: '',
          max: '',
          minMaxAlarm: false,
          step: '',
          stepAlarm: false,
          dataLength: '',
          dataLengthAlarm: false,
        },
      }

      if (postForm.type !== 3) {
        postForm.alarmParams['min'] = ctxData.propertyForm.min
        postForm.alarmParams['max'] = ctxData.propertyForm.max
        postForm.alarmParams['minMaxAlarm'] = ctxData.propertyForm.minMaxAlarm
        postForm.alarmParams['step'] = ctxData.propertyForm.step
        postForm.alarmParams['stepAlarm'] = ctxData.propertyForm.stepAlarm
      } else {
        postForm.alarmParams['dataLength'] = ctxData.propertyForm.dataLength
        postForm.alarmParams['dataLengthAlarm'] = ctxData.propertyForm.dataLengthAlarm
      }
      const pData = {
        token: users.token,
        data: postForm,
      }
      if (ctxData.pTitle.includes('添加')) {
        VirtualDeviceApi.addProperty(pData).then((res) => {
          handleResult(res, getPropertyList)
          cancelSubmit()
        })
      } else {
        VirtualDeviceApi.editProperty(pData).then((res) => {
          handleResult(res, getPropertyList)
          cancelSubmit()
        })
      }
    } else {
      return false
    }
  })
}
//取消保存
const cancelSubmit = () => {
  ctxData.pFlag = false
  initPropertyForm()
}
//删除属性
const deleteProperty = () => {
  let dList = []
  if (ctxData.selectedNodes.length === 0) {
    ElMessage.info('请至少选择一个属性！')
    return
  } else {
    ctxData.selectedNodes.forEach((item) => {
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
          names: dList,
        },
      }
      VirtualDeviceApi.deletePorperties(pData).then((res) => {
        handleResult(res, getPropertyList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}

const initPropertyForm = () => {
  ctxData.propertyForm = {
    name: '',
    label: '',
    type: 0,
    decimals: '0',
    unit: '',
    collName: '',
    deviceName: '',
    propertyName: '',
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
