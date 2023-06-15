<template>
  <div class="main">
    <div class="title" style="justify-content: space-between">
      <div class="title-left">
        <el-button type="primary" plain @click="toVirtualDevice()" style="margin-right: 20px">
          <el-icon class="el-input__icon"><back /></el-icon>
          返回上报模型
        </el-button>
      </div>
    </div>
    <div class="title" style="top: 60px; height: 76px; padding: 20px 0; justify-content: space-between">
      <div class="tName">{{ props.curVirtualDevice.label }}</div>
      <div style="display: flex; align-items: center">
        <el-input style="width: 200px" placeholder="请输入属性名称或者标签" clearable v-model="ctxData.PropertyInfo">
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
    <div class="content" ref="contentRef" style="top: 136px">
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
        <el-table-column sortable prop="name" label="属性名称" width="auto" min-width="180" align="center"> </el-table-column>
        <el-table-column sortable prop="label" label="属性标签" width="auto" min-width="180" align="center"> </el-table-column>
        <el-table-column sortable label="数据类型" width="auto" min-width="100" align="center">
          <template #default="scope">
            {{ ctxData.typeNames['t' + scope.row.type] }}
          </template>
        </el-table-column>
        <el-table-column sortable label="小数位数" width="auto" min-width="80" align="center">
          <template #default="scope">
            {{ scope.row.decimals === undefined || scope.row.decimals === '' ? 0 : scope.row.decimals }}
          </template>
        </el-table-column>
        <el-table-column sortable prop="unit" label="单位" width="auto" min-width="80" align="center"> </el-table-column>
        <el-table-column sortable label="映射参数详情" width="auto" min-width="300" align="center">
          <template #default="scope">
            <el-popover
              trigger="hover"
              :show-after="500"
              :auto-close="500"
              effect="light"
              width="auto"
              placement="left"
            >
              <template #default>
                <div class="param-content">
                  <div class="pc-title">
                    <div class="pct-info">
                      <b> {{ scope.row.name }} </b>
                      {{ '参数详情' }}
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
              <template #reference>
                <el-tag size="large">{{ scope.row.params }}</el-tag>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="报警参数详情" width="auto" min-width="300" align="center">
          <template #default="scope">
            <el-popover
              trigger="hover"
              :show-after="500"
              :auto-close="500"
              effect="light"
              width="auto"
              placement="left"
            >
              <template #default>
                <div class="param-content">
                  <div class="pc-title">
                    <div class="pct-info">
                      <b> {{ scope.row.name }} </b>
                      {{ '参数详情' }}
                    </div>
                  </div>
                  <div class="pc-content">
                    <div class="param-item" v-for="(item, key, index) of scope.row.alarmParams" :key="index">
                      <div class="param-value">{{ ctxData.paramName[key] }}：</div>
                      <div class="param-name">{{ typeof item === 'boolean' ? (item ? '是' : '否') : item }}</div>
                    </div>
                  </div>
                </div>
              </template>
              <template #reference>
                <el-tag size="large">{{ scope.row.alarmParams }}</el-tag>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="auto" min-width="150" align="center" fixed="right">
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
    <el-dialog
      v-model="ctxData.pFlag"
      :title="ctxData.pTitle"
      width="800px"
      :before-close="handleCloseProperty"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.propertyForm"
          :rules="ctxData.propertyRules"
          ref="propertyFormRef"
          status-icon
          label-position="right"
          label-width="120px"
          :inline="true"
        >
          <el-form-item label="属性名称" prop="name">
            <el-input
              style="width: 220px"
              :disabled="ctxData.pTitle.includes('编辑')"
              :readonly="true"
              type="text"
              v-model="ctxData.propertyForm.name"
              autocomplete="off"
              placeholder="请选择属性名称"
            >
              <template #append>
                <el-button @click="showModelProperty()">同步</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="属性标签" prop="label">
            <el-input
              style="width: 220px"
              type="text"
              :readonly="true"
              v-model="ctxData.propertyForm.label"
              autocomplete="off"
              placeholder="请选择属性名称同步属性标签"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="数据类型" prop="type">
            <el-select v-model.number="ctxData.propertyForm.type" style="width: 220px" placeholder="请选择数据类型">
              <el-option
                v-for="item in ctxData.typeOptions"
                :key="'type_' + item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type === 2" label="小数位数" prop="decimals">
            <el-input
              style="width: 220px"
              type="text"
              v-model.number="ctxData.propertyForm.decimals"
              autocomplete="off"
              placeholder="请输入小数位数"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type !== 3" label="单位" prop="unit">
            <el-input
              style="width: 220px"
              type="text"
              v-model="ctxData.propertyForm.unit"
              autocomplete="off"
              placeholder="请输入单位"
            >
            </el-input>
          </el-form-item>
          <el-divider content-position="center"> 映射参数 </el-divider>
          <el-form-item label="采集接口名称" prop="collName">
            <el-select
              v-model="ctxData.propertyForm.collName"
              style="width: 220px"
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
              style="width: 220px"
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
              style="width: 220px"
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
          <el-form-item label="掉电保存">
            <el-switch
              style="width: 220px"
              v-model="ctxData.propertyForm.saveEnable"
              inline-prompt
              active-text="是"
              inactive-text="否"
            />
          </el-form-item>
          <el-form-item label="计算公式">
            <el-input
              style="width: 220px"
              type="text"
              readonly
              v-model="ctxData.propertyForm.formula"
              autocomplete="off"
              placeholder="请编辑计算公式"
            >
              <template #append><el-button @click="showFormula()">编辑</el-button></template>
            </el-input>
          </el-form-item>
          <el-divider content-position="center"> 报警参数 </el-divider>
          <el-form-item v-if="ctxData.propertyForm.type !== 3" label="范围报警" prop="minMaxAlarm">
            <el-switch
              style="width: 220px"
              v-model="ctxData.propertyForm.minMaxAlarm"
              inline-prompt
              active-text="是"
              inactive-text="否"
            />
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type !== 3" label="步长报警" prop="stepAlarm">
            <el-switch
              style="width: 220px"
              v-model="ctxData.propertyForm.stepAlarm"
              inline-prompt
              active-text="是"
              inactive-text="否"
            />
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.minMaxAlarm"
            label="最小值"
            prop="min"
          >
            <el-input
              style="width: 220px"
              type="text"
              v-model="ctxData.propertyForm.min"
              autocomplete="off"
              placeholder="请输入最小值"
            >
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.minMaxAlarm"
            label="最大值"
            prop="max"
          >
            <el-input
              style="width: 220px"
              type="text"
              v-model="ctxData.propertyForm.max"
              autocomplete="off"
              placeholder="请输入最大值"
            >
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.stepAlarm"
            label="步长"
            prop="step"
          >
            <el-input
              style="width: 220px"
              type="text"
              v-model="ctxData.propertyForm.step"
              autocomplete="off"
              placeholder="请输入步长"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type === 3" label="字符串长度报警" prop="dataLengthAlarm">
            <el-switch
              style="width: 220px"
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
              style="width: 220px"
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
          <el-button @click="cancelPorpertySubmit()">取消</el-button>
          <el-button type="primary" @click="submitPorpertyForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- 导入属性 -->
    <!-- <el-dialog
      v-model="ctxData.psFlag"
      title="导入属性"
      width="600px"
      :before-close="beforeCloseUploadProperty"
      :close-on-click-modal="false"
    >
      <el-upload
        ref="uploadPropertyRef"
        action=""
        :auto-upload="false"
        :http-request="myRequestProperty"
        :limit="1"
        :on-exceed="handleExceedProperty"
        :before-upload="beforeUploadProperty"
      >
        <el-button type="primary">选择文件</el-button>
        <template #tip>
          <div class="el-upload__tip">只能上传一个文件，只支持xlsx格式文件！</div>
        </template>
      </el-upload>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelUploadProperty">取消</el-button>
          <el-button type="primary" @click="submitUploadProperty">上传</el-button>
        </span>
      </template>
    </el-dialog> -->
    <!-- 编辑计算公式 -->
    <el-dialog v-model="ctxData.fFlag" title="编辑计算公式" width="600px">
      <div class="dialog-content">
        <el-form :model="ctxData.formulaForm" label-position="right" label-width="120px">
          <el-form-item label="运算符号">
            <div style="width: 100%">
              <div v-for="item in ctxData.operationList" class="operation" :key="item.name">
                <div style="padding-right: 10px">
                  <el-button style="width: 100%" @click="setOperation(item)">{{ item.label }}</el-button>
                </div>
              </div>
            </div>
          </el-form-item>
          <el-form-item label="设备属性">
            <el-select
              v-model="ctxData.pItem"
              filterable
              placeholder="请选择设备属性"
              @change="setProperty(ctxData.pItem)"
              style="width: 100%"
            >
              <el-option v-for="item in ctxData.propertyInfo" :key="item.name" :label="item.label" :value="item.name" />
            </el-select>
          </el-form-item>

          <el-form-item label="当前公式">
            <el-input
              type="text"
              clearable
              v-model="ctxData.curFormula"
              autocomplete="off"
              placeholder="请编辑计算公式"
              @blur="onBulr"
            ></el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelSave()">取消</el-button>
          <el-button type="primary" @click="saveFormula()">保存</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- 选择模型属性 -->
    <el-dialog v-model="ctxData.mpFlag" title="模型属性选择" width="1000px">
      <div class="dialog-content">
        <div class="dialog-cLeft">
          <el-card class="card-box" shadow="hover">
            <template #header>
              <div class="card-header">
                <span>请选择采集模型</span>
              </div>
            </template>
            <div class="card-main">
              <el-radio-group v-model="ctxData.curModelName" style="width: 100%">
                <div
                  class="upload-item"
                  v-for="item in ctxData.modelList"
                  :key="'tml' + item.name"
                  :class="{ isClicked: item.isClicked }"
                >
                  <div class="upload-item-text" @click="modelClicked(item)">
                    <el-radio :label="item.name">{{ item.label + ' (' + item.name + ')' }}</el-radio>
                  </div>
                </div>
              </el-radio-group>
            </div>
          </el-card>
        </div>
        <div class="dialog-cRight">
          <el-card class="card-box" shadow="hover">
            <template #header>
              <div class="card-header">
                <span>请选择模型属性</span>
              </div>
            </template>
            <div class="card-main">
              <el-radio-group v-model="ctxData.curPropertyName" style="width: 100%">
                <div
                  class="upload-item"
                  v-for="item in ctxData.propertyList"
                  :key="'tml' + item.name"
                  :class="{ isClicked: item.isClicked }"
                >
                  <div class="upload-item-text" @click="propertyClicked(item)">
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
          <el-button @click="cancelSelect()">取消</el-button>
          <el-button type="primary" @click="selectModelProperty()">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { Search, Back } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import variables from 'styles/variables.module.scss'
import VirtualServiceApi from 'api/virtualService.js'
import InterfaceApi from 'api/interface.js'
import DeviceModelApi from 'api/deviceModel.js'
import { userStore } from 'stores/user'
const users = userStore()
const props = defineProps({
  curVirtualDevice: {
    type: Object,
    default: {},
  },
})
console.log(props)
const emit = defineEmits(['changeDpFlag'])
const toVirtualDevice = () => {
  emit('changeDpFlag')
}
const contentRef = ref(null)
const ctxData = reactive({
  formulaForm: {},
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
  PropertyTableData: [],
  pFlag: false,
  pTitle: '添加属性',
  PropertyInfo: '',
  propertyForm: {
    name: '', // 属性名称
    label: '', // 属性标签
    type: 0, // 数据类型
    unit: '', // 单位，只有uint32，int32，double有效
    decimals: 0, // 小数位数，只有double有效
    //params
    collName: '', //采集接口名称
    deviceName: '', //设备名称
    propertyName: '', //设备属性名称
    saveEnable: false, //掉电保存
    formula: '', //计算公式
    min: '', // 属性最小值，只有uint32，int32，double有效
    max: '', // 属性最大值，只有uint32，int32，double有效
    minMaxAlarm: false, // 范围报警，只有uint32，int32，double有效
    step: '', // 步长，只有uint32，int32，double有效
    stepAlarm: false, // 步长报警，只有uint32，int32，double有效
    dataLength: '', // 字符串长度，只有string有效
    dataLengthAlarm: false, // 字符串长度报警，只有string有效
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
  paramName: {
    name: '属性名称',
    label: '属性标签',
    type: '数据类型',
    min: '最小值',
    max: '最大值',
    minMaxAlarm: '范围报警',
    step: '步长',
    stepAlarm: '步长报警',
    unit: '单位',
    decimals: '小数位数',
    dataLength: '字符串长度',
    dataLengthAlarm: '字符串长度报警',
    collName: '采集接口名称',
    deviceName: '设备名称',
    propertyName: '属性名称',
    saveEnable: '掉电保存',
    formula: '计算公式',
  },
  propertyRules: {
    name: [
      {
        required: true,
        message: '属性名称不能为空',
        trigger: 'blur',
      },
    ],
    label: [
      {
        required: true,
        message: '属性标签不能为空',
        trigger: 'blur',
      },
    ],
    type: [
      {
        required: true,
        message: '数据类型不能为空',
        trigger: 'blur',
      },
    ],
    decimals: [
      {
        type: 'number',
        message: '小数位数只能输入数字',
      },
    ],
    collName: [
      {
        required: true,
        message: '采集接口名称不能为空',
        trigger: 'blur',
      },
    ],
    deviceName: [
      {
        required: true,
        message: '设备名称不能为空',
        trigger: 'blur',
      },
    ],
    propertyName: [
      {
        required: true,
        message: '设备属性名称不能为空',
        trigger: 'blur',
      },
    ],
  },
  collOptions: [], //采集接口选项
  nodeOptions: [], //设备选项
  propertyOptions: [], //设备属性选项
  nameTotsl: {}, //设备名称对应设备模板名称
  propertyInfo: [],
  operationList: [
    { name: '+', label: '+', describe: '加' },
    { name: '-', label: '-', describe: '减' },
    { name: '*', label: '*', describe: '乘' },
    { name: '/', label: '/', describe: '除' },
    { name: '%', label: '%', describe: '模' },
    { name: '^', label: '^', describe: '次幂' },
    { name: '(', label: '(', describe: '左括号' },
    { name: ')', label: ')', describe: '右括号' },
  ],
  curFormula: '',
  selectionStart: 0,
  pItem: '',
  fFlag: false,
  selectedProperties: [],
  modelList: [], //模型列表
  curModelName: '', //当前模型名称
  curPropertyName: '', //当前模型属性名称
  propertyList: [], //属性列表
  mpFlag: false, //模型属性弹窗标识
  curProperty: {},
})
// 获取虚拟设备属性列表
const getPropertiesList = (flag) => {
  const pData = {
    token: users.token,
    data: {
      deviceName: props.curVirtualDevice.name,
    },
  }
  VirtualServiceApi.getPropertyList(pData).then(async (res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.PropertyTableData = res.data
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
getPropertiesList()
// 过滤表格数据
const filterTableData = computed(() => {
  let PropertyInfo = ctxData.PropertyInfo
  return ctxData.PropertyTableData.filter(
    (item) =>
      !PropertyInfo ||
      item.name.toLowerCase().includes(PropertyInfo.toLowerCase()) ||
      item.label.toLowerCase().includes(PropertyInfo.toLowerCase())
  ).slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  let PropertyInfo = ctxData.PropertyInfo
  return ctxData.PropertyTableData.filter(
    (item) =>
      !PropertyInfo ||
      item.name.toLowerCase().includes(PropertyInfo.toLowerCase()) ||
      item.label.toLowerCase().includes(PropertyInfo.toLowerCase())
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

// 获取采集模型列表
const getModelList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  DeviceModelApi.getDeviceModelList(pData).then((res) => {
    console.log('getDeviceModelList -> res = ', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.modelList = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
const modelClicked = (item) => {
  ctxData.curModelName = item.name
  getDeviceModelProperty()
}

const propertyClicked = (item) => {
  ctxData.curProperty = item
}
const selectModelProperty = () => {
  ctxData.propertyForm.name = ctxData.curProperty.name
  ctxData.propertyForm.label = ctxData.curProperty.label
  cancelSelect()
}
const cancelSelect = () => {
  ctxData.mpFlag = false
}
// 获取设备模型属性
const getDeviceModelProperty = () => {
  const pData = {
    token: users.token,
    data: {
      name: ctxData.curModelName,
    },
  }
  DeviceModelApi.getDeviceModelProperty(pData).then((res) => {
    console.log('getDeviceModelProperty -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.propertyList = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
// 显示选择属性弹出框
const showModelProperty = () => {
  ctxData.mpFlag = true
  getModelList()
}

// 添加属性
const addProperty = () => {
  ctxData.pFlag = true
  ctxData.pTitle = '添加属性'
}
// 编辑属性
const editProperty = (row) => {
  console.log('editProperty -> row ', row)
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
  ctxData.propertyForm.formula = row.params.formula
  ctxData.propertyForm.saveEnable = row.params.saveEnable
  if (row.type !== 3) {
    ctxData.propertyForm['min'] = row.alarmParams.min
    ctxData.propertyForm['max'] = row.alarmParams.max
    ctxData.propertyForm['minMaxAlarm'] = row.alarmParams.minMaxAlarm
    ctxData.propertyForm['step'] = row.alarmParams.step
    ctxData.propertyForm['stepAlarm'] = row.alarmParams.stepAlarm
  } else {
    ctxData.propertyForm['dataLength'] = row.alarmParams.dataLength
    ctxData.propertyForm['dataLengthAlarm'] = row.alarmParams.dataLengthAlarm
  }
  console.log('ctxData.propertyForm', ctxData.propertyForm)
  getCollDevices(ctxData.propertyForm.collName).then((res) => {
    handleDevice(res)
    getNodeProperty(ctxData.propertyForm.deviceName)
  })
}
const propertyFormRef = ref(null)
const submitPorpertyForm = () => {
  propertyFormRef.value.validate((valid) => {
    console.log('valid', valid)
    if (valid) {
      //
      let property = {}
      property['name'] = ctxData.propertyForm.name
      property['label'] = ctxData.propertyForm.label
      property['type'] = ctxData.propertyForm.type
      property['decimals'] = ctxData.propertyForm.decimals
      property['unit'] = ctxData.propertyForm.unit
      let params = {
        collName: ctxData.propertyForm.collName,
        deviceName: ctxData.propertyForm.deviceName,
        propertyName: ctxData.propertyForm.propertyName,
        formula: ctxData.propertyForm.formula,
        saveEnable: ctxData.propertyForm.saveEnable,
      }
      let alarmParams = {}
      if (ctxData.propertyForm.type !== 3) {
        alarmParams['min'] = ctxData.propertyForm.min
        alarmParams['max'] = ctxData.propertyForm.max
        alarmParams['minMaxAlarm'] = ctxData.propertyForm.minMaxAlarm
        alarmParams['step'] = ctxData.propertyForm.step
        alarmParams['stepAlarm'] = ctxData.propertyForm.stepAlarm
      } else {
        alarmParams['dataLength'] = ctxData.propertyForm.dataLength
        alarmParams['dataLengthAlarm'] = ctxData.propertyForm.dataLengthAlarm
      }
      property['alarmParams'] = alarmParams
      property['params'] = params
      console.log('submitPorpertyForm -> property', property)
      const pData = {
        token: users.token,
        data: {
          deviceName: props.curVirtualDevice.name,
          property: property,
        },
      }
      if (ctxData.pTitle.includes('添加')) {
        console.log('添加属性')
        VirtualServiceApi.addProperty(pData).then((res) => {
          handleResult(res, getPropertiesList)
          cancelPorpertySubmit()
        })
      } else {
        console.log('编辑属性')
        VirtualServiceApi.editProperty(pData).then((res) => {
          handleResult(res, getPropertiesList)
          cancelPorpertySubmit()
        })
      }
    } else {
      return false
    }
  })
}
//取消提交
const cancelPorpertySubmit = () => {
  ctxData.pFlag = false
  propertyFormRef.value.resetFields()
  initPropertyForm()
}
//初始化接口表单
const initPropertyForm = () => {
  ctxData.propertyForm = {
    name: '', // 属性名称
    label: '', // 属性标签
    type: 0, // 数据类型
    unit: '', // 单位，只有uint32，int32，double有效
    decimals: 0, // 小数位数，只有double有效
    //params
    collName: '', //采集接口名称
    deviceName: '', //设备名称
    propertyName: '', //设备属性名称
    saveEnable: false, //掉电保存
    formula: '', //计算公式
    min: '', // 属性最小值，只有uint32，int32，double有效
    max: '', // 属性最大值，只有uint32，int32，double有效
    minMaxAlarm: false, // 范围报警，只有uint32，int32，double有效
    step: '', // 步长，只有uint32，int32，double有效
    stepAlarm: false, // 步长报警，只有uint32，int32，double有效
    dataLength: '', // 字符串长度，只有string有效
    dataLengthAlarm: false, // 字符串长度报警，只有string有效
  }
}
//处理关闭属性弹窗
const handleCloseProperty = (done) => {
  cancelPorpertySubmit()
}

const handleSelectionChange = (val) => {
  ctxData.selectedProperties = val
}

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
}
// 选择采集接口触发事件
const selectDeviceName = (deviceName) => {
  getNodeProperty(deviceName)
}
// 获取设备属性列表
const getNodeProperty = (deviceName) => {
  const pData = {
    token: users.token,
    data: {
      name: ctxData.nameTotsl[deviceName],
    },
  }
  DeviceModelApi.getDeviceModelProperty(pData).then((res) => {
    if (res.code === '0') {
      ctxData.propertyOptions = []
      res.data.forEach((item) => {
        const temp = {
          label: item.name + '(' + item.label + ')',
          name: item.name,
        }
        ctxData.propertyOptions.push(temp)
      })
    } else {
      showOneResMsg(res)
    }
  })
}
//显示计算公式弹窗
const showFormula = () => {
  ctxData.fFlag = true
  ctxData.curFormula = ctxData.propertyForm.formula
  getDeviceProperty()
}
//获取设备属性
const getDeviceProperty = () => {
  ctxData.propertyInfo = [{ name: 't', label: '当前属性' }]
  ctxData.PropertyTableData.forEach((item) => {
    ctxData.propertyInfo.push({
      name: item.name,
      label: item.label,
    })
  })
}
//设置运算符号
const setOperation = (item) => {
  let length = ctxData.curFormula.length
  if (length !== ctxData.selectionStart) {
    let a = ctxData.curFormula.substr(0, ctxData.selectionStart)
    let b = ctxData.curFormula.substr(ctxData.selectionStart, length - ctxData.selectionStart)
    ctxData.curFormula = a + item.name + b
    ctxData.selectionStart += item.name.length
  } else {
    ctxData.curFormula += item.name
    ctxData.selectionStart += item.name.length
  }
}
//设置属性
const setProperty = (item) => {
  let length = ctxData.curFormula.length
  if (length !== ctxData.selectionStart) {
    let a = ctxData.curFormula.substr(0, ctxData.selectionStart)
    let b = ctxData.curFormula.substr(ctxData.selectionStart, length - ctxData.selectionStart)
    ctxData.curFormula = a + '(' + item + ')' + b
    ctxData.selectionStart += item.length + 2
  } else {
    ctxData.curFormula += '(' + item + ')'
    ctxData.selectionStart += item.length + 2
  }
}
const onBulr = (obj) => {
  ctxData.selectionStart = obj.srcElement.selectionStart
}
//保存公式
const saveFormula = () => {
  ctxData.fFlag = false
  ctxData.propertyForm.formula = ctxData.curFormula
}
//取消保存
const cancelSave = () => {
  ctxData.fFlag = false
  ctxData.curFormula = ''
}
//
const refresh = () => {
  getPropertiesList(1)
}
// 删除属性,可以批量删除
const deleteProperty = () => {
  let pList = []
  if (ctxData.selectedProperties.length === 0) {
    ElMessage.info('请至少选择一个属性！')
    return
  } else {
    ctxData.selectedProperties.forEach((item) => {
      pList.push(item.name)
    })
  }
  console.log('pList', pList)
  ElMessageBox.confirm('确认要删除这些属性吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          deviceName: props.curVirtualDevice.name,
          PropertyNames: pList,
        },
      }
      VirtualServiceApi.deletePorperties(pData).then((res) => {
        handleResult(res, getPropertiesList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
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
.tName {
  line-height: 14px;
  font-size: 14px;
  border-left: 3px solid #3054eb;
  padding-left: 15px;
}
.operation {
  display: inline-block;
  width: 25%;
  height: 46px;
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
</style>
