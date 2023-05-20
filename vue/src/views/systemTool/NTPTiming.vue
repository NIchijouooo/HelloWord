<template>
  <div class="main-container">
    <div class="main" style="background-color: inherit">
      <el-row :gutter="20" style="height: 100%">
        <el-col :span="12">
          <el-card class="box-card" :shadow="'hover'">
            <template #header>
              <div class="card-header">
                <span>NTP校时</span>
              </div>
            </template>
            <div class="item">
              <el-card style="height: 100%">
                <div style="padding-right: 14px; height: 100%; overflow: auto">
                  <el-form
                    :model="ctxData.ntpForm"
                    :rules="ctxData.ntpRules"
                    ref="ntpRef"
                    status-icon
                    label-position="right"
                    label-width="120px"
                  >
                    <el-form-item label="启用状态">
                      <el-switch v-model="ctxData.ntpForm.enable" inline-prompt active-text="是" inactive-text="否" />
                    </el-form-item>
                    <el-form-item label="时区">
                      <el-select v-model.number="ctxData.ntpForm.timeZone" style="width: 100%" placeholder="请选时区">
                        <el-option
                          v-for="item in ctxData.timeZoneOptions"
                          :key="'type_' + item.value"
                          :label="item.label"
                          :value="item.value"
                        >
                        </el-option>
                      </el-select>
                    </el-form-item>
                    <el-form-item label="主服务器地址" prop="urlMaster">
                      <el-input
                        type="text"
                        v-model="ctxData.ntpForm.urlMaster"
                        autocomplete="off"
                        placeholder="请输入主服务器地址"
                      >
                      </el-input>
                    </el-form-item>
                    <el-form-item label="主服务器端口" prop="portMaster">
                      <el-input
                        type="text"
                        v-model.number="ctxData.ntpForm.portMaster"
                        autocomplete="off"
                        placeholder="请输入主服务器端口"
                      >
                      </el-input>
                    </el-form-item>
                    <el-form-item label="次服务器地址" prop="urlSlave">
                      <el-input
                        type="text"
                        v-model="ctxData.ntpForm.urlSlave"
                        autocomplete="off"
                        placeholder="请输入次服务器地址"
                      >
                      </el-input>
                    </el-form-item>
                    <el-form-item label="次服务器端口" prop="portSlave">
                      <el-input
                        type="text"
                        v-model.number="ctxData.ntpForm.portSlave"
                        autocomplete="off"
                        placeholder="请输入次服务器端口"
                      >
                      </el-input>
                    </el-form-item>
                    <el-form-item>
                      <el-button :disabled="ctxData.showBtn" type="primary" @click="submitNtpForm()">保存</el-button>
                      <el-button type="primary" @click="sendTimingCmd()">立即校时</el-button>
                    </el-form-item>
                  </el-form>
                </div>
              </el-card>

              <div class="remark">
                <el-row :gutter="16">
                  <el-col :span="3">操作步骤:</el-col>
                  <el-col :span="21">1、选择时区（默认是UTC+8时区）</el-col>
                </el-row>
                <el-row :gutter="16">
                  <el-col :offset="3" :span="21">2、输入主服务器地址、端口</el-col>
                </el-row>
                <el-row :gutter="16">
                  <el-col :offset="3" :span="21">3、输入次服务器地址、端口</el-col>
                </el-row>
                <el-row :gutter="16">
                  <el-col :offset="3" :span="21">4、保存配置，重启后生肖</el-col>
                </el-row>
                <el-row :gutter="16">
                  <el-col :offset="3" :span="21">5、立即校时，系统时间立即更新</el-col>
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

const users = userStore()
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
const ctxData = reactive({
  showBtn: false,
  timeZoneOptions: [
    {
      value: 'UTC+8',
      label: 'UTC+8',
    },
    {
      value: 'UTC+7',
      label: 'UTC+7',
    },
    {
      value: 'UTC+6',
      label: 'UTC+6',
    },
    {
      value: 'UTC+5',
      label: 'UTC+5',
    },
    {
      value: 'UTC+4',
      label: 'UTC+4',
    },
    {
      value: 'UTC+3',
      label: 'UTC+3',
    },
    {
      value: 'UTC+2',
      label: 'UTC+2',
    },
    {
      value: 'UTC+1',
      label: 'UTC+1',
    },
    {
      value: 'UTC+0',
      label: 'UTC+0',
    },
    {
      value: 'UTC-1',
      label: 'UTC-1',
    },
    {
      value: 'UTC-2',
      label: 'UTC-2',
    },
    {
      value: 'UTC-3',
      label: 'UTC-3',
    },
    {
      value: 'UTC-4',
      label: 'UTC-4',
    },
    {
      value: 'UTC-5',
      label: 'UTC-5',
    },
    {
      value: 'UTC-6',
      label: 'UTC-6',
    },
    {
      value: 'UTC-7',
      label: 'UTC-7',
    },
    {
      value: 'UTC-8',
      label: 'UTC-8',
    },
  ],
  ntpForm: {
    enable: '',
    timeZone: '',
    urlMaster: '',
    portMaster: null,
    urlSlave: '',
    portSlave: null,
  },
  ntpRules: {
    urlMaster: [
      {
        validator: validateIP,
        trigger: 'blur',
      },
    ],
    portMaster: [
      {
        type: 'number',
        message: '必须是数字',
        trigger: 'blur',
      },
    ],
    urlSlave: [
      {
        validator: validateIP,
        trigger: 'blur',
      },
    ],
    portSlave: [
      {
        type: 'number',
        message: '必须是数字',
        trigger: 'blur',
      },
    ],
  },
})
//获取NTP信息
const getNTPInfo = () => {
  const pData = {
    token: users.token,
    data: {},
    mock: true,
  }
  SysToolApi.getNTPInfo(pData).then((res) => {
    if (res.code === '0') {
      ctxData.ntpForm = res.data
      ElMessage({
        type: 'success',
        message: '获取NTP服务成功！',
      })
      ctxData.showBtn = false
    } else {
      showOneResMsg(res)
    }
  })
}
getNTPInfo()
// 立即校时
const sendTimingCmd = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  SysToolApi.sendTimingCmd(pData).then((res) => {
    handleResult(res, getNTPInfo)
  })
}
const ntpRef = ref(null)
const submitNtpForm = () => {
  ntpRef.value.validate((valid) => {
    if (valid) {
      ctxData.showBtn = true
      const pData = {
        token: users.token,
        data: ctxData.ntpForm,
        mock: true,
      }
      SysToolApi.updateNTP(pData).then((res) => {
        handleResult(res, getNTPInfo)
      })
    } else {
      return false
    }
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
.remark {
  margin-top: 16px;
  font-size: 12px;
  color: #f56c6c;
}
</style>
