package mqttEmqx

import (
	"encoding/json"
	"gateway/setting"
	"strconv"
	"time"
)

type MQTTNodeLogoutParamTemplate struct {
	ClientID  string `json:"clientID"`
	Timestamp int64  `json:"timestamp"`
}

type MQTTNodeLogoutTemplate struct {
	ID      string                        `json:"id"`
	Version string                        `json:"version"`
	Params  []MQTTNodeLogoutParamTemplate `json:"params"`
}

type MQTTEmqxLogOutDataTemplate struct {
	Code       int32  `json:"code"`
	Message    string `json:"message"`
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

type MQTTEmqxLogOutAckTemplate struct {
	ID      string                       `json:"id"`
	Code    int32                        `json:"code"`
	Message string                       `json:"message"`
	Data    []MQTTEmqxLogOutDataTemplate `json:"data"`
}

func MQTTEmqxNodeLogOut(param ReportServiceGWParamEmqxTemplate, nodeMap []string) int {

	nodeLogout := MQTTNodeLogoutTemplate{
		ID:      strconv.Itoa(MsgID),
		Version: "V1.0",
	}
	MsgID++

	for _, v := range nodeMap {
		nodeLogoutParam := MQTTNodeLogoutParamTemplate{
			ClientID:  v,
			Timestamp: time.Now().Unix(),
		}
		nodeLogout.Params = append(nodeLogout.Params, nodeLogoutParam)
	}

	//批量注册
	logoutTopic := "/sys/thing/node/logout/post/" + param.Param.ClientID

	sJson, _ := json.Marshal(nodeLogout)
	if len(nodeLogout.Params) > 0 {
		setting.ZAPS.Infof("上报服务[%s]发布节点下线消息主题%s", param.ServiceName, logoutTopic)
		setting.ZAPS.Debugf("上报服务[%s]发布节点下线消息内容%s", param.ServiceName, sJson)

		if param.MQTTClient != nil {
			//token := param.MQTTClient.Publish(loginTopic, 1, false, sJson)
			//token.WaitTimeout(1000 * time.Millisecond)

			token := param.MQTTClient.Publish(logoutTopic, 0, false, sJson)
			select {
			case <-token.Done():
				{
					if token.Error() != nil {
						setting.ZAPS.Errorf("EMQX上报服务[%s]发布节点下线消息失败 %v", param.ServiceName, token.Error())
						return MsgID
					}
				}
			case <-time.After(time.Duration(EMQXTimeOutLogout) * time.Millisecond):
				{
					setting.ZAPS.Errorf("EMQX上报服务[%s]发布节点下线消息失败 %v", param.ServiceName, token.Error())
					return MsgID
				}
			}
		}
		setting.ZAPS.Debugf("EMQX上报服务[%s]发布节点下线消息成功", param.ServiceName)
	}

	return MsgID
}

func (r *ReportServiceParamEmqxTemplate) NodeLogOut(name []string) bool {

	nodeMap := make([]string, 0)
	status := false

	setting.ZAPS.Debugf("上报服务[%s]节点%s离线", r.GWParam.ServiceName, name)
	for _, d := range name {
		for _, v := range r.NodeList {
			if d == v.Name {
				nodeMap = append(nodeMap, v.Param.DeviceCode)

				MQTTEmqxNodeLogOut(r.GWParam, nodeMap)
				select {
				case frame := <-r.ReceiveLogOutAckFrameChan:
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
