package mqttZxJs

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
)

/*
*
设备控制对象
*/
type MQTTZxjsControlTemplate struct {
	DeviceSN string   `json:"deviceSN"`
	Seq      int      `json:"seq"`
	Ts       int64    `json:"ts"`
	Data     DataInfo `json:"data"`
	Type     string   `json:"type"`
}

type DataInfo struct {
	Valid   int                    `json:"valid"`
	Model   int                    `json:"model"`
	Content map[string]interface{} `json:"content"`
}

func (r *ReportServiceParamZxjsTemplate) ReportServiceZxjsProcessWriteProperty(serviceInfo MQTTZxjsControlTemplate) {

	content := serviceInfo.Data.Content
	if len(content) == 0 {
		setting.ZAPS.Infof("set content is empty serviceInfo = %v", serviceInfo)
		return
	}

	param := make(map[string]interface{}, 1)
	var collName, deviceName string
	for _, v := range r.NodeList {
		//获取采集接口数据
		coll, collErr := device.CollectInterfaceMap.Coll[v.CollInterfaceName]
		if !collErr {
			setting.ZAPS.Debugf("coll接口[%s]不存在", v.CollInterfaceName)
			continue
		}

		//获取节点数据
		node, nodeErr := coll.DeviceNodeMap[v.Name]
		if !nodeErr {
			setting.ZAPS.Debugf("coll接口[%s]下的设备[%s]不存在", v.CollInterfaceName, v.Name)
			continue
		}
		for _, properties := range node.Properties {
			val, ok := content[properties.Identity]
			if !ok {
				continue
			}
			collName = v.CollInterfaceName
			deviceName = v.Name
			param[properties.Label] = val
		}
	}

	cmd := device.CommunicationCmdTemplate{}
	cmd.CollInterfaceName = collName
	cmd.DeviceName = deviceName
	cmd.FunName = "SetVariables"
	paramStr, _ := json.Marshal(param)
	cmd.FunPara = string(paramStr)

	setting.ZAPS.Infof("数据解析完毕，即将修改属性。FunName [%s] 设备名[%s] 采集接口名称[%s] 属性参数 %v", cmd.FunPara, cmd.DeviceName, cmd.CollInterfaceName, cmd.FunPara)
	coll, ok := device.CollectInterfaceMap.Coll[collName]
	if !ok {
		setting.ZAPS.Errorf("ReportServiceFeisjyProcessWriteProperty eer")
		return
	}
	cmdRX := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
	valid := 0
	if cmdRX.Status == true {
		setting.ZAPS.Infof("修改属性成功 DevName [%s]", cmd.DeviceName)
		valid = 1
	} else {
		setting.ZAPS.Infof("修改属性失败 DevName [%s]", cmd.DeviceName)
	}
	r.ZxjsPublishSetCack(serviceInfo, valid)
}
