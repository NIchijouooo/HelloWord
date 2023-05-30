package mqttFeisjy

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
)

type MQTTFeisjyWritePropertyRequestParamPropertyTemplate struct {
	Code  string      `json:"code"`
	Value interface{} `json:"value"`
}

type MQTTFeisjyWritePropertyTemplate struct {
	Uuid       string                                                `json:"uuid"`
	DeviceAddr string                                                `json:"deviceAddr"`
	Properties []MQTTFeisjyWritePropertyRequestParamPropertyTemplate `json:"properties"`
}

func (r *ReportServiceParamFeisjyTemplate) ReportServiceFeisjyProcessWriteProperty(reqFrame MQTTFeisjyWritePropertyTemplate) {
	serviceInfo := struct {
		CollInterfaceName string                 `json:"collInterfaceName"`
		DeviceName        string                 `json:"deviceName"`
		ServiceName       string                 `json:"serviceName"`
		ServiceParam      map[string]interface{} `json:"serviceParam"`
	}{}

	setting.ZAPS.Infof("reqFrame %v", reqFrame)

	serviceInfo.ServiceParam = make(map[string]interface{}, 1)
	for _, node := range r.NodeList {
		if node.Param.DeviceID == reqFrame.DeviceAddr {
			serviceInfo.CollInterfaceName = node.CollInterfaceName
			serviceInfo.DeviceName = node.Name
			serviceInfo.ServiceName = node.ServiceName
			for _, properties := range reqFrame.Properties {
				serviceInfo.ServiceParam[properties.Code] = properties.Value
			}
		}
	}

	cmd := device.CommunicationCmdTemplate{}
	cmd.CollInterfaceName = serviceInfo.CollInterfaceName
	cmd.DeviceName = serviceInfo.DeviceName
	cmd.FunName = "SetVariables"
	paramStr, _ := json.Marshal(serviceInfo.ServiceParam)
	cmd.FunPara = string(paramStr)

	setting.ZAPS.Infof("cmd %v", cmd)
	coll, ok := device.CollectInterfaceMap.Coll[serviceInfo.CollInterfaceName]
	if !ok {
		setting.ZAPS.Errorf("ReportServiceFeisjyProcessWriteProperty eer")
		return
	}
	cmdRX := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
	if cmdRX.Status == true {
		setting.ZAPS.Info("成功")
	} else {
		setting.ZAPS.Info("失败")
	}
}
