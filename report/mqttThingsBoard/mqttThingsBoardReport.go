package mqttThingsBoard

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
	"time"
)

type MQTTThingsBoardReportPropertyTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
}

type MQTTThingsBoardPropertyPostParamPropertyTemplate struct {
}

type MQTTThingsBoardPropertyPostParamTemplate struct {
}

type MQTTThingsBoardNodePropertyPostTemplate struct {
	Values map[string]interface{} `json:"values"`
}

type MQTTThingsBoardGWPropertyPostTemplate struct {
}

func MQTTThingsBoardPropertyPost(gwParam ReportServiceGWParamThingsBoardTemplate, propertyParam []MQTTThingsBoardPropertyPostParamTemplate) (int, bool) {

	return 0, true
}

func (r *ReportServiceParamThingsBoardTemplate) GWPropertyPost() {

	properties := make(map[string]interface{})

	properties["memTotal"] = setting.SystemState.MemTotal

	properties["memUse"] = setting.SystemState.MemUse

	properties["diskTotal"] = setting.SystemState.DiskTotal

	properties["diskUse"] = setting.SystemState.DiskUse

	//properties["Name"] = setting.SystemState.Name

	properties["SN"] = setting.SystemState.SN

	//properties["HardVer"] = setting.SystemState.HardVer

	properties["softVer"] = setting.SystemState.SoftVer

	properties["systemRTC"] = setting.SystemState.SystemRTC

	properties["runTime"] = setting.SystemState.RunTime

	properties["deviceOnline"] = setting.SystemState.DeviceOnline

	properties["devicePacketLoss"] = setting.SystemState.DevicePacketLoss

	//清空接收缓存
	//for i := 0; i < len(r.ReceiveReportPropertyAckFrameChan); i++ {
	//	<-r.ReceiveReportPropertyAckFrameChan
	//}

	sJson, _ := json.Marshal(properties)
	propertyPostTopic := "v1/devices/me/telemetry"

	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 1, false, sJson)
		if token.WaitTimeout(time.Duration(TimeOutReportProperty)*time.Second) == false {
			setting.ZAPS.Errorf("上报服务[%s]上报网关属性失败 %v", r.GWParam.ServiceName, token.Error())
		} else {
			setting.ZAPS.Infof("上报服务[%s]上报网关属性成功 上报主题%v", r.GWParam.ServiceName, propertyPostTopic)
			setting.ZAPS.Infof("上报服务[%s]上报网关属性内容%v", r.GWParam.ServiceName, string(sJson))
		}
	}
}

//指定设备上传属性
func (r *ReportServiceParamThingsBoardTemplate) NodePropertyPost(name []string) {

	nodesProperties := make(map[string]interface{})
	for _, n := range name {
		for k, v := range r.NodeList {
			if n == v.Name {
				//上报故障计数值先加，收到正确回应后清0
				r.NodeList[k].ReportErrCnt++
				coll, ok := device.CollectInterfaceMap.Coll[v.CollInterfaceName]
				if !ok {
					continue
				}
				node, ok := coll.DeviceNodeMap[v.Name]
				if !ok {
					continue
				}

				nodes := make([]map[string]interface{}, 0)
				nodeValues := make(map[string]interface{})
				for _, p := range node.Properties {
					if len(p.Value) >= 1 {
						index := len(p.Value) - 1
						nodeValues[p.Name] = p.Value[index].Value
					}
				}
				nodes = append(nodes, nodeValues)
				nodesProperties[v.Param.DeviceName] = nodes
			}
		}
	}

	//清空接收缓存
	//for i := 0; i < len(r.ReceiveReportPropertyAckFrameChan); i++ {
	//	<-r.ReceiveReportPropertyAckFrameChan
	//}

	sJson, _ := json.Marshal(nodesProperties)
	propertyPostTopic := "v1/gateway/telemetry"

	setting.ZAPS.Infof("上报服务[%s]发布上报消息主题%s", r.GWParam.ServiceName, propertyPostTopic)
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 1, false, sJson)
		if token.WaitTimeout(time.Duration(TimeOutReportProperty)*time.Second) == false {
			setting.ZAPS.Errorf("上报服务[%s]上报节点属性失败 %v", r.GWParam.ServiceName, token.Error())
			setting.ZAPS.Debugf("上报服务[%s]上报节点属性内容%v", r.GWParam.ServiceName, string(sJson))
		} else {
			setting.ZAPS.Debugf("上报服务[%s]上报节点属性成功 上报内容%v", r.GWParam.ServiceName, string(sJson))
		}
	}
}
