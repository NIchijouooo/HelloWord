package mqttZxJs

import (
	"encoding/json"
	"gateway/device"
	"gateway/report/mqttFeisjy"
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
		r.ZxjsPublishSetCack(serviceInfo, 0)
		setting.ZAPS.Infof("set content is empty serviceInfo = %v", serviceInfo)
		return
	}

	// 初始化命令集合
	var cmdList []device.CommunicationCmdTemplate
	for _, v := range r.NodeList {
		initCmdList(&cmdList, content, v.CollInterfaceName, v.Name)
	}
	for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
		for _, v := range v.NodeList {
			initCmdList(&cmdList, content, v.CollInterfaceName, v.Name)
		}
	}

	if len(cmdList) == 0 {
		r.ZxjsPublishSetCack(serviceInfo, 0)
		setting.ZAPS.Infof("set content is empty serviceInfo = %v", serviceInfo)
		return
	}

	valid := 0
	for _, cmd := range cmdList {
		//获取采集接口数据
		coll, collErr := device.CollectInterfaceMap.Coll[cmd.CollInterfaceName]
		if !collErr {
			setting.ZAPS.Debugf("coll接口[%s]不存在", cmd.CollInterfaceName)
			continue
		}
		setting.ZAPS.Infof("数据解析完毕，即将修改属性。FunName [%s] 设备名[%s] 采集接口名称[%s] 属性参数 %v", cmd.FunPara, cmd.DeviceName, cmd.CollInterfaceName, cmd.FunPara)
		cmdRX := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
		if cmdRX.Status == true {
			setting.ZAPS.Infof("修改属性成功 DevName [%s]", cmd.DeviceName)
			valid = 1
		} else {
			setting.ZAPS.Infof("修改属性失败 DevName [%s]", cmd.DeviceName)
		}
	}

	r.ZxjsPublishSetCack(serviceInfo, valid)
}

/*
*
初始化控制指令集合
*/
func initCmdList(cmdList *[]device.CommunicationCmdTemplate, content map[string]interface{}, collInterfaceName string, deviceName string) {
	param := make(map[string]interface{}, 0)
	//获取采集接口数据
	coll, collErr := device.CollectInterfaceMap.Coll[collInterfaceName]
	if !collErr {
		setting.ZAPS.Debugf("coll接口[%s]不存在", collInterfaceName)
		return
	}

	//获取节点数据
	node, nodeErr := coll.DeviceNodeMap[deviceName]
	if !nodeErr {
		setting.ZAPS.Debugf("coll接口[%s]下的设备[%s]不存在", collInterfaceName, deviceName)
		return
	}
	for _, properties := range node.Properties {
		val, ok := content[properties.Identity]
		if !ok {
			continue
		}
		param[properties.Label] = val
	}

	cmd := device.CommunicationCmdTemplate{}
	cmd.CollInterfaceName = collInterfaceName
	cmd.DeviceName = deviceName
	cmd.FunName = "SetVariables"
	paramStr, _ := json.Marshal(param)
	cmd.FunPara = string(paramStr)
	*cmdList = append(*cmdList, cmd)
}
