package mqttRT

import (
	"encoding/json"
	"gateway/device"
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

type MQTTRTLogInDataTemplate struct {
	ProductKey string `json:"productKey,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
}

type MQTTRTLogInAckTemplate struct {
	ID      string                    `json:"id"`
	Code    int32                     `json:"code"`
	Message string                    `json:"message"`
	Data    []MQTTRTLogInDataTemplate `json:"data,omitempty"`
}

var MsgID int = 0

func (r *ReportServiceParamRTTemplate) MQTTRTOnConnectHandler(client MQTT.Client) {

	setting.ZAPS.Debug("MQTTRT 链接成功")
	r.GWParam.ReportStatus = "onLine"
	MQTTRTAddCommunicationMessage(r, "MQTTRT 链接成功", Direction_RX, "")

	setting.LedOnOff(setting.LEDServer, setting.LEDON)

	r.GWParam.MQTTClient = client

	subTopic := ""
	//订阅子设备属性下发请求
	subTopic = "/device/data/set/" + r.GWParam.Param.ClientID
	MQTTRTSubscribeTopic(r.GWParam.ServiceName, r.GWParam.MQTTClient, subTopic)

	//订阅子设备属性查询请求
	subTopic = "/device/data/get/" + r.GWParam.Param.ClientID
	MQTTRTSubscribeTopic(r.GWParam.ServiceName, r.GWParam.MQTTClient, subTopic)

	//订阅子设备服务调用请求
	subTopic = "/device/data/cmd/" + r.GWParam.Param.ClientID
	MQTTRTSubscribeTopic(r.GWParam.ServiceName, r.GWParam.MQTTClient, subTopic)
}

func (r *ReportServiceParamRTTemplate) MQTTRTConnectionLostHandler(client MQTT.Client, err error) {

	MQTTRTAddCommunicationMessage(r, "MQTTRT 链接断开", Direction_RX, "")
	setting.LedOnOff(setting.LEDServer, setting.LEDOFF)
	setting.ZAPS.Infof("MQTTRT 上报服务名称[%v]链接断开 %v", r.GWParam.ServiceName, err)
	r.GWParam.ReportStatus = "offLine"
}

func MQTTRTGWLogin(service *ReportServiceParamRTTemplate, publishHandler MQTT.MessageHandler) (bool, *MQTT.ClientOptions, MQTT.Client) {

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
	opts.SetConnectTimeout(time.Duration(RTTimeOutLogin) * time.Second)
	opts.SetAutoReconnect(true)

	mqttClient := MQTT.NewClient(opts)

	setting.ZAPS.Infof("上报服务[%s]网关登录", service.GWParam.ServiceName)
	MQTTRTAddCommunicationMessage(service, "网关登录", Direction_TX, "")
	token := mqttClient.Connect()
	token.Wait()
	select {
	case <-token.Done():
		{
			err := token.Error()
			if err != nil {
				setting.ZAPS.Errorf("上报服务[%s]网关登录失败 %v", service.GWParam.ServiceName, err)
				setting.LedOnOff(setting.LEDServer, setting.LEDOFF)
				return false, opts, nil
			}
		}
	}
	setting.ZAPS.Infof("上报服务[%s]网关登录成功", service.GWParam.ServiceName)
	//event.AddReportEvent(service.GWParam.ServiceName, "logIN")
	//MQTTRTAddCommunicationMessage(service, "网关登录成功", Direction_RX, nil)

	return true, opts, mqttClient
}

func MQTTRTSubscribeTopic(serviceName string, client MQTT.Client, topic string) {

	token := client.Subscribe(topic, 1, nil)
	if token.WaitTimeout(time.Duration(RTTimeOutSubscribe)*time.Second) == false {
		setting.ZAPS.Errorf("RT上报服务[%s]订阅主题[%s]失败 %v", serviceName, topic, token.Error())
		return
	}
	setting.ZAPS.Infof("RT上报服务[%s]订阅主题[%s]成功", serviceName, topic)
}

func (r *ReportServiceParamRTTemplate) GWLogin() bool {
	status := false
	status, r.GWParam.MQTTClientOptions, r.GWParam.MQTTClient = MQTTRTGWLogin(r, ReceiveMessageHandler)

	return status
}

func MQTTRTNodeLoginIn(param ReportServiceGWParamRTTemplate, nodeMap []string) int {

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
			//token := param.MQTTClient.Publish(loginTopic, 1, false, sJson)
			//token.WaitTimeout(1000 * time.Millisecond)

			token := param.MQTTClient.Publish(loginTopic, 0, false, sJson)
			select {
			case <-token.Done():
				{
					if token.Error() != nil {
						setting.ZAPS.Errorf("RT上报服务[%s]发布节点上线消息失败 %v", param.ServiceName, token.Error())
						return MsgID
					}
				}
			case <-time.After(time.Duration(RTTimeOutLogin) * time.Second):
				{
					setting.ZAPS.Errorf("RT上报服务[%s]发布节点上线消息失败 %v", param.ServiceName, token.Error())
					return MsgID
				}
			}
		}
	}
	setting.ZAPS.Infof("RT上报服务[%s]发布节点上线消息成功", param.ServiceName)

	return MsgID
}

func (r *ReportServiceParamRTTemplate) NodeLogIn(name string) error {

	setting.ZAPS.Infof("上报服务[%s]节点%s上线", r.GWParam.ServiceName, name)
	propertyPostParamMap := make([]MQTTRTPropertyPostParamTemplate, 0)
	for k, v := range r.NodeList {
		if name == v.Name {
			//传输编码为空不上报
			if v.Param.DeviceCode == "" {
				continue
			}

			//上报故障计数值先加，收到正确回应后清0
			r.NodeList[k].ReportErrCnt++
			onLine := false
			if v.CommStatus == "onLine" {
				onLine = true
			}
			propertyPostParam := MQTTRTPropertyPostParamTemplate{
				DeviceCode: v.Param.DeviceCode,
				OnLine:     onLine,
				VarData:    make([]MQTTRTPropertyPostParamPropertyTemplate, 0),
			}

			property := MQTTRTPropertyPostParamPropertyTemplate{}
			for _, p := range v.Properties {
				coll, ok := device.CollectInterfaceMap.Coll[v.CollInterfaceName]
				if !ok {
					continue
				}
				node, ok := coll.DeviceNodeMap[v.Name]
				if !ok {
					continue
				}
				for _, n := range node.Properties {
					if n.Name == p.Name {
						if len(n.Value) == 0 {
							continue
						}
						property.Name = p.Name
						property.TimeStamp = n.Value[len(n.Value)-1].TimeStamp.Unix()
						property.Value = n.Value[len(n.Value)-1].Value
						propertyPostParam.VarData = append(propertyPostParam.VarData, property)
					}
				}
			}

			//单个设备发送
			propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)
			r.MQTTRTPropertyPost(DeviceTypeNode, r.GWParam, propertyPostParamMap[len(propertyPostParamMap)-1:], false)
		}
	}

	return nil
}
