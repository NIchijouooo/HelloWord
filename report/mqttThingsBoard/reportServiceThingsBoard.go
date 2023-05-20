package mqttThingsBoard

import (
	"context"
	"encoding/json"
	"errors"
	"gateway/device"
	"gateway/device/eventBus"
	"gateway/report/reportModel"
	"gateway/setting"
	"gateway/utils"
	"gateway/virtual"
	"github.com/jasonlvhit/gocron"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//上报节点参数结构体
type ReportServiceNodeParamThingsBoardTemplate struct {
	Index             int                                                 `json:"index"`
	ServiceName       string                                              `json:"serviceName"`
	CollInterfaceName string                                              `json:"collInterfaceName"`
	Name              string                                              `json:"deviceName"`
	Label             string                                              `json:"deviceLabel"`
	Addr              string                                              `json:"deviceAddr"`
	CommStatus        string                                              `json:"commStatus"`
	ReportErrCnt      int                                                 `json:"-"`
	ReportStatus      string                                              `json:"reportStatus"`
	UploadModel       string                                              `json:"uploadModel"`
	Protocol          string                                              `json:"protocol"`
	Properties        map[string]*reportModel.ReportModelPropertyTemplate `json:"properties"`
	Param             struct {
		DeviceName string `json:"deviceName"`
	}
}

//上报网关参数结构体
type ReportServiceGWParamThingsBoardTemplate struct {
	Index        int    `json:"index"`
	ServiceName  string `json:"serviceName"`
	IP           string `json:"ip"`
	Port         string `json:"port"`
	ReportStatus string `json:"reportStatus"`
	ReportTime   int    `json:"reportTime"`
	ReportErrCnt int    `json:"reportErrCnt"`
	Protocol     string `json:"protocol"`
	Param        struct {
		UserName     string `json:"userName"`
		Password     string `json:"password"`
		ClientID     string `json:"clientID"`
		KeepAlive    string `json:"keepAlive"`
		CleanSession bool   `json:"cleanSession"`
	}
	MQTTClient MQTT.Client `json:"-"`
}

//上报服务参数，网关参数，节点参数
type ReportServiceParamThingsBoardTemplate struct {
	Index                 int `json:"index"`
	GWParam               ReportServiceGWParamThingsBoardTemplate
	NodeList              []ReportServiceNodeParamThingsBoardTemplate
	ReceiveFrameChan      chan MQTTThingsBoardReceiveFrameTemplate `json:"-"`
	LogInRequestFrameChan chan []string                            `json:"-"` //上线
	//ReceiveLogInAckFrameChan             chan MQTTThingsBoardLogInAckTemplate             `json:"-"`
	LogOutRequestFrameChan chan []string `json:"-"`
	//ReceiveLogOutAckFrameChan            chan MQTTThingsBoardLogOutAckTemplate            `json:"-"`
	ReportPropertyRequestFrameChan chan MQTTThingsBoardReportPropertyTemplate `json:"-"`
	//ReceiveReportPropertyAckFrameChan    chan MQTTThingsBoardReportPropertyAckTemplate    `json:"-"`
	ReceiveInvokeServiceRequestFrameChan chan InvokeServiceRequestTemplate `json:"-"`
	ReceiveInvokeServiceAckFrameChan     chan InvokeServiceAckTemplate     `json:"-"`
	//ReceiveWritePropertyRequestFrameChan chan MQTTThingsBoardWritePropertyRequestTemplate `json:"-"`
	//ReceiveReadPropertyRequestFrameChan  chan MQTTThingsBoardReadPropertyRequestTemplate  `json:"-"`
	CancelFunc context.CancelFunc `json:"-"`
	EventSub   eventBus.Sub       `json:"-"`
}

type ReportServiceParamListThingsBoardTemplate struct {
	ServiceList []*ReportServiceParamThingsBoardTemplate
}

const (
	TimeOutLogin          int = 60
	TimeOutLogout         int = 5
	TimeOutSubscribe      int = 5
	TimeOutReportProperty int = 5
	TimeOutReadProperty   int = 5
	TimeOutWriteProperty  int = 5
	TimeOutService        int = 5
)

//实例化上报服务
var ReportServiceParamListThingsBoard = &ReportServiceParamListThingsBoardTemplate{
	ServiceList: make([]*ReportServiceParamThingsBoardTemplate, 0),
}

var writeTimer *time.Timer

func (s *ReportServiceParamListThingsBoardTemplate) ReadParamFromJson() bool {

	configParam := ReportServiceParamListThingsBoardTemplate{}

	data, err := utils.FileRead("./selfpara/reportServiceParamListThingsBoard.json")
	if err != nil {
		setting.ZAPS.Debugf("上报服务[ThingsBoard]配置json文件读取失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &configParam)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[ThingsBoard]配置json文件格式化失败")
		return false
	}
	setting.ZAPS.Info("上报服务[ThingsBoard]配置json文件读取成功")

	//初始化
	for _, v := range configParam.ServiceList {
		//v.ReceiveFrameChan = make(chan MQTTThingsBoardReceiveFrameTemplate, 100)
		//v.LogInRequestFrameChan = make(chan []string, 0)
		//v.ReceiveLogInAckFrameChan = make(chan MQTTThingsBoardLogInAckTemplate, 5)
		//v.LogOutRequestFrameChan = make(chan []string, 0)
		//v.ReceiveLogOutAckFrameChan = make(chan MQTTThingsBoardLogOutAckTemplate, 5)
		//v.ReportPropertyRequestFrameChan = make(chan MQTTThingsBoardReportPropertyTemplate, 50)
		//v.ReceiveReportPropertyAckFrameChan = make(chan MQTTThingsBoardReportPropertyAckTemplate, 50)
		//v.ReceiveInvokeServiceRequestFrameChan = make(chan InvokeServiceRequestTemplate, 50)
		//v.ReceiveInvokeServiceAckFrameChan = make(chan InvokeServiceAckTemplate, 50)
		//v.ReceiveWritePropertyRequestFrameChan = make(chan MQTTThingsBoardWritePropertyRequestTemplate, 50)
		//v.ReceiveReadPropertyRequestFrameChan = make(chan MQTTThingsBoardReadPropertyRequestTemplate, 50)

		v.GWParam.ReportStatus = "offLine"
		for _, n := range v.NodeList {
			n.CommStatus = "offLine"
			n.ReportStatus = "offLine"
		}
		ReportServiceParamListThingsBoard.ServiceList = append(ReportServiceParamListThingsBoard.ServiceList, NewReportServiceParamThingsBoard(v.GWParam, v.NodeList))
	}
	return true
}

func (s *ReportServiceParamListThingsBoardTemplate) WriteParamToJson() {
	utils.DirIsExist("./selfpara")
	sJson, _ := json.Marshal(*s)
	err := utils.FileWrite("./selfpara/reportServiceParamListThingsBoard.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[ThingsBoard]配置json文件写入失败 %v", err)
		return
	}
	setting.ZAPS.Info("上报服务[ThingsBoard]配置json文件写入成功")
}

func (s *ReportServiceParamListThingsBoardTemplate) AddReportService(param ReportServiceGWParamThingsBoardTemplate) error {

	index := -1
	for k, v := range s.ServiceList {
		//存在相同的，表示修改;不存在表示增加
		if v.GWParam.ServiceName == param.ServiceName {
			index = k
			break
		}
	}

	if index != -1 {
		return errors.New("服务名称已经存在")
	}

	nodeList := make([]ReportServiceNodeParamThingsBoardTemplate, 0)
	ReportServiceParam := NewReportServiceParamThingsBoard(param, nodeList)
	s.ServiceList = append(s.ServiceList, ReportServiceParam)

	s.WriteParamToJson()

	return nil
}

func (s *ReportServiceParamListThingsBoardTemplate) ModifyReportService(param ReportServiceGWParamThingsBoardTemplate) error {

	index := -1
	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == param.ServiceName {
			index = k
			break
		}
	}

	//不存在表示增加
	if index == -1 {
		return errors.New("服务名称不存在")
	}

	//旧上报服务退出
	s.ServiceList[index].CancelFunc()
	s.ServiceList[index].GWParam = param

	s.WriteParamToJson()

	//启动新上报服务
	s.ServiceList[index] = NewReportServiceParamThingsBoard(param, s.ServiceList[index].NodeList)

	return nil
}

func (s *ReportServiceParamListThingsBoardTemplate) DeleteReportService(serviceName string) {

	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == serviceName {
			s.ServiceList = append(s.ServiceList[:k], s.ServiceList[k+1:]...)
			s.WriteParamToJson()
			//旧协程退出
			v.CancelFunc()
			return
		}
	}
}

func (r *ReportServiceParamThingsBoardTemplate) AddReportNode(param ReportServiceNodeParamThingsBoardTemplate) {

	param.Index = len(r.NodeList)
	param.CommStatus = "offLine"
	param.ReportStatus = "offLine"
	param.ReportErrCnt = 0

	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == param.Name {
			r.NodeList[k] = param
			ReportServiceParamListThingsBoard.WriteParamToJson()
			return
		}
	}

	//节点不存在则新建
	r.NodeList = append(r.NodeList, param)
	ReportServiceParamListThingsBoard.WriteParamToJson()

	//setting.ZAPS.Debugf("param %v", ReportServiceParamListThingsBoard)
}

func (r *ReportServiceParamThingsBoardTemplate) DeleteReportNode(name string) int {

	index := -1
	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == name {
			index = k
			r.NodeList = append(r.NodeList[:k], r.NodeList[k+1:]...)
			ReportServiceParamListThingsBoard.WriteParamToJson()
			return index
		}
	}
	return index
}

func (r *ReportServiceParamThingsBoardTemplate) ProcessUpLinkFrame(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case reqFrame := <-r.LogInRequestFrameChan:
			{
				r.LogIn(reqFrame)
			}
		case reqFrame := <-r.LogOutRequestFrameChan:
			{
				r.LogOut(reqFrame)
			}
		case reqFrame := <-r.ReportPropertyRequestFrameChan:
			{
				if reqFrame.DeviceType == "gw" {
					r.GWPropertyPost()
				} else if reqFrame.DeviceType == "node" {
					r.NodePropertyPost(reqFrame.DeviceName)
				}
			}
		//case reqFrame := <-r.ReceiveWritePropertyRequestFrameChan:
		//	{
		//		r.ReportServiceThingsBoardProcessWriteProperty(reqFrame)
		//	}
		//case reqFrame := <-r.ReceiveReadPropertyRequestFrameChan:
		//	{
		//		r.ReportServiceThingsBoardProcessReadProperty(reqFrame)
		//	}
		case reqFrame := <-r.ReceiveInvokeServiceRequestFrameChan:
			{
				r.ReportServiceThingsBoardProcessInvokeService(reqFrame)
			}
		}
	}
}

func (r *ReportServiceParamThingsBoardTemplate) ProcessDownLinkFrame(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case frame := <-r.ReceiveFrameChan:
			{
				if strings.Contains(frame.Topic, "v1/gateway/rpc") { //设备服务调用
					serviceFrame := InvokeServiceRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &serviceFrame)
					if err != nil {
						setting.ZAPS.Errorf("上报服务[%s]接收调用设备服务命令JSON格式化错误 %v", err)
						continue
					}
					r.ReceiveInvokeServiceRequestFrameChan <- serviceFrame
				}
				//} else if strings.Contains(frame.Topic, "/sys/thing/event/property/set") { //设置属性请求
				//	writePropertyFrame := MQTTThingsBoardWritePropertyRequestTemplate{}
				//	err := json.Unmarshal(frame.Payload, &writePropertyFrame)
				//	if err != nil {
				//		setting.ZAPS.Errorf("writePropertyFrame json unmarshal err")
				//		continue
				//	}
				//	r.ReceiveWritePropertyRequestFrameChan <- writePropertyFrame
				//} else if strings.Contains(frame.Topic, "/sys/thing/event/property/get") { //获取属性请求
				//	readPropertyFrame := MQTTThingsBoardReadPropertyRequestTemplate{}
				//	err := json.Unmarshal(frame.Payload, &readPropertyFrame)
				//	if err != nil {
				//		setting.ZAPS.Errorf("readPropertyFrame json unmarshal err")
				//		continue
				//	}
				//	r.ReceiveReadPropertyRequestFrameChan <- readPropertyFrame
				//}
			}
		}
	}
}

func (r *ReportServiceParamThingsBoardTemplate) ProcessCollEvent(ctx context.Context, sub eventBus.Sub) {
	for {
		select {
		case <-ctx.Done():
			{
				//取消订阅采集接口消息
				device.CollectInterfaceMap.Lock.Lock()
				for _, coll := range device.CollectInterfaceMap.Coll {
					setting.ZAPS.Infof("上报服务[%s]取消订阅采集接口[%v]", r.GWParam.ServiceName, coll.CollInterfaceName)
					coll.CollEventBus.UnSubscribe("onLine", r.EventSub)
					coll.CollEventBus.UnSubscribe("offLine", r.EventSub)
					coll.CollEventBus.UnSubscribe("update", r.EventSub)
				}
				device.CollectInterfaceMap.Lock.Unlock()
				return
			}
		case msg := <-sub.Out():
			{
				subMsg := msg.(device.CollectInterfaceEventTemplate)
				//判断设备在该上报服务中
				index := -1
				for k, v := range r.NodeList {
					if v.Name == subMsg.NodeName {
						index = k
					}
				}
				if index == -1 {
					continue
				}
				setting.ZAPS.Debugf("上报服务[%s] 采集接口[%s] 设备[%s] 主题[%s] 消息内容[%v]",
					r.GWParam.ServiceName,
					subMsg.CollName,
					subMsg.NodeName,
					subMsg.Topic,
					subMsg.Content)
				nodeName := make([]string, 0)
				switch subMsg.Topic {
				case "onLine":
					{
						nodeName = append(nodeName, subMsg.NodeName)
						r.NodeList[index].CommStatus = "onLine"
						r.LogInRequestFrameChan <- nodeName
					}
				case "offLine":
					{
						nodeName = append(nodeName, subMsg.NodeName)
						r.NodeList[index].CommStatus = "offLine"
						r.LogOutRequestFrameChan <- nodeName
					}
				case "update":
					{
						//更新设备的通信状态
						r.NodeList[index].CommStatus = "onLine"

						//coll, ok := device.CollectInterfaceMap.Coll[subMsg.CollName]
						//if !ok {
						//	continue
						//}
						//node, ok := coll.DeviceNodeMap[subMsg.NodeName]
						//if !ok {
						//	return
						//}

						reportStatus := false
						//for _, v := range node.Properties {
						//	if v.Params.StepAlarm == true {
						//		valueCnt := len(v.Value)
						//		if valueCnt >= 2 { //阶跃报警必须是2个值
						//			if v.Type == device.PropertyTypeInt32 {
						//				pValueCur := v.Value[valueCnt-1].Value.(int32)
						//				pValuePre := v.Value[valueCnt-2].Value.(int32)
						//				step, _ := strconv.Atoi(v.Params.Step)
						//				if math.Abs(float64(pValueCur-pValuePre)) > float64(step) {
						//					reportStatus = true //满足报警条件，上报
						//					nodeName = append(nodeName, node.Name)
						//				}
						//			} else if v.Type == device.PropertyTypeUInt32 {
						//				pValueCur := v.Value[valueCnt-1].Value.(uint32)
						//				pValuePre := v.Value[valueCnt-2].Value.(uint32)
						//				step, _ := strconv.Atoi(v.Params.Step)
						//				if math.Abs(float64(pValueCur-pValuePre)) > float64(step) {
						//					reportStatus = true //满足报警条件，上报
						//					nodeName = append(nodeName, node.Name)
						//				}
						//			} else if v.Type == device.PropertyTypeDouble {
						//				pValueCur := v.Value[valueCnt-1].Value.(float64)
						//				pValuePre := v.Value[valueCnt-2].Value.(float64)
						//				step, err := strconv.ParseFloat(v.Params.Step, 64)
						//				if err != nil {
						//					continue
						//				}
						//				if math.Abs(pValueCur-pValuePre) > float64(step) {
						//					reportStatus = true //满足报警条件，上报
						//					nodeName = append(nodeName, node.Name)
						//				}
						//			}
						//		}
						//	} else if v.Params.MinMaxAlarm == true {
						//		valueCnt := len(v.Value)
						//		if v.Type == device.PropertyTypeInt32 {
						//			pValueCur := v.Value[valueCnt-1].Value.(int32)
						//			min, _ := strconv.Atoi(v.Params.Min)
						//			max, _ := strconv.Atoi(v.Params.Max)
						//			if pValueCur < int32(min) || pValueCur > int32(max) {
						//				reportStatus = true //满足报警条件，上报
						//				nodeName = append(nodeName, node.Name)
						//			}
						//		} else if v.Type == device.PropertyTypeUInt32 {
						//			pValueCur := v.Value[valueCnt-1].Value.(uint32)
						//			min, _ := strconv.Atoi(v.Params.Min)
						//			max, _ := strconv.Atoi(v.Params.Max)
						//			if pValueCur < uint32(min) || pValueCur > uint32(max) {
						//				reportStatus = true //满足报警条件，上报
						//				nodeName = append(nodeName, node.Name)
						//			}
						//		} else if v.Type == device.PropertyTypeDouble {
						//			pValueCur := v.Value[valueCnt-1].Value.(float64)
						//			min, err := strconv.ParseFloat(v.Params.Min, 64)
						//			if err != nil {
						//				continue
						//			}
						//			max, err := strconv.ParseFloat(v.Params.Max, 64)
						//			if err != nil {
						//				continue
						//			}
						//			if pValueCur < min || pValueCur > max {
						//				reportStatus = true //满足报警条件，上报
						//				nodeName = append(nodeName, node.Name)
						//			}
						//		}
						//	}
						//}

						if reportStatus == true {
							reportNodeProperty := MQTTThingsBoardReportPropertyTemplate{
								DeviceType: "node",
								DeviceName: nodeName,
							}
							r.ReportPropertyRequestFrameChan <- reportNodeProperty
						}
					}
				}
			}
		}
	}
}

func (r *ReportServiceParamThingsBoardTemplate) LogIn(nodeName []string) {

	////清空接收chan，避免出现有上次接收的缓存
	//for i := 0; i < len(r.ReceiveLogInAckFrameChan); i++ {
	//	<-r.ReceiveLogInAckFrameChan
	//}

	r.NodeLogIn(nodeName)
}

func (r *ReportServiceParamThingsBoardTemplate) LogOut(nodeName []string) {

	////清空接收chan，避免出现有上次接收的缓存
	//for i := 0; i < len(r.ReceiveLogOutAckFrameChan); i++ {
	//	<-r.ReceiveLogOutAckFrameChan
	//}

	r.NodeLogOut(nodeName)
}

func (r *ReportServiceParamThingsBoardTemplate) PeriodReportHandler() {

	//网关上报
	reportGWProperty := MQTTThingsBoardReportPropertyTemplate{
		DeviceType: "gw",
	}
	r.ReportPropertyRequestFrameChan <- reportGWProperty
	if r.GWParam.ReportStatus == "onLine" {
		//全部末端设备上报
		nodeName := make([]string, 0)
		for _, v := range r.NodeList {
			if v.CommStatus == "onLine" {
				nodeName = append(nodeName, v.Name)
			}
		}
		setting.ZAPS.Debugf("上报服务[%s]定时上报任务中上报节点%v", r.GWParam.ServiceName, nodeName)
		if len(nodeName) > 0 {
			reportNodeProperty := MQTTThingsBoardReportPropertyTemplate{
				DeviceType: "node",
				DeviceName: nodeName,
			}
			r.ReportPropertyRequestFrameChan <- reportNodeProperty
		}
	}
}

//查看上报服务中设备是否离线
func (r *ReportServiceParamThingsBoardTemplate) ReportOfflineTimeFun() {

	//setting.ZAPS.Infof("上报服务[%s] 巡检网关和节点是否离线任务", r.GWParam.ServiceName)
	//if r.GWParam.ReportErrCnt >= 5 {
	//	r.GWParam.ReportStatus = "offLine"
	//	setting.ZAPS.Warnf("上报服务[%s] 网关离线", r.GWParam.ServiceName)
	//	r.GWParam.MQTTClient.Disconnect(0)
	//
	//	if r.GWLogin() == true {
	//		setting.ZAPS.Debugf("上报服务[%s] 网关 重新登录成功", r.GWParam.ServiceName)
	//	} else {
	//		setting.ZAPS.Warnf("上报服务[%s] 网关 重新登录失败", r.GWParam.ServiceName)
	//	}
	//}
	//r.GWParam.ReportErrCnt = 0
}

func NewReportServiceParamThingsBoard(gw ReportServiceGWParamThingsBoardTemplate, nodeList []ReportServiceNodeParamThingsBoardTemplate) *ReportServiceParamThingsBoardTemplate {

	for _, v := range nodeList {
		rModel, ok := reportModel.ReportModels[v.UploadModel]
		if !ok {
			continue
		}
		v.Properties = rModel.Properties
	}

	ThingsBoardParam := &ReportServiceParamThingsBoardTemplate{
		Index:                 len(ReportServiceParamListThingsBoard.ServiceList),
		GWParam:               gw,
		NodeList:              nodeList,
		ReceiveFrameChan:      make(chan MQTTThingsBoardReceiveFrameTemplate, 100),
		LogInRequestFrameChan: make(chan []string, 0),
		//ReceiveLogInAckFrameChan:             make(chan MQTTThingsBoardLogInAckTemplate, 5),
		LogOutRequestFrameChan: make(chan []string, 0),
		//ReceiveLogOutAckFrameChan:            make(chan MQTTThingsBoardLogOutAckTemplate, 5),
		ReportPropertyRequestFrameChan: make(chan MQTTThingsBoardReportPropertyTemplate, 50),
		//ReceiveReportPropertyAckFrameChan:    make(chan MQTTThingsBoardReportPropertyAckTemplate, 50),
		ReceiveInvokeServiceRequestFrameChan: make(chan InvokeServiceRequestTemplate, 50),
		ReceiveInvokeServiceAckFrameChan:     make(chan InvokeServiceAckTemplate, 50),
		//ReceiveWritePropertyRequestFrameChan: make(chan MQTTThingsBoardWritePropertyRequestTemplate, 50),
		//ReceiveReadPropertyRequestFrameChan:  make(chan MQTTThingsBoardReadPropertyRequestTemplate, 50),
	}

	ctx, cancel := context.WithCancel(context.Background())
	ThingsBoardParam.CancelFunc = cancel

	go ReportServiceThingsBoardPoll(ctx, ThingsBoardParam)

	return ThingsBoardParam
}

func ReportServiceThingsBoardInit() {

	writeTimer = time.AfterFunc(time.Second, func() {
		ReportServiceParamListThingsBoard.WriteParamToJson()
	})
	writeTimer.Stop()

	ReportServiceParamListThingsBoard.ReadParamFromJson()
}

func ReportServiceThingsBoardPoll(ctx context.Context, r *ReportServiceParamThingsBoardTemplate) {

	reportState := 0

	// 定义一个cron运行器
	scheduler := gocron.NewScheduler()

	setting.ZAPS.Infof("上报服务[%s]定时上报任务周期为[%v]", r.GWParam.ServiceName, r.GWParam.ReportTime)
	_ = scheduler.Every(uint64(r.GWParam.ReportTime)).Second().Do(r.PeriodReportHandler)

	//订阅采集接口消息
	device.CollectInterfaceMap.Lock.Lock()
	r.EventSub = eventBus.NewSub()
	for _, coll := range device.CollectInterfaceMap.Coll {
		setting.ZAPS.Infof("上报服务[%s]订阅采集接口[%v]", r.GWParam.ServiceName, coll.CollInterfaceName)
		coll.CollEventBus.Subscribe("onLine", r.EventSub)
		coll.CollEventBus.Subscribe("offLine", r.EventSub)
		coll.CollEventBus.Subscribe("update", r.EventSub)
	}
	device.CollectInterfaceMap.Lock.Unlock()
	go r.ProcessCollEvent(ctx, r.EventSub)

	//订阅虚拟接口消息
	virtualEventSub := eventBus.NewSub()
	virtual.VirtualDevice.EventBus.Subscribe("onLine", virtualEventSub)
	virtual.VirtualDevice.EventBus.Subscribe("offLine", virtualEventSub)
	virtual.VirtualDevice.EventBus.Subscribe("update", virtualEventSub)

	go r.ProcessUpLinkFrame(ctx)
	go r.ProcessDownLinkFrame(ctx)

	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		default:
			{
				if reportState == 0 {
					if r.GWLogin() == true {
						reportState = 1
						scheduler.Start()
					}
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
