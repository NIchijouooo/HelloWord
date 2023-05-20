<template>
  <div class="main-container">
    <div class="main">
      <div class="title" style="justify-content: space-between">
        <el-input style="width: 200px" placeholder="请输入采集接口名称" v-model="ctxData.interfaceName">
          <template #prefix>
            <el-icon class="el-input__icon"><search /></el-icon>
          </template>
        </el-input>
        <div>
          <el-button type="primary" bg class="right-btn" @click="addInterface()">
            <el-icon class="btn-icon">
              <Icon name="local-add" size="14px" color="#ffffff" />
            </el-icon>
            添加
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
          @row-dblclick="editInterface"
        >
          <el-table-column prop="collInterfaceName" label="采集接口名称" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column prop="commInterfaceName" label="通讯接口" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column prop="protocolTypeName" label="通讯协议" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column prop="pollPeriod" label="采集周期(秒)" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column prop="offlinePeriod" label="离线判断次数" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column prop="deviceNodeCnt" label="设备总数" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column prop="deviceNodeOnlineCnt" label="设备在线数" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column label="操作" width="auto" min-width="300" align="center" fixed="right">
            <template #default="scope">
              <el-button @click="editInterface(scope.row)" text type="primary">编辑</el-button>
              <el-button @click="deleteInterface(scope.row)" text type="danger">删除</el-button>
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
    <!-- 添加编辑采集接口 -->
    <el-dialog
      v-model="ctxData.iFlag"
      :title="ctxData.iTitle"
      width="600px"
      :before-close="handleInterfaceClose"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.interfaceForm"
          :rules="ctxData.interfaceRules"
          ref="interfaceFormRef"
          status-icon
          label-position="right"
          label-width="120px"
        >
          <el-form-item label="采集接口名称" prop="collInterfaceName">
            <el-input
              :disabled="ctxData.iTitle.includes('编辑')"
              type="text"
              v-model="ctxData.interfaceForm.collInterfaceName"
              autocomplete="off"
              placeholder="请输入采集接口名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="通讯接口名称" prop="commInterfaceName">
            <el-select
              v-model="ctxData.interfaceForm.commInterfaceName"
              style="width: 100%"
              placeholder="请选择通讯接口名称"
            >
              <el-option
                v-for="(item, index) of ctxData.commInterfaceList"
                :key="'comm_' + index"
                :label="item.name"
                :value="item.name"
              >
              </el-option>
            </el-select>
          </el-form-item>

            <el-form-item label="通讯协议" prop="protocolTypeName">
              <el-select
                v-model="ctxData.interfaceForm.protocolTypeName"
                style="width: 100%"
                placeholder="请选择通讯协议(GWAI ADD 2023-05-10)"
              >
                <el-option
                  v-for="item of ctxData.protocolTypeNameList"
                  :key="item.name"
                  :label="item.label"
                  :value="item.name"
                >
                </el-option>
              </el-select>
          </el-form-item>

          <el-form-item label="采集周期" prop="pollPeriod">
            <el-input
              type="text"
              v-model.number="ctxData.interfaceForm.pollPeriod"
              autocomplete="off"
              placeholder="请输入采集周期"
            >
              <template #append>单位秒</template>
            </el-input>
          </el-form-item>
          <el-form-item label="离线判断次数" prop="offlinePeriod">
            <el-input
              type="text"
              v-model.number="ctxData.interfaceForm.offlinePeriod"
              autocomplete="off"
              placeholder="请输入离线判断次数"
            >
              <template #append>单位次</template>
            </el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelInterfaceSubmit()">取消</el-button>
          <el-button type="primary" @click="submitInterfaceForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import variables from 'styles/variables.module.scss'
import InterfaceApi from 'api/interface.js'
import CommInterfaceApi from 'api/commInterface.js'
import { userStore } from 'stores/user'
import {reactive} from "vue";
const users = userStore()

const ctxData = reactive({
  interfaceName: '',
  checkAll: true,
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
  interfaceTableData: [],
  iFlag: false,
  iTitle: '添加采集接口',
  interfaceForm: {
    collInterfaceName: '', //采集接口名称，只能是字母+数字的组合，不可以是中文
    commInterfaceName: '', //通讯接口名称
    protocolTypeName:'',//协议类型名称,新增加的协议通过这个变量区分协议，不在通信接口里面做了--gwai add 2023-05-10
    pollPeriod: 0, // 采集周期
    offlinePeriod: 0, // 离线判断次数
  },
  interfaceRules: {
    collInterfaceName: [
      {
        required: true,
        message: '采集接口不能为空',
        trigger: 'blur',
      },
    ],
    commInterfaceName: [
      {
        required: true,
        message: '通讯接口不能为空',
        trigger: 'blur',
      },
    ],
    pollPeriod: [
      {
        required: true,
        message: '采集周期不能为空',
        trigger: 'blur',
      },
      {
        type: 'number',
        message: '采集周期只能输入数字',
      },
    ],
    offlinePeriod: [
      {
        required: true,
        message: '离线判断次数不能为空',
        trigger: 'blur',
      },
      {
        type: 'number',
        message: '离线判断次数只能输入数字',
      },
    ],
  },
  commInterfaceList: [],
  protocolTypeNameList: [   //GWAI ADD 2023-05-10
    { name: 'NULL', label: 'NULL' },
    { name: 'LUA', label: 'LUA' },
    { name: 'DLT645-2007', label: 'DLT645-2007' },
    { name: 'MODBUS-RTU', label: 'MODBUS-RTU' },
    { name: 'MODBUS-TCP', label: 'MODBUS-TCP' },
  ],
})
const contentRef = ref(null)
// 获取采集接口列表
const getInterfaceList = (flag) => {
  const pData = {
    token: users.token,
    data: {},
  }
  InterfaceApi.getInterfaceList(pData).then(async (res) => {
    console.log('getInterfaceList -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.interfaceTableData = res.data

      // 循环处理通讯协议无数据，设置默认值  gwai add 2023-05-10
      ctxData.interfaceTableData.forEach(item => {
        if (!item.protocolTypeName){   //GWAI ADD 2023-05-10
          item.protocolTypeName = ctxData.protocolTypeNameList[0].name
        }
      })

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
getInterfaceList()
const filterTableData = computed(() => {
  return ctxData.interfaceTableData
    .filter((item) => {
      var a = !ctxData.interfaceName
      var b = item.collInterfaceName.toLowerCase().includes(ctxData.interfaceName.toLowerCase())
      return a || b
    })
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  return ctxData.interfaceTableData.filter((item) => {
    var a = !ctxData.interfaceName
    var b = item.collInterfaceName.toLowerCase().includes(ctxData.interfaceName.toLowerCase())
    return a || b
  })
})
const refresh = () => {
  getInterfaceList(1)
}
//获取通讯接口列表
const getCommInterfaceList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  CommInterfaceApi.getCommInterfaceList(pData).then((res) => {
    if (!res) return
    console.log('getCommInterfaceList -> res', res)
    if (res.code === '0') {
      ctxData.commInterfaceList = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
getCommInterfaceList()
// 添加采集接口
const addInterface = () => {
  ctxData.iFlag = true
  ctxData.iTitle = '添加采集接口'
}
// 编辑采集接口
const editInterface = (row) => {
  ctxData.iFlag = true
  ctxData.iTitle = '编辑采集接口'
  ctxData.interfaceForm.collInterfaceName = row.collInterfaceName
  ctxData.interfaceForm.commInterfaceName = row.commInterfaceName
  ctxData.interfaceForm.protocolTypeName = row.protocolTypeName    //gwai add 2023-05-10
  ctxData.interfaceForm.pollPeriod = row.pollPeriod
  ctxData.interfaceForm.offlinePeriod = row.offlinePeriod
}
const interfaceFormRef = ref(null)
const submitInterfaceForm = () => {
  interfaceFormRef.value.validate((valid) => {
    if (valid) {
      const pData = {
        token: users.token,
        data: ctxData.interfaceForm,
      }
      if (ctxData.iTitle.includes('添加')) {
        InterfaceApi.addInterface(pData).then((res) => {
          handleResult(res, getInterfaceList)
          cancelInterfaceSubmit()
        })
      } else {
        InterfaceApi.editInterface(pData).then((res) => {
          handleResult(res, getInterfaceList)
          cancelInterfaceSubmit()
        })
      }
    } else {
      return false
    }
  })
}

//全选操作
const handleCheckAllChange = (val) => {
  ctxData.checkedInterfaceList = val ? ctxData.interfaces : []
  ctxData.isIndeterminate = false
  ctxData.interfaceTableData.forEach((item) => {
    item.checkFlag = val
  })
}

const cancelInterfaceSubmit = () => {
  ctxData.iFlag = false
  interfaceFormRef.value.resetFields()
  initInterfaceForm()
}
//处理弹出框右上角关闭图标
const handleInterfaceClose = (done) => {
  cancelInterfaceSubmit()
}
const initInterfaceForm = () => {
  ctxData.interfaceForm = {
    collInterfaceName: '', //采集接口名称，只能是字母+数字的组合，不可以是中文
    commInterfaceName: '', //通讯接口名称
    protocolTypeName: '', //协议接口  gwai add 2023-05-10
    pollPeriod: 0, // 采集周期
    offlinePeriod: 0, // 离线判断次数
  }
}

// 删除采集接口
const deleteInterface = (row) => {
  ElMessageBox.confirm('确定要删除这个采集接口吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          name: row.collInterfaceName,
        },
      }
      InterfaceApi.deleteInterface(pData).then((res) => {
        handleResult(res, getInterfaceList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
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
    type: res.code === '0' ? 'success' : res.code === '1' ? 'error' : 'warning',
    message: res.message,
  })
  if (res.code === '0' && doFunction) {
    doFunction()
  }
}
</script>
<style lang="scss" scoped>
@use 'styles/custom-scoped.scss' as *;
.mclc-hItem {
  display: flex;
  justify-content: center;
  align-items: center;
}
.head-tag {
  float: right;
  border: 0;
  margin-left: 4px;
  height: 22px;
  padding: 0 2px;
}
:deep(.el-card__body) {
  padding: 10px;
}
</style>
