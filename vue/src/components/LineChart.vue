<template>
  <div class="lineChart">
    <v-chart :option="lineChartOption" :height="'400px'" />
  </div>
</template>
<script setup>
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import VChart from 'vue-echarts'
use([CanvasRenderer, LineChart, GridComponent, TooltipComponent, LegendComponent])

const props = defineProps({
  chartData: {
    type: Object,
    required: true,
  },
})
let lineChartOption = ref()
const setOptions = (chartData) => {
  const { data, time, legend } = chartData
  console.log('data', data)
  if (!data) return
  lineChartOption = ref({
    xAxis: {
      type: 'category',
      // 动态时间
      data: time,
      boundaryGap: false,
      axisTick: {
        show: false,
      },
      axisLabel: {
        rotate: 0,
        color: '#3054eb',
        show: true,
        formatter: (value) => {
          return value.split(' ').join('\n')
        },
      },
      axisLine: {
        lineStyle: {
          color: '#3054eb',
        },
      },
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#3054eb',
      },
      axisLine: {
        show: true,
        lineStyle: {
          color: '#3054eb',
        },
      },
    },
    grid: {
      left: 10,
      right: 34,
      bottom: 20,
      top: 30,
      containLabel: true,
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
      },
      padding: [5, 10],
    },
    legend: {
      // 动态图例
      data: [legend],
    },
    series: [
      {
        // 动态name
        name: legend,
        itemStyle: {
          color: '#FF005A',
        },
        lineStyle: {
          color: '#FF005A',
          width: 2,
        },
        smooth: true,
        type: 'line',
        data: data,
        animationDuration: 1000,
      },
    ],
  })
}
console.log('lineChartOption', lineChartOption)
if (props.chartDat !== null) {
  setOptions(props.chartData)
}
</script>
<style lang="scss" scoped>
.lineChart {
  width: 100%;
  height: 100%;
}
:deep(x-vue-echarts div) {
  width: 100% !important;
  height: 100% !important;
}
:deep(x-vue-echarts div canvas) {
  width: 100% !important;
  height: 400px !important;
}
</style>
