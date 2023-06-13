<template>
  <div class="main-container">
    <div class="main" v-if="ctxData.isModel">
      <div class="search-bar">
        <el-form :inline="true" ref="searchFormRef" status-icon label-width="120px">
          <el-form-item label="上报模型名称">
            <el-input style="width: 200px" placeholder="请输入上报模型名称" clearable v-model="ctxData.transferModelName">
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
        <el-button type="primary" bg class="right-btn" @click="addTransferModel()">
          <el-icon class="btn-icon">
            <Icon name="local-add" size="14px" color="#ffffff" />
          </el-icon>
          添加
        </el-button>
      </div>
      <div class="content" ref="contentRef">
        <el-table
          :data="filterTableData"
          :cell-style="ctxData.cellStyle"
          :header-cell-style="ctxData.headerCellStyle"
          :max-height="ctxData.tableMaxHeight"
          style="width: 100%"
          stripe
          @row-dblclick="editTransferModel"
        >
          <el-table-column sortable prop="name" label="上报模型名称" width="auto" min-width="200" align="center">
          </el-table-column>
          <el-table-column sortable prop="label" label="上报模型标签" width="auto" min-width="200" align="center">
          </el-table-column>
          <el-table-column sortable prop="code" label="上报模型编号" width="auto" min-width="200" align="center">
          </el-table-column>
          <el-table-column label="操作" width="auto" min-width="200" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="showTransferModelInfo(scope.row)" text type="success">变量详情</el-button>
              <el-button @click="editTransferModel(scope.row)" text type="primary">编辑</el-button>
              <el-button @click="deleteTransferModel(scope.row)" text type="danger">删除</el-button>
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
            :page-sizes="[5, 10, 20, 50]"
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
    <ModelProperty v-else :curTransferModel="ctxData.curTransferModel" @changeTmFlag="changeTmFlag()" style="width: 100%; height: 100%;overflow:hidden;"></ModelProperty>
    <el-dialog
      v-model="ctxData.tFlag"
      :title="ctxData.tTitle"
      width="600px"
      :before-close="handleClose"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.transferForm"
          :rules="ctxData.transferRules"
          ref="transferFormRef"
          status-icon
          label-position="right"
          label-width="120px"
        >
          <el-form-item label="上报模型名称" prop="name">
            <el-select
              v-model="ctxData.transferForm.name"
              filterable
              allow-create
              default-first-option
              placeholder="请选择或者输入上报模型名称"
              style="width: 100%"
            >
              <el-option
                v-for="item in ctxData.collectOptions"
                :key="item.name"
                :label="item.name"
                :value="item.name"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="上报模型标签" prop="label">
            <el-select
              v-model="ctxData.transferForm.label"
              filterable
              allow-create
              default-first-option
              placeholder="请选择或者输入上报模型标签"
              style="width: 100%"
            >
              <el-option
                v-for="item in ctxData.collectOptions"
                :key="item.label"
                :label="item.label"
                :value="item.label"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="上报模型编号" prop="code">
            <el-input
              type="text"
              v-model="ctxData.transferForm.code"
              autocomplete="off"
              placeholder="请输入上报模型编号"
            >
            </el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelSubmit()">取消</el-button>
          <el-button type="primary" @click="submitTransferForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import variables from 'styles/variables.module.scss'
import TransferModelApi from 'api/transferModel.js'
import DeviceModelApi from 'api/deviceModel.js'
import { userStore } from 'stores/user'
import ModelProperty from './transferModel/ModelProperty.vue'
const users = userStore()

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
  pagesize: 10, // 每页数据个数
  transferName: '',
  transferTableData: [
    {
      name: '测试01',
      label: '测试01',
    },
  ],
  isModel: true,
  curTransferModel: null,
  tFlag: false,
  tTitle: '添加上报模型',
  transferForm: {
    name: '',
    label: '',
    code: '',
  },
  transferRules: {
    name: [
      {
        required: true,
        message: '上报模型名称不能为空',
        trigger: 'blur',
      },
    ],
    label: [
      {
        required: true,
        message: '上报模型标签不能为空',
        trigger: 'blur',
      },
    ],
    code: [
      {
        required: true,
        message: '上报模型编号不能为空',
        trigger: 'blur',
      },
    ],
  },
  collectOptions: [],
})
// 获取上报模型列表
const getTransferList = (flag) => {
  const pData = {
    token: users.token,
    data: {},
  }
  TransferModelApi.getModelList(pData).then(async (res) => {
    console.log('getModelList -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.transferTableData = res.data
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
getTransferList()

const getCollectList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  DeviceModelApi.getDeviceModelList(pData).then((res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.collectOptions = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
getCollectList()
// 刷新
const refresh = () => {
  getTransferList(1)
}
//处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
//处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
// 过滤表格数据
const filterTableData = computed(() =>
  ctxData.transferTableData
    .filter(
      (item) => !ctxData.transferModelName || item.name.toLowerCase().includes(ctxData.transferModelName.toLowerCase())
    )
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
)
const filterTableDataPage = computed(() =>
  ctxData.transferTableData.filter(
    (item) => !ctxData.transferModelName || item.name.toLowerCase().includes(ctxData.transferModelName.toLowerCase())
  )
)
const initTranserModel = () => {
  ctxData.transferForm = {
    name: '',
    label: '',
    code: '',
  }
}
//编辑上报模型
const addTransferModel = () => {
  initTranserModel()
  ctxData.tTitle = '添加上报模型'
  ctxData.tFlag = true
}
//编辑上报模型
const editTransferModel = (row) => {
  ctxData.tTitle = '编辑上报模型'
  ctxData.tFlag = true
  ctxData.transferForm.name = row.name
  ctxData.transferForm.label = row.label
  ctxData.transferForm.code = row.code
}
const transferFormRef = ref(null)
const submitTransferForm = () => {
  transferFormRef.value.validate((valid) => {
    if (valid) {
      const pData = {
        token: users.token,
        data: ctxData.transferForm,
      }
      if (ctxData.tTitle.includes('添加')) {
        TransferModelApi.addModel(pData).then((res) => {
          handleResult(res, getTransferList)
          cancelSubmit()
        })
      } else {
        TransferModelApi.editModel(pData).then((res) => {
          handleResult(res, getTransferList)
          cancelSubmit()
        })
      }
    } else {
      return false
    }
  })
}

const cancelSubmit = () => {
  ctxData.tFlag = false
}

const handleClose = () => {
  cancelSubmit()
}
//删除上报模型
const deleteTransferModel = (row) => {
  ElMessageBox.confirm('确定要删除这个上报模型吗?', '警告', {
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
      TransferModelApi.deleteModel(pData).then((res) => {
        handleResult(res, getTransferList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}
//查看变量信息
const showTransferModelInfo = (row) => {
  ctxData.curTransferModel = row
  ctxData.isModel = false
}
//更改页面状态
const changeTmFlag = () => {
  ctxData.isModel = true
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
