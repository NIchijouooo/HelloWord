package mqttFeisjy

import (
	"encoding/json"
	"fmt"
	"gateway/device"
	"gateway/setting"
	"net"
	"strconv"
	"time"
)

type MQTTFeisjyReportFrameTemplate struct {
	Topic   string
	Payload interface{}
}

type MQTTFeisjyReportDataTemplate struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type MQTTFeisjyReportYxTemplate struct {
	Time       string                         `json:"time"`
	CommStatus string                         `json:"commStatus,omitempty"`
	YxList     []MQTTFeisjyReportDataTemplate `json:"yxList"`
}

type MQTTFeisjyReportYcTemplate struct {
	Time       string                         `json:"time"`
	CommStatus string                         `json:"commStatus,omitempty"`
	YcList     []MQTTFeisjyReportDataTemplate `json:"ycList"`
}

type MQTTFeisjyReportSettingTemplate struct {
	Time        string                         `json:"time"`
	CommStatus  string                         `json:"commStatus,omitempty"`
	SettingList []MQTTFeisjyReportDataTemplate `json:"settingList"`
}
type MQTTFeisjyReportGPSTemplate struct {
	Time      string `json:"time"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}
type MQTTFeisjyReportPropertyTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
}

func (r *ReportServiceParamFeisjyTemplate) FeisjyPublishData(msg *MQTTFeisjyReportFrameTemplate) bool {
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

func (r *ReportServiceParamFeisjyTemplate) FeisjyPublishYxData(msg *MQTTFeisjyReportYxTemplate, id string) bool {
	status := false

	//propertyPostTopic := "iot/rx/" + r.GWParam.Param.AppKey + "/" + id + "/resultYx"
	propertyPostTopic := fmt.Sprintf(FeisjyMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "resultYx")
	sJson, _ := json.Marshal(msg)

	data := &MQTTFeisjyReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.FeisjyPublishData(data)

	return status
}

func (r *ReportServiceParamFeisjyTemplate) FeisjyPublishYcData(msg *MQTTFeisjyReportYcTemplate, id string) bool {
	status := false

	//propertyPostTopic := "iot/rx/" + r.GWParam.Param.AppKey + "/" + id + "/resultYc"
	propertyPostTopic := fmt.Sprintf(FeisjyMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "resultYc")
	sJson, _ := json.Marshal(msg)

	data := &MQTTFeisjyReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.FeisjyPublishData(data)

	return status
}

func (r *ReportServiceParamFeisjyTemplate) FeisjyPublishSettingData(msg *MQTTFeisjyReportSettingTemplate, id string) bool {
	status := false

	//propertyPostTopic := "iot/rx/" + r.GWParam.Param.AppKey + "/" + id + "/setting"
	propertyPostTopic := fmt.Sprintf(FeisjyMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "setting")
	sJson, _ := json.Marshal(msg)

	data := &MQTTFeisjyReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.FeisjyPublishData(data)

	return status
}

func (r *ReportServiceParamFeisjyTemplate) FeisjyPublishLocationData(msg *MQTTFeisjyReportGPSTemplate, id string) bool {
	status := false

	propertyPostTopic := fmt.Sprintf(FeisjyMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "location")
	sJson, _ := json.Marshal(msg)

	data := &MQTTFeisjyReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.FeisjyPublishData(data)

	return status
}

//上传网关属性
var count uint32 = 0

func (r *ReportServiceParamFeisjyTemplate) GWPropertyPost() {
	//YC
	ycPropertyMap := make([]MQTTFeisjyReportDataTemplate, 0)
	ycProperty := MQTTFeisjyReportDataTemplate{}

	ycProperty.ID = 1
	ycProperty.Name = "M2M上报计数"
	ycProperty.Value = count
	ycPropertyMap = append(ycPropertyMap, ycProperty)

	ycPropertyPostParam := MQTTFeisjyReportYcTemplate{
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		YcList: ycPropertyMap,
	}

	//SETTING
	settingPropertyMap := make([]MQTTFeisjyReportDataTemplate, 0)
	settingProperty := MQTTFeisjyReportDataTemplate{}

	//当前设备时间  id =1
	settingProperty.ID = 1
	settingProperty.Name = "当前时间"
	settingProperty.Value = time.Now().Format("2006-01-02 15:04:05")
	settingPropertyMap = append(settingPropertyMap, settingProperty)

	//设备的网卡名称 id =2
	data := setting.GetNetworkNames()
	str := ""
	for _, v := range data {
		str += v
		str += " & "
	}
	settingProperty.ID = 2
	settingProperty.Name = "网卡名称"
	settingProperty.Value = str // "eth0,usb0"
	settingPropertyMap = append(settingPropertyMap, settingProperty)

	inters, _ := net.Interfaces()
	var params net.Interface
	for _, v := range inters {
		if "eth0" == v.Name {
			params = v
		}
	}
	ip, mask := setting.GetIPAndMask(params.Name)
	gateway := setting.GetGateway(params.Name)

	//设备的eth0网卡ip id =3
	settingProperty.ID = 3
	settingProperty.Name = "eth0 IP"
	settingProperty.Value = ip
	settingPropertyMap = append(settingPropertyMap, settingProperty)
	//设备的eth0网卡mask id =4
	settingProperty.ID = 4
	settingProperty.Name = "eth0 mask"
	settingProperty.Value = mask
	settingPropertyMap = append(settingPropertyMap, settingProperty)
	//设备的eth0网卡gateway id =5
	settingProperty.ID = 5
	settingProperty.Name = "eth0 gateway"
	settingProperty.Value = gateway
	settingPropertyMap = append(settingPropertyMap, settingProperty)

	settingPropertyPostParam := MQTTFeisjyReportSettingTemplate{
		Time:        time.Now().Format("2006-01-02 15:04:05"),
		SettingList: settingPropertyMap,
	}

	if true == r.FeisjyPublishYcData(&ycPropertyPostParam, r.GWParam.Param.DeviceID) {
		r.GWParam.HeartBeatMark = true
	}

	if true == r.FeisjyPublishSettingData(&settingPropertyPostParam, r.GWParam.Param.DeviceID) {
		r.GWParam.HeartBeatMark = true
	}
}

//指定设备上传属性，可以是多个设备
func (r *ReportServiceParamFeisjyTemplate) NodePropertyPost(name []string) {

	for _, n := range name {
		for k, v := range r.NodeList {
			if n == v.Name {
				r.NodeList[k].ReportErrCnt++ //上报故障计数值先加，收到正确回应后清0

				ycPropertyMap := make([]MQTTFeisjyReportDataTemplate, 0)
				ycPropertyPostParam := MQTTFeisjyReportYcTemplate{
					Time:       time.Now().Format("2006-01-02 15:04:05"),
					CommStatus: "onLink",
				}

				//获取采集接口数据
				coll, collErr := device.CollectInterfaceMap.Coll[v.CollInterfaceName]
				if !collErr {
					ycPropertyPostParam.CommStatus = fmt.Sprintf("coll接口[%s]不存在", v.CollInterfaceName)
					if true == r.FeisjyPublishYcData(&ycPropertyPostParam, v.Param.DeviceID) {
						r.NodeList[k].HeartBeatMark = true
						r.NodeList[k].ReportErrCnt = 0
					}
					continue
				}

				//获取节点数据
				node, nodeErr := coll.DeviceNodeMap[v.Name]
				if !nodeErr {
					ycPropertyPostParam.CommStatus = fmt.Sprintf("coll接口[%s]下的设备[%s]不存在", v.CollInterfaceName, v.Name)
					if true == r.FeisjyPublishYcData(&ycPropertyPostParam, v.Param.DeviceID) {
						r.NodeList[k].HeartBeatMark = true
						r.NodeList[k].ReportErrCnt = 0
					}
					continue
				}

				//获取节点属性数据
				if node.CommStatus == "offLine" {
					ycPropertyPostParam.CommStatus = "offLine"
				} else {
					for _, v := range node.Properties {
						if num, err := strconv.Atoi(v.Name); err == nil {
							ycProperty := MQTTFeisjyReportDataTemplate{
								ID:   num,
								Name: v.Label,
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

				if true == r.FeisjyPublishYcData(&ycPropertyPostParam, v.Param.DeviceID) {
					r.NodeList[k].HeartBeatMark = true
					r.NodeList[k].ReportErrCnt = 0
				}
			}
		}
	}
}
