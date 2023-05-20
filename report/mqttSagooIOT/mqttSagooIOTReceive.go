package mqttSagooIOT

import (
	"gateway/setting"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTSagooIOTReceiveFrameTemplate struct {
	Topic   string
	Payload []byte
}

//发送数据回调函数
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListSagooIOT.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTSagooIOTReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}

			setting.ZAPS.Debugf("SagooIOT MQTT接收消息主题:%s", receiveFrame.Topic)
			setting.ZAPS.Debugf("SagooIOT MQTT接收消息内容:%s", receiveFrame.Payload)
			ReportServiceParamListSagooIOT.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
