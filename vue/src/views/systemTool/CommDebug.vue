<template>
  <div class="main-container">
    <div class="main" style="background-color: inherit">
      <el-row :gutter="20" style="height: 100%">
        <el-col :span="24">
          <el-card class="box-card" :shadow="'hover'">
            <template #header>
              <div class="card-header">
                <span>通讯调试助手</span>
              </div>
            </template>
            <div class="item">
              <div class="st-receive">
                <div class="str-content">
                  <div class="str-noData" v-if="ctxData.receiveData.length == 0">显示区</div>
                  <div v-for="(item, index) in ctxData.receiveData" :key="index">
                    <div v-if="ctxData.receiveTime">【{{ item.date }}】</div>
                    <div :class="item.type == 1 ? 'TxInfo' : 'RxInfo'">
                      {{ item.type == 1 ? 'Tx：' + item.data : 'Rx：' + item.data }}
                    </div>
                  </div>
                </div>
              </div>
              <div class="st-send">
                <el-input
                  v-model="ctxData.sendData"
                  class="sts-sendArea"
                  rows="5"
                  type="textarea"
                  :placeholder="ctxData.sendData == '' ? '发送区' : ''"
                />
              </div>
              <div class="st-operate">
                <div>
                  <el-select v-model="ctxData.curCollect" style="width:100%;margin-bottom: 10px" placeholder="请选择采集接口">
                    <el-option
                      v-for="item in ctxData.interfaceList"
                      :key="item.collInterfaceName"
                      :label="item.collInterfaceName"
                      :value="item.collInterfaceName"
                    />
                  </el-select>
                  <div style="margin-bottom: 10px">
                    <el-checkbox
                      v-model="ctxData.receiveTime"
                      label="接收时间"
                      border
                      style="width: 100%; height: 34px"
                    />
                  </div>
                  <el-button type="danger" style="height: 34px; width: 100%" plain @click="clearReceive()"
                    >清空显示区</el-button
                  >
                </div>
                <div>
                  <el-button
                    type="primary"
                    style="height: 34px; width: 100%; margin-bottom: 10px"
                    plain
                    @click="sendMessage()"
                  >
                    发送
                  </el-button>
                  <el-select v-model="ctxData.curCheck" style="width:100%;margin-bottom: 10px" placeholder="请选择数据校验">
                    <el-option
                      v-for="item in ctxData.checkOptions"
                      :key="item.name"
                      :label="item.label"
                      :value="item.name"
                    />
                  </el-select>
                  <el-button type="danger" style="height: 34px; width: 100%" plain @click="clearSend()"
                    >清空发送区</el-button
                  >
                </div>
              </div>
              <div class="st-remark">
                <el-row :gutter="16">
                  <el-col :span="1">操作步骤:</el-col>
                  <el-col :span="21">1、选择采集接口</el-col>
                </el-row>
                <el-row :gutter="16">
                  <el-col :offset="1" :span="21">2、选择校验方式</el-col>
                </el-row>
                <el-row :gutter="16">
                  <el-col :offset="1" :span="21">3、发送区编写发送的报文</el-col>
                </el-row>
                <el-row :gutter="16">
                  <el-col :offset="1" :span="21">4、显示区显示发送的报文和返回的报文</el-col>
                </el-row>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>
<script setup>
import { userStore } from 'stores/user'
import SysToolApi from 'api/sysTool.js'
import InterfaceApi from 'api/interface.js'

const users = userStore()
const ctxData = reactive({
  receiveData: [],
  sendData: '',
  curCollect: '',
  interfaceList: [
    { name: 'collect01', label: '采集接口01' },
    { name: 'modbusTCP', label: '采集接口02' },
  ],
  receiveTime: true,
  curCheck: 0,
  checkOptions: [
    { name: 0, label: '不校验' },
    { name: 1, label: 'CRC16' },
    { name: 2, label: 'SUM' },
  ],
})
// 获取采集接口列表
const getInterfaceList = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  InterfaceApi.getInterfaceList(pData).then((res) => {
    console.log('getInterfaceList -> res', res)
    if (!res) return
    if (res.code === '0') {
      ctxData.interfaceList = res.data
    } else {
      showOneResMsg(res)
    }
  })
}
getInterfaceList()
//发送
const sendMessage = () => {
  if (ctxData.curCollect === '') {
    showOneResMsg({ message: '请选择采集接口！' })
    return
  }
  if (ctxData.curCheck === '') {
    showOneResMsg({ message: '请选择校验方式！' })
    return
  }
  let sendMessage = ''
  const regMessage = /^[0-9a-fA-F]+$/
  if (ctxData.sendData === '') {
    showOneResMsg({ message: '发送数据不能为空！' })
    return
  }

  let temp = ctxData.sendData.replaceAll(' ', '')
  if (!regMessage.test(temp)) {
    showOneResMsg({ message: '报文格式错误，只能输入字符【0-9,a-f,A-F或者空格】！' })
    ctxData.sendData = ''
    temp = ''
    return
  }
  const len = temp.length
  for (var i = 0; i < len; i = i + 2) {
    if (i < len) {
      sendMessage += temp.substr(i, 2) + (i + 2 >= len ? '' : ' ')
    }
  }

  //调用接口
  const pData = {
    token: users.token,
    data: {
      collInterfaceName: ctxData.curCollect,
      directData: sendMessage,
      checkSum: ctxData.curCheck,
    },
  }
  SysToolApi.sendMessage(pData).then((res) => {
    if (res.code === '0') {
      res.data.forEach((item) => {
        ctxData.receiveData.push(item)
      })
    } else {
      showOneResMsg(res)
    }
  })
}
//清空显示区
const clearReceive = () => {
  ctxData.receiveData = []
}
//清空发送区
const clearSend = () => {
  ctxData.sendData = ''
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
@use 'styles/custom-scoped.scss' as *;
.box-card {
  padding-bottom: 20px;
  padding-right: 20px;
}
.item {
  position: relative;
  height: 780px;
  width: 100%;
  min-width: 540px;
}
.st-receive {
  position: absolute;
  left: 0;
  right: 192px;
  top: 0;
  width: calc(100% - 290px);
  bottom: 241px;
  box-sizing: border-box;
  overflow-y: auto;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 12px;
}
.st-send {
  position: absolute;
  left: 0;
  right: 192px;
  height: 129px;
  bottom: 100px;
  width: calc(100% - 290px);
  box-sizing: border-box;
  overflow-y: auto;
  border-radius: 4px;
}
:deep(.sts-sendArea .el-textarea__inner) {
  padding: 12px;
}
.st-receive:hover {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

.str-content {
  position: relative;
  width: 100%;
  font-size: 14px;
}
.str-noData {
  color: #a8abb2;
}
.st-operate {
  position: absolute;
  right: 0;
  width: 280px;
  top: 0;
  bottom: 100px;
  box-sizing: border-box;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 12px;
  display: flex;
  justify-content: space-between;
  flex-flow: column;
}

.st-operate:hover {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}
.st-remark {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 84px;
  font-size: 12px;
  color: #f56c6c;
}
.RxInfo {
  color: #67c23a;
}
.TxInfo {
  color: #409eff;
}
:deep(.el-col) {
  position: relative;
}
:deep(.el-checkbox) {
  justify-content: center;
}
:deep(.el-card__body) {
  height: 100%;
  position: relative;
  box-sizing: border-box;
  overflow: auto;
}
</style>
