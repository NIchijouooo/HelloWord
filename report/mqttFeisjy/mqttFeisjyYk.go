package mqttFeisjy

import (
	"gateway/setting"
	"os"
	"time"
)

// Yk遥控遥调控制命令处理状态机
func (r *ReportServiceParamFeisjyTemplate) FeisjyYkMachine(reqFrame MQTTFeisjyWritePropertyTemplate) bool {

	if reqFrame.DeviceAddr == r.GWParam.Param.DeviceID {
		for _, v := range reqFrame.Properties {
			switch v.Code {
			case "3":
				{
					if v.Value == 1 {
						setting.ZAPS.Infof("地址[%v] 网关将在0.5s后重启,程序退出码为[9]...", reqFrame.DeviceAddr)
						r.ReportServiceFeisjyWritePropertyAck(reqFrame, 1)
						time.Sleep(500 * time.Millisecond)
						os.Exit(9)

						return true
					}
				}
			}
		}
	}

	r.ReportServiceFeisjyWritePropertyAck(reqFrame, 0)

	return false
}
