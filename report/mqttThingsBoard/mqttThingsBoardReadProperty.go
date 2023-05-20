package mqttThingsBoard

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
	"time"
)

type MQTTThingsBoardReadPropertyRequestParamPropertyTemplate struct {
	Name string `json:"name"`
}

type MQTTThingsBoardReadPropertyRequestParamTemplate struct {
	ClientID   string                                                    `json:"clientID"`
	Properties []MQTTThingsBoardReadPropertyRequestParamPropertyTemplate `json:"properties"`
}

type MQTTThingsBoardReadPropertyRequestTemplate struct {
	ID      string                                            `json:"id"`
	Version string                                            `json:"version"`
	Ack     int                                               `json:"ack"`
	Params  []MQTTThingsBoardReadPropertyRequestParamTemplate `json:"params"`
}

type MQTTThingsBoardReadPropertyAckParamPropertyTemplate struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

type MQTTThingsBoardReadPropertyAckParamTemplate struct {
	ClientID   string                                                `json:"clientID"`
	Properties []MQTTThingsBoardReadPropertyAckParamPropertyTemplate `json:"properties"`
}

type MQTTThingsBoardReadPropertyAckTemplate struct {
	ID      string                                        `json:"id"`
	Version string                                        `json:"version"`
	Code    int                                           `json:"code"`
	Params  []MQTTThingsBoardReadPropertyAckParamTemplate `json:"params"`
}

func (r *ReportServiceParamThingsBoardTemplate) ReportServiceThingsBoardReadPropertyAck(reqFrame MQTTThingsBoardReadPropertyRequestTemplate, code int, ackParams []MQTTThingsBoardReadPropertyAckParamTemplate) {

	ackFrame := MQTTThingsBoardReadPropertyAckTemplate{
		ID:      reqFrame.ID,
		Version: reqFrame.Version,
		Code:    code,
		Params:  ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	propertyPostTopic := "/sys/thing/event/property/get_reply/" + r.GWParam.Param.ClientID

	setting.ZAPS.Infof("property get_reply topic: %s", propertyPostTopic)
	setting.ZAPS.Debugf("property get_reply: %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson)
		token.Wait()
	}
}

func (r *ReportServiceParamThingsBoardTemplate) ReportServiceThingsBoardProcessReadProperty(reqFrame MQTTThingsBoardReadPropertyRequestTemplate) {

	ReadStatus := false

	ackParams := make([]MQTTThingsBoardReadPropertyAckParamTemplate, 0)

	for _, v := range reqFrame.Params {
		for _, node := range r.NodeList {
			if v.ClientID == node.Param.DeviceName {
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
						ackParam := MQTTThingsBoardReadPropertyAckParamTemplate{
							ClientID: node.Param.DeviceName,
						}
						property := MQTTThingsBoardReadPropertyAckParamPropertyTemplate{}
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
		r.ReportServiceThingsBoardReadPropertyAck(reqFrame, 0, ackParams)
	} else {
		r.ReportServiceThingsBoardReadPropertyAck(reqFrame, 1, ackParams)
	}
}
