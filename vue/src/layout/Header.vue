<template>
  <el-header height="66px">
    <div class="header_right">
      <span>
        {{ ctxData.configInfo.name }}
        <b>{{ ctxData.configInfo.version }}</b>
      </span>
      <div style="display: flex; align-items: center">
        <el-tooltip :content="'CPU使用率：' + ctxData.sysParams.cpuUse + '%'">
          <el-image @click="changeIndex(0)" class="header-item" :src="dbIcon00" fit="cover" />
        </el-tooltip>
        <el-tooltip :content="'内存使用率：'+ctxData.sysParams.memUse + '%'">
          <el-image @click="changeIndex(1)" class="header-item" :src="dbIcon04" fit="cover" />
        </el-tooltip>
        <el-tooltip :content="'硬盘使用率：'+ctxData.sysParams.diskUse + '%'">
          <el-image @click="changeIndex(2)" class="header-item" :src="dbIcon03" fit="cover" />
        </el-tooltip>
        <el-tooltip :content="'设备在线率：'+ctxData.sysParams.deviceOnline + '%'">
          <el-image @click="changeIndex(3)" class="header-item" :src="dbIcon01" fit="cover" />
        </el-tooltip>
        <el-tooltip :content="'通讯丢包率：'+ctxData.sysParams.devicePacketLoss + '%'">
          <el-image @click="changeIndex(4)" class="header-item" :src="dbIcon02" fit="cover" />
        </el-tooltip>
        <el-tooltip :content="ctxData.isFullScreen ? '退出全屏' : '全屏'">
          <el-icon @click="handleFullScreen">
            <full-screen></full-screen>
          </el-icon>
        </el-tooltip>
        <el-dropdown size="medium" @command="handleCommand">
          <div class="user_info">
            <el-avatar :size="40">
              <span class="username">{{ ctxData.userName }}</span>
            </el-avatar>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="user">
                <template #default>
                  <el-icon :size="14"><user /></el-icon>
                  修改密码
                </template>
              </el-dropdown-item>
              <el-dropdown-item command="logout">
                <template #default>
                  <el-icon :size="14"><circle-close /></el-icon>
                  退出登录
                </template>
              </el-dropdown-item>
              <el-dropdown-item command="restartServer">
                <template #default>
                  <el-icon :size="14"><switch-button /></el-icon>
                  重启服务
                </template>
              </el-dropdown-item>
              <el-dropdown-item command="restartSystem">
                <template #default>
                  <el-icon :size="14"><switch-button /></el-icon>
                  重启系统
                </template>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
    <!-- 更改密码 -->
    <el-dialog
      title="更改密码"
      v-model="ctxData.pFlag"
      width="500px"
      :before-close="handleClose"
      :close-on-click-modal="false"
    >
      <div class="dialog-updatePwd">
        <el-form
          :model="ctxData.pwdForm"
          status-icon
          :rules="ctxData.pwdRules"
          ref="pwdFormRef"
          label-position="right"
          label-width="100px"
          style="margin: 0 12px"
        >
          <el-form-item label="用户名" prop="role">
            <el-input v-model="ctxData.pwdForm.role" :disabled="true"></el-input>
          </el-form-item>
          <el-form-item label="旧的密码" prop="oldPwd">
            <el-input
              type="password"
              v-model="ctxData.pwdForm.oldPwd"
              autocomplete="off"
              placeholder="请输入旧密码"
            ></el-input>
          </el-form-item>
          <el-form-item label="新的密码" prop="newPwd">
            <el-input
              type="password"
              v-model="ctxData.pwdForm.newPwd"
              autocomplete="off"
              placeholder="请输入新密码"
            ></el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="submitPwdForm()">保存</el-button>
          <el-button @click="handleClose()">取消</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- 查看使用率 -->
    <el-dialog 
      v-model="ctxData.bFlag" 
      :title="ctxData.headItemName['hin' + ctxData.activeIndex]" 
      :before-close="handleFlagClose" 
      :close-on-click-modal="false"
      width="50%"
    >
      <div style="width: 100%;height: calc(410px);">
        <line-chart :chart-data="ctxData.curChartData" :key="ctxData.curChartData" style="width: 100%;height: calc(410px);"></line-chart>
      </div>
    </el-dialog>
  </el-header>
</template>

<script setup>
import { userStore } from '@/stores/user.js'
import { useRouter } from 'vue-router'
import { configStore } from '@/stores/app.js'
import { User, CircleClose, SwitchButton } from '@element-plus/icons-vue'
import screenfull from 'screenfull'
import { FullScreen } from '@element-plus/icons-vue'
import setLoginInfo from 'utils/setLoginInfo.js'
import { reactive } from 'vue'
import LoginApi from 'api/login.js'
import SystemApi from 'api/sysMaintenance.js'
import LineChart from 'comps/LineChart.vue'
import DashboardApi from 'api/dashboard.js'
import { ElMessage, ElLoading } from 'element-plus'
import dbIcon00 from '@/assets/images/icon/db-icon00.png'
import dbIcon01 from '@/assets/images/icon/db-icon01.png'
import dbIcon02 from '@/assets/images/icon/db-icon02.png'
import dbIcon03 from '@/assets/images/icon/db-icon03.png'
import dbIcon04 from '@/assets/images/icon/db-icon04.png'

const router = useRouter()
const users = userStore()
const config = configStore()
const validatePass1 = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('旧的密码不能为空！'))
  } else {
    if (value.length <= 3) {
      callback(new Error('旧的密码长度不小于3！'))
    } else {
      callback()
    }
  }
}
const validatePass2 = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('新的密码不能为空！'))
  } else {
    if (value.length <= 3) {
      callback(new Error('新的密码长度不小于3！'))
    } else {
      callback()
    }
  }
}
const ctxData = reactive({
  activeIndex: 0,
  headItemName: {
    hin0: 'CPU占用率',
    hin1: '内存使用率',
    hin2: '硬盘使用率',
    hin3: '设备在线率',
    hin4: '通讯丢包率',
  },
  curChartData: [],
  bFlag: false,
  sysParams: {
    cpuUse: '',
    memUse: '',
    diskUse: '',
    deviceOnline: '',
    devicePacketLoss: '',
  },

  configInfo: config.configInfo,
  isFullScreen: false,
  userName: '',
  pFlag: false,
  pwdForm: {
    role: '',
    oldPwd: '',
    newPwd: '',
  },
  pwdRules: {
    role: [
      {
        required: true,
        message: '用户名不能为空',
        trigger: 'blur',
      },
    ],
    oldPwd: [
      {
        required: true,
        validator: validatePass1,
        trigger: 'blur',
      },
    ],
    newPwd: [
      {
        required: true,
        validator: validatePass2,
        trigger: 'blur',
      },
    ],
  },
})
if (!users.userInfo) {
  setLoginInfo(users)
} else {
  ctxData.userName = users.userInfo.userName
}

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
    }
  })
}
getSysParams()
const changeIndex = (indexValue) => {
  ctxData.activeIndex = indexValue
  const legend = ctxData.headItemName['hin' + indexValue]
  const pData = {
    token: users.token,
    data: {},
  }
  if (indexValue === 0) {
    DashboardApi.getCpuList(pData).then((res) => {
      handleDatas(res, legend)
    })
  } else if (indexValue === 1) {
    DashboardApi.getMemoryList(pData).then((res) => {
      handleDatas(res, legend)
    })
  } else if (indexValue === 2) {
    DashboardApi.getDiskList(pData).then((res) => {
      handleDatas(res, legend)
    })
  } else if (indexValue === 3) {
    DashboardApi.getDeviceOnlineList(pData).then((res) => {
      handleDatas(res, legend)
    })
  } else if (indexValue === 4) {
    DashboardApi.getDevicePacketLossList(pData).then((res) => {
      handleDatas(res, legend)
    })
  }
}
// 处理请求返回的数据
const handleDatas = (res, legend) => {
  if (res.code === '0') {
    const dataPoint = res.data
    var data = []
    var time = []
    for (var i = 0; i < dataPoint.length; i++) {
      const currentPoint = dataPoint[i]
      data.push(parseFloat(currentPoint.value))
      time.push(currentPoint.time)
    }
    nextTick(() => {})
    ctxData.curChartData = { data, time, legend }
    ctxData.bFlag = true
  }
}

const cancelDialog = () => {
  ctxData.curChartData = []
  ctxData.bFlag = false
}
//处理弹出框右上角关闭图标
const handleFlagClose = (done) => {
  cancelDialog()
}

const handleCommand = (command) => {
  if (command === 'user') {
    ctxData.pwdForm = {
      role: users.userInfo.userName,
      oldPwd: '',
      newPwd: '',
    }
    ctxData.pFlag = true
  } else if (command === 'logout') {
    users.logout()
    router.replace('/login')
  } else if (command === 'restartSystem') {
    const pData = {
      token: users.token,
      data: {},
    }
    SystemApi.rebootSystem(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success('重启系统命令下发成功！')
        // rebootTip()
      }
    })
    router.replace('/login')
  } else if (command === 'restartServer') {
    const pData = {
      token: users.token,
      data: {},
    }
    SystemApi.rebootServer(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success('重启服务命令下发成功！')
        // rebootTip()
      }
    })
    router.replace('/login')
  }
}
// const rebootTip = () => {
//   const loading = ElLoading.service({
//     lock: true,
//     text: '网关正在重启...',
//     background: 'rgba(0, 0, 0, 0.7)',
//   })
//   setTimeout(() => {
//     loading.close()
//   }, 3000)
// }
const handleFullScreen = () => {
  if (screenfull.isEnabled) {
    ctxData.isFullScreen = !ctxData.isFullScreen
    screenfull.toggle()
  }
}
const pwdFormRef = ref(null)
const submitPwdForm = () => {
  pwdFormRef.value.validate((valid) => {
    if (valid) {
      const pData = {
        token: users.token,
        data: ctxData.pwdForm,
      }
      LoginApi.updatePassword(pData).then((res) => {
        console.log('updatePassword -> res ', res)
        if (res.code === '0') {
          users.updatePwd()
          router.replace('/login')
        } else {
          showOneResMsg(res)
        }
      })
    } else {
      return false
    }
  })
}
const cancelSubmit = () => {
  ctxData.pFlag = false
  pwdFormRef.value.resetFields()
}
//处理弹出框右上角关闭图标
const handleClose = (done) => {
  cancelSubmit()
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
.el-header {
  padding: 0 16px;
  border-bottom: 1px solid #ddd;
  display: flex;
  display: -webkit-flex;
  align-items: center;
  justify-content: flex-end;
  overflow: hidden;
  .header_left {
    flex: 1;
    font-size: 24px;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  .header-item {
    width: 20px; 
    height: 20px;
    margin-right: 20px;
    cursor: pointer;

    &:hover {
      color: #409EFF;
    }
  }

  .header_right {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-left: 8px;
    width: 100%;
    & > div i {
      padding: 8px 6px;
      font-size: 20px;
      cursor: pointer;
    }

    .user_info {
      margin: 0 8px;
      width: auto;
      text-align: right;
      cursor: pointer;
    }
    .el-avatar {
      vertical-align: middle;
    }
    span.username {
      vertical-align: middle;
      font-size: 12px;
      &:hover {
        color: #3054eb;
      }
    }
  }
}
</style>
