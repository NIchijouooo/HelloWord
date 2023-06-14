<template>
<div class="main-container">
  <div class="main">
  <div class="mc-left">
    <div class="title">
      <el-button type="primary" style="width: 100%" plain @click="toDeviceModel()">
        <el-icon class="el-input__icon"><back /></el-icon>
        返回设备模型
      </el-button>
    </div>
    <div class="mcl-content">
      <el-card
        v-for="(item, index) in props.deviceModelList"
        :key="'i' + index"
        shadow="hover"
        :class="{ 'mclc-card': true, activeModel: item.name == ctxData.curDeviceModel.name }"
        @click="selectModel(item)"
      >
        <div class="mclc-content">
          <span>{{ item.name }}</span>
        </div>
      </el-card>
    </div>
  </div>
  <div class="mc-right">
    <div class="mcr-top">
      <div style="display: flex; justify-content: space-between;">
        <div class="title" style="position: relative;width: 40%;justify-content: flex-start;">
          <div class="tName">{{ ctxData.curDeviceModel.label }}：命令列表</div>
        </div>
        <div class="search-bar" style="text-align:right;">
          <el-form :inline="true" ref="searchFormRef" status-icon label-width="90px">
            <el-form-item style="margin-left: 20px;">
              <el-input style="width: 200px" placeholder="请输入命令名称" v-model="ctxData.deviceModelBlock">
                <template #prefix>
                  <el-icon class="el-input__icon"><search /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" plain class="right-btn" @click="importDPS()">
                <el-icon class="el-input__icon"><download /></el-icon>
                导入命令
              </el-button>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" plain class="right-btn" @click="exportDPS()">
                <el-icon class="el-input__icon"><upload /></el-icon>
                导出命令
              </el-button>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" bg class="right-btn" @click="addDeviceModelBlock()">
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
              <el-button type="danger" bg class="right-btn" @click="deleteDeviceModelBlock()">
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
          :max-height="ctxData.tableMaxHeight"
          style="width: 100%"
          stripe
          @selection-change="handleSelectionChange"
          @row-dblclick="editDeviceModelBlock"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column sortable prop="name" label="命令名称" width="auto" min-width="150" align="center"> </el-table-column>
          <el-table-column sortable prop="label" label="命令标签" width="auto" min-width="150" align="center"> </el-table-column>
          <el-table-column sortable label="功能码" width="auto" min-width="150" align="center">
            <template #default="scope">
              {{ ctxData.funCodeNames['f' + scope.row.funCode] }}
            </template>
          </el-table-column>
          <el-table-column sortable prop="startRegAddr" label="寄存器地址" width="auto" min-width="120" align="center">
          </el-table-column>
          <el-table-column sortable prop="regCnt" label="寄存器数量" width="auto" min-width="100" align="center">
          </el-table-column>
          <el-table-column label="操作" width="auto" min-width="200" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="editDeviceModelBlock(scope.row)" text type="primary">编辑</el-button>
              <el-button @click="showModelPropertyBlock(scope.row, 2)" text type="success">命令参数详情</el-button>
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
          ></el-pagination>
        </div>
      </div>
    </div>
    <div class="mcr-bottom" v-if="ctxData.porpFlag">
      <BlockPropertyRtu :curModelBlock="ctxData.curModelBlock"> </BlockPropertyRtu>
    </div>
  </div>

  <!-- 添加编辑命令 -->
  <el-dialog
    v-model="ctxData.pFlag"
    :title="ctxData.pTitle"
    width="600px"
    :before-close="handleCloseBlock"
    :close-on-click-modal="false"
  >
    <div class="dialog-content">
      <el-form
        :model="ctxData.blockForm"
        :rules="ctxData.propertyRules"
        ref="blockFormRef"
        status-icon
        label-position="right"
        label-width="120px"
      >
        <el-form-item label="块名称" prop="name">
          <el-input
            :disabled="ctxData.pTitle.includes('编辑')"
            type="text"
            v-model="ctxData.blockForm.name"
            autocomplete="off"
            placeholder="请输入块名称"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="块标签" prop="label">
          <el-input type="text" v-model="ctxData.blockForm.label" autocomplete="off" placeholder="请输入块标签">
          </el-input>
        </el-form-item>
        <el-form-item label="功能码" prop="funCode">
          <el-select v-model="ctxData.blockForm.funCode" style="width: 100%" placeholder="请选择功能码">
            <el-option
              v-for="item in ctxData.funCodeOptions"
              :key="'funCode_' + item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="寄存器地址" prop="startRegAddr">
          <el-input
            type="text"
            v-model.number="ctxData.blockForm.startRegAddr"
            autocomplete="off"
            placeholder="请输入寄存器地址"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="寄存器数量" prop="regCnt">
          <el-input
            type="text"
            v-model.number="ctxData.blockForm.regCnt"
            autocomplete="off"
            placeholder="请输入寄存器数量"
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

  <!-- 导入命令 -->
  <el-dialog
    v-model="ctxData.psFlag"
    title="导入命令"
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
import ModelBlockApi from 'api/modelBlock.js'
import { userStore } from 'stores/user'
import BlockPropertyRtu from './BlockPropertyRtu.vue'
const users = userStore()

const props = defineProps({
  curDeviceModel: {
    type: Object,
    default: {},
  },
  deviceModelList: {
    type: Array,
    default: [],
  },
})
console.log('id -> props', props)
// 返回设备模型
const emit = defineEmits(['changeShowFlag'])
const toDeviceModel = () => {
  //lp update 2023-06-12 首页设备模型跳转进入显示命令详情后，对点击返回按钮进行操作标识，防止死循环跳转显示详情页面
  emit('changeShowFlag', 'goBack')
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
  curDeviceModel: props.curDeviceModel,
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
  deviceModelBlock: '',
  typeOptions: [
    { label: 'uint32', value: 0 },
    { label: 'int32', value: 1 },
    { label: 'double', value: 2 },
    { label: 'string', value: 3 },
  ],
  funCodeOptions: [
    { label: '01(读线圈状态)', value: 1 },
    { label: '02(读离散输入状态)', value: 2 },
    { label: '03(读保持寄存器)', value: 3 },
    { label: '04(读输入寄存器)', value: 4 },
    { label: '05(写单个线圈)', value: 5 },
    { label: '06(写单个保持寄存器)', value: 6 },
    { label: '15(写多个线圈)', value: 15 },
    { label: '16(写多个保持寄存器)', value: 16 },
  ],
  funCodeNames: {
    f1: '01(读线圈状态)',
    f2: '02(读离散输入状态)',
    f3: '03(读保持寄存器)',
    f4: '04(读输入寄存器)',
    f5: '05(写单个线圈)',
    f6: '06(写单个保持寄存器)',
    f15: '15(写多个线圈)',
    f16: '16(写多个保持寄存器)',
  },
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
  },
  pFlag: false, //属性对话框标识
  pTitle: '添加属性', //属性对话框titleName
  blockForm: {
    name: '',
    label: '',
    funCode: null, //功能码
    startRegAddr: 0, //寄存器PLC地址
    regCnt: 0, //数量
  },
  paramName: {
    name: '块名称',
    label: '块标签',
    funCode: '功能码',
    type: '属性类型',
    startRegAddr: '寄存器地址',
    regCnt: '寄存器数量',
  },
  propertyRules: {
    name: [
      {
        required: true,
        message: '块名称不能为空',
        trigger: 'blur',
      },
    ],
    label: [
      {
        required: true,
        message: '块标签不能为空',
        trigger: 'blur',
      },
    ],
    funCode: [
      {
        required: true,
        message: '功能码不能为空',
        trigger: 'blur',
      },
    ],
    startRegAddr: [
      {
        required: true,
        message: '寄存器地址不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateRegAddr,
      },
    ],
    regCnt: [
      {
        required: true,
        message: '寄存器数量不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateRegCnt,
      },
    ],
  },
  psFlag: false,
  selectedProperties: [],
  curModelBlock: {},
  porpFlag: false,
})
// 获取命令
const getDeviceModelBlock = (flag) => {
  const pData = {
    token: users.token,
    data: {
      tslName: ctxData.curDeviceModel.name,
    },
  }
  ModelBlockApi.getDeviceModelBlock(pData).then(async (res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.propertyList = res.data
      if (ctxData.propertyList.length > 0) {
        ctxData.curModelBlock = ctxData.propertyList[0]
        ctxData.curModelBlock['tslName'] = ctxData.curDeviceModel.name
        console.log('getDeviceModelBlock -> ctxData.curModelBlock', ctxData.curModelBlock)
        ctxData.porpFlag = true
      } else {
        ctxData.porpFlag = false
      }
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
      ctxData.tableMaxHeight = contentRef.value.clientHeight - 34 - 22 - 82
      if (flag === 2) showModelPropertyBlock(ctxData.curModelBlock, 1)
    })
  })
}

getDeviceModelBlock(2)
const refresh = () => {
  getDeviceModelBlock(1)
}
const selectModel = (item) => {
  if (ctxData.curDeviceModel.name == item.name) return
  ctxData.curDeviceModel = item
  getDeviceModelBlock()
}
const filterDMPTableData = computed(() => {
  console.log('ctxData.propertyList ->', ctxData.propertyList)
  return ctxData.propertyList
    .filter((item) => {
      var a = !ctxData.deviceModelBlock
      var b = item.name.toLowerCase().includes(ctxData.deviceModelBlock.toLowerCase())

      return a || b
    })
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  return ctxData.propertyList.filter((item) => {
    var a = !ctxData.deviceModelBlock
    var b = item.name.toLowerCase().includes(ctxData.deviceModelBlock.toLowerCase())

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
      name: ctxData.curDeviceModel.name,
    },
  }
  ModelBlockApi.exportDeviceModelBlock(pData).then((res) => {
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
    formData.append('name', ctxData.curDeviceModel.name)
    formData.append('fileName', file)
    const pData = {
      token: users.token,
      contentType: 'multipart/form-data',
      data: formData,
    }
    ModelBlockApi.importDeviceModelBlock(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success(res.message)
        ctxData.psFlag = false
        getDeviceModelBlock()
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
const addDeviceModelBlock = () => {
  console.log('addDeviceModelBlock')
  ctxData.pFlag = true
  ctxData.pTitle = '添加属性'
}
// 编辑属性
const editDeviceModelBlock = (row) => {
  ctxData.pFlag = true
  ctxData.pTitle = '编辑属性'
  ctxData.blockForm.name = row.name
  ctxData.blockForm.label = row.label
  ctxData.blockForm.funCode = row.funCode
  ctxData.blockForm.startRegAddr = row.startRegAddr
  ctxData.blockForm.regCnt = row.regCnt
}
const blockFormRef = ref(null)
const submitPorpertyForm = () => {
  blockFormRef.value.validate((valid) => {
    if (valid) {
      const pData = {
        token: users.token,
        data: {
          tslName: ctxData.curDeviceModel.name,
          name: ctxData.blockForm.name,
          label: ctxData.blockForm.label,
          funCode: ctxData.blockForm.funCode,
          startRegAddr: ctxData.blockForm.startRegAddr,
          regCnt: ctxData.blockForm.regCnt,
        },
      }
      if (ctxData.pTitle.includes('添加')) {
        ModelBlockApi.addDeviceModelBlock(pData).then((res) => {
          handleResult(res, getDeviceModelBlock)
          cancelPorpertySubmit()
        })
      } else {
        ModelBlockApi.editDeviceModelBlock(pData).then((res) => {
          handleResult(res, getDeviceModelBlock)
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
  blockFormRef.value.resetFields()
  initBlockForm()
}
const handleCloseBlock = (done) => {
  cancelPorpertySubmit()
}

const handleSelectionChange = (val) => {
  ctxData.selectedProperties = val
  console.log('handleSelectionChange -> val =', val)
}
// 删除命令,可以批量删除
const deleteDeviceModelBlock = (row) => {
  let pList = []
  if (ctxData.selectedProperties.length === 0) {
    ElMessage.info('请至少选择一个命令！')
    return
  } else {
    ctxData.selectedProperties.forEach((item) => {
      pList.push(item.name)
    })
  }
  console.log('pList', pList)
  ElMessageBox.confirm('确认要删除这些命令吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          tslName: ctxData.curDeviceModel.name,
          names: pList,
        },
      }
      ModelBlockApi.deleteDeviceModelBlock(pData).then((res) => {
        handleResult(res, getDeviceModelBlock)
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
const initBlockForm = () => {
  ctxData.blockForm = {
    name: '', //命令名称
    label: '', // 命令标签
    funCode: null, // 读写属性
    startRegAddr: 0,
    regCnt: 0,
  }
}

const showModelPropertyBlock = (row) => {
  ctxData.curModelBlock = row
  ctxData.curModelBlock['tslName'] = ctxData.curDeviceModel.name
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
.mc-left {
  width: 230px;
}
.mc-right {
  left: 266px;
  background-color: initial;
}
.mcr-top {
  position: absolute;
  top: 0;
  height: calc(50% - 10px);
  left: 0;
  right: 0;
  background-color: #f5f5f5;
  border-radius: 4px;
}
.mcr-bottom {
  position: absolute;
  bottom: 0;
  height: calc(50% - 10px);
  left: 0;
  right: 0;
  background-color: #f5f5f5;
  border-radius: 4px;
}
.activeModel {
  box-shadow: 0px 0px 12px rgba(48, 84, 235, 0.3);
}
.main-container .pagination {
  bottom: 12px;
}
:deep(.el-card__body) {
  cursor: pointer;
  padding: 10px;
}
.activeModel:hover {
  box-shadow: 0px 0px 12px rgba(48, 84, 235, 0.3);
}
</style>
