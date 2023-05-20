package mqttEmqx

import (
	"encoding/json"
	"gateway/event"
	"gateway/setting"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTNodeLoginParamTemplate struct {
	ClientID  string `json:"clientID"`
	Timestamp int64  `json:"timestamp"`
}

type MQTTNodeLoginTemplate struct {
	ID      string                       `json:"id"`
	Version string                       `json:"version"`
	Params  []MQTTNodeLoginParamTemplate `json:"params"`
}

type MQTTEmqxLogInDataTemplate struct {
	ProductKey string `json:"productKey,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
}

type MQTTEmqxLogInAckTemplate struct {
	ID      string                      `json:"id"`
	Code    int32                       `json:"code"`
	Message string                      `json:"message"`
	Data    []MQTTEmqxLogInDataTemplate `json:"data,omitempty"`
}

var MsgID int = 0

func (r *ReportServiceParamEmqxTemplate) MQTTEmqxOnConnectHandler(client MQTT.Client) {

	setting.ZAPS.Info("MQTTEmqx 链接成功")
	r.GWParam.ReportStatus = "onLine"
	r.GWParam.MQTTClient = client

	subTopic := ""
	//订阅子设备上线回应
	subTopic = "/sys/thing/node/login/post_reply/" + r.GWParam.Param.ClientID
	MQTTEmqxSubscribeTopic(r.GWParam.ServiceName, client, subTopic)

	//订阅子设备下线回应
	subTopic = "/sys/thing/node/logout/post_reply/" + r.GWParam.Param.ClientID
	MQTTEmqxSubscribeTopic(r.GWParam.ServiceName, client, subTopic)

	//订阅子设备属性上报回应
	subTopic = "/sys/thing/node/property/post_reply/" + r.GWParam.Param.ClientID
	MQTTEmqxSubscribeTopic(r.GWParam.ServiceName, client, subTopic)

	//订阅子设备属性下发请求
	subTopic = "/sys/thing/node/property/set/" + r.GWParam.Param.ClientID
	MQTTEmqxSubscribeTopic(r.GWParam.ServiceName, client, subTopic)

	//订阅子设备属性查询请求
	subTopic = "/sys/thing/node/property/get/" + r.GWParam.Param.ClientID
	MQTTEmqxSubscribeTopic(r.GWParam.ServiceName, client, subTopic)

	//订阅子设备服务调用请求
	subTopic = "/sys/thing/node/service/invoke/" + r.GWParam.Param.ClientID
	MQTTEmqxSubscribeTopic(r.GWParam.ServiceName, client, subTopic)

	//订阅子设备获取设备状态请求
	subTopic = "/sys/thing/node/status/get/" + r.GWParam.Param.ClientID
	MQTTEmqxSubscribeTopic(r.GWParam.ServiceName, client, subTopic)

	//------------------------------------------------------------

	//订阅网关属性上报回应
	subTopic = "/sys/thing/gw/property/post_reply/" + r.GWParam.Param.ClientID
	MQTTEmqxSubscribeTopic(r.GWParam.ServiceName, client, subTopic)

	//订阅网关服务调用请求
	subTopic = "/sys/thing/gw/service/invoke/" + r.GWParam.Param.ClientID
	MQTTEmqxSubscribeTopic(r.GWParam.ServiceName, client, subTopic)
}

func (r *ReportServiceParamEmqxTemplate) MQTTEmqxConnectionLostHandler(client MQTT.Client, err error) {

	setting.ZAPS.Debugf("MQTTEmqx上报服务[%v] 链接断开", r.GWParam.ServiceName)
	r.GWParam.ReportStatus = "offLine"

}

func MQTTEmqxGWLogin(r *ReportServiceParamEmqxTemplate, publishHandler MQTT.MessageHandler) (bool, *MQTT.ClientOptions, MQTT.Client) {

	opts := MQTT.NewClientOptions().AddBroker(r.GWParam.IP + ":" + r.GWParam.Port)

	opts.SetClientID(r.GWParam.Param.ClientID)
	opts.SetUsername(r.GWParam.Param.UserName)
	opts.SetPassword(r.GWParam.Param.Password)
	keepAlive, _ := strconv.Atoi(r.GWParam.Param.KeepAlive)
	opts.SetKeepAlive(time.Duration(keepAlive) * time.Second)
	opts.SetDefaultPublishHandler(publishHandler)
	opts.SetOnConnectHandler(r.MQTTEmqxOnConnectHandler)
	opts.SetConnectionLostHandler(r.MQTTEmqxConnectionLostHandler)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(30 * time.Second)
	opts.SetConnectTimeout(60 * time.Second)
	opts.SetAutoReconnect(true)

	mqttClient := MQTT.NewClient(opts)
	token := mqttClient.Connect()
	token.Wait()
	select {
	case <-token.Done():
		{
			err := token.Error()
			if err != nil {
				setting.ZAPS.Errorf("上报服务[%s]链接EMQX Broker失败 %v", r.GWParam.ServiceName, token.Error())
				return false, opts, nil
			}
		}
	}

	setting.ZAPS.Infof("上报服务[%s]链接EMQX Broker成功", r.GWParam.ServiceName)
	event.AddReportEvent(r.GWParam.ServiceName, "logIN")

	return true, opts, mqttClient

}

func MQTTEmqxSubscribeTopic(serviceName string, client MQTT.Client, topic string) {

	token := client.Subscribe(topic, 1, nil)
	if token.WaitTimeout(5*time.Second) == false {
		setting.ZAPS.Errorf("EMQX上报服务[%s]订阅主题[%s]失败 %v", serviceName, topic, token.Error())
		return
	}
	setting.ZAPS.Infof("EMQX上报服务[%s]订阅主题[%s]成功", serviceName, topic)
}

func (r *ReportServiceParamEmqxTemplate) GWLogin() bool {
	status := false
	status, r.GWParam.MQTTClientOptions, r.GWParam.MQTTClient = MQTTEmqxGWLogin(r, ReceiveMessageHandler)
	if status == true {
		r.GWParam.ReportStatus = "onLine"
	}

	return status
}

func MQTTEmqxNodeLoginIn(param ReportServiceGWParamEmqxTemplate, nodeMap []string) int {

	nodeLogin := MQTTNodeLoginTemplate{
		ID:      strconv.Itoa(MsgID),
		Version: "V1.0",
	}
	MsgID++

	for _, v := range nodeMap {
		nodeLoginParam := MQTTNodeLoginParamTemplate{
			ClientID:  v,
			Timestamp: time.Now().Unix(),
		}
		nodeLogin.Params = append(nodeLogin.Params, nodeLoginParam)
	}

	//批量注册
	loginTopic := "/sys/thing/node/login/post/" + param.Param.ClientID

	sJson, _ := json.Marshal(nodeLogin)
	if len(nodeLogin.Params) > 0 {
		setting.ZAPS.Infof("上报服务[%s]发布节点上线消息主题%s", param.ServiceName, loginTopic)
		setting.ZAPS.Debugf("上报服务[%s]发布节点上线消息内容%s", param.ServiceName, sJson)

		if param.MQTTClient != nil {
			token := param.MQTTClient.Publish(loginTopic, 1, false, sJson)
			if token.WaitTimeout(5*time.Second) == false {
				setting.ZAPS.Errorf("EMQX上报服务[%s]发布节点上线消息失败 %v", param.ServiceName, token.Error())
				return MsgID
			}
		} else {
			setting.ZAPS.Errorf("EMQX上报服务[%s]发布节点上线消息失败", param.ServiceName)
			return MsgID
		}
		setting.ZAPS.Infof("EMQX上报服务[%s]发布节点上线消息成功", param.ServiceName)
	}

	return MsgID
}

func (r *ReportServiceParamEmqxTemplate) NodeLogIn(name []string) bool {

	nodeMap := make([]string, 0)
	status := false

	setting.ZAPS.Debugf("上报服务[%s]节点%s上线", r.GWParam.ServiceName, name)
	for _, d := range name {
		for _, v := range r.NodeList {
			if d == v.Name {
				nodeMap = append(nodeMap, v.Param.DeviceCode)

				MQTTEmqxNodeLoginIn(r.GWParam, nodeMap)
				select {
				case frame := <-r.ReceiveLogInAckFrameChan:
					{
						if frame.Code == 200 {
							status = true
						}
					}
				case <-time.After(time.Millisecond * 2000):
					{
						status = false
					}
				}
			}
		}
	}

	return status
}
