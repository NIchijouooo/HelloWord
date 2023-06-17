<template>
  <div class="main-container">
        <div class="sn-container">
          <div class="sn-title">设置网关SN</div>
          <div >
            <el-form
              :model="ctxData.snForm"
              :rules="ctxData.snRules"
              ref="snFormRef"
              status-icon
              label-position="right"
              label-width="120px"
              style="width: 50%;margin: 0 auto;padding: 40px 40px 20px 0;border: 1px solid #e4e7ed;background:rgb(255 255 255 / 60%);border-radius:4px;"
            >
              <el-form-item label="网关名称" prop="name">
                <el-input type="text" v-model="ctxData.snForm.name" autocomplete="off" placeholder="请输入网关名称！">
                </el-input>
              </el-form-item>
              <el-form-item label="SN" prop="sn">
                <el-input type="text" v-model="ctxData.snForm.sn" autocomplete="off" placeholder="请输入SN！">
                </el-input>
              </el-form-item>
              <el-form-item>
                <div style="display: flex; justify-content: center;width: 90%;">
                  <el-button type="primary" @click="setGatewaySN()" :loading="ctxData.isLoading">设置网关SN</el-button>
                  <el-button type="success" @click="refresh()">刷新</el-button>
                </div>
              </el-form-item>
            </el-form>
          </div>
        </div>
  </div>
</template>
<script setup>
import ProductApi from 'api/product.js'
import { ElMessage } from 'element-plus'
import { userStore } from 'stores/user'
const users = userStore()

const ctxData = reactive({
  typeFlag: '2',
  isLoading: false,
  snForm: {
    name: '',
    sn: '',
  },
  snRules: {
    name: [
      {
        required: true,
        message: '网关名称不能为空',
        trigger: 'blur',
      },
    ],
    sn: [
      {
        required: true,
        message: 'sn不能为空',
        trigger: 'blur',
      },
    ],
  },
})

const getGatewaySN = (flag) => {
  const pData = {
    token: users.token,
    data: {},
    mock: true,
  }
  ProductApi.getGatewaySN(pData).then((res) => {
    if (res.code === '0') {
      ctxData.snForm = res.data
      if (flag === 1) {
        ElMessage.success('刷新成功！')
      }
    } else {
      showOneResMsg(res)
    }
  })
}
getGatewaySN()
const refresh = () => {
  getGatewaySN(1)
}

const snFormRef = ref(null)
const setGatewaySN = () => {
  snFormRef.value.validate((valid) => {
    if (valid) {
      ctxData.isLoading = true
      const pData = {
        token: users.token,
        data: ctxData.snForm,
        mock: true,
      }
      ProductApi.setGatewaySN(pData).then((res) => {
        if (res.code === '0') {
          ctxData.isLoading = false
          ElMessage.success(res.message)
        } else {
          showOneResMsg(res)
        }
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
.my-card-tabs {
  position: absolute;
  top: 20px;
  left: 20px;
  right: 20px;
  bottom: 20px;
  background-color: #fff;
}
.tab-content {
  position: relative;
  margin: 10px;
}
.comItem {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  &-info {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }
  &-name {
    color: #909399;
    margin: 20px;
    color: #333;
  }
}
.sn-title {
  height: 80px;
  line-height: 80px;
  text-align: center;
  color: #303133;
  font-size: 16px;
}

.sn-container {
  position: relative;
  height: 97%;
  width: auto;
  margin-left: 20px;
  margin-top: 20px;
  border-radius: 4px;
  background-color: #f5f8fa;
}
</style>
