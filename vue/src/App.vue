<template>
  <el-config-provider :locale="zhCn">
    <div class="app">
      <router-view /></div
  ></el-config-provider>
</template>
<script setup>
import zhCn from 'element-plus/lib/locale/lang/zh-cn'
</script>
<style lang="scss">
#app {
  margin: 0;
  padding: 0;
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
}

// 页面主体结构
.main-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  .main {
    position: absolute;
    top: 20px;
    left: 20px;
    right: 20px;
    bottom: 20px;
    background-color: #fff;
    border-radius: 4px;
    overflow-x: hidden;
  }
  .title {
    position: absolute;
    top: 0;
    height: 80px;
    left: 24px;
    right: 24px;
    display: flex;
    justify-content: right;
    align-items: center;
    padding: 24px 0 20px;
    box-sizing: border-box;
    .tName {
      line-height: 16px;
      font-size: 18px;
      border-left: 4px solid #3054eb;
      padding-left: 20px;
    }
    .title-left {
      display: flex;
      align-items: center;
    }
  }
  .content {
    width: auto;
    height: auto;
    min-width: 700px;
    overflow-x: hidden;
  }
  .search-bar {
    width: 100%;
    margin-top: 15px;
  }
  .tool-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 10px;
    margin-bottom: 15px;
  }
  .btn-icon {
    display: flex;
    width: 14px;
    height: 14px;
    background-size: 100% 100%;
    background-repeat: no-repeat;
  }
  .pagination {
    position: absolute;
    bottom: 34px;
    left: 0;
    right: 0;
    display: flex;
    justify-content: center;
  }
  .dialog-pagination {
    bottom: 0px;
  }
  .dialog-content {
    position: relative;
    width: 100%;
    min-height: 400px;
    max-height: 550px;
    overflow: auto;
    padding: 0px 8px 0 0;
  }
}
// 自定义弹窗内容
.param-content {
  position: relative;
  width: 400px;
  padding-left: 40px;
  .pc-title {
    position: relative;
    display: flex;
    align-items: center;
    height: 50px;
    .pct-info {
      line-height: 16px;
      font-size: 18px;
      border-left: 4px solid #3054eb;
      padding-left: 20px;
    }
  }
  .pc-content {
    position: relative;
    padding: 0 24px;
    height: 100%;
    .param-item {
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
  }
}
//左右页
.mc-left {
  position: absolute;
  left: 20px;
  width: 384px;
  top: 20px;
  bottom: 20px;
  background-color: #fff;
  border-radius: 4px;
}
.mcl-content {
  position: absolute;
  top: 74px;
  left: 18px;
  right: 18px;
  bottom: 20px;
  padding: 6px;
  overflow: auto;
}
.mclc-card {
  margin-bottom: 10px;
}
.mclc-head {
  display: flex;
  justify-content: space-between;
  margin: 4px;
  border-bottom: 1px solid #e4e7ed;
}
.mclc-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  padding: 6px 10px;
}
.mc-right {
  position: absolute;
  left: 420px;
  right: 20px;
  top: 20px;
  bottom: 20px;
  background-color: #fff;
  border-radius: 4px;
}
</style>
