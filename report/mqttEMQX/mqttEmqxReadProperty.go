package mqttEmqx

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
	"time"
)

type MQTTEmqxReadPropertyRequestParamPropertyTemplate struct {
	Name string `json:"name"`
}

type MQTTEmqxReadPropertyRequestParamTemplate struct {
	ClientID   string                                             `json:"clientID"`
	Properties []MQTTEmqxReadPropertyRequestParamPropertyTemplate `json:"properties"`
}

type MQTTEmqxReadPropertyRequestTemplate struct {
	ID      string                                     `json:"id"`
	Version string                                     `json:"version"`
	Ack     int                                        `json:"ack"`
	Params  []MQTTEmqxReadPropertyRequestParamTemplate `json:"params"`
}

type MQTTEmqxReadPropertyAckParamPropertyTemplate struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

type MQTTEmqxReadPropertyAckParamTemplate struct {
	ClientID   string                                         `json:"clientID"`
	Properties []MQTTEmqxReadPropertyAckParamPropertyTemplate `json:"properties"`
}

type MQTTEmqxReadPropertyAckTemplate struct {
	ID      string                                 `json:"id"`
	Version string                                 `json:"version"`
	Code    int                                    `json:"code"`
	Params  []MQTTEmqxReadPropertyAckParamTemplate `json:"params"`
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxReadPropertyAck(reqFrame MQTTEmqxReadPropertyRequestTemplate, code int, ackParams []MQTTEmqxReadPropertyAckParamTemplate) {

	ackFrame := MQTTEmqxReadPropertyAckTemplate{
		ID:      reqFrame.ID,
		Version: reqFrame.Version,
		Code:    code,
		Params:  ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	propertyPostTopic := "/sys/thing/node/property/get_reply/" + r.GWParam.Param.ClientID

	setting.ZAPS.Infof("property get_reply topic: %s", propertyPostTopic)
	setting.ZAPS.Debugf("property get_reply: %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 1, false, sJson)
		if token.WaitTimeout(5*time.Second) == false {
			setting.ZAPS.Errorf("EMQX上报服务[%s]发布回复读属性消息失败 %v", r.GWParam.ServiceName, token.Error())
			return
		}
	} else {
		setting.ZAPS.Errorf("EMQX上报服务[%s]发布回复读属性消息失败", r.GWParam.ServiceName)
		return
	}
	setting.ZAPS.Debugf("EMQX上报服务[%s]发布读属性消息成功", r.GWParam.ServiceName)
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxProcessReadProperty(reqFrame MQTTEmqxReadPropertyRequestTemplate) {

	ReadStatus := false

	ackParams := make([]MQTTEmqxReadPropertyAckParamTemplate, 0)

	for _, v := range reqFrame.Params {
		for _, node := range r.NodeList {
			if v.ClientID == node.Param.DeviceCode {
				coll, ok := device.CollectInterfaceMap.Coll[node.CollInterfaceName]
				if !ok {
					continue
				}
				//从上报节点中找到相应节点
				for _, n := range coll.DeviceNodeMap {
					if n.Name == node.Name {
						//从采集服务中找到相应节点
						cmd := device.CommunicationCmdTemplate{}
						cmd.CollInterfaceName = node.CollInterfaceName
						cmd.DeviceName = node.Name
						cmd.FunName = "GetRealVariables"
						nameMap := make([]string, 0)
						for _, pro := range v.Properties {
							nameMap = append(nameMap, pro.Name)
						}
						paramStr, _ := json.Marshal(nameMap)
						cmd.FunPara = string(paramStr)
						ackParam := MQTTEmqxReadPropertyAckParamTemplate{
							ClientID: node.Param.DeviceCode,
						}
						property := MQTTEmqxReadPropertyAckParamPropertyTemplate{}
						timeStamp := time.Now().Unix()
						ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
						if ackData.Status {
							ReadStatus = true
							for _, p := range v.Properties {
								for _, variable := range n.Properties {
									if p.Name == variable.Name {
										if len(variable.Value) >= 1 {
											index := len(variable.Value) - 1
											property.Name = variable.Name
											property.Timestamp = timeStamp
											property.Value = variable.Value[index].Value
											ackParam.Properties = append(ackParam.Properties, property)
										}
									}
								}
							}
						} else {
							ReadStatus = false
							for _, p := range v.Properties {
								property.Name = p.Name
								property.Value = -1
								ackParam.Properties = append(ackParam.Properties, property)
							}
						}
						ackParams = append(ackParams, ackParam)
					}
				}
			}
		}
	}

	if ReadStatus == true {
		r.ReportServiceEmqxReadPropertyAck(reqFrame, 0, ackParams)
	} else {
		r.ReportServiceEmqxReadPropertyAck(reqFrame, 1, ackParams)
	}
}
