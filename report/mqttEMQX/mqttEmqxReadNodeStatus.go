package mqttEmqx

import (
	"encoding/json"
	"gateway/setting"
	"time"
)

type MQTTEmqxReadNodeStatusRequestParamTemplate struct {
}

type MQTTEmqxReadNodeStatusRequestTemplate struct {
	ID      string                                     `json:"id"`
	Version string                                     `json:"version"`
	Ack     int                                        `json:"ack"`
	Params  MQTTEmqxReadNodeStatusRequestParamTemplate `json:"params"`
}

type MQTTEmqxReadNodeStatusAckParamTemplate struct {
	ClientID   string `json:"clientID"`
	CommStatus string `json:"commStatus"`
}

type MQTTEmqxReadNodeStatusAckTemplate struct {
	ID      string                                   `json:"id"`
	Version string                                   `json:"version"`
	Code    int                                      `json:"code"`
	Params  []MQTTEmqxReadNodeStatusAckParamTemplate `json:"params"`
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxReadNodeStatusAck(reqFrame MQTTEmqxReadNodeStatusRequestTemplate, code int, ackParams []MQTTEmqxReadNodeStatusAckParamTemplate) {

	ackFrame := MQTTEmqxReadNodeStatusAckTemplate{
		ID:      reqFrame.ID,
		Version: reqFrame.Version,
		Code:    code,
		Params:  ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	nodeStatusGetReplyTopic := "/sys/thing/node/status/get_reply/" + r.GWParam.Param.ClientID

	setting.ZAPS.Infof("nodeStatus get_reply topic: %s", nodeStatusGetReplyTopic)
	setting.ZAPS.Debugf("nodeStatus get_reply: %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(nodeStatusGetReplyTopic, 1, false, sJson)
		if token.WaitTimeout(5*time.Second) == false {
			setting.ZAPS.Errorf("EMQX上报服务[%s]发布回复读取设备状态消息失败 %v", r.GWParam.ServiceName, token.Error())
			return
		}
	} else {
		setting.ZAPS.Errorf("EMQX上报服务[%s]发布回复读取设备状态消息失败", r.GWParam.ServiceName)
		return
	}
	setting.ZAPS.Infof("EMQX上报服务[%s]发布读取设备状态消息成功", r.GWParam.ServiceName)
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxProcessReadNodeStatus(reqFrame MQTTEmqxReadNodeStatusRequestTemplate) {

	ackParams := make([]MQTTEmqxReadNodeStatusAckParamTemplate, 0)
	nodeStatus := MQTTEmqxReadNodeStatusAckParamTemplate{}
	for _, node := range r.NodeList {
		nodeStatus.ClientID = node.Param.DeviceCode
		nodeStatus.CommStatus = node.CommStatus
		ackParams = append(ackParams, nodeStatus)
	}

	r.ReportServiceEmqxReadNodeStatusAck(reqFrame, 0, ackParams)

}
