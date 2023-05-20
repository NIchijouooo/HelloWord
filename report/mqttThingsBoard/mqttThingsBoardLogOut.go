package mqttThingsBoard

import (
	"encoding/json"
	"gateway/setting"
	"time"
)

type MQTTNodeLogoutTemplate struct {
	Device string `json:"device"`
}

func (r *ReportServiceParamThingsBoardTemplate) MQTTThingsBoardNodeLogOut(nodeMap []string) int {

	nodeLogout := MQTTNodeLogoutTemplate{}

	for _, v := range nodeMap {
		nodeLogout.Device = v

		sJson, _ := json.Marshal(nodeLogout)
		logoutTopic := "v1/gateway/disconnect"

		setting.ZAPS.Infof("上报服务[%s]发布节点离线线消息主题%s", r.GWParam.ServiceName, logoutTopic)
		setting.ZAPS.Debugf("上报服务[%s]发布节点离线消息内容%s", r.GWParam.ServiceName, sJson)

		if r.GWParam.MQTTClient != nil {
			token := r.GWParam.MQTTClient.Publish(logoutTopic, 1, false, sJson)
			if token.WaitTimeout(time.Duration(TimeOutLogout)*time.Second) == false {
				setting.ZAPS.Errorf("上报服务[%s]发布主题[%s]失败 %v", r.GWParam.ServiceName, logoutTopic, token.Error())
				continue
			}
			setting.ZAPS.Infof("上报服务[%s]发布主题[%s]成功", r.GWParam.ServiceName, logoutTopic)
		}
	}
	return 0
}

func (r *ReportServiceParamThingsBoardTemplate) NodeLogOut(name []string) bool {

	nodeMap := make([]string, 0)
	status := false

	setting.ZAPS.Debugf("上报服务[%s]节点%s离线", r.GWParam.ServiceName, name)
	for _, d := range name {
		for _, v := range r.NodeList {
			if d == v.Name {
				nodeMap = append(nodeMap, v.Param.DeviceName)
			}
		}
	}

	if len(nodeMap) > 0 {
		r.MQTTThingsBoardNodeLogOut(nodeMap)
	}

	return status
}
