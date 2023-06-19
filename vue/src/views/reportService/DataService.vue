<template>
  <div class="main-container">
    <div class="main" v-if="ctxData.dnFlag === 1">
      <div class="search-bar">
        <el-form :inline="true" ref="searchFormRef" status-icon label-width="120px">
          <el-form-item label="上报服务名称">
            <el-input style="width: 200px" placeholder="请输入上报服务名称" clearable v-model="ctxData.dataServiceName">
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
        <el-button type="primary" bg class="right-btn" @click="addGateway()">
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
          style="width: 100%"
          :max-height="ctxData.tableMaxHeight"
          stripe
          @row-dblclick="editGateway"
        >
          <el-table-column type="expand" min-width="80">
            <template #default="scope">
              <div class="param-content">
                <div class="pc-title">
                  <div class="pct-info">
                    <b> {{ scope.row.serviceName }} </b>
                    参数详情
                  </div>
                </div>
                <div class="pc-content">
                  <div class="param-item" v-for="(item, key, index) of scope.row.param" :key="index">
                    <div class="param-value">{{ ctxData.paramName[key] }}：</div>
                    <div class="param-name">{{ item || '-' }}</div>
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column sortable prop="serviceName" label="服务名称" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column sortable prop="protocol" label="协议名称" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column sortable prop="reportNetSW" label="指定上报通道" width="auto" min-width="100" align="center">
            <template #default="scope">
              {{ scope.row.reportNetSW === false ? '否' : '是' }}
            </template>
          </el-table-column>
          <el-table-column sortable prop="reportNet" label="上报网卡" width="auto" min-width="100" align="center">
            <template #default="scope">
              {{ scope.row.reportNetSW === false ? '-' : scope.row.reportNet }}
            </template>
          </el-table-column>
          <el-table-column sortable prop="ip" label="上报地址" width="auto" min-width="150" align="center"> </el-table-column>
          <el-table-column sortable prop="port" label="上报端口" width="auto" min-width="100" align="center"> </el-table-column>
          <el-table-column sortable prop="reportTime" label="上报周期（秒）" width="auto" min-width="150" align="center">
          </el-table-column>
          <el-table-column sortable label="在线状态" width="auto" min-width="100" align="center">
            <template #default="scope">
              {{ scope.row.reportStatus === 'onLine' ? '在线' : '离线' }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="auto" min-width="200" align="center" fixed="right">
            <template #default="scope">
              <el-button v-if="scope.row.protocol !== 'ModbusTCP'" @click="showGateway(scope.row)" text type="success">
                上报节点
              </el-button>
              <el-button v-if="scope.row.protocol == 'ModbusTCP'" @click="showRegInfo(scope.row)" text type="success">
                寄存器详情
              </el-button>
              <el-button @click="editGateway(scope.row)" text type="primary">编辑</el-button>
              <el-button @click="deleteGateway(scope.row)" text type="danger">删除</el-button>
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
    <NodeService
      v-if="ctxData.dnFlag === 2"
      :curGateway="ctxData.curGateway"
      @changeDnFlag="changeDnFlag"
      style="width: 100%; height: 100%;overflow:hidden;"
    ></NodeService>

    <RegInfo 
      v-if="ctxData.dnFlag === 3" 
      :curGateway="ctxData.curGateway" 
      @changeDnFlag="changeDnFlag"
      style="width: 100%; height: 100%;overflow:hidden;"
    ></RegInfo>

    <el-dialog
      v-model="ctxData.dFlag"
      :title="ctxData.dTitle"
      width="600px"
      :before-close="handleClose"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-form
          :model="ctxData.gatewayForm"
          :rules="ctxData.gatewayRules"
          ref="gatewayFormRef"
          status-icon
          label-position="right"
          label-width="120px"
          style="height: 450px"
        >
          <el-form-item label="服务名称" prop="serviceName">
            <el-input
              type="text"
              v-model="ctxData.gatewayForm.serviceName"
              autocomplete="off"
              placeholder="请输入服务名称"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="协议类型" prop="protocol">
            <el-select
              v-model="ctxData.gatewayForm.protocol"
              @change="changeProtocol(ctxData.gatewayForm.protocol)"
              style="width: 100%"
              placeholder="请选择协议类型"
            >
              <el-option
                v-for="item in ctxData.protocolOptions"
                :key="'protocol' + item.name"
                :label="item.label"
                :value="item.name"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="指定上报通道">
            <el-switch v-model="ctxData.gatewayForm.reportNetSw" inline-prompt active-text="是" inactive-text="否" />
          </el-form-item>
          <el-form-item label="上报网卡" v-if="ctxData.gatewayForm.reportNetSw" prop="reportNet">
            <el-input type="text" v-model.number="ctxData.gatewayForm.reportNet" autocomplete="off" placeholder="请输入上报网卡"></el-input>
          </el-form-item>
          <el-form-item
            label="上报地址"
            :prop="ctxData.gatewayForm.protocol !== 'ModbusTCP' ? 'ip' : ''"
            v-show="ctxData.gatewayForm.protocol !== 'ModbusTCP'"
          >
            <el-input type="text" v-model="ctxData.gatewayForm.ip" autocomplete="off" placeholder="请输入上报地址">
            </el-input>
          </el-form-item>
          <el-form-item label="上报端口" prop="port">
            <el-input type="text" v-model="ctxData.gatewayForm.port" autocomplete="off" placeholder="请输入上报端口">
            </el-input>
          </el-form-item>
          <el-form-item
            label="上报周期"
            :prop="ctxData.gatewayForm.protocol !== 'ModbusTCP' ? 'reportTime' : ''"
            v-show="ctxData.gatewayForm.protocol !== 'ModbusTCP'"
          >
            <el-input
              type="text"
              v-model.number="ctxData.gatewayForm.reportTime"
              autocomplete="off"
              placeholder="请输入上报周期"
            >
              <template #append>单位秒</template>
            </el-input>
          </el-form-item>

          <el-form-item v-if="ctxData.gatewayForm.protocol.includes('MQTT')" label="用户名" prop="userName">
            <el-input type="text" v-model="ctxData.gatewayForm.userName" autocomplete="off" placeholder="请输入用户名">
            </el-input>
          </el-form-item>

          <el-form-item v-if="ctxData.gatewayForm.protocol.includes('MQTT')" label="密码" prop="password">
            <el-input type="text" v-model="ctxData.gatewayForm.password" autocomplete="off" placeholder="请输入密码">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.gatewayForm.protocol.includes('MQTT')" label="客户端名称" prop="clientID">
            <el-input
              type="text"
              v-model="ctxData.gatewayForm.clientID"
              autocomplete="off"
              placeholder="请输入客户端名称"
            >
              <template #append><el-button @click="getGatewaySN()">同步SN</el-button></template>
            </el-input>
          </el-form-item>

          <el-form-item v-if="ctxData.gatewayForm.protocol.includes('FSJY')" label="APP-KEY" prop="appKey">
            <el-input type="text" v-model="ctxData.gatewayForm.appKey" autocomplete="off" placeholder="请输入appKey">
            </el-input>
          </el-form-item>

          <el-form-item v-if="ctxData.gatewayForm.protocol.includes('FSJY')" label="产品密钥" prop="productKey">
            <el-input type="text" v-model="ctxData.gatewayForm.productKey" autocomplete="off" placeholder="请输入网关产品密钥">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.gatewayForm.protocol.includes('FSJY')" label="通讯地址" prop="deviceID">
            <el-input type="text" v-model="ctxData.gatewayForm.deviceID" autocomplete="off" placeholder="请输入网关通讯地址">
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.gatewayForm.protocol.includes('FSJY')" label="设备密钥" prop="deviceSecret">
            <el-input type="text" v-model="ctxData.gatewayForm.deviceSecret" autocomplete="off" placeholder="请输入网关设备密钥">
            </el-input>
          </el-form-item>


          <el-form-item v-if="ctxData.gatewayForm.protocol === 'RT.STLV2'" label="项目编码" prop="projectCode">
            <el-input
              type="text"
              v-model="ctxData.gatewayForm.projectCode"
              autocomplete="off"
              placeholder="请输入项目编码"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.dnList.includes(ctxData.gatewayForm.protocol)" label="通信编码" prop="deviceName">
            <el-input
              type="text"
              v-model="ctxData.gatewayForm.deviceName"
              autocomplete="off"
              placeholder="请输入通信编码"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.kaList.includes(ctxData.gatewayForm.protocol)" label="保活时间" prop="keepAlive">
            <el-input
              type="text"
              v-model="ctxData.gatewayForm.keepAlive"
              autocomplete="off"
              placeholder="请输入保活时间"
            >
              <template #append>单位秒</template>
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.gatewayForm.protocol.includes('MQTT')" label="是否清除会话" prop="cleanSession">
            <el-switch v-model="ctxData.gatewayForm.cleanSession" inline-prompt active-text="是" inactive-text="否" />
          </el-form-item>
          <el-form-item v-if="ctxData.gatewayForm.protocol === 'RT.STLV2'" label="秘钥" prop="token">
            <el-input type="text" v-model="ctxData.gatewayForm.token" autocomplete="off" placeholder="请输入秘钥">
            </el-input>
          </el-form-item>
          <!-- ModbusTCP -->
          <el-form-item v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'" label="从机ID" prop="slaveID">
            <el-input
              type="text"
              v-model.number="ctxData.gatewayForm.slaveID"
              autocomplete="off"
              placeholder="请输入从机ID"
            >
            </el-input>
          </el-form-item>
          <el-divider content-position="center" v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'">
            线圈状态寄存器
          </el-divider>
          <el-form-item v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'" label="起始地址" prop="coilStatusRegStart">
            <el-input
              type="text"
              v-model.number="ctxData.gatewayForm.coilStatusRegStart"
              autocomplete="off"
              placeholder="请输入线圈状态寄存器起始地址"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'" label="数量" prop="coilStatusRegCnt">
            <el-input
              type="text"
              v-model.number="ctxData.gatewayForm.coilStatusRegCnt"
              autocomplete="off"
              placeholder="请输入线圈状态寄存器数量"
            >
            </el-input>
          </el-form-item>

          <el-divider content-position="center" v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'">
            离散输入状态寄存器
          </el-divider>
          <el-form-item v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'" label="起始地址" prop="inputStatusRegStart">
            <el-input
              type="text"
              v-model.number="ctxData.gatewayForm.inputStatusRegStart"
              autocomplete="off"
              placeholder="请输入离散输入状态寄存器起始地址"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'" label="数量" prop="inputStatusRegCnt">
            <el-input
              type="text"
              v-model.number="ctxData.gatewayForm.inputStatusRegCnt"
              autocomplete="off"
              placeholder="请输入离散输入状态寄存器数量"
            >
            </el-input>
          </el-form-item>
          <el-divider content-position="center" v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'">
            保持寄存器
          </el-divider>
          <el-form-item v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'" label="起始地址" prop="holdingRegStart">
            <el-input
              type="text"
              v-model.number="ctxData.gatewayForm.holdingRegStart"
              autocomplete="off"
              placeholder="请输入保持寄存器起始地址"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'" label="数量" prop="holdingRegCnt">
            <el-input
              type="text"
              v-model.number="ctxData.gatewayForm.holdingRegCnt"
              autocomplete="off"
              placeholder="请输入保持寄存器数量"
            >
            </el-input>
          </el-form-item>
          <el-divider content-position="center" v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'">
            输入寄存器
          </el-divider>
          <el-form-item v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'" label="起始地址" prop="inputRegStart">
            <el-input
              type="text"
              v-model.number="ctxData.gatewayForm.inputRegStart"
              autocomplete="off"
              placeholder="请输入输入寄存器起始地址"
            >
            </el-input>
          </el-form-item>
          <el-form-item v-if="ctxData.gatewayForm.protocol === 'ModbusTCP'" label="数量" prop="inputRegCnt">
            <el-input
              type="text"
              v-model.number="ctxData.gatewayForm.inputRegCnt"
              autocomplete="off"
              placeholder="请输入输入寄存器数量"
            >
            </el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelSubmit()">取消</el-button>
          <el-button type="primary" @click="submitGatewayForm()">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { Search } from '@element-plus/icons-vue'
import variables from 'styles/variables.module.scss'
import ServiceApi from 'api/service.js'
import { userStore } from 'stores/user'
import { useRoute } from 'vue-router'
import NodeService from './dataService/NodeService.vue'
import RegInfo from './dataService/RegInfo.vue'
import ProductApi from 'api/product.js'
import {ElMessageBox} from "element-plus";
const users = userStore()
console.log('window.location', window.location)
const refExpIP =
  /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/
const refExpYM =
  /^(?=^.{3,255}$)(http(s)?:\/\/)?(www\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:\d+)*(\/\w+\.\w+)*$/
const validateIP = (rule, value, callback) => {
  console.log('validateIP')
  if (ctxData.gatewayForm.reportNetSw) {
    if (refExpIP.test(value)) {
      callback()
    } else {
      callback(new Error('请输入正确的IP地址！'))
    }
  } else {
    if (refExpIP.test(value) || refExpYM.test(value)) {
      callback()
    } else {
      callback(new Error('请输入正确的IP地址、域名地址！'))
    }
  }
}
const regCnt = /^[0-9]*[1-9][0-9]*$/
const validateRegAddr = (rule, value, callback) => {
  if (value !== 0 && !regCnt.test(value)) {
    callback(new Error('只能输入自然数！'))
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
  dataServiceName: '',
  gatewayTableData: [],
  dFlag: false,
  dTitle: '添加上报服务',
  networkNames: [],
  gatewayForm: {
    serviceName: '',
    ip: '',
    port: '',
    reportTime: 0,
    protocol: '',
    //EMQX.MQTT
    userName: '',
    password: '',
    clientID: '',
    keepAlive: '',
    cleanSession: false,
    deviceName: '',
    projectCode: '',
    token: '',
    slaveID: 0,
    coilStatusRegStart: 0,
    coilStatusRegCnt: 0,
    inputStatusRegStart: 0,
    inputStatusRegCnt: 0,
    holdingRegStart: 0,
    holdingRegCnt: 0,
    inputRegStart: 0,
    inputRegCnt: 0,
    //gwai add 2023-04-05
    appKey: '',
    productKey: '',
    deviceID: '',
    deviceSecret: '',
    //lp add 2023-06-17
    reportNetSw: false,
    reportNet: ''
  },
  paramName: {
    serviceName: '服务名称',
    ip: '服务名称',
    port: '服务名称',
    reportTime: '上报时间',
    protocol: '协议名称',
    //EMQX.MQTT
    UserName: '用户名',
    Password: '密码',
    ClientID: '客户端名称',
    keepAlive: '保活时间',
    cleanSession: '是否清除会话',
    deviceName: '通信编码',
    projectCode: '项目编码',
    token: '密钥',
    //gwai add 2023-04-05
    AppKey: 'appKey',
    ProductKey: '产品密钥',
    DeviceID: '通讯地址',
    DeviceSecret: '设备密钥',
  },
  gatewayRules: {
    serviceName: [
      {
        required: true,
        message: '服务名称不能为空',
        trigger: 'blur',
      },
    ],
    ip: [
      {
        required: true,
        validator: validateIP,
        trigger: 'blur',
      },
    ],
    port: [
      {
        required: true,
        message: '端口不能为空',
        trigger: 'blur',
      },
    ],
    reportTime: [
      {
        required: true,
        message: '上报时间不能为空',
        trigger: 'blur',
      },
      {
        type: 'number',
        message: '上报时间只能输入数字',
      },
    ],
    protocol: [
      {
        required: true,
        message: '协议名称不能为空',
        trigger: 'blur',
      },
    ],
    ModbusTCP: [
      {
        required: true,
        message: '上报地址不能为空',
        trigger: 'blur',
      },
    ],
    userName: [
      {
        required: true,
        message: '用户名不能为空',
        trigger: 'blur',
      },
    ],
    password: [
      {
        required: true,
        message: '密码不能为空',
        trigger: 'blur',
      },
    ],
    clientID: [
      {
        required: true,
        message: '客户端名称不能为空',
        trigger: 'blur',
      },
    ],
    // lp update 2023-06-17
    /*keepAlive: [
      {
        required: true,
        message: '保活时间不能为空',
        trigger: 'blur',
      },
    ],*/
    deviceName: [
      {
        required: true,
        message: '通讯编码不能为空',
        trigger: 'blur',
      },
    ],
    projectCode: [
      {
        required: true,
        message: '项目编码不能为空',
        trigger: 'blur',
      },
    ],
    token: [
      {
        required: true,
        message: '密钥不能为空',
        trigger: 'blur',
      },
    ],
    // lp update 2023-06-17
    /*cleanSession: [
      {
        required: true,
        message: '清除会话不能为空',
        trigger: 'blur',
      },
    ],*/
    slaveID: [
      {
        required: true,
        message: '从机ID不能为空',
        trigger: 'blur',
      },
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    coilStatusRegStart: [
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    coilStatusRegCnt: [
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    inputStatusRegStart: [
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    inputStatusRegCnt: [
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    holdingRegStart: [
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    holdingRegCnt: [
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    inputRegStart: [
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    inputRegCnt: [
      {
        validator: validateRegAddr,
        trigger: 'blur',
      },
    ],
    //gwai add 2023-04-05
    appKey: [
      {
        required: true,
        message: 'appKey不能为空',
        trigger: 'blur',
      },
    ],
    productKey: [
      {
        required: true,
        message: '产品密钥不能为空',
        trigger: 'blur',
      },
    ],
    deviceID: [
      {
        required: true,
        message: '产品密钥不能为空',
        trigger: 'blur',
      },
    ],
    deviceSecret: [
      {
        required: true,
        message: '产品密钥不能为空',
        trigger: 'blur',
      },
    ],
    // lp add 2023-06-17
    reportNet: [
      {
        required: true,
        message: '上报网卡不能为空',
        trigger: 'blur',
      },
    ],
  },
  protocolOptions: [
    {
      name: 'EMQX.MQTT',
      label: 'EMQX.MQTT',
    },
    {
      name: 'RT.STLV1',
      label: 'RT.STLV1',
    },
    {
      name: 'RT.STLV2',
      label: 'RT.STLV2',
    },
    {
      name: 'RT.MQTT',
      label: 'RT.MQTT',
    },
    {
      name: 'ModbusTCP',
      label: 'ModbusTCP',
    },
    //gwai add 2023-04-05
    {
      name: 'FSJY.MQTT',
      label: 'FSJY.MQTT',
    },
  ],
  dnList: ['RT.STLV1', 'RT.STLV2'],
  kaList: ['EMQX.MQTT', 'RT.MQTT', 'RT.STLV2','FSJY.MQTT'],
  dnFlag: 1,
  serviceName: '',
  handleFlag: '',//lp update 2023-06-12 首页上报信息查看详情，对点击返回按钮进行操作标识
})

//接收参数，处理命令详情页面的显示
const route = useRoute()
ctxData.serviceName = route.query.serviceName;

// 获取所有上报服务信息
const getGatewayList = (flag) => {
  const pData = {
    token: users.token,
    data: {},
  }
  ServiceApi.getGatewayList(pData).then(async (res) => {
    console.log('getGatewayList -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.gatewayTableData = res.data
      // 首页传参跳转到列表页面，根据参数自动进入上报节点页面中
      if (ctxData.serviceName && !ctxData.handleFlag) { //lp update 2023-06-12 返回按钮进行操作标识后，回到列表页面，不再在进行属性页面跳转
        const detail = ctxData.gatewayTableData.find(item => item.serviceName === ctxData.serviceName)
        showGateway(detail);
      }
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
getGatewayList()
//获取当前网关上报接口的通讯协议
const getReportProtocolList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  ServiceApi.getReportProtocolList(pData).then((res) => {
    if (!res) return
    console.log('getReportProtocolList -> res', res)
    if (res.code === '0') {
      ctxData.protocolOptions = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
getReportProtocolList()
// 过滤表格数据
const filterTableData = computed(() =>
  ctxData.gatewayTableData
    .filter(
      (item) =>
        !ctxData.dataServiceName || item.serviceName.toLowerCase().includes(ctxData.dataServiceName.toLowerCase())
    )
    .slice((ctxData.currentPage - 1) * ctxData.pagesize, ctxData.currentPage * ctxData.pagesize)
)
const filterTableDataPage = computed(() =>
  ctxData.gatewayTableData.filter(
    (item) => !ctxData.dataServiceName || item.serviceName.toLowerCase().includes(ctxData.dataServiceName.toLowerCase())
  )
)
// 刷新
const refresh = () => {
  getGatewayList(1)
}
//处理当前页变化
const handleCurrentChange = (val) => {
  ctxData.currentPage = val
}
//处理每页大小变化
const handleSizeChange = (val) => {
  ctxData.pagesize = val
}
// 添加上报服务
const addGateway = () => {
  ctxData.dFlag = true
  ctxData.dTitle = '添加上报服务'
}
// 编辑上报服务
const editGateway = (row) => {
  //
  ctxData.dFlag = true
  ctxData.dTitle = '编辑上报服务'
  ctxData.gatewayForm = {
    serviceName: row.serviceName,
    ip: row.ip,
    port: row.port,
    reportTime: row.reportTime,
    protocol: row.protocol,
  }
  // lp add 2023-06-17
  ctxData.gatewayForm['reportNetSw'] = row.reportNetSw
  ctxData.gatewayForm['reportNet'] = row.reportNet

  if (row.protocol.includes('FSJY.MQTT')) {  //gwai add 2023-04-05
    ctxData.gatewayForm['userName'] = row.param.UserName
    ctxData.gatewayForm['password'] = row.param.Password
    ctxData.gatewayForm['clientID'] = row.param.ClientID
    ctxData.gatewayForm['keepAlive'] = row.param.KeepAlive
    ctxData.gatewayForm['cleanSession'] = row.param.CleanSession
    ctxData.gatewayForm['appKey'] = row.param.AppKey
    ctxData.gatewayForm['productKey'] = row.param.ProductKey
    ctxData.gatewayForm['deviceID'] = row.param.DeviceID
    ctxData.gatewayForm['deviceSecret'] = row.param.DeviceSecret
  }else if (row.protocol.includes('MQTT')) {
    ctxData.gatewayForm['userName'] = row.param.userName
    ctxData.gatewayForm['password'] = row.param.password
    ctxData.gatewayForm['clientID'] = row.param.clientID
    ctxData.gatewayForm['keepAlive'] = row.param.keepAlive
    ctxData.gatewayForm['cleanSession'] = row.param.cleanSession
  } else if (row.protocol === 'RT.STLV1') {
    ctxData.gatewayForm['deviceName'] = row.param.deviceName
  } else if (row.protocol === 'RT.STLV2') {
    ctxData.gatewayForm['projectCode'] = row.param.projectCode
    ctxData.gatewayForm['deviceName'] = row.param.deviceName
    ctxData.gatewayForm['keepAlive'] = row.param.keepAlive
    ctxData.gatewayForm['token'] = row.param.token
  } else if (row.protocol === 'ModbusTCP') {
    ctxData.gatewayForm['ip'] = ''
    ctxData.gatewayForm['reportTime'] = 0

    ctxData.gatewayForm['slaveID'] = row.param.slaveID
    ctxData.gatewayForm['coilStatusRegStart'] = row.param.coilStatusRegStart
    ctxData.gatewayForm['coilStatusRegCnt'] = row.param.coilStatusRegCnt

    ctxData.gatewayForm['inputStatusRegStart'] = row.param.inputStatusRegStart
    ctxData.gatewayForm['inputStatusRegCnt'] = row.param.inputStatusRegCnt
    ctxData.gatewayForm['holdingRegStart'] = row.param.holdingRegStart
    ctxData.gatewayForm['holdingRegCnt'] = row.param.holdingRegCnt
    ctxData.gatewayForm['inputRegStart'] = row.param.inputRegStart
    ctxData.gatewayForm['inputRegCnt'] = row.param.inputRegCnt
  }
}
const showGateway = (row) => {
  ctxData.dnFlag = 2
  ctxData.curGateway = row
}

const showRegInfo = (row) => {
  ctxData.dnFlag = 3
  ctxData.curGateway = row
}
const changeDnFlag = (flag) => {
  //lp update 2023-06-12 首页上报信息查看详情，对点击返回按钮进行操作标识
  ctxData.handleFlag = flag
  ctxData.dnFlag = 1
  getGatewayList()
}

const changeProtocol = (protocol) => {
  if (protocol === 'RT.MQTT') {
    ctxData.gatewayForm.ip = 'mqtt.reatgreen.com'
    ctxData.gatewayForm.port = '1883'
    ctxData.gatewayForm.userName = 'reatgreen'
    ctxData.gatewayForm.password = 'Reat@2022'
    ctxData.gatewayForm.reportTime = 900
    ctxData.gatewayForm.keepAlive = '30'
    ctxData.gatewayForm.cleanSession = true
  }
}

const getGatewaySN = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  ProductApi.getGatewaySN(pData).then((res) => {
    if (res.code === '0') {
      ctxData.gatewayForm.clientID = res.data.sn
    } else {
      showOneResMsg(res)
    }
  })
}
// 初始化上报服务表单
const initGatewayForm = () => {
  ctxData.gatewayForm = {
    serviceName: '',
    ip: '',
    port: '',
    reportTime: 0,
    protocol: '',
    //EMQX.MQTT
    userName: '',
    password: '',
    clientID: '',
    keepAlive: '',
    cleanSession: false,
    deviceName: '',
    projectCode: '',
    token: '',
    appKey: '',
    productKey: '',
    deviceID: '',
    deviceSecret: '',

    // lp add 2023-06-17
    reportNetSw: false,
    reportNet: ''
  }
}
const gatewayFormRef = ref(null)
// 提交上报服务表单
const submitGatewayForm = () => {
  gatewayFormRef.value.validate((valid) => {
    if (valid) {
      const pForm = {
        serviceName: ctxData.gatewayForm.serviceName,
        ip: ctxData.gatewayForm.ip,
        port: ctxData.gatewayForm.port,
        reportTime: ctxData.gatewayForm.reportTime,
        protocol: ctxData.gatewayForm.protocol,
        reportNetSw: ctxData.gatewayForm.reportNetSw,
        reportNet: ctxData.gatewayForm.reportNet,
        param: {},
      }
      let param = {}
     if (ctxData.gatewayForm.protocol.includes('FSJY.MQTT')) {   //gwai add 2023-04-05
        param['userName'] = ctxData.gatewayForm.userName
        param['password'] = ctxData.gatewayForm.password
        param['clientID'] = ctxData.gatewayForm.clientID
        param['keepAlive'] = ctxData.gatewayForm.keepAlive
        param['cleanSession'] = ctxData.gatewayForm.cleanSession
        param['appKey'] = ctxData.gatewayForm.appKey
        param['productKey'] = ctxData.gatewayForm.productKey
        param['deviceID'] = ctxData.gatewayForm.deviceID
        param['deviceSecret'] = ctxData.gatewayForm.deviceSecret
      } else if (ctxData.gatewayForm.protocol.includes('MQTT')) {
        param['userName'] = ctxData.gatewayForm.userName
        param['password'] = ctxData.gatewayForm.password
        param['clientID'] = ctxData.gatewayForm.clientID
        param['keepAlive'] = ctxData.gatewayForm.keepAlive
        param['cleanSession'] = ctxData.gatewayForm.cleanSession
      } else if (ctxData.gatewayForm.protocol === 'RT.STLV1') {
        param['deviceName'] = ctxData.gatewayForm.deviceName
      } else if (ctxData.gatewayForm.protocol === 'RT.STLV2') {
        param['projectCode'] = ctxData.gatewayForm.projectCode
        param['deviceName'] = ctxData.gatewayForm.deviceName
        param['keepAlive'] = ctxData.gatewayForm.keepAlive
        param['token'] = ctxData.gatewayForm.token
      } else if (ctxData.gatewayForm.protocol === 'ModbusTCP') {
        param['slaveID'] = ctxData.gatewayForm.slaveID
        param['coilStatusRegStart'] = ctxData.gatewayForm.coilStatusRegStart
        param['coilStatusRegCnt'] = ctxData.gatewayForm.coilStatusRegCnt
        param['inputStatusRegStart'] = ctxData.gatewayForm.inputStatusRegStart
        param['inputStatusRegCnt'] = ctxData.gatewayForm.inputStatusRegCnt
        param['holdingRegStart'] = ctxData.gatewayForm.holdingRegStart
        param['holdingRegCnt'] = ctxData.gatewayForm.holdingRegCnt
        param['inputRegStart'] = ctxData.gatewayForm.inputRegStart
        param['inputRegCnt'] = ctxData.gatewayForm.inputRegCnt
      }
      pForm.param = param
      const pData = {
        token: users.token,
        data: pForm,
      }
      if (ctxData.dTitle.includes('添加')) {
        ServiceApi.addGateway(pData).then((res) => {
          handleResult(res, getGatewayList)
          cancelSubmit()
        })
      } else {
        ServiceApi.editGateway(pData).then((res) => {
          handleResult(res, getGatewayList)
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
  ctxData.dFlag = false
  gatewayFormRef.value.resetFields()
  initGatewayForm()
}
// 删除上报服务
const deleteGateway = (row) => {
  //
  ElMessageBox.confirm('确定要删除这个上报服务吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const pData = {
        token: users.token,
        data: {
          serviceName: row.serviceName,
        },
      }
      ServiceApi.deleteGateway(pData).then((res) => {
        handleResult(res, getGatewayList)
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
</style>
