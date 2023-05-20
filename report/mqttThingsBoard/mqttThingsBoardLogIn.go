package mqttThingsBoard

import (
	"encoding/json"
	"gateway/setting"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTNodeLoginTemplate struct {
	Device string `json:"device"`
}

func (r *ReportServiceParamThingsBoardTemplate) MQTTRTOnConnectHandler(client MQTT.Client) {

	setting.ZAPS.Debug("ThingsBoard.MQTT 链接成功")
	r.GWParam.ReportStatus = "onLine"
	//MQTTRTAddCommunicationMessage(r, "ThingsBoard.MQTT 链接成功", Direction_RX, "")

	setting.LedOnOff(setting.LEDServer, setting.LEDON)

	r.GWParam.MQTTClient = client

	subTopic := ""
	////订阅属性下发请求
	//subTopic = "v1/gateway/attributes/request"
	//MQTTThingsBoardSubscribeTopic(r.GWParam.ServiceName, r.GWParam.MQTTClient, subTopic)

	//订阅RPC下发请求
	//subTopic = "v1/devices/me/rpc/request/+"
	subTopic = "v1/gateway/rpc"
	MQTTThingsBoardSubscribeTopic(r.GWParam.ServiceName, r.GWParam.MQTTClient, subTopic)
}

func (r *ReportServiceParamThingsBoardTemplate) MQTTRTConnectionLostHandler(client MQTT.Client, err error) {

	//MQTTTAddCommunicationMessage(r, "ThingsBoard.MQTT 链接断开", Direction_RX, "")
	setting.ZAPS.Infof("ThingsBoard.MQTT 上报服务名称[%v]链接断开 %v", r.GWParam.ServiceName, err)
	r.GWParam.ReportStatus = "offLine"
}

func MQTTThingsBoardGWLogin(service *ReportServiceParamThingsBoardTemplate, publishHandler MQTT.MessageHandler) (bool, MQTT.Client) {

	opts := MQTT.NewClientOptions().AddBroker(service.GWParam.IP + ":" + service.GWParam.Port)

	opts.SetClientID(service.GWParam.Param.ClientID)
	opts.SetUsername(service.GWParam.Param.UserName)
	opts.SetPassword(service.GWParam.Param.Password)
	keepAlive, _ := strconv.Atoi(service.GWParam.Param.KeepAlive)
	opts.SetKeepAlive(time.Duration(keepAlive) * time.Second)
	opts.SetDefaultPublishHandler(publishHandler)
	opts.SetOnConnectHandler(service.MQTTRTOnConnectHandler)
	opts.SetConnectionLostHandler(service.MQTTRTConnectionLostHandler)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(30 * time.Second)
	opts.SetConnectTimeout(time.Duration(TimeOutLogin) * time.Second)
	opts.SetAutoReconnect(true)
	// create and start a client using the above ClientOptions
	mqttClient := MQTT.NewClient(opts)
	token := mqttClient.Connect()
	token.Wait()
	select {
	case <-token.Done():
		{
			err := token.Error()
			if err != nil {
				setting.ZAPS.Errorf("上报服务[%s]网关登录失败 %v", service.GWParam.ServiceName, err)
				return false, nil
			}
		}
	}
	setting.ZAPS.Infof("上报服务[%s]网关登录成功", service.GWParam.ServiceName)

	return true, mqttClient

}

func MQTTThingsBoardSubscribeTopic(serviceName string, client MQTT.Client, topic string) {

	token := client.Subscribe(topic, 1, nil)
	if token.WaitTimeout(time.Duration(TimeOutSubscribe)*time.Second) == false {
		setting.ZAPS.Errorf("上报服务[%s]订阅主题[%s]失败 %v", serviceName, topic, token.Error())
		return
	}
	setting.ZAPS.Infof("上报服务[%s]订阅主题[%s]成功", serviceName, topic)
}

func (r *ReportServiceParamThingsBoardTemplate) GWLogin() bool {
	status := false
	status, r.GWParam.MQTTClient = MQTTThingsBoardGWLogin(r, ReceiveMessageHandler)

	return status
}

func (r *ReportServiceParamThingsBoardTemplate) MQTTThingsBoardNodeLoginIn(nodeMap []string) int {

	nodeLogin := MQTTNodeLoginTemplate{}

	for _, v := range nodeMap {
		nodeLogin.Device = v

		sJson, _ := json.Marshal(nodeLogin)
		loginTopic := "v1/gateway/connect"

		setting.ZAPS.Infof("上报服务[%s]发布节点上线消息主题%s", r.GWParam.ServiceName, loginTopic)
		setting.ZAPS.Debugf("上报服务[%s]发布节点上线消息内容%s", r.GWParam.ServiceName, sJson)

		if r.GWParam.MQTTClient != nil {
			token := r.GWParam.MQTTClient.Publish(loginTopic, 1, false, sJson)
			if token.WaitTimeout(time.Duration(TimeOutReportProperty)*time.Second) == false {
				setting.ZAPS.Errorf("上报服务[%s]发布主题[%s]失败 %v", r.GWParam.ServiceName, loginTopic, token.Error())
				continue
			}
			setting.ZAPS.Infof("上报服务[%s]发布主题[%s]成功", r.GWParam.ServiceName, loginTopic)
		}
	}
	return 0
}

func (r *ReportServiceParamThingsBoardTemplate) NodeLogIn(name []string) bool {

	nodeMap := make([]string, 0)
	status := false

	setting.ZAPS.Debugf("上报服务[%s]节点%s上线", r.GWParam.ServiceName, name)
	for _, d := range name {
		for _, v := range r.NodeList {
			if d == v.Name {
				nodeMap = append(nodeMap, v.Param.DeviceName)
			}
		}
	}

	if len(nodeMap) > 0 {
		r.MQTTThingsBoardNodeLoginIn(nodeMap)
	}

	return status
}
