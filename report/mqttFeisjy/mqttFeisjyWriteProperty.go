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
	CmdType    string                                                `json:"cmdType"`
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

	setting.ZAPS.Infof("[%s]数据解析完毕，即将修改属性。FunName [%s] 设备名[%s] 采集接口名称[%s] 属性参数 %v", reqFrame.CmdType, cmd.FunPara, cmd.DeviceName, cmd.CollInterfaceName, cmd.FunPara)
	coll, ok := device.CollectInterfaceMap.Coll[serviceInfo.CollInterfaceName]
	if !ok {
		setting.ZAPS.Errorf("ReportServiceFeisjyProcessWriteProperty eer")
		return
	}
	cmdRX := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
	if cmdRX.Status == true {
		setting.ZAPS.Infof("[%s] 修改属性成功 DevName [%s]", reqFrame.CmdType, cmd.DeviceName)
		r.ReportServiceFeisjyWritePropertyAck(reqFrame, 1)
	} else {
		setting.ZAPS.Infof("[%s] 修改属性失败 DevName [%s]", reqFrame.CmdType, cmd.DeviceName)
		r.ReportServiceFeisjyWritePropertyAck(reqFrame, 0)
	}
}

func (r *ReportServiceParamFeisjyTemplate) ReportServiceFeisjyWritePropertyAck(reqFrame MQTTFeisjyWritePropertyTemplate, status int) {

	type deviceControlResult struct {
		Uuid   string `json:"uuid"`
		Status int    `json:"status"`
	}

	v := deviceControlResult{
		Uuid:   reqFrame.Uuid,
		Status: status,
	}
	sJson, _ := json.Marshal(v)

	setting.ZAPS.Info(v)

	r.FeisjyPublishdeviceControlResult(sJson, reqFrame.DeviceAddr)
}
