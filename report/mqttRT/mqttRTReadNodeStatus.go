package mqttRT

import (
	"encoding/json"
	"gateway/setting"
	"time"
)

type MQTTRTReadNodeStatusRequestParamTemplate struct {
}

type MQTTRTReadNodeStatusRequestTemplate struct {
	ID      string                                   `json:"id"`
	Version string                                   `json:"version"`
	Ack     int                                      `json:"ack"`
	Params  MQTTRTReadNodeStatusRequestParamTemplate `json:"params"`
}

type MQTTRTReadNodeStatusAckParamTemplate struct {
	ClientID   string `json:"clientID"`
	CommStatus string `json:"commStatus"`
}

type MQTTRTReadNodeStatusAckTemplate struct {
	ID      string                                 `json:"id"`
	Version string                                 `json:"version"`
	Code    int                                    `json:"code"`
	Params  []MQTTRTReadNodeStatusAckParamTemplate `json:"params"`
}

func (r *ReportServiceParamRTTemplate) ReportServiceRTReadNodeStatusAck(reqFrame MQTTRTReadNodeStatusRequestTemplate, code int, ackParams []MQTTRTReadNodeStatusAckParamTemplate) {

	ackFrame := MQTTRTReadNodeStatusAckTemplate{
		ID:      reqFrame.ID,
		Version: reqFrame.Version,
		Code:    code,
		Params:  ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	nodeStatusGetReplyTopic := "/sys/thing/node/status/get_reply/" + r.GWParam.Param.ClientID

	setting.ZAPS.Infof("RT上报服务发布读取设备状态应答消息主题 %s", nodeStatusGetReplyTopic)
	setting.ZAPS.Debugf("RT上报服务发布读取设备状态应答消息内容 %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(nodeStatusGetReplyTopic, 1, false, sJson)
		if token.WaitTimeout(time.Duration(RTTimeOutReadNode)*time.Second) == false {
			setting.ZAPS.Errorf("RT上报服务[%s]发布回复读取设备状态消息失败 %v", r.GWParam.ServiceName, token.Error())
			return
		}
	} else {
		setting.ZAPS.Errorf("RT上报服务[%s]发布回复读取设备状态消息失败", r.GWParam.ServiceName)
		return
	}
	setting.ZAPS.Infof("RT上报服务[%s]发布读取设备状态消息成功", r.GWParam.ServiceName)
}

func (r *ReportServiceParamRTTemplate) ReportServiceRTProcessReadNodeStatus(reqFrame MQTTRTReadNodeStatusRequestTemplate) {

	ackParams := make([]MQTTRTReadNodeStatusAckParamTemplate, 0)
	nodeStatus := MQTTRTReadNodeStatusAckParamTemplate{}
	for _, node := range r.NodeList {
		nodeStatus.ClientID = node.Param.DeviceCode
		nodeStatus.CommStatus = node.CommStatus
		ackParams = append(ackParams, nodeStatus)
	}

	r.ReportServiceRTReadNodeStatusAck(reqFrame, 0, ackParams)

}
