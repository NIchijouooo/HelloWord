package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/device/commInterface"
	"gateway/device/eventBus"
	"gateway/setting"
	"gateway/utils"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"

	"github.com/robfig/cron/v3"
)

type CollectInterfaceEventTemplate struct {
	Topic    string      `json:"Topic"`    //事件主题，online，offline，update
	CollName string      `json:"CollName"` //采集接口名称
	NodeName string      `json:"NodeName"` //设备节点名称
	Content  interface{} `json:"Content"`  //事件内容
}

// 采集接口模板
type CollectInterfaceTemplate struct {
	Index               int                                  `json:"index"`
	CollInterfaceName   string                               `json:"collInterfaceName"` //采集接口
	CommInterfaceName   string                               `json:"commInterfaceName"` //通信接口
	ProtocolTypeName    string                               `json:"protocolTypeName"`  //协议类型名称,新增加的协议通过这个变量区分协议，不在通信接口里面做了--gwai add 2023-05-05
	CommInterface       commInterface.CommunicationInterface `json:"-"`
	CommInterfaceUpdate chan bool                            `json:"-"`                   //通信接口更新
	CommQueueManage     *CommunicationManageTemplate         `json:"-"`                   //通信队列
	PollPeriod          int                                  `json:"pollPeriod"`          //采集周期
	OfflinePeriod       int                                  `json:"offlinePeriod"`       //离线超时周期
	DeviceNodeCnt       int                                  `json:"deviceNodeCnt"`       //设备数量
	DeviceNodeOnlineCnt int                                  `json:"deviceNodeOnlineCnt"` //设备在线数量
	DeviceNodeMap       map[string]*DeviceNodeTemplate       `json:"deviceNodeMap"`       //节点表
	CollEventBus        eventBus.Bus                         `json:"-"`                   //事件总线
	Cron                *cron.Cron                           `json:"-"`                   //定时管理
	CronId              cron.EntryID                         `json:"-"`                   //定时器ID 2023/5/31 QJHui ADD
	ContextCancelFun    context.CancelFunc                   `json:"-"`
	WG                  sync.WaitGroup                       `json:"-"`
	TSLLuaStateMap      map[string]*lua.LState               `json:"-"`
	TSLEventSub         eventBus.Sub                         `json:"-"`
	MessageEventBus     eventBus.Bus                         `json:"-"` //通信报文总线
}

type CollectInterfaceMapTemplate struct {
	Lock sync.RWMutex
	Coll map[string]*CollectInterfaceTemplate
}

var CollectInterfaceMap = CollectInterfaceMapTemplate{
	Lock: sync.RWMutex{},
	Coll: make(map[string]*CollectInterfaceTemplate),
}

var writeTimer *time.Timer

func CollectInterfaceInit() {

	//通信接口
	commInterface.CommInterfaceInit()
	//采集接口
	if ReadCollectInterfaceManageFromJson() == true {
		//CollectInterfaceMap.Lock.Lock()
		//for _, v := range CollectInterfaceMap.Coll {
		//	//立马进行一次采集，加快设备第一次通信速度
		//	v.CommunicationManagePoll()
		//}
		//CollectInterfaceMap.Lock.Unlock()
	}
	writeTimer = time.AfterFunc(time.Second, func() {
		WriteCollectInterfaceManageToJson()
	})
	writeTimer.Stop()
}

func WriteCollectInterfaceManageToJson() {

	//采集接口配置参数
	type ConfigParamTemplate struct {
		CollInterfaceName  string   `json:"CollInterfaceName"`  //采集接口
		CommInterfaceName  string   `json:"CommInterfaceName"`  //通信接口
		ProtocolTypeName   string   `json:"protocolTypeName"`   //协议类型名称,新增加的协议通过这个变量区分协议，不在通信接口里面做了--gwai add 2023-05-05
		PollPeriod         int      `json:"PollPeriod"`         //采集周期
		OfflinePeriod      int      `json:"OfflinePeriod"`      //离线超时周期
		DeviceNodeCnt      int      `json:"DeviceNodeCnt"`      //设备数量
		DeviceNodeNameMap  []string `json:"DeviceNodeNameMap"`  //节点名称
		DeviceNodeLabelMap []string `json:"DeviceNodeLabelMap"` //节点标签
		DeviceNodeAddrMap  []string `json:"DeviceNodeAddrMap"`  //节点地址
		DeviceNodeTypeMap  []string `json:"DeviceNodeTypeMap"`  //节点类型
	}

	utils.DirIsExist("./selfpara")

	//定义采集接口参数结构体
	configParamMap := struct {
		CollectInterfaceParam []ConfigParamTemplate
	}{
		CollectInterfaceParam: make([]ConfigParamTemplate, 0),
	}

	CollectInterfaceMap.Lock.Lock()
	for _, v := range CollectInterfaceMap.Coll {
		if v == nil {
			continue
		}
		ParamTemplate := ConfigParamTemplate{
			CollInterfaceName: v.CollInterfaceName,
			CommInterfaceName: v.CommInterfaceName,
			ProtocolTypeName:  v.ProtocolTypeName,
			PollPeriod:        v.PollPeriod,
			OfflinePeriod:     v.OfflinePeriod,
			DeviceNodeCnt:     v.DeviceNodeCnt,
		}

		ParamTemplate.DeviceNodeNameMap = make([]string, 0)
		ParamTemplate.DeviceNodeLabelMap = make([]string, 0)
		ParamTemplate.DeviceNodeAddrMap = make([]string, 0)
		ParamTemplate.DeviceNodeTypeMap = make([]string, 0)

		for _, d := range v.DeviceNodeMap {
			ParamTemplate.DeviceNodeNameMap = append(ParamTemplate.DeviceNodeNameMap, d.Name)
			ParamTemplate.DeviceNodeLabelMap = append(ParamTemplate.DeviceNodeLabelMap, d.Label)
			ParamTemplate.DeviceNodeAddrMap = append(ParamTemplate.DeviceNodeAddrMap, d.Addr)
			ParamTemplate.DeviceNodeTypeMap = append(ParamTemplate.DeviceNodeTypeMap, d.TSL)
		}

		configParamMap.CollectInterfaceParam = append(configParamMap.CollectInterfaceParam,
			ParamTemplate)
	}
	CollectInterfaceMap.Lock.Unlock()

	sJson, _ := json.Marshal(configParamMap)
	err := utils.FileWrite("./selfpara/collInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("采集接口配置json文件写入失败 %v", err)
		return
	}
	setting.ZAPS.Debug("采集接口配置json文件写入成功")
}

func ReadCollectInterfaceManageFromJson() bool {
	//采集接口配置参数
	type ConfigParamTemplate struct {
		CollInterfaceName  string   `json:"CollInterfaceName"`  //采集接口
		CommInterfaceName  string   `json:"CommInterfaceName"`  //通信接口
		ProtocolTypeName   string   `json:"ProtocolTypeName"`   //协议
		PollPeriod         int      `json:"PollPeriod"`         //采集周期
		OfflinePeriod      int      `json:"OfflinePeriod"`      //离线超时周期
		DeviceNodeCnt      int      `json:"DeviceNodeCnt"`      //设备数量
		DeviceNodeNameMap  []string `json:"DeviceNodeNameMap"`  //节点名称
		DeviceNodeLabelMap []string `json:"DeviceNodeLabelMap"` //节点标签
		DeviceNodeAddrMap  []string `json:"DeviceNodeAddrMap"`  //节点地址
		DeviceNodeTypeMap  []string `json:"DeviceNodeTypeMap"`  //节点类型
	}

	data, err := utils.FileRead("./selfpara/collInterface.json")
	if err != nil {
		setting.ZAPS.Debugf("采集接口配置json文件读取失败 %v", err)
		return false
	}
	//定义采集接口参数结构体
	configParamMap := struct {
		CollectInterfaceParam []ConfigParamTemplate
	}{
		CollectInterfaceParam: make([]ConfigParamTemplate, 0),
	}

	err = json.Unmarshal(data, &configParamMap)
	if err != nil {
		setting.ZAPS.Errorf("采集接口配置json文件格式化失败 %v", err)
		return false
	}
	setting.ZAPS.Info("采集接口配置json文件读取成功")

	for _, v := range configParamMap.CollectInterfaceParam {
		index := -1
		for k, c := range commInterface.CommunicationInterfaceMap {
			if c.GetName() == v.CommInterfaceName {
				index = k
				break
			}
		}
		if index == -1 {
			setting.ZAPS.Errorf("采集接口[%s]中通信接口[%s]不存在", v.CollInterfaceName, v.CommInterfaceName)
			continue
		}
		nodeCnt := len(v.DeviceNodeNameMap)
		//创建接口实例
		CollectInterfaceMap.Coll[v.CollInterfaceName], err = NewCollectInterface(v.CollInterfaceName,
			v.CommInterfaceName,
			v.ProtocolTypeName,
			v.PollPeriod,
			v.OfflinePeriod,
			nodeCnt)

		nodeTypeMap := make(map[string]string)

		//创建设备实例
		//setting.ZAPS.Debugf("collName %s,nodeCnt %v", v.CollInterfaceName, nodeCnt)
		for i := 0; i < nodeCnt; i++ {
			CollectInterfaceMap.Coll[v.CollInterfaceName].NewDeviceNode(
				i,
				v.DeviceNodeNameMap[i],
				v.DeviceNodeLabelMap[i],
				v.DeviceNodeTypeMap[i],
				v.DeviceNodeAddrMap[i])
			//采集接口下设备模型去重复
			nodeTypeMap[v.DeviceNodeTypeMap[i]] = v.DeviceNodeTypeMap[i]

		}

		setting.ZAPS.Infof("采集接口[%s]物模型种类[%v]", v.CollInterfaceName, nodeTypeMap)
		//将lua文件加载到虚拟机中,每个采集接口单独一组，多个采集接口就不会冲突
		for _, d := range TSLLuaMap {
			//订阅物模型的修改和删除事件
			d.Event.Subscribe("modify", CollectInterfaceMap.Coll[v.CollInterfaceName].TSLEventSub)
			d.Event.Subscribe("delete", CollectInterfaceMap.Coll[v.CollInterfaceName].TSLEventSub)
			_, ok := nodeTypeMap[d.Name]
			if !ok {
				continue
			}
			//订阅物模型的修改和删除事件
			//d.Event.Subscribe("modify", CollectInterfaceMap.Coll[v.CollInterfaceName].TSLEventSub)
			//d.Event.Subscribe("delete", CollectInterfaceMap.Coll[v.CollInterfaceName].TSLEventSub)

			err, lState := d.DeviceTSLOpenPlugin()
			if err != nil {
				setting.ZAPS.Infof("采集接口[%s]打开物模型[%s]插件失败", v.CollInterfaceName, d.Name)
				continue
			}
			setting.ZAPS.Infof("采集接口[%s]打开物模型[%s]插件成功", v.CollInterfaceName, d.Name)
			CollectInterfaceMap.Coll[v.CollInterfaceName].TSLLuaStateMap[d.Name] = lState
		}

		for _, d := range TSLModelS7Map {
			//订阅物模型的修改和删除事件
			d.Event.Subscribe("modify", CollectInterfaceMap.Coll[v.CollInterfaceName].TSLEventSub)
			d.Event.Subscribe("delete", CollectInterfaceMap.Coll[v.CollInterfaceName].TSLEventSub)
			_, ok := nodeTypeMap[d.Name]
			if !ok {
				continue
			}
		}

		for _, d := range TSLModbusMap {
			//订阅物模型的修改和删除事件
			d.Event.Subscribe("modify", CollectInterfaceMap.Coll[v.CollInterfaceName].TSLEventSub)
			d.Event.Subscribe("delete", CollectInterfaceMap.Coll[v.CollInterfaceName].TSLEventSub)
			_, ok := nodeTypeMap[d.Name]
			if !ok {
				continue
			}
		}
		for _, d := range TSLDLT6452007Map {
			//订阅物模型的修改和删除事件
			d.Event.Subscribe("modify", CollectInterfaceMap.Coll[v.CollInterfaceName].TSLEventSub)
			d.Event.Subscribe("delete", CollectInterfaceMap.Coll[v.CollInterfaceName].TSLEventSub)
			_, ok := nodeTypeMap[d.Name]
			if !ok {
				continue
			}
		}
	}

	return true
}

/*
*******************************************************
功能描述：	增加接口
参数说明：
返回说明：
调用方式：
全局变量：
读写时间：
注意事项：
日期    ：
*******************************************************
*/
func NewCollectInterface(collInterfaceName, commInterfaceName, protocolTypeName string,
	pollPeriod, offlinePeriod int, deviceNodeCnt int) (*CollectInterfaceTemplate, error) {

	index := -1
	for k, v := range commInterface.CommunicationInterfaceMap {
		if v.GetName() == commInterfaceName {
			index = k
			break
		}
	}
	if index == -1 {
		return nil, errors.New("通信接口不存在")
	}

	coll := &CollectInterfaceTemplate{
		Index:               len(CollectInterfaceMap.Coll),
		CollInterfaceName:   collInterfaceName,
		CommInterfaceName:   commInterfaceName,
		ProtocolTypeName:    protocolTypeName,
		CommInterface:       commInterface.CommunicationInterfaceMap[index],
		CommInterfaceUpdate: make(chan bool),
		CommQueueManage:     NewCommunicationManageTemplate(),
		PollPeriod:          pollPeriod,
		OfflinePeriod:       offlinePeriod,
		DeviceNodeCnt:       deviceNodeCnt,
		DeviceNodeMap:       make(map[string]*DeviceNodeTemplate),
		CollEventBus:        eventBus.NewBus(),
		Cron:                cron.New(),
		TSLLuaStateMap:      make(map[string]*lua.LState),
		TSLEventSub:         eventBus.NewSub(),
		MessageEventBus:     eventBus.NewBus(),
	}

	////将lua文件加载到虚拟机中,每个采集接口单独一组，多个采集接口就不会冲突
	//for _, v := range TSLLuaMap {
	//	//订阅物模型的修改和删除事件
	//	v.Event.Subscribe("modify", coll.TSLEventSub)
	//	v.Event.Subscribe("delete", coll.TSLEventSub)
	//
	//	setting.ZAPS.Infof("采集接口[%s]打开物模型[%s]", collInterfaceName, v.Name)
	//	err, lState := v.DeviceTSLOpenPlugin()
	//	if err != nil {
	//		setting.ZAPS.Infof("采集接口[%s]打开物模型[%s]失败", collInterfaceName, v.Name)
	//		continue
	//	}
	//	coll.TSLLuaStateMap[v.Name] = lState
	//}

	ctx, cancel := context.WithCancel(context.Background())
	coll.ContextCancelFun = cancel

	//打开通信接口
	//coll.CommInterface.Open()

	str := fmt.Sprintf("@every %dm%ds", coll.PollPeriod/60, coll.PollPeriod%60)
	setting.ZAPS.Infof("采集接口[%s]定时轮询任务开启 %dm%ds执行一次", collInterfaceName, coll.PollPeriod/60, coll.PollPeriod%60)
	//添加定时任务
	coll.CronId, _ = coll.Cron.AddFunc(str, coll.CommunicationManagePoll)
	coll.Cron.Start()
	//创建通信接口接收协程
	switch coll.CommInterface.GetType() {
	case commInterface.CommTypeSerial:
		fallthrough
	case commInterface.CommTypeTcpClient:
		fallthrough
	case commInterface.CommTypeDR504:
		fallthrough
	case commInterface.CommTypeTcpServer:
		go coll.CommQueueManage.CommunicationManageProcessReceiveData(ctx, coll.CommInterface)
	default:
	}

	//创建通信队列处理协程
	go coll.CommunicationManageDel(ctx)
	//创建通信接口更新协程
	go coll.CommInterfaceProcessUpdate(ctx)
	//创建物模型更新协程
	go coll.TSLProcessUpdate(ctx)

	return coll, nil
}

func AddCollectInterface(collName string, commName string, protocolTypeName string, PollPeriod, OfflinePeriod int) error {

	CollectInterfaceMap.Lock.Lock()
	_, ok := CollectInterfaceMap.Coll[collName]
	CollectInterfaceMap.Lock.Unlock()
	if ok {
		return errors.New("采集接口名称已经存在")
	}

	var err error
	CollectInterfaceMap.Coll[collName], err = NewCollectInterface(collName, commName, protocolTypeName, PollPeriod, OfflinePeriod, 0)
	if err != nil {
		return err
	}
	WriteCollectInterfaceManageToJson()

	return nil
}

/*
*******************************************************
功能描述：	修改接口
参数说明：
返回说明：
调用方式：
全局变量：
读写时间：
注意事项：
日期    ：
*******************************************************
*/
func ModifyCollectInterface(collName string, commName string, protocolTypeName string, pollPeriod, offlinePeriod int) error {

	coll, ok := CollectInterfaceMap.Coll[collName]
	if !ok {
		return errors.New("采集接口不存在")
	}

	//轮询周期发生了变化
	if coll.PollPeriod != pollPeriod {
		//停止
		coll.Cron.Stop()
		//重启
		coll.Cron = cron.New()
		str := fmt.Sprintf("@every %dm%ds", pollPeriod/60, pollPeriod%60)
		setting.ZAPS.Infof("采集任务[%s] %+v", coll.CollInterfaceName, str)
		//添加定时任务
		coll.CronId, _ = coll.Cron.AddFunc(str, coll.CommunicationManagePoll)
		coll.Cron.Start()

		// 2023/5/31 QJHui ADD
		coll.PollPeriod = pollPeriod
	}

	//通信接口发生了变化
	if coll.CommInterfaceName != commName {
		//关闭旧的通信接口
		coll.CommInterface.Close()

		//相关协程退出
		coll.ContextCancelFun()

		index := -1
		for k, v := range commInterface.CommunicationInterfaceMap {
			if v.GetName() == commName {
				index = k
				break
			}
		}
		if index == -1 {
			return nil
		}
		coll.CommInterface = commInterface.CommunicationInterfaceMap[index]
		ctx, cancel := context.WithCancel(context.Background())
		coll.ContextCancelFun = cancel

		//打开通信接口
		coll.CommInterface.Open()

		//创建通信接口接收协程
		go coll.CommQueueManage.CommunicationManageProcessReceiveData(ctx, coll.CommInterface)
		//创建通信队列处理协程
		go coll.CommunicationManageDel(ctx)
		//创建通信接口更新协程
		go coll.CommInterfaceProcessUpdate(ctx)
		//创建物模型更新协程
		go coll.TSLProcessUpdate(ctx)
	}

	coll.CommInterfaceName = commName
	coll.ProtocolTypeName = protocolTypeName //GWAI ADD 2023-05-10
	coll.PollPeriod = pollPeriod
	coll.OfflinePeriod = offlinePeriod
	WriteCollectInterfaceManageToJson()

	return nil
}

/*
*******************************************************
功能描述：	删除接口
参数说明：
返回说明：
调用方式：
全局变量：
读写时间：
注意事项：
日期    ：
*******************************************************
*/
func DeleteCollectInterface(collName string) error {

	CollectInterfaceMap.Lock.Lock()
	coll, ok := CollectInterfaceMap.Coll[collName]
	CollectInterfaceMap.Lock.Unlock()
	if !ok {
		return errors.New("采集接口不存在")
	}

	if len(coll.DeviceNodeMap) > 0 {
		return errors.New("采集接口已添加设备，请先删除设备")
	}
	coll.Cron.Stop()

	//相关协程退出
	coll.ContextCancelFun()

	delete(CollectInterfaceMap.Coll, collName)

	WriteCollectInterfaceManageToJson()

	return nil
}

/*
*******************************************************
功能描述：	增加单个节点
参数说明：
返回说明：
调用方式：
全局变量：
读写时间：
注意事项：
日期    ：
*******************************************************
*/
func (d *CollectInterfaceTemplate) NewDeviceNode(index int, dName string, dLabel string, dTSL string, dAddr string) {

	node := &DeviceNodeTemplate{
		Index:          index,
		Name:           dName,
		Label:          dLabel,
		Addr:           dAddr,
		TSL:            dTSL,
		LastCommRTC:    "1970-01-01 00:00:00",
		CommTotalCnt:   0,
		CommSuccessCnt: 0,
		CurCommFailCnt: 0,
		CommStatus:     "offLine",
	}

	properties := node.NewVariables()
	node.Properties = append(node.Properties, properties...)
	services := node.NewServices()
	node.Services = append(node.Services, services...)

	d.DeviceNodeMap[dName] = node
}

func (d *CollectInterfaceTemplate) AddDeviceNode(dName string, dTSL string, dAddr string, dLabel string) error {

	_, ok := d.TSLLuaStateMap[dTSL]
	//物模型在采集中不存在
	if !ok {
		//将lua文件加载到虚拟机中
		for _, v := range TSLLuaMap {
			//订阅物模型的修改和删除事件
			v.Event.Subscribe("modify", d.TSLEventSub)
			v.Event.Subscribe("delete", d.TSLEventSub)

			err, lState := v.DeviceTSLOpenPlugin()
			if err != nil {
				setting.ZAPS.Infof("采集接口[%s]打开物模型[%s]失败", d.CollInterfaceName, v.Name)
				continue
			}
			setting.ZAPS.Infof("采集接口[%s]打开物模型[%s]成功", d.CollInterfaceName, v.Name)
			d.TSLLuaStateMap[v.Name] = lState
		}
	}

	node := &DeviceNodeTemplate{}
	node.Index = len(d.DeviceNodeMap)
	node.Name = dName
	node.Addr = dAddr
	node.TSL = dTSL
	node.Label = dLabel
	node.LastCommRTC = "1970-01-01 00:00:00"
	node.CommTotalCnt = 0
	node.CommSuccessCnt = 0
	node.CurCommFailCnt = 0
	node.CommStatus = "offLine"
	//node.VariableMap = make([]VariableTemplate, 0)
	//variables := node.NewVariables()
	//node.VariableMap = append(node.VariableMap, variables...)

	properties := node.NewVariables()
	node.Properties = append(node.Properties, properties...)
	services := node.NewServices()
	node.Services = append(node.Services, services...)

	d.DeviceNodeMap[dName] = node

	d.DeviceNodeCnt++

	writeTimer.Reset(time.Second)

	return nil
}

func (d *CollectInterfaceTemplate) ModifyDeviceNode(dName string, dTSL string, dAddr string, dLabel string) error {

	_, ok := d.TSLLuaStateMap[dTSL]
	//物模型在采集中不存在
	if !ok {
		//将lua文件加载到虚拟机中
		for _, v := range TSLLuaMap {
			//订阅物模型的修改和删除事件
			v.Event.Subscribe("modify", d.TSLEventSub)
			v.Event.Subscribe("delete", d.TSLEventSub)

			err, lState := v.DeviceTSLOpenPlugin()
			if err != nil {
				setting.ZAPS.Infof("采集接口[%s]打开物模型[%s]失败", d.CollInterfaceName, v.Name)
				continue
			}
			setting.ZAPS.Infof("采集接口[%s]打开物模型[%s]成功", d.CollInterfaceName, v.Name)
			d.TSLLuaStateMap[v.Name] = lState
		}
	}

	node, ok := d.DeviceNodeMap[dName]
	if !ok {
		return errors.New("设备名称不存在")
	}

	node.Addr = dAddr
	node.TSL = dTSL
	node.Label = dLabel

	writeTimer.Reset(time.Second)

	return nil
}

func (d *CollectInterfaceTemplate) DeleteDeviceNodes(deviceNames []string) {

	for _, v := range deviceNames {
		_, ok := d.DeviceNodeMap[v]
		if ok {
			d.DeviceNodeCnt--
			delete(d.DeviceNodeMap, v)
		}
	}

	writeTimer.Reset(time.Second)
}

func (d *CollectInterfaceTemplate) GetDeviceNode(dAddr string) *DeviceNodeTemplate {

	for _, v := range d.DeviceNodeMap {
		if v.Addr == dAddr {
			return v
		}
	}

	return nil
}

func (d *CollectInterfaceTemplate) CommInterfaceProcessUpdate(ctx context.Context) {
	setting.ZAPS.Debugf("采集接口[%s]处理更新协程3/4进入", d.CollInterfaceName)
	for {
		select {
		case <-ctx.Done():
			setting.ZAPS.Debugf("采集接口[%s]处理更新协程3/4退出", d.CollInterfaceName)
			return
		case <-d.CommInterfaceUpdate:
			{
				setting.ZAPS.Debugf("通信接口[%s]发生更新事件", d.CommInterfaceName)
				//旧采集队列接收数据协程退出
				d.CommQueueManage.QuitChan <- true
				rt := d.CommInterface.Open()
				if rt != true {
					setting.ZAPS.Debugf("通信接口[%s]重新打开失败", d.CommInterfaceName)
				} else {
					setting.ZAPS.Debugf("通信接口[%s]重新打开成功", d.CommInterfaceName)
				}
				time.Sleep(100 * time.Millisecond)
				//创建新通信接口接收协程
				setting.ZAPS.Debugf("采集接口[%s]处理更新协程重新打开", d.CollInterfaceName)
				if d.CommInterface.GetType() != commInterface.CommTypeIoIn {
					go d.CommQueueManage.CommunicationManageProcessReceiveData(ctx, d.CommInterface)
				}
			}
		}
	}
}

func (d *CollectInterfaceTemplate) TSLProcessUpdate(ctx context.Context) {
	setting.ZAPS.Debugf("采集接口[%s]处理物模型更新协程4/4进入", d.CollInterfaceName)
	for {
		select {
		case <-ctx.Done():
			setting.ZAPS.Debugf("采集接口[%s]处理物模型更新协程4/4退出", d.CollInterfaceName)
			return
		case msg := <-d.TSLEventSub.Out():
			{
				eventMsg := msg.(TSLEventTemplate)
				setting.ZAPS.Debugf("采集接口[%s]处理物模型更新消息 %v", d.CollInterfaceName, msg)
				nodeNames := make([]string, 0)
				for k, v := range d.DeviceNodeMap {
					if v.TSL == eventMsg.TSL {
						nodeNames = append(nodeNames, k)
					}
				}
				setting.ZAPS.Debugf("采集接口[%s] 物模型[%s] 发生事件[%s]", d.CollInterfaceName, eventMsg.TSL, eventMsg.Topic)
				for _, nodeName := range nodeNames {
					switch eventMsg.Topic {
					case "modify":
						{
							for _, v := range TSLLuaMap {
								if v.Name == eventMsg.TSL {
									_, d.TSLLuaStateMap[eventMsg.TSL] = v.DeviceTSLOpenPlugin()
								}
							}
							d.DeviceNodeMap[nodeName].Properties = d.DeviceNodeMap[nodeName].Properties[0:0]
							properties := d.DeviceNodeMap[nodeName].NewVariables()
							d.DeviceNodeMap[nodeName].Properties = append(d.DeviceNodeMap[nodeName].Properties, properties...)
							d.DeviceNodeMap[nodeName].Services = d.DeviceNodeMap[nodeName].Services[0:0]
							d.DeviceNodeMap[nodeName].NewServices()
						}
					case "delete":
						{
							d.DeviceNodeMap[nodeName].Properties = d.DeviceNodeMap[nodeName].Properties[0:0]
							d.DeviceNodeMap[nodeName].Services = d.DeviceNodeMap[nodeName].Services[0:0]
						}
					}
				}
			}
		default:
			{
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (d *CollectInterfaceTemplate) CommunicationManageDel(ctx context.Context) {
	setting.ZAPS.Debugf("采集接口[%s]命令对列监听协程2/4进入", d.CollInterfaceName)
	for {
		select {
		case <-ctx.Done():
			setting.ZAPS.Debugf("采集接口[%s]命令对列监听协程2/4退出", d.CollInterfaceName)
			return
		case cmd := <-d.CommQueueManage.EmergencyRequestChan:
			{
				setting.ZAPS.Infof("采集接口[%s]处理紧急任务 节点名称[%v] 命令函数名词[%v]", d.CollInterfaceName, cmd.DeviceName, cmd.FunName)
				rxResult := ReceiveDataTemplate{}
				node, ok := d.DeviceNodeMap[cmd.DeviceName]
				if !ok {
					setting.ZAPS.Warnf("采集接口[%s]处理紧急任务失败 节点名称[%v]不存在", d.CollInterfaceName, cmd.DeviceName)
					d.CommQueueManage.EmergencyAckChan <- rxResult
					continue
				}

				if d.ProtocolTypeName == "DLT645-2007" {
					rxResult = d.CommQueueManage.CommunicationStateMachineDLT64507(cmd,
						d.CollInterfaceName,
						d.CommInterface,
						node,
						&d.CollEventBus,
						d.TSLLuaStateMap,
						d.OfflinePeriod)
				} else if d.CommInterface.GetType() == commInterface.CommTypeIoIn {
					rxResult = d.CommQueueManage.CommunicationStateMachineIoIn(cmd,
						d.CollInterfaceName,
						d.CommInterface,
						node,
						&d.CollEventBus,
						d.TSLLuaStateMap,
						d.OfflinePeriod)
				} else if d.CommInterface.GetType() == commInterface.CommTypeIoOut {
					rxResult = d.CommQueueManage.CommunicationStateMachineIoOut(cmd,
						d.CollInterfaceName,
						d.CommInterface,
						node,
						&d.CollEventBus,
						d.TSLLuaStateMap,
						d.OfflinePeriod)
				} else if (d.CommInterface.GetType() == commInterface.CommTypeHTTPSmartNode) ||
					(d.CommInterface.GetType() == commInterface.CommTypeHTTPTianGang) {
					rxResult = d.CommQueueManage.CommunicationStateMachineHTTP(cmd,
						d.CollInterfaceName,
						d.CommInterface,
						node,
						&d.CollEventBus,
						d.TSLLuaStateMap,
						d.OfflinePeriod)
				} else if d.CommInterface.GetType() == commInterface.CommTypeS7 {
					rxResult = d.CommQueueManage.CommunicationStateMachineS7(cmd,
						d.CollInterfaceName,
						d.CommInterface,
						node,
						&d.CollEventBus,
						d.TSLLuaStateMap,
						d.OfflinePeriod)
				} else if d.CommInterface.GetType() == commInterface.CommTypeSAC009 {
					rxResult = d.CommQueueManage.CommunicationStateMachineSAC009(cmd,
						d.CollInterfaceName,
						d.CommInterface,
						node,
						&d.CollEventBus,
						d.TSLLuaStateMap,
						d.OfflinePeriod)
				} else if d.CommInterface.GetType() == commInterface.CommTypeDR504 {
					rxResult = d.CommQueueManage.CommunicationStateMachineDR504(cmd,
						d.CollInterfaceName,
						d.CommInterface,
						node,
						&d.CollEventBus,
						d.TSLLuaStateMap,
						d.OfflinePeriod)
				} else if d.CommInterface.GetType() == commInterface.CommTypeModbusTCP || d.CommInterface.GetType() == commInterface.CommTypeModbusRTU {
					rxResult = d.CommQueueManage.CommunicationStateMachineModbusTCP(cmd,
						d.CollInterfaceName,
						d.CommInterface,
						node,
						&d.CollEventBus,
						d.TSLLuaStateMap,
						d.OfflinePeriod)
				} else if d.CommInterface.GetType() == commInterface.CommTypeHDKJ {
					rxResult = d.CommQueueManage.CommunicationStateMachineHDKJ(cmd,
						d.CollInterfaceName,
						d.CommInterface,
						node,
						&d.CollEventBus,
						d.TSLLuaStateMap,
						d.OfflinePeriod)
				} else {
					rxResult = d.CommQueueManage.CommunicationStateMachine(cmd,
						d.CollInterfaceName,
						d.CommInterface,
						node,
						&d.CollEventBus,
						d.TSLLuaStateMap,
						d.OfflinePeriod)
				}
				//更新设备在线数量
				d.DeviceNodeOnlineCnt = 0
				for _, v := range d.DeviceNodeMap {
					if v.CommStatus == "onLine" {
						d.DeviceNodeOnlineCnt++
					}
				}

				setting.ZAPS.Debugf("采集接口[%s]处理紧急任务完成，处理结果%v", d.CollInterfaceName, rxResult.Status)

				d.CommQueueManage.EmergencyAckChan <- rxResult
			}
		default:
			{
				select {
				case req := <-d.CommQueueManage.DirectDataRequestChan:
					{
						ack := d.CommQueueManage.CommunicationDirectDataStateMachine(req, d.CommInterface)
						d.CommQueueManage.DirectDataAckChan <- ack
					}
				case cmd := <-d.CommQueueManage.CommonRequestChan:
					{
						node, ok := d.DeviceNodeMap[cmd.DeviceName]
						if !ok {
							continue
						}

						if d.CommInterface.GetType() != commInterface.CommTypeIoIn && d.CommInterface.GetType() != commInterface.CommTypeIoOut {
							//							setting.ZAPS.Debugf("采集接口[%s]通信队列剩余节点数%d", d.CollInterfaceName, len(d.CommQueueManage.CommonRequestChan))
						}

						if d.ProtocolTypeName == "DLT645-2007" {
							d.CommQueueManage.CommunicationStateMachineDLT64507(cmd,
								d.CollInterfaceName,
								d.CommInterface,
								node,
								&d.CollEventBus,
								d.TSLLuaStateMap,
								d.OfflinePeriod)

						} else if d.CommInterface.GetType() == commInterface.CommTypeIoIn {
							d.CommQueueManage.CommunicationStateMachineIoIn(cmd,
								d.CollInterfaceName,
								d.CommInterface,
								node,
								&d.CollEventBus,
								d.TSLLuaStateMap,
								d.OfflinePeriod)
						} else if d.CommInterface.GetType() == commInterface.CommTypeIoOut {
							d.CommQueueManage.CommunicationStateMachineIoOut(cmd,
								d.CollInterfaceName,
								d.CommInterface,
								node,
								&d.CollEventBus,
								d.TSLLuaStateMap,
								d.OfflinePeriod)
						} else if (d.CommInterface.GetType() == commInterface.CommTypeHTTPSmartNode) ||
							(d.CommInterface.GetType() == commInterface.CommTypeHTTPTianGang) {
							d.CommQueueManage.CommunicationStateMachineHTTP(cmd,
								d.CollInterfaceName,
								d.CommInterface,
								node,
								&d.CollEventBus,
								d.TSLLuaStateMap,
								d.OfflinePeriod)
						} else if d.CommInterface.GetType() == commInterface.CommTypeS7 {
							d.CommQueueManage.CommunicationStateMachineS7(cmd,
								d.CollInterfaceName,
								d.CommInterface,
								node,
								&d.CollEventBus,
								d.TSLLuaStateMap,
								d.OfflinePeriod)
						} else if d.CommInterface.GetType() == commInterface.CommTypeSAC009 {
							d.CommQueueManage.CommunicationStateMachineSAC009(cmd,
								d.CollInterfaceName,
								d.CommInterface,
								node,
								&d.CollEventBus,
								d.TSLLuaStateMap,
								d.OfflinePeriod)
						} else if d.CommInterface.GetType() == commInterface.CommTypeDR504 {
							d.CommQueueManage.CommunicationStateMachineDR504(cmd,
								d.CollInterfaceName,
								d.CommInterface,
								node,
								&d.CollEventBus,
								d.TSLLuaStateMap,
								d.OfflinePeriod)
						} else if d.CommInterface.GetType() == commInterface.CommTypeModbusTCP || d.CommInterface.GetType() == commInterface.CommTypeModbusRTU {
							d.CommQueueManage.CommunicationStateMachineModbusTCP(cmd,
								d.CollInterfaceName,
								d.CommInterface,
								node,
								&d.CollEventBus,
								d.TSLLuaStateMap,
								d.OfflinePeriod)
						} else if d.CommInterface.GetType() == commInterface.CommTypeHDKJ {
							d.CommQueueManage.CommunicationStateMachineHDKJ(cmd,
								d.CollInterfaceName,
								d.CommInterface,
								node,
								&d.CollEventBus,
								d.TSLLuaStateMap,
								d.OfflinePeriod)
						} else {
							d.CommQueueManage.CommunicationStateMachine(cmd,
								d.CollInterfaceName,
								d.CommInterface,
								node,
								&d.CollEventBus,
								d.TSLLuaStateMap,
								d.OfflinePeriod)
						}

						//更新设备在线数量,当本次采集最后一个设备时进行更新
						if len(d.CommQueueManage.CommonRequestChan) == 0 {
							d.DeviceNodeOnlineCnt = 0
							for _, v := range d.DeviceNodeMap {
								if v.CommStatus == "onLine" {
									d.DeviceNodeOnlineCnt++
								}
							}
						}
					}
				default:

					// 2023/5/31 QJHui ADD 解决设备采集一段时间后不再采集问题
					if d.Cron != nil {
						entry := d.Cron.Entry(d.CronId)
						// 判断当前条目是否在运行
						lastTime := entry.Prev
						nextTime := entry.Next
						if nextTime.IsZero() || nextTime.Before(lastTime) {
							setting.ZAPS.Errorf("检测到采集接口[%s]定时未运行,将重启采集定时器", d.CollInterfaceName)

							//停止
							d.Cron.Stop()
							//重启
							d.Cron = cron.New()

							str := fmt.Sprintf("@every %dm%ds", d.PollPeriod/60, d.PollPeriod%60)
							setting.ZAPS.Infof("采集任务[%s] %+v", d.CollInterfaceName, str)
							//添加定时任务
							d.CronId, _ = d.Cron.AddFunc(str, d.CommunicationManagePoll)
							d.Cron.Start()
						}
					}

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}
}

func (d *CollectInterfaceTemplate) CommunicationManagePoll() {

	cmd := CommunicationCmdTemplate{}
	for _, v := range d.DeviceNodeMap {
		cmd.CollInterfaceName = d.CollInterfaceName
		cmd.DeviceName = v.Name
		cmd.FunName = "GetDeviceRealVariables"
		d.CommQueueManage.CommunicationManageAddCommon(cmd)
	}
	if d.CommInterface.GetType() == commInterface.CommTypeIoIn || d.CommInterface.GetType() == commInterface.CommTypeIoOut {
		return
	}
	// setting.ZAPS.Debugf("采集接口[%s]通信对列节点总数为%d", d.CollInterfaceName, len(d.CommQueueManage.CommonRequestChan))
}

func (d *CollectInterfaceTemplate) ReadDeviceRealVariable(name string) (error, []TSLPropertiesTemplate) {

	//发送命令到响应的采集接口
	cmd := CommunicationCmdTemplate{}
	cmd.CollInterfaceName = d.CollInterfaceName
	cmd.DeviceName = d.DeviceNodeMap[name].Name
	cmd.FunName = "GetRealVariables"
	cmd.FunPara = ""
	cmdRX := d.CommQueueManage.CommunicationManageAddEmergency(cmd)
	if cmdRX.Status == true {
		return nil, d.DeviceNodeMap[name].Properties
	}

	return errors.New("读取错误"), make([]TSLPropertiesTemplate, 0)
}
