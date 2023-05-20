package mqttRT

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
	"gateway/virtual"
	"time"
)

type MQTTRTWritePropertyRequestParamPropertyTemplate struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type MQTTRTWritePropertyRequestParamTemplate struct {
	DeviceCode string                                            `json:"devCode"`
	Properties []MQTTRTWritePropertyRequestParamPropertyTemplate `json:"varData"`
}

type MQTTRTWritePropertyRequestTemplate struct {
	MsgID    string                                    `json:"msgID"`
	ClientID string                                    `json:"clientID"`
	Params   []MQTTRTWritePropertyRequestParamTemplate `json:"params"`
}

type MQTTRTWritePropertyAckTemplate struct {
	Code     int                                       `json:"code"`
	MsgID    string                                    `json:"msgID"`
	ClientID string                                    `json:"clientID"`
	Params   []MQTTRTWritePropertyRequestParamTemplate `json:"params"`
}

func (r *ReportServiceParamRTTemplate) ReportServiceRTWritePropertyAck(reqFrame MQTTRTWritePropertyRequestTemplate, code int, ackParams []MQTTRTWritePropertyRequestParamTemplate) {

	ackFrame := MQTTRTWritePropertyAckTemplate{
		Code:     code,
		MsgID:    reqFrame.MsgID,
		ClientID: reqFrame.ClientID,
		Params:   ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	propertyPostTopic := "/device/data/set_replay/" + r.GWParam.Param.ClientID

	MQTTRTAddCommunicationMessage(r, "MQTT写属性应答包", Direction_TX, string(sJson))
	setting.ZAPS.Infof("RT上报服务发布回复写属性应答消息主题 %s", propertyPostTopic)
	setting.ZAPS.Debugf("RT上报服务发布回复些属性应答消息内容 %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 1, false, sJson)
		if token.WaitTimeout(time.Duration(RTTimeOutWriteProperty)*time.Second) == false {
			setting.ZAPS.Errorf("RT上报服务[%s]发布回应写属性消息失败 %v", r.GWParam.ServiceName, token.Error())
			return
		}
	} else {
		setting.ZAPS.Errorf("RT上报服务[%s]发布回应写属性消息失败", r.GWParam.ServiceName)
		return
	}
	setting.ZAPS.Debugf("RT上报服务[%s]发布回应写属性消息成功", r.GWParam.ServiceName)
}

func (r *ReportServiceParamRTTemplate) ReportServiceRTProcessWriteProperty(reqFrame MQTTRTWritePropertyRequestTemplate) {

	ackParams := make([]MQTTRTWritePropertyRequestParamTemplate, 0)
	writeStatus := false

	rData, _ := json.Marshal(reqFrame)
	MQTTRTAddCommunicationMessage(r, "MQTT写属性请求包", Direction_RX, string(rData))

	for _, v := range reqFrame.Params {
		if v.DeviceCode == r.GWParam.Param.ClientID {
			setting.ZAPS.Debugf("RT上报服务[%s]写网关变量", r.GWParam.ServiceName)
			writeStatus, ackParams = r.ReportServiceRTProcessWriteGWProperty(v)
		} else {
			writeStatus, ackParams = r.ReportServiceRTProcessWriteNodeProperty(v)
		}
	}

	if writeStatus == true {
		r.ReportServiceRTWritePropertyAck(reqFrame, 0, ackParams)
	} else {
		r.ReportServiceRTWritePropertyAck(reqFrame, 1, ackParams)
	}
}

func (r *ReportServiceParamRTTemplate) ReportServiceRTProcessWriteGWProperty(reqFrame MQTTRTWritePropertyRequestParamTemplate) (bool, []MQTTRTWritePropertyRequestParamTemplate) {
	ackParams := make([]MQTTRTWritePropertyRequestParamTemplate, 0)
	writeStatus := false

	reqValueMap := make(map[string]interface{})
	for _, pro := range reqFrame.Properties {
		reqValueMap[pro.Name] = pro.Value
	}
	for _, n := range virtual.VirtualDevice.Nodes {
		for _, v := range n.Properties {
			reqProperty, ok := reqValueMap[v.Name]
			if !ok {
				continue
			}

			//从上报节点中找到相应节点
			coll, ok := device.CollectInterfaceMap.Coll[v.Params.CollName]
			if !ok {
				continue
			}

			for _, d := range coll.DeviceNodeMap {
				if d.Name == v.Params.DeviceName {
					//从采集服务中找到相应节点
					cmd := device.CommunicationCmdTemplate{}
					cmd.CollInterfaceName = v.Params.CollName
					cmd.DeviceName = v.Params.DeviceName
					cmd.FunName = "SetVariables"
					valueMap := make(map[string]interface{})
					valueMap[v.Params.PropertyName] = reqProperty
					//setting.ZAPS.Debugf("value %+v", valueMap)

					paramStr, _ := json.Marshal(valueMap)
					cmd.FunPara = string(paramStr)
					param := MQTTRTWritePropertyRequestParamTemplate{
						DeviceCode: reqFrame.DeviceCode,
					}
					property := MQTTRTWritePropertyRequestParamPropertyTemplate{}

					ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
					if ackData.Status {
						writeStatus = true
						for _, p := range reqFrame.Properties {
							property.Name = p.Name
							property.Value = 0
							param.Properties = append(param.Properties, property)
						}
					} else {
						writeStatus = false
						for _, p := range reqFrame.Properties {
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

	return writeStatus, ackParams
}

func (r *ReportServiceParamRTTemplate) ReportServiceRTProcessWriteNodeProperty(reqFrame MQTTRTWritePropertyRequestParamTemplate) (bool, []MQTTRTWritePropertyRequestParamTemplate) {

	ackParams := make([]MQTTRTWritePropertyRequestParamTemplate, 0)
	writeStatus := false

	wProperyMap := make(map[string]interface{})
	for _, pro := range reqFrame.Properties {
		wProperyMap[pro.Name] = pro.Value
	}
	setting.ZAPS.Debugf("上报服务[%v]写变量列表 %+v", r.GWParam.ServiceName, wProperyMap)

	for _, node := range r.NodeList {
		if reqFrame.DeviceCode == node.Param.DeviceCode {
			param := MQTTRTWritePropertyRequestParamTemplate{
				DeviceCode: node.Param.DeviceCode,
			}
			property := MQTTRTWritePropertyRequestParamPropertyTemplate{}

			if node.CollInterfaceName == "virtual" {
				vrNode, ok := virtual.VirtualDevice.Nodes[node.Name]
				if !ok {
					continue
				}
				wProperyResultMap := vrNode.VirtualDeviceSetVariables(wProperyMap)
				for _, p := range reqFrame.Properties {
					value, ok := wProperyResultMap[p.Name]
					if !ok {
						continue
					}
					property.Name = p.Name
					property.Value = value
					param.Properties = append(param.Properties, property)
				}
				ackParams = append(ackParams, param)
				writeStatus = true
			} else {
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
						paramStr, _ := json.Marshal(wProperyMap)
						cmd.FunPara = string(paramStr)

						ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
						if ackData.Status {
							writeStatus = true
							for _, p := range reqFrame.Properties {
								property.Name = p.Name
								property.Value = 0
								param.Properties = append(param.Properties, property)
							}
						} else {
							writeStatus = false
							for _, p := range reqFrame.Properties {
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

	return writeStatus, ackParams
}
