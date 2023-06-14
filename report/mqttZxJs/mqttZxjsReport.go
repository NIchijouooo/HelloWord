package mqttZxJs

import (
	"encoding/json"
	"fmt"
	"gateway/device"
	"gateway/setting"
	"strconv"
	"time"
)

type MQTTZxjsReportFrameTemplate struct {
	Topic   string
	Payload interface{}
}

type MQTTZxjsReportDataTemplate struct {
	ID    int         `json:"id"`
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value"`
}

type MQTTZxjsReportYxTemplate struct {
	Time       string                       `json:"time"`
	CommStatus string                       `json:"commStatus,omitempty"`
	YxList     []MQTTZxjsReportDataTemplate `json:"yxList"`
}

type MQTTZxjsReportYcTemplate struct {
	Time       string                       `json:"time"`
	CommStatus string                       `json:"commStatus,omitempty"`
	YcList     []MQTTZxjsReportDataTemplate `json:"ycList"`
}

type MQTTZxjsReportSettingTemplate struct {
	Time        string                       `json:"time"`
	CommStatus  string                       `json:"commStatus,omitempty"`
	SettingList []MQTTZxjsReportDataTemplate `json:"settingList"`
}
type MQTTZxjsReportGPSTemplate struct {
	Time      string `json:"time"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}
type MQTTZxjsReportPropertyTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
}

func (r *ReportServiceParamZxjsTemplate) ZxjsPublishData(msg *MQTTZxjsReportFrameTemplate) bool {
	status := false

	if r.GWParam.MQTTClient != nil {
		if token := r.GWParam.MQTTClient.Publish(msg.Topic, 0, false, msg.Payload); token.WaitTimeout(2000*time.Millisecond) && token.Error() != nil {
			status = false
			setting.ZAPS.Debugf("上报服务[%s]发布消息失败 %v", r.GWParam.ServiceName, token.Error())
		} else {
			status = true
			setting.ZAPS.Debugf("上报服务[%s]发布消息成功 内容%s", r.GWParam.ServiceName, msg.Payload)
		}
	}

	return status
}

func (r *ReportServiceParamZxjsTemplate) ZxjsPublishYcData(msg *MQTTZxjsReportYcTemplate, id string) bool {
	status := false

	//propertyPostTopic := "iot/rx/" + r.GWParam.Param.AppKey + "/" + id + "/resultYc"
	propertyPostTopic := fmt.Sprintf(ZxjsMQTTTopicRxFormat, r.GWParam.Param.ProductSn, id, "resultYc")
	sJson, _ := json.Marshal(msg)

	data := &MQTTZxjsReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.ZxjsPublishData(data)

	return status
}

func (r *ReportServiceParamZxjsTemplate) ZxjsPublishdeviceControlResult(sJson []byte, id string) bool {
	status := false

	//propertyPostTopic := "iot/rx/" + r.GWParam.Param.AppKey + "/" + id + "/resultYc"
	propertyPostTopic := fmt.Sprintf(ZxjsMQTTTopicRxFormat, r.GWParam.Param.ProductSn, id, "/upload/cack")
	//sJson, _ := json.Marshal(msg)

	data := &MQTTZxjsReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.ZxjsPublishData(data)

	return status
}

func (r *ReportServiceParamZxjsTemplate) ZxjsPublishSettingData(msg *MQTTZxjsReportSettingTemplate, id string) bool {
	status := false

	//propertyPostTopic := "iot/rx/" + r.GWParam.Param.AppKey + "/" + id + "/setting"
	propertyPostTopic := fmt.Sprintf(ZxjsMQTTTopicRxFormat, r.GWParam.Param.ProductSn, id, "setting")
	sJson, _ := json.Marshal(msg)

	data := &MQTTZxjsReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.ZxjsPublishData(data)

	return status
}

func (r *ReportServiceParamZxjsTemplate) ZxjsPublishLocationData(msg *MQTTZxjsReportGPSTemplate, id string) bool {
	status := false

	propertyPostTopic := fmt.Sprintf(ZxjsMQTTTopicRxFormat, r.GWParam.Param.ProductSn, id, "location")
	sJson, _ := json.Marshal(msg)

	data := &MQTTZxjsReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.ZxjsPublishData(data)

	return status
}

// 上传网关属性
var count uint32 = 0

// 指定设备上传属性，可以是多个设备
func (r *ReportServiceParamZxjsTemplate) NodePropertyPost(name []string) {

	for _, n := range name {
		for k, v := range r.NodeList {
			if n == v.Name {
				r.NodeList[k].ReportErrCnt++ //上报故障计数值先加，收到正确回应后清0

				ycPropertyMap := make([]MQTTZxjsReportDataTemplate, 0)
				ycPropertyPostParam := MQTTZxjsReportYcTemplate{
					Time:       time.Now().Format("2006-01-02 15:04:05"),
					CommStatus: "onLink",
				}

				//获取采集接口数据
				coll, collErr := device.CollectInterfaceMap.Coll[v.CollInterfaceName]
				if !collErr {
					ycPropertyPostParam.CommStatus = fmt.Sprintf("coll接口[%s]不存在", v.CollInterfaceName)
					if true == r.ZxjsPublishYcData(&ycPropertyPostParam, v.Param.DeviceSn) {
						r.NodeList[k].HeartBeatMark = true
						r.NodeList[k].ReportErrCnt = 0
					}
					continue
				}

				//获取节点数据
				node, nodeErr := coll.DeviceNodeMap[v.Name]
				if !nodeErr {
					ycPropertyPostParam.CommStatus = fmt.Sprintf("coll接口[%s]下的设备[%s]不存在", v.CollInterfaceName, v.Name)
					if true == r.ZxjsPublishYcData(&ycPropertyPostParam, v.Param.DeviceSn) {
						r.NodeList[k].HeartBeatMark = true
						r.NodeList[k].ReportErrCnt = 0
					}
					continue
				}

				//获取节点属性数据
				if node.CommStatus == "offLine" {
					ycPropertyPostParam.CommStatus = "offLine"
					r.NodeList[k].CommStatus = "offLine"
				} else {
					r.NodeList[k].CommStatus = "onLink"
					for _, v := range node.Properties {
						if num, err := strconv.Atoi(v.Name); err == nil {
							ycProperty := MQTTZxjsReportDataTemplate{
								ID: num,
								// QJH Delect 2023/6/6 去掉name上报，解决数据包上报流量过大
								//Name: v.Label,
							}
							if len(v.Value) >= 1 { //当前属性有数据
								ycProperty.Value = v.Value[len(v.Value)-1].Value
							} else {
								ycProperty.Value = "当前属性无数据"
							}

							ycPropertyMap = append(ycPropertyMap, ycProperty)
						}
					}

					ycPropertyPostParam.YcList = ycPropertyMap
				}

				if true == r.ZxjsPublishYcData(&ycPropertyPostParam, v.Param.DeviceSn) {
					r.NodeList[k].HeartBeatMark = true
					r.NodeList[k].ReportErrCnt = 0
				}
			}
		}
	}
}
