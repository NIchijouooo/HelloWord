package mqttZxJs

import (
	"encoding/json"
	"fmt"
	"gateway/device"
	"gateway/report/mqttFeisjy"
	"gateway/setting"
	"time"
)

type MQTTZxjsReportFrameTemplate struct {
	Topic   string
	Payload interface{}
}

type MQTTZxjsReportPropertyTemplate struct {
	Seq      int //设备类型，"gw" "node"
	DeviceSN string
}

type MQTTZxjsReportValueTemplate struct {
	ID    string      `json:"id"`
	Value interface{} `json:"value"`
}

type MQTTZxjsReportDataTemplate struct {
	DeviceSN string `json:"deviceSN"`
	Seq      int    `json:"seq"`
	Ts       int64  `json:"ts"`
	Type     string `json:"type"`
}

type CACKTemplate struct {
	Seq int `json:"seq"`
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

func (r *ReportServiceParamZxjsTemplate) ZxjsPublishYcData(msg *MQTTZxjsReportDataTemplate, ycList []MQTTZxjsReportValueTemplate, id string) bool {
	ycLength := len(ycList)
	setting.ZAPS.Infof("zxjs report data length = [%v]", ycLength)
	if ycLength == 0 {
		return false
	}
	// 每包上送数量
	index := 1000
	// 模与余数
	remainder := ycLength % index
	// 需要上送次数
	count := 0
	if remainder == 0 {
		count = ycLength / index
	} else {
		count = ycLength/index + 1
	}
	// 防止一包数据过大,分包上送
	for i := 0; i < count; i++ {
		start := i * index
		end := start + index
		if i == count-1 {
			end = ycLength
		}
		list := ycList[start:end]
		sendData := make(map[string]interface{}, index)
		sendData["deviceSN"] = msg.DeviceSN
		sendData["seq"] = msg.Seq
		sendData["ts"] = msg.Ts
		sendData["type"] = msg.Type
		for _, v := range list {
			sendData[v.ID] = v.Value
		}
		//propertyPostTopic := "/" + ProductSn + "/" + id + "/upload"
		propertyPostTopic := fmt.Sprintf(ZxjsMQTTTopicRxFormat, r.GWParam.Param.ProductSn, id, "upload")
		sJson, _ := json.Marshal(sendData)

		data := &MQTTZxjsReportFrameTemplate{
			Topic:   propertyPostTopic,
			Payload: sJson,
		}

		publish := r.ZxjsPublishData(data)
		setting.ZAPS.Infof("zxjs report data index[%v] size[%v] result[%v]", i, len(list), publish)
	}

	return true
}

// 上传网关属性
var count uint32 = 0

// 指定设备上传属性，可以是多个设备
func (r *ReportServiceParamZxjsTemplate) NodePropertyPost(property MQTTZxjsReportPropertyTemplate) {
	seq := property.Seq
	if seq != 0 {
		r.ZxjsPublishCallCack(seq)
	}

	ycPropertyPostParam := MQTTZxjsReportDataTemplate{
		DeviceSN: property.DeviceSN,
		Seq:      seq,
		Ts:       time.Now().Unix(),
		Type:     "rtg",
	}

	ycList := make([]MQTTZxjsReportValueTemplate, 0)

	// 将iot的数据上报到第三方平台
	for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
		for _, d := range v.NodeList {
			initYcList(&ycList, d.CollInterfaceName, d.Name)
		}
	}
	r.ZxjsPublishYcData(&ycPropertyPostParam, ycList, property.DeviceSN)
}

/*
*
初始化要上报的点位集合
*/
func initYcList(ycList *[]MQTTZxjsReportValueTemplate, collInterfaceName string, nodeName string) {
	//获取采集接口数据
	coll, collErr := device.CollectInterfaceMap.Coll[collInterfaceName]
	if !collErr {
		setting.ZAPS.Infof("coll接口[%s]不存在", collInterfaceName)
		return
	}

	//获取节点数据
	node, nodeErr := coll.DeviceNodeMap[nodeName]
	if !nodeErr {
		setting.ZAPS.Infof("coll接口[%s]下的设备[%s]不存在", collInterfaceName, nodeName)
		return
	}
	if len(node.Properties) == 0 || node.CommStatus == "offLine" {
		setting.ZAPS.Infof("节点数据为空或通信状态离线,coll接口[%s],设备[%s]", collInterfaceName, nodeName)
		return
	}
	// 判断节点唯一标识和数据不为空则上报
	for _, v := range node.Properties {
		if len(v.Identity) == 0 || len(v.Value) == 0 {
			continue
		}
		val := MQTTZxjsReportValueTemplate{
			ID:    v.Identity,
			Value: v.Value[len(v.Value)-1].Value,
		}
		*ycList = append(*ycList, val)
		////本地测试
		//if len(v.Identity) == 0 {
		//	continue
		//}
		//val := MQTTZxjsReportValueTemplate{
		//	ID: v.Identity,
		//}
		//if len(v.Value) == 0 {
		//	val.Value = 0
		//	*ycList = append(*ycList, val)
		//} else {
		//	val.Value = v.Value[len(v.Value)-1].Value
		//	*ycList = append(*ycList, val)
		//}
	}
}

/*
*召读响应
 */
func (r *ReportServiceParamZxjsTemplate) ZxjsPublishCallCack(seq int) {
	msg := CACKTemplate{
		Seq: seq,
	}
	//propertyPostTopic := "/" + r.GWParam.Param.ProductSn + "/" + r.GWParam.Param.DeviceSn + "/upload/cack"
	propertyPostTopic := fmt.Sprintf(ZxjsMQTTTopicRxFormat, r.GWParam.Param.ProductSn, r.GWParam.Param.DeviceSn, "upload/cack")
	sJson, _ := json.Marshal(msg)

	data := &MQTTZxjsReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	r.ZxjsPublishData(data)
}

/*
*控制响应
 */
func (r *ReportServiceParamZxjsTemplate) ZxjsPublishSetCack(ctrlInfo MQTTZxjsControlTemplate, valid int) {
	ctrlInfo.Data.Valid = valid
	//propertyPostTopic := "/" + r.GWParam.Param.ProductSn + "/" + r.GWParam.Param.DeviceSn + "/set/cack"
	sJson, _ := json.Marshal(ctrlInfo)
	propertyPostTopic := fmt.Sprintf(ZxjsMQTTTopicRxFormat, r.GWParam.Param.ProductSn, r.GWParam.Param.DeviceSn, "set/cack")
	data := &MQTTZxjsReportFrameTemplate{
		Topic:   propertyPostTopic,
		Payload: sJson,
	}

	r.ZxjsPublishData(data)
}
