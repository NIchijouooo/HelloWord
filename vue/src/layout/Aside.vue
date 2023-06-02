<template>
  <el-aside :width="ctxData.isCollapse ? variables.sideBarCloseWidth : variables.sideBarWidth">
    <div class="logo">
      <el-image :src="logoUrl"></el-image>
    </div>
    <el-menu
      ref="menu"
      :background-color="variables.menuBg"
      :text-color="variables.menuText"
      :active-text-color="variables.menuActiveText"
      router
      :default-active="activeIndex"
      :collapse="ctxData.isCollapse"
    >
      <template v-for="item in routers">
        <template v-if="!item['hidden']">
          <el-sub-menu v-if="item.children && item.children.length" :key="item.name" :index="concatPath(item.path)">
            <template #title>
              <div class="menu-icon" :class="itemClass(item.meta['icon'])"></div>
              <span>{{ item.meta.title }}</span>
            </template>
            <template v-for="sub in item.children">
              <el-sub-menu class="in-sub-menu" :key="sub.name" v-if="sub.children && sub.children.length" :index="concatPath(sub.path)">
                <template #title>
                  <span class="in-sub-menu-item"> {{ sub.meta.title }}</span>
                </template>
                <el-menu-item v-for="tub in sub.children" :key="tub.name" :index="concatPath(tub.path)">
                  <template #title>
                    <span style="margin-left: 10px">{{ tub.meta.title }}</span>
                  </template>
                </el-menu-item>
              </el-sub-menu>
              <el-menu-item v-else :index="concatPath(sub.path)" :key="sub.name">
                <template #title>{{ sub.meta.title }}</template>
              </el-menu-item>
            </template>
          </el-sub-menu>
          <el-menu-item v-else :index="concatPath(item.path)" :key="item.name">
            <div class="menu-icon" :class="itemClass(item.meta['icon'])"></div>
            <template #title>{{ item.meta.title }}</template>
          </el-menu-item>
        </template>
      </template>
    </el-menu>

    <div class="sidebar-toggle" @click="ctxData.isCollapse = !ctxData.isCollapse">
      <el-icon
        :size="18"
        color="#1890ff"
        :class="{ 'menu-collapse': true, isOpened: !ctxData.isCollapse }"
      >
        <expand />
      </el-icon>
      <span class="text" v-show="!ctxData.isCollapse">收缩侧边栏</span>
    </div>
  </el-aside>
</template>

<script setup>
import { computed, reactive } from 'vue'
import { useRoute } from 'vue-router'
import variables from 'styles/variables.module.scss'
import { Expand } from '@element-plus/icons-vue'
import { userStore } from '@/stores/user.js'
import setLoginInfo from 'utils/setLoginInfo.js'
import logoUrl from '@/assets/logoLight.png'//'@/assets/logo.png'

const user = userStore()
const ctxData = reactive({
  isCollapse: false,
})
const route = useRoute()
console.log('Aside route', route)
if (!user.userInfo) {
  setLoginInfo(user)
}
const routers = user.routers
const concatPath = (p_path) => {
  const path = `${p_path !== '' ? p_path : '/'}`
  return path
}
const itemClass = (name) => {
  return name
}
const activeIndex = computed(() => route.path)
</script>

<style lang="scss" scoped>
@mixin noScrollBar {
  overflow: hidden;
  overflow-y: scroll;
  -ms-overflow-style: none;
  overflow: -moz-scrollbars-none;
  scrollbar-width: none;
  &::-webkit-scrollbar {
    width: 0 !important;
  }
}

@mixin noSelect {
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -khtml-user-select: none;
  -moz-user-select: none;
  user-select: none;
}

@mixin line($n) {
  height: $n + px;
  line-height: $n + px;
}

.el-aside {
  box-sizing: border-box;
  height: 100vh;
  display: flex;
  display: -webkit-flex;
  flex-direction: column;
 /*background: linear-gradient(180deg, #384969, #22304c);
  box-shadow: 0 10px 10px -10px #35304c inset;*/
  box-shadow: 1px 0px 6px rgba(0, 21, 41, 0.35);
  overflow: hidden;
  transition: width 0.3s ease-in-out;
  -moz-transition: width 0.3s ease-in-out;
  -webkit-transition: width 0.3s ease-in-out;
  @include noSelect;

  .el-menu {
    flex: 1;
    border-right: none;
    @include noScrollBar;
    &:not(.el-menu--collapse) {
      width: 200px;
    }
  }
}
.logo {
  position: relative;
  display: flex;
  justify-content: space-around;
  align-items: center;
  width: 100%;
  height: 56px;
  box-sizing: border-box;
  border-bottom: solid 1px #e6e6e6;
}

.dashboard {
  background: url(assets/images/menu/dashboard-default.svg);
}
.is-active .dashboard {
  background: url(assets/images/menu/dashboard-active.svg);
}
.collectService {
  background: url(assets/images/menu/collectService-default.svg);
}
.is-active .collectService {
  background: url(assets/images/menu/collectService-active.svg);
}
.reportService {
  background: url(assets/images/menu/reportService-default.svg);
}
.is-active .reportService {
  background: url(assets/images/menu/reportService-active.svg);
}
.systemService {
  background: url(assets/images/menu/systemService-default.svg);
}
.is-active .systemService {
  background: url(assets/images/menu/systemService-active.svg);
}
.virtualService {
  background: url(assets/images/menu/virtualService-default.svg);
}
.is-active .virtualService {
  background: url(assets/images/menu/virtualService-active.svg);
}
.production {
  background: url(assets/images/menu/production-default.svg);
}
.is-active .production {
  background: url(assets/images/menu/production-active.svg);
}
.networkService {
  background: url(assets/images/menu/networkService-default.svg);
}
.is-active .networkService {
  background: url(assets/images/menu/networkService-active.svg);
}
.systemTool {
  background: url(assets/images/menu/systemTool-default.svg);
}
.is-active .systemTool {
  background: url(assets/images/menu/systemTool-active.svg);
}
.menu-icon {
  display: flex;
  width: 20px;
  height: 20px;
  margin-right: 20px;
  margin-left: 16px;
  background-size: 100% 100%;
  background-repeat: no-repeat;
}
:deep(.in-sub-menu .el-sub-menu__title) {
  height: 50px;
}

:deep(.el-sub-menu__title) {
  height: 68px;
  font-size: 16px;
  letter-spacing: 1px; //字间距
  background-color: inherit !important;
}
:deep(.el-menu-item) {
  height: 68px;
  font-size: 16px;
  letter-spacing: 1px; //字间距
}
:deep(.el-sub-menu .el-menu-item) {
  padding-left: 24px !important;
  font-size: 14px;
  height: 50px;
  justify-content: center;
}
:deep(.in-sub-menu .el-menu-item) {
  padding-left: 60px !important;
  justify-content: start;
}
:deep(.el-menu) {
  background: initial;
}
:deep(.el-menu-item:hover) {
  background-color: #f0f0f0;
}
:deep(.el-menu-item.is-active:hover) {
  background-color: #f0f0f0;
}
/*:deep(.is-active) {
  background-color: #f0f0f0;
}*/
:deep(.el-sub-menu .el-menu-item.is-active) {
  background-color: #f0f0f0;
}
:deep(.el-sub-menu .el-menu-item .menu-icon) {
  height: 18px;
  width: 18px;
  background-size: 100% 100%;
  background-repeat: no-repeat;
}
:deep(.el-tooltip__trigger) {
  display: flex !important;
  // justify-content: center;
  align-items: center;
  justify-content: center;
  .menu-icon {
    margin: 0;
  }
}
:deep(.el-sub-menu__title .el-sub-menu__title) {
  height: 50px;
}

:deep(.el-popper.is-light) {
  background: #f0f0f0;
}
:deep(.el-menu--popup-right-start .el-menu-item) {
  font-size: 14px;
}
.in-sub-menu-item {
  font-size: 14px;
  margin-left: 19px;
}
.el-menu--popup-right-start .in-sub-menu-item {
  margin-left: 0;
}
.menu-collapse {
  cursor: pointer;
  font-size: 24px;
  margin-right: 10px;
}
.isOpened {
  transform: rotate(180deg);
}
.sidebar-toggle {
  height: 56px;
  width: 100%;
  margin-bottom: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 14px;
  color: #1890ff;
  background: transparent;
  box-shadow: inset 0 -50px 30px rgba(123, 240, 255,  0.02);
  transition: all .15s;
  cursor: pointer;
  user-select: none;

  .text {
    margin: 0 0px;
    word-break: keep-all;
  }
}
</style>
