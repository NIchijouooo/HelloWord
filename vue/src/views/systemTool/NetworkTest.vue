<template>
  <div class="main-container">
    <div class="main" style="background-color: inherit">
      <el-row :gutter="20" style="height: 100%">
        <el-col :span="24">
          <el-card class="box-card" :shadow="'never'">
            <template #header>
              <div class="card-header">
                <span>网络诊断</span>
              </div>
            </template>
            <div class="item">
              <el-card style="height: 100%">
                <div style="padding: 14px">
                  <el-form
                    :model="ctxData.net1Form"
                    :rules="ctxData.net1Rules"
                    ref="net1Ref"
                    status-icon
                    label-width="120px"
                  >
                    <el-form-item label="IP地址" prop="ip">
                      <el-input v-model="ctxData.net1Form.ip" placeholder="请输入ip地址" />
                    </el-form-item>
                    <el-form-item>
                      <el-button type="primary" @click="checkIP" :loading="ctxData.loadingFlag">检测</el-button>
                    </el-form-item>
                    <el-form-item label="返回结果">
                      <el-input rows="10" type="textarea" v-model="ctxData.resultData" />
                    </el-form-item>
                  </el-form>
                </div>
              </el-card>
              <div class="remark">
                <el-row :gutter="16">
                  <el-col :span="1">操作步骤:</el-col>
                  <el-col :span="21">1、输入需要ping设备的ip地址</el-col>
                </el-row>
                <el-row :gutter="16">
                  <el-col :offset="1" :span="21">2、点击检测按钮</el-col>
                </el-row>
                <el-row :gutter="16">
                  <el-col :offset="1" :span="21">3、返回结果框返回ping命令的结果</el-col>
                </el-row>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>
<script setup>
import { userStore } from 'stores/user'
import SysToolApi from 'api/sysTool.js'
const refExpIP =
  /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/
const refExpYM =
  /^(?=^.{3,255}$)(http(s)?:\/\/)?(www\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:\d+)*(\/\w+\.\w+)*$/
const validateIP = (rule, value, callback) => {
  console.log('validateIP')
  if (value !== '') {
    if (refExpIP.test(value) || refExpYM.test(value)) {
      callback()
    } else {
      callback(new Error('格式错误！'))
    }
  } else {
    callback()
  }
}
const users = userStore()
const ctxData = reactive({
  loadingFlag: false,
  net1Form: {
    ip: '',
  },
  net1Rules: {
    ip: [
      {
        validator: validateIP,
        trigger: 'blur',
      },
    ],
  },
  resultData: '',
})
const net1Ref = ref(null)
const checkIP = () => {
  net1Ref.value.validate((valid) => {
    if (valid) {
      ctxData.loadingFlag = true
      const pData = {
        token: users.token,
        data: {
          ip: ctxData.net1Form.ip,
        },
        mock: true,
      }
      SysToolApi.sendPingCmd(pData).then((res) => {
        if (res && res.code === '0') {
          ctxData.resultData = res.data
        } else {
          showOneResMsg(res)
        }
        ctxData.loadingFlag = false
      })
    } else {
      return false
    }
  })
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
  padding-bottom: 20px;
  padding-right: 20px;
}
.remark {
  margin-top: 16px;
  font-size: 12px;
  color: #f56c6c;
}
.item {
  position: relative;
  height: 700px;
  width: 100%;
  min-width: 540px;
  margin-bottom: 50px;
}
:deep(.el-card.is-always-shadow) {
  box-shadow: none;
}
</style>
