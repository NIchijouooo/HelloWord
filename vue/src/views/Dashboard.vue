<template>
  <div class="dashboard">
    <div class="head">
      <div class="head-item" @click="changeIndex(0)">
        <div class="head-item-content bg0">
          <div class="hic-left">
            <div class="hicl-title">CPU使用率</div>
            <div class="hicl-main">
              {{ ctxData.sysParams.cpuUse }}
              <span class="hicl-unit">%</span>
            </div>
          </div>
          <div class="hic-right">
            <div class="hic-icon bg-icon0"></div>
          </div>
        </div>
        <div v-show="ctxData.activeIndex == 0" class="head-item-active"></div>
      </div>
      <div class="head-item" @click="changeIndex(1)">
        <div class="head-item-content bg1">
          <div class="hic-left">
            <div class="hicl-title">内存使用率</div>
            <div class="hicl-main">
              {{ ctxData.sysParams.memUse }}
              <span class="hicl-unit">%</span>
            </div>
          </div>
          <div class="hic-right">
            <div class="hic-icon bg-icon4"></div>
          </div>
        </div>
        <div v-show="ctxData.activeIndex == 1" class="head-item-active"></div>
      </div>
      <div class="head-item" @click="changeIndex(2)">
        <div class="head-item-content bg2">
          <div class="hic-left">
            <div class="hicl-title">硬盘使用率</div>
            <div class="hicl-main">
              {{ ctxData.sysParams.diskUse }}
              <span class="hicl-unit">%</span>
            </div>
          </div>
          <div class="hic-right">
            <div class="hic-icon bg-icon3"></div>
          </div>
        </div>
        <div v-show="ctxData.activeIndex == 2" class="head-item-active"></div>
      </div>
      <div class="head-item" @click="changeIndex(3)">
        <div class="head-item-content bg3">
          <div class="hic-left">
            <div class="hicl-title">设备在线率</div>
            <div class="hicl-main">
              {{ ctxData.sysParams.deviceOnline }}
              <span class="hicl-unit">%</span>
            </div>
          </div>
          <div class="hic-right">
            <div class="hic-icon bg-icon1"></div>
          </div>
        </div>
        <div v-show="ctxData.activeIndex == 3" class="head-item-active"></div>
      </div>
      <div class="head-item" @click="changeIndex(4)">
        <div class="head-item-content bg4">
          <div class="hic-left">
            <div class="hicl-title">通讯丢包率</div>
            <div class="hicl-main">
              {{ ctxData.sysParams.devicePacketLoss }}
              <span class="hicl-unit">%</span>
            </div>
          </div>
          <div class="hic-right">
            <div class="hic-icon bg-icon2"></div>
          </div>
        </div>
        <div v-show="ctxData.activeIndex == 4" class="head-item-active"></div>
      </div>
    </div>
    <div class="content">
      <div class="basic-info">
        <div class="bi-title">基本信息</div>
        <div class="bi-body">
          <div class="bib-item">
            <div class="bibi-label">网关名称：</div>
            <div class="bibi-info">{{ ctxData.sysParams.name }}</div>
          </div>
          <div class="bib-item">
            <div class="bibi-label">网关编号：</div>
            <div class="bibi-info" style="display: inline-flex; align-items: center">
              {{ ctxData.sysParams.SN }}

              <el-button
                v-show="ctxData.showFlag"
                style="color: #3054eb; border: 1px solid #3054eb; margin-left: 10px; letter-spacing: 1px"
                color="#DFE6FE"
                @click="showBarcode()"
              >
                查看条形码
              </el-button>
            </div>
          </div>
          <div class="bib-item">
            <div class="bibi-label">硬件版本：</div>
            <div class="bibi-info">{{ ctxData.sysParams.hardVer }}</div>
          </div>
          <div class="bib-item">
            <div class="bibi-label">软件版本：</div>
            <div class="bibi-info">{{ ctxData.sysParams.softVer }}</div>
          </div>
          <div class="bib-item">
            <div class="bibi-label">系统时间：</div>
            <div class="bibi-info">{{ ctxData.sysParams.systemRTC }}</div>
          </div>
          <div class="bib-item">
            <div class="bibi-label">内存总量：</div>
            <div class="bibi-info">{{ ctxData.sysParams.memTotal }}</div>
          </div>
          <div class="bib-item">
            <div class="bibi-label">硬盘总量：</div>
            <div class="bibi-info">{{ ctxData.sysParams.diskTotal }}</div>
          </div>
          <div class="bib-item">
            <div class="bibi-label">运行时间：</div>
            <div class="bibi-info">{{ ctxData.sysParams.runTime }}</div>
          </div>
        </div>
      </div>
      <div class="chart-info">
        <div class="bi-title">{{ ctxData.headItemName['hin' + ctxData.activeIndex] }}</div>
        <div class="bi-body">
          <line-chart :chart-data="ctxData.curChartData" :key="ctxData.curChartData"></line-chart>
        </div>
      </div>
    </div>
  </div>
  <el-dialog v-model="ctxData.bFlag" title="网关SN条形码" width="30%">
    <div style="display: flex; width: 100%; justify-content: center; text-align: center">
      <vue3-barcode :value="ctxData.barcodeVale" :height="60" />
    </div>
  </el-dialog>
</template>
<script setup>
import LineChart from 'comps/LineChart.vue'
import DashboardApi from 'api/dashboard.js'
import { nextTick } from 'vue'
import { userStore } from 'stores/user'
import Vue3Barcode from 'vue3-barcode'
const users = userStore()
// 自定义响应数据
const ctxData = reactive({
  activeIndex: 0,
  headItemName: {
    hin0: 'CPU占用率',
    hin1: '内存使用率',
    hin2: '硬盘使用率',
    hin3: '设备在线率',
    hin4: '通讯丢包率',
  },
  curChartData: [],
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
})

const changeIndex = (indexValue) => {
  ctxData.activeIndex = indexValue
  const legend = ctxData.headItemName['hin' + indexValue]
  const pData = {
    token: users.token,
    data: {},
  }
  if (indexValue === 0) {
    DashboardApi.getCpuList(pData).then((res) => {
      handleDatas(res, legend)
    })
  } else if (indexValue === 1) {
    DashboardApi.getMemoryList(pData).then((res) => {
      handleDatas(res, legend)
    })
  } else if (indexValue === 2) {
    DashboardApi.getDiskList(pData).then((res) => {
      handleDatas(res, legend)
    })
  } else if (indexValue === 3) {
    DashboardApi.getDeviceOnlineList(pData).then((res) => {
      handleDatas(res, legend)
    })
  } else if (indexValue === 4) {
    DashboardApi.getDevicePacketLossList(pData).then((res) => {
      handleDatas(res, legend)
    })
  }
}
// 处理请求返回的数据
const handleDatas = (res, legend) => {
  if (res.code === '0') {
    const dataPoint = res.data
    var data = []
    var time = []
    for (var i = 0; i < dataPoint.length; i++) {
      const currentPoint = dataPoint[i]
      data.push(parseFloat(currentPoint.value))
      time.push(currentPoint.time)
    }
    nextTick(() => {})
    ctxData.curChartData = { data, time, legend }
  }
}
changeIndex(0)
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
</script>
<style lang="scss" scoped>
.dashboard {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  .head {
    position: absolute;
    top: 20px;
    height: 180px;
    left: 20px;
    right: 20px;
    margin: 0 -10px;
    display: flex;
    justify-content: flex-start;
    align-items: center;
    &-item {
      position: relative;
      height: 100%;
      width: 25%;
      padding: 0 10px;
      cursor: pointer;
      &-content {
        position: relative;
        width: 100%;
        height: 100%;
        border-radius: 8px;
        box-shadow: 0px 4px 6px 0px rgba(101, 75, 255, 0.2);
        display: flex;
        align-items: center;
        .hic-left {
          position: relative;
          height: 90px;
          width: calc(100% - 130px);
          .hicl-title {
            position: relative;
            height: 30px;
            line-height: 30px;
            margin-left: 40px;
            color: #ffffff;
            letter-spacing: 1px;
            font-size: 17px;
          }
          .hicl-main {
            position: relative;
            height: 60px;
            margin-left: 40px;
            line-height: 70px;
            color: #ffffff;
            font-size: 50px;
            letter-spacing: 1px;
            font-family: BebasNeue;
            margin-top: 5px;
          }
          .hicl-unit {
            font-size: 16px;
            font-family: auto;
            margin-left: 5px;
          }
        }
        .hic-right {
          position: relative;
          height: 120px;
          width: 130px;
          display: flex;
          align-items: center;
          justify-content: flex-start;
          .hic-icon {
            height: 90px;
            width: 90px;
          }
        }
      }
      .head-item-active {
        position: absolute;
        bottom: 0;
        height: 36px;
        background-color: #fff;
        left: 10px;
        right: 10px;
        border-radius: 0 0 4px 4px;
        opacity: 0.2;
      }
    }
    .bg0 {
      background: linear-gradient(49deg, #b87bee 0%, #7551fd 100%);
    }
    .bg1 {
      background: linear-gradient(49deg, #a87aff 0%, #7151ed 100%);
    }
    .bg2 {
      background: linear-gradient(49deg, #7a81ff 0%, #6345ff 100%);
    }
    .bg3 {
      background: linear-gradient(47deg, #5c8bff 0%, #0095ff 100%);
    }
    .bg4 {
      background: linear-gradient(52deg, #5c8bff 0%, #356df9 100%);
    }
    .bg-icon0 {
      background: url(assets/images/db-icon00.png) no-repeat;
      background-size: 100% 100%;
    }
    .bg-icon1 {
      background: url(assets/images/db-icon01.png) no-repeat;
      background-size: 100% 100%;
    }
    .bg-icon2 {
      background: url(assets/images/db-icon02.png) no-repeat;
      background-size: 100% 100%;
    }
    .bg-icon3 {
      background: url(assets/images/db-icon03.png) no-repeat;
      background-size: 100% 100%;
    }
    .bg-icon4 {
      background: url(assets/images/db-icon04.png) no-repeat;
      background-size: 100% 100%;
    }
  }
  .content {
    position: absolute;
    top: 220px;
    left: 20px;
    right: 20px;
    bottom: 20px;
    .basic-info {
      position: absolute;
      top: 0;
      left: 0;
      width: 580px;
      bottom: 0;
      background-color: #fff;
      border-radius: 4px;
    }
    .chart-info {
      position: absolute;
      top: 0;
      left: 600px;
      right: 0;
      bottom: 0;
      background-color: #fff;
      border-radius: 4px;
    }
    .bi-title {
      position: absolute;
      top: 40px;
      left: 32px;
      height: 16px;
      right: 32px;
      line-height: 16px;
      font-size: 18px;
      border-left: 4px solid #3054eb;
      padding-left: 20px;
    }
    .bi-body {
      position: absolute;
      top: 80px;
      left: 56px;
      bottom: 40px;
      right: 32px;
      overflow: auto;
      .bib-item {
        position: relative;
        height: 60px;
        width: 100%;
        display: flex;
        align-items: center;
        .bibi-label {
          width: 120px;
          color: #666;
          height: 60px;
          line-height: 60px;
        }
        .bibi-info {
          color: #000;
        }
      }
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
  .dashboard .content {
    top: 160px;
  }
  .dashboard .content .basic-info {
    width: 400px;
  }
  .dashboard .content .chart-info {
    left: 420px;
  }
  .dashboard .content .bi-title {
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
  .dashboard .content .bi-body {
    position: absolute;
    top: 60px;
    left: 36px;
    bottom: 20px;
    right: 20px;
    overflow: auto;
  }
  .dashboard .content .bi-body .bib-item {
    height: 36px;
    font-size: 14px;
  }
  .dashboard .content .bi-body .bib-item .bibi-label[data-v-22ba47ca] {
    height: 36px;
    line-height: 36px;
  }
}
</style>
