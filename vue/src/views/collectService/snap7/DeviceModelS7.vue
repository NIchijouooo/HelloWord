<template>
  <div class="main-container">
    <!-- 模型页 -->
    <div class="main" v-if="ctxData.dpFlag">
      <div class="search-bar">
        <el-form :inline="true" ref="searchFormRef" status-icon label-width="120px">
          <el-form-item label="采集模型名称">
            <el-input style="width: 200px" placeholder="请输入采集模型名称" v-model="ctxData.deviceModelInfo">
              <template #prefix>
                <el-icon class="el-input__icon"><search /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button style="color: #fff; margin-left: 20px" color="#2EA554" class="right-btn" @click="refresh()">
              <el-icon class="btn-icon">
                <Icon name="local-refresh" size="14px" color="#ffffff" />
              </el-icon>
              刷新
            </el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="tool-bar">
        <el-button type="primary" bg class="right-btn" @click="addDeviceModel()">
          <el-icon class="btn-icon">
            <Icon name="local-add" size="14px" color="#ffffff" />
          </el-icon>
          添加
        </el-button>
      </div>
      <div class="content" ref="contentRef">
        <el-table
          :data="filterDMTableData"
          :cell-style="ctxData.cellStyle"
          :header-cell-style="ctxData.headerCellStyle"
          :max-height="ctxData.tableMaxHeight"
          style="width: 100%"
          stripe
          @row-dblclick="editDeviceModel"
        >
          <el-table-column prop="name" label="采集模型名称" width="auto" min-width="200" align="center">
          </el-table-column>
          <el-table-column prop="label" label="采集模型标签" width="auto" min-width="200" align="center">
          </el-table-column>
          <el-table-column label="操作" width="auto" min-width="300" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="showVariableDetail(scope.row)" text type="success">变量详情</el-button>
              <el-button @click="editDeviceModel(scope.row)" text type="primary">编辑</el-button>
              <el-button @click="deleteDeviceModel(scope.row)" text type="danger">删除</el-button>
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
    <!-- s7 -->
    <PropertyS7 v-else :curDeviceModel="ctxData.curDeviceModel" @changeDpFlag="changeDpFlag()" style="width: 100%; height: 100%;overflow:hidden;"></PropertyS7>
    <!-- dialog 内容 -->
    <!-- 添加编辑采集模型 -->
    <el-dialog
      v-model="ctxData.dmFlag"
      :title="ctxData.dmTitle"
      width="600px"
      :before-close="handleDMClose"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.modelForm"
          :rules="ctxData.modelRules"
          ref="modelFormRef"
          status-icon
          label-position="right"
          label-width="120px"
        >
          <el-form-item label="采集模型名称" prop="name">
            <el-input
              type="text"
              :disabled="ctxData.dmTitle.includes('编辑')"
              v-model="ctxData.modelForm.name"
              autocomplete="off"
              placeholder="请输入采集模型名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="采集模型标签" prop="label">
            <el-input type="text" v-model="ctxData.modelForm.label" autocomplete="off" placeholder="请输入采集模型标签">
            </el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelModelSubmit()">取消</el-button>
          <el-button type="primary" @click="submitModelForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { Search } from '@element-plus/icons-vue'
import variables from 'styles/variables.module.scss'
import DeviceModelApi from 'api/deviceModel.js'

import PropertyS7 from './PropertyS7.vue'
import { userStore } from 'stores/user'
const users = userStore()
const contentRef = ref(null)
const ctxData = reactive({
  modelType: 1,
  dpFlag: true, //模型页和变量页切换标志
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
  deviceModelInfo: '',
  tableData: [],
  // 模型表单
  modelForm: {
    name: '', //HL7031,名称，只能是字母+数字的组合，不可以是中文
    label: '', //海林风机盘管控制器，标签，可以是中文
    type: 0, //模型类型，0：自定义插件LUA;1：西门子PLC协议
    param: '', //HL7031，插件名称
  },
  modelRules: {
    name: [
      {
        required: true,
        message: '采集模型名称不能为空',
        trigger: 'blur',
      },
    ],
    label: [
      {
        required: true,
        message: '采集模型标识不能为空',
        trigger: 'blur',
      },
    ],
  },
  udFlag: false,
  dmFlag: false, //采集模型弹窗表示
  dmTitle: '添加采集模型',
  pluginList: [], //
  curDeviceModel: '', //当前采集模型
})
// 获取采集模型列表
const getDeviceModelList = (flag) => {
  //
  const pData = {
    token: users.token,
    data: {
      type: ctxData.modelType,
    },
  }
  DeviceModelApi.getDeviceModelList(pData).then(async (res) => {
    console.log('getDeviceModelList -> res = ', res)
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
getDeviceModelList()
const filterDMTableData = computed(() => {
  return ctxData.tableData
    .filter((item) => {
      var a = !ctxData.deviceModelInfo
      var b =
        item.name.toLowerCase().includes(ctxData.deviceModelInfo.toLowerCase()) ||
        item.label.toLowerCase().includes(ctxData.deviceModelInfo.toLowerCase())
      return a || b
    })
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  return ctxData.tableData.filter((item) => {
    var a = !ctxData.deviceModelInfo
    var b =
      item.name.toLowerCase().includes(ctxData.deviceModelInfo.toLowerCase()) ||
      item.label.toLowerCase().includes(ctxData.deviceModelInfo.toLowerCase())
    return a || b
  })
})
// 刷新
const refresh = () => {
  getDeviceModelList(1)
}
// 处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
// 添加采集模型
const addDeviceModel = () => {
  ctxData.dmFlag = true
  ctxData.dmTitle = '添加采集模型'
}
// 编辑采集模型
const editDeviceModel = (row) => {
  ctxData.dmFlag = true
  ctxData.dmTitle = '编辑采集模型'
  ctxData.modelForm = {
    name: row.name,
    label: row.label,
    type: row.type === undefined ? 1 : row.type,
    param: '',
  }
}

const modelFormRef = ref(null)
// 提交采集模型表单
const submitModelForm = () => {
  modelFormRef.value.validate((valid) => {
    if (valid) {
      ctxData.modelForm.type = ctxData.modelType
      const pData = {
        token: users.token,
        data: ctxData.modelForm,
      }
      if (ctxData.dmTitle.includes('添加')) {
        DeviceModelApi.addDeviceModel(pData).then((res) => {
          handleResult(res, getDeviceModelList)
          cancelModelSubmit()
        })
      } else {
        DeviceModelApi.editDeviceModel(pData).then((res) => {
          handleResult(res, getDeviceModelList)
          cancelModelSubmit()
        })
      }
    } else {
      return false
    }
  })
}
// 取消采集模型提交
const cancelModelSubmit = () => {
  ctxData.dmFlag = false
  modelFormRef.value.resetFields()
  initDeviceModelForm()
}
const handleDMClose = () => {
  cancelModelSubmit()
}
const deleteDeviceModel = (row) => {
  ElMessageBox.confirm('确定要删除这个采集模型吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          name: row.name,
          type: 1,
        },
      }
      DeviceModelApi.deleteDeviceModel(pData).then((res) => {
        handleResult(res, getDeviceModelList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}
const initDeviceModelForm = () => {
  ctxData.modelForm = {
    name: '',
    label: '',
    type: 0,
  }
}
// 查看变量详情
const showVariableDetail = (row) => {
  ctxData.curDeviceModel = row
  ctxData.dpFlag = false
}

const changeDpFlag = () => {
  ctxData.dpFlag = true
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
</style>
<style lang="scss">
.el-popover.el-popper {
  padding: 20px;
}
</style>
