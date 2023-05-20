import Icon from 'comps/icon/Index.vue'
import * as elIcons from '@element-plus/icons-vue'
export function registerIcons(app) {
  /*
   * 全局注册 Icon
   * 使用方式: <Icon name="name" size="size" color="color" />
   * 详见<待完善>
   */
  app.component('Icon', Icon)

  /*
   * 全局注册element Plus的icon
   */
  const icons = elIcons
  for (const i in icons) {
    app.component(`el-icon-${icons[i].name}`, icons[i])
  }
}
/**
 * 是否是外部链接
 * @param {string} path
 * @return {Boolean}
 */
export function isExternal(path) {
  return /^(https?|ftp|mailto|tel):/.test(path)
}
/* 加载网络css文件 */
export function loadCss(url) {
  const link = document.createElement('link')
  link.rel = 'stylesheet'
  link.href = url
  link.crossOrigin = 'anonymous'
  document.getElementsByTagName('head')[0].appendChild(link)
}
