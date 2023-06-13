<template>
  <div class="main-container">
    <div class="main">
      <div class="title" style="justify-content: space-between">
        <div>
          <el-input
            v-model="ctxData.queryParams.dictName"
            placeholder="请输入字典名称"
            clearable
            style="width: 200px"
            @change="handleQuery"
          >
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
        </div>
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button type="primary" bg @click="handleAdd">
              <el-icon class="btn-icon">
                <Icon name="local-add" size="14px" color="#ffffff" />
              </el-icon>
              添加
            </el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button style="color: #fff" color="#2EA554" @click="getList">
              <el-icon class="btn-icon">
                <Icon name="local-refresh" size="14px" color="#ffffff" />
              </el-icon>
              刷新
            </el-button>
          </el-col>
        </el-row>
      </div>

      <div class="content" ref="contentRef">
        <el-table
          :data="ctxData.typeList"
          :cell-style="ctxData.cellStyle"
          :header-cell-style="ctxData.headerCellStyle"
          :max-height="ctxData.tableMaxHeight"
          style="width: 100%"
          stripe
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" align="center"></el-table-column>
          <el-table-column label="字典编号" align="center" width="100" prop="dictId"></el-table-column>
          <el-table-column label="字典名称" align="center" prop="dictName"></el-table-column>
          <el-table-column label="字典类型" align="center">
            <template #default="scope">
              <span>{{ scope.row.dictType }}</span>
            </template>
          </el-table-column>
          <el-table-column label="状态" align="center" prop="status" :formatter="statusFormat"></el-table-column>
          <el-table-column label="备注" align="center" prop="remark"></el-table-column>
          <el-table-column label="操作" fixed="right" align="center" class-name="small-padding fixed-width" width="250">
            <template #default="scope">
              <el-button text type="success" @click="toData(scope.row.dictId)">查看</el-button>
              <el-button text type="primary" @click="handleUpdate(scope.row)">修改</el-button>
              <el-button text type="danger" @click="handleDelete(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="pagination">
          <el-pagination
            :current-page="ctxData.queryParams.pageNum"
            :page-size="ctxData.queryParams.pageSize"
            :page-sizes="[20, 50, 200, 500]"
            :total="ctxData.total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            background
            layout="total, sizes, prev, pager, next, jumper"
            style="margin-top: 46px"
          ></el-pagination>
        </div>
      </div>
    </div>

    <!-- 添加或修改参数配置对话框 -->
    <el-dialog :title="ctxData.title" v-model="ctxData.open" width="500px" append-to-body :close-on-click-modal="false">
      <el-form ref="formRef" :model="ctxData.form" :rules="ctxData.rules" label-width="80px">
        <el-form-item label="字典名称" prop="dictName">
          <el-input v-model="ctxData.form.dictName" placeholder="请输入字典名称"></el-input>
        </el-form-item>
        <el-form-item label="字典类型" prop="dictType">
          <el-input v-model="ctxData.form.dictType" placeholder="请输入字典类型"></el-input>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="ctxData.form.status">
            <el-radio
              v-for="dict in ctxData.statusOptions"
              :key="dict.dictValue"
              :label="dict.dictValue"
            >{{ dict.dictLabel }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="ctxData.form.remark" type="textarea" placeholder="请输入内容" ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancel()">取消</el-button>
          <el-button type="primary" @click="submitForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import typeApi from '../../../api/dict/type'
import dataApi from '../../../api/dict/data'
import { Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from "vue-router";
const router = useRouter();
import { userStore } from 'stores/user'
import variables from 'styles/variables.module.scss'
const users = userStore()

// 回显数据字典
const selectDictLabel = (datas, value) => {
  var actions = [];
  Object.keys(datas).some((key) => {
    if (datas[key].dictValue == ('' + value)) {
      actions.push(datas[key].dictLabel);
      return true;
    }
  })
  return actions.join('');
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
  // 遮罩层
  loading: true,
    // 选中数组
    ids: [],
  // 非单个禁用
  single: true,
  // 非多个禁用
  multiple: true,
  // 总条数
  total: 0,
  // 字典表格数据
  typeList: [],
  // 弹出层标题
  title: '',
  // 是否显示弹出层
  open: false,
  // 状态数据字典
  statusOptions: [],
  // 日期范围
  dateRange: [],
  // 查询参数
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    dictName: undefined,
    dictType: undefined,
    status: undefined
  },
  // 表单参数
  form: {},
  // 表单校验
  rules: {
    dictName: [
      { required: true, message: '字典名称不能为空', trigger: 'blur' }
    ],
      dictType: [
      { required: true, message: '字典类型不能为空', trigger: 'blur' }
    ]
  }
})

nextTick(() => {
  ctxData.tableMaxHeight = contentRef.value.clientHeight - 34 - 36 - 22
})

const getStatusOptions = () => {
  const pdata = {
    token: users.token,
    data: {
      dictType: 'sys_normal_disable'
    }
  }
  dataApi.getDicts(pdata).then(response => {
    ctxData.statusOptions = response.data
  })
}
getStatusOptions()

const getList = () => {
  ctxData.loading = true
  const pdata = {
    token: users.token,
    data: ctxData.queryParams
  }
  typeApi.listType(pdata).then(response => {
    ctxData.typeList = response.data.data
    ctxData.total = response.data.total
    ctxData.loading = false
  })
}
getList()

const handleCurrentChange = (value) => {
  ctxData.queryParams.pageNum = value
  getList()
}


const handleSizeChange = (value) => {
  ctxData.queryParams.pageSize = value
  getList()
}

const toData = (id) => {
  router.push({ path: '/dictData', query: {dictId: id} })
}

// 字典状态字典翻译
const statusFormat = (row, column) => {
  return selectDictLabel(ctxData.statusOptions, row.status)
}

// 取消按钮
const cancel = () => {
  ctxData.open = false
  reset()
}

// 表单重置
const formRef = ref(null)
const reset = () => {
  ctxData.form = {
    dictId: undefined,
    dictName: undefined,
    dictType: undefined,
    status: '0',
    remark: undefined
  }
  formRef.value && formRef.value.resetFields()
}

// 搜索按钮操作
const handleQuery = () =>  {
  ctxData.queryParams.pageNum = 1
  getList()
}

// 重置按钮操作
const queryFormRef = ref(null)
const resetQuery = () => {
  ctxData.dateRange = []
  queryFormRef.value.resetFields()
  handleQuery()
}

// 新增按钮操作
const handleAdd = () => {
  reset()
  ctxData.open = true
  ctxData.title = '添加字典类型'
}


// 多选框选中数据
const handleSelectionChange = (selection) => {
  ctxData.ids = selection.map(item => item.dictId)
  // eslint-disable-next-line eqeqeq
  ctxData.single = selection.length != 1
  ctxData.multiple = !selection.length
}

// 修改按钮操作
const handleUpdate = (row) => {
  reset()
  const dictId = row.dictId || ctxData.ids
  const pdata = {
    token: users.token,
    data: {dictId}
  }
  typeApi.getType(pdata).then(response => {
    ctxData.form = response.data
    ctxData.open = true
    ctxData.title = '修改字典类型'
  })
}

// 提交按钮
const submitForm = () => {
  formRef.value.validate(valid => {
    if (valid) {
      const pdata = {
        token: users.token,
        data: ctxData.form
      }
      // eslint-disable-next-line eqeqeq
      if (ctxData.form.dictId != undefined) {
        typeApi.updateType(pdata).then(response => {
          ElMessage({type: 'success', message: '修改成功',})
          ctxData.open = false
          getList()
        })
      } else {
        typeApi.addType(pdata).then(response => {
          ElMessage({type: 'success', message: '新增成功',})
          ctxData.open = false
          getList()
        })
      }
    }
  })
}

// 删除按钮操作
const handleDelete = (row) => {
  const dictIds = row.dictId || ctxData.ids
  ElMessageBox.confirm('是否确认删除字典编号为"' + dictIds + '"的数据项?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(function() {
    const pdata = {
      token: users.token,
      dictIds
    }
    return typeApi.delType(pdata)
  }).then(() => {
    getList()
    ElMessage({type: 'success', message: '删除成功',})
  }).catch(function() {})
}
</script>

<style lang="scss" scoped>
@use 'styles/custom-scoped.scss' as *;
</style>
