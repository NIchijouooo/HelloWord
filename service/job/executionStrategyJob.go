package job

import (
	"encoding/json"
	"fmt"
	"gateway/device"
	"gateway/models"
	"gateway/repositories"
	"gateway/setting"
	"strconv"
	"strings"
	"time"
)

/*
*
执行策略
*/
func ExecutionStrategy() {
	strategyRepo := repositories.NewCentralizedRepository()
	// 查询正在运行的策略列表
	strategyList, err := strategyRepo.GetRunStrategyList()
	if err != nil || len(strategyList) == 0 {
		return
	}
	layout := "2006-01-0215:04:05"
	nowTime := time.Now()
	milli := nowTime.UnixMilli()
	for _, strategy := range strategyList {
		startDate := strategy.StartDate
		endDate := strategy.EndDate
		startTime := strategy.StartTime
		endTime := strategy.EndTime
		// 参数错误
		if len(startDate) == 0 || len(endDate) == 0 || len(startTime) == 0 || len(endTime) == 0 {
			setting.ZAPS.Errorf("策略配置参数错误,id=%d", strategy.Id)
			continue
		}
		// 判断是否在执行时间内
		startTimeStr := startDate + startTime + ":00"
		endTimeStr := endDate + endTime + ":00"
		start, startErr := time.ParseInLocation(layout, startTimeStr, nowTime.Location())
		end, endErr := time.ParseInLocation(layout, endTimeStr, nowTime.Location())
		if startErr != nil || endErr != nil {
			setting.ZAPS.Errorf("策略配置时间错误,id=%d,startErr=%v,endErr=%v", strategy.Id, startErr, endErr)
			continue
		}
		startMilli := start.UnixMilli()
		endMilli := end.UnixMilli()
		// 当前时间在执行时间内,下发策略
		if startMilli <= milli && milli < endMilli {
			sendStrategyVal(strategy.ActivePower, nowTime.Hour(), nowTime.Minute())
		}
		// 当前时间等于结束时间,下发功率重置为0
		if nowTime.Unix() == end.Unix() {
			sendStrategyVal(0, nowTime.Hour(), nowTime.Minute())
		}
	}
}

/*
*
重置功率,原来策略状态为开启,改为关闭或者删除策略时需要重置功率
*/
func ResetPower(strategy models.EmStrategy) {
	startTimeStr := strategy.StartTime
	endTimeStr := strategy.EndTime
	if len(startTimeStr) == 0 || len(endTimeStr) == 0 {
		return
	}
	now := time.Now()
	// 截取小时和分钟并转换成int
	startArray := strings.Split(startTimeStr, ":")
	endArray := strings.Split(endTimeStr, ":")
	startHour, startHourErr := strconv.Atoi(startArray[0])
	endHour, endHourErr := strconv.Atoi(endArray[0])
	startMinute, startMinuteErr := strconv.Atoi(startArray[1])
	endMinute, endMinuteErr := strconv.Atoi(endArray[1])
	if startHourErr != nil || endHourErr != nil || startMinuteErr != nil || endMinuteErr != nil {
		setting.ZAPS.Errorf("策略配置时间错误,id=%d,startHourErr=%v,endHourErr=%v,startMinuteErr=%v,endMinuteErr=%v", strategy.Id, startHourErr, endHourErr, startMinuteErr, endMinuteErr)
		return
	}
	// 按当天时间设置小时分钟,计算需要下发多少次重置
	startTime := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMinute, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), endHour, endMinute, 0, 0, now.Location())
	startMilli := startTime.UnixMilli()
	endMilli := endTime.UnixMilli()
	// 每次加半小时,从开始时间到结束时间前都重置为0
	for i := startMilli; i < endMilli && i < now.UnixMilli(); i += 30 * 60 * 1000 {
		resetTime := time.UnixMilli(i)
		// 重置功率
		sendStrategyVal(0, resetTime.Hour(), resetTime.Minute())
	}

}

/*
*
下发策略配置的有功功率
activePower=有功功率数值
isReset=是否重置有功功率
hour=当前小时
minute=当前分钟
*/
func sendStrategyVal(activePower float64, hour int, minute int) {
	dictLabel := fmt.Sprintf("plc_active_power_yc_code_%v_%v", hour, minute)
	dictRepo := repositories.NewDictDataRepository()
	dictData, err := dictRepo.SelectDictValue("energy_product_code_setting", dictLabel)
	if err != nil || len(dictData.DictValue) == 0 {
		setting.ZAPS.Errorf("未查询到需要下发的点位,dictLabel=%s", dictLabel)
		return
	}
	deviceRepo := repositories.NewDeviceRepository()
	deviceList, err := deviceRepo.GetDeviceListByType(models.DeviceParam{
		DeviceType: "plcMaster",
	})
	if err != nil || len(deviceList) == 0 {
		setting.ZAPS.Errorf("获取plc主机设备失败,err=%v", err)
		return
	}
	emRepo := repositories.NewEmRepository()
	// 查询全部采集接口
	collList, err := emRepo.GetAllCollInterface()
	collMap := make(map[int]models.EmCollInterface)
	if err == nil && len(collList) > 0 {
		for _, coll := range collList {
			collMap[coll.Id] = *coll
		}
	}
	for _, dev := range deviceList {
		collId := dev.CollInterfaceId
		collInterface, collOk := collMap[collId]
		if !collOk {
			setting.ZAPS.Errorf("设备采集接口id为空,deviceId=%v", dev.Id)
			continue
		}

		cmd := device.CommunicationCmdTemplate{}
		cmd.CollInterfaceName = collInterface.Name
		cmd.DeviceName = dev.Name
		cmd.FunName = "SetVariables"
		param := make(map[string]interface{}, 0)
		param[dictData.DictValue] = activePower
		paramStr, _ := json.Marshal(param)
		cmd.FunPara = string(paramStr)

		setting.ZAPS.Infof("即将修改属性。FunName [%s] 设备名[%s] 采集接口名称[%s] 属性参数 %v", cmd.FunPara, cmd.DeviceName, cmd.CollInterfaceName, cmd.FunPara)
		coll, ok := device.CollectInterfaceMap.Coll[collInterface.Name]
		if !ok {
			setting.ZAPS.Errorf("execution strategy eer")
			return
		}
		cmdRX := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
		if cmdRX.Status == true {
			setting.ZAPS.Infof("修改属性成功。FunName [%s] 设备名[%s] 采集接口名称[%s] 属性参数 %v", cmd.FunPara, cmd.DeviceName, cmd.CollInterfaceName, cmd.FunPara)
		} else {
			setting.ZAPS.Infof("修改属性失败。FunName [%s] 设备名[%s] 采集接口名称[%s] 属性参数 %v", cmd.FunPara, cmd.DeviceName, cmd.CollInterfaceName, cmd.FunPara)
		}
	}
}
