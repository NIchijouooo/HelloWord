package mqttM2M

import (
	"encoding/json"
	"fmt"
	"gateway/setting"
	"net"
	"time"
)

type MQTTM2MReportFrameTemplate struct {
	Topic   string
	Payload interface{}
}

type MQTTM2MReportDataTemplate struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type MQTTM2MReportYxTemplate struct {
	Time       string                      `json:"time"`
	CommStatus string                      `json:"commStatus,omitempty"`
	YxList     []MQTTM2MReportDataTemplate `json:"yxList"`
}

type MQTTM2MReportYcTemplate struct {
	Time       string                      `json:"time"`
	CommStatus string                      `json:"commStatus,omitempty"`
	YcList     []MQTTM2MReportDataTemplate `json:"ycList"`
}

type MQTTM2MReportSettingTemplate struct {
	Time        string                      `json:"time"`
	CommStatus  string                      `json:"commStatus,omitempty"`
	SettingList []MQTTM2MReportDataTemplate `json:"settingList"`
}
type MQTTM2MReportGPSTemplate struct {
	Time      string `json:"time"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}
type MQTTM2MReportPropertyTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
}

func (r *ReportServiceParamM2MTemplate) M2MPublishData(msg *MQTTM2MReportFrameTemplate) bool {
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

func (r *ReportServiceParamM2MTemplate) M2MPublishYxData(msg *MQTTM2MReportYxTemplate, id string) bool {
	status := false

	//propertyPostTopic := "iot/rx/" + r.GWParam.Param.AppKey + "/" + id + "/resultYx"
	propertyPostTopic := fmt.Sprintf(M2MMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "resultYx")
	sJson, _ := json.Marshal(msg)

	data := &MQTTM2MReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.M2MPublishData(data)

	return status
}

func (r *ReportServiceParamM2MTemplate) M2MPublishYcData(msg *MQTTM2MReportYcTemplate, id string) bool {
	status := false

	//propertyPostTopic := "iot/rx/" + r.GWParam.Param.AppKey + "/" + id + "/resultYc"
	propertyPostTopic := fmt.Sprintf(M2MMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "resultYc")
	sJson, _ := json.Marshal(msg)

	data := &MQTTM2MReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.M2MPublishData(data)

	return status
}

func (r *ReportServiceParamM2MTemplate) FM2MPublishSettingData(msg *MQTTM2MReportSettingTemplate, id string) bool {
	status := false

	//propertyPostTopic := "iot/rx/" + r.GWParam.Param.AppKey + "/" + id + "/setting"
	propertyPostTopic := fmt.Sprintf(M2MMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "setting")
	sJson, _ := json.Marshal(msg)

	data := &MQTTM2MReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.M2MPublishData(data)

	return status
}

func (r *ReportServiceParamM2MTemplate) M2MPublishLocationData(msg *MQTTM2MReportGPSTemplate, id string) bool {
	status := false

	propertyPostTopic := fmt.Sprintf(M2MMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "location")
	sJson, _ := json.Marshal(msg)

	data := &MQTTM2MReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	status = r.M2MPublishData(data)

	return status
}

//上传网关属性
var count uint32 = 0

func (r *ReportServiceParamM2MTemplate) GWPropertyPost() {

	count++

	ycPropertyMap := make([]MQTTM2MReportDataTemplate, 0)
	ycProperty := MQTTM2MReportDataTemplate{}

	ycProperty.ID = 1
	ycProperty.Name = "M2M上报计数"
	ycProperty.Value = count
	ycPropertyMap = append(ycPropertyMap, ycProperty)

	ycPropertyPostParam := MQTTM2MReportYcTemplate{
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		YcList: ycPropertyMap,
	}

	settingPropertyMap := make([]MQTTM2MReportDataTemplate, 0)
	settingProperty := MQTTM2MReportDataTemplate{}

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
	settingProperty.Value = str
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
	fmt.Println(params)
	settingProperty.ID = 3
	settingProperty.Name = "eth0 IP"
	settingProperty.Value = ip
	settingPropertyMap = append(settingPropertyMap, settingProperty)
	//设备的eth0网卡mask id =4
	fmt.Println(params)
	settingProperty.ID = 4
	settingProperty.Name = "eth0 mask"
	settingProperty.Value = mask
	settingPropertyMap = append(settingPropertyMap, settingProperty)
	//设备的eth0网卡gateway id =5
	fmt.Println(params)
	settingProperty.ID = 5
	settingProperty.Name = "eth0 gateway"
	settingProperty.Value = gateway
	settingPropertyMap = append(settingPropertyMap, settingProperty)
	//设备的eth0网卡dhcp id =6
	fmt.Println(params)
	settingProperty.ID = 6
	settingProperty.Name = "eth0 DHCP配置(不上报)"
	settingProperty.Value = "不上报"
	settingPropertyMap = append(settingPropertyMap, settingProperty)

	settingPropertyPostParam := MQTTM2MReportSettingTemplate{
		Time:        time.Now().Format("2006-01-02 15:04:05"),
		SettingList: settingPropertyMap,
	}

	if true == r.M2MPublishYcData(&ycPropertyPostParam, r.GWParam.Param.DeviceID) {
		r.GWParam.HeartBeatMark = true
	}

	if true == r.FM2MPublishSettingData(&settingPropertyPostParam, r.GWParam.Param.DeviceID) {
		r.GWParam.HeartBeatMark = true
	}
}
