<template>
  <div class="main-container">
    <div class="main" style="background-color: inherit">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="box-card" :shadow="'hover'">
            <template #header>
              <div class="card-header">
                <span>本地备份</span>
              </div>
            </template>
            <div class="item">
              <el-card class="box-card" :shadow="'hover'">
                <template #header>
                  <div class="card-header">
                    <span>备份系统文件</span>
                  </div>
                </template>
                <div class="item">
                  <div>将当前系统配置导出并保存在您的计算机中。今后若需要恢复此配置，直接导入该备份文件即可。</div>
                  <br />
                  <div>提示：导出的配置文件会加密存储您的个人数据，请妥善保存。</div>
                  <br />
                  <br />
                  <br />
                  <div style="text-align: center">
                    <el-button type="primary" @click="showFileList()">
                      <el-icon class="el-input__icon"><download /></el-icon>
                      导出系统配置文件
                    </el-button>
                  </div>
                </div>
              </el-card>
              <el-card class="box-card" :shadow="'hover'">
                <template #header>
                  <div class="card-header">
                    <span>还原系统文件</span>
                  </div>
                </template>
                <div class="item">
                  <div>请选择您要恢复的配置文件。</div>
                  <br />
                  <br />
                  <br />
                  <div style="text-align: center">
                    <el-button type="primary" @click="importSysConfig()">
                      <el-icon class="el-input__icon"><upload /></el-icon>
                      导入系统配置文件
                    </el-button>
                  </div>
                </div>
              </el-card>
            </div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="box-card" :shadow="'hover'">
            <template #header>
              <div class="card-header">
                <span>远程备份</span>
              </div>
            </template>
            <div class="item">
              <el-card class="box-card" :shadow="'hover'">
                <template #header>
                  <div class="card-header">
                    <span>备份系统文件</span>
                  </div>
                </template>
                <div class="item">
                  <div>将当前系统配置导出并保存在云平台。今后若需要恢复此配置，可以从云平台导入需要的版本即可。</div>
                  <br />
                  <div>提示：云平台只保存最新的5次备份文件。</div>
                  <br />
                  <br />
                  <br />
                  <div style="text-align: center">
                    <el-button type="primary" @click="remoteBackupFile()">
                      <el-icon class="el-input__icon"><upload /></el-icon>
                      备份系统配置文件至云平台
                    </el-button>
                  </div>
                </div>
              </el-card>
              <el-card class="box-card" :shadow="'hover'">
                <template #header>
                  <div class="card-header">
                    <span>还原系统文件</span>
                  </div>
                </template>
                <div class="item">
                  <div>1、点击“获取远程备份文件”按钮，获取备份文件列表；</div>
                  <div>2、请选择您需要的版本，点击“还原系统配置文件”按钮，来恢复系统配置文件；</div>
                  <div>3、请重启网关。</div>
                  <br />
                  <br />
                  <br />
                  <div style="text-align: center">
                    <el-button type="primary" @click="getRemoteFileList()">获取远程备份文件</el-button>
                    <el-select v-model="ctxData.fileName" style="margin: 0 20px">
                      <el-option
                        v-for="item in ctxData.fileList"
                        :key="item.backupId + item.backupTime"
                        :label="item.backupId + '(' + item.backupTime + ')'"
                        :value="item.backupId"
                      />
                    </el-select>
                    <el-button type="primary" @click="getRemoteFile()">
                      <el-icon class="el-input__icon"><download /></el-icon>
                      还原系统配置文件
                    </el-button>
                  </div>
                </div>
              </el-card>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
    <!-- 导入系统配置文件 -->
    <el-dialog
      v-model="ctxData.uFlag"
      title="上传系统配置文件"
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
    <!-- 显示需要导出的内容 -->
    <el-dialog v-model="ctxData.sFlag" title="需要导出的内容" width="500px">
      <el-checkbox-group v-model="ctxData.checkList">
        <el-checkbox label="comm">通讯接口</el-checkbox>
        <el-checkbox label="coll">采集接口</el-checkbox>
        <el-checkbox label="plugin">驱动插件</el-checkbox>
        <el-checkbox label="report">上报接口</el-checkbox>
      </el-checkbox-group>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="exportSysConfig">导出</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { userStore } from 'stores/user'
import SysMaintenance from 'api/sysMaintenance.js'
import { Download, Upload } from '@element-plus/icons-vue'

const users = userStore()
const ctxData = reactive({
  uFlag: false,
  fileList: [],
  fileName: '',
  sn: '',
  sFlag: false,
  checkList: ['comm', 'coll', 'plugin', 'report'],
})
let snInfo = sessionStorage.getItem('SN')
if (snInfo !== '') {
  ctxData.sn = snInfo
  console.log('SN = ', ctxData.sn)
}
// 弹出提示框
const showFileList = () => {
  ctxData.sFlag = true
}

// 导出系统配置文件
const exportSysConfig = () => {
  const pData = {
    token: users.token,
    responseType: 'blob',
    data: {
      names: ctxData.checkList,
    },
  }
  SysMaintenance.exportSystemConfig(pData).then((res) => {
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
// 导入系统配置文件
const importSysConfig = () => {
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
  const fileTypeList = ['application/x-zip-compressed', 'application/zip']
  const typeFlag = fileTypeList.includes(file.type)
  if (!typeFlag) {
    ElMessage({
      type: 'error',
      message: '文件格式不正确,必须是zip文件！',
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
    SysMaintenance.importSystemConfig(pData).then((res) => {
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
//取消上传系统配置文件
const cancelUploadSysConfig = () => {
  ctxData.uFlag = false
  uploadRef.value.clearFiles()
}
// 弹窗右上角关闭事件处理
const beforeCloseUploadSysConfig = () => {
  cancelUploadSysConfig()
}
// 远程备份文件
const remoteBackupFile = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  SysMaintenance.backupRemote(pData).then((res) => {
    if (res.code === '0') {
      ElMessage({
        type: 'success',
        message: res.message,
      })
    } else {
      showOneResMsg(res)
    }
  })
}
// 获取远程文件列表
const getRemoteFileList = () => {
  if (ctxData.sn === '') {
    ElMessage.warning('本网关为配置SN号！')
  } else {
    const pData = {
      token: users.token,
      data: {
        sn: ctxData.sn,
      },
    }
    SysMaintenance.getRemoteFileInfo(pData).then((res) => {
      if (res.code === '0') {
        ctxData.fileList = res.data
      } else {
        showOneResMsg(res)
      }
    })
  }
}
//
const getRemoteFile = () => {
  if (ctxData.fileName !== '') {
    const pData = {
      token: users.token,
      data: {
        sn: ctxData.sn,
        backupId: ctxData.fileName,
      },
    }
    SysMaintenance.getRemoteFile(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success('获取备份文件成功！')
      } else {
        showOneResMsg(res)
      }
    })
  } else {
    ElMessage.warning('请选择备份文件的版本！')
  }
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
