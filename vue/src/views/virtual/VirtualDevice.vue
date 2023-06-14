<template>
  <div class="main-container">
    <div class="main" v-if="ctxData.dpFlag">
      <div class="title" style="justify-content: space-between">
        <el-input style="width: 200px" placeholder="请输入虚拟设备名/标签" v-model="ctxData.virtualDeviceInfo">
          <template #prefix>
            <el-icon class="el-input__icon"><search /></el-icon>
          </template>
        </el-input>
        <div>
          <el-button type="primary" bg class="right-btn" @click="addVirtualDevice()">
            <el-icon class="btn-icon">
              <Icon name="local-add" size="14px" color="#ffffff" />
            </el-icon>
            添加
          </el-button>
          <el-button type="danger" bg class="right-btn" @click="deleteDeviceList()">
            <el-icon class="btn-icon">
              <Icon name="local-delete" size="14px" color="#ffffff" />
            </el-icon>
            删除
          </el-button>
          <el-button style="color: #fff" color="#2EA554" class="right-btn" @click="refresh()">
            <el-icon class="btn-icon">
              <Icon name="local-refresh" size="14px" color="#ffffff" />
            </el-icon>
            刷新
          </el-button>
        </div>
      </div>
      <div class="content" ref="contentRef">
        <el-table
          :data="filterTableData"
          :cell-style="ctxData.cellStyle"
          :header-cell-style="ctxData.headerCellStyle"
          :max-height="ctxData.tableMaxHeight"
          style="width: 100%"
          stripe
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column sortable prop="name" label="虚拟设备名称" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column sortable prop="label" label="虚拟设备标签" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column label="操作" width="auto" min-width="300" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="showVariableDetail(scope.row)" text type="success">变量详情</el-button>
              <el-button @click="editVirtualDevice(scope.row)" text type="primary">编辑</el-button>
              <el-button @click="deleteVirtualDevice(scope.row)" text type="danger">删除</el-button>
            </template>
          </el-table-column>
          <template #empty>
            <div>无数据</div>
          </template>
        </el-table>
        <div class="pagination">
          <el-pagination
            :current-page="ctxData.currentPage"
            :page-size="ctxData.pagesize"
            :page-sizes="[20, 50, 200, 500]"
            :total="filterTableDataPage.length"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            background
            layout="total, sizes, prev, pager, next, jumper"
            style="margin-top: 46px"
          ></el-pagination>
        </div>
      </div>
    </div>
    <!-- 变量页 -->
    <DeviceProperty v-else :curVirtualDevice="ctxData.curVirtualDevice" @changeDpFlag="changeDpFlag()"></DeviceProperty>
    <!-- 添加编辑虚拟设备 -->
    <el-dialog
      v-model="ctxData.vdFlag"
      :title="ctxData.vdTitle"
      width="600px"
      :before-close="handleVDClose"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.virtualForm"
          :rules="ctxData.virtualRules"
          ref="virtualFormRef"
          status-icon
          label-position="right"
          label-width="120px"
        >
          <el-form-item label="虚拟设备名称" prop="name">
            <el-input
              type="text"
              :disabled="ctxData.vdTitle.includes('编辑')"
              v-model="ctxData.virtualForm.name"
              autocomplete="off"
              placeholder="请输入虚拟设备名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="虚拟设备标签" prop="label">
            <el-input
              type="text"
              v-model="ctxData.virtualForm.label"
              autocomplete="off"
              placeholder="请输入虚拟设备标签"
            >
            </el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelVirtualSubmit()">取消</el-button>
          <el-button type="primary" @click="submitVirtualForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { Search } from '@element-plus/icons-vue'
import variables from 'styles/variables.module.scss'
import VirtualServiceApi from 'api/virtualService.js'
import DeviceProperty from './DeviceProperty.vue'
import { userStore } from 'stores/user'
const users = userStore()

const contentRef = ref(null)
const ctxData = reactive({
  virtualDeviceInfo: '',
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
  tableData: [],
  dpFlag: true,
  vdFlag: false, //虚拟设备弹窗标识
  vdTitle: '添加虚拟设备',
  curVirtualDevice: {}, //当前虚拟设备
  virtualForm: {
    name: '',
    label: '',
  },
  virtualRules: {
    name: [
      {
        required: true,
        message: '虚拟设备名称不能为空',
        trigger: 'blur',
      },
    ],
    label: [
      {
        required: true,
        message: '虚拟设备标识不能为空',
        trigger: 'blur',
      },
    ],
  },
  selectedDevice: [],
})

const getVirtualDeviceList = (flag) => {
  const pData = {
    token: users.token,
    data: {},
  }
  VirtualServiceApi.getDeviceList(pData).then(async (res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.tableData = res.data
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
getVirtualDeviceList()
const filterTableData = computed(() => {
  return ctxData.tableData
    .filter((item) => {
      var a = !ctxData.virtualDeviceInfo
      var b = item.name.toLowerCase().includes(ctxData.virtualDeviceInfo.toLowerCase())
      return a || b
    })
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  return ctxData.tableData.filter((item) => {
    var a = !ctxData.virtualDeviceInfo
    var b = item.name.toLowerCase().includes(ctxData.virtualDeviceInfo.toLowerCase())
    return a || b
  })
})

// 刷新
const refresh = () => {
  getVirtualDeviceList(1)
}
// 处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
// 添加虚拟设备
const addVirtualDevice = () => {
  ctxData.vdFlag = true
  ctxData.vdTitle = '添加虚拟设备'
}
// 编辑虚拟设备
const editVirtualDevice = (row) => {
  ctxData.vdFlag = true
  ctxData.vdTitle = '编辑虚拟设备'
  ctxData.virtualForm = {
    name: row.name,
    label: row.label,
  }
}

const virtualFormRef = ref(null)
// 提交虚拟设备表单
const submitVirtualForm = () => {
  virtualFormRef.value.validate((valid) => {
    if (valid) {
      const pData = {
        token: users.token,
        data: ctxData.virtualForm,
      }
      if (ctxData.vdTitle.includes('添加')) {
        VirtualServiceApi.addDevice(pData).then((res) => {
          handleResult(res, getVirtualDeviceList)
          cancelVirtualSubmit()
        })
      } else {
        VirtualServiceApi.editDevice(pData).then((res) => {
          handleResult(res, getVirtualDeviceList)
          cancelVirtualSubmit()
        })
      }
    } else {
      return false
    }
  })
}
// 取消虚拟设备提交
const cancelVirtualSubmit = () => {
  ctxData.vdFlag = false
  virtualFormRef.value.resetFields()
  initVirtualDeviceForm()
}
const handleVDClose = () => {
  cancelVirtualSubmit()
}

const handleSelectionChange = (val) => {
  ctxData.selectedDevice = val
  console.log('handleSelectionChange -> val =', val)
}

const deleteVirtualDevice = (row) => {
  ElMessageBox.confirm('确定要删除这个虚拟设备吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          names: [row.name],
        },
      }
      VirtualServiceApi.deleteDevice(pData).then((res) => {
        handleResult(res, getVirtualDeviceList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}

const deleteDeviceList = () => {
  let pList = []
  if (ctxData.selectedDevice.length === 0) {
    ElMessage.info('请至少选择一个虚拟设备！')
    return
  } else {
    ctxData.selectedDevice.forEach((item) => {
      pList.push(item.name)
    })
  }
  ElMessageBox.confirm('确定要删除这些虚拟设备吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          names: pList,
        },
      }
      VirtualServiceApi.deleteDevice(pData).then((res) => {
        handleResult(res, getVirtualDeviceList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}

const initVirtualDeviceForm = () => {
  ctxData.virtualForm = {
    name: '',
    label: '',
  }
}
// 显示变量详情
const showVariableDetail = (row) => {
  console.log('showVariableDetail ->', row)
  ctxData.curVirtualDevice = row
  ctxData.dpFlag = false
}

//显示单个res结果，code不等于 '0' 的message
const showOneResMsg = (res) => {
  ElMessage({
    type: 'error',
    message: res.message,
  })
}
//更改页面状态
const changeDpFlag = () => {
  ctxData.dpFlag = true
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
</style>
