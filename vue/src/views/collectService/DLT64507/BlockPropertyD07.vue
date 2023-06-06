<template>
<div>
  <div style="display: flex; justify-content: space-between;">
    <div class="title" style="position: relative;width: 40%;justify-content: flex-start;">
      <div class="tName">{{ props.curModelBlock.label }}：命令参数列表</div>
    </div>
    <div class="search-bar" style="text-align:right;">
      <el-form :inline="true" ref="searchFormRef" status-icon label-width="90px">
        <el-form-item style="margin-left: 20px;">
          <el-input style="width: 200px" placeholder="请输入命令参数名称" v-model="ctxData.deviceModelProperty">
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" plain class="right-btn" @click="importDPS()">
            <el-icon class="el-input__icon"><download /></el-icon>
            导入命令参数
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" plain class="right-btn" @click="exportDPS()">
            <el-icon class="el-input__icon"><upload /></el-icon>
            导出命令参数
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" bg class="right-btn" @click="addDeviceModelProperty()">
            <el-icon class="btn-icon"> <Icon name="local-add" size="14px" color="#ffffff" /> </el-icon> 添加
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
          <el-button type="danger" bg class="right-btn" @click="deleteDeviceModelProperty()">
            <el-icon class="btn-icon">
              <Icon name="local-delete" size="14px" color="#ffffff" />
            </el-icon>
            删除
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
  <div class="content" ref="contentRef">
    <el-table
      :data="filterDMPTableData"
      :cell-style="ctxData.cellStyle"
      :header-cell-style="ctxData.headerCellStyle"
      style="width: 100%"
      height="300"
      stripe
      @selection-change="handleSelectionChange"
      @row-dblclick="editDeviceModelProperty"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="属性名称" width="auto" min-width="150" align="center"> </el-table-column>
      <el-table-column prop="label" label="属性标签" width="auto" min-width="150" align="center"> </el-table-column>

      <el-table-column prop="rulerId" label="数据标识" width="auto" min-width="150" align="center"> </el-table-column>
      <el-table-column prop="format" label="数据格式" width="auto" min-width="150" align="center"> </el-table-column>
      <el-table-column prop="len" label="数据长度" width="auto" min-width="150" align="center"> </el-table-column>
      <el-table-column label="读写属性" width="auto" min-width="80" align="center">
        <template #default="scope">
          {{ ctxData.accessModeNames['am' + scope.row.accessMode] }}
        </template>
      </el-table-column>
      <el-table-column prop="type" label="数据类型" width="auto" min-width="100" align="center">
        <template #default="scope">
          {{ ctxData.typeNames['t' + scope.row.type] }}
        </template>
      </el-table-column>
      <el-table-column prop="unit" label="单位" width="auto" min-width="80" align="center" ></el-table-column>
      <el-table-column prop="blockAddOffset" label="块偏移地址" width="auto" min-width="150" align="center"> </el-table-column>
      <el-table-column prop="rulerAddOffset" label="标识偏移地址" width="auto" min-width="150" align="center"> </el-table-column>

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
  <!-- 添加编辑设备命令参数 -->
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



        <el-form-item label="" prop="" style="width: 220px" v-if="props.curModelBlock.format != 1"> </el-form-item>
        <el-form-item label="数据格式" prop="len">
          <el-input
            type="text"
            v-if="props.curModelBlock.format == 1"
            style="width: 220px"
            disabled
            v-model.number="ctxData.propertyForm.format"
            autocomplete="off"
            placeholder="请输入数据格式"
          >
          </el-input>

          <el-select
            v-else
            v-model.number="ctxData.propertyForm.format"
            style="width: 220px"

            placeholder="请选择数据格式"
          >
            <el-option
              v-for="item in ctxData.formatOptions"
              :key= item.value
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>


        <el-form-item label="" prop="" style="width: 220px" v-if="props.curModelBlock.len != 1"> </el-form-item>
        <el-form-item label="数据长度" prop="len">
         <el-input
            type="text"
            v-if="props.curModelBlock.len == 1"
            style="width: 220px"
            disabled
            v-model.number="ctxData.propertyForm.len"
            autocomplete="off"
            placeholder="请输入数据长度"
          >
          </el-input>

          <el-select
            v-else
            v-model.number="ctxData.propertyForm.len"
            style="width: 220px"

            placeholder="请选择数据长度"
          >
            <el-option
              v-for="item in ctxData.lenOptions"
              :key="'regCount_' + item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>


        <div class="form-title"><div class="tName">配置参数</div></div>

        <el-form-item label="数据标识" prop="rulerId">
          <el-input
            type="text"
            style="width: 220px"
            v-model="ctxData.propertyForm.rulerId"
            autocomplete="off"
            placeholder="请输入数据标识"
          >
          </el-input>
        </el-form-item>


        <el-form-item label="" prop="" style="width: 220px" v-if="props.curModelBlock.blockAddOffset != 1"> </el-form-item>

        <el-form-item label="块偏移地址" prop="blockAddOffset">
          <el-input
            type="text"
            style="width: 220px"
            v-model.number="ctxData.propertyForm.blockAddOffset"
            autocomplete="off"
            placeholder="请输入数据块偏移地址"
          >
          </el-input>
        </el-form-item>

        <el-form-item label="标识偏移地址" prop="rulerAddOffset">
          <el-input
            type="text"
            style="width: 220px"
            v-model.number="ctxData.propertyForm.rulerAddOffset"
            autocomplete="off"
            placeholder="请输入标识偏移地址"
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
  <!-- 编辑计算公式 -->
  <el-dialog v-model="ctxData.fFlag" title="编辑计算公式" width="600px">
    <div class="dialog-content">
      <el-form :model="ctxData.formulaForm" label-position="right" label-width="120px">
        <el-form-item label="运算符号">
          <div style="width: 100%">
            <div v-for="(item, index) in ctxData.operationList" class="operation" :key="index">
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
  <!-- 导入命令参数 -->
  <el-dialog
    v-model="ctxData.psFlag"
    title="导入命令参数"
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
import ModelBlockApiD07 from 'api/modelBlockD07.js'
import DeviceModelApi from 'api/deviceModel.js'
import { userStore } from 'stores/user'
const users = userStore()

const props = defineProps({
  curModelBlock: {
    type: Object,
    default: {},
  },
})
// 返回设备块
const emit = defineEmits(['changeBShowFlag'])
const toModelBlock = () => {
  emit('changeBShowFlag')
}

const regCnt = /^[0-9]*[1-9][0-9]*$/
const validateRegCnt = (rule, value, callback) => {
  if (!regCnt.test(value)) {
    callback(new Error('只能输入正整数数字！'))
  } else {
    callback()
  }
}
const validateRegAddr = (rule, value, callback) => {
  if (value !== 0 && !regCnt.test(value)) {
    callback(new Error('只能输入自然数！'))
  } else {
    callback()
  }
}
const contentRef = ref(null)
const ctxData = reactive({
  formulaForm: {},
  headerCellStyle: {
    background: variables.primaryColor,
    color: variables.fontWhiteColor,
    height: '48px',
  },
  cellStyle: {
    height: '42px',
    padding: '2px 0',
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
  typeNames: {
    t0: 'uint32',
    t1: 'int32',
    t2: 'double',
    t3: 'string',
  },
  accessModeNames: {
    am0: '只读',
    am1: '只写',
    am2: '读写',
  },

  formatOptions: [
    { label: "X.XXX", value: "X.XXX" },
    { label: "XX", value: "XX" },
    { label: "XX.XX", value: "XX.XX" },
    { label: "XX.XXXX", value: "XX.XXXX" },
    { label: "XXX.X", value: "XXX.X" },
    { label: "XXX.XXX", value: "XXX.XXX" },
    { label: "XXXXXX", value: "XXXXXX" },
    { label: "XXXXXX.XX", value: "XXXXXX.XX" },
    { label: "XXXXXXXX", value: "XXXXXXXX" },
    { label: "XXXXXXXXXXXX", value: "XXXXXXXXXXXX" },
    { label: "DDhh", value: "DDhh" },
    { label: "MMDDNN", value: "MMDDNN" },
    { label: "hhmmNN", value: "hhmmNN" },
    { label: "YYMMDDWW", value: "YYMMDDWW" },
    { label: "YYMMDDhhmm", value: "YYMMDDhhmm" },
    { label: "YYMMDDhhmmss", value: "YYMMDDhhmmss" },
    { label: "C0C1C2C3", value: "C0C1C2C3" },
  ],

  lenOptions: [
    { label: 1, value: 1 },
    { label: 2, value: 2 },
    { label: 3, value: 3 },
    { label: 4, value: 4 },
    { label: 5, value: 5 },
    { label: 6, value: 6 },
  ],

  pFlag: false, //属性对话框标识
  pTitle: '添加属性', //属性对话框titleName
  propertyForm: {
    name: '', // 属性名称，只能是字母+数字的组合，不可以是中文
    label: '', // 属性标签
    accessMode: 0, // 读写属性
    type: 0,
    unit: '', // 单位，只有uint32，int32，double有效

    rulerId:'',
    len:0,
    blockAddOffset:0,
    rulerAddOffset:0,
  },

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
  paramName: {
    name: '属性名称',
    label: '属性标签',
    accessMode: '读写模式',
    type: '属性类型',
    decimals: '小数位数',
    regCnt: '寄存器数量',
    ruleType: '解析规则',
    regAddr: '寄存器地址',
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

    rulerId: [
      {
        required: true,
        message: '标识数据长度不能为空',
        trigger: 'blur',
      },
    ],

    len: [
      {
        required: true,
        type: 'number',
        message: '数据长度只能输入数字',
        trigger: 'blur',
      },
    ],
    blockAddOffset: [
      {
        required: true,
        type: 'number',
        message: '块数据偏移地址只能输入数字',
        trigger: 'blur',
      },
    ],
    rulerAddOffset: [
      {
        required: true,
        type: 'number',
        message: '标识数据地址只能输入数字',
        trigger: 'blur',
      },
    ],

  },
  psFlag: false,
  selectedProperties: [],
  fFlag: false,
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
})
// 获取设备命令参数
const getDeviceModelBlockProperty = (flag) => {
  const pData = {
    token: users.token,
    data: {
      tslName: props.curModelBlock.tslName,
      cmdName: props.curModelBlock.name,
    },
  }
  ModelBlockApiD07.getDeviceModelBlockPropertyD07(pData).then(async (res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.propertyList = res.data == null ? [] : res.data
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
      ctxData.tableMaxHeight = contentRef.value.clientHeight - 34 - 22
    })
  })
}

getDeviceModelBlockProperty()
const refresh = () => {
  getDeviceModelBlockProperty(1)
}
const filterDMPTableData = computed(() => {
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
  ctxData.psFlag = true
}
// 导出设备属性和服务
const exportDPS = () => {
  const pData = {
    token: users.token,
    responseType: 'blob',
    data: {
      tslName: props.curModelBlock.tslName,
      cmdName: props.curModelBlock.name,
    },
  }
  ModelBlockApiD07.exportDeviceModelBlockPropertyD07(pData).then((res) => {
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
    formData.append('tslName', props.curModelBlock.tslName)
    formData.append('cmdName', props.curModelBlock.name)
    formData.append('fileName', file)
    const pData = {
      token: users.token,
      contentType: 'multipart/form-data',
      data: formData,
    }
    ModelBlockApiD07.importDeviceModelBlockPropertyD07(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success(res.message)
        ctxData.psFlag = false
        getDeviceModelBlockProperty()
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
  ctxData.pFlag = true
  ctxData.pTitle = '添加属性'
}
// 编辑属性
const editDeviceModelProperty = (row) => {
  ctxData.pFlag = true
  ctxData.pTitle = '编辑属性'
  ctxData.propertyForm.name = row.name
  ctxData.propertyForm.label = row.label
  ctxData.propertyForm.rulerId = row.rulerId
  ctxData.propertyForm.format = row.format
  ctxData.propertyForm.len = row.len
  ctxData.propertyForm.unit = row.unit
  ctxData.propertyForm.accessMode = row.accessMode
  ctxData.propertyForm.type = row.type
  ctxData.propertyForm.blockAddOffset = row.blockAddOffset
  ctxData.propertyForm.rulerAddOffset = row.rulerAddOffset
}
const propertyFormRef = ref(null)
const submitPorpertyForm = () => {
  propertyFormRef.value.validate((valid) => {
    if (valid) {
      const pData = {
        token: users.token,
        data: {
          tslName: props.curModelBlock.tslName,
          cmdName: props.curModelBlock.name,
          name: ctxData.propertyForm.name,
          label: ctxData.propertyForm.label,
          rulerId:ctxData.propertyForm.rulerId,
          format:ctxData.propertyForm.format,
          len:ctxData.propertyForm.len,
          unit:ctxData.propertyForm.unit,
          accessMode:ctxData.propertyForm.accessMode,
          type:ctxData.propertyForm.type,
          blockAddOffset:ctxData.propertyForm.blockAddOffset,
          rulerAddOffset:ctxData.propertyForm.rulerAddOffset,
        },
      }
      if (ctxData.pTitle.includes('添加')) {
        ModelBlockApiD07.addDeviceModelBlockPropertyD07(pData).then((res) => {
          handleResult(res, getDeviceModelBlockProperty)
          cancelPorpertySubmit()
        })
      } else {
        ModelBlockApiD07.editDeviceModelBlockPropertyD07(pData).then((res) => {
          handleResult(res, getDeviceModelBlockProperty)
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
  ElMessageBox.confirm('确认要删除这些属性吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          tslName: props.curModelBlock.tslName,
          cmdName: props.curModelBlock.name,
          propertyNames: pList,
        },
      }
      ModelBlockApiD07.deleteDeviceModelBlockPropertyD07(pData).then((res) => {
        handleResult(res, getDeviceModelBlockProperty)
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
    name: '', // 属性名称，只能是字母+数字的组合，不可以是中文
    label: '', // 属性标签
    accessMode: 0, // 读写属性
    type: 0,
    unit: '', // 单位，只有uint32，int32，double有效
    len:0,
    rulerId:'',
    blockAddOffset:0,
    rulerAddOffset:0,
  }
}
const getDeviceProperty = () => {
  const pData = {
    token: users.token,
    data: {
      name: props.curModelBlock.tslName,
    },
  }
  DeviceModelApi.getDeviceModelProperty(pData).then((res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.propertyInfo = [{ name: 't', label: '当前属性' }]
      res.data.forEach((item) => {
        ctxData.propertyInfo.push({
          name: item.name,
          label: item.label,
        })
      })
      //ctxData.propertyInfo = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
const showFormula = () => {
  ctxData.fFlag = true
  ctxData.curFormula = ctxData.propertyForm.formula
  getDeviceProperty()
}
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

const saveFormula = () => {
  ctxData.fFlag = false
  ctxData.propertyForm.formula = ctxData.curFormula
}
const cancelSave = () => {
  ctxData.fFlag = false
  ctxData.curFormula = ''
}
const changeRegCnt = () => {
  const val = ctxData.propertyForm.regCnt
  if (val === 1) {
    ctxData.propertyForm.ruleType = 'Int_AB'
  }
  if (val === 2) {
    ctxData.propertyForm.ruleType = 'Long_ABCD'
  }
  if (val === 4) {
    ctxData.propertyForm.ruleType = 'Double_ABCDEFGH'
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
//非第一次
watch(
  () => props.curModelBlock,
  (newVal, oldVal) => {
    getDeviceModelBlockProperty()
  }
)
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
.operation {
  display: inline-block;
  width: 25%;
  height: 46px;
}
.main-container .pagination {
  bottom: 12px;
}
</style>
