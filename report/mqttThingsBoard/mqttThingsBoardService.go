package mqttThingsBoard

import (
	"encoding/json"
	"gateway/device"
	"gateway/setting"
	"time"
)

type InvokeServiceAckParamTemplate struct {
	Success bool `json:"success"`
}

type InvokeServiceAckTemplate struct {
	Device string                        `json:"device"`
	ID     int                           `json:"id"`
	Data   InvokeServiceAckParamTemplate `json:"data"`
}

type InvokeServiceRequestParamTemplate struct {
	ID     int                    `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type InvokeServiceRequestTemplate struct {
	Device string                            `json:"device"`
	Data   InvokeServiceRequestParamTemplate `json:"data"`
}

func (r *ReportServiceParamThingsBoardTemplate) ReportServiceThingsBoardInvokeServiceAck(reqFrame InvokeServiceRequestTemplate, ackParam InvokeServiceAckParamTemplate) {

	ackFrame := InvokeServiceAckTemplate{
		Device: reqFrame.Device,
		ID:     reqFrame.Data.ID,
		Data:   ackParam,
	}

	sJson, _ := json.Marshal(ackFrame)
	serviceInvokeTopic := "v1/gateway/rpc"

	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(serviceInvokeTopic, 0, false, sJson)
		if token.WaitTimeout(time.Duration(TimeOutService)*time.Second) == false {
			setting.ZAPS.Errorf("上报服务[%s]发布调用设备服务应答消息主题失败 %v", r.GWParam.ServiceName, token.Error())
		} else {
			setting.ZAPS.Infof("上报服务[%s]发布调用设备服务应答消息主题%v成功 上报内容%v", r.GWParam.ServiceName, serviceInvokeTopic, string(sJson))
			setting.ZAPS.Infof("上报服务[%s]发布调用设备服务应答消息内容%v", r.GWParam.ServiceName, string(sJson))
		}
	}
}

func (r *ReportServiceParamThingsBoardTemplate) ReportServiceThingsBoardProcessInvokeService(reqFrame InvokeServiceRequestTemplate) {

	ackParam := InvokeServiceAckParamTemplate{
		Success: false,
	}

	for _, node := range r.NodeList {
		if reqFrame.Device == node.Param.DeviceName {
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
					cmd.FunName = reqFrame.Data.Method
					paramStr, _ := json.Marshal(reqFrame.Data.Params)
					cmd.FunPara = string(paramStr)

					ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
					if ackData.Status {
						ackParam.Success = true
					} else {
						ackParam.Success = false
					}
				}
			}
		}
	}

	r.ReportServiceThingsBoardInvokeServiceAck(reqFrame, ackParam)
}
