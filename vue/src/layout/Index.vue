<template>
  <el-container class="layout">
    <Aside />
    <el-container>
      <Header />
      <div class="app_wrap">
        <el-main>
          <router-view />
        </el-main>
      </div>
    </el-container>
  </el-container>
</template>

<script setup>
import Aside from './Aside.vue'
import Header from './Header.vue'
</script>

<style lang="scss" scoped>
.layout {
  box-sizing: border-box;
  width: 100vw;
  height: 100vh;
  overflow: auto;
  background-color: #f3f3f3;
  .el-container {
    height: 100vh;
    min-width: 1370px;
    display: flex;
    display: -webkit-flex;
    flex-direction: column;

    .app_wrap {
      flex: 1;
      height: auto;
      background-color: #f5f7fa;
      overflow: hidden;
      overflow-y: auto;
      .el-main {
        position: relative;
        height: 100%;
        width: 100%;
      }

      .el-footer {
        line-height: 48px;
        font-size: 14px;
        color: rgba(0, 0, 0, 0.5);
        text-align: center;
        overflow: hidden;
        user-select: none;
      }
    }
  }
}
</style>
