package mqttEmqx

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
	"strconv"
	"time"
)

type MQTTEmqxReportPropertyTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
}

type MQTTEmqxReportAlarmTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
	Properties []MQTTEmqxPropertyPostParamPropertyTemplate
}

type MQTTEmqxPropertyPostParamPropertyTemplate struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	TimeStamp int64       `json:"timestamp"`
}

type MQTTEmqxPropertyPostParamTemplate struct {
	ClientID   string                                      `json:"clientID"`
	Properties []MQTTEmqxPropertyPostParamPropertyTemplate `json:"properties"`
}

type MQTTEmqxPropertyPostTemplate struct {
	ID      string                              `json:"id"`
	Version string                              `json:"version"`
	Ack     int                                 `json:"ack"`
	Params  []MQTTEmqxPropertyPostParamTemplate `json:"params"`
}

type MQTTEmqxReportPropertyAckTemplate struct {
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

func MQTTEmqxPropertyPost(rType int, gwParam ReportServiceGWParamEmqxTemplate, propertyParam []MQTTEmqxPropertyPostParamTemplate, alarm bool) (int, bool) {

	propertyPost := MQTTEmqxPropertyPostTemplate{
		ID:      strconv.Itoa(MsgID),
		Version: "V1.0",
		Ack:     1,
		Params:  propertyParam,
	}
	MsgID++

	sJson, _ := json.Marshal(propertyPost)
	propertyPostTopic := ""
	if rType == DeviceTypeGW {
		propertyPostTopic = "/sys/thing/gw/property/post/" + gwParam.Param.ClientID
	} else {
		if alarm == false {
			propertyPostTopic = "/sys/thing/node/property/post/" + gwParam.Param.ClientID
		} else {
			propertyPostTopic = "/sys/thing/node/alarm/post/" + gwParam.Param.ClientID
		}
	}

	setting.ZAPS.Infof("上报服务[%s]发布上报消息主题%s", gwParam.ServiceName, propertyPostTopic)
	if gwParam.MQTTClient != nil {
		token := gwParam.MQTTClient.Publish(propertyPostTopic, 1, false, sJson)
		if token.WaitTimeout(5*time.Second) == false {
			if alarm == false {
				setting.ZAPS.Errorf("Emqx上报服务[%s]发布上报属性消息失败 %v", gwParam.ServiceName, token.Error())
			} else {
				setting.ZAPS.Errorf("Emqx上报服务[%s]发布上报告警消息失败 %v", gwParam.ServiceName, token.Error())
			}
			return MsgID, false
		}
	} else {
		if alarm == false {
			setting.ZAPS.Errorf("Emqx上报服务[%s]发布上报属性消息失败", gwParam.ServiceName)
		} else {
			setting.ZAPS.Errorf("Emqx上报服务[%s]发布上报告警消息失败", gwParam.ServiceName)
		}
		return MsgID, false
	}
	if alarm == false {
		setting.ZAPS.Debugf("Emqx上报服务[%s]发布上报属性消息成功", gwParam.ServiceName)
	} else {
		setting.ZAPS.Debugf("Emqx上报服务[%s]发布上报告警消息成功", gwParam.ServiceName)
	}
	return MsgID, true
}

func (r *ReportServiceParamEmqxTemplate) GWPropertyPost() {

	propertyMap := make([]MQTTEmqxPropertyPostParamPropertyTemplate, 0)

	property := MQTTEmqxPropertyPostParamPropertyTemplate{}

	timeStamp := time.Now().Unix()

	property.Name = "MemTotal"
	property.Value = setting.SystemState.MemTotal
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "MemUse"
	property.Value = setting.SystemState.MemUse
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "DiskTotal"
	property.Value = setting.SystemState.DiskTotal
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "DiskUse"
	property.Value = setting.SystemState.DiskUse
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "Name"
	property.Value = setting.SystemState.Name
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "SN"
	property.Value = setting.SystemState.SN
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "HardVer"
	property.Value = setting.SystemState.HardVer
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "SoftVer"
	property.Value = setting.SystemState.SoftVer
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "SystemRTC"
	property.Value = setting.SystemState.SystemRTC
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "RunTime"
	property.Value = setting.SystemState.RunTime
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "DeviceOnline"
	property.Value = setting.SystemState.DeviceOnline
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "DevicePacketLoss"
	property.Value = setting.SystemState.DevicePacketLoss
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)
	//清空接收缓存
	for i := 0; i < len(r.ReceiveReportPropertyAckFrameChan); i++ {
		<-r.ReceiveReportPropertyAckFrameChan
	}

	propertyPostParam := MQTTEmqxPropertyPostParamTemplate{
		ClientID:   r.GWParam.Param.ClientID,
		Properties: propertyMap,
	}

	propertyPostParamMap := make([]MQTTEmqxPropertyPostParamTemplate, 0)
	propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)
	_, rt := MQTTEmqxPropertyPost(DeviceTypeGW, r.GWParam, propertyPostParamMap, false)
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

			token := r.GWParam.MQTTClient.Connect()
			select {
			case <-token.Done():
				{
					if token.Error() != nil {
						setting.ZAPS.Warnf("上报服务[%s] 网关 重新登录失败", r.GWParam.ServiceName)
					}
				}
			case <-time.After(time.Duration(EMQXTimeOutLogin) * time.Second):
				{
					r.GWParam.ReportStatus = "onLine"
					setting.ZAPS.Debugf("上报服务[%s] 网关 重新登录成功", r.GWParam.ServiceName)
				}
			}
			//if r.GWLogin() == true {
			//	r.GWParam.ReportStatus = "onLine"
			//	setting.ZAPS.Debugf("上报服务[%s] 网关 重新登录成功", r.GWParam.ServiceName)
			//} else {
			//	setting.ZAPS.Warnf("上报服务[%s] 网关 重新登录失败", r.GWParam.ServiceName)
			//}
		}
	}
}

//指定设备上传属性
//指定设备上传属性
func (r *ReportServiceParamEmqxTemplate) NodePropertyPost(name []string, properties []MQTTEmqxPropertyPostParamPropertyTemplate, alarm bool) {

	propertyPostParamMap := make([]MQTTEmqxPropertyPostParamTemplate, 0)
	for _, n := range name {
		for k, v := range r.NodeList {
			if n == v.Name {
				//上报故障计数值先加，收到正确回应后清0
				r.NodeList[k].ReportErrCnt++
				propertyPostParam := MQTTEmqxPropertyPostParamTemplate{
					ClientID: v.Param.DeviceCode,
				}
				if alarm == true {
					//单个设备发送
					propertyPostParam.Properties = properties
					propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)
					MQTTEmqxPropertyPost(DeviceTypeNode, r.GWParam, propertyPostParamMap[len(propertyPostParamMap)-1:], alarm)
				} else {
					coll, ok := device.CollectInterfaceMap.Coll[v.CollInterfaceName]
					if !ok {
						continue
					}
					node, ok := coll.DeviceNodeMap[v.Name]
					if !ok {
						continue
					}

					property := MQTTEmqxPropertyPostParamPropertyTemplate{}
					for _, p := range v.Properties {
						for _, n := range node.Properties {
							if n.Name == p.Name {
								if len(n.Value) == 0 {
									continue
								}
								property.Name = p.Name
								property.TimeStamp = n.Value[len(n.Value)-1].TimeStamp.Unix()
								property.Value = n.Value[len(n.Value)-1].Value
								propertyPostParam.Properties = append(propertyPostParam.Properties, property)
							}
						}
					}

					//单个设备发送
					propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)
					MQTTEmqxPropertyPost(DeviceTypeNode, r.GWParam, propertyPostParamMap[len(propertyPostParamMap)-1:], alarm)
				}
			}
		}
	}
}
