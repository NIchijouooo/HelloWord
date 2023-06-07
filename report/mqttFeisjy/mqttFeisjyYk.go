package mqttFeisjy

import (
	"fmt"
	"gateway/setting"
	"os"
	"time"
)

// Yk遥控遥调控制命令处理状态机
func (r *ReportServiceParamFeisjyTemplate) FeisjyYkMachine(reqFrame MQTTFeisjyWritePropertyTemplate) bool {

	if reqFrame.DeviceAddr == r.GWParam.Param.DeviceID {
		fmt.Println(reqFrame.DeviceAddr, r.GWParam.Param.DeviceID, reqFrame.Properties)
		for _, v := range reqFrame.Properties {
			switch v.Code {
			case "3":
				{
					// QJHui UPDate 2023/6/7 value修改为interface类型后,其数值类型默认为float64
					if int(v.Value.(float64)) == 1 {
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
