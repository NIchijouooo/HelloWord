package mqttThingsBoard

import (
	"gateway/setting"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTThingsBoardReceiveFrameTemplate struct {
	Topic   string
	Payload []byte
}

//发送数据回调函数
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListThingsBoard.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTThingsBoardReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}
			setting.ZAPS.Debugf("上报服务[%v]接收主题 %s", v.GWParam.ServiceName, receiveFrame.Topic)
			setting.ZAPS.Debugf("上报服务[%v]接收内容 %s", v.GWParam.ServiceName, receiveFrame.Payload)
			ReportServiceParamListThingsBoard.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
