package mqttThingsBoard

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
)

type MQTTThingsBoardWritePropertyRequestParamPropertyTemplate struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type MQTTThingsBoardWritePropertyRequestParamTemplate struct {
	ClientID   string                                                     `json:"clientID"`
	Properties []MQTTThingsBoardWritePropertyRequestParamPropertyTemplate `json:"properties"`
}

type MQTTThingsBoardWritePropertyRequestTemplate struct {
	ID      string                                             `json:"id"`
	Version string                                             `json:"version"`
	Ack     int                                                `json:"ack"`
	Params  []MQTTThingsBoardWritePropertyRequestParamTemplate `json:"params"`
}

type MQTTThingsBoardWritePropertyAckTemplate struct {
	ID      string                                             `json:"id"`
	Version string                                             `json:"version"`
	Code    int                                                `json:"code"`
	Params  []MQTTThingsBoardWritePropertyRequestParamTemplate `json:"params"`
}

func (r *ReportServiceParamThingsBoardTemplate) ReportServiceThingsBoardWritePropertyAck(reqFrame MQTTThingsBoardWritePropertyRequestTemplate, code int, ackParams []MQTTThingsBoardWritePropertyRequestParamTemplate) {

	ackFrame := MQTTThingsBoardWritePropertyAckTemplate{
		ID:      reqFrame.ID,
		Version: reqFrame.Version,
		Code:    code,
		Params:  ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	propertyPostTopic := "/sys/thing/event/property/set_reply/" + r.GWParam.Param.ClientID

	setting.ZAPS.Infof("property set_reply topic: %s", propertyPostTopic)
	setting.ZAPS.Debugf("property set_reply: %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson)
		token.Wait()
	}
}

func (r *ReportServiceParamThingsBoardTemplate) ReportServiceThingsBoardProcessWriteProperty(reqFrame MQTTThingsBoardWritePropertyRequestTemplate) {

	writeStatus := false

	ackParams := make([]MQTTThingsBoardWritePropertyRequestParamTemplate, 0)

	for _, v := range reqFrame.Params {
		for _, node := range r.NodeList {
			if v.ClientID == node.Param.DeviceName {
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
						cmd.FunName = "SetVariables"
						valueMap := make(map[string]interface{})
						for _, pro := range v.Properties {
							valueMap[pro.Name] = pro.Value
						}
						paramStr, _ := json.Marshal(valueMap)
						cmd.FunPara = string(paramStr)
						param := MQTTThingsBoardWritePropertyRequestParamTemplate{
							ClientID: node.Param.DeviceName,
						}
						property := MQTTThingsBoardWritePropertyRequestParamPropertyTemplate{}

						ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
						if ackData.Status {
							writeStatus = true
							for _, p := range v.Properties {
								property.Name = p.Name
								property.Value = 0
								param.Properties = append(param.Properties, property)
							}
						} else {
							writeStatus = false
							for _, p := range v.Properties {
								property.Name = p.Name
								property.Value = 1
								param.Properties = append(param.Properties, property)
							}
						}
						ackParams = append(ackParams, param)
					}
				}
			}
		}
	}

	if writeStatus == true {
		r.ReportServiceThingsBoardWritePropertyAck(reqFrame, 0, ackParams)
	} else {
		r.ReportServiceThingsBoardWritePropertyAck(reqFrame, 1, ackParams)
	}
}
