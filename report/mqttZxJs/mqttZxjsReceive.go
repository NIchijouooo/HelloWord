package mqttZxJs

import (
	"gateway/setting"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTZxjsReceiveFrameTemplate struct {
	Topic   string
	Payload []byte
}

// ReceiveMessageHandler 是MQTT订阅的主题回调函数，订阅的主题有数据过来，会执行这个回调，把数据传入ReceiveFrameChan通道
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListZxjs.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTZxjsReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}
			setting.ZAPS.Debugf("Zxjs MQTT接收消息主题:%s, 内容:%s", receiveFrame.Topic, receiveFrame.Payload)
			ReportServiceParamListZxjs.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
