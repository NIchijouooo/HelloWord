<template>
  <div class="main-container">
    <div class="main" style="background-color: inherit">
      <el-card class="box-card" :shadow="'hover'">
        <template #header>
          <div class="card-header">
            <span>在线升级</span>
          </div>
        </template>
        <div class="item">
          <div>当前系统版本：{{ ctxData.sysParams.softVer }} 。</div>
          <br />
          <div>最新系统版本：{{ ctxData.sysParams.softVer }} 。</div>
          <br />
          <br />
          <br />
          <div style="text-align: center">
            <el-button type="primary" @click="updateSysOnline()">
              <el-icon class="el-input__icon"><refresh /></el-icon>
              在线升级系统到最新版本
            </el-button>
          </div>
        </div>
      </el-card>
      <el-card class="box-card" :shadow="'hover'">
        <template #header>
          <div class="card-header">
            <span>本地升级</span>
          </div>
        </template>
        <div class="item">
          <div>请选择您要升级的系统文件。</div>
          <br />
          <br />
          <br />
          <div style="text-align: center">
            <el-button type="primary" @click="importSysConfig()">
              <el-icon class="el-input__icon"><upload /></el-icon>
              导入系统升级文件
            </el-button>
          </div>
        </div>
      </el-card>
    </div>
    <!-- 导入系统升级文件 -->
    <el-dialog
      v-model="ctxData.uFlag"
      title="上传系统升级文件"
      width="600px"
      :before-close="beforeCloseUploadSysConfig"
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
          <div class="el-upload__tip">只能上传一个文件，只支持zip格式文件！</div>
        </template>
      </el-upload>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelUploadSysConfig">取消</el-button>
          <el-button type="primary" @click="submitUploadSysConfig">上传</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { userStore } from 'stores/user'
import SysMaintenance from 'api/sysMaintenance.js'
import DashboardApi from 'api/dashboard.js'
import { Refresh, Upload } from '@element-plus/icons-vue'

const users = userStore()

const ctxData = reactive({
  uFlag: false,
  sysParams: {},
})
// 获取系统参数
const getSysParams = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  DashboardApi.getSysParams(pData).then((res) => {
    console.log('getSysParams -> res', res)
    if (res.code === '0') {
      ctxData.sysParams = res.data
      console.log('getSysParams -> ctxData.sysParams', ctxData.sysParams)
    }
  })
}
getSysParams()
// 在线升级系统
const updateSysOnline = () => {
  //
  ElMessage.info('功能升级中...')
}
// 导入系统升级文件
const importSysConfig = () => {
  console.log('importSysConfig')
  ctxData.uFlag = true
}
/**
 * 提交上传的插件
 */
const submitUploadSysConfig = () => {
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
  const fileTypeList = ['application/x-tar']
  const typeFlag = fileTypeList.includes(file.type)
  if (!typeFlag) {
    ElMessage({
      type: 'error',
      message: '文件格式不正确,必须是tar文件！',
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
    SysMaintenance.updateSystem(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success(res.message)
        ctxData.uFlag = false
      } else {
        showOneResMsg(res)
      }
    })
  }
}
//取消上传系统配置文件
const cancelUploadSysConfig = () => {
  ctxData.uFlag = false
  uploadRef.value.clearFiles()
}
// 弹窗右上角关闭事件处理
const beforeCloseUploadSysConfig = () => {
  cancelUploadSysConfig()
}

//显示单个res结果，code不等于 '0' 的message
const showOneResMsg = (res) => {
  ElMessage({
    type: 'error',
    message: res.message,
  })
}
</script>
<style lang="scss" scoped>
@use 'styles/custom-scoped.scss' as *;
.box-card {
  margin-bottom: 20px;
}
</style>
