package mqttSagooIOT

import (
	"context"
	"encoding/json"
	"gateway/device"
	"gateway/device/eventBus"
	"gateway/report/reportModel"
	"gateway/setting"
	"gateway/virtual"
	"github.com/jasonlvhit/gocron"
	"math"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//上报节点参数结构体
type ReportServiceNodeParamSagooIOTTemplate struct {
	Index             int    `json:"index"`
	ServiceName       string `json:"serviceName"`
	CollInterfaceName string `json:"collInterfaceName"`
	Name              string `json:"deviceName"`
	Label             string `json:"deviceLabel"`
	Addr              string `json:"deviceAddr"`
	CommStatus        string
	ReportErrCnt      int `json:"-"`
	ReportStatus      string
	UploadModel       string                                              `json:"uploadModel"`
	Protocol          string                                              `json:"protocol"`
	Properties        map[string]*reportModel.ReportModelPropertyTemplate `json:"properties"`
	Param             struct {
		ProductKey string `json:"productKey"`
		DeviceCode string `json:"deviceCode"`
	}
}

//上报网关参数结构体
type ReportServiceGWParamSagooIOTTemplate struct {
	ServiceName  string
	IP           string
	Port         string
	ReportStatus string
	ReportTime   int
	ReportErrCnt int
	Protocol     string
	Param        struct {
		UserName     string `json:"userName"`
		Password     string `json:"password"`
		ClientID     string `json:"clientID"`
		KeepAlive    string `json:"keepAlive"`
		CleanSession bool   `json:"cleanSession"`
	}
	MQTTClient        MQTT.Client         `json:"-"`
	MQTTClientOptions *MQTT.ClientOptions `json:"-"`
	MQTTClientID      int                 `json:"-"`
	OfflineChan       chan bool           `json:"-"`
}

//上报服务参数，网关参数，节点参数
type ReportServiceParamSagooIOTTemplate struct {
	Index                    int `json:"index"`
	GWParam                  ReportServiceGWParamSagooIOTTemplate
	NodeList                 []*ReportServiceNodeParamSagooIOTTemplate
	ReceiveFrameChan         chan MQTTSagooIOTReceiveFrameTemplate `json:"-"`
	LogInRequestFrameChan    chan string                           `json:"-"` //上线
	ReceiveLogInAckFrameChan chan MQTTSagooIOTLogInAckTemplate     `json:"-"`
	LogOutRequestFrameChan   chan string                           `json:"-"`
	//ReceiveLogOutAckFrameChan              chan MQTTSagooIOTLogOutAckTemplate             `json:"-"`
	ReportPropertyRequestFrameChan         chan MQTTSagooIOTReportPropertyTemplate        `json:"-"`
	ReceiveReportPropertyAckFrameChan      chan MQTTSagooIOTReportPropertyAckTemplate     `json:"-"`
	ReportAlarmRequestFrameChan            chan MQTTSagooIOTReportAlarmTemplate           `json:"-"`
	ReceiveReportAlarmAckFrameChan         chan MQTTSagooIOTReportPropertyAckTemplate     `json:"-"`
	ReceiveInvokeServiceRequestFrameChan   chan MQTTSagooIOTInvokeServiceRequestTemplate  `json:"-"`
	ReceiveInvokeServiceAckFrameChan       chan MQTTSagooIOTInvokeServiceAckTemplate      `json:"-"`
	ReceiveWritePropertyRequestFrameChan   chan MQTTSagooIOTWritePropertyRequestTemplate  `json:"-"`
	ReceiveReadPropertyRequestFrameChan    chan MQTTSagooIOTReadPropertyRequestTemplate   `json:"-"`
	ReceiveReadNodeStatusRequestFrameChan  chan MQTTSagooIOTReadNodeStatusRequestTemplate `json:"-"`
	ReceiveInvokeGWServiceRequestFrameChan chan MQTTSagooIOTInvokeServiceRequestTemplate  `json:"-"`
	ReceiveInvokeGWServiceAckFrameChan     chan MQTTSagooIOTInvokeServiceAckTemplate      `json:"-"`
	CancelFunc                             context.CancelFunc                             `json:"-"`
	EventSub                               eventBus.Sub                                   `json:"-"`
	MessageEventBus                        eventBus.Bus                                   `json:"-"` //通信报文总线
}

type ReportServiceParamListSagooIOTTemplate struct {
	ServiceList []*ReportServiceParamSagooIOTTemplate
}

const (
	SagooIOTTimeOutLogin          int = 60
	SagooIOTTimeOutLogout         int = 1000
	SagooIOTTimeOutSubscribe      int = 5
	SagooIOTTimeOutReportProperty int = 5
	SagooIOTTimeOutReadProperty   int = 5
	SagooIOTTimeOutWriteProperty  int = 5
	SagooIOTTimeOutService        int = 5
	SagooIOTTimeOutReadNode       int = 5
)

//实例化上报服务
var ReportServiceParamListSagooIOT ReportServiceParamListSagooIOTTemplate

func (r *ReportServiceParamSagooIOTTemplate) ProcessUpLinkFrame(ctx context.Context) {

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
					r.NodePropertyPost(reqFrame.DeviceName, nil, false)
				}
			}
		case reqFrame := <-r.ReportAlarmRequestFrameChan:
			{
				if reqFrame.DeviceType == "gw" {
					r.GWPropertyPost()
				} else if reqFrame.DeviceType == "node" {
					if setting.SystemState.LockStatus == setting.LockCmdDisable {
						r.NodePropertyPost(reqFrame.DeviceName, reqFrame.Properties, true)
					}
				}
			}
		case reqFrame := <-r.ReceiveWritePropertyRequestFrameChan:
			{
				r.ReportServiceSagooIOTProcessWriteProperty(reqFrame)
			}
		case reqFrame := <-r.ReceiveReadPropertyRequestFrameChan:
			{
				r.ReportServiceSagooIOTProcessReadProperty(reqFrame)
			}
		case reqFrame := <-r.ReceiveInvokeServiceRequestFrameChan:
			{
				r.ReportServiceSagooIOTProcessInvokeService(reqFrame)
			}
		case reqFrame := <-r.ReceiveReadNodeStatusRequestFrameChan:
			{
				r.ReportServiceSagooIOTProcessReadNodeStatus(reqFrame)
			}
		case reqFrame := <-r.ReceiveInvokeGWServiceRequestFrameChan:
			{
				r.ReportServiceSagooIOTProcessInvokeGWService(reqFrame)
			}
		}
	}
}

func (r *ReportServiceParamSagooIOTTemplate) ProcessDownLinkFrame(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case frame := <-r.ReceiveFrameChan:
			{
				//setting.ZAPS.Debugf("Recv TOPIC: %s", frame.Topic)
				//setting.ZAPS.Debugf("Recv MSG: %s", frame.Payload)
				if strings.Contains(frame.Topic, "/sys/thing/node/property/post_reply") { //子设备上报属性回应

					ackFrame := MQTTSagooIOTReportPropertyAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.ZAPS.Errorf("ReportPropertyAck json unmarshal err")
						continue
					}
					r.ReceiveReportPropertyAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/node/login/post_reply") { //子设备上线回应

					ackFrame := MQTTSagooIOTLogInAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.ZAPS.Warnf("LogInAck json unmarshal err")
						continue
					}
					r.ReceiveLogInAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/node/logout/post_reply") { //子设备下线回应

				} else if strings.Contains(frame.Topic, "/device/data/cmd") { //设备服务调用
					serviceFrame := MQTTSagooIOTInvokeServiceRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &serviceFrame)
					if err != nil {
						setting.ZAPS.Errorf("serviceFrame json unmarshal err")
						continue
					}
					r.ReceiveInvokeServiceRequestFrameChan <- serviceFrame
				} else if strings.Contains(frame.Topic, "/device/data/set") { //子设备设置属性请求
					writePropertyFrame := MQTTSagooIOTWritePropertyRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &writePropertyFrame)
					if err != nil {
						setting.ZAPS.Errorf("writePropertyFrame json unmarshal err")
						continue
					}
					r.ReceiveWritePropertyRequestFrameChan <- writePropertyFrame
				} else if strings.Contains(frame.Topic, "/device/data/get") { //子设备获取属性请求
					readPropertyFrame := MQTTSagooIOTReadPropertyRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &readPropertyFrame)
					if err != nil {
						setting.ZAPS.Errorf("readPropertyFrame json unmarshal err")
						continue
					}
					r.ReceiveReadPropertyRequestFrameChan <- readPropertyFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/node/status/get") { //子设备状态查询
					readNodeStatusFrame := MQTTSagooIOTReadNodeStatusRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &readNodeStatusFrame)
					if err != nil {
						setting.ZAPS.Errorf("readNodeStatusFrame json unmarshal err")
						continue
					}
					r.ReceiveReadNodeStatusRequestFrameChan <- readNodeStatusFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/gw/service/invoke") { //网关服务调用
					serviceFrame := MQTTSagooIOTInvokeServiceRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &serviceFrame)
					if err != nil {
						setting.ZAPS.Errorf("serviceFrame json unmarshal err")
						continue
					}
					r.ReceiveInvokeGWServiceRequestFrameChan <- serviceFrame
				}
			}
		}
	}
}

func (r *ReportServiceParamSagooIOTTemplate) ProcessCollEvent(ctx context.Context, sub eventBus.Sub) {
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
				switch subMsg.Topic {
				case "onLine":
					{
						r.NodeList[index].CommStatus = "onLine"

						//判断告警
						r.ProcessAlarmEvent(index, subMsg.CollName, subMsg.NodeName)

						if len(r.LogInRequestFrameChan) == 50 {
							<-r.LogInRequestFrameChan
						}
						r.LogInRequestFrameChan <- subMsg.NodeName
					}
				case "offLine":
					{
						r.NodeList[index].CommStatus = "offLine"
						if len(r.LogOutRequestFrameChan) == 50 {
							<-r.LogOutRequestFrameChan
						}
						r.LogOutRequestFrameChan <- subMsg.NodeName
					}
				case "update":
					{
						//更新设备的通信状态
						r.NodeList[index].CommStatus = "onLine"

						//判断告警
						r.ProcessAlarmEvent(index, subMsg.CollName, subMsg.NodeName)
					}
				}
			}
		}
	}
}

func (r *ReportServiceParamSagooIOTTemplate) ProcessVirtualEvent(ctx context.Context, sub eventBus.Sub) {
	for {
		select {
		case <-ctx.Done():
			{
				//取消订阅
				virtual.VirtualDevice.EventBus.UnSubscribe("update", sub)
				return
			}
		case msg := <-sub.Out():
			{
				subMsg := msg.(virtual.VirtualEventTemplate)

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

				setting.ZAPS.Debugf("上报服务[%s] 主题[%s] 消息内容[%v]",
					r.GWParam.ServiceName,
					subMsg.Topic,
					subMsg.PropertyNames)
				switch subMsg.Topic {
				case "onLine":
					{
						r.NodeList[index].CommStatus = "onLine"

						if len(r.LogInRequestFrameChan) == 50 {
							<-r.LogInRequestFrameChan
						}
						r.LogInRequestFrameChan <- subMsg.NodeName
					}
				case "offLine":
					{
						r.NodeList[index].CommStatus = "offLine"
						if len(r.LogOutRequestFrameChan) == 50 {
							<-r.LogOutRequestFrameChan
						}
						r.LogOutRequestFrameChan <- subMsg.NodeName
					}
				case "update":
					{
						//判断告警
						reportNodeProperty := MQTTSagooIOTReportPropertyTemplate{
							DeviceType: "node",
						}
						r.ReportPropertyRequestFrameChan <- reportNodeProperty
					}
				}
			}
		}
	}
}

//func (r *ReportServiceParamSagooIOTTemplate) ProcessUpdateEvent(nodeIndex int, collName string, nodeName string) {
//
//	_, ok := reportModel.ReportModels[r.NodeList[nodeIndex].UploadModel]
//	if !ok {
//		setting.ZAPS.Errorf("上报服务[%s]上报设备[%s]上报模型不存在", r.GWParam.ServiceName, r.NodeList[nodeIndex].Name)
//		return
//	}
//
//	for _, v := range device.CollectInterfaceMap.Coll[collName].DeviceNodeMap[nodeName].Properties {
//		setting.ZAPS.Debugf("nodeName %v reportProperty %+v", r.NodeList[nodeIndex].Name, v)
//		p, ok := r.NodeList[nodeIndex].Properties[v.Name]
//		if !ok {
//			continue
//		}
//
//		if len(v.Value) > 0 {
//			setting.ZAPS.Debugf("当前属性 %+v", v.Value[len(v.Value)-1])
//			property := reportModel.ReportModelPropertyValueTemplate{}
//			property.Index = len(p.Value)
//			property.Value = v.Value[len(v.Value)-1].Value
//			property.Explain = v.Value[len(v.Value)-1].Explain
//			property.TimeStamp = v.Value[len(v.Value)-1].TimeStamp
//			if len(v.Value) > 2 {
//				r.NodeList[nodeIndex].Properties[v.Name].Value = append(r.NodeList[nodeIndex].Properties[v.Name].Value[:0], r.NodeList[nodeIndex].Properties[v.Name].Value[1:]...)
//			}
//			r.NodeList[nodeIndex].Properties[v.Name].Value = append(r.NodeList[nodeIndex].Properties[v.Name].Value, property)
//		}
//		setting.ZAPS.Debugf("nodeName %v name %v,添加后 %+v", r.NodeList[nodeIndex].Name, v.Name, p)
//	}
//}

func (r *ReportServiceParamSagooIOTTemplate) ProcessAlarmEvent(index int, collName string, nodeName string) {

	reportStatus := false

	properties := make([]MQTTSagooIOTPropertyPostParamPropertyTemplate, 0)

	for _, v := range device.CollectInterfaceMap.Coll[collName].DeviceNodeMap[nodeName].Properties {
		//从上报模型中查找属性
		property, ok := r.NodeList[index].Properties[v.Name]
		if !ok {
			continue
		}
		if property.Params.StepAlarm == true {
			valueCnt := len(v.Value)
			if valueCnt >= 2 { //阶跃报警必须是2个值
				switch v.Type {
				case device.PropertyTypeInt32:
					{
						pValueCur := v.Value[valueCnt-1].Value.(int32)
						pValuePre := v.Value[valueCnt-2].Value.(int32)
						step, _ := strconv.Atoi(property.Params.Step)
						if math.Abs(float64(pValueCur-pValuePre)) > float64(step) {
							reportStatus = true //满足报警条件，上报
							setting.ZAPS.Infof("设备[%v]阶跃报警", r.NodeList[index].Name)
							//转换时间
							//timeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", v.Value[valueCnt-1].TimeStamp, time.Local)
							property := MQTTSagooIOTPropertyPostParamPropertyTemplate{
								Name:      v.Name,
								Value:     v.Value[valueCnt-1].Value.(int32),
								TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
							}
							properties = append(properties, property)
						}
					}
				case device.PropertyTypeUInt32:
					{
						pValueCur := v.Value[valueCnt-1].Value.(uint32)
						pValuePre := v.Value[valueCnt-2].Value.(uint32)
						step, _ := strconv.Atoi(property.Params.Step)
						if math.Abs(float64(pValueCur-pValuePre)) > float64(step) {
							reportStatus = true //满足报警条件，上报
							setting.ZAPS.Infof("设备[%v]阶跃报警", r.NodeList[index].Name)
							//转换时间
							//timeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", v.Value[valueCnt-1].TimeStamp, time.Local)
							property := MQTTSagooIOTPropertyPostParamPropertyTemplate{
								Name:      v.Name,
								Value:     v.Value[valueCnt-1].Value.(uint32),
								TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
							}
							properties = append(properties, property)
						}
					}
				case device.PropertyTypeDouble:
					{
						pValueCur := v.Value[valueCnt-1].Value.(float64)
						pValuePre := v.Value[valueCnt-2].Value.(float64)
						step, err := strconv.ParseFloat(property.Params.Step, 64)
						if err != nil {
							continue
						}
						if math.Abs(pValueCur-pValuePre) > float64(step) {
							reportStatus = true //满足报警条件，上报
							setting.ZAPS.Infof("设备[%v]阶跃报警", r.NodeList[index].Name)
							//转换时间
							//timeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", v.Value[valueCnt-1].TimeStamp, time.Local)
							property := MQTTSagooIOTPropertyPostParamPropertyTemplate{
								Name:      v.Name,
								Value:     v.Value[valueCnt-1].Value.(float64),
								TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
							}
							properties = append(properties, property)
						}
					}
				}
			}
		}
		if property.Params.MinMaxAlarm == true {
			valueCnt := len(v.Value)
			if v.Type == device.PropertyTypeInt32 {
				pValueCur := v.Value[valueCnt-1].Value.(int32)
				min, _ := strconv.Atoi(property.Params.Min)
				max, _ := strconv.Atoi(property.Params.Max)
				if pValueCur < int32(min) || pValueCur > int32(max) {
					reportStatus = true //满足报警条件，上报
					setting.ZAPS.Infof("设备[%v]范围报警", r.NodeList[index].Name)
					property := MQTTSagooIOTPropertyPostParamPropertyTemplate{
						Name:      v.Name,
						Value:     v.Value[valueCnt-1].Value.(int32),
						TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
					}
					properties = append(properties, property)
				}
			} else if v.Type == device.PropertyTypeUInt32 {
				pValueCur := v.Value[valueCnt-1].Value.(uint32)
				min, _ := strconv.Atoi(property.Params.Min)
				max, _ := strconv.Atoi(property.Params.Max)
				if pValueCur < uint32(min) || pValueCur > uint32(max) {
					reportStatus = true //满足报警条件，上报
					setting.ZAPS.Infof("设备[%v]范围报警", r.NodeList[index].Name)
					property := MQTTSagooIOTPropertyPostParamPropertyTemplate{
						Name:      v.Name,
						Value:     v.Value[valueCnt-1].Value.(uint32),
						TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
					}
					properties = append(properties, property)
				}
			} else if v.Type == device.PropertyTypeDouble {
				pValueCur := v.Value[valueCnt-1].Value.(float64)
				min, err := strconv.ParseFloat(property.Params.Min, 64)
				if err != nil {
					continue
				}
				max, err := strconv.ParseFloat(property.Params.Max, 64)
				if err != nil {
					continue
				}
				if pValueCur < min || pValueCur > max {
					reportStatus = true //满足报警条件，上报
					setting.ZAPS.Infof("设备[%v]范围报警", r.NodeList[index].Name)
					property := MQTTSagooIOTPropertyPostParamPropertyTemplate{
						Name:      v.Name,
						Value:     v.Value[valueCnt-1].Value.(float64),
						TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
					}
					properties = append(properties, property)
				}
			}
		}
	}

	nodeNames := make([]string, 0)
	if reportStatus == true {
		nodeNames = append(nodeNames, r.NodeList[index].Name)
		reportNodeProperty := MQTTSagooIOTReportAlarmTemplate{
			DeviceType: "node",
			DeviceName: nodeNames,
			Properties: properties,
		}
		if len(r.ReportAlarmRequestFrameChan) == 50 {
			<-r.ReportAlarmRequestFrameChan
		}
		r.ReportAlarmRequestFrameChan <- reportNodeProperty
	}
}

func (r *ReportServiceParamSagooIOTTemplate) LogIn(nodeName string) {

	//清空接收chan，避免出现有上次接收的缓存
	for i := 0; i < len(r.ReceiveLogInAckFrameChan); i++ {
		<-r.ReceiveLogInAckFrameChan
	}

	r.NodeLogIn(nodeName)
}

func (r *ReportServiceParamSagooIOTTemplate) LogOut(nodeName string) {

	//清空接收chan，避免出现有上次接收的缓存

	r.NodeLogOut(nodeName)
}

func (r *ReportServiceParamSagooIOTTemplate) ReportTimeFun() {

	//网关上报
	reportGWProperty := MQTTSagooIOTReportPropertyTemplate{
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
			reportNodeProperty := MQTTSagooIOTReportPropertyTemplate{
				DeviceType: "node",
				DeviceName: nodeName,
			}
			r.ReportPropertyRequestFrameChan <- reportNodeProperty
		}
	}
}

func (r *ReportServiceParamSagooIOTTemplate) ReportNode(name string) {

	if r.GWParam.ReportStatus == "onLine" {
		//全部末端设备上报
		nodeName := make([]string, 0)

		nodeName = append(nodeName, name)
		setting.ZAPS.Debugf("上报服务[%s]手动设置节点上报%v", r.GWParam.ServiceName, nodeName)
		if len(nodeName) > 0 {
			reportNodeProperty := MQTTSagooIOTReportPropertyTemplate{
				DeviceType: "node",
				DeviceName: nodeName,
			}
			r.ReportPropertyRequestFrameChan <- reportNodeProperty
		}
	}
}

func (r *ReportServiceParamSagooIOTTemplate) ReportNodeNoCheck(name string) {

	//全部末端设备上报
	nodeName := make([]string, 0)

	nodeName = append(nodeName, name)
	setting.ZAPS.Debugf("上报服务[%s]手动设置节点上报%v", r.GWParam.ServiceName, nodeName)
	if len(nodeName) > 0 {
		reportNodeProperty := MQTTSagooIOTReportPropertyTemplate{
			DeviceType: "node",
			DeviceName: nodeName,
		}
		setting.ZAPS.Debugf("上报服务[%s]上报请求队列数量%v", r.GWParam.ServiceName, len(r.ReportPropertyRequestFrameChan))
		r.ReportPropertyRequestFrameChan <- reportNodeProperty
	}
}

func ReportServiceSagooIOTInit() {
	ReportServiceSagooIOTReadParamFromJson()

	ReportServiceSagooIOTConfigInit()
}

func NewReportServiceParamSagooIOT(gw ReportServiceGWParamSagooIOTTemplate, nodeList []*ReportServiceNodeParamSagooIOTTemplate) *ReportServiceParamSagooIOTTemplate {

	for _, v := range nodeList {
		rModel, ok := reportModel.ReportModels[v.UploadModel]
		if !ok {
			continue
		}
		v.Properties = rModel.Properties
	}

	SagooIOTParam := &ReportServiceParamSagooIOTTemplate{
		Index:                    len(ReportServiceParamListSagooIOT.ServiceList),
		GWParam:                  gw,
		NodeList:                 nodeList,
		ReceiveFrameChan:         make(chan MQTTSagooIOTReceiveFrameTemplate, 500),
		LogInRequestFrameChan:    make(chan string, 50),
		ReceiveLogInAckFrameChan: make(chan MQTTSagooIOTLogInAckTemplate, 5),
		LogOutRequestFrameChan:   make(chan string, 50),
		//ReceiveLogOutAckFrameChan:              make(chan MQTTSagooIOTLogOutAckTemplate, 5),
		ReportPropertyRequestFrameChan:         make(chan MQTTSagooIOTReportPropertyTemplate, 500),
		ReceiveReportPropertyAckFrameChan:      make(chan MQTTSagooIOTReportPropertyAckTemplate, 50),
		ReportAlarmRequestFrameChan:            make(chan MQTTSagooIOTReportAlarmTemplate, 50),
		ReceiveReportAlarmAckFrameChan:         make(chan MQTTSagooIOTReportPropertyAckTemplate, 50),
		ReceiveInvokeServiceRequestFrameChan:   make(chan MQTTSagooIOTInvokeServiceRequestTemplate, 500),
		ReceiveInvokeServiceAckFrameChan:       make(chan MQTTSagooIOTInvokeServiceAckTemplate, 50),
		ReceiveWritePropertyRequestFrameChan:   make(chan MQTTSagooIOTWritePropertyRequestTemplate, 500),
		ReceiveReadPropertyRequestFrameChan:    make(chan MQTTSagooIOTReadPropertyRequestTemplate, 500),
		ReceiveReadNodeStatusRequestFrameChan:  make(chan MQTTSagooIOTReadNodeStatusRequestTemplate, 50),
		ReceiveInvokeGWServiceRequestFrameChan: make(chan MQTTSagooIOTInvokeServiceRequestTemplate, 50),
		ReceiveInvokeGWServiceAckFrameChan:     make(chan MQTTSagooIOTInvokeServiceAckTemplate, 50),
		MessageEventBus:                        eventBus.NewBus(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	SagooIOTParam.CancelFunc = cancel

	go ReportServiceSagooIOTPoll(ctx, SagooIOTParam)

	return SagooIOTParam
}

func ReportServiceSagooIOTPoll(ctx context.Context, r *ReportServiceParamSagooIOTTemplate) {

	reportState := 0

	// 定义一个cron运行器
	scheduler := gocron.NewScheduler()

	setting.ZAPS.Infof("上报服务[%s]定时上报任务周期为[%v]", r.GWParam.ServiceName, r.GWParam.ReportTime)
	_ = scheduler.Every(uint64(r.GWParam.ReportTime)).Second().Do(r.ReportTimeFun)

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
	go r.ProcessVirtualEvent(ctx, virtualEventSub)

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
