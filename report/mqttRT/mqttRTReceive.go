package mqttRT

import (
	"gateway/setting"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTRTReceiveFrameTemplate struct {
	Topic   string
	Payload []byte
}

//发送数据回调函数
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListRT.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTRTReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}

			setting.ZAPS.Debugf("RT MQTT接收消息主题:%s", receiveFrame.Topic)
			setting.ZAPS.Debugf("RT MQTT接收消息内容:%s", receiveFrame.Payload)
			ReportServiceParamListRT.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
