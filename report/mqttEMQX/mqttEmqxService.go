package mqttEmqx

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
	"time"
)

type MQTTEmqxInvokeServiceAckParamTemplate struct {
	ClientID  string `json:"clientID"`
	CmdName   string `json:"cmdName"`
	CmdStatus int    `json:"cmdStatus"`
}

type MQTTEmqxInvokeServiceAckTemplate struct {
	ID      string                                  `json:"id"`
	Version string                                  `json:"version"`
	Code    int                                     `json:"code"`
	Params  []MQTTEmqxInvokeServiceAckParamTemplate `json:"params"`
}

type MQTTEmqxInvokeServiceRequestParamTemplate struct {
	ClientID  string                 `json:"clientID"`
	CmdName   string                 `json:"cmdName"`
	CmdParams map[string]interface{} `json:"cmdParams"`
}

type MQTTEmqxInvokeServiceRequestTemplate struct {
	ID      string                                      `json:"id"`
	Version string                                      `json:"version"`
	Ack     int                                         `json:"ack"`
	Params  []MQTTEmqxInvokeServiceRequestParamTemplate `json:"params"`
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxInvokeServiceAck(reqFrame MQTTEmqxInvokeServiceRequestTemplate, deviceType int, code int, ackParams []MQTTEmqxInvokeServiceAckParamTemplate) {

	ackFrame := MQTTEmqxInvokeServiceAckTemplate{
		ID:      reqFrame.ID,
		Version: reqFrame.Version,
		Code:    code,
		Params:  ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	serviceInvokeTopic := ""
	if deviceType == DeviceTypeGW {
		serviceInvokeTopic = "/sys/thing/gw/service/invoke_reply/" + r.GWParam.Param.ClientID
	} else {
		serviceInvokeTopic = "/sys/thing/node/service/invoke_reply/" + r.GWParam.Param.ClientID
	}

	setting.ZAPS.Infof("service invoke_reply topic: %s", serviceInvokeTopic)
	setting.ZAPS.Debugf("service invoke_reply: %v", string(sJson))

	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(serviceInvokeTopic, 1, false, sJson)
		if token.WaitTimeout(5*time.Second) == false {
			setting.ZAPS.Errorf("EMQX上报服务[%s]发布回复服务调用消息失败 %v", r.GWParam.ServiceName, token.Error())
			return
		}

	} else {
		setting.ZAPS.Errorf("EMQX上报服务[%s]发布回复服务调用消息失败", r.GWParam.ServiceName)
		return
	}
	setting.ZAPS.Debugf("EMQX上报服务[%s]发布回复服务调用消息成功", r.GWParam.ServiceName)
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxProcessInvokeService(reqFrame MQTTEmqxInvokeServiceRequestTemplate) {

	ReadStatus := false

	ackParams := make([]MQTTEmqxInvokeServiceAckParamTemplate, 0)

	for _, v := range reqFrame.Params {
		for _, node := range r.NodeList {
			if v.ClientID == node.Param.DeviceCode {
				//从上报节点中找到相应节点
				coll, ok := device.CollectInterfaceMap.Coll[node.CollInterfaceName]
				if !ok {
					continue
				}

				for _, n := range coll.DeviceNodeMap {
					if n.Name == node.Name {
						//从采集服务中找到相应节点
						cmd := device.CommunicationCmdTemplate{}
						cmd.CollInterfaceName = node.CollInterfaceName
						cmd.DeviceName = node.Name
						cmd.FunName = v.CmdName
						paramStr, _ := json.Marshal(v.CmdParams)
						cmd.FunPara = string(paramStr)
						ackParam := MQTTEmqxInvokeServiceAckParamTemplate{
							ClientID: node.Param.DeviceCode,
							CmdName:  v.CmdName,
						}

						ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
						if ackData.Status {
							ReadStatus = true
							ackParam.CmdStatus = 0
						} else {
							ReadStatus = false
							ackParam.CmdStatus = 1
						}
						ackParams = append(ackParams, ackParam)
					}
				}
			}
		}
	}

	if ReadStatus == true {
		r.ReportServiceEmqxInvokeServiceAck(reqFrame, DeviceTypeNode, 0, ackParams)
	} else {
		r.ReportServiceEmqxInvokeServiceAck(reqFrame, DeviceTypeNode, 1, ackParams)
	}
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxProcessInvokeGWService(reqFrame MQTTEmqxInvokeServiceRequestTemplate) {

	ackParams := make([]MQTTEmqxInvokeServiceAckParamTemplate, 0)

	index := -1
	for k, v := range reqFrame.Params {
		if v.ClientID == r.GWParam.Param.ClientID {
			index = k
			break
		}
	}
	if index != -1 {
		return
	}

	r.ReportServiceEmqxInvokeServiceAck(reqFrame, DeviceTypeGW, 0, ackParams)

	//if reqFrame.Params[index].CmdName == "reboot" {
	//	setting.SystemReboot()
	//} else if reqFrame.Params[index].CmdName == "lockCmd" {
	//	cmdTmp, ok := reqFrame.Params[index].CmdParams["cmd"]
	//	if !ok {
	//		return
	//	}
	//	switch cmd := cmdTmp.(type) {
	//	case string:
	//		{
	//			cmdInt, _ := strconv.Atoi(cmd)
	//			setting.SystemLock(cmdInt)
	//		}
	//	case int:
	//		{
	//			setting.SystemLock(cmd)
	//		}
	//	}
	//}
}
