package mqttFeisjy

import (
	"context"
	"gateway/device"
	"gateway/device/eventBus"
	"gateway/setting"
	"strconv"
	"time"
)

// 定义字典类型管理的存储库
type RealtimeDataSubRepository struct {
	ServiceList []*ReportServiceParamFeisjyTemplate
}

//用的是ReportServiceParamFeisjyTemplate
func NewRealtimeDataSubRepository() *ReportServiceParamFeisjyTemplate {
	feisjyParam := &ReportServiceParamFeisjyTemplate{
		ReceiveFrameChan:               make(chan MQTTFeisjyReceiveFrameTemplate, 100),
		ReceiveLogInAckFrameChan:       make(chan MQTTFeisjyLogInAckTemplate, 2),
		ReportPropertyRequestFrameChan: make(chan MQTTFeisjyReportPropertyTemplate, 50),
		ReceiveDevUpGradeChan:          make(chan MQTTFeisjyUpGradeTemplate, 2),
		ReceiveFileListChan:            make(chan FileListFeisjyTemplate, 50),
		MessageEventBus:                eventBus.NewBus(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	feisjyParam.CancelFunc = cancel
	go RealtimeDataFeisjyPoll(ctx, feisjyParam)
	return feisjyParam
}

func RealtimeDataFeisjyPoll(ctx context.Context, r *ReportServiceParamFeisjyTemplate) {
	setting.ZAPS.Debugf("RealtimeDataFeisjyPoll start......")

	//订阅采集接口消息
	device.CollectInterfaceMap.Lock.Lock()

	//device.CollectInterfaceMap.Coll {if coll.CommInterfaceName -
	//	node := coll.DeviceNodeMap["123"]
	//	node.Properties

	/**
	todo20230615网关设备上线更新sqlite
	*/
	//NewRealtimeDataRepository().UpdateGatewayDeviceConnetStatus(r.GWParam)

	for _, coll := range device.CollectInterfaceMap.Coll {
		sub := eventBus.NewSub()
		coll.CollEventBus.Subscribe("onLine", sub)
		coll.CollEventBus.Subscribe("offLine", sub)
		coll.CollEventBus.Subscribe("update", sub)
		go r.ProcessCollRealtimeData(ctx, sub)
	}
	device.CollectInterfaceMap.Lock.Unlock()

}

func (r *ReportServiceParamFeisjyTemplate) ProcessCollRealtimeData(ctx context.Context, sub eventBus.Sub) {
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case msg := <-sub.Out():
			{
				subMsg := msg.(device.CollectInterfaceEventTemplate)
				setting.ZAPS.Debugf("ProcessCollRealtimeData start......[%v]", device.CollectInterfaceMap)
				for _, collParam := range device.CollectInterfaceMap.Coll {
					if collParam.CollInterfaceName == subMsg.CollName {
						node := collParam.DeviceNodeMap[subMsg.NodeName]
						if node != nil {
							setting.ZAPS.Debugf("ProcessCollRealtimeData start......[%v]", subMsg)
							//去sqlite查数据判断设备在该上报服务中
							//var result, err = repositories.NewDeviceRepository().GetDeviceListByCollAndName(subMsg.CollName, subMsg.NodeName)
							//if result.Id <= 0{
							//	continue
							//}
							//if err != nil {
							//	continue
							//}
							switch subMsg.Topic {
								case "onLine":
									{
										//nodeName = append(nodeName, subMsg.NodeName)
										//r.NodeList[index].CommStatus = "onLine"
										////判断告警
										//r.ProcessAlarmEvent(index, subMsg.CollName, subMsg.NodeName)
									}
								case "offLine":
									{
										//nodeName = append(nodeName, subMsg.NodeName)
										//r.NodeList[index].CommStatus = "offLine"
									}
								case "update":
									setting.ZAPS.Debugf("ProcessCollRealtimeData update......[%v]", subMsg.Topic)
									//更新设备的通信状态
									//r.NodeList[index].CommStatus = "onLine"
									r.ProcessTaosRealtimeData(node, subMsg.CollName, subMsg.NodeName)
							}
						}
					}
				}
			}
		}
	}
}

func (r *ReportServiceParamFeisjyTemplate) ProcessTaosRealtimeData(devProperties *device.DeviceNodeTemplate, collName string, nodeName string) {
	setting.ZAPS.Infof("ProcessTaosRealtimeData start......[%v]", device.CollectInterfaceMap)

	//1、查找到对应的设备
	coll, ok := device.CollectInterfaceMap.Coll[collName]
	if !ok {
		return
	}
	node, ok := coll.DeviceNodeMap[nodeName]
	if !ok {
		return
	}

	//2、初始化上报结构体
	ycPropertyMap := make([]MQTTFeisjyReportDataTemplate, 0)
	ycPropertyPostParam := MQTTFeisjyReportYcTemplate{
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		CommStatus: "onLink",
	}
	//3、判断设备是否在线，不在线退出，在线判断报警上报状态
	setting.ZAPS.Infof("[%v]设备[%v]时值更新[%v]", collName, nodeName, node.CommStatus)
	for _, v := range devProperties.Properties {
		if !ok {
			continue
		}
		//3.1、判断步长报警
		valueCnt := len(v.Value)
		if valueCnt > 0 { //阶跃报警必须是2个值
			switch v.Type {
			case device.PropertyTypeInt32:
				{
					if num, err := strconv.Atoi(v.Name); err == nil {
						ycProperty := MQTTFeisjyReportDataTemplate{
							ID:    num,
							Value: v.Value[valueCnt-1].Value.(int32),
						}
						ycPropertyMap = append(ycPropertyMap, ycProperty)
					}
					continue
				}
			case device.PropertyTypeUInt32:
				{
					if num, err := strconv.Atoi(v.Name); err == nil {
						ycProperty := MQTTFeisjyReportDataTemplate{
							ID:    num,
							Value: v.Value[valueCnt-1].Value.(uint32),
						}
						ycPropertyMap = append(ycPropertyMap, ycProperty)
					}
					continue
				}
			case device.PropertyTypeDouble:
				{
					if num, err := strconv.Atoi(v.Name); err == nil {
						ycProperty := MQTTFeisjyReportDataTemplate{
							ID:    num,
							Value: v.Value[valueCnt-1].Value.(float64),
						}
						ycPropertyMap = append(ycPropertyMap, ycProperty)
					}
					continue
				}
			}
		}
	}
	ycPropertyPostParam.YcList = ycPropertyMap
	/**
	20230615实时值更新taos
	*/
	setting.ZAPS.Infof("[%v]设备[%v]时值更新taos", collName, nodeName)
	NewRealtimeDataRepository().SaveRealtimeDataList(nodeName, collName, ycPropertyPostParam)

}
