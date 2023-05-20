<template>
  <div class="main-container">
    <el-tabs type="border-card" class="my-card-tabs" v-model="ctxData.typeFlag">
      <el-tab-pane label="设置网关SN" name="2">
        <div class="sn-container">
          <el-card class="box-card">
            <div>
              <el-form
                :model="ctxData.snForm"
                :rules="ctxData.snRules"
                ref="snFormRef"
                status-icon
                label-position="right"
                label-width="120px"
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
                  <el-button type="primary" @click="setGatewaySN()" :loading="ctxData.isLoading">设置网关SN</el-button>
                  <el-button type="success" @click="refresh()">刷新</el-button>
                </el-form-item>
              </el-form>
            </div>
          </el-card>
        </div>
      </el-tab-pane>
    </el-tabs>
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
.tips {
  background-color: #d3dce6;
  border-radius: 4px;
  line-height: 48px;
  height: 48px;
  padding: 0 20px;
  margin: 20px;
}

.sn-container {
  position: relative;
  height: 100%;
  width: 50%;
  background-color: #d3dce6;
}
:deep(.el-tabs--border-card > .el-tabs__header .el-tabs__item:not(.is-disabled):hover) {
  background-color: #3054eb;
  color: #fff;
}
:deep(.el-tabs--border-card > .el-tabs__header .el-tabs__item.is-active) {
  background-color: #3054eb;
  color: #fff;
}
</style>
