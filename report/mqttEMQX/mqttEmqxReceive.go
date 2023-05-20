package mqttEmqx

import (
	"gateway/setting"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTEmqxReceiveFrameTemplate struct {
	Topic   string
	Payload []byte
}

//发送数据回调函数
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListEmqx.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTEmqxReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}

			setting.ZAPS.Debugf("EMQX MQTT接收消息主题:%s", receiveFrame.Topic)
			setting.ZAPS.Debugf("EMQX MQTT接收消息内容:%s", receiveFrame.Payload)
			ReportServiceParamListEmqx.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
