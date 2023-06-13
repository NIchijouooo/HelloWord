import Layout from '../layout/Index.vue'
import Login from '../views/Login.vue'
const layoutMap = [
  {
    path: '/dashboard',
    name: 'Dashboard',
    META: { title: '首页', icon: 'role' },
    component: () => import('../views/Dashboard.vue'),
    children: [],
  },
  // 采集服务 start
  {
    path: '/commInterface',
    name: 'CommInterface',
    meta: { title: '通讯接口', icon: 'commInterface' },
    component: () => import('../views/collectService/CommInterface.vue'),
    children: [],
  },
  {
    path: '/interface',
    name: 'Interface',
    meta: { title: '采集接口', icon: 'interface' },
    component: () => import('../views/collectService/Interface.vue'),
    children: [],
  },
  {
    path: '/deviceManage',
    name: 'DeviceManage',
    meta: { title: '设备管理', icon: 'deviceManage' },
    component: () => import('../views/collectService/DeviceManage.vue'),
    children: [],
  },
  {
    path: '/deviceModelLua',
    name: 'DeviceModelLua',
    meta: { title: '通用', icon: 'deviceModelLua' },
    component: () => import('../views/collectService/general/DeviceModelLua.vue'),
    children: [],
  },
  {
    path: '/deviceModelS7',
    name: 'DeviceModelS7',
    meta: { title: 'Snap7', icon: 'deviceModelS7' },
    component: () => import('../views/collectService/snap7/DeviceModelS7.vue'),
    children: [],
  },
  {
    path: '/deviceModelRtu',
    name: 'DeviceModelRtu',
    meta: { title: 'Modbus', icon: 'deviceModelRtu' },
    component: () => import('../views/collectService/modbusRtu/DeviceModelRtu.vue'),
    children: [],
  },
  {
    path: '/deviceModelD07',
    name: 'deviceModelD07',
    meta: { title: 'DL/T645-2007', icon: 'deviceModelD07' },
    component: () => import('../views/collectService/DLT64507/DeviceModelD07.vue'),
    children: [],
  },
  // 采集服务 end
  // 上报服务 start
  {
    path: '/transferModel',
    name: 'TransferModel',
    meta: { title: '上报模型', icon: 'transferModel' },
    component: () => import('../views/reportService/TransferModel.vue'),
    children: [],
  },
  {
    path: '/dataService',
    name: 'DataService',
    meta: { title: '上报服务', icon: 'dataService' },
    component: () => import('../views/reportService/DataService.vue'),
    children: [],
  },
  // 上报服务 end
  // 网络服务 start
  {
    path: '/network',
    name: 'Network',
    meta: { title: '网口设置', icon: 'network' },
    component: () => import('../views/networkService/Network.vue'),
    children: [],
  },
  {
    path: '/mobile',
    name: 'Mobile',
    meta: { title: '移动网络', icon: 'mobile' },
    component: () => import('../views/networkService/Mobile.vue'),
    children: [],
  },
  // 网络服务 end
  // 虚拟设备 start
  {
    path: '/virtualDevice',
    name: 'VirtualDevice',
    meta: { title: '虚拟设备', icon: 'virtualDevice' },
    component: () => import('../views/virtual/VirtualDevice.vue'),
    children: [],
  },
  // 虚拟设备 end
  // 系统服务 start
  {
    path: '/backupRestore',
    name: 'BackupRestore',
    meta: { title: '备份还原', icon: 'backupRestore' },
    component: () => import('../views/systemService/BackupRestore.vue'),
    children: [],
  },
  {
    path: '/sysUpgrade',
    name: 'SysUpgrade',
    meta: { title: '系统升级', icon: 'sysUpgrade' },
    component: () => import('../views/systemService/SysUpgrade.vue'),
    children: [],
  },
  {
    path: '/sysLog',
    name: 'SysLog',
    meta: { title: '系统日志', icon: 'sysLog' },
    component: () => import('../views/systemService/SysLog.vue'),
    children: [],
  },
  {
    path: '/dictType',
    name: 'DictType',
    meta: { title: '字典管理', icon: 'sysDict' },
    component: () => import('../views/systemService/dict/index.vue'),
    children: [],
  },
  {
    path: '/dictData',
    name: 'dictData',
    meta: { title: '字典详情', icon: 'sysDictData' },
    component: () => import('../views/systemService/dict/data.vue'),
    children: [],
  },
  // 系统服务 end
  // 系统工具 start
  {
    path: '/commDebug',
    name: 'CommDebug',
    meta: { title: '通讯调试', icon: 'commDebug' },
    component: () => import('../views/systemTool/CommDebug.vue'),
    children: [],
  },
  {
    path: '/networkTest',
    name: 'NetworkTest',
    meta: { title: '网络诊断', icon: 'networkTest' },
    component: () => import('../views/systemTool/NetworkTest.vue'),
    children: [],
  },
  {
    path: '/ntpTiming',
    name: 'NTPTiming',
    meta: { title: '通讯调试', icon: 'ntpTiming' },
    component: () => import('../views/systemTool/NTPTiming.vue'),
    children: [],
  },
  // 系统工具 end
  // 生产服务 start
  {
    path: '/production',
    name: 'Production',
    meta: { title: '生产服务', icon: 'production' },
    component: () => import('../views/Production.vue'),
    children: [],
  },
  // 生产服务end
  {
    path: '/error',
    name: 'NotFound',
    hidden: true,
    meta: { title: '404' },
    component: () => import('../components/NotFound.vue'),
  },
  {
    path: '/:w+',
    hidden: true,
    redirect: { name: 'NotFound' },
  },
]

const routes = [
  {
    path: '/',
    redirect: { name: 'Login' },
  },
  {
    path: '/login',
    name: 'Login',
    meta: { title: '登录' },
    component: Login,
  },
  {
    path: '/layout',
    name: 'Layout',
    component: Layout,
    redirect: { name: 'Test1' },
    children: [...layoutMap],
  },
]

export { routes, layoutMap }
