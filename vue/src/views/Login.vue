<template>
  <div class="login">
    <el-card class="login_center">
      <template #header>
        <div class="card_header">
          <span>{{ ctxData.configInfo.name }}登录</span>
        </div>
      </template>
      <el-form :model="ctxData.loginForm" :rules="ctxData.loginRules" ref="loginFormRef">
        <el-form-item prop="userName">
          <el-input
            v-model.trim="ctxData.loginForm.userName"
            maxlength="32"
            placeholder="请输入账号"
            clearable
            style="height: 50px"
          >
            <template #prefix>
              <el-icon :size="20" class="el-input__icon"><user /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model.trim="ctxData.loginForm.password"
            maxlength="16"
            show-password
            placeholder="请输入密码"
            clearable
            style="height: 50px"
            @keyup.enter.exact="handleLogin"
          >
            <template #prefix>
              <el-icon :size="20" class="el-input__icon"><lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            bg
            style="width: 100%; margin-top: 20px; font-size: 16px"
            :loading="ctxData.loading"
            @click="handleLogin"
            >登 录</el-button
          >
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>
<script setup>
import LoginApi from 'api/login.js'
import { User, Lock } from '@element-plus/icons-vue'
import router from '@/router'
import { userStore } from 'stores/user.js'
import { ElMessage } from 'element-plus'
import 'element-plus/es/components/message/style/css'
import { configStore } from '@/stores/app.js'

const config = configStore()
const users = userStore()
const logout = localStorage.getItem('logout')

if (logout != null && logout != '') {
  try {
    const lo = JSON.parse(logout)
    ElMessage({
      type: lo.type,
      message: lo.message,
    })
  } catch (error) {
    ElMessage({
      type: 'error',
      message: error,
    })
  }
  localStorage.setItem('logout', '')
} else {
  users.setLoginInfo(null, [], '')
  localStorage.setItem('loginInfo', '')
}

const ctxData = reactive({
  loginForm: {
    userName: '',
    password: '',
    roleType: '',
  },
  loginRules: {
    userName: [{ required: true, message: '账号不能为空', trigger: 'blur' }],
    password: [
      { required: true, message: '密码不能为空', trigger: 'blur' },
      { min: 4, max: 16, message: '密码长度为4-16位', trigger: 'blur' },
    ],
  },
  loading: false,
  configInfo: config.configInfo,
})

const loginFormRef = ref()
const handleLogin = () => {
  console.log('handleLogin')
  loginFormRef.value.validate((valid) => {
    if (valid) {
      ctxData.loading = true

      // 登录请求
      const pData = {
        loginFlag: true,
        data: {
          username: ctxData.loginForm.userName,
          password: ctxData.loginForm.password,
        },
        token: '',
      }
      LoginApi.login(pData).then((res) => {
        console.log('LoginApi.login -> res', res)
        ctxData.loading = false
        if (res.code === '0') {
          const userInfo = {
            userName: res.data.username,
          }
          const roots = res.data.permissions
          const token = res.data.token
          users.setLoginInfo(userInfo, roots, token)
          // 登录成功后将信息存入本地缓存
          localStorage.setItem('loginInfo', JSON.stringify(res.data))
          router.push('/dashboard')
        } else {
          showOneResMsg(res)
        }
      })
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
.login-partic {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
}
.login {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
  background: url(assets/images/login-bg.png);
  background-repeat: no-repeat;
  background-size: 100% 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  .login_center {
    width: 490px;
    height: auto;
    background-color: rgba(0,0,0,.5);
    box-shadow: none;
    border: 0;
  }

  .card_header {
    font-size: 20px;
    text-align: left;
    color: #f0f7ff;
    letter-spacing: 1px;
    text-align: center;
  }
  .el-input__icon {
    width: 40px;
  }
}
:deep(.el-card__header) {
  border-bottom: 0;
  padding: 48px 46px 16px 46px;
}
:deep(.el-card__body) {
  padding: 16px 46px;
}
:deep(.el-form-item) {
  margin-bottom: 20px;
}
:deep(.el-input__inner) {
  height: 50px;
  line-height: 50px;
  font-size: 16px;
  background-color: #f5f8fa !important;
}
:deep(.el-input__prefix) {
  left: 0;
}
:deep(.el-button) {
  height: 50px;
}
:deep(.el-input__wrapper) {
  background-color: #f5f8fa !important;
}
:deep(.el-input__prefix) {
  background-color: #f5f8fa !important;
}
</style>
