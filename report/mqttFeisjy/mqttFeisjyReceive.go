package mqttFeisjy

import (
	"fmt"
	"gateway/setting"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTFeisjyReceiveFrameTemplate struct {
	Topic   string
	Payload []byte
}

//ReceiveMessageHandler 是MQTT订阅的主题回调函数，订阅的主题有数据过来，会执行这个回调，把数据传入ReceiveFrameChan通道
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListFeisjy.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTFeisjyReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}
			setting.ZAPS.Debugf("Feisjy MQTT接收消息主题:%s", receiveFrame.Topic)
			setting.ZAPS.Debugf("Feisjy MQTT接收消息内容:%s", receiveFrame.Payload)
			MQTTFeisjyAddCommunicationMessage(v, fmt.Sprintf("%s", receiveFrame.Topic), Direction_RX, fmt.Sprintf("%s", receiveFrame.Payload))
			ReportServiceParamListFeisjy.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
