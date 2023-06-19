<template>
<div class="main-container">
  <div style="display: flex; justify-content: space-between;">
    <div class="title" style="position: relative;width: 40%;justify-content: flex-start;">
      <div class="tName">{{ props.curModelBlock.label }}：命令参数列表</div>
    </div>
    <div class="search-bar" style="text-align:right;">
      <el-form :inline="true" ref="searchFormRef" status-icon label-width="90px">
        <el-form-item style="margin-left: 20px;">
          <el-input style="width: 200px" placeholder="请输入命令参数名称" v-model="ctxData.deviceModelProperty">
            <template #prefix>
              <el-icon class="el-input__icon"><search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" plain class="right-btn" @click="importDPS()">
            <el-icon class="el-input__icon"><download /></el-icon>
            导入命令参数
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" plain class="right-btn" @click="exportDPS()">
            <el-icon class="el-input__icon"><upload /></el-icon>
            导出命令参数
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" bg class="right-btn" @click="addDeviceModelProperty()">
            <el-icon class="btn-icon"> <Icon name="local-add" size="14px" color="#ffffff" /> </el-icon> 添加
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
        <el-form-item>
          <el-button type="danger" bg class="right-btn" @click="deleteDeviceModelProperty()">
            <el-icon class="btn-icon">
              <Icon name="local-delete" size="14px" color="#ffffff" />
            </el-icon>
            删除
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
  <div class="content" ref="contentRef">
    <el-table
      :data="filterDMPTableData"
      :cell-style="ctxData.cellStyle"
      :header-cell-style="ctxData.headerCellStyle"
      style="width: 100%"
      stripe
      :max-height="ctxData.tableMaxHeight"
      @selection-change="handleSelectionChange"
      @row-dblclick="editDeviceModelProperty"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column type="expand" width="60">
        <template #default="scope">
          <div class="param-content">
            <div class="pc-title">
              <div class="pct-info">
                <b> {{ scope.row.name }} </b>
                参数详情
              </div>
            </div>
            <div class="pc-content">
              <div class="param-item" v-for="(item, key, index) of scope.row.params" :key="index">
                <div class="param-value">{{ ctxData.paramName[key] }}：</div>
                <div class="param-name">{{ item || '-' }}</div>
              </div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column sortable prop="name" label="属性名称" width="auto" min-width="150" align="center"> </el-table-column>
      <el-table-column sortable prop="label" label="属性标签" width="auto" min-width="150" align="center"> </el-table-column>
      <el-table-column sortable label="读写属性" width="auto" min-width="100" align="center">
        <template #default="scope">
          {{ ctxData.accessModeNames['am' + scope.row.accessMode] }}
        </template>
      </el-table-column>
      <el-table-column sortable prop="type" label="数据类型" width="auto" min-width="100" align="center">
        <template #default="scope">
          {{ ctxData.typeNames['t' + scope.row.type] }}
        </template>
      </el-table-column>
      <el-table-column sortable label="小数位数" width="auto" min-width="100" align="center">
        <template #default="scope">
          {{ scope.row.decimals === '' ? 0 : scope.row.decimals }}
        </template>
      </el-table-column>
      <el-table-column sortable prop="unit" label="单位" width="auto" min-width="80" align="center" />
      <el-table-column sortable prop="regAddr" label="寄存器地址" width="auto" min-width="120" align="center" />
      <el-table-column sortable prop="regCnt" label="寄存器数量" width="auto" min-width="120" align="center" />
      <el-table-column sortable prop="ruleType" label="解析规则" width="auto" min-width="120" align="center" />
      <el-table-column sortable prop="formula" label="计算公式" width="auto" min-width="200" align="center" />
      <el-table-column sortable label="按位解析" width="auto" min-width="100" align="center">
          <template #default="scope">
            {{ scope.row.bitOffsetSw === false ? '否' : '是' }}
          </template>
        </el-table-column>
        <el-table-column sortable prop="bitOffset" label="位偏移" width="auto" min-width="120" align="center">
          <template #default="scope">
            {{ scope.row.bitOffsetSw === false ? '-' : scope.row.bitOffset }}
          </template>
        </el-table-column>
      <el-table-column sortable prop="identity" label="唯一标识" width="auto" min-width="120" align="center" />
      <el-table-column label="操作" width="auto" min-width="200" align="center" fixed="right">
        <template #default="scope">
          <el-button @click="editDeviceModelProperty(scope.row)" text type="primary">编辑</el-button>
          <!-- <el-button @click="deleteDeviceModelProperty(scope.row)" text type="danger">删除</el-button> -->
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
  <!-- 添加编辑设备命令参数 -->
  <el-dialog
    v-model="ctxData.pFlag"
    :title="ctxData.pTitle"
    width="800px"
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
        :inline="true"
      >
        <el-form-item label="属性名称" prop="name">
          <el-input
            style="width: 220px"
            :disabled="ctxData.pTitle.includes('编辑')"
            type="text"
            v-model="ctxData.propertyForm.name"
            autocomplete="off"
            placeholder="请输入属性名称"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="属性标签" prop="label">
          <el-input
            style="width: 220px"
            type="text"
            v-model="ctxData.propertyForm.label"
            autocomplete="off"
            placeholder="请输入属性标签"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="读写属性" prop="accessMode">
          <el-select v-model="ctxData.propertyForm.accessMode" style="width: 220px" placeholder="请选择读写属性">
            <el-option
              v-for="item in ctxData.accessModeOptions"
              :key="'accessMode_' + item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="单位" prop="unit">
          <el-input
            style="width: 220px"
            type="text"
            v-model="ctxData.propertyForm.unit"
            autocomplete="off"
            placeholder="请输入单位"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="属性类型" prop="type">
          <el-select v-model.number="ctxData.propertyForm.type" style="width: 220px" placeholder="请选择属性类型">
            <el-option
              v-for="item in ctxData.typeOptions"
              :key="'type_' + item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="计算公式" prop="formula">
          <el-input
            style="width: 220px"
            type="text"
            readonly
            v-model="ctxData.propertyForm.formula"
            autocomplete="off"
            placeholder="请编辑计算公式"
          >
            <template #append><el-button @click="showFormula()">编辑</el-button></template>
          </el-input>
        </el-form-item>
        <el-form-item label="小数位数" prop="decimals" v-show="ctxData.propertyForm.type === 2">
          <el-input
            style="width: 220px"
            type="text"
            v-model.number="ctxData.propertyForm.decimals"
            autocomplete="off"
            placeholder="请输入小数位数"
          >
          </el-input>
        </el-form-item>

        <!-- lp add 2023-06-15-->
        <div>
          <el-form-item v-if="ctxData.propertyForm.type !== 3" label="范围报警" prop="minMaxAlarm">
          <el-switch v-model="ctxData.propertyForm.minMaxAlarm" inline-prompt active-text="是" inactive-text="否" />
        </el-form-item>
        </div>
        <div>
        <el-form-item v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.minMaxAlarm" label="最小值" prop="min">
          <el-input type="text" v-model="ctxData.propertyForm.min" autocomplete="off" placeholder="请输入最小值"></el-input>
        </el-form-item>
        </div>
        <el-form-item v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.minMaxAlarm" label="最大值"  prop="max">
          <el-input type="text" v-model="ctxData.propertyForm.max" autocomplete="off" placeholder="请输入最大值"></el-input>
        </el-form-item>

        <div>
        <el-form-item v-if="ctxData.propertyForm.type !== 3" label="步长报警" prop="stepAlarm">
          <el-switch v-model="ctxData.propertyForm.stepAlarm" inline-prompt active-text="是" inactive-text="否" />
        </el-form-item>
        </div>
        <el-form-item v-if="ctxData.propertyForm.type !== 3 && ctxData.propertyForm.stepAlarm" label="步长" prop="step">
          <el-input type="text" v-model="ctxData.propertyForm.step" autocomplete="off" placeholder="请输入步长"></el-input>
        </el-form-item>

        <div>
        <el-form-item v-if="ctxData.propertyForm.type === 3" label="字符串长度报警" prop="dataLengthAlarm">
          <el-switch v-model="ctxData.propertyForm.dataLengthAlarm" inline-prompt active-text="是" inactive-text="否" />
        </el-form-item>
        </div>
        <el-form-item v-if="ctxData.propertyForm.type === 3 && ctxData.propertyForm.dataLengthAlarm" label="字符串长度" prop="dataLength">
          <el-input type="text" v-model="ctxData.propertyForm.dataLength" autocomplete="off" placeholder="请输入字符串长度" ></el-input>
        </el-form-item>

        <div class="form-title"><div class="tName">配置参数</div></div>

        <el-form-item label="寄存器地址" prop="regAddr">
          <el-input
            type="text"
            style="width: 220px"
            v-model.number="ctxData.propertyForm.regAddr"
            autocomplete="off"
            placeholder="请输入寄存器地址"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="" prop="" style="width: 220px" v-if="props.curModelBlock.funCode != 1"> </el-form-item>
        <el-form-item label="寄存器数量" prop="regCnt">
          <el-input
            type="text"
            v-if="props.curModelBlock.funCode == 1"
            style="width: 220px"
            disabled
            v-model.number="ctxData.propertyForm.regCnt"
            autocomplete="off"
            placeholder="请输入寄存器数量"
          >
          </el-input>
          <el-select
            v-else
            v-model.number="ctxData.propertyForm.regCnt"
            style="width: 220px"
            @change="changeRegCnt()"
            placeholder="请选择寄存器数量"
          >
            <el-option
              v-for="item in ctxData.regCntOptions"
              :key="'regCount_' + item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="解析规则" prop="ruleType" v-if="props.curModelBlock.funCode != 1">
          <el-select v-model.number="ctxData.propertyForm.ruleType" style="width: 220px" placeholder="请选择解析规则">
            <el-option
              v-for="item in ctxData.ruleTypeOptions['rt' + ctxData.propertyForm.regCnt]"
              :key="'type_' + item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-row>
            <el-form-item v-if="props.curModelBlock.funCode == 3 || props.curModelBlock.funCode == 4" label="按位解析">
              <el-tooltip class="item" effect="dark" content="开启后按字节位解析数据,位偏移从0开始" placement="top">
              <el-switch style="width: 220px" v-model="ctxData.propertyForm.bitOffsetSw" inline-prompt active-text="是" inactive-text="否" />
              </el-tooltip>
            </el-form-item>
            <el-form-item label="位偏移" v-if="ctxData.propertyForm.bitOffsetSw" prop="bitOffset">
              <el-input
                type="text"
                style="width: 220px"
                v-model.number="ctxData.propertyForm.bitOffset"
                autocomplete="off"
                placeholder="请输入位偏移"
              >
              </el-input>
            </el-form-item>
        </el-row>
        <!-- <el-form-item label="步长" v-if="props.curModelBlock.funCode == 3 || props.curModelBlock.funCode == 4" prop="step">
          <el-tooltip class="item" effect="dark" content="数据变化超过所配步长则变化上送" placement="top">
            <el-input
              type="text"
              style="width: 220px"
              v-model="ctxData.propertyForm.step"
              autocomplete="off"
              placeholder="请输入步长"
            >
            </el-input>
          </el-tooltip>
        </el-form-item> -->

        <el-form-item label="唯一标识">
            <el-input
              type="text"
              style="width: 220px"
              v-model="ctxData.propertyForm.identity"
              autocomplete="off"
              placeholder="请输入唯一标识"
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
  <!-- 编辑计算公式 -->
  <el-dialog v-model="ctxData.fFlag" title="编辑计算公式" width="600px">
    <div class="dialog-content">
      <el-form :model="ctxData.formulaForm" label-position="right" label-width="120px">
        <el-form-item label="运算符号">
          <div style="width: 100%">
            <div v-for="(item, index) in ctxData.operationList" :key="index" class="operation">
              <div style="padding-right: 10px">
                <el-button style="width: 100%" @click="setOperation(item)">{{ item.label }}</el-button>
              </div>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="设备属性">
          <el-select
            v-model="ctxData.pItem"
            filterable
            placeholder="请选择设备属性"
            @change="setProperty(ctxData.pItem)"
            style="width: 100%"
          >
            <el-option v-for="item in ctxData.propertyInfo" :key="item.name" :label="item.label" :value="item.name" />
          </el-select>
        </el-form-item>

        <el-form-item label="当前公式">
          <el-input
            type="text"
            clearable
            v-model="ctxData.curFormula"
            autocomplete="off"
            placeholder="请编辑计算公式"
            @blur="onBulr"
          ></el-input>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="cancelSave()">取消</el-button>
        <el-button type="primary" @click="saveFormula()">保存</el-button>
      </span>
    </template>
  </el-dialog>
  <!-- 导入命令参数 -->
  <el-dialog
    v-model="ctxData.psFlag"
    title="导入命令参数"
    width="600px"
    :before-close="beforeCloseUploadDPS"
    :close-on-click-modal="false"
  >
    <el-upload
      ref="uploadDPSRef"
      action=""
      :auto-upload="false"
      :http-request="myRequestDPS"
      :limit="1"
      :on-exceed="handleExceedDPS"
      :before-upload="beforeUploadDPS"
    >
      <el-button type="primary">选择文件</el-button>
      <template #tip>
        <div class="el-upload__tip">只能上传一个文件，只支持xlsx格式文件！</div>
      </template>
    </el-upload>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="cancelUploadDPS">取消</el-button>
        <el-button type="primary" @click="submitUploadDPS">上传</el-button>
      </span>
    </template>
  </el-dialog>
</div>
</template>
<script setup>
import { Search, Back, Download, Upload } from '@element-plus/icons-vue'
import variables from 'styles/variables.module.scss'
import ModelBlockApi from 'api/modelBlock.js'
import DeviceModelApi from 'api/deviceModel.js'
import { userStore } from 'stores/user'
const users = userStore()

const props = defineProps({
  curModelBlock: {
    type: Object,
    default: {},
  },
})
// 返回设备块
const emit = defineEmits(['changeBShowFlag'])
const toModelBlock = () => {
  emit('changeBShowFlag')
}

const regCnt = /^[0-9]*[1-9][0-9]*$/
const validateRegCnt = (rule, value, callback) => {
  if (!regCnt.test(value)) {
    callback(new Error('只能输入正整数数字！'))
  } else {
    callback()
  }
}
const validateRegAddr = (rule, value, callback) => {
  if (value !== 0 && !regCnt.test(value)) {
    callback(new Error('只能输入自然数！'))
  } else {
    callback()
  }
}
const regOffset = /^[0-9]\d*$/
const validateBitOffset = (rule, value, callback) => {
  if (!regOffset.test(value) || value > 65) {
    callback(new Error('只能输入0-64的整数数字！'))
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
const contentRef = ref(null)
const ctxData = reactive({
  formulaForm: {},
  headerCellStyle: {
    background: variables.primaryColor,
    color: variables.fontWhiteColor,
    height: '48px',
  },
  cellStyle: {
    height: '42px',
    padding: '2px 0',
  },
  tableMaxHeight: 0,
  currentPage: 1, // 默认当前页是第一页
  pagesize: 20, // 每页数据个数
  propertyList: [],
  deviceModelProperty: '',
  accessModeOptions: [
    //读写模式
    { label: '只读', value: 0 },
    { label: '只写', value: 1 },
    { label: '读写', value: 2 },
  ],
  typeOptions: [
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
  accessModeNames: {
    am0: '只读',
    am1: '只写',
    am2: '读写',
  },
  regCntOptions: [
    { label: 1, value: 1 },
    { label: 2, value: 2 },
    { label: 4, value: 4 },
  ],
  pFlag: false, //属性对话框标识
  pTitle: '添加属性', //属性对话框titleName
  propertyForm: {
    name: '', // 属性名称，只能是字母+数字的组合，不可以是中文
    label: '', // 属性标签
    accessMode: 0, // 读写属性
    type: 0,
    //params
    regCnt: 1,
    regAddr: null,
    ruleType: 'Int_AB',
    unit: '', // 单位，只有uint32，int32，double有效
    decimals: 0, // 小数位数，只有double有效
    formula: '', // 计算公式
    bitOffsetSw: false, // 按位解析
    bitOffset: 0, // 位偏移


    min: 0, // 属性最小值，只有uint32，int32，double有效
    max: 0, // 属性最大值，只有uint32，int32，double有效
    minMaxAlarm: false, // 范围报警，只有uint32，int32，double有效
    step: 0, // 步长，只有uint32，int32，double有效
    stepAlarm: false, // 步长报警，只有uint32，int32，double有效
    dataLength: 0, // 字符串长度，只有string有效
    dataLengthAlarm: false, // 字符串长度报警，只有string有效
    identity: '', // 唯一标识
  },

  //数据类型
  ruleTypeOptions: {
    rt1: [
      { label: 'Int_AB', value: 'Int_AB' },
      { label: 'Int_BA', value: 'Int_BA' },
    ],
    rt2: [
      { label: 'Long_ABCD', value: 'Long_ABCD' },
      { label: 'Long_BADC', value: 'Long_BADC' },
      { label: 'Long_DCBA', value: 'Long_DCBA' },
      { label: 'Long_CDAB', value: 'Long_CDAB' },
      { label: 'Float_ABCD', value: 'Float_ABCD' },
      { label: 'Float_BADC', value: 'Float_BADC' },
      { label: 'Float_DCBA', value: 'Float_DCBA' },
      { label: 'Float_CDAB', value: 'Float_CDAB' },
    ],
    rt4: [
      { label: 'Double_ABCDEFGH', value: 'Double_ABCDEFGH' },
      { label: 'Double_GHEFCDAB', value: 'Double_GHEFCDAB' },
      { label: 'Double_BADCFEHG', value: 'Double_BADCFEHG' },
      { label: 'Double_HGFEDCBA', value: 'Double_HGFEDCBA' },
    ],
  },
  paramName: {
    name: '属性名称',
    label: '属性标签',
    accessMode: '读写模式',
    type: '属性类型',
    decimals: '小数位数',
    regCnt: '寄存器数量',
    ruleType: '解析规则',
    regAddr: '寄存器地址',
    formula: '计算公式',
    min: '最小值',
    max: '最大值',
    minMaxAlarm: '范围报警',
    step: '步长',
    stepAlarm: '步长报警',
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
    accessMode: [
      {
        required: true,
        message: '读写属性不能为空',
        trigger: 'blur',
      },
    ],
    type: [
      {
        required: true,
        message: '属性类型不能为空',
        trigger: 'blur',
      },
    ],
    ruleType: [
      {
        required: true,
        message: '解析规则不能为空',
        trigger: 'blur',
      },
    ],
    decimals: [
      {
        type: 'number',
        message: '小数位数只能输入数字',
      },
    ],
    regCnt: [
      {
        required: true,
        message: '寄存器数量不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateRegCnt,
      },
    ],
    regAddr: [
      {
        required: true,
        message: '寄存器地址不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateRegAddr,
      },
    ],
    bitOffset: [
      {
        required: true,
        message: '位偏移不能为空',
        trigger: 'blur',
      },
      {
        trigger: 'blur',
        validator: validateBitOffset,
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
  fFlag: false,
  propertyInfo: [],
  operationList: [
    { name: '+', label: '+', describe: '加' },
    { name: '-', label: '-', describe: '减' },
    { name: '*', label: '*', describe: '乘' },
    { name: '/', label: '/', describe: '除' },
    { name: '%', label: '%', describe: '模' },
    { name: '^', label: '^', describe: '次幂' },
    { name: '(', label: '(', describe: '左括号' },
    { name: ')', label: ')', describe: '右括号' },
  ],
  curFormula: '',
  selectionStart: 0,
  pItem: '',
})
// 获取设备命令参数
const getDeviceModelBlockProperty = (flag) => {
  const pData = {
    token: users.token,
    data: {
      tslName: props.curModelBlock.tslName,
      cmdName: props.curModelBlock.name,
    },
  }
  ModelBlockApi.getDeviceModelBlockProperty(pData).then(async (res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.propertyList = res.data == null ? [] : res.data
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
      ctxData.tableMaxHeight = contentRef.value.clientHeight - 34 - 22 - 82
    })
  })
}

getDeviceModelBlockProperty()
const refresh = () => {
  getDeviceModelBlockProperty(1)
}
const filterDMPTableData = computed(() => {
  return ctxData.propertyList
    .filter((item) => {
      var a = !ctxData.deviceModelProperty
      var b = item.name.toLowerCase().includes(ctxData.deviceModelProperty.toLowerCase())

      return a || b
    })
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
})
const filterTableDataPage = computed(() => {
  return ctxData.propertyList.filter((item) => {
    var a = !ctxData.deviceModelProperty
    var b = item.name.toLowerCase().includes(ctxData.deviceModelProperty.toLowerCase())

    return a || b
  })
})
// 处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
// 处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
// 导入设备属性和服务
const importDPS = () => {
  ctxData.psFlag = true
}
// 导出设备属性和服务
const exportDPS = () => {
  const pData = {
    token: users.token,
    responseType: 'blob',
    data: {
      tslName: props.curModelBlock.tslName,
      cmdName: props.curModelBlock.name,
    },
  }
  ModelBlockApi.exportDeviceModelBlockProperty(pData).then((res) => {
    if (res && res.code === '1') {
      showOneResMsg(res)
      return
    }
    const blob = new Blob([res.blob])
    const fileName = res.fileName

    const elink = document.createElement('a')
    elink.download = fileName
    elink.style.display = 'none'
    elink.href = URL.createObjectURL(blob)
    document.body.appendChild(elink)
    elink.click()
    URL.revokeObjectURL(elink.href) // 释放URL 对象
    document.body.removeChild(elink)
  })
}
/**
 * 提交上传的插件
 */
const submitUploadDPS = () => {
  uploadDPSRef.value && uploadDPSRef.value.submit()
}

// 取消插件自带的xhr
const myRequestDPS = (obj) => {}
const uploadDPSRef = ref(null)
/**
 * 上传文件大于limit时事件
 * @param {要上传的文件} filesimportPlugin
 */
const handleExceedDPS = (files) => {
  uploadDPSRef.value.clearFiles()
  //超过limit取第一个
  const file = files[0]
  uploadDPSRef.value.handleStart(file)
}
/**
 * 文件上传前事件
 * @param {要上传的文件} file
 */
const beforeUploadDPS = (file) => {
  const fileTypeList = ['application/vnd.openxmlformats-officedocument.spreadsheetml.sheet']
  let typeFlag = fileTypeList.includes(file.type)
  if (!typeFlag) {
    ElMessage({
      type: 'error',
      message: '文件格式不正确,必须是xlsx文件！',
    })
    return
  }
  // 调上传接口
  if (typeFlag) {
    let formData = new FormData()
    formData.append('tslName', props.curModelBlock.tslName)
    formData.append('cmdName', props.curModelBlock.name)
    formData.append('fileName', file)
    const pData = {
      token: users.token,
      contentType: 'multipart/form-data',
      data: formData,
    }
    ModelBlockApi.importDeviceModelBlockProperty(pData).then((res) => {
      if (res.code === '0') {
        ElMessage.success(res.message)
        ctxData.psFlag = false
        getDeviceModelBlockProperty()
      } else {
        showOneResMsg(res)
      }
    })
  }
}
//取消上传属性和服务
const cancelUploadDPS = () => {
  ctxData.psFlag = false
  uploadRef.value.clearFiles()
}
// 弹窗右上角关闭事件处理
const beforeCloseUploadDPS = () => {
  cancelUploadDPS()
}

// 添加属性
const addDeviceModelProperty = () => {
  ctxData.pFlag = true
  ctxData.pTitle = '添加属性'
}
// 编辑属性
const editDeviceModelProperty = (row) => {
  ctxData.pFlag = true
  ctxData.pTitle = '编辑属性'
  ctxData.propertyForm.name = row.name
  ctxData.propertyForm.label = row.label
  ctxData.propertyForm.accessMode = row.accessMode
  ctxData.propertyForm.type = row.type
  ctxData.propertyForm.decimals = row.decimals
  ctxData.propertyForm.unit = row.unit
  ctxData.propertyForm.regCnt = row.regCnt
  ctxData.propertyForm.regAddr = row.regAddr
  ctxData.propertyForm.ruleType = row.ruleType
  ctxData.propertyForm.formula = row.formula === undefined || row.formula === null ? '' : row.formula
  ctxData.propertyForm.bitOffsetSw = row.bitOffsetSw === undefined || row.bitOffsetSw === null ? false : row.bitOffsetSw
  ctxData.propertyForm.bitOffset = row.bitOffset === undefined || row.bitOffset === null ? 0 : row.bitOffset
  // ctxData.propertyForm.step = row.step === undefined || row.step === null ? 0 : row.step

  if (row.type !== 3) {
    ctxData.propertyForm['min'] = row.params.min
    ctxData.propertyForm['max'] = row.params.max
    ctxData.propertyForm['minMaxAlarm'] = row.params.minMaxAlarm
    ctxData.propertyForm['step'] = row.params.step === undefined || row.params.step === null ? 0 : row.params.step
    ctxData.propertyForm['stepAlarm'] = row.params.stepAlarm
  } else {
    ctxData.propertyForm['dataLength'] = row.params.dataLength
    ctxData.propertyForm['dataLengthAlarm'] = row.params.dataLengthAlarm
  }
  ctxData.propertyForm.identity = row.identity === undefined || row.identity === null ? '' : row.identity
}
const propertyFormRef = ref(null)
const submitPorpertyForm = () => {
  if (ctxData.propertyForm.type !== 3 && Number(ctxData.propertyForm.min) > Number(ctxData.propertyForm.max)) {
    ElMessage.warning('最大值必须大于最小值')
    return;
  }
  propertyFormRef.value.validate((valid) => {
    if (valid) {
      const pData = {
        token: users.token,
        data: {
          tslName: props.curModelBlock.tslName,
          cmdName: props.curModelBlock.name,
          name: ctxData.propertyForm.name,
          label: ctxData.propertyForm.label,
          accessMode: ctxData.propertyForm.accessMode,
          type: ctxData.propertyForm.type,
          decimals: ctxData.propertyForm.decimals,
          unit: ctxData.propertyForm.unit,
          regCnt: ctxData.propertyForm.regCnt,
          regAddr: ctxData.propertyForm.regAddr,
          ruleType: ctxData.propertyForm.ruleType,
          formula: ctxData.propertyForm.formula,
          bitOffsetSw: ctxData.propertyForm.bitOffsetSw,
          bitOffset: ctxData.propertyForm.bitOffsetSw ? ctxData.propertyForm.bitOffset.toString() : '0',
          // step: +ctxData.propertyForm.step,
          identity: ctxData.propertyForm.identity,
        },
      }
      let params = {}
      if (ctxData.propertyForm.type !== 3) {
        params['min'] = ctxData.propertyForm.min.toString()
        params['max'] = ctxData.propertyForm.max.toString()
        params['minMaxAlarm'] = ctxData.propertyForm.minMaxAlarm
        //params['step'] = +ctxData.propertyForm.step
        params['step'] = ctxData.propertyForm.step.toString()   //ltg del 2023-06-15
        params['stepAlarm'] = ctxData.propertyForm.stepAlarm
      } else {
        params['dataLength'] = ctxData.propertyForm.dataLength.toString()
        params['dataLengthAlarm'] = ctxData.propertyForm.dataLengthAlarm
      }
      pData.data['params'] = params
      if (ctxData.pTitle.includes('添加')) {
        ModelBlockApi.addDeviceModelBlockProperty(pData).then((res) => {
          handleResult(res, getDeviceModelBlockProperty)
          cancelPorpertySubmit()
        })
      } else {
        ModelBlockApi.editDeviceModelBlockProperty(pData).then((res) => {
          handleResult(res, getDeviceModelBlockProperty)
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
}
// 删除属性,可以批量删除
const deleteDeviceModelProperty = (row) => {
  let pList = []
  if (ctxData.selectedProperties.length === 0) {
    ElMessage.info('请至少选择一个属性！')
    return
  } else {
    ctxData.selectedProperties.forEach((item) => {
      pList.push(item.name)
    })
  }
  ElMessageBox.confirm('确认要删除这些属性吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          tslName: props.curModelBlock.tslName,
          cmdName: props.curModelBlock.name,
          propertyNames: pList,
        },
      }
      ModelBlockApi.deleteDeviceModelBlockProperty(pData).then((res) => {
        handleResult(res, getDeviceModelBlockProperty)
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
    accessMode: 0, // 读写属性
    type: 0,
    decimals: 0,
    unit: '',
    regCnt: 1,
    regAddr: null,
    ruleType: 'Int_AB',
    formula: '',
    bitOffsetSw: false,
    bitOffset: 0,
    step: 0,
    min: 0,
    max: 0,
    dataLength: 0,
    identity: '',
  }
}
const getDeviceProperty = () => {
  const pData = {
    token: users.token,
    data: {
      name: props.curModelBlock.tslName,
    },
  }
  DeviceModelApi.getDeviceModelProperty(pData).then((res) => {
    if (!res) return
    if (res.code === '0') {
      ctxData.propertyInfo = [{ name: 't', label: '当前属性' }]
      res.data.forEach((item) => {
        ctxData.propertyInfo.push({
          name: item.name,
          label: item.label,
        })
      })
      //ctxData.propertyInfo = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
const showFormula = () => {
  ctxData.fFlag = true
  ctxData.curFormula = ctxData.propertyForm.formula
  getDeviceProperty()
}
const setOperation = (item) => {
  let length = ctxData.curFormula.length
  if (length !== ctxData.selectionStart) {
    let a = ctxData.curFormula.substr(0, ctxData.selectionStart)
    let b = ctxData.curFormula.substr(ctxData.selectionStart, length - ctxData.selectionStart)
    ctxData.curFormula = a + item.name + b
    ctxData.selectionStart += item.name.length
  } else {
    ctxData.curFormula += item.name
    ctxData.selectionStart += item.name.length
  }
}
const setProperty = (item) => {
  let length = ctxData.curFormula.length
  if (length !== ctxData.selectionStart) {
    let a = ctxData.curFormula.substr(0, ctxData.selectionStart)
    let b = ctxData.curFormula.substr(ctxData.selectionStart, length - ctxData.selectionStart)
    ctxData.curFormula = a + '(' + item + ')' + b
    ctxData.selectionStart += item.length + 2
  } else {
    ctxData.curFormula += '(' + item + ')'
    ctxData.selectionStart += item.length + 2
  }
}
const onBulr = (obj) => {
  ctxData.selectionStart = obj.srcElement.selectionStart
}

const saveFormula = () => {
  ctxData.fFlag = false
  ctxData.propertyForm.formula = ctxData.curFormula
}
const cancelSave = () => {
  ctxData.fFlag = false
  ctxData.curFormula = ''
}
const changeRegCnt = () => {
  const val = ctxData.propertyForm.regCnt
  if (val === 1) {
    ctxData.propertyForm.ruleType = 'Int_AB'
  }
  if (val === 2) {
    ctxData.propertyForm.ruleType = 'Long_ABCD'
  }
  if (val === 4) {
    ctxData.propertyForm.ruleType = 'Double_ABCDEFGH'
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
//非第一次
watch(
  () => props.curModelBlock,
  (newVal, oldVal) => {
    getDeviceModelBlockProperty()
  }
)
</script>
<style lang="scss" scoped>
@use 'styles/custom-scoped.scss' as *;
.form-title {
  margin-bottom: 12px;
}
.tName {
  line-height: 14px;
  font-size: 14px;
  border-left: 3px solid #3054eb;
  padding-left: 15px;
}
.operation {
  display: inline-block;
  width: 25%;
  height: 46px;
}
.main-container .pagination {
  bottom: 12px;
}
</style>
