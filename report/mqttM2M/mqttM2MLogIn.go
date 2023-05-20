package mqttM2M

import (
	"encoding/json"
	"fmt"
	"gateway/setting"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"time"
)

const (
	M2MMQTTTopicRxFormat = "feisjy/rx/%s/%s/%s"
	M2MMQTTTopicTxFormat = "feisjy/tx/%s/%s/%s"
)

// 设备发送到平台的登录请求报文
type MQTTM2MLogInAskTemplate struct {
	User   string `json:"user"`
	Pwd    string `json:"pwd"`
	Status string `json:"status"`
	ID     int32  `json:"id"`
}

// 平台发送到设备的登录确认报文
type MQTTM2MLogInAckTemplate struct {
	Status string `json:"status"`
	ID     int32  `json:"id"`
}

// 设备发给平台的登录完成确认报文
type MQTTM2MLogInFinTemplate struct {
	Status string `json:"status"`
	ID     int32  `json:"id"`
}

// 连接MQTT Broker成功回调
func MQTTM2MOnConnectHandler(client MQTT.Client) {
	setting.ZAPS.Debug("MQTTM2M 链接成功")
	for _, v := range ReportServiceParamListM2M.ServiceList {
		if v.GWParam.MQTTClient == client {
			v.MQTTM2MSubTopicList(client, v.GWParam)
		}
	}
}

// MQTT Broker连接断开回调
func MQTTM2MConnectionLostHandler(client MQTT.Client, err error) {
	setting.ZAPS.Debug("MQTTM2M 链接断开")
}

// MQTT Broker链接重新链接回调
func MQTTM2MReconnectHandler(client MQTT.Client, opt *MQTT.ClientOptions) {
	setting.ZAPS.Debug("MQTTM2M 链接重新链接")
}

// 用于指定默认的消息处理函数。如果客户端订阅了一个主题但没有指定回调函数，则会使用默认的消息处理函数。
func MQTTM2MSubscribeTopic(client MQTT.Client, topic string) {

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		setting.ZAPS.Warnf("M2M上报服务订阅主题%s失败 %v", topic, token.Error())
		return
	}
	setting.ZAPS.Infof("M2M上报服务订阅主题%s成功", topic)
}

// M2M平台需要订阅的主题统一在这里添加
func (r *ReportServiceParamM2MTemplate) MQTTM2MSubTopicList(client MQTT.Client, param ReportServiceGWParamM2MTemplate) {

	subTopic := ""

	//平台登陆确认报文
	subTopic = fmt.Sprintf(M2MMQTTTopicTxFormat, param.Param.AppKey, param.Param.DeviceID, "login")
	MQTTM2MSubscribeTopic(client, subTopic)

	//设备总招命令下发
	subTopic = fmt.Sprintf(M2MMQTTTopicTxFormat, param.Param.AppKey, param.Param.DeviceID, "allcall")
	MQTTM2MSubscribeTopic(client, subTopic)

	subTopic = fmt.Sprintf(M2MMQTTTopicTxFormat, param.Param.AppKey, param.Param.DeviceID, "resultYx")
	MQTTM2MSubscribeTopic(client, subTopic)

	subTopic = fmt.Sprintf(M2MMQTTTopicTxFormat, param.Param.AppKey, param.Param.DeviceID, "resultYc")
	MQTTM2MSubscribeTopic(client, subTopic)

	subTopic = fmt.Sprintf(M2MMQTTTopicTxFormat, param.Param.AppKey, param.Param.DeviceID, "setting")
	MQTTM2MSubscribeTopic(client, subTopic)

}

// MQTT Broker配置，连接
func (r *ReportServiceParamM2MTemplate) MQTTM2MGWLogin(param ReportServiceGWParamM2MTemplate, publishHandler MQTT.MessageHandler) (bool, *MQTT.ClientOptions, MQTT.Client) {

	setting.ZAPS.Info(param)

	opts := MQTT.NewClientOptions().AddBroker(param.IP + ":" + param.Port)

	opts.SetClientID(param.Param.ClientID)
	opts.SetUsername(param.Param.UserName)
	opts.SetPassword(param.Param.Password)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler) //用于指定默认的消息处理函数。如果客户端订阅了一个主题但没有指定回调函数，则会使用默认的消息处理函数。
	opts.SetAutoReconnect(true)
	//opts.SetConnectRetry(false)
	opts.SetConnectTimeout(2 * time.Second)
	opts.SetOnConnectHandler(MQTTM2MOnConnectHandler)
	opts.SetConnectionLostHandler(MQTTM2MConnectionLostHandler)
	opts.SetReconnectingHandler(MQTTM2MReconnectHandler)

	// create and start a client using the above ClientOptions
	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		setting.ZAPS.Errorf("上报服务[%s]链接M2M Broker失败 %v", param.ServiceName, token.Error())
		return false, opts, nil
	}
	setting.ZAPS.Infof("上报服务[%s]链接M2M Broker成功", param.ServiceName)

	r.MQTTM2MSubTopicList(mqttClient, param)

	return true, opts, mqttClient
}

func (r *ReportServiceParamM2MTemplate) GWLogin() bool {
	status := false
	status, r.GWParam.MQTTClientOptions, r.GWParam.MQTTClient = r.MQTTM2MGWLogin(r.GWParam, ReceiveMessageHandler)

	return status
}

// 登录请求报文
func (r *ReportServiceParamM2MTemplate) M2MLogInAsk(msg *MQTTM2MLogInAskTemplate, id string) bool {
	status := false

	sJson, _ := json.Marshal(msg)
	propertyPostTopic := fmt.Sprintf(M2MMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "login")

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
func (r *ReportServiceParamM2MTemplate) M2MLogInFin(msg *MQTTM2MLogInFinTemplate, id string) bool {
	status := false
	sJson, _ := json.Marshal(msg)
	propertyPostTopic := fmt.Sprintf(M2MMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "login")

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
func (r *ReportServiceParamM2MTemplate) M2MLogInMachine(id string) bool {
	status := false
	cnt := 0
	var loginId int32

	setting.ZAPS.Infof("登录设备%s[%s]", id)

	// 使用当前时间的纳秒数作为随机数生成器的种子
	rand.Seed(time.Now().UnixNano())

	// 生成一个0到65536之间的随机数
	loginId = int32(rand.Intn(65536)) //产生一个随机数

	//msgAsk := &MQTTM2MLogInAskTemplate{
	//	User:   r.GWParam.Param.UserName,
	//	Pwd:    r.GWParam.Param.Password,
	//	Status: "ask",
	//	ID:     loginId,
	//}

	msgAsk := &MQTTM2MLogInAskTemplate{
		User:   "wrc",
		Pwd:    "123456",
		Status: "ask",
		ID:     loginId,
	}

	msgFin := &MQTTM2MLogInFinTemplate{
		Status: "fin",
		ID:     loginId + 2,
	}

	//发起登录ask(loginId)
	if id == r.GWParam.Param.ClientID {
		status = r.M2MLogInAsk(msgAsk, r.GWParam.Param.DeviceID)
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
					return r.M2MLogInFin(msgFin, r.GWParam.Param.DeviceID)
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
func (r *ReportServiceParamM2MTemplate) FeisjyHeartBeat(id string, heartBeatCnt uint32) bool {
	status := false

	msg := struct {
		ID uint32 `json:"id"`
	}{
		ID: heartBeatCnt,
	}

	sJson, _ := json.Marshal(msg)
	propertyPostTopic := fmt.Sprintf(M2MMQTTTopicRxFormat, r.GWParam.Param.AppKey, id, "heartbeat")

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
