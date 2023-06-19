package mqttZxJs

import (
	"fmt"
	"gateway/setting"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net"
	"time"
)

const (
	ZxjsMQTTTopicRxFormat = "/%s/%s/%s"
	ZxjsMQTTTopicTxFormat = "/%s/%s/%s"
)

// 设备发送到平台的登录请求报文
type MQTTZxjsLogInAskTemplate struct {
	ProductKey   string `json:"productKey"`
	DeviceSecret string `json:"deviceSecret"`
	User         string `json:"user"`
	Pwd          string `json:"pwd"`
	Status       string `json:"status"`
	ID           int32  `json:"id"`
}

// 平台发送到设备的登录确认报文
type MQTTZxjsLogInAckTemplate struct {
	Status string `json:"status"`
	ID     int32  `json:"id"`
}

// 设备发给平台的登录完成确认报文
type MQTTZxjsLogInFinTemplate struct {
	Status string `json:"status"`
	ID     int32  `json:"id"`
}

// 连接MQTT Broker成功回调
func MQTTZxjsOnConnectHandler(client MQTT.Client) {
	setting.ZAPS.Debug("MQTTZxjs 链接成功")
	for _, v := range ReportServiceParamListZxjs.ServiceList {
		if v.GWParam.MQTTClient == client {
			v.MQTTZxjsSubTopicList(client, v.GWParam)
		}
	}
}

// MQTT Broker连接断开回调
func MQTTZxjsConnectionLostHandler(client MQTT.Client, err error) {
	setting.ZAPS.Debug("MQTTZxjs 链接断开")
}

// MQTT Broker链接重新链接回调
func MQTTZxjsReconnectHandler(client MQTT.Client, opt *MQTT.ClientOptions) {
	setting.ZAPS.Debug("MQTTZxjs 链接重新链接")
}

// 用于指定默认的消息处理函数。如果客户端订阅了一个主题但没有指定回调函数，则会使用默认的消息处理函数。
func MQTTZxjsSubscribeTopic(client MQTT.Client, topic string) {

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		setting.ZAPS.Warnf("Zxjs上报服务订阅主题%s失败 %v", topic, token.Error())
		return
	}
	setting.ZAPS.Infof("Zxjs上报服务订阅主题%s成功", topic)
}

// 用于指定默认的消息处理函数。如果客户端订阅了一个主题但没有指定回调函数，则会使用默认的消息处理函数。
func MQTTZxjsUnsubscribeTopic(client MQTT.Client, topic string) {

	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		setting.ZAPS.Warnf("Zxjs上报服务解除订阅主题%s失败 %v", topic, token.Error())
		return
	}
	setting.ZAPS.Infof("Zxjs上报服务解除订阅主题%s成功", topic)
}

// Zxjs平台需要订阅的主题统一在这里添加
func (r *ReportServiceParamZxjsTemplate) MQTTZxjsSubTopicList(client MQTT.Client, param ReportServiceGWParamZxjsTemplate) {

	subTopic := ""

	// 云端召读实时数据
	subTopic = fmt.Sprintf(ZxjsMQTTTopicTxFormat, param.Param.ProductSn, param.Param.DeviceSn, "upload/call")
	MQTTZxjsSubscribeTopic(client, subTopic)

	// 云端召读历史数据
	subTopic = fmt.Sprintf(ZxjsMQTTTopicTxFormat, param.Param.ProductSn, param.Param.DeviceSn, "history/call")
	MQTTZxjsSubscribeTopic(client, subTopic)

	// 云端控制指令
	subTopic = fmt.Sprintf(ZxjsMQTTTopicTxFormat, param.Param.ProductSn, param.Param.DeviceSn, "set")
	MQTTZxjsSubscribeTopic(client, subTopic)
}

func (r *ReportServiceParamZxjsTemplate) MQTTZxjsUnsubNodeTopic(DeviceID string) {

	subTopic := ""

	//平台登陆确认报文
	subTopic = fmt.Sprintf(ZxjsMQTTTopicTxFormat, r.GWParam.Param.ProductSn, DeviceID, "upload/call")
	MQTTZxjsUnsubscribeTopic(r.GWParam.MQTTClient, subTopic)

	// 云端召读历史数据
	subTopic = fmt.Sprintf(ZxjsMQTTTopicTxFormat, r.GWParam.Param.ProductSn, DeviceID, "history/call")
	MQTTZxjsUnsubscribeTopic(r.GWParam.MQTTClient, subTopic)

	// 云端控制指令
	subTopic = fmt.Sprintf(ZxjsMQTTTopicTxFormat, r.GWParam.Param.ProductSn, DeviceID, "set")
	MQTTZxjsUnsubscribeTopic(r.GWParam.MQTTClient, subTopic)
}

// MQTT Broker配置，连接
func (r *ReportServiceParamZxjsTemplate) MQTTZxjsGWLogin(param ReportServiceGWParamZxjsTemplate, publishHandler MQTT.MessageHandler) (bool, *MQTT.ClientOptions, MQTT.Client) {

	//配置上报平台使用的网卡
	localIP := setting.GetIPByNetCardName(setting.ReportNet)
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP: net.ParseIP(localIP), // 指定使用的网卡 IP 地址 hxd modity
		},
	}
	setting.ZAPS.Infof("setting report name ->%s  lcoalIP -> %s", setting.ReportNet, localIP)

	opts := MQTT.NewClientOptions().AddBroker(param.IP + ":" + param.Port)

	// 客户端id生成规则${ProductSN}.${DeviceSN}
	opts.SetClientID(fmt.Sprintf(param.Param.ProductSn, ".", param.Param.DeviceSn))
	// 用户名生成规则${ProductSN}|${DeviceSN}|${authmode} authmode: 1表⽰静态注册；2表⽰动态注册
	opts.SetUsername(fmt.Sprintf(param.Param.ProductSn, "|", param.Param.DeviceSn, "|1"))
	// 本地测试
	//opts.SetUsername(param.Param.DeviceSn)
	// 使用设备密码连接
	opts.SetPassword(param.Param.DevicePwd)
	opts.SetDialer(dialer) // hxd modity
	opts.SetKeepAlive(30 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler) //用于指定默认的消息处理函数。如果客户端订阅了一个主题但没有指定回调函数，则会使用默认的消息处理函数。
	opts.SetAutoReconnect(true)
	//opts.SetConnectRetry(false)
	opts.SetConnectTimeout(2 * time.Second)
	opts.SetOnConnectHandler(MQTTZxjsOnConnectHandler)
	opts.SetConnectionLostHandler(MQTTZxjsConnectionLostHandler)
	opts.SetReconnectingHandler(MQTTZxjsReconnectHandler)

	// create and start a client using the above ClientOptions
	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		setting.ZAPS.Errorf("上报服务[%s]链接Zxjs Broker失败 %v", param.ServiceName, token.Error())
		return false, opts, nil
	}
	setting.ZAPS.Infof("上报服务[%s]链接Zxjs Broker成功", param.ServiceName)

	r.MQTTZxjsSubTopicList(mqttClient, param)
	r.GWParam.ReportStatus = "onLine"

	return true, opts, mqttClient
}

func (r *ReportServiceParamZxjsTemplate) GWLogin() bool {
	status := false
	status, r.GWParam.MQTTClientOptions, r.GWParam.MQTTClient = r.MQTTZxjsGWLogin(r.GWParam, ReceiveMessageHandler)

	return status
}
