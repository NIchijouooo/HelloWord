package mqttFeisjy

import (
	"encoding/json"
	"fmt"
	"gateway/setting"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"net"
	"time"
)

const (
	FeisjyMQTTTopicRxFormat = "iot/rx/%s/%s/%s"
	FeisjyMQTTTopicTxFormat = "iot/tx/%s/%s/%s"
)

// 设备发送到平台的登录请求报文
type MQTTFeisjyLogInAskTemplate struct {
	ProductKey   string `json:"productKey"`
	DeviceSecret string `json:"deviceSecret"`
	User         string `json:"user"`
	Pwd          string `json:"pwd"`
	Status       string `json:"status"`
	ID           int32  `json:"id"`
}

// 平台发送到设备的登录确认报文
type MQTTFeisjyLogInAckTemplate struct {
	Status string `json:"status"`
	ID     int32  `json:"id"`
}

// 设备发给平台的登录完成确认报文
type MQTTFeisjyLogInFinTemplate struct {
	Status string `json:"status"`
	ID     int32  `json:"id"`
}

// 连接MQTT Broker成功回调
func MQTTFeisjyOnConnectHandler(client MQTT.Client) {
	setting.ZAPS.Debug("MQTTFeisjy 链接成功")
	for _, v := range ReportServiceParamListFeisjy.ServiceList {
		if v.GWParam.MQTTClient == client {
			v.MQTTFeisjySubTopicList(client, v.GWParam)
			for _, node := range v.NodeList {
				v.MQTTFeisjySubNodeTopic(node.Param.DeviceID)
			}
		}
	}
}

// MQTT Broker连接断开回调
func MQTTFeisjyConnectionLostHandler(client MQTT.Client, err error) {
	setting.ZAPS.Debug("MQTTFeisjy 链接断开")
}

// MQTT Broker链接重新链接回调
func MQTTFeisjyReconnectHandler(client MQTT.Client, opt *MQTT.ClientOptions) {
	setting.ZAPS.Debug("MQTTFeisjy 链接重新链接")
}

// 用于指定默认的消息处理函数。如果客户端订阅了一个主题但没有指定回调函数，则会使用默认的消息处理函数。
func MQTTFeisjySubscribeTopic(client MQTT.Client, topic string) {

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		setting.ZAPS.Warnf("Feisjy上报服务订阅主题%s失败 %v", topic, token.Error())
		return
	}
	setting.ZAPS.Infof("Feisjy上报服务订阅主题%s成功", topic)
}

// 用于指定默认的消息处理函数。如果客户端订阅了一个主题但没有指定回调函数，则会使用默认的消息处理函数。
func MQTTFeisjyUnsubscribeTopic(client MQTT.Client, topic string) {

	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		setting.ZAPS.Warnf("Feisjy上报服务解除订阅主题%s失败 %v", topic, token.Error())
		return
	}
	setting.ZAPS.Infof("Feisjy上报服务解除订阅主题%s成功", topic)
}

// Feisjy平台需要订阅的主题统一在这里添加
func (r *ReportServiceParamFeisjyTemplate) MQTTFeisjySubTopicList(client MQTT.Client, param ReportServiceGWParamFeisjyTemplate) {

	subTopic := ""

	//平台登陆确认报文
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, param.Param.AppKey, param.Param.DeviceID, "login")
	MQTTFeisjySubscribeTopic(client, subTopic)

	//设备总招命令下发
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, param.Param.AppKey, param.Param.DeviceID, "deviceControl")
	MQTTFeisjySubscribeTopic(client, subTopic)

	// 设备固件升级命令下发
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, param.Param.AppKey, param.Param.DeviceID, "deviceUpgrade")
	MQTTFeisjySubscribeTopic(client, subTopic)

	// 获取文件列表命令下发
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, param.Param.AppKey, param.Param.DeviceID, "fileList")
	MQTTFeisjySubscribeTopic(client, subTopic)

	// 平台回复信息
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, param.Param.AppKey, param.Param.DeviceID, "resultMsg")
	MQTTFeisjySubscribeTopic(client, subTopic)

	//for _, node := range r.NodeList {
	//	//平台登陆确认报文
	//	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, param.Param.AppKey, node.Param.DeviceID, "login")
	//	MQTTFeisjySubscribeTopic(client, subTopic)
	//
	//	//设备总招命令下发
	//	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, param.Param.AppKey, node.Param.DeviceID, "deviceControl")
	//	MQTTFeisjySubscribeTopic(client, subTopic)
	//
	//	// 平台回复信息
	//	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, param.Param.AppKey, node.Param.DeviceID, "resultMsg")
	//	MQTTFeisjySubscribeTopic(client, subTopic)
	//}
}
func (r *ReportServiceParamFeisjyTemplate) MQTTFeisjySubNodeTopic(DeviceID string) {

	subTopic := ""

	//平台登陆确认报文
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, r.GWParam.Param.AppKey, DeviceID, "login")
	MQTTFeisjySubscribeTopic(r.GWParam.MQTTClient, subTopic)

	//设备总招命令下发
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, r.GWParam.Param.AppKey, DeviceID, "deviceControl")
	MQTTFeisjySubscribeTopic(r.GWParam.MQTTClient, subTopic)

	// 平台回复信息
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, r.GWParam.Param.AppKey, DeviceID, "resultMsg")
	MQTTFeisjySubscribeTopic(r.GWParam.MQTTClient, subTopic)
}

func (r *ReportServiceParamFeisjyTemplate) MQTTFeisjyUnsubNodeTopic(DeviceID string) {

	subTopic := ""

	//平台登陆确认报文
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, r.GWParam.Param.AppKey, DeviceID, "login")
	MQTTFeisjyUnsubscribeTopic(r.GWParam.MQTTClient, subTopic)

	//设备总招命令下发
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, r.GWParam.Param.AppKey, DeviceID, "deviceControl")
	MQTTFeisjyUnsubscribeTopic(r.GWParam.MQTTClient, subTopic)

	// 平台回复信息
	subTopic = fmt.Sprintf(FeisjyMQTTTopicTxFormat, r.GWParam.Param.AppKey, DeviceID, "resultMsg")
	MQTTFeisjyUnsubscribeTopic(r.GWParam.MQTTClient, subTopic)
}

// MQTT Broker配置，连接
func (r *ReportServiceParamFeisjyTemplate) MQTTFeisjyGWLogin(param ReportServiceGWParamFeisjyTemplate, publishHandler MQTT.MessageHandler) (bool, *MQTT.ClientOptions, MQTT.Client) {

	//配置上报平台使用的网卡
	localIP := setting.GetIPByNetCardName(setting.ReportNet)
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP: net.ParseIP(localIP), // 指定使用的网卡 IP 地址 hxd modity
		},
	}
	setting.ZAPS.Infof("setting report name ->%s  lcoalIP -> %s", setting.ReportNet, localIP)

	//QJHui add 2023/6/16 新增根据网卡配置路由表信息
	if localIP != "" {
		defaultRoute := fmt.Sprintf("ip route replace 0.0.0.0/0 via 0.0.0.0 dev %s", setting.ReportNet)
		setting.Exec_shell(defaultRoute)
	}

	opts := MQTT.NewClientOptions().AddBroker(param.IP + ":" + param.Port)

	opts.SetClientID(param.Param.ClientID)
	opts.SetUsername(param.Param.UserName)
	opts.SetPassword(param.Param.Password)
	opts.SetDialer(dialer) // hxd modity
	opts.SetKeepAlive(30 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler) //用于指定默认的消息处理函数。如果客户端订阅了一个主题但没有指定回调函数，则会使用默认的消息处理函数。
	opts.SetAutoReconnect(true)
	//opts.SetConnectRetry(false)
	opts.SetConnectTimeout(2 * time.Second)
	opts.SetOnConnectHandler(MQTTFeisjyOnConnectHandler)
	opts.SetConnectionLostHandler(MQTTFeisjyConnectionLostHandler)
	opts.SetReconnectingHandler(MQTTFeisjyReconnectHandler)

	// create and start a client using the above ClientOptions
	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		setting.ZAPS.Errorf("上报服务[%s]链接Feisjy Broker失败 %v", param.ServiceName, token.Error())
		return false, opts, nil
	}
	setting.ZAPS.Infof("上报服务[%s]链接Feisjy Broker成功", param.ServiceName)

	r.MQTTFeisjySubTopicList(mqttClient, param)

	return true, opts, mqttClient
}

func (r *ReportServiceParamFeisjyTemplate) GWLogin() bool {
	status := false
	status, r.GWParam.MQTTClientOptions, r.GWParam.MQTTClient = r.MQTTFeisjyGWLogin(r.GWParam, ReceiveMessageHandler)

	return status
}

// 登录请求报文
func (r *ReportServiceParamFeisjyTemplate) FeisjyLogInAsk(msg *MQTTFeisjyLogInAskTemplate, id string) bool {
	status := false

	sJson, _ := json.Marshal(msg)
	propertyPostTopic := fmt.Sprintf(FeisjyMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "login")

	if r.GWParam.MQTTClient != nil {
		if token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson); token.WaitTimeout(2000*time.Millisecond) && token.Error() != nil {
			status = false
			setting.ZAPS.Debugf("上报服务[%s]发布[%s]登录请求消息失败 %v", r.GWParam.ServiceName, id, token.Error())
		} else {
			status = true
			setting.ZAPS.Debugf("上报服务[%s]发布[%s]登录请求消息成功 内容%v", r.GWParam.ServiceName, id, string(sJson))
		}
	}

	return status
}

// 完成登录确认报文
func (r *ReportServiceParamFeisjyTemplate) FeisjyLogInFin(msg *MQTTFeisjyLogInFinTemplate, id string) bool {
	status := false
	sJson, _ := json.Marshal(msg)
	propertyPostTopic := fmt.Sprintf(FeisjyMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "login")

	if r.GWParam.MQTTClient != nil {
		if token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson); token.WaitTimeout(2000*time.Millisecond) && token.Error() != nil {
			status = false
			setting.ZAPS.Debugf("上报服务[%s]发布[%s]登录确认消息失败 %v", r.GWParam.ServiceName, id, token.Error())
		} else {
			status = true
			setting.ZAPS.Debugf("上报服务[%s]发布[%s]登录确认消息成功 内容%v", r.GWParam.ServiceName, id, string(sJson))
		}
	}

	return status
}

// 设备登录流程状态机
func (r *ReportServiceParamFeisjyTemplate) FeisjyLogInMachine(id string) bool {
	status := false
	cnt := 0
	idType := 0
	var loginId int32

	// 使用当前时间的纳秒数作为随机数生成器的种子
	rand.Seed(time.Now().UnixNano())

	// 生成一个0到65536之间的随机数
	loginId = int32(rand.Intn(65536)) //产生一个随机数

	msgAsk := &MQTTFeisjyLogInAskTemplate{
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
		User:         r.GWParam.Param.UserName,
		Pwd:          r.GWParam.Param.Password,
		Status:       "ask",
		ID:           loginId,
	}

	msgFin := &MQTTFeisjyLogInFinTemplate{
		Status: "fin",
		ID:     loginId + 2,
	}

	//发起登录ask(loginId)
	if id == r.GWParam.Param.ClientID {
		status = r.FeisjyLogInAsk(msgAsk, r.GWParam.Param.DeviceID)
		idType = 0
	} else {
		for _, node := range r.NodeList {
			if node.Param.DeviceID == id {
				msgAsk.ProductKey = node.Param.ProductKey
				msgAsk.DeviceSecret = node.Param.DeviceSecret
				status = r.FeisjyLogInAsk(msgAsk, id)
				idType = 1
			}
		}
	}

	if status == false {
		return status
	}

	//等待平台返回确认req(loginId+1)并发送登陆完成fin(loginId+2),5秒内收不到就退出
	for {
		select {
		case finFrame := <-r.ReceiveLogInAckFrameChan: //收到平台发过来的登录确认帧
			{
				if finFrame.ID == loginId+1 {
					if idType == 0 {
						return r.FeisjyLogInFin(msgFin, r.GWParam.Param.DeviceID)
					} else if idType == 1 {
						return r.FeisjyLogInFin(msgFin, id)
					} else {
						return false
					}
				}
			}
		default:
			{
				cnt++
				if cnt > 50 {
					return false
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// 设备心跳信息
func (r *ReportServiceParamFeisjyTemplate) FeisjyHeartBeat(id string, heartBeatCnt uint32) bool {
	status := false

	msg := struct {
		ID uint32 `json:"id"`
	}{
		ID: heartBeatCnt,
	}

	sJson, _ := json.Marshal(msg)
	propertyPostTopic := fmt.Sprintf(FeisjyMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "heartbeat")

	if r.GWParam.MQTTClient != nil {
		if token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson); token.WaitTimeout(2000*time.Millisecond) && token.Error() != nil {
			status = false
			setting.ZAPS.Debugf("上报服务[%s]发布[%s]心跳消息失败 %v", r.GWParam.ServiceName, id, token.Error())
		} else {
			status = true
			setting.ZAPS.Debugf("上报服务[%s]发布[%s]心跳确认消息成功 内容%v", r.GWParam.ServiceName, id, string(sJson))
		}
	}

	return status
}
