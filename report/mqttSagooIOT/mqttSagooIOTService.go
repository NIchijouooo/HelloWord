package mqttSagooIOT

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
	"net/http"
	"strings"
	"time"
)

type MQTTSagooIOTInvokeServiceAckTemplate struct {
	Code       int                    `json:"code"`
	ClientID   string                 `json:"clientID"`
	MsgID      string                 `json:"msgID"`
	DeviceCode string                 `json:"devCode"`
	CmdName    string                 `json:"cmdName"`
	CmdReturn  map[string]interface{} `json:"cmdReturn"`
}

type MQTTSagooIOTInvokeServiceRequestTemplate struct {
	ClientID   string                 `json:"clientID"`
	MsgID      string                 `json:"msgID"`
	DeviceCode string                 `json:"devCode"`
	CmdName    string                 `json:"cmdName"`
	CmdData    map[string]interface{} `json:"cmdData"`
}

func (r *ReportServiceParamSagooIOTTemplate) ReportServiceSagooIOTInvokeServiceAck(reqFrame MQTTSagooIOTInvokeServiceRequestTemplate, deviceType int, ackParam MQTTSagooIOTInvokeServiceAckTemplate) {

	sJson, _ := json.Marshal(ackParam)
	serviceInvokeTopic := ""
	if deviceType == DeviceTypeGW {
		serviceInvokeTopic = "/device/data/cmd_replay/" + r.GWParam.Param.ClientID
	} else {
		serviceInvokeTopic = "/device/data/cmd_replay/" + r.GWParam.Param.ClientID
	}
	setting.ZAPS.Infof("SagooIOT上报服务发布回复写自定义命令应答消息主题 %s", serviceInvokeTopic)
	setting.ZAPS.Debugf("SagooIOT上报服务发布回复写自定义命令应答消息内容 %v", string(sJson))
	MQTTSagooIOTAddCommunicationMessage(r, "MQTT写自定义命令应答包", Direction_TX, string(sJson))

	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(serviceInvokeTopic, 1, false, sJson)
		if token.WaitTimeout(time.Duration(SagooIOTTimeOutService)*time.Second) == false {
			setting.ZAPS.Errorf("SagooIOT上报服务[%s]发布回复服务调用消息失败 %v", r.GWParam.ServiceName, token.Error())
			return
		}
	} else {
		setting.ZAPS.Errorf("SagooIOT上报服务[%s]发布回复服务调用消息失败", r.GWParam.ServiceName)
		return
	}
	setting.ZAPS.Debugf("SagooIOT上报服务[%s]发布回复服务调用消息成功", r.GWParam.ServiceName)
}

func (r *ReportServiceParamSagooIOTTemplate) ReportServiceSagooIOTProcessInvokeService(reqFrame MQTTSagooIOTInvokeServiceRequestTemplate) {

	ackParam := MQTTSagooIOTInvokeServiceAckTemplate{
		ClientID:   r.GWParam.Param.ClientID,
		MsgID:      reqFrame.MsgID,
		DeviceCode: reqFrame.DeviceCode,
		CmdName:    reqFrame.CmdName,
		CmdReturn:  make(map[string]interface{}),
	}

	rData, _ := json.Marshal(reqFrame)
	MQTTSagooIOTAddCommunicationMessage(r, "MQTT写自定义命令请求包", Direction_RX, string(rData))

	//命令是针对网关的命令
	if reqFrame.DeviceCode == r.GWParam.Param.ClientID {
		ackParam.Code = r.ReportServiceSagooIOTProcessInvokeGWService(reqFrame)
		//r.ReportServiceSagooIOTInvokeServiceAck(reqFrame, DeviceTypeNode, ackParam)
	} else {
		//命令是针对设备的命令
		for _, node := range r.NodeList {
			if reqFrame.DeviceCode == node.Param.DeviceCode {
				//从上报节点中找到相应节点
				coll, ok := device.CollectInterfaceMap.Coll[node.CollInterfaceName]
				if !ok {
					continue
				}

				for _, n := range coll.DeviceNodeMap {
					if n.Name == node.Name {
						//从采集服务中找到相应节点
						cmd := device.CommunicationCmdTemplate{}
						cmd.CollInterfaceName = node.CollInterfaceName
						cmd.DeviceName = node.Name
						cmd.FunName = reqFrame.CmdName
						paramStr, _ := json.Marshal(reqFrame.CmdData)
						cmd.FunPara = string(paramStr)

						ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
						if ackData.Status {
							ackParam.Code = 0
							for _, v := range ackData.Properties {
								ackParam.CmdReturn[v.Name] = v.Value
							}
						} else {
							ackParam.Code = 1
						}
						r.ReportServiceSagooIOTInvokeServiceAck(reqFrame, DeviceTypeNode, ackParam)
					}
				}
			}
		}
	}
}

func (r *ReportServiceParamSagooIOTTemplate) ReportServiceSagooIOTProcessInvokeGWService(reqFrame MQTTSagooIOTInvokeServiceRequestTemplate) int {

	ackParam := MQTTSagooIOTInvokeServiceAckTemplate{
		ClientID:   r.GWParam.Param.ClientID,
		MsgID:      reqFrame.MsgID,
		DeviceCode: reqFrame.DeviceCode,
		CmdName:    reqFrame.CmdName,
		CmdReturn:  make(map[string]interface{}),
	}

	switch reqFrame.CmdName {
	case "SetIAP": //升级任务
		{
			url, ok := reqFrame.CmdData["URL"].(string)
			if !ok {
				return 1
			}

			reader := setting.NewHTTPDownFile()
			go setting.HttpDownFileFromURL(url, "./gatewayUpdate.tar", reader)
			for {
				select {
				//case percnt := <-reader.PercentOut:
				//	{
				//		setting.ZAPS.Debugf("percent %v", percnt)
				//		//param := MQTTTuoQiaoOTAPrecentRequestParamTemplate{
				//		//	FileVersion: reqFrame.Params.FileVersion,
				//		//	FileName:    reqFrame.Params.FileName,
				//		//	Percent:     strconv.Itoa(percnt),
				//		//	Message:     "",
				//		//}
				//		//r.ReportServiceTuoQiaoOTAPercent(reqFrame, param)
				//	}
				case <-reader.Finish:
					{
						index := strings.Index(url, "taskId=")
						if index > 0 {
							taskID := url[index+7:]
							setting.ZAPS.Debugf("taskID %v", taskID)
							_ = r.ReportServiceSagooIOTGWOTA(taskID, reqFrame.ClientID, 100, 0)
						}
						return 0
					}
				case err := <-reader.Error:
					{
						if err != nil {
							setting.ZAPS.Debugf("http读取失败 %v", err)
							return 1
						}
						return 0
					}
				}
			}
		}
	case "GetRealVariables":
		{
			for _, node := range r.NodeList {
				if reqFrame.DeviceCode == node.Param.DeviceCode {
					//从上报节点中找到相应节点
					coll, ok := device.CollectInterfaceMap.Coll[node.CollInterfaceName]
					if !ok {
						continue
					}

					for _, n := range coll.DeviceNodeMap {
						if n.Name == node.Name {
							//从采集服务中找到相应节点
							cmd := device.CommunicationCmdTemplate{}
							cmd.CollInterfaceName = node.CollInterfaceName
							cmd.DeviceName = node.Name
							cmd.FunName = reqFrame.CmdName
							paramStr, _ := json.Marshal(reqFrame.CmdData)
							cmd.FunPara = string(paramStr)

							ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
							if ackData.Status {
								ackParam.Code = 0
								for _, v := range ackData.Properties {
									ackParam.CmdReturn[v.Name] = v.Value
								}
							} else {
								ackParam.Code = 1
							}
							r.ReportServiceSagooIOTInvokeServiceAck(reqFrame, DeviceTypeNode, ackParam)
						}
					}
				}
			}
		}
	case "BrandSet": //设置品牌
		{
			//命令是针对设备的命令
			for _, node := range r.NodeList {
				if reqFrame.DeviceCode == node.Param.DeviceCode {
					//从上报节点中找到相应节点
					coll, ok := device.CollectInterfaceMap.Coll[node.CollInterfaceName]
					if !ok {
						continue
					}

					for _, n := range coll.DeviceNodeMap {
						if n.Name == node.Name {
							//从采集服务中找到相应节点
							cmd := device.CommunicationCmdTemplate{}
							cmd.CollInterfaceName = node.CollInterfaceName
							cmd.DeviceName = node.Name
							cmd.FunName = "SetVariables"
							paramStr, _ := json.Marshal(reqFrame.CmdData)
							cmd.FunPara = string(paramStr)

							ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
							if ackData.Status {
								ackParam.Code = 0
								for _, v := range ackData.Properties {
									ackParam.CmdReturn[v.Name] = v.Value
								}
							} else {
								ackParam.Code = 1
							}
							r.ReportServiceSagooIOTInvokeServiceAck(reqFrame, DeviceTypeNode, ackParam)
						}
					}
				}
			}
		}
	}

	return 0
}

func (r *ReportServiceParamSagooIOTTemplate) ReportServiceSagooIOTGWOTA(taskId string, clientID string, percent int, code int) error {

	url := "http://" + "cloud.reatgreen.com" + "/gateway/device-maintenance/deviceUpgrade/updateDeviceUpgradeResult"

	param := struct {
		TaskId  string `json:"taskId"`
		SnCode  string `json:"snCode"`
		Percent int    `json:"percent"`
		ErrCode int    `json:"errCode"`
		ErrMsg  string `json:"errMsg"`
	}{
		TaskId:  taskId,
		SnCode:  clientID,
		Percent: percent,
		ErrCode: code,
		ErrMsg:  "",
	}

	reqParam, _ := json.Marshal(&param)

	// 准备: HTTP请求
	reqBody := strings.NewReader(string(reqParam))
	httpReq, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		setting.ZAPS.Error("上报服务[%v]构建POST请求失败url: %s, reqBody: %s, err: %v", r.GWParam.ServiceName, url, reqBody, err)
		return err
	}
	httpReq.Header.Add("Content-Type", "application/json")

	// DO: HTTP请求
	response, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		setting.ZAPS.Error("上报服务[%v]发送POST请求失败url: %s, reqBody: %s, err: %v", r.GWParam.ServiceName, url, reqBody, err)
		return err
	}
	defer response.Body.Close()

	return nil
}
