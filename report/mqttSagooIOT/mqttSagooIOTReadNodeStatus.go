package mqttSagooIOT

import (
	"encoding/json"
	"gateway/setting"
	"time"
)

type MQTTSagooIOTReadNodeStatusRequestParamTemplate struct {
}

type MQTTSagooIOTReadNodeStatusRequestTemplate struct {
	ID      string                                         `json:"id"`
	Version string                                         `json:"version"`
	Ack     int                                            `json:"ack"`
	Params  MQTTSagooIOTReadNodeStatusRequestParamTemplate `json:"params"`
}

type MQTTSagooIOTReadNodeStatusAckParamTemplate struct {
	ClientID   string `json:"clientID"`
	CommStatus string `json:"commStatus"`
}

type MQTTSagooIOTReadNodeStatusAckTemplate struct {
	ID      string                                       `json:"id"`
	Version string                                       `json:"version"`
	Code    int                                          `json:"code"`
	Params  []MQTTSagooIOTReadNodeStatusAckParamTemplate `json:"params"`
}

func (r *ReportServiceParamSagooIOTTemplate) ReportServiceSagooIOTReadNodeStatusAck(reqFrame MQTTSagooIOTReadNodeStatusRequestTemplate, code int, ackParams []MQTTSagooIOTReadNodeStatusAckParamTemplate) {

	ackFrame := MQTTSagooIOTReadNodeStatusAckTemplate{
		ID:      reqFrame.ID,
		Version: reqFrame.Version,
		Code:    code,
		Params:  ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	nodeStatusGetReplyTopic := "/sys/thing/node/status/get_reply/" + r.GWParam.Param.ClientID

	setting.ZAPS.Infof("SagooIOT上报服务发布读取设备状态应答消息主题 %s", nodeStatusGetReplyTopic)
	setting.ZAPS.Debugf("SagooIOT上报服务发布读取设备状态应答消息内容 %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(nodeStatusGetReplyTopic, 1, false, sJson)
		if token.WaitTimeout(time.Duration(SagooIOTTimeOutReadNode)*time.Second) == false {
			setting.ZAPS.Errorf("SagooIOT上报服务[%s]发布回复读取设备状态消息失败 %v", r.GWParam.ServiceName, token.Error())
			return
		}
	} else {
		setting.ZAPS.Errorf("SagooIOT上报服务[%s]发布回复读取设备状态消息失败", r.GWParam.ServiceName)
		return
	}
	setting.ZAPS.Infof("SagooIOT上报服务[%s]发布读取设备状态消息成功", r.GWParam.ServiceName)
}

func (r *ReportServiceParamSagooIOTTemplate) ReportServiceSagooIOTProcessReadNodeStatus(reqFrame MQTTSagooIOTReadNodeStatusRequestTemplate) {

	ackParams := make([]MQTTSagooIOTReadNodeStatusAckParamTemplate, 0)
	nodeStatus := MQTTSagooIOTReadNodeStatusAckParamTemplate{}
	for _, node := range r.NodeList {
		nodeStatus.ClientID = node.Param.DeviceCode
		nodeStatus.CommStatus = node.CommStatus
		ackParams = append(ackParams, nodeStatus)
	}

	r.ReportServiceSagooIOTReadNodeStatusAck(reqFrame, 0, ackParams)

}
