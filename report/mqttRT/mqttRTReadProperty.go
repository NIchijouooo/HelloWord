package mqttRT

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
	"gateway/virtual"
	"time"
)

type MQTTRTReadPropertyRequestParamTemplate struct {
	DeviceCode string   `json:"devCode"`
	Properties []string `json:"properties"`
}

type MQTTRTReadPropertyRequestTemplate struct {
	ClientID string                                   `json:"clientID"`
	MsgID    string                                   `json:"msgID"`
	Params   []MQTTRTReadPropertyRequestParamTemplate `json:"params"`
}

type MQTTRTReadPropertyAckParamPropertyTemplate struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

type MQTTRTReadPropertyAckParamTemplate struct {
	DeviceCode string                                       `json:"devCode"`
	Properties []MQTTRTReadPropertyAckParamPropertyTemplate `json:"varData"`
}

type MQTTRTReadPropertyAckTemplate struct {
	Code     int                                  `json:"code"`
	ClientID string                               `json:"clientID"`
	MsgID    string                               `json:"msgID"`
	Params   []MQTTRTReadPropertyAckParamTemplate `json:"data"`
}

func (r *ReportServiceParamRTTemplate) ReportServiceRTReadPropertyAck(reqFrame MQTTRTReadPropertyRequestTemplate, code int, ackParams []MQTTRTReadPropertyAckParamTemplate) {

	ackFrame := MQTTRTReadPropertyAckTemplate{
		Code:     code,
		ClientID: reqFrame.ClientID,
		MsgID:    reqFrame.MsgID,
		Params:   ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	propertyPostTopic := "/device/data/get_replay/" + r.GWParam.Param.ClientID
	aData, _ := json.Marshal(ackFrame)
	MQTTRTAddCommunicationMessage(r, "MQTT读属性应答包", Direction_TX, string(aData))

	setting.ZAPS.Infof("RT上报服务发布回复读属性应答消息主题 %s", propertyPostTopic)
	setting.ZAPS.Debugf("RT上报服务发布回复读属性应答消息内容 %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 1, false, sJson)
		if token.WaitTimeout(time.Duration(RTTimeOutReadProperty)*time.Second) == false {
			setting.ZAPS.Errorf("RT上报服务[%s]发布回复读属性消息失败 %v", r.GWParam.ServiceName, token.Error())
			return
		}
	} else {
		setting.ZAPS.Errorf("RT上报服务[%s]发布回复读属性消息失败", r.GWParam.ServiceName)
		return
	}
	setting.ZAPS.Debugf("RT上报服务[%s]发布回复读属性消息成功", r.GWParam.ServiceName)
}

func (r *ReportServiceParamRTTemplate) ReportServiceRTProcessReadProperty(reqFrame MQTTRTReadPropertyRequestTemplate) {

	ReadStatus := false

	ackParams := make([]MQTTRTReadPropertyAckParamTemplate, 0)
	rData, _ := json.Marshal(reqFrame)
	MQTTRTAddCommunicationMessage(r, "MQTT读属性请求包", Direction_RX, string(rData))

	for _, req := range reqFrame.Params {
		if req.DeviceCode == r.GWParam.Param.ClientID { //获取网关数据
			ReadStatus = true
			ackParam := r.ReportServiceRTProcessReadGWProperty(req)
			ackParams = append(ackParams, ackParam)
		} else { //获取网关设备数据
			status, params := r.ReportServiceRTProcessReadNodeProperty(req)
			ackParams = append(ackParams, params...)
			ReadStatus = status
		}

		if ReadStatus == true {
			r.ReportServiceRTReadPropertyAck(reqFrame, 0, ackParams)
		} else {
			r.ReportServiceRTReadPropertyAck(reqFrame, 1, ackParams)
		}
	}
}

func (r *ReportServiceParamRTTemplate) ReportServiceRTProcessReadGWProperty(reqParam MQTTRTReadPropertyRequestParamTemplate) MQTTRTReadPropertyAckParamTemplate {
	ackParam := MQTTRTReadPropertyAckParamTemplate{
		DeviceCode: r.GWParam.Param.ClientID,
	}
	property := MQTTRTReadPropertyAckParamPropertyTemplate{}
	timeStamp := time.Now().Unix()
	for _, name := range reqParam.Properties {
		switch name {
		case "ICCID":
			{
				property.Name = "ICCID"
				property.Timestamp = timeStamp
				property.Value = setting.MobileModule.RunParam.ICCID
				ackParam.Properties = append(ackParam.Properties, property)
			}
		}
	}

	return ackParam
}

func (r *ReportServiceParamRTTemplate) ReportServiceRTProcessReadNodeProperty(reqParam MQTTRTReadPropertyRequestParamTemplate) (bool, []MQTTRTReadPropertyAckParamTemplate) {

	status := false
	ackParams := make([]MQTTRTReadPropertyAckParamTemplate, 0)

	for _, rpNode := range r.NodeList {
		//从上报节点中找到相应节点
		if reqParam.DeviceCode == rpNode.Param.DeviceCode {
			nameMap := make([]string, 0)
			ackParam := MQTTRTReadPropertyAckParamTemplate{
				DeviceCode: rpNode.Param.DeviceCode,
			}
			property := MQTTRTReadPropertyAckParamPropertyTemplate{}
			timeStamp := time.Now().Unix()

			for _, proName := range reqParam.Properties {
				nameMap = append(nameMap, proName)
			}
			if rpNode.CollInterfaceName == "virtual" {
				vrNode, ok := virtual.VirtualDevice.Nodes[rpNode.Name]
				if !ok {
					continue
				}
				vrNode.VirtualDeviceGetVariables(nameMap)
				status = true
				for _, name := range reqParam.Properties {
					for _, p := range rpNode.Properties {
						//查找上报模型属性和mqtt通信属性对应关系
						if name == p.UploadName {
							for _, variable := range virtual.VirtualDevice.Nodes[rpNode.Name].Properties {
								//setting.ZAPS.Debugf("上报属性名称 %v，设备属性名称 %v",name,variable.Name)
								//查找上报模型属性和采集模型属性对应关系
								if p.Name == variable.Name {
									//setting.ZAPS.Debugf("上报属性名称 %v，设备属性名称 %v,设备属性个数 %v",name,variable.Name,len(variable.Value))
									if len(variable.Value) >= 1 {
										index := len(variable.Value) - 1
										property.Name = p.UploadName
										property.Timestamp = timeStamp
										property.Value = variable.Value[index].Value
										ackParam.Properties = append(ackParam.Properties, property)
									}
								}
							}
						}
					}
				}
			} else {
				coll, ok := device.CollectInterfaceMap.Coll[rpNode.CollInterfaceName]
				if !ok {
					continue
				}

				node, ok := coll.DeviceNodeMap[rpNode.Name]
				if !ok {
					continue
				}
				//从采集服务中找到相应节点
				cmd := device.CommunicationCmdTemplate{}
				cmd.CollInterfaceName = rpNode.CollInterfaceName
				cmd.DeviceName = rpNode.Name
				cmd.FunName = "GetRealVariables"
				paramStr, _ := json.Marshal(nameMap)
				cmd.FunPara = string(paramStr)

				ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
				if ackData.Status {
					status = true
					for _, name := range reqParam.Properties {
						for _, p := range rpNode.Properties {
							//查找上报模型属性和mqtt通信属性对应关系
							if name == p.UploadName {
								for _, variable := range node.Properties {
									//setting.ZAPS.Debugf("上报属性名称 %v，设备属性名称 %v",name,variable.Name)
									//查找上报模型属性和采集模型属性对应关系
									if p.Name == variable.Name {
										//setting.ZAPS.Debugf("上报属性名称 %v，设备属性名称 %v,设备属性个数 %v",name,variable.Name,len(variable.Value))
										if len(variable.Value) >= 1 {
											index := len(variable.Value) - 1
											property.Name = p.UploadName
											property.Timestamp = timeStamp
											property.Value = variable.Value[index].Value
											ackParam.Properties = append(ackParam.Properties, property)
										}
									}
								}
							}
						}
					}
				} else {
					status = false
					for _, name := range reqParam.Properties {
						property.Name = name
						property.Value = -1
						ackParam.Properties = append(ackParam.Properties, property)
					}
				}
				ackParams = append(ackParams, ackParam)
			}
		}
	}

	return status, ackParams
}
