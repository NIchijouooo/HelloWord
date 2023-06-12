<template>
  <div class="dashboard">
    <el-container style="height: calc(100% - 10px);">
      <el-main style="height: calc(50% - 10px);">
        <DeviceList style="width: 100%; height: calc(63.7% - 1px);"></DeviceList>
      </el-main>
      <el-footer style="height: calc(36% - 25px);padding: 0px;">
        <el-row :gutter="20" style="height: calc(100%);padding:0px;margin: 0px;">
          <el-col :span="8" style="height:calc(100%);padding-left:0px;padding-right: 15px;">
            <div class="basic-info">
              <div class="bi-title">基本信息</div>
              <div class="bi-body">
                <div class="cui-content">
                  <el-card class="box-card" shadow="hover">
                    <template #header>
                      <div class="card-header">
                        <div class="cui-header">
                          <span style="font-weight: 600">网关名称：{{ ctxData.sysParams.name }}</span>
                        </div>
                        <div>
                          <span>{{ ctxData.sysParams.SN }}</span>
                          <el-button v-show="ctxData.showFlag" class="head-tag" text type="primary" @click="showBarcode()" size="small" style="margin-left: 8px">查看条形码</el-button>
                        </div>
                      </div>
                    </template>
                    <div class="card-content">
                      <div class="cc-body">
                        <div class="ccb-item" style="height: 31px;line-height: 31px;">
                          <div class="">硬件版本：</div>
                          <div class="head-name">{{ ctxData.sysParams.hardVer }}</div>
                        </div>
                        <div class="ccb-item" style="height: 31px;line-height: 31px;">
                          <div class="">软件版本：</div>
                          <div class="head-name">{{ ctxData.sysParams.softVer }}</div>
                        </div>
                        <div class="ccb-item" style="height: 31px;line-height: 31px;">
                          <div class="">系统时间：</div>
                          <div class="head-name">{{ ctxData.sysParams.systemRTC }}</div>
                        </div>
                        <div class="ccb-item" style="height: 31px;line-height: 31px;">
                          <div class="">内存总量：</div>
                          <div class="head-name">{{ ctxData.sysParams.memTotal }}</div>
                        </div>
                        <div class="ccb-item" style="height: 31px;line-height: 31px;">
                          <div class="">硬盘总量：</div>
                          <div class="head-name">{{ ctxData.sysParams.diskTotal }}</div>
                        </div>
                        <div class="ccb-item" style="margin-bottom: 15px;height: 31px;line-height: 31px;">
                          <div class="">运行时间：</div>
                          <div class="head-name">{{ ctxData.sysParams.runTime }}</div>
                        </div>
                      </div>
                    </div>
                  </el-card>
                </div>
              </div>
            </div>
          </el-col>
          <el-col :span="16" style="height: calc(104%);background-color: #fff;border-radius: 4px;">
            <div style="padding-top: 10px;">
              <div class="bi-title">上报信息</div>
              <div class="bi-body">
                <div class="info-list">
                  <div class="info-item" v-for="(item, index) in ctxData.gatewayTableData" :key="index">
                    <div class="cui-content">
                    <el-card class="box-card" shadow="hover">
                      <template #header>
                        <div class="card-header">
                          <div class="cui-header">
                            <span style="font-weight: 600">{{ item.serviceName }}</span>
                            <el-tag size="small" style="margin-left: 8px" :type="item.reportStatus === 'onLine' ? 'success' : 'danger'">{{
                              item.reportStatus === 'onLine' ? '在线' : '离线'
                            }}</el-tag>
                          </div>
                          <div>
                            <el-button class="head-tag" text type="primary" @click="showGatewayInfo(item)" size="small">查看详情</el-button>
                          </div>
                        </div>
                      </template>
                      <div class="card-content">
                        <div class="cc-body">
                          <div class="ccb-item">
                            <div class="">协议名称：</div>
                            <div class="head-name">{{ item.protocol }}</div>
                          </div>
                          <div class="ccb-item">
                            <div class="">上报周期：</div>
                            <div class="head-name">{{ item.reportTime }} 秒</div>
                          </div>
                          <div class="ccb-item">
                            <div class="">产品密钥：</div>
                            <div class="head-name">{{ item.param.ProductKey }}</div>
                          </div>
                          <div class="ccb-item">
                            <div class="">通讯地址：</div>
                            <div class="head-name">{{ item.param.DeviceID }}</div>
                          </div>
                          <div class="ccb-item">
                            <div class="">设备密钥：</div>
                            <div class="head-name">{{ item.param.DeviceSecret }}</div>
                          </div>
                        </div>
                      </div>
                    </el-card>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </el-footer>
    </el-container>

    <el-dialog v-model="ctxData.bFlag" title="网关SN条形码" width="30%">
      <div style="display: flex; width: 100%; justify-content: center; text-align: center">
        <vue3-barcode :value="ctxData.barcodeVale" :height="60" />
      </div>
    </el-dialog>
  </div>
</template>
<script setup>
import DashboardApi from 'api/dashboard.js'
import ServiceApi from 'api/service.js'
import { nextTick } from 'vue'
import { userStore } from 'stores/user'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import Vue3Barcode from 'vue3-barcode'
import DeviceList from './collectService/DeviceList.vue'

const users = userStore()
// 自定义响应数据
const ctxData = reactive({
  sysParams: {
    SN: '',
    deviceOnline: '',
    devicePacketLoss: '',
    diskTotal: '',
    diskUse: '',
    hardVer: '',
    lockStatus: 0,
    memTotal: '',
    memUse: '',
    name: '',
    runTime: '',
    softVer: '',
    systemRTC: '',
  },
  showFlag: false,
  bFlag: false,
  barcodeVale: 0,

  gatewayTableData: []
})

// 获取系统参数
const getSysParams = () => {
  const pData = {
    token: users.token,
    data: {},
  }
  DashboardApi.getSysParams(pData).then((res) => {
    console.log('getSysParams -> res', res)
    if (res.code === '0') {
      ctxData.sysParams = res.data
      console.log('getSysParams -> ctxData.sysParams', ctxData.sysParams)
      ctxData.showFlag = true
      ctxData.barcodeVale = ctxData.sysParams.SN
      sessionStorage.setItem('SN', ctxData.barcodeVale)
    }
  })
}
getSysParams()
const showBarcode = () => {
  ctxData.bFlag = true
}
const router = useRouter()
const showGatewayInfo = (data) => {
  router.push({ path: 'dataService', query: { serviceName: data.serviceName }})
}
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
      if (flag === 1) {
        ElMessage({
          type: 'success',
          message: '刷新成功！',
        })
      }
    } else {
      showOneResMsg(res)
    }
  })
}
getGatewayList()
</script>
<style lang="scss" scoped>
.dashboard {
  width: 100%;
  height: 100%;
    .basic-info {
      width: 96%;
      height: calc(90%);
      padding-top: 10px;
      padding-left: 20px;
      padding-bottom: 30px;
      background-color: #fff;
      border-radius: 4px;
    }
    .bi-title {
      line-height: 16px;
      font-size: 18px;
      border-left: 4px solid #3054eb;
      padding-left: 20px;
    }
        
    .cui-content {
      position: relative;
      height: 100%;
      margin: 8px 8px 16px 8px;
    }
        
    .cui-header {
      display: flex;
      align-items: center;
      width: 190px;
      white-space: nowrap;
      overflow: auto;
    }
    .card-content {
      position: relative;
      font-size: 14px;
      .cc-body {
        position: relative;
        top: 10px;
        left: 18px;
        right: 18px;
        bottom: 8px;
        overflow: auto;
      }
      .ccb-item {
        float: left;
        height: 33px;
        width: calc(100% - 44px);
        padding: 0 4px;
        line-height: 33px;
        display: flex;
        justify-content: space-between;
      }
    }
    :deep(.card-header) {
      display: flex;
      justify-content: space-between;
    }
    :deep(.el-card) {
      height: 100%;
    }
    :deep(.el-card__header) {
      padding: 10px 20px;
    }
    :deep(.el-card__body) {
      height: calc(100% - 48px);
      position: relative;
      box-sizing: border-box;
      padding: 0;
    }
    :deep(.el-collapse-item__content) {
      padding-bottom: 15px;
    }

    .bi-body {
      overflow: auto;
      .bib-item {
        position: relative;
        height: 30px;
        width: 100%;
        display: flex;
        align-items: center;
        .bibi-label {
          width: 100px;
          color: #666;
        }
        .bibi-info {
          color: #000;
        }
      }

      .info-list {
        display: flex;
        flex-wrap: wrap;
        gap: 20px;
        padding: 10px;
      }

      .info-item {
        width: calc(50% - 10px);
        margin-bottom: 10px;
      }
    }
  
}
@media screen and (max-width: 1634px) {
  .dashboard .head {
    height: 120px;
  }
  .dashboard .head-item-content .hic-left .hicl-title {
    position: relative;
    height: 20px;
    line-height: 20px;
    margin-left: 20px;
    font-size: 14px;
  }
  .dashboard .head-item .head-item-active {
    height: 22px;
  }
  .dashboard .head-item-content .hic-left .hicl-main {
    height: 50px;
    margin-left: 20px;
    line-height: 60px;
    font-size: 40px;
    margin-top: 5px;
  }
  .dashboard .head-item-content .hic-left {
    width: calc(100% - 40px);
  }
  .dashboard .head-item-content .hic-right .hic-icon {
    height: 60px;
    width: 60px;
  }
  .dashboard .basic-info {
    width: 400px;
  }
  .dashboard .bi-title {
    position: absolute;
    top: 24px;
    left: 20px;
    height: 16px;
    right: 20px;
    line-height: 16px;
    font-size: 18px;
    border-left: 4px solid #3054eb;
    padding-left: 20px;
  }
  .dashboard .bi-body {
    position: absolute;
    top: 60px;
    left: 36px;
    bottom: 20px;
    right: 20px;
    overflow: hidden;
  }
  .dashboard .bi-body .bib-item {
    height: 36px;
    font-size: 14px;
  }
  .dashboard .bi-body .bib-item .bibi-label[data-v-22ba47ca] {
    height: 36px;
    line-height: 36px;
  }
}
</style>
