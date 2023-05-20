export const configStore = defineStore({
  id: 'config',
  state: () => ({
    // 配置信息
    configInfo: {
      name: '沃瑞珂边缘网关',
      version: 'V2.0.1',
    },
  }),
  actions: {
    setConfigInfo(configInfo) {
      this.configInfo = configInfo
    },
  },
})
