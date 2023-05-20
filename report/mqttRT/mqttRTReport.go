package mqttRT

import (
	"encoding/json"
	"gateway/device"
	"gateway/report/reportModel"
	"gateway/setting"
	"gateway/virtual"
	"time"
)

type MQTTRTReportPropertyTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
}

type MQTTRTPropertyPostParamPropertyTemplate struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	TimeStamp int64       `json:"timestamp"`
}

type MQTTRTReportAlarmTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
	Properties []MQTTRTPropertyPostParamPropertyTemplate
}

type MQTTRTPropertyPostParamTemplate struct {
	DeviceCode string                                    `json:"devCode"`
	OnLine     bool                                      `json:"onLine"`
	TypeCode   string                                    `json:"typeCode"`
	VarData    []MQTTRTPropertyPostParamPropertyTemplate `json:"varData"`
}

type MQTTRTPropertyPostTemplate struct {
	ClientID string                            `json:"clientID"`
	Params   []MQTTRTPropertyPostParamTemplate `json:"params"`
}

type MQTTRTReportPropertyAckTemplate struct {
	Code    int32  `json:"code"`
	Data    string `json:"-"`
	ID      string `json:"id"`
	Message string `json:"message"`
	Method  string `json:"method"`
	Version string `json:"version"`
}

const (
	DeviceTypeGW = iota
	DeviceTypeNode
)

func (r *ReportServiceParamRTTemplate) MQTTRTPropertyPost(rType int, gwParam ReportServiceGWParamRTTemplate, propertyParam []MQTTRTPropertyPostParamTemplate, alarm bool) (int, bool) {

	propertyPost := MQTTRTPropertyPostTemplate{
		ClientID: gwParam.Param.ClientID,
		Params:   propertyParam,
	}
	MsgID++

	sJson, _ := json.Marshal(propertyPost)
	propertyPostTopic := ""
	if rType == DeviceTypeGW {
		propertyPostTopic = "/device/data/post/gw/" + gwParam.Param.ClientID
	} else {
		propertyPostTopic = "/device/data/post/" + gwParam.Param.ClientID
	}

	setting.ZAPS.Infof("上报服务[%s]发布上报消息主题%s", gwParam.ServiceName, propertyPostTopic)
	MQTTRTAddCommunicationMessage(r, propertyPostTopic, Direction_TX, string(sJson))

	if gwParam.MQTTClient != nil {
		token := gwParam.MQTTClient.Publish(propertyPostTopic, 1, false, sJson)
		if token.WaitTimeout(time.Duration(RTTimeOutReportProperty)*time.Second) == false {
			if alarm == false {
				setting.ZAPS.Errorf("RT上报服务[%s]发布上报属性消息失败 %v", gwParam.ServiceName, token.Error())
			} else {
				setting.ZAPS.Errorf("RT上报服务[%s]发布上报告警消息失败 %v", gwParam.ServiceName, token.Error())
			}
			return MsgID, false
		}
	} else {
		if alarm == false {
			setting.ZAPS.Errorf("RT上报服务[%s]发布上报属性消息失败", gwParam.ServiceName)
		} else {
			setting.ZAPS.Errorf("RT上报服务[%s]发布上报告警消息失败", gwParam.ServiceName)
		}
		return MsgID, false
	}

	if alarm == false {
		setting.ZAPS.Debugf("RT上报服务[%s]发布上报属性消息成功", gwParam.ServiceName)
	} else {
		setting.ZAPS.Debugf("RT上报服务[%s]发布上报告警消息成功", gwParam.ServiceName)
	}
	return MsgID, true
}

func (r *ReportServiceParamRTTemplate) GWPropertyPost() {

	propertyMap := make([]MQTTRTPropertyPostParamPropertyTemplate, 0)

	property := MQTTRTPropertyPostParamPropertyTemplate{}

	timeStamp := time.Now().Unix()

	//property.Name = "MemTotal"
	//property.Value = setting.SystemState.MemTotal
	//property.TimeStamp = timeStamp
	//propertyMap = append(propertyMap, property)
	//
	//property.Name = "MemUse"
	//property.Value = setting.SystemState.MemUse
	//property.TimeStamp = timeStamp
	//propertyMap = append(propertyMap, property)
	//
	//property.Name = "DiskTotal"
	//property.Value = setting.SystemState.DiskTotal
	//property.TimeStamp = timeStamp
	//propertyMap = append(propertyMap, property)
	//
	//property.Name = "DiskUse"
	//property.Value = setting.SystemState.DiskUse
	//property.TimeStamp = timeStamp
	//propertyMap = append(propertyMap, property)

	property.Name = "Name"
	property.Value = setting.SystemState.Name
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "SN"
	property.Value = setting.SystemState.SN
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)
	//
	//property.Name = "HardVer"
	//property.Value = setting.SystemState.HardVer
	//property.TimeStamp = timeStamp
	//propertyMap = append(propertyMap, property)

	property.Name = "SoftVer"
	property.Value = setting.SystemState.SoftVer
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "GOARCH"
	property.Value = setting.SystemState.GOARCH
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	//
	//property.Name = "SystemRTC"
	//property.Value = setting.SystemState.SystemRTC
	//property.TimeStamp = timeStamp
	//propertyMap = append(propertyMap, property)
	//
	//property.Name = "RunTime"
	//property.Value = setting.SystemState.RunTime
	//property.TimeStamp = timeStamp
	//propertyMap = append(propertyMap, property)
	//
	//property.Name = "DeviceOnline"
	//property.Value = setting.SystemState.DeviceOnline
	//property.TimeStamp = timeStamp
	//propertyMap = append(propertyMap, property)
	//
	//property.Name = "DevicePacketLoss"
	//property.Value = setting.SystemState.DevicePacketLoss
	//property.TimeStamp = timeStamp
	//propertyMap = append(propertyMap, property)

	property.Name = "CSQ"
	property.Value = setting.MobileModule.RunParam.CSQ
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "ICCID"
	property.Value = setting.MobileModule.RunParam.ICCID
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	//清空接收缓存
	for i := 0; i < len(r.ReceiveReportPropertyAckFrameChan); i++ {
		<-r.ReceiveReportPropertyAckFrameChan
	}

	propertyPostParam := MQTTRTPropertyPostParamTemplate{
		DeviceCode: r.GWParam.Param.ClientID,
		OnLine:     true,
		TypeCode:   "",
		VarData:    propertyMap,
	}

	propertyPostParamMap := make([]MQTTRTPropertyPostParamTemplate, 0)
	propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)
	_, rt := r.MQTTRTPropertyPost(DeviceTypeGW, r.GWParam, propertyPostParamMap, false)
	if rt == true {
		r.GWParam.ReportErrCnt = 0
		r.GWParam.ReportStatus = "onLine"
		setting.ZAPS.Debugf("上报服务[%s] 上报网关属性成功", r.GWParam.ServiceName)
	} else {
		setting.ZAPS.Debugf("上报服务[%s] 上报网关属性失败 失败计数%d/3", r.GWParam.ServiceName, r.GWParam.ReportErrCnt)
		r.GWParam.ReportErrCnt++
		if r.GWParam.ReportErrCnt >= 3 {
			r.GWParam.ReportErrCnt = 0

			r.GWParam.ReportStatus = "offLine"
			setting.ZAPS.Warnf("上报服务[%s] 网关离线", r.GWParam.ServiceName)
		}
	}
}

//指定设备上传属性
func (r *ReportServiceParamRTTemplate) NodePropertyPost(name []string, properties []MQTTRTPropertyPostParamPropertyTemplate, alarm bool) {

	propertyPostParamMap := make([]MQTTRTPropertyPostParamTemplate, 0)
	for _, n := range name {
		for k, v := range r.NodeList {
			if n == v.Name {
				//上报故障计数值先加，收到正确回应后清0
				r.NodeList[k].ReportErrCnt++
				onLine := false
				if v.CommStatus == "onLine" {
					onLine = true
				}

				model, ok := reportModel.ReportModels[v.UploadModel]
				if !ok {
					continue
				}
				propertyPostParam := MQTTRTPropertyPostParamTemplate{
					DeviceCode: v.Param.DeviceCode,
					OnLine:     onLine,
					TypeCode:   model.Code,
					VarData:    make([]MQTTRTPropertyPostParamPropertyTemplate, 0),
				}
				if alarm == true {
					//单个设备发送
					propertyPostParam.VarData = properties
					propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)
					r.MQTTRTPropertyPost(DeviceTypeNode, r.GWParam, propertyPostParamMap[len(propertyPostParamMap)-1:], alarm)
				} else {
					property := MQTTRTPropertyPostParamPropertyTemplate{}
					if v.CollInterfaceName != "virtual" {
						coll, ok := device.CollectInterfaceMap.Coll[v.CollInterfaceName]
						if !ok {
							continue
						}
						node, ok := coll.DeviceNodeMap[v.Name]
						if !ok {
							continue
						}

						for _, p := range v.Properties {
							for _, n := range node.Properties {
								if n.Name == p.Name {
									if len(n.Value) == 0 {
										continue
									}
									property.Name = p.UploadName
									property.TimeStamp = n.Value[len(n.Value)-1].TimeStamp.Unix()
									property.Value = n.Value[len(n.Value)-1].Value
									propertyPostParam.VarData = append(propertyPostParam.VarData, property)
								}
							}
						}
					} else {
						node, ok := virtual.VirtualDevice.Nodes[v.Name]
						if !ok {
							continue
						}
						for _, p := range v.Properties {
							for _, n := range node.Properties {
								if n.Name == p.Name {
									if len(n.Value) == 0 {
										continue
									}
									property.Name = p.UploadName
									property.TimeStamp = n.Value[len(n.Value)-1].TimeStamp.Unix()
									property.Value = n.Value[len(n.Value)-1].Value
									propertyPostParam.VarData = append(propertyPostParam.VarData, property)
								}
							}
						}
					}
					//单个设备发送
					propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)

					r.MQTTRTPropertyPost(DeviceTypeNode, r.GWParam, propertyPostParamMap[len(propertyPostParamMap)-1:], alarm)
				}
			}
		}
	}
}
