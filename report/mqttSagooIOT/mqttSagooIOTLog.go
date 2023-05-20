package mqttSagooIOT

import (
	"time"
)

type MQTTSagooIOTCommunicationCmdTemplate struct {
	CmdName           string `json:"CmdName"`           //命令名称
	CmdRequestContent string `json:"CmdRequestContent"` //请求内容
	CmdAckContent     string `json:"CmdAckContent"`     //回应内容
	CmdStartTime      string `json:"CmdStartTime"`      //开始时间戳
	CmdStopTime       string `json:"CmdStopTime"`       //结束时间戳
}

type MessageTemplate struct {
	TimeStamp string `json:"timeStamp"` //时间戳
	Direction int    `json:"direction"` //数据方向
	Label     string `json:"label"`     //数据标识
	Content   string `json:"content"`   //数据内容
}

const (
	Direction_TX = 1
	Direction_RX = 0
)

func MQTTSagooIOTAddCommunicationMessage(service *ReportServiceParamSagooIOTTemplate, label string, dir int, buf string) {

	msg := MessageTemplate{
		TimeStamp: time.Now().Format("2006-01-02 15:04:05.1234"),
		Direction: dir,
		Label:     label,
		//Content:   fmt.Sprintf("%X", buf),
		Content: buf,
	}

	service.MessageEventBus.Publish(service.GWParam.ServiceName, msg)
}
