<template>
  <el-header height="56px">
    <div class="header_right">
      <span>
        {{ ctxData.configInfo.name }}
        <b>{{ ctxData.configInfo.version }}</b>
      </span>
      <div style="display: flex; align-items: center">
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
          <el-button @click="closetPwdFormDialog()">取消</el-button>
        </span>
      </template>
    </el-dialog>
  </el-header>
</template>

<script setup>
import { userStore } from '@/stores/user.js'
import { configStore } from '@/stores/app.js'
import { User, CircleClose, SwitchButton } from '@element-plus/icons-vue'
import screenfull from 'screenfull'
import { FullScreen } from '@element-plus/icons-vue'
import setLoginInfo from 'utils/setLoginInfo.js'
import { reactive } from 'vue'
import LoginApi from 'api/login.js'
import SystemApi from 'api/sysMaintenance.js'
import { ElMessage, ElLoading } from 'element-plus'
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
