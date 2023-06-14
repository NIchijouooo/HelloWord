package mqttZxJs

import "gateway/setting"

// DeviceControl控制命令处理状态机
func (r *ReportServiceParamZxjsTemplate) ZxjsDeviceControlMachine(reqFrame MQTTZxjsWritePropertyTemplate) bool {

	switch reqFrame.CmdType {
	case "allcall":

		if reqFrame.DeviceAddr == r.GWParam.Param.DeviceSn {
			reportGWProperty := MQTTZxjsReportPropertyTemplate{
				DeviceType: "gw",
			}
			r.ReportPropertyRequestFrameChan <- reportGWProperty
		} else {
			name := make([]string, 0)
			for _, v := range r.NodeList {
				if v.Param.DeviceSn == reqFrame.DeviceAddr {
					name = append(name, v.Name)
				}
			}

			reportNodeProperty := MQTTZxjsReportPropertyTemplate{
				DeviceType: "node",
				DeviceName: name,
			}

			r.ReportPropertyRequestFrameChan <- reportNodeProperty
		}
	case "yk":
		{
			setting.ZAPS.Infof("摇控遥调下发信息 %v", reqFrame)

			r.ZxjsYkMachine(reqFrame)
		}
	case "yc":
		{
			setting.ZAPS.Infof("遥测遥调下发信息 %v", reqFrame)

			r.ReportServiceZxjsProcessWriteProperty(reqFrame)
		}
	case "setting":
		{
			setting.ZAPS.Infof("配置参数(setting)下发信息 %v", reqFrame)

			r.ReportServiceZxjsProcessWriteProperty(reqFrame)
		}
	}

	return true
}
