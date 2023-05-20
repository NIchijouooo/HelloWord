// 单纯的使⽤⾃⼰的返回数据话，可以不⽤引⼊mockjs
// test
export default [
  {
    url: '/api/test01',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: 'ok',
        data: [1, 2, 3, 4],
      }
    },
  },
  {
    url: '/tsls',
    method: 'get',
    response: () => {
      return {
        code: '0',
        message: '获取物模型列表成功！',
        data: [
          {
            name: 'HL7031',
            label: '海林风机盘管控制器1',
            type: 0,
            blockNum: 2,
            param: {
              name: 'MaiSiHVAC',
              label: '迈斯VRV网关',
              version: 'V0.0.5',
              author: 'leiyang',
              date: '2022/7/6',
              message: '',
            },
          },
          {
            name: 'HL7032',
            label: '海林风机盘管控制器2',
            param: {
              name: '',
            },
            blockNum: 2,
            type: 1,
          },
        ],
      }
    },
  },
  {
    url: '/tsl',
    method: 'post',
    response: () => {
      return {
        code: '0',
        message: '添加物模型成功！',
        data: {},
      }
    },
  },
  {
    url: '/tsl',
    method: 'put',
    response: () => {
      return {
        code: '0',
        message: '修改物模型成功！',
        data: {},
      }
    },
  },
  {
    url: '/tsl',
    method: 'delete',
    response: () => {
      return {
        code: '0',
        message: '删除物模型成功！',
        data: {},
      }
    },
  },
  //上报模型
  {
    url: '/service/report/models',
    method: 'get',
    response: () => {
      return {
        code: '0',
        message: '获取上报模型信息成功！',
        data: [
          {
            name: 'model01',
            label: '上报模型01',
          },
          {
            name: 'model02',
            label: '上报模型02',
          },
        ],
      }
    },
  },
  {
    url: '/service/report/model',
    method: 'post',
    response: () => {
      return {
        code: '0',
        message: '添加上报模型成功！',
        data: {},
      }
    },
  },
  {
    url: '/service/report/model',
    method: 'put',
    response: () => {
      return {
        code: '0',
        message: '修改上报模型成功！',
        data: {},
      }
    },
  },
  {
    url: '/service/report/model',
    method: 'delete',
    response: () => {
      return {
        code: '0',
        message: '删除上报模型成功！',
        data: {},
      }
    },
  },
  {
    url: '/service/report/model/properties',
    method: 'get',
    response: () => {
      return {
        code: '0',
        message: '获取成功！',
        data: [
          {
            uploadName: 0,
            isChecked: true,
            label: '2',
            name: '2',
            params: {
              dataLength: '',
              dataLengthAlarm: false,
              decimals: '',
              max: '',
              min: '',
              minMaxAlarm: false,
              step: '',
              stepAlarm: false,
              unit: '',
            },
            type: 0,
            value: null,
          },
        ],
      }
    },
  },
  {
    url: '/service/report/model/property',
    method: 'post',
    response: () => {
      return {
        code: '0',
        message: '添加属性成功！',
        data: {},
      }
    },
  },
  {
    url: '/service/report/model/property',
    method: 'put',
    response: () => {
      return {
        code: '0',
        message: '修改属性成功！',
        data: {},
      }
    },
  },
  {
    url: '/service/report/model/property',
    method: 'delete',
    response: () => {
      return {
        code: '0',
        message: '删除属性成功！',
        data: {},
      }
    },
  },
  {
    url: '/tsl/content/properties',
    method: 'get',
    response: () => {
      return {
        code: '0',
        message: '',
        data: [
          {
            name: 'nUa',
            label: 'A相电压',
            accessMode: 0,
            type: 1,
            params: {
              min: 0,
              max: 300,
              minMaxAlarm: true,
              step: 100,
              stepAlarm: true,
              decimals: 1,
              unit: 'V',
            },
          },
        ],
      }
    },
  },
  {
    url: '/network/mobiles',
    method: 'get',
    response: () => {
      return {
        code: '0',
        message: '',
        data: [
          {
            name: 'EC600N',
            model: 'EC600N',
            runParam: {
              iccid: '123',
              imei: '456',
              csq: '11',
              flow: '22',
              lac: '33',
              ci: '44',
              simInsert: true,
              netRegister: true,
              operator: 0,
            },
            configParam: {
              flowAlarm: false, //流量报警
              flowAlarmValue: 50, //
            },
            commParam: {
              name: '',
              baudRate: '',
              dataBits: '',
              stopBits: '',
              parity: '',
              timeout: '',
              interval: '',
              polling: '',
            },
          },
        ],
      }
    },
  },
  {
    url: '/network/mobile',
    method: 'post',
    response: () => {
      return {
        code: '0',
        message: '添加移动模块成功！',
        data: {},
      }
    },
  },
  {
    url: '/network/mobile',
    method: 'put',
    response: () => {
      return {
        code: '0',
        message: '修改移动模块成功！',
        data: {},
      }
    },
  },
  {
    url: '/network/mobile',
    method: 'delete',
    response: () => {
      return {
        code: '0',
        message: '删除移动模块成功！',
        data: {},
      }
    },
  },
  {
    url: '/ntp/param',
    method: 'get',
    response: () => {
      return {
        code: '0',
        message: '',
        data: {
          enable: false, //是否启用
          timeZone: 'UTC+8', //时区
          urlMaster: 'ntp1.aliyun.com', //主服务器URL
          urlSlave: 'xxx.xxx.com', //次服务器URL
        },
      }
    },
  },
  {
    url: '/ntp/param',
    method: 'put',
    response: () => {
      return {
        code: '0',
        message: '修改NTP服务成功！',
        data: {},
      }
    },
  },
  {
    url: '/product/sn',
    method: 'get',
    response: () => {
      return {
        code: '0',
        message: '',
        data: {
          sn: '11111',
        },
      }
    },
  },
  {
    url: '/product/sn',
    method: 'post',
    response: () => {
      return {
        code: '0',
        message: '操作成功！',
        data: '',
      }
    },
  },
  {
    url: '/tsl/content/block/properties',
    method: 'get',
    response: () => {
      return {
        code: '0',
        message: '操作成功！',
        data: [
          {
            name: 'property01',
            label: '属性01',
            accessMode: 0,
            type: 0,
            unit: '',
            decimals: 0,
            address: '30100',
            num: 10,
            ruleType: 'AABB',
          },
          {
            name: 'property02',
            label: '属性02',
            accessMode: 0,
            type: 0,
            unit: '',
            decimals: 0,
            address: '40100',
            num: 20,
            ruleType: 'AABB',
          },
        ],
      }
    },
  },
  {
    url: '/tsl/content/block/property',
    method: 'post',
    response: () => {
      return {
        code: '0',
        message: '添加块属性成功！',
        data: '',
      }
    },
  },
  {
    url: '/tsl/content/block/property',
    method: 'put',
    response: () => {
      return {
        code: '0',
        message: '修改块属性成功！',
        data: '',
      }
    },
  },
  {
    url: '/tsl/content/block/properties',
    method: 'delete',
    response: () => {
      return {
        code: '0',
        message: '删除块属性成功！',
        data: '',
      }
    },
  },
  {
    url: '/system/ping',
    method: 'post',
    response: () => {
      return {
        code: '0',
        message: '发送ping命令成功！',
        data: 'PING 58.215.221.98 (58.215.221.98): 56 data bytes\n64 bytes from 58.215.221.98: icmp_seq=0 ttl=128 time=4.041 ms\n64 bytes from 58.215.221.98: icmp_seq=2 ttl=128 time=4.482 ms\n\n--- 58.215.221.98 ping statistics ---\n4 packets transmitted, 4 packets received, 0.0% packet loss, 2 packets out of wait time\nround-trip min/avg/max/stddev = 4.041/11.184/30.259/11.036 ms\n',
      }
    },
  },

  {
    url: '/virtual/properties',
    method: 'get',
    response: () => {
      return {
        code: '0',
        message: '操作成功！',
        data: [
          {
            name: 'property01',
            label: '属性01',
            type: 0,
            unit: 't',
            decimals: 0,
            params: {
              collName: '1',
              deviceName: '1',
              propertyName: '1',
            },
          },
          {
            name: 'property02',
            label: '属性02',
            type: 0,
            unit: 'k',
            decimals: 0,
            params: {
              collName: '2',
              deviceName: '2',
              propertyName: '2',
            },
          },
        ],
      }
    },
  },
  {
    url: '/virtual/property',
    method: 'post',
    response: () => {
      return {
        code: '0',
        message: '添加虚拟设备属性成功！',
        data: '',
      }
    },
  },
  {
    url: '/virtual/property',
    method: 'put',
    response: () => {
      return {
        code: '0',
        message: '修改虚拟设备属性成功！',
        data: '',
      }
    },
  },
  {
    url: '/virtual/properties',
    method: 'delete',
    response: () => {
      return {
        code: '0',
        message: '删除虚拟设备属性成功！',
        data: '',
      }
    },
  },
]
