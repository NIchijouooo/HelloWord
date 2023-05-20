package mqttSagooIOT

import (
	"gateway/setting"
)

func (r *ReportServiceParamSagooIOTTemplate) NodeLogOut(name string) error {

	setting.ZAPS.Debugf("上报服务[%s]节点%s离线", r.GWParam.ServiceName, name)
	propertyPostParamMap := make([]MQTTSagooIOTPropertyPostParamTemplate, 0)
	for k, v := range r.NodeList {
		if name == v.Name {
			//传输编码为空不上报
			if v.Param.DeviceCode == "" {
				continue
			}
			//上报故障计数值先加，收到正确回应后清0
			r.NodeList[k].ReportErrCnt++
			propertyPostParam := MQTTSagooIOTPropertyPostParamTemplate{
				DeviceCode: v.Param.DeviceCode,
				OnLine:     false,
				VarData:    make([]MQTTSagooIOTPropertyPostParamPropertyTemplate, 0),
			}

			//单个设备发送
			propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)
			r.MQTTSagooIOTPropertyPost(DeviceTypeNode, r.GWParam, propertyPostParamMap[len(propertyPostParamMap)-1:], false)
		}
	}

	return nil
}
