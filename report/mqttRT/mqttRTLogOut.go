package mqttRT

import (
	"gateway/setting"
)

func (r *ReportServiceParamRTTemplate) NodeLogOut(name string) error {

	setting.ZAPS.Debugf("上报服务[%s]节点%s离线", r.GWParam.ServiceName, name)
	propertyPostParamMap := make([]MQTTRTPropertyPostParamTemplate, 0)
	for k, v := range r.NodeList {
		if name == v.Name {
			//传输编码为空不上报
			if v.Param.DeviceCode == "" {
				continue
			}
			//上报故障计数值先加，收到正确回应后清0
			r.NodeList[k].ReportErrCnt++
			propertyPostParam := MQTTRTPropertyPostParamTemplate{
				DeviceCode: v.Param.DeviceCode,
				OnLine:     false,
				VarData:    make([]MQTTRTPropertyPostParamPropertyTemplate, 0),
			}

			//单个设备发送
			propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)
			r.MQTTRTPropertyPost(DeviceTypeNode, r.GWParam, propertyPostParamMap[len(propertyPostParamMap)-1:], false)
		}
	}

	return nil
}
