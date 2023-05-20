package mqttM2M

import (
	"gateway/setting"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTM2MReceiveFrameTemplate struct {
	Topic   string
	Payload []byte
}

//ReceiveMessageHandler 是MQTT订阅的主题回调函数，订阅的主题有数据过来，会执行这个回调，把数据传入ReceiveFrameChan通道
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListM2M.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTM2MReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}
			setting.ZAPS.Debugf("m2m MQTT接收消息主题:%s", receiveFrame.Topic)
			setting.ZAPS.Debugf("m2m MQTT接收消息内容:%s", receiveFrame.Payload)
			ReportServiceParamListM2M.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
