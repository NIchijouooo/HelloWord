<template>
  <div class="main-container">
    <div class="main">
      <div class="tool-bar">
        <el-button type="primary" bg class="right-btn" @click="addNetwork()">
          <el-icon class="btn-icon">
            <Icon name="local-add" size="14px" color="#ffffff" />
          </el-icon>
          添加
        </el-button>
        <el-button style="color: #fff; margin-right: 20px" color="#2EA554" class="right-btn" @click="refresh()">
          <el-icon class="btn-icon">
            <Icon name="local-refresh" size="14px" color="#ffffff" />
          </el-icon>
          刷新
        </el-button>
      </div>
      <div class="content" ref="contentRef">
        <el-table
          :data="ctxData.networkTableData"
          :cell-style="ctxData.cellStyle"
          :header-cell-style="ctxData.headerCellStyle"
          :max-height="ctxData.tableMaxHeight"
          style="width: 100%"
          stripe
        >
          <el-table-column sortable prop="name" label="网卡名称" width="auto" min-width="180" align="center"> </el-table-column>
          <el-table-column sortable label="配置状态" width="auto" min-width="120" align="center">
            <template #default="scope">
              <el-tag :type="scope.row.configParam.configEnable ? 'success' : 'warning'">{{
                scope.row.configParam.configEnable ? '已配置' : '未配置'
              }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column sortable prop="mtu" label="MTU" width="auto" min-width="150" align="center"> </el-table-column>
          <el-table-column sortable prop="mac" label="MAC地址" width="auto" min-width="200" align="center"> </el-table-column>
          <el-table-column sortable prop="flags" label="网卡标志" width="auto" min-width="200" align="center"> </el-table-column>
          <el-table-column sortable prop="ip" label="网络地址" width="auto" min-width="150" align="center"> </el-table-column>
          <el-table-column sortable prop="netmask" label="子网掩码" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column sortable prop="gateway" label="默认网关" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column label="操作" width="auto" min-width="250" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="editNetwork(scope.row)" text type="primary">偏好设置</el-button>
              <el-button @click="deleteNetwork(scope.row)" text type="danger">删除</el-button>
            </template>
          </el-table-column>
          <template #empty>
            <div>无数据</div>
          </template>
        </el-table>
        <!-- <div class="pagination">
          <el-pagination
            :current-page="ctxData.currentPage"
            :page-size="ctxData.pagesize"
            :page-sizes="[20, 50, 200, 500]"
            :total="ctxData.networkTableData.length"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            background
            layout="total, sizes, prev, pager, next, jumper"
            style="margin-top: 46px"
          ></el-pagination>
        </div>
         -->
        <div class="mobile-tips">
          <el-tag class="ml-2" type="danger">注：网口设置配置完成后，必须重启网关，才能生效！</el-tag>
        </div>
      </div>
    </div>
    <el-dialog
      v-model="ctxData.nFlag"
      :title="ctxData.nTitle"
      width="600px"
      :before-close="handleClose"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.configForm"
          :rules="ctxData.configRules"
          ref="configFormRef"
          status-icon
          label-position="right"
          label-width="100px"
        >
          <el-form-item label="网卡名称" prop="name">
            <el-select
              :disabled="ctxData.nTitle.includes('编辑')"
              v-model="ctxData.configForm.name"
              style="width: 100%"
              placeholder="请选择网卡名称"
            >
              <el-option v-for="(item, key) of ctxData.networkNames" :key="'net_' + key" :label="item" :value="item">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="启用配置" prop="configEnable">
            <el-switch v-model="ctxData.configForm.configEnable" />
          </el-form-item>
          <el-form-item v-if="ctxData.configForm.configEnable" label="获取方式" prop="dhcpEnable">
            <el-radio-group v-model="ctxData.configForm.dhcpEnable">
              <el-radio-button :label="true">自动获取</el-radio-button>
              <el-radio-button :label="false">手动设置</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item
            v-if="ctxData.configForm.configEnable"
            label="IP地址"
            :prop="ctxData.configForm.dhcpEnable ? '' : 'configIP'"
          >
            <el-input
              :disabled="ctxData.configForm.dhcpEnable"
              type="text"
              v-model="ctxData.configForm.configIP"
              autocomplete="off"
              placeholder="请输入IP地址"
            >
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="ctxData.configForm.configEnable"
            label="子网掩码"
            :prop="ctxData.configForm.dhcpEnable ? '' : 'configNetmask'"
          >
            <el-input
              :disabled="ctxData.configForm.dhcpEnable"
              type="text"
              v-model="ctxData.configForm.configNetmask"
              autocomplete="off"
              placeholder="请输入子网掩码"
            >
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="ctxData.configForm.configEnable"
            label="默认网关"
            :prop="!ctxData.configForm.dhcpEnable && ctxData.configForm.configGateway !== '' ? 'configGateway' : ''"
          >
            <el-input
              :disabled="ctxData.configForm.dhcpEnable"
              type="text"
              v-model="ctxData.configForm.configGateway"
              autocomplete="off"
              placeholder="请输入默认网关"
            >
            </el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelSubmit()">取消</el-button>
          <el-button type="primary" @click="submitNetworkForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import variables from 'styles/variables.module.scss'
import NetworkApi from 'api/network.js'
import { userStore } from 'stores/user'
const users = userStore()

const refExpIP = /^((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))$/
const validateIP = (rule, value, callback) => {
  console.log('validateIP')
  if (!refExpIP.test(value)) {
    callback(new Error('格式错误！'))
  } else {
    callback()
  }
}

const contentRef = ref(null)
const ctxData = reactive({
  headerCellStyle: {
    background: variables.primaryColor,
    color: variables.fontWhiteColor,
    height: '54px',
  },
  cellStyle: {
    height: '48px',
  },
  tableMaxHeight: 0,
  currentPage: 1, // 默认当前页是第一页
  pagesize: 20, // 每页数据个数
  networkTableData: [],
  nFlag: false,
  nTitle: '添加网卡',
  networkNames: [],
  configForm: {
    name: '',
    configEnable: false,
    dhcpEnable: false,
    configIP: '',
    configNetmask: '',
    configGateway: '',
  },
  configRules: {
    name: [
      {
        required: true,
        message: '网卡名称不能为空',
        trigger: 'blur',
      },
    ],
    configEnable: [
      {
        required: true,
        message: '启用配置不能为空',
        trigger: 'blur',
      },
    ],
    dhcpEnable: [
      {
        required: true,
        message: '获取方式不能为空',
        trigger: 'blur',
      },
    ],
    configIP: [
      {
        required: true,
        message: 'IP地址不能为空',
        trigger: 'blur',
      },
      {
        validator: validateIP,
        trigger: 'blur',
      },
    ],
    configNetmask: [
      {
        required: true,
        message: '子网掩码不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateIP,
      },
    ],
    configGateway: [
      {
        trigger: 'blur',
        validator: validateIP,
      },
    ],
  },
})
// 获取所有网卡信息
const getNetworkList = (flag) => {
  const pData = {
    token: users.token,
    data: {},
  }
  NetworkApi.getNetworkList(pData).then(async (res) => {
    console.log('getNetworkList -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.networkTableData = res.data
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
      ctxData.tableMaxHeight = contentRef.value.clientHeight - 34 - 36 - 22
    })
  })
}
getNetworkList()
//获取所有网卡名称
const getNetworkNames = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  NetworkApi.getNetworkNames(pData).then((res) => {
    if (res.code === '0') {
      ctxData.networkNames = res.data
      console.log('ctxData.networkNames -> ', ctxData.networkNames)
    } else {
      showOneResMsg(res)
    }
  })
}
getNetworkNames()
// 刷新
const refresh = () => {
  getNetworkList(1)
}
//处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
//处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
// 添加网口
const addNetwork = () => {
  ctxData.nFlag = true
  ctxData.nTitle = '添加网卡'
}
// 编辑网卡
const editNetwork = (row) => {
  //
  ctxData.nFlag = true
  ctxData.nTitle = '编辑网卡'
  ctxData.configForm = {
    name: row.configParam.name,
    configEnable: row.configParam.configEnable,
    dhcpEnable: row.configParam.dhcpEnable,
    configIP: row.configParam.configIP,
    configNetmask: row.configParam.configNetmask,
    configGateway: row.configParam.configGateway,
  }
}
// 初始化网卡表单
const initNetworkForm = () => {
  ctxData.configForm = {
    name: '',
    configEnable: false,
    dhcpEnable: false,
    configIP: '',
    configNetmask: '',
    configGateway: '',
  }
}
const configFormRef = ref(null)
// 提交网卡表单
const submitNetworkForm = () => {
  configFormRef.value.validate((valid) => {
    if (valid) {
      const pData = {
        token: users.token,
        data: ctxData.configForm,
      }
      if (ctxData.nTitle.includes('添加')) {
        NetworkApi.addNetwork(pData).then((res) => {
          handleResult(res, getNetworkList)
          cancelSubmit()
        })
      } else {
        NetworkApi.editNetwork(pData).then((res) => {
          handleResult(res, getNetworkList)
          cancelSubmit()
        })
      }
    } else {
      return false
    }
  })
}
// 取消提交表单
const cancelSubmit = () => {
  ctxData.nFlag = false
  configFormRef.value.resetFields()
  initNetworkForm()
}
// 删除网卡
const deleteNetwork = (row) => {
  //
  ElMessageBox.confirm('确定要删除这个网卡配置吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          name: row.name,
        },
      }
      NetworkApi.deleteNetwork(pData).then((res) => {
        handleResult(res, getNetworkList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
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
.mobile-tips {
  position: absolute;
  bottom: 20px;
}
</style>
