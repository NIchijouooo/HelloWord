<template>
  <div class="main">
    <div class="title" style="justify-content: space-between">
      <div>
        <el-button type="primary" plain @click="toDeviceModel()">
          <el-icon class="el-input__icon"><back /></el-icon>
          返回设备模型
        </el-button>
      </div>
      <div style="display: flex; align-items: center">
        <el-button type="primary" plain class="right-btn" @click="importDPS()">
          <el-icon class="el-input__icon"><download /></el-icon>
          导入模型属性
        </el-button>
        <el-button type="primary" plain class="right-btn" @click="exportDPS()">
          <el-icon class="el-input__icon"><upload /></el-icon>
          导出模型属性
        </el-button>
      </div>
    </div>
    <div class="title" style="top: 60px; height: 76px; padding: 20px 0; justify-content: space-between">
      <div class="tName">{{ props.curDeviceModel.label }}：属性列表</div>
      <div style="display: flex">
        <el-input style="width: 200px" placeholder="请输入属性名称" v-model="ctxData.deviceModelProperty">
          <template #prefix>
            <el-icon class="el-input__icon"><search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" bg class="right-btn" @click="addDeviceModelProperty()">
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
        <el-button type="danger" bg class="right-btn" @click="deleteDeviceModelProperty()">
          <el-icon class="btn-icon">
            <Icon name="local-delete" size="14px" color="#ffffff" />
          </el-icon>
          删除
        </el-button>
      </div>
    </div>
    <div class="content" ref="contentRef" style="top: 136px">
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
        <el-table-column prop="name" label="属性名称" width="auto" min-width="150" align="center"> </el-table-column>
        <el-table-column prop="label" label="属性标签" width="auto" min-width="150" align="center"> </el-table-column>
        <el-table-column label="读写属性" width="auto" min-width="80" align="center">
          <template #default="scope">
            {{ ctxData.accessModeNames['am' + scope.row.accessMode] }}
          </template>
        </el-table-column>
        <el-table-column prop="type" label="属性类型" width="auto" min-width="100" align="center">
          <template #default="scope">
            {{ ctxData.typeNames['t' + scope.row.type] }}
          </template>
        </el-table-column>
        <el-table-column label="小数位数" width="auto" min-width="80" align="center">
          <template #default="scope">
            {{ scope.row.decimals === '' ? 0 : scope.row.decimals }}
          </template>
        </el-table-column>
        <el-table-column prop="unit" label="单位" width="auto" min-width="80" align="center" />
        <el-table-column label="参数详情" width="auto" min-width="500" align="center">
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
                    <div class="param-item" v-for="(item, key, index) of scope.row.params">
                      <div class="param-value">{{ ctxData.paramName[key] }}：</div>
                      <div v-if="key === 'dataType'" class="param-name">{{ ctxData.dataTypeNames['dt' + item] }}</div>
                      <div v-if="key !== 'dataType'" class="param-name">
                        {{ typeof item === 'boolean' ? (item ? '是' : '否') : item }}
                      </div>
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
        <el-table-column label="操作" width="auto" min-width="200" align="center" fixed="right">
          <template #default="scope">
            <el-button @click="editDeviceModelProperty(scope.row)" text type="primary">编辑</el-button>
            <!-- <el-button @click="deleteDeviceModelProperty(scope.row)" text type="danger">删除</el-button> -->
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
    <!-- 添加编辑设备模型属性 -->
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
          inline="true"
        >
          <el-form-item label="属性名称" prop="name">
            <el-input
              style="width: 220px"
              :disabled="ctxData.pTitle.includes('编辑')"
              type="text"
              v-model="ctxData.propertyForm.name"
              autocomplete="off"
              placeholder="请输入属性名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="属性标签" prop="label">
            <el-input
              style="width: 220px"
              type="text"
              v-model="ctxData.propertyForm.label"
              autocomplete="off"
              placeholder="请输入属性标签"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="读写属性" prop="accessMode">
            <el-select v-model="ctxData.propertyForm.accessMode" style="width: 220px" placeholder="请选择读写属性">
              <el-option
                v-for="item in ctxData.accessModeOptions"
                :key="'accessMode_' + item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="单位" prop="unit">
            <el-input
              style="width: 220px"
              type="text"
              v-model="ctxData.propertyForm.unit"
              autocomplete="off"
              placeholder="请输入单位"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="属性类型" prop="type">
            <el-select v-model.number="ctxData.propertyForm.type" style="width: 220px" placeholder="请选择属性类型">
              <el-option
                v-for="item in ctxData.typeOptions"
                :key="'type_' + item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="小数位数" prop="decimals" v-show="ctxData.propertyForm.type === 2">
            <el-input
              style="width: 220px"
              type="text"
              v-model.number="ctxData.propertyForm.decimals"
              autocomplete="off"
              placeholder="请输入小数位数"
            >
            </el-input>
          </el-form-item>
          <div class="form-title"><div class="tName">配置参数</div></div>
          <el-form-item label="数据块号" prop="dbNumber">
            <el-input
              type="text"
              style="width: 220px"
              v-model="ctxData.propertyForm.dbNumber"
              autocomplete="off"
              placeholder="请输入数据块号"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="" prop="" style="width: 220px"> </el-form-item>
          <el-form-item label="PLC地址" prop="startAddr">
            <el-input
              type="text"
              style="width: 220px"
              v-model="ctxData.propertyForm.startAddr"
              autocomplete="off"
              placeholder="请输入PLC地址"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="数据类型" prop="dataType">
            <el-select v-model.number="ctxData.propertyForm.dataType" style="width: 220px" placeholder="请选择数据类型">
              <el-option
                v-for="item in ctxData.dataTypeOptions"
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
          <el-button @click="cancelPorpertySubmit()">取消</el-button>
          <el-button type="primary" @click="submitPorpertyForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 导入模型属性 -->
    <el-dialog
      v-model="ctxData.psFlag"
      title="导入模型属性"
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
  propertyList: [],
  deviceModelProperty: '',
  accessModeOptions: [
    //读写模式
    { label: '只读', value: 0 },
    { label: '只写', value: 1 },
    { label: '读写', value: 2 },
  ],
  typeOptions: [
    { label: 'uint32', value: 0 },
    { label: 'int32', value: 1 },
    { label: 'double', value: 2 },
    { label: 'string', value: 3 },
  ],
  dataTypeOptions: [
    //数据类型
    { label: 'uint8', value: 0 },
    { label: 'int8', value: 1 },
    { label: 'uint16', value: 2 },
    { label: 'int16', value: 3 },
    { label: 'uint32', value: 4 },
    { label: 'int32', value: 5 },
    { label: 'float', value: 6 },
    { label: 'double', value: 7 },
    { label: 'bool', value: 8 },
  ],
  typeNames: {
    t0: 'uint32',
    t1: 'int32',
    t2: 'double',
    t3: 'string',
  },
  dataTypeNames: {
    dt0: 'uint8',
    dt1: 'int8',
    dt2: 'uint16',
    dt3: 'int16',
    dt4: 'uint32',
    dt5: 'int32',
    dt6: 'float',
    dt7: 'double',
    dt8: 'bool',
  },
  accessModeNames: {
    am0: '只读',
    am1: '只写',
    am2: '读写',
  },
  pFlag: false, //属性对话框标识
  pTitle: '添加属性', //属性对话框titleName
  propertyForm: {
    name: '', // 属性名称，只能是字母+数字的组合，不可以是中文
    label: '', // 属性标签
    accessMode: 0, // 读写属性
    type: 0,
    //params
    dbNumber: '',
    startAddr: '',
    dataType: 0,
    unit: '', // 单位，只有uint32，int32，double有效
    decimals: 0, // 小数位数，只有double有效
  },
  paramName: {
    name: '属性名称',
    label: '属性标签',
    accessMode: '读写模式',
    type: '属性类型',
    decimals: '小数位数',
    dbNumber: '数据块号',
    dataType: '数据类型',
    startAddr: 'PLC地址',
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
        message: '属性类型不能为空',
        trigger: 'blur',
      },
    ],
    dataType: [
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
    dbNumber: [
      {
        required: true,
        message: '数据块号不能为空',
        trigger: 'blur',
      },
    ],
    startAddr: [
      {
        required: true,
        message: 'PLC地址不能为空',
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
const refresh = () => {
  getDeviceModelProperty(1)
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
    console.log('exportDeviceModelProptyAndService -> res', res)
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
  ctxData.propertyForm.decimals = row.decimals
  ctxData.propertyForm.unit = row.unit
  // s7
  ctxData.propertyForm.dbNumber = row.params.dbNumber
  ctxData.propertyForm.startAddr = row.params.startAddr
  ctxData.propertyForm.dataType = row.params.dataType
  console.log('ctxData.propertyForm', ctxData.propertyForm)
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
      property['decimals'] = ctxData.propertyForm.decimals
      property['unit'] = ctxData.propertyForm.unit

      let params = {
        dbNumber: ctxData.propertyForm.dbNumber,
        startAddr: ctxData.propertyForm.startAddr,
        dataType: ctxData.propertyForm.dataType,
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
    type: 0,
    decimals: 0,
    unit: '',
    dbNumber: '',
    startAddr: '',
    dataType: 0,
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
