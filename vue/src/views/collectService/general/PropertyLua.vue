<template>
<div class="main-container">
  <div class="main">
    <div class="search-bar">
      <el-form :inline="true" ref="searchFormRef" status-icon label-width="120px">
        <el-form-item style="margin-right: 20px;margin-left: 20px;">
          <el-button type="primary" plain @click="toDeviceModel()">
            <el-icon class="el-input__icon"><back /></el-icon>
            返回设备模型
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-radio-group v-model="ctxData.varParamType" style="margin-right: 20px">
            <el-radio label="Properties" border>属性</el-radio>
            <!-- <el-radio label="Services" border>服务</el-radio> -->
          </el-radio-group>
          <el-divider direction="vertical" style="height: 2em" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" plain class="right-btn" @click="importDPS()">
            <el-icon class="el-input__icon"><download /></el-icon>
            导入属性
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" plain class="right-btn" @click="exportDPS()">
            <el-icon class="el-input__icon"><upload /></el-icon>
            导出属性
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div v-if="ctxData.varParamType === 'Properties'" style="height: calc(100% - 20px)">
      <div class="search-bar" style="display: flex;">
        <div class="title" style="position: relative;margin-right: 40px;justify-content: flex-start;padding: 0px 0px;height: 40px;">
          <div class="tName">
            {{ props.curDeviceModel.label }}：{{ ctxData.varParamType === 'Properties' ? '属性列表' : '服务列表' }}
          </div>
        </div>
        <el-form :inline="true" ref="searchFormRef" status-icon label-width="90px">
          <el-form-item label="">
            <el-input style="width: 200px" placeholder="请输入属性名称" v-model="ctxData.deviceModelProperty">
              <template #prefix>
                <el-icon class="el-input__icon"><search /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" bg class="right-btn" @click="syncDeviceModelProperty()">
              <el-icon class="btn-icon">
                <Icon name="local-tongbu" size="14px" color="#ffffff" />
              </el-icon>
              同步
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" bg class="right-btn" @click="addDeviceModelProperty()">
              <el-icon class="btn-icon">
                <Icon name="local-add" size="14px" color="#ffffff" />
              </el-icon>
              添加
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-button type="danger" bg class="right-btn" @click="deleteDeviceModelProperty()">
              <el-icon class="btn-icon">
                <Icon name="local-delete" size="14px" color="#ffffff" />
              </el-icon>
              删除
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-button style="color: #fff; margin-left: 20px" color="#2EA554" class="right-btn" @click="refresh('deviceModelProperty')">
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
          :data="filterDMPTableData"
          :cell-style="ctxData.cellStyle"
          :header-cell-style="ctxData.headerCellStyle"
          :max-height="ctxData.tableMaxHeight"
          style="width: 100%"
          stripe
          @selection-change="handleSelectionChange"
          @row-dblclick="editDeviceModelProperty"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column sortable prop="name" label="属性名称" width="auto" min-width="150" align="center" />
          <el-table-column sortable prop="label" label="属性标签" width="auto" min-width="150" align="center" />
          <el-table-column sortable label="读写属性" width="auto" min-width="80" align="center">
            <template #default="scope">
              {{ ctxData.accessModeNames['am' + scope.row.accessMode] }}
            </template>
          </el-table-column>
          <el-table-column sortable label="属性类型" width="auto" min-width="80" align="center">
            <template #default="scope">
              {{ ctxData.typeNames['t' + scope.row.type] }}
            </template>
          </el-table-column>
          <el-table-column sortable label="小数位数" width="auto" min-width="80" align="center">
            <template #default="scope">
              {{ scope.row.decimals === undefined || scope.row.decimals === '' ? 0 : scope.row.decimals }}
            </template>
          </el-table-column>
          <el-table-column sortable prop="unit" label="单位" width="auto" min-width="80" align="center" />

          <el-table-column label="操作" width="auto" min-width="200" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="editDeviceModelProperty(scope.row)" text type="primary">编辑</el-button>
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
    <div e-else>
      <!-- <div class="title" style="top: 60px; height: 76px; padding: 20px 0; justify-content: space-between">
          <div class="tName">{{ ctxData.varParamType === 'Properties' ? '属性列表' : '服务列表' }}</div>
          <div style="display: flex">
            <el-input style="width: 200px" placeholder="请输入服务名称" v-model="ctxData.deviceModelService">
              <template #prefix>
                <el-icon class="el-input__icon"><search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" bg class="right-btn" @click="addDeviceModelService()">
              <el-icon class="btn-icon tianjia"></el-icon>
              添加
            </el-button>
            <el-button style="color: #fff" color="#2EA554" class="right-btn" @click="refresh('deviceModelService')">
              <el-icon class="btn-icon shuaxin"></el-icon>
              刷新
            </el-button>
          </div>
        </div>
        <div class="content" ref="contentRef" style="top: 136px">
          <el-table
            :data="filterDMSTableData"
            :cell-style="ctxData.cellStyle"
            :header-cell-style="ctxData.headerCellStyle"
            :max-height="ctxData.tableMaxHeight"
            style="width: 100%"
            stripe
          >
            <el-table-column prop="name" label="属性标签" width="auto" min-width="150" align="center">
            </el-table-column>
            <el-table-column prop="label" label="属性属性标签" width="auto" min-width="200" align="center">
            </el-table-column>
            <el-table-column label="读写属性" width="auto" min-width="100" align="center">
              <template #default="scope">
                {{ ctxData.accessModeNames[scope.row.accessMode] }}
              </template>
            </el-table-column>
            <el-table-column label="属性类型" width="auto" min-width="100" align="center">
              <template #default="scope">
                {{ ctxData.typeNames[scope.row.type] }}
              </template>
            </el-table-column>
            <el-table-column label="参数" width="auto" min-width="500" align="center">
              <template #default="scope">
                <el-popover @show="showParam(scope.row.type)" trigger="hover" effect="light" width="auto">
                  <template #default>
                    <div class="param-content">
                      <div class="pc-title">
                        <div class="pct-info">
                          <b> {{ scope.row.name }} </b>
                          {{ '参数详情' }}
                        </div>
                      </div>
                      <div class="pc-content">
                        <div class="param-item" v-for="(item, key, index) of scope.row.params">
                          <div class="param-value">{{ ctxData.paramName[key] }}：</div>
                          <div class="param-name">{{ item }}</div>
                        </div>
                      </div>
                    </div>
                  </template>
                  <template #reference>
                    <el-tag size="large">{{ scope.row.param }}</el-tag>
                  </template>
                </el-popover>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="auto" min-width="200" align="center">
              <template #default="scope">
                <el-button @click="editDeviceModelPorperty(scope.row)" text type="primary">编辑</el-button>
                <el-button @click="deleteDeviceModelPorperty(scope.row)" text type="danger">删除</el-button>
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
              :total="ctxData.serviceList.length"
              @current-change="handleCurrentChange"
              @size-change="handleSizeChange"
              background
              layout="total, sizes, prev, pager, next, jumper"
              style="margin-top: 46px"
            ></el-pagination>
          </div>
        </div> -->
    </div>
    <!-- 添加编辑设备模型属性 -->
    <el-dialog
      v-model="ctxData.pFlag"
      :title="ctxData.pTitle"
      width="600px"
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
        >
          <el-form-item label="属性名称" prop="name">
            <el-input
              :disabled="ctxData.pTitle.includes('编辑')"
              type="text"
              v-model="ctxData.propertyForm.name"
              autocomplete="off"
              placeholder="请输入属性名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="属性标签" prop="label">
            <el-input type="text" v-model="ctxData.propertyForm.label" autocomplete="off" placeholder="请输入属性标签">
            </el-input>
          </el-form-item>
          <el-form-item label="读写属性" prop="accessMode">
            <el-select v-model="ctxData.propertyForm.accessMode" style="width: 100%" placeholder="请选择读写属性">
              <el-option
                v-for="item in ctxData.accessModeOptions"
                :key="'accessMode_' + item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="属性类型" prop="type">
            <el-select v-model.number="ctxData.propertyForm.type" style="width: 100%" placeholder="请选择属性类型">
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
              type="text"
              v-model.number="ctxData.propertyForm.decimals"
              autocomplete="off"
              placeholder="请输入小数位数"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type !== 3" label="单位" prop="unit">
            <el-input type="text" v-model="ctxData.propertyForm.unit" autocomplete="off" placeholder="请输入单位">
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

    <!-- 导入设备属性和服务 -->
    <el-dialog
      v-model="ctxData.psFlag"
      title="导入设备属性"
      width="600px"
      :before-close="beforeCloseUploadDPS"
      :close-on-click-modal="false"
    >
      <el-upload
        ref="uploadDPSRef"
        action=""
        :auto-upload="false"
        :http-request="myRequestDPS"
        :limit="1"
        :on-exceed="handleExceedDPS"
        :before-upload="beforeUploadDPS"
      >
        <el-button type="primary">选择文件</el-button>
        <template #tip>
          <div class="el-upload__tip">只能上传一个文件，只支持xlsx格式文件！</div>
        </template>
      </el-upload>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelUploadDPS">取消</el-button>
          <el-button type="primary" @click="submitUploadDPS">上传</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</div>
</template>
<script setup>
import { Search, Back, Download, Upload } from '@element-plus/icons-vue'
import variables from 'styles/variables.module.scss'
import DeviceModelApi from 'api/deviceModel.js'
import { userStore } from 'stores/user'
const users = userStore()

const props = defineProps({
  curDeviceModel: {
    type: Object,
    default: {},
  },
})
console.log('id -> props', props)
// 返回设备模型
const emit = defineEmits(['changeDpFlag'])
const toDeviceModel = () => {
  console.log('toDeviceModel')
  emit('changeDpFlag')
}

const contentRef = ref(null)
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
  varParamType: 'Properties', //变量的参数类型
  propertyList: [],
  serviceList: [],
  deviceModelProperty: '',
  deviceModelService: '',
  accessModeOptions: [
    //读写模式
    { label: '只读', value: 0 },
    { label: '只写', value: 1 },
    { label: '读写', value: 2 },
  ],
  typeOptions: [
    //属性类型
    { label: 'uint32', value: 0 },
    { label: 'int32', value: 1 },
    { label: 'double', value: 2 },
    { label: 'string', value: 3 },
  ],
  accessModeNames: {
    am0: '只读',
    am1: '只写',
    am2: '读写',
  },
  typeNames: {
    t0: 'uint32',
    t1: 'int32',
    t2: 'double',
    t3: 'string',
  },
  pFlag: false, //属性对话框标识
  pTitle: '添加属性', //属性对话框titleName
  propertyForm: {
    name: '', // 属性名称，只能是字母+数字的组合，不可以是中文
    label: '', // 属性标签
    accessMode: 0, // 读写属性
    type: 0, // 属性类型
    unit: '', // 单位，只有uint32，int32，double有效
    decimals: 0, // 小数位数，只有double有效
    //params
    min: '', // 属性最小值，只有uint32，int32，double有效
    max: '', // 属性最大值，只有uint32，int32，double有效
    minMaxAlarm: false, // 范围报警，只有uint32，int32，double有效
    step: '', // 布长，只有uint32，int32，double有效
    stepAlarm: false, // 布长报警，只有uint32，int32，double有效
    dataLength: '', // 字符串长度，只有string有效
    dataLengthAlarm: false, // 字符串长度报警，只有string有效
  },
  paramName: {
    name: '属性名称',
    label: '属性标签',
    accessMode: '读写模式',
    type: '属性类型',
    min: '最小值',
    max: '最大值',
    minMaxAlarm: '范围报警',
    step: '布长',
    stepAlarm: '布长报警',
    unit: '单位',
    decimals: '小数位数',
    dataLength: '字符串长度',
    dataLengthAlarm: '字符串长度报警',
  },
  propertyRules: {
    name: [
      {
        required: true,
        message: '属性名称不能为空',
        trigger: 'blur',
      },
    ],
    accessMode: [
      {
        required: true,
        message: '读写属性不能为空',
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
        trigger: 'blur',
      },
    ],
  },
  psFlag: false,
  selectedProperties: [],
})
// 获取设备模型属性
const getDeviceModelProperty = (flag) => {
  const pData = {
    token: users.token,
    data: {
      name: props.curDeviceModel.name,
    },
  }
  DeviceModelApi.getDeviceModelProperty(pData).then(async (res) => {
    console.log('getDeviceModelProperty -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.propertyList = res.data
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

getDeviceModelProperty()

const refresh = (value) => {
  getDeviceModelProperty(1)
}

// 获取设备模型服务
const getDeviceModelService = (flag) => {
  const pData = {
    token: users.token,
    data: {
      name: props.curDeviceModel.name,
    },
  }
  DeviceModelApi.getDeviceModelService(pData).then(async (res) => {
    console.log('getDeviceModelService -> res', res)
    if (res.code === '0') {
      ctxData.serviceList = res.data
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

const filterDMPTableData = computed(() => {
  console.log('ctxData.propertyList ->', ctxData.propertyList)
  return ctxData.propertyList
    .filter((item) => {
      var a = !ctxData.deviceModelProperty
      var b = item.name.toLowerCase().includes(ctxData.deviceModelProperty.toLowerCase())

      return a || b
    })
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  return ctxData.propertyList.filter((item) => {
    var a = !ctxData.deviceModelProperty
    var b = item.name.toLowerCase().includes(ctxData.deviceModelProperty.toLowerCase())

    return a || b
  })
})
const filterDMSTableData = computed(() => {
  console.log('ctxData.serviceList ->', ctxData.serviceList)
  return ctxData.serviceList
    .filter((item) => {
      var a = !ctxData.deviceModelService
      var b = item.name.toLowerCase().includes(ctxData.deviceModelService.toLowerCase())
      return a || b
    })
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
// 处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
// 导入设备属性和服务
const importDPS = () => {
  console.log('importDPS')
  ctxData.psFlag = true
}
// 导出设备属性和服务
const exportDPS = () => {
  console.log('exportDPS')
  const pData = {
    token: users.token,
    responseType: 'blob',
    data: {
      name: props.curDeviceModel.name,
    },
  }
  DeviceModelApi.exportDeviceModelProptyAndService(pData).then((res) => {
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
 * 提交上传的插件
 */
const submitUploadDPS = () => {
  uploadDPSRef.value && uploadDPSRef.value.submit()
}

// 取消插件自带的xhr
const myRequestDPS = (obj) => {}
const uploadDPSRef = ref(null)
/**
 * 上传文件大于limit时事件
 * @param {要上传的文件} filesimportPlugin
 */
const handleExceedDPS = (files) => {
  uploadDPSRef.value.clearFiles()
  //超过limit取第一个
  const file = files[0]
  uploadDPSRef.value.handleStart(file)
}
/**
 * 文件上传前事件
 * @param {要上传的文件} file
 */
const beforeUploadDPS = (file) => {
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
    formData.append('name', props.curDeviceModel.name)
    formData.append('fileName', file)
    const pData = {
      token: users.token,
      contentType: 'multipart/form-data',
      data: formData,
    }
    DeviceModelApi.importDeviceModelProptyAndService(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success(res.message)
        ctxData.psFlag = false
        getDeviceModelProperty()
      } else {
        showOneResMsg(res)
      }
    })
  }
}
//取消上传属性和服务
const cancelUploadDPS = () => {
  ctxData.psFlag = false
  uploadRef.value.clearFiles()
}
// 弹窗右上角关闭事件处理
const beforeCloseUploadDPS = () => {
  cancelUploadDPS()
}
// 同步属性
const syncDeviceModelProperty = () => {
  const pData = {
    token: users.token,
    data: {
      name: props.curDeviceModel.name,
    },
  }
  DeviceModelApi.syncDeviceModelProperty(pData).then((res) => {
    console.log(res)
    handleResult(res, getDeviceModelProperty)
  })
}
// 添加属性
const addDeviceModelProperty = () => {
  console.log('addDeviceModelProperty')
  ctxData.pFlag = true
  ctxData.pTitle = '添加属性'
}
// 编辑属性
const editDeviceModelProperty = (row) => {
  ctxData.pFlag = true
  ctxData.pTitle = '编辑属性'
  ctxData.propertyForm.name = row.name
  ctxData.propertyForm.label = row.label
  ctxData.propertyForm.accessMode = row.accessMode
  ctxData.propertyForm.type = row.type
  ctxData.propertyForm.unit = row.unit
  if (row.type === 2) {
    ctxData.propertyForm.decimals = row.decimals
  }
  if (row.type !== 3) {
    ctxData.propertyForm['min'] = row.params.min
    ctxData.propertyForm['max'] = row.params.max
    ctxData.propertyForm['minMaxAlarm'] = row.params.minMaxAlarm
    ctxData.propertyForm['step'] = row.params.step
    ctxData.propertyForm['stepAlarm'] = row.params.stepAlarm
  } else {
    ctxData.propertyForm['dataLength'] = row.params.dataLength
    ctxData.propertyForm['dataLengthAlarm'] = row.params.dataLengthAlarm
  }
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
      property['accessMode'] = ctxData.propertyForm.accessMode
      property['type'] = ctxData.propertyForm.type
      property['unit'] = ctxData.propertyForm.unit
      if (ctxData.propertyForm.type === 2) {
        property['decimals'] = ctxData.propertyForm.decimals
      }
      let params = {}
      if (ctxData.propertyForm.type !== 3) {
        params['min'] = ctxData.propertyForm.min
        params['max'] = ctxData.propertyForm.max
        params['minMaxAlarm'] = ctxData.propertyForm.minMaxAlarm
        params['step'] = ctxData.propertyForm.step
        params['stepAlarm'] = ctxData.propertyForm.stepAlarm
      } else {
        params['dataLength'] = ctxData.propertyForm.dataLength
        params['dataLengthAlarm'] = ctxData.propertyForm.dataLengthAlarm
      }
      property['params'] = params
      console.log('submitPorpertyForm -> property', property)
      const pData = {
        token: users.token,
        data: {
          name: props.curDeviceModel.name,
          property: property,
        },
      }
      if (ctxData.pTitle.includes('添加')) {
        console.log('添加属性')
        DeviceModelApi.addDeviceModelProperty(pData).then((res) => {
          handleResult(res, getDeviceModelProperty)
          cancelPorpertySubmit()
        })
      } else {
        console.log('编辑属性')
        DeviceModelApi.editDeviceModelProperty(pData).then((res) => {
          handleResult(res, getDeviceModelProperty)
          cancelPorpertySubmit()
        })
      }
    } else {
      return false
    }
  })
}
const cancelPorpertySubmit = () => {
  ctxData.pFlag = false
  propertyFormRef.value.resetFields()
  initPropertyForm()
}
const handleCloseProperty = (done) => {
  cancelPorpertySubmit()
}

const handleSelectionChange = (val) => {
  ctxData.selectedProperties = val
  console.log('handleSelectionChange -> val =', val)
}
// 删除属性,可以批量删除
const deleteDeviceModelProperty = (row) => {
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
          name: props.curDeviceModel.name,
          properties: pList,
        },
      }
      DeviceModelApi.deleteDeviceModelProperty(pData).then((res) => {
        handleResult(res, getDeviceModelProperty)
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
const initPropertyForm = () => {
  ctxData.propertyForm = {
    name: '', //属性名称
    label: '', // 属性标签
    accessMode: 0, // 读写属性
    type: 0, // 属性类型
    //params
    min: '', // 属性最小值，只有uint32，int32，double有效
    max: '', // 属性最大值，只有uint32，int32，double有效
    minMaxAlarm: false, // 范围报警使能，只有uint32，int32，double有效
    step: '', // 布长，只有uint32，int32，double有效
    stepAlarm: false, // 布长报警，只有uint32，int32，double有效
    unit: '', // 单位，只有uint32，int32，double有效
    decimals: 0, // 小数位数，只有double有效
    dataLength: '', // 字符串长度，只有string有效
    dataLengthAlarm: false, // 字符串长度报警使能，只有string有效
  }
}
const addDeviceModelService = () => {
  //
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
