<template>
<div class="main-container">
  <div class="main">
    <div class="search-bar">
      <el-form :inline="true" ref="searchFormRef" status-icon label-width="120px">
        <el-form-item style="margin-left: 20px;">
          <el-button type="primary" plain @click="toTransferModel()" style="margin-right: 20px">
            <el-icon class="el-input__icon"><back /></el-icon>
            返回上报模型
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="search-bar" style="display: flex;">
      <div class="title" style="position: relative;margin-right: 40px;height: 40px; padding: 0px 0; justify-content: flex-start">
        <div class="tName">{{ props.curTransferModel.label }}</div>
      </div>
      <el-form :inline="true" ref="searchFormRef2" status-icon label-width="90px">
        <el-form-item label="">
          <el-input style="width: 200px" placeholder="请输入属性名称或者标签" clearable v-model="ctxData.PropertyInfo">
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" bg class="right-btn" @click="addProperty()">
            <el-icon class="btn-icon">
              <Icon name="local-add" size="14px" color="#ffffff" />
            </el-icon>
            手动添加
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" bg class="right-btn" @click="addProperties()">
            <el-icon class="btn-icon">
              <Icon name="local-tongbu" size="14px" color="#ffffff" />
            </el-icon>
            同步
          </el-button>
        </el-form-item>
        <el-form-item>
          <!-- <el-button type="primary" bg class="right-btn" @click="editProperties()">
            <el-icon class="btn-icon tianjia"></el-icon>
            批量修改
          </el-button> -->
          <el-button type="danger" bg class="right-btn" @click="deleteProperty()">
            <el-icon class="btn-icon">
              <Icon name="local-delete" size="14px" color="#ffffff" />
            </el-icon>
            删除
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button style="color: #fff" color="#2EA554" class="right-btn" @click="refresh()">
            <el-icon class="btn-icon">
              <Icon name="local-refresh" size="14px" color="#ffffff" />
            </el-icon>
            刷新
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="content" ref="contentRef" style="top: 136px">
      <el-table
        :data="filterTableData"
        :cell-style="ctxData.cellStyle"
        :header-cell-style="ctxData.headerCellStyle"
        :max-height="ctxData.tableMaxHeight"
        style="width: 100%"
        stripe
        @selection-change="handleSelectionChange"
        @row-dblclick="editProperty"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column type="expand">
          <template #default="props">
            <div class="param-content">
              <div class="pc-title">
                <div class="pct-info">
                  <b> {{ props.row.name }} </b>
                  {{ '参数详情' }}
                </div>
              </div>
              <div class="pc-content">
                <div class="param-item" v-for="(item, key, index) of props.row.params" :key="index">
                  <div class="param-value">{{ ctxData.paramName[key] }}：</div>
                  <div class="param-name">{{ typeof item === 'boolean' ? (item ? '是' : '否') : item || '-' }}</div>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column sortable prop="name" label="属性名称" width="auto" min-width="180" align="center"> </el-table-column>
        <el-table-column sortable prop="label" label="属性标签" width="auto" min-width="180" align="center"> </el-table-column>
        <el-table-column sortable prop="uploadName" label="上报属性名称" width="auto" min-width="180" align="center">
        </el-table-column>
        <el-table-column sortable label="数据类型" width="auto" min-width="100" align="center">
          <template #default="scope">
            {{ ctxData.typeNames['t' + scope.row.type] }}
          </template>
        </el-table-column>
        <el-table-column sortable label="小数位数" width="auto" min-width="80" align="center">
          <template #default="scope">
            {{ scope.row.decimals === undefined || scope.row.decimals === '' ? 0 : scope.row.decimals }}
          </template>
        </el-table-column>
        <el-table-column sortable prop="unit" label="单位" width="auto" min-width="80" align="center"> </el-table-column>
        <el-table-column label="操作" width="auto" min-width="150" align="center" fixed="right">
          <template #default="scope">
            <el-button @click="editProperty(scope.row)" text type="primary">编辑</el-button>
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
    <el-dialog
      v-model="ctxData.pFlag"
      :title="ctxData.pTitle"
      width="600px"
      :before-close="handleCloseProperty"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.propertyForm"
          :rules="ctxData.propertyRules"
          ref="propertyFormRef"
          status-icon
          label-position="right"
          label-width="120px"
        >
          <el-form-item label="属性名称" prop="name">
            <el-input
              :disabled="ctxData.pTitle.includes('编辑')"
              type="text"
              v-model="ctxData.propertyForm.name"
              autocomplete="off"
              placeholder="请输入属性名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="属性标签" prop="label">
            <el-input type="text" v-model="ctxData.propertyForm.label" autocomplete="off" placeholder="请输入属性标签">
            </el-input>
          </el-form-item>

          <el-form-item label="上报属性名称" prop="uploadName">
            <el-input
              type="text"
              v-model="ctxData.propertyForm.uploadName"
              autocomplete="off"
              placeholder="请输入上报属性名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="数据类型" prop="type">
            <el-select v-model.number="ctxData.propertyForm.type" style="width: 100%" placeholder="请选择数据类型">
              <el-option
                v-for="item in ctxData.typeOptions"
                :key="'type_' + item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type !== 3" label="范围报警" prop="minMaxAlarm">
            <el-switch v-model="ctxData.propertyForm.minMaxAlarm" inline-prompt active-text="是" inactive-text="否" />
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.minMaxAlarm"
            label="最小值"
            prop="min"
          >
            <el-input type="text" v-model="ctxData.propertyForm.min" autocomplete="off" placeholder="请输入最小值">
            </el-input>
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.minMaxAlarm"
            label="最大值"
            prop="max"
          >
            <el-input type="text" v-model="ctxData.propertyForm.max" autocomplete="off" placeholder="请输入最大值">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type !== 3" label="步长报警" prop="stepAlarm">
            <el-switch v-model="ctxData.propertyForm.stepAlarm" inline-prompt active-text="是" inactive-text="否" />
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.stepAlarm"
            label="步长"
            prop="step"
          >
            <el-input type="text" v-model="ctxData.propertyForm.step" autocomplete="off" placeholder="请输入步长">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type === 2" label="小数位数" prop="decimals">
            <el-input
              type="text"
              v-model.number="ctxData.propertyForm.decimals"
              autocomplete="off"
              placeholder="请输入小数位数"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type !== 3" label="单位" prop="unit">
            <el-input type="text" v-model="ctxData.propertyForm.unit" autocomplete="off" placeholder="请输入单位">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.propertyForm.type === 3" label="字符串长度报警" prop="dataLengthAlarm">
            <el-switch
              v-model="ctxData.propertyForm.dataLengthAlarm"
              inline-prompt
              active-text="是"
              inactive-text="否"
            />
          </el-form-item>
          <el-form-item
            v-if="ctxData.propertyForm.type === 3 && ctxData.propertyForm.dataLengthAlarm"
            label="字符串长度"
            prop="dataLength"
          >
            <el-input
              type="text"
              v-model="ctxData.propertyForm.dataLength"
              autocomplete="off"
              placeholder="请输入字符串长度"
            >
            </el-input>
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelPorpertySubmit()">取消</el-button>
          <el-button type="primary" @click="submitPorpertyForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- 导入属性 -->
    <el-dialog
      v-model="ctxData.psFlag"
      title="导入属性"
      width="600px"
      :before-close="beforeCloseUploadProperty"
      :close-on-click-modal="false"
    >
      <el-upload
        ref="uploadPropertyRef"
        action=""
        :auto-upload="false"
        :http-request="myRequestProperty"
        :limit="1"
        :on-exceed="handleExceedProperty"
        :before-upload="beforeUploadProperty"
      >
        <el-button type="primary">选择文件</el-button>
        <template #tip>
          <div class="el-upload__tip">只能上传一个文件，只支持xlsx格式文件！</div>
        </template>
      </el-upload>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelUploadProperty">取消</el-button>
          <el-button type="primary" @click="submitUploadProperty">上传</el-button>
        </span>
      </template>
    </el-dialog>
    <!--同步采集属性-->
    <el-dialog
      title="同步采集属性"
      v-model="ctxData.plFlag"
      width="1000px"
      :before-close="beforeCloseProperties"
      :close-on-click-modal="false"
    >
      <div class="dialog-content" style="min-height: 408px; overflow: unset; padding: 0">
        <div class="dialog-content-head">
          <el-select v-model="ctxData.curCollName" placeholder="请选择采集模型" @change="selectColl()">
            <el-option
              v-for="items in ctxData.collDeviceModelList"
              :key="items.name"
              :label="items.label"
              :value="items.name"
            >
            </el-option>
          </el-select>
          <el-input placeholder="请输入属性名称或者标签" style="width: 200px" clearable v-model="ctxData.pName">
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
        </div>
        <div class="dialog-content-content">
          <el-table
            :data="
              ctxData.pTableData
                .filter(
                  (item) => !ctxData.pName || item.name.includes(ctxData.pName) || item.label.includes(ctxData.pName)
                )
                .slice((ctxData.pCurrentPage - 1) * ctxData.pPagesize, ctxData.pCurrentPage * ctxData.pPagesize)
            "
            :header-cell-style="ctxData.pHeadStyle"
            :cell-style="ctxData.pCellStyle"
            style="border: 1px solid #c0c4cc"
            max-height="290"
            @selection-change="handleSelect"
            stripe
          >
            <el-table-column type="selection" width="60" :selectable="selectedProp" />
            <el-table-column prop="name" label="属性名称" min-width="200" width="auto" align="center">
            </el-table-column>
            <el-table-column prop="label" label="属性标签" min-width="200" width="auto" align="center">
            </el-table-column>
            <el-table-column label="状态" min-width="120" width="auto" align="center">
              <template #default="scope">
                <span
                  :style="{
                    color: ctxData.AllPropertyNameList.includes(scope.row.name)
                      ? variables.successColor
                      : variables.dangerColor,
                  }"
                  >{{ ctxData.AllPropertyNameList.includes(scope.row.name) ? '已选择' : '未选择' }}</span
                >
              </template>
            </el-table-column>
          </el-table>
          <div class="pagination dialog-pagination">
            <el-pagination
              :current-page="ctxData.pCurrentPage"
              :page-size="ctxData.pPagesize"
              :page-sizes="[5, 20, 50, 200]"
              :total="ctxData.pTableData.length"
              @current-change="handlePCurrentChange"
              @size-change="handlePSizeChange"
              background
              layout="total, sizes, prev, pager, next, jumper"
              style="margin-top: 22px"
            ></el-pagination>
          </div>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelProperties">取消</el-button>
          <el-button type="primary" @click="saveProperties">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</div>
</template>
<script setup>
import { Search, Back } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import variables from 'styles/variables.module.scss'
import TransferModelApi from 'api/transferModel.js'
import DeviceModelApi from 'api/deviceModel.js'
import { userStore } from 'stores/user'
const users = userStore()
const props = defineProps({
  curTransferModel: {
    type: Object,
    default: {},
  },
})
console.log('id -> props', props)

const emit = defineEmits(['changeTmFlag'])
const toTransferModel = () => {
  emit('changeTmFlag')
}

const regCnt = /^[0-9]*[1-9][0-9]*$/
const validateRegCnt = (rule, value, callback) => {
  if (!regCnt.test(value)) {
    callback(new Error('只能输入正整数数字！'))
  } else {
    callback()
  }
}
const regStep = /^[0-9]+(.[0-9]{1,2})?$/
const validateStep = (rule, value, callback) => {
  if (!regStep.test(value)) {
    callback(new Error('只能输入大于等于0,最多两位小数的数字！'))
  } else {
    callback()
  }
}

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
  PropertyTableData: [],
  pFlag: false,
  pTitle: '添加属性',
  PropertyInfo: '',
  propertyForm: {
    name: '', // 属性名称
    label: '', // 属性标签
    uploadName: '', // 上报属性名称
    type: 0, // 数据类型
    //params
    min: 0, // 属性最小值，只有uint32，int32，double有效
    max: 0, // 属性最大值，只有uint32，int32，double有效
    minMaxAlarm: false, // 范围报警，只有uint32，int32，double有效
    step: 0, // 步长，只有uint32，int32，double有效
    stepAlarm: false, // 步长报警，只有uint32，int32，double有效
    unit: 0, // 单位，只有uint32，int32，double有效
    decimals: 0, // 小数位数，只有double有效
    dataLength: '', // 字符串长度，只有string有效
    dataLengthAlarm: false, // 字符串长度报警，只有string有效
  },
  typeOptions: [
    //数据类型
    { label: 'uint32', value: 0 },
    { label: 'int32', value: 1 },
    { label: 'double', value: 2 },
    { label: 'string', value: 3 },
  ],
  typeNames: {
    t0: 'uint32',
    t1: 'int32',
    t2: 'double',
    t3: 'string',
  },
  paramName: {
    name: '属性名称',
    label: '属性标签',
    uploadName: '上报属性名称',
    type: '数据类型',
    min: '最小值',
    max: '最大值',
    minMaxAlarm: '范围报警',
    step: '步长',
    stepAlarm: '步长报警',
    unit: '单位',
    decimals: '小数位数',
    dataLength: '字符串长度',
    dataLengthAlarm: '字符串长度报警',
  },
  propertyRules: {
    name: [
      {
        required: true,
        message: '属性名称不能为空',
        trigger: 'blur',
      },
    ],
    label: [
      {
        required: true,
        message: '属性标签不能为空',
        trigger: 'blur',
      },
    ],
    uploadName: [
      {
        required: true,
        message: '上报属性名称不能为空',
        trigger: 'blur',
      },
    ],
    type: [
      {
        required: true,
        message: '数据类型不能为空',
        trigger: 'blur',
      },
    ],
    decimals: [
      {
        type: 'number',
        message: '小数位数只能输入数字',
      },
    ],
    step: [
      {
        required: true,
        message: '步长不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateStep,
      },
    ],
    min: [
      {
        required: true,
        message: '最小值不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateStep,
      },
    ],
    max: [
      {
        required: true,
        message: '最大值不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateStep,
      },
    ],
    dataLength: [
      {
        required: true,
        message: '字符串长度不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: regCnt,
      },
    ],
  },
  psFlag: false,
  selectedProperties: [],
  collDeviceModelList: [],
  curCollName: '', //当前采集服务名称
  selectList: [],
  AllPropertyNameList: [],
  pName: '',
  pTableData: [],
  pHeadStyle: {
    color: variables.headerTextColor,
    borderColor: '#C0C4CC',
    height: '48px',
  },
  pCellStyle: {
    borderColor: '#C0C4CC',
    height: '48px',
  },
  pCurrentPage: 1,
  pPagesize: 5,
})

const contentRef = ref(null)

const getModelPropertiesList = (flag) => {
  const pData = {
    token: users.token,
    data: {
      name: props.curTransferModel.name,
    },
  }
  TransferModelApi.getPropertiesByModelIdList(pData).then(async (res) => {
    console.log('getPropertiesByModelIdList -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.PropertyTableData = res.data
      ctxData.AllPropertyNameList = []
      ctxData.PropertyTableData.forEach((item) => {
        ctxData.AllPropertyNameList.push(item.name)
      })
      console.log('AllPropertyNameList ->', ctxData.AllPropertyNameList)
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
      ctxData.tableMaxHeight = contentRef.value.clientHeight - 34 - 36 - 132
    })
  })
}
getModelPropertiesList()
const refresh = () => {
  getModelPropertiesList(1)
}
// 过滤表格数据
const filterTableData = computed(() => {
  let PropertyInfo = ctxData.PropertyInfo
  return ctxData.PropertyTableData.filter(
    (item) =>
      !PropertyInfo ||
      item.name.toLowerCase().includes(PropertyInfo.toLowerCase()) ||
      item.label.toLowerCase().includes(PropertyInfo.toLowerCase())
  ).slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  let PropertyInfo = ctxData.PropertyInfo
  return ctxData.PropertyTableData.filter(
    (item) =>
      !PropertyInfo ||
      item.name.toLowerCase().includes(PropertyInfo.toLowerCase()) ||
      item.label.toLowerCase().includes(PropertyInfo.toLowerCase())
  )
})
// 处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
// 获取设备模型列表
const getDeviceModelList = () => {
  //
  const pData = {
    token: users.token,
    data: {},
  }
  DeviceModelApi.getDeviceModelList(pData).then((res) => {
    console.log('getDeviceModelList -> res = ', res)
    if (res.code === '0') {
      ctxData.collDeviceModelList = res.data
      if (res.data.length > 0) {
        ctxData.curCollName = res.data[0].name
      }
      getPropertiesByCollName()
    } else {
      showOneResMsg(res)
    }
  })
}

//根据采集模型name查询属性列表
const getPropertiesByCollName = () => {
  const pData = {
    token: users.token,
    data: {
      name: ctxData.curCollName,
    },
  }
  DeviceModelApi.getDeviceModelProperty(pData).then((res) => {
    console.log('getDeviceModelProperty -> res', res)
    if (res.code === '0') {
      ctxData.pTableData = res.data
    } else {
      ctxData.pTableData = []
      showOneResMsg(res)
    }
  })
}
const handleSelect = (val) => {
  ctxData.selectList = val
}

//批量添加属性
const addProperties = () => {
  ctxData.plFlag = true
  getDeviceModelList()
}
const selectColl = () => {
  getPropertiesByCollName()
}
// 处理当前页变化
const handlePCurrentChange = (val) => {
  ctxData.pCurrentPage = val
}
// 处理每页大小变化
const handlePSizeChange = (val) => {
  ctxData.pPagesize = val
}
const selectedProp = (row) => {
  return !ctxData.AllPropertyNameList.includes(row.name)
}
const saveProperties = async() => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  console.log('saveProperties')

  if (ctxData.selectList.length > 0) {
    let properties = [];
    ctxData.selectList.map(item => {
      let property = {}
      property['name'] = item.name
      property['label'] = item.label
      property['uploadName'] = item.name
      property['type'] = item.type
      property['decimals'] = item.decimals
      property['unit'] = item.unit
      let params = {}
      if (item.type !== 3) {
        params['min'] = item.params.min
        params['max'] = item.params.max
        params['minMaxAlarm'] = item.params.minMaxAlarm
        params['step'] = item.params.step
        params['stepAlarm'] = item.params.stepAlarm
      } else {
        params['dataLength'] = item.params.dataLength
        params['dataLengthAlarm'] = item.params.dataLengthAlarm
      }
      property['params'] = params

      properties.push(property)
    })

    console.log('saveProperties -> properties', properties)

    const pData = {
      token: users.token,
      data: {
        name: props.curTransferModel.name,
        property: properties,
      },
    }
    const res = await TransferModelApi.addProperties(pData)
    const count = properties.length;
    if (res.data == count) {// 全部操作成功
      loading.close()
      handleResult(res, getModelPropertiesList)
      cancelProperties() // 添加成功，关闭弹窗
    } else if (!res.data) { // 操作失败，同时刷新同步属性的弹窗列表
      loading.close()
      let message = '添加失败';
      ElMessage({
        type: 'error',
        message: message,
      })
      selectColl()
    } else { // 部分成功，同时刷新同步属性的弹窗列表
      loading.close()
      let message = res.data + '个添加成功，' + (count - res.data) + '个添加失败';
      ElMessage({
        type: 'error',
        message: message,
      })
      selectColl()
    }
  }
}
const cancelProperties = () => {
  ctxData.plFlag = false
}
const beforeCloseProperties = (done) => {
  cancelProperties()
}

// 添加属性
const addProperty = () => {
  console.log('addDeviceModelProperty')
  ctxData.pFlag = true
  ctxData.pTitle = '添加属性'
}
// 编辑属性
const editProperty = (row) => {
  ctxData.pFlag = true
  ctxData.pTitle = '编辑属性'
  ctxData.propertyForm.name = row.name
  ctxData.propertyForm.label = row.label
  ctxData.propertyForm.uploadName = row.uploadName
  console.log('editProperty -> row ', row)
  ctxData.propertyForm.type = row.type
  ctxData.propertyForm.decimals = row.decimals
  ctxData.propertyForm.unit = row.unit
  if (row.type !== 3) {
    ctxData.propertyForm['min'] = row.params.min
    ctxData.propertyForm['max'] = row.params.max
    ctxData.propertyForm['minMaxAlarm'] = row.params.minMaxAlarm
    ctxData.propertyForm['step'] = row.params.step
    ctxData.propertyForm['stepAlarm'] = row.params.stepAlarm
  } else {
    ctxData.propertyForm['dataLength'] = row.params.dataLength
    ctxData.propertyForm['dataLengthAlarm'] = row.params.dataLengthAlarm
  }
  console.log('ctxData.propertyForm', ctxData.propertyForm)
}
const propertyFormRef = ref(null)
const submitPorpertyForm = () => {
  if (ctxData.propertyForm.type !== 3 && Number(ctxData.propertyForm.min) > Number(ctxData.propertyForm.max)) {
    ElMessage.warning('最大值必须大于最小值')
    return;
  }
  propertyFormRef.value.validate((valid) => {
    console.log('valid', valid)
    if (valid) {
      //
      let property = {}
      property['name'] = ctxData.propertyForm.name
      property['label'] = ctxData.propertyForm.label
      property['uploadName'] = ctxData.propertyForm.uploadName
      property['type'] = ctxData.propertyForm.type
      property['decimals'] = ctxData.propertyForm.decimals
      property['unit'] = ctxData.propertyForm.unit.toString()
      let params = {}
      if (ctxData.propertyForm.type !== 3) {
        params['min'] = ctxData.propertyForm.min.toString()
        params['max'] = ctxData.propertyForm.max.toString()
        params['minMaxAlarm'] = ctxData.propertyForm.minMaxAlarm
        params['step'] = ctxData.propertyForm.step.toString()
        params['stepAlarm'] = ctxData.propertyForm.stepAlarm
      } else {
        params['dataLength'] = ctxData.propertyForm.dataLength.toString()
        params['dataLengthAlarm'] = ctxData.propertyForm.dataLengthAlarm
      }
      property['params'] = params
      console.log('submitPorpertyForm -> property', property)
      const pData = {
        token: users.token,
        data: {
          name: props.curTransferModel.name,
          property: property,
        },
      }
      if (ctxData.pTitle.includes('添加')) {
        console.log('添加属性')
        TransferModelApi.addProperty(pData).then((res) => {
          handleResult(res, getModelPropertiesList)
          cancelPorpertySubmit()
        })
      } else {
        console.log('编辑属性')
        TransferModelApi.editProperty(pData).then((res) => {
          handleResult(res, getModelPropertiesList)
          cancelPorpertySubmit()
        })
      }
    } else {
      return false
    }
  })
}
const cancelPorpertySubmit = () => {
  ctxData.pFlag = false
  propertyFormRef.value.resetFields()
  initPropertyForm()
}
const handleCloseProperty = (done) => {
  cancelPorpertySubmit()
}

const handleSelectionChange = (val) => {
  ctxData.selectedProperties = val
  console.log('handleSelectionChange -> val =', val)
}
// 删除属性,可以批量删除
const deleteProperty = () => {
  let pList = []
  if (ctxData.selectedProperties.length === 0) {
    ElMessage.info('请至少选择一个属性！')
    return
  } else {
    ctxData.selectedProperties.forEach((item) => {
      pList.push(item.name)
    })
  }
  console.log('pList', pList)
  ElMessageBox.confirm('确认要删除这些属性吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          name: props.curTransferModel.name,
          properties: pList,
        },
      }
      TransferModelApi.deleteProperty(pData).then((res) => {
        handleResult(res, getModelPropertiesList)
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消删除',
      })
    })
}
//初始化接口表单
const initPropertyForm = () => {
  ctxData.propertyForm = {
    name: '', //属性名称
    label: '', // 属性标签
    type: 0, // 数据类型
    //params
    min: 0, // 属性最小值，只有uint32，int32，double有效
    max: 0, // 属性最大值，只有uint32，int32，double有效
    minMaxAlarm: false, // 范围报警使能，只有uint32，int32，double有效
    step: 0, // 步长，只有uint32，int32，double有效
    stepAlarm: false, // 步长报警，只有uint32，int32，double有效
    unit: '', // 单位，只有uint32，int32，double有效
    decimals: 0, // 小数位数，只有double有效
    dataLength: 0, // 字符串长度，只有string有效
    dataLengthAlarm: false, // 字符串长度报警使能，只有string有效
  }
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
.dialog-content-head {
  position: relative;
  width: 100%;
  height: 60px;
  display: flex;
  justify-content: space-between;
  box-sizing: border-box;
  padding-bottom: 24px;
}
.dialog-content-content {
  position: relative;
  width: 100%;
  text-align: center;
  height: 348px;
}
:deep(.dialog-content .el-card__body) {
  position: absolute;
  top: 56px;
  left: 0;
  right: 0;
  bottom: 0;
}
:deep(.el-card__header) {
  height: 55px;
}
</style>
