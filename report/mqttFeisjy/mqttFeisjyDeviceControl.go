package mqttFeisjy

import "gateway/setting"

// DeviceControl控制命令处理状态机
func (r *ReportServiceParamFeisjyTemplate) FeisjyDeviceControlMachine(reqFrame MQTTFeisjyWritePropertyTemplate) bool {

	switch reqFrame.CmdType {
	case "allcall":

		if reqFrame.DeviceAddr == r.GWParam.Param.DeviceID {
			reportGWProperty := MQTTFeisjyReportPropertyTemplate{
				DeviceType: "gw",
			}
			r.ReportPropertyRequestFrameChan <- reportGWProperty
		} else {
			name := make([]string, 0)
			for _, v := range r.NodeList {
				if v.Param.DeviceID == reqFrame.DeviceAddr {
					name = append(name, v.Name)
				}
			}

			reportNodeProperty := MQTTFeisjyReportPropertyTemplate{
				DeviceType: "node",
				DeviceName: name,
			}

			r.ReportPropertyRequestFrameChan <- reportNodeProperty
		}
	case "yk":
		{
			setting.ZAPS.Infof("摇控遥调下发信息 %v", reqFrame)

			r.FeisjyYkMachine(reqFrame)
		}
	case "yc":
		{
			setting.ZAPS.Infof("遥测遥调下发信息 %v", reqFrame)

			r.ReportServiceFeisjyProcessWriteProperty(reqFrame)
		}
	case "setting":
		{
			setting.ZAPS.Infof("配置参数(setting)下发信息 %v", reqFrame)

			r.ReportServiceFeisjyProcessWriteProperty(reqFrame)
		}
	}

	return true
}
