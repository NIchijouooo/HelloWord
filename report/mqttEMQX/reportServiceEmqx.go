package mqttEmqx

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"gateway/device"
	"gateway/device/eventBus"
	"gateway/report/reportModel"
	"gateway/setting"
	"gateway/utils"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jasonlvhit/gocron"
	"math"
	"os"
	"strconv"
	"strings"
)

//上报节点参数结构体
type ReportServiceNodeParamEmqxTemplate struct {
	ServiceName       string `json:"serviceName"`
	CollInterfaceName string `json:"collInterfaceName"`
	Name              string `json:"deviceName"`
	Label             string `json:"Label"`
	Addr              string `json:"deviceAddr"`
	CommStatus        string
	ReportErrCnt      int `json:"-"`
	ReportStatus      string
	UploadModel       string                                              `json:"uploadModel"`
	Protocol          string                                              `json:"protocol"`
	Properties        map[string]*reportModel.ReportModelPropertyTemplate `json:"properties"`
	Param             struct {
		DeviceCode string `json:"deviceCode"`
	}
}

//上报网关参数结构体
type ReportServiceGWParamEmqxTemplate struct {
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
type ReportServiceParamEmqxTemplate struct {
	GWParam                                ReportServiceGWParamEmqxTemplate
	NodeList                               []ReportServiceNodeParamEmqxTemplate
	ReceiveFrameChan                       chan MQTTEmqxReceiveFrameTemplate          `json:"-"`
	LogInRequestFrameChan                  chan []string                              `json:"-"` //上线
	ReceiveLogInAckFrameChan               chan MQTTEmqxLogInAckTemplate              `json:"-"`
	LogOutRequestFrameChan                 chan []string                              `json:"-"`
	ReceiveLogOutAckFrameChan              chan MQTTEmqxLogOutAckTemplate             `json:"-"`
	ReportPropertyRequestFrameChan         chan MQTTEmqxReportPropertyTemplate        `json:"-"`
	ReceiveReportPropertyAckFrameChan      chan MQTTEmqxReportPropertyAckTemplate     `json:"-"`
	ReportAlarmRequestFrameChan            chan MQTTEmqxReportAlarmTemplate           `json:"-"`
	ReceiveReportAlarmAckFrameChan         chan MQTTEmqxReportPropertyAckTemplate     `json:"-"`
	ReceiveInvokeServiceRequestFrameChan   chan MQTTEmqxInvokeServiceRequestTemplate  `json:"-"`
	ReceiveInvokeServiceAckFrameChan       chan MQTTEmqxInvokeServiceAckTemplate      `json:"-"`
	ReceiveWritePropertyRequestFrameChan   chan MQTTEmqxWritePropertyRequestTemplate  `json:"-"`
	ReceiveReadPropertyRequestFrameChan    chan MQTTEmqxReadPropertyRequestTemplate   `json:"-"`
	ReceiveReadNodeStatusRequestFrameChan  chan MQTTEmqxReadNodeStatusRequestTemplate `json:"-"`
	ReceiveInvokeGWServiceRequestFrameChan chan MQTTEmqxInvokeServiceRequestTemplate  `json:"-"`
	ReceiveInvokeGWServiceAckFrameChan     chan MQTTEmqxInvokeServiceAckTemplate      `json:"-"`
	CancelFunc                             context.CancelFunc                         `json:"-"`
	EventSub                               eventBus.Sub                               `json:"-"`
}

type ReportServiceParamListEmqxTemplate struct {
	ServiceList []*ReportServiceParamEmqxTemplate
}

const (
	EMQXTimeOutLogin          int = 10
	EMQXTimeOutLogout         int = 5000
	EMQXTimeOutReportProperty int = 2000
	EMQXTimeOutReadProperty   int = 2000
	EMQXTimeOutWriteProperty  int = 2000
	EMQXTimeOutService        int = 2000
	EMQXTimeOutReadNode       int = 2000
)

//实例化上报服务
var ReportServiceParamListEmqx ReportServiceParamListEmqxTemplate

func ReportServiceEmqxReadParamFromJson() bool {
	type ReportServiceConfigParamEmqxTemplate struct {
		ServiceList []ReportServiceParamEmqxTemplate `json:"ServiceList"`
	}

	configParam := ReportServiceConfigParamEmqxTemplate{}
	data, err := utils.FileRead("./selfpara/reportServiceParamListEmqx.json")
	if err != nil {
		setting.ZAPS.Debugf("上报服务[Emqx]配置json文件读取失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &configParam)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[Emqx]配置json文件格式化失败")
		return false
	}
	setting.ZAPS.Debug("上报服务[Emqx]配置json文件读取成功")

	//初始化
	for _, v := range configParam.ServiceList {
		v.GWParam.ReportStatus = "offLine"
		for _, n := range v.NodeList {
			n.CommStatus = "offLine"
			n.ReportStatus = "offLine"
		}
		ReportServiceParamListEmqx.ServiceList = append(ReportServiceParamListEmqx.ServiceList, NewReportServiceParamEmqx(v.GWParam, v.NodeList))
	}

	return true
}

func ReportServiceEmqxWriteParamToJson() {
	utils.DirIsExist("./selfpara")
	sJson, _ := json.Marshal(ReportServiceParamListEmqx)
	err := utils.FileWrite("./selfpara/reportServiceParamListEmqx.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[Emqx]配置json文件写入失败")
		return
	}
	setting.ZAPS.Debugf("上报服务[Emqx]配置json文件写入成功")
}

func (s *ReportServiceParamListEmqxTemplate) AddReportService(param ReportServiceGWParamEmqxTemplate) error {

	index := -1
	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == param.ServiceName {
			index = k
			break
		}
	}

	if index != -1 {
		return errors.New("服务名称已经存在")
	}

	nodeList := make([]ReportServiceNodeParamEmqxTemplate, 0)
	ReportServiceParam := NewReportServiceParamEmqx(param, nodeList)
	s.ServiceList = append(s.ServiceList, ReportServiceParam)

	ReportServiceEmqxWriteParamToJson()
	return nil
}

func (s *ReportServiceParamListEmqxTemplate) ModifyReportService(param ReportServiceGWParamEmqxTemplate) error {

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
	//存在相同的，表示修改;
	s.ServiceList[index].GWParam.MQTTClientOptions.SetAutoReconnect(false)
	if s.ServiceList[index].GWParam.MQTTClient != nil {
		if s.ServiceList[index].GWParam.MQTTClient.IsConnected() {
			s.ServiceList[index].GWParam.MQTTClient.Disconnect(0)
		}
	}

	//旧上报服务退出
	s.ServiceList[index].CancelFunc()
	s.ServiceList[index].GWParam = param
	ReportServiceEmqxWriteParamToJson()

	//启动新上报服务
	s.ServiceList[index] = NewReportServiceParamEmqx(param, s.ServiceList[index].NodeList)

	return nil
}

func (s *ReportServiceParamListEmqxTemplate) DeleteReportService(serviceName string) {

	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == serviceName {
			s.ServiceList = append(s.ServiceList[:k], s.ServiceList[k+1:]...)
			ReportServiceEmqxWriteParamToJson()
			//协程退出
			v.CancelFunc()
			return
		}
	}
}

func (r *ReportServiceParamEmqxTemplate) AddReportNode(param ReportServiceNodeParamEmqxTemplate) error {

	param.CommStatus = "offLine"
	param.ReportStatus = "offLine"
	param.ReportErrCnt = 0

	//节点存在则进行修改
	for _, v := range r.NodeList {
		//节点已经存在
		if v.Name == param.Name {
			return errors.New("设备名称已经存在")
		}
	}

	//节点不存在则新建
	rModel, ok := reportModel.ReportModels[param.UploadModel]
	if ok {
		param.Properties = rModel.Properties
	}
	r.NodeList = append(r.NodeList, param)
	ReportServiceEmqxWriteParamToJson()

	return nil
}

func (r *ReportServiceParamEmqxTemplate) ModifyReportNode(param ReportServiceNodeParamEmqxTemplate) error {

	param.CommStatus = "offLine"
	param.ReportStatus = "offLine"
	param.ReportErrCnt = 0

	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == param.Name {
			rModel, ok := reportModel.ReportModels[param.UploadModel]
			if ok {
				param.Properties = rModel.Properties
			}
			r.NodeList[k] = param
			ReportServiceEmqxWriteParamToJson()
			return nil
		}
	}

	//节点不存在则新建
	return errors.New("设备名称不存在")
}

func (r *ReportServiceParamEmqxTemplate) DeleteReportNode(name string) int {

	index := -1
	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == name {
			index = k
			r.NodeList = append(r.NodeList[:k], r.NodeList[k+1:]...)
			ReportServiceEmqxWriteParamToJson()
			return index
		}
	}
	return index
}

func (r *ReportServiceParamEmqxTemplate) ExportParamToCsv() (bool, string) {

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + r.GWParam.ServiceName + ".csv"

	fs, err := os.Create(fileName)
	if err != nil {
		setting.ZAPS.Errorf("创建csv文件错误 %v", err)
		return false, ""
	}
	defer fs.Close()

	//创建一个新的写入文件流
	w := csv.NewWriter(fs)
	csvData := [][]string{
		{"上报服务名称", "采集接口名称", "设备名称", "设备通信地址", "上报模型", "上报服务协议", "通信编码"},
		{"ServiceName", "CollInterfaceName", "Name", "Addr", "UploadModel", "Protocol", "DeviceName"},
	}

	for _, n := range r.NodeList {
		param := make([]string, 0)
		param = append(param, n.ServiceName)
		param = append(param, n.CollInterfaceName)
		param = append(param, n.Name)
		param = append(param, n.Addr)
		param = append(param, n.UploadModel)
		param = append(param, n.Protocol)
		param = append(param, n.Param.DeviceCode)
		csvData = append(csvData, param)
	}

	//写入数据
	err = w.WriteAll(csvData)
	if err != nil {
		setting.ZAPS.Errorf("写csv文件错误 %v", err)
		return false, ""
	}
	w.Flush()

	return true, fileName
}

func (r *ReportServiceParamEmqxTemplate) ExportParamToXlsx() (bool, string) {

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + r.GWParam.ServiceName + ".xlsx"

	cellData := [][]string{
		{"上报服务名称", "采集接口名称", "设备名称", "设备标签", "设备通信地址", "上报模型", "上报服务协议", "通信编码"},
		{"ServiceName", "CollInterfaceName", "Name", "Label", "Addr", "UploadModel", "Protocol", "DeviceName"},
	}

	for _, n := range r.NodeList {
		param := make([]string, 0)
		param = append(param, n.ServiceName)
		param = append(param, n.CollInterfaceName)
		param = append(param, n.Name)
		param = append(param, n.Label)
		param = append(param, n.Addr)
		param = append(param, n.UploadModel)
		param = append(param, n.Protocol)
		param = append(param, n.Param.DeviceCode)
		cellData = append(cellData, param)
	}

	//写入数据
	err := setting.WriteExcel(fileName, cellData)
	if err != nil {
		return false, ""
	}

	return true, fileName
}

func (r *ReportServiceParamEmqxTemplate) ProcessUpLinkFrame(ctx context.Context) {

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
				r.ReportServiceEmqxProcessWriteProperty(reqFrame)
			}
		case reqFrame := <-r.ReceiveReadPropertyRequestFrameChan:
			{
				r.ReportServiceEmqxProcessReadProperty(reqFrame)
			}
		case reqFrame := <-r.ReceiveInvokeServiceRequestFrameChan:
			{
				r.ReportServiceEmqxProcessInvokeService(reqFrame)
			}
		case reqFrame := <-r.ReceiveReadNodeStatusRequestFrameChan:
			{
				r.ReportServiceEmqxProcessReadNodeStatus(reqFrame)
			}
		case reqFrame := <-r.ReceiveInvokeGWServiceRequestFrameChan:
			{
				r.ReportServiceEmqxProcessInvokeGWService(reqFrame)
			}
		}
	}
}

func (r *ReportServiceParamEmqxTemplate) ReportNode(name string) {

	if r.GWParam.ReportStatus == "onLine" {
		nodeName := make([]string, 0)
		nodeName = append(nodeName, name)
		setting.ZAPS.Debugf("上报服务[%s]手动设置节点上报%v", r.GWParam.ServiceName, nodeName)
		if len(nodeName) > 0 {
			r.NodePropertyPost(nodeName, nil, false)
		}
	}
}

func (r *ReportServiceParamEmqxTemplate) ProcessDownLinkFrame(ctx context.Context) {

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

					ackFrame := MQTTEmqxReportPropertyAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.ZAPS.Errorf("ReportPropertyAck json unmarshal err")
						continue
					}
					r.ReceiveReportPropertyAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/node/login/post_reply") { //子设备上线回应

					ackFrame := MQTTEmqxLogInAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.ZAPS.Warnf("LogInAck json unmarshal err")
						continue
					}
					r.ReceiveLogInAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/node/logout/post_reply") { //子设备下线回应

					ackFrame := MQTTEmqxLogOutAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.ZAPS.Errorf("LogOutAck json unmarshal err")
						continue
					}
					r.ReceiveLogOutAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/node/service/invoke") { //设备服务调用
					serviceFrame := MQTTEmqxInvokeServiceRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &serviceFrame)
					if err != nil {
						setting.ZAPS.Errorf("serviceFrame json unmarshal err")
						continue
					}
					r.ReceiveInvokeServiceRequestFrameChan <- serviceFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/node/property/set") { //子设备设置属性请求
					writePropertyFrame := MQTTEmqxWritePropertyRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &writePropertyFrame)
					if err != nil {
						setting.ZAPS.Errorf("writePropertyFrame json unmarshal err")
						continue
					}
					r.ReceiveWritePropertyRequestFrameChan <- writePropertyFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/node/property/get") { //子设备获取属性请求
					readPropertyFrame := MQTTEmqxReadPropertyRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &readPropertyFrame)
					if err != nil {
						setting.ZAPS.Errorf("readPropertyFrame json unmarshal err")
						continue
					}
					r.ReceiveReadPropertyRequestFrameChan <- readPropertyFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/node/status/get") { //子设备状态查询
					readNodeStatusFrame := MQTTEmqxReadNodeStatusRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &readNodeStatusFrame)
					if err != nil {
						setting.ZAPS.Errorf("readNodeStatusFrame json unmarshal err")
						continue
					}
					r.ReceiveReadNodeStatusRequestFrameChan <- readNodeStatusFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/gw/service/invoke") { //网关服务调用
					serviceFrame := MQTTEmqxInvokeServiceRequestTemplate{}
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

func (r *ReportServiceParamEmqxTemplate) ProcessCollEvent(ctx context.Context, sub eventBus.Sub) {
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
						if len(r.LogInRequestFrameChan) == 50 {
							<-r.LogInRequestFrameChan
						}
						r.LogInRequestFrameChan <- nodeName

						//判断告警
						r.ProcessAlarmEvent(index, subMsg.CollName, subMsg.NodeName)
					}
				case "offLine":
					{
						nodeName = append(nodeName, subMsg.NodeName)
						r.NodeList[index].CommStatus = "offLine"
						if len(r.LogOutRequestFrameChan) == 50 {
							<-r.LogOutRequestFrameChan
						}
						r.LogOutRequestFrameChan <- nodeName
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

func (r *ReportServiceParamEmqxTemplate) ProcessAlarmEvent(index int, collName string, nodeName string) {

	reportStatus := false

	properties := make([]MQTTEmqxPropertyPostParamPropertyTemplate, 0)

	coll, ok := device.CollectInterfaceMap.Coll[collName]
	if !ok {
		return
	}
	node, ok := coll.DeviceNodeMap[nodeName]
	if !ok {
		return
	}

	for _, v := range node.Properties {
		//从上报模型中查找属性
		rProperty, ok := r.NodeList[index].Properties[v.Name]
		if !ok {
			continue
		}
		if rProperty.Params.StepAlarm == true {
			valueCnt := len(v.Value)
			if valueCnt >= 2 { //阶跃报警必须是2个值
				switch v.Type {
				case device.PropertyTypeInt32:
					{
						pValueCur := v.Value[valueCnt-1].Value.(int32)
						pValuePre := v.Value[valueCnt-2].Value.(int32)
						step, _ := strconv.Atoi(rProperty.Params.Step)
						if math.Abs(float64(pValueCur-pValuePre)) > float64(step) {
							reportStatus = true //满足报警条件，上报
							setting.ZAPS.Infof("设备[%v]阶跃报警", r.NodeList[index].Name)
							//转换时间
							//timeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", v.Value[valueCnt-1].TimeStamp, time.Local)
							property := MQTTEmqxPropertyPostParamPropertyTemplate{
								Name:      v.Name,
								Value:     v.Value[valueCnt-1].Value.(int32),
								TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
							}
							properties = append(properties, property)
							continue
						}
					}
				case device.PropertyTypeUInt32:
					{
						pValueCur := v.Value[valueCnt-1].Value.(uint32)
						pValuePre := v.Value[valueCnt-2].Value.(uint32)
						step, _ := strconv.Atoi(rProperty.Params.Step)
						if math.Abs(float64(pValueCur-pValuePre)) > float64(step) {
							reportStatus = true //满足报警条件，上报
							setting.ZAPS.Infof("设备[%v]阶跃报警", r.NodeList[index].Name)
							//转换时间
							//timeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", v.Value[valueCnt-1].TimeStamp, time.Local)
							property := MQTTEmqxPropertyPostParamPropertyTemplate{
								Name:      v.Name,
								Value:     v.Value[valueCnt-1].Value.(uint32),
								TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
							}
							properties = append(properties, property)
							continue
						}
					}
				case device.PropertyTypeDouble:
					{
						pValueCur := v.Value[valueCnt-1].Value.(float64)
						pValuePre := v.Value[valueCnt-2].Value.(float64)
						step, err := strconv.ParseFloat(rProperty.Params.Step, 64)
						if err != nil {
							continue
						}
						if math.Abs(pValueCur-pValuePre) > float64(step) {
							reportStatus = true //满足报警条件，上报
							setting.ZAPS.Infof("设备[%v]阶跃报警", r.NodeList[index].Name)
							//转换时间
							//timeStamp, _ := time.ParseInLocation("2006-01-02 15:04:05", v.Value[valueCnt-1].TimeStamp, time.Local)
							property := MQTTEmqxPropertyPostParamPropertyTemplate{
								Name:      v.Name,
								Value:     v.Value[valueCnt-1].Value.(float64),
								TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
							}
							properties = append(properties, property)
							continue
						}
					}
				}
			}
		}
		if rProperty.Params.MinMaxAlarm == true {
			valueCnt := len(v.Value)
			if v.Type == device.PropertyTypeInt32 {
				pValueCur := v.Value[valueCnt-1].Value.(int32)
				min, _ := strconv.Atoi(rProperty.Params.Min)
				max, _ := strconv.Atoi(rProperty.Params.Max)
				if pValueCur < int32(min) || pValueCur > int32(max) {
					reportStatus = true //满足报警条件，上报
					setting.ZAPS.Infof("设备[%v]范围报警", r.NodeList[index].Name)
					property := MQTTEmqxPropertyPostParamPropertyTemplate{
						Name:      v.Name,
						Value:     v.Value[valueCnt-1].Value.(int32),
						TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
					}
					properties = append(properties, property)
				}
			} else if v.Type == device.PropertyTypeUInt32 {
				pValueCur := v.Value[valueCnt-1].Value.(uint32)
				min, _ := strconv.Atoi(rProperty.Params.Min)
				max, _ := strconv.Atoi(rProperty.Params.Max)
				if pValueCur < uint32(min) || pValueCur > uint32(max) {
					reportStatus = true //满足报警条件，上报
					setting.ZAPS.Infof("设备[%v]范围报警", r.NodeList[index].Name)
					property := MQTTEmqxPropertyPostParamPropertyTemplate{
						Name:      v.Name,
						Value:     v.Value[valueCnt-1].Value.(uint32),
						TimeStamp: v.Value[valueCnt-1].TimeStamp.Unix(),
					}
					properties = append(properties, property)
				}
			} else if v.Type == device.PropertyTypeDouble {
				pValueCur := v.Value[valueCnt-1].Value.(float64)
				min, err := strconv.ParseFloat(rProperty.Params.Min, 64)
				if err != nil {
					continue
				}
				max, err := strconv.ParseFloat(rProperty.Params.Max, 64)
				if err != nil {
					continue
				}
				if pValueCur < min || pValueCur > max {
					reportStatus = true //满足报警条件，上报
					setting.ZAPS.Infof("设备[%v]范围报警", r.NodeList[index].Name)
					property := MQTTEmqxPropertyPostParamPropertyTemplate{
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
		reportNodeProperty := MQTTEmqxReportAlarmTemplate{
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

func (r *ReportServiceParamEmqxTemplate) LogIn(nodeName []string) {

	//清空接收chan，避免出现有上次接收的缓存
	for i := 0; i < len(r.ReceiveLogInAckFrameChan); i++ {
		<-r.ReceiveLogInAckFrameChan
	}

	r.NodeLogIn(nodeName)
}

func (r *ReportServiceParamEmqxTemplate) LogOut(nodeName []string) {

	//清空接收chan，避免出现有上次接收的缓存
	for i := 0; i < len(r.ReceiveLogOutAckFrameChan); i++ {
		<-r.ReceiveLogOutAckFrameChan
	}

	r.NodeLogOut(nodeName)
}

func (r *ReportServiceParamEmqxTemplate) ReportTimeFun() {

	//网关上报
	reportGWProperty := MQTTEmqxReportPropertyTemplate{
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
			reportNodeProperty := MQTTEmqxReportPropertyTemplate{
				DeviceType: "node",
				DeviceName: nodeName,
			}
			r.ReportPropertyRequestFrameChan <- reportNodeProperty
		}
	}
}

func ReportServiceEmqxInit() {
	ReportServiceEmqxReadParamFromJson()
}

func NewReportServiceParamEmqx(gw ReportServiceGWParamEmqxTemplate, nodeList []ReportServiceNodeParamEmqxTemplate) *ReportServiceParamEmqxTemplate {
	for _, v := range nodeList {
		rModel, ok := reportModel.ReportModels[v.UploadModel]
		if !ok {
			continue
		}
		v.Properties = rModel.Properties
	}

	emqxParam := &ReportServiceParamEmqxTemplate{
		GWParam:                                gw,
		NodeList:                               nodeList,
		ReceiveFrameChan:                       make(chan MQTTEmqxReceiveFrameTemplate, 100),
		LogInRequestFrameChan:                  make(chan []string, 50),
		ReceiveLogInAckFrameChan:               make(chan MQTTEmqxLogInAckTemplate, 5),
		LogOutRequestFrameChan:                 make(chan []string, 50),
		ReceiveLogOutAckFrameChan:              make(chan MQTTEmqxLogOutAckTemplate, 5),
		ReportPropertyRequestFrameChan:         make(chan MQTTEmqxReportPropertyTemplate, 50),
		ReceiveReportPropertyAckFrameChan:      make(chan MQTTEmqxReportPropertyAckTemplate, 50),
		ReportAlarmRequestFrameChan:            make(chan MQTTEmqxReportAlarmTemplate, 50),
		ReceiveReportAlarmAckFrameChan:         make(chan MQTTEmqxReportPropertyAckTemplate, 50),
		ReceiveInvokeServiceRequestFrameChan:   make(chan MQTTEmqxInvokeServiceRequestTemplate, 50),
		ReceiveInvokeServiceAckFrameChan:       make(chan MQTTEmqxInvokeServiceAckTemplate, 50),
		ReceiveWritePropertyRequestFrameChan:   make(chan MQTTEmqxWritePropertyRequestTemplate, 50),
		ReceiveReadPropertyRequestFrameChan:    make(chan MQTTEmqxReadPropertyRequestTemplate, 50),
		ReceiveReadNodeStatusRequestFrameChan:  make(chan MQTTEmqxReadNodeStatusRequestTemplate, 50),
		ReceiveInvokeGWServiceRequestFrameChan: make(chan MQTTEmqxInvokeServiceRequestTemplate, 50),
		ReceiveInvokeGWServiceAckFrameChan:     make(chan MQTTEmqxInvokeServiceAckTemplate, 50),
	}

	ctx, cancel := context.WithCancel(context.Background())
	emqxParam.CancelFunc = cancel

	go ReportServiceEmqxPoll(ctx, emqxParam)

	return emqxParam
}

func ReportServiceEmqxPoll(ctx context.Context, r *ReportServiceParamEmqxTemplate) {

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
	//处理采集接口的消息
	go r.ProcessCollEvent(ctx, r.EventSub)

	go r.ProcessUpLinkFrame(ctx)
	go r.ProcessDownLinkFrame(ctx)

	_ = r.GWLogin()
	scheduler.Start()
}
