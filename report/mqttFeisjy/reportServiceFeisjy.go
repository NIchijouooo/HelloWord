package mqttFeisjy

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/device"
	"gateway/device/eventBus"
	"gateway/report/reportModel"
	"gateway/setting"
	"gateway/utils"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/robfig/cron/v3"
	"os"
	"strings"
	"time"
)

// 上报节点参数结构体
type ReportServiceNodeParamFeisjyTemplate struct {
	ServiceName       string `json:"serviceName"`
	CollInterfaceName string `json:"collInterfaceName"`
	Name              string `json:"deviceName"`
	Label             string `json:"Label"`
	Addr              string `json:"deviceAddr"`
	CommStatus        string
	ReportErrCnt      int `json:"-"`
	ReportStatus      string
	HeartBeatMark     bool                                                `json:"-"` //定时上报时间到了，是否需要推送心跳，如果定时过程中有其他命令已经发过数据了，就不再发心跳
	HeartBeatCnt      uint32                                              `json:"-"`
	UploadModel       string                                              `json:"uploadModel"`
	Protocol          string                                              `json:"protocol"`
	Properties        map[string]*reportModel.ReportModelPropertyTemplate `json:"properties"`
	Param             struct {
		ProductKey   string
		DeviceID     string
		DeviceSecret string
	}
}

// 上报网关参数结构体
type ReportServiceGWParamFeisjyTemplate struct {
	ServiceName  string
	ReportNetSW  bool
	ReportNet    string
	IP           string
	Port         string
	ReportStatus string
	ReportTime   int
	ReportErrCnt int
	Protocol     string
	Param        struct {
		UserName     string
		Password     string
		ClientID     string //MQTT客户端ID，用DeviceID+一个随机数(1000内)
		AppKey       string
		ProductKey   string
		DeviceID     string
		DeviceSecret string
	}
	MQTTClient        MQTT.Client         `json:"-"`
	MQTTClientOptions *MQTT.ClientOptions `json:"-"`
	OfflineChan       chan bool           `json:"-"`
	HeartBeatMark     bool                `json:"-"` //定时上报时间到了，是否需要推送心跳，如果定时过程中有其他命令已经发过数据了，就不再发心跳
	HeartBeatCnt      uint32              `json:"-"`
}

// 上报服务参数，网关参数，节点参数
type ReportServiceParamFeisjyTemplate struct {
	GWParam                        ReportServiceGWParamFeisjyTemplate
	NodeList                       []ReportServiceNodeParamFeisjyTemplate
	ReceiveFrameChan               chan MQTTFeisjyReceiveFrameTemplate   `json:"-"`
	ReceiveLogInAckFrameChan       chan MQTTFeisjyLogInAckTemplate       `json:"-"` //平台下发的登陆确认报文
	ReportPropertyRequestFrameChan chan MQTTFeisjyReportPropertyTemplate `json:"-"`
	ReceiveDevUpGradeChan          chan MQTTFeisjyUpGradeTemplate        `json:"-"` // 平台下发固件升级
	ReceiveFileListChan            chan FileListFeisjyTemplate           `json:"-"` // FileList相关
	CancelFunc                     context.CancelFunc                    `json:"-"`
	MessageEventBus                eventBus.Bus                          `json:"-"` //通信报文总线
}

type ReportServiceParamListFeisjyTemplate struct {
	ServiceList []*ReportServiceParamFeisjyTemplate
}

// 实例化上报服务
var ReportServiceParamListFeisjy ReportServiceParamListFeisjyTemplate

func ReportServiceFeisjyInit() {
	ReportServiceFeisjyReadParamFromJson()
}

func NewReportServiceParamFeisjy(gw ReportServiceGWParamFeisjyTemplate, nodeList []ReportServiceNodeParamFeisjyTemplate) *ReportServiceParamFeisjyTemplate {
	feisjyParam := &ReportServiceParamFeisjyTemplate{
		GWParam:                        gw,
		NodeList:                       nodeList,
		ReceiveFrameChan:               make(chan MQTTFeisjyReceiveFrameTemplate, 100),
		ReceiveLogInAckFrameChan:       make(chan MQTTFeisjyLogInAckTemplate, 2),
		ReportPropertyRequestFrameChan: make(chan MQTTFeisjyReportPropertyTemplate, 50),
		ReceiveDevUpGradeChan:          make(chan MQTTFeisjyUpGradeTemplate, 2),
		ReceiveFileListChan:            make(chan FileListFeisjyTemplate, 50),
		MessageEventBus:                eventBus.NewBus(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	feisjyParam.CancelFunc = cancel

	go ReportServiceFeisjyPoll(ctx, feisjyParam)

	return feisjyParam
}

func ReportServiceFeisjyPoll(ctx context.Context, r *ReportServiceParamFeisjyTemplate) {
	reportState := 0

	// 定义一个cron运行器
	cronProcess := cron.New()

	reportTime := fmt.Sprintf("@every %dm%ds", r.GWParam.ReportTime/60, r.GWParam.ReportTime%60)
	setting.ZAPS.Infof("上报服务[%s]定时上报周期为%v", r.GWParam.ServiceName, reportTime)

	cronProcess.AddFunc(reportTime, r.ReportTimeFun)

	//订阅采集接口消息
	device.CollectInterfaceMap.Lock.Lock()
	for _, coll := range device.CollectInterfaceMap.Coll {
		sub := eventBus.NewSub()
		coll.CollEventBus.Subscribe("onLine", sub)
		coll.CollEventBus.Subscribe("offLine", sub)
		coll.CollEventBus.Subscribe("update", sub)
		go r.ProcessCollEvent(ctx, sub)
	}
	device.CollectInterfaceMap.Lock.Unlock()

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
				switch reportState {
				case 0: //  配置MQTT,连接MQTT BREAK,
					if r.GWLogin() == true {
						reportState = 1
					} else {
						time.Sleep(120 * time.Second)
					}

				case 1: //  网关登录流程
					for i, _ := range r.NodeList {
						r.NodeList[i].ReportStatus = "offLine" //网关未登陆，节点配置为离线状态
					}
					if r.FeisjyLogInMachine(r.GWParam.Param.DeviceID) == true {
						r.GWParam.ReportStatus = "onLine"
						cronProcess.Start()
						reportState = 2
						/**
						todo20230615网关设备上线更新sqlite
						*/
						//NewRealtimeDataRepository().UpdateGatewayDeviceConnetStatus(r.GWParam)

					} else {
						time.Sleep(10 * time.Second)
					}
					break

				case 2: //  节点登录流程
					if r.GWParam.ReportStatus == "onLine" {
						for i, node := range r.NodeList {
							if node.ReportStatus == "offLine" {
								r.MQTTFeisjySubNodeTopic(node.Param.DeviceID)
								if r.FeisjyLogInMachine(node.Param.DeviceID) == true {
									r.NodeList[i].ReportStatus = "onLine"
									/**
									todo20230615设备上线更新sqlite
									*/
									//NewRealtimeDataRepository().UpdateDeviceConnetStatus(r.NodeList[i])

								}
							}
						}
						time.Sleep(30 * time.Second)
					} else {
						reportState = 1
					}
					break
				}

				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (r *ReportServiceParamFeisjyTemplate) ReportTimeFun() {

	if r.GWParam.ReportStatus == "offLine" {
		return
	}

	//发送网关心跳到平台
	if r.GWParam.HeartBeatMark == false { //如果有其它数据推送到平台了，可以不发心路
		r.GWParam.HeartBeatCnt++
		r.FeisjyHeartBeat(r.GWParam.Param.DeviceID, r.GWParam.HeartBeatCnt)
	}
	r.GWParam.HeartBeatMark = false

	//发送节点心跳到平台
	for i, node := range r.NodeList {
		if node.ReportStatus == "onLine" {
			if node.HeartBeatMark == false { //如果有其它数据推送到平台了，可以不发心路
				r.NodeList[i].HeartBeatCnt++
				r.FeisjyHeartBeat(node.Param.DeviceID, r.NodeList[i].HeartBeatCnt)
			}
			r.NodeList[i].HeartBeatMark = false
		}
	}
}

func (r *ReportServiceParamFeisjyTemplate) ProcessCollEvent(ctx context.Context, sub eventBus.Sub) {
	for {
		select {
		case <-ctx.Done():
			{
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
						//判断告警
						r.ProcessAlarmEvent(index, subMsg.CollName, subMsg.NodeName)
					}
				case "offLine":
					{
						nodeName = append(nodeName, subMsg.NodeName)
						r.NodeList[index].CommStatus = "offLine"
					}
				case "update":
					{
						//更新设备的通信状态
						r.NodeList[index].CommStatus = "onLine"
						//判断告警
						r.ProcessAlarmEvent(index, subMsg.CollName, subMsg.NodeName)
					}
				}

				//if subMsg.Topic == "onLine" || subMsg.Topic == "offLine" || subMsg.Topic == "update" {
				//	{
				//		nodeName = append(nodeName, subMsg.NodeName)
				//		if r.NodeList[index].CommStatus == "onLine" {
				//			reportNodeProperty := MQTTFeisjyReportPropertyTemplate{
				//				DeviceType: "node",
				//				DeviceName: nodeName,
				//			}
				//			r.ReportPropertyRequestFrameChan <- reportNodeProperty
				//		}
				//	}
				//}
			}
		}
	}
}

func (r *ReportServiceParamFeisjyTemplate) ProcessUpLinkFrame(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case reqFrame := <-r.ReportPropertyRequestFrameChan:
			{
				if reqFrame.DeviceType == "gw" {
					r.GWPropertyPost()
				} else if reqFrame.DeviceType == "node" {
					r.NodePropertyPost(reqFrame.DeviceName)
				}
			}
		case reqFrame := <-r.ReceiveDevUpGradeChan:
			{
				r.FeisjyDevUpGradeaMachine(&reqFrame)
			}
		case reqFrame := <-r.ReceiveFileListChan:
			{
				r.FeisjyFileListMachine(reqFrame)
			}
		}
	}
}

func (r *ReportServiceParamFeisjyTemplate) ProcessDownLinkFrame(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case frame := <-r.ReceiveFrameChan:
			{
				parts := strings.Split(frame.Topic, "/") //字符串按照指定的分隔符进行分割
				if len(parts) < 5 {
					continue
				}

				ID := parts[3]
				Command := parts[4]
				setting.ZAPS.Infof("DownLinkFrame ID[%s] Command[%s]", ID, Command)

				switch Command {
				case "login": //平台下发的登录相关报文
					{
						ackFrame := MQTTFeisjyLogInAckTemplate{}
						err := json.Unmarshal(frame.Payload, &ackFrame)
						if err != nil {
							setting.ZAPS.Errorf("LogIn Ack json unmarshal err")
							continue
						}
						r.ReceiveLogInAckFrameChan <- ackFrame
					}

				case "deviceControl": //平台下发的deviceControl报文
					{
						type CmdItemsFeisjyTemplate struct {
							Code  int         `json:"code"`
							Value interface{} `json:"value"`
						}

						var deviceControlData struct {
							Uuid       string                   `json:"uuid"`
							DeviceAddr string                   `json:"deviceAddr"`
							CmdType    string                   `json:"cmdType"`
							Time       int64                    `json:"time"`
							CmdItems   []CmdItemsFeisjyTemplate `json:"cmdItems"`
						}

						setting.ZAPS.Infof("%s", frame.Payload)
						err := json.Unmarshal(frame.Payload, &deviceControlData)
						if err != nil {
							setting.ZAPS.Error("Unmarshal deviceControl err", err)
							continue
						}

						reqFrame := MQTTFeisjyWritePropertyTemplate{
							CmdType:    deviceControlData.CmdType,
							Uuid:       deviceControlData.Uuid,
							DeviceAddr: deviceControlData.DeviceAddr,
							Properties: nil,
						}

						for _, v := range deviceControlData.CmdItems {
							reqFrame.Properties = append(reqFrame.Properties, MQTTFeisjyWritePropertyRequestParamPropertyTemplate{fmt.Sprintf("%d", v.Code), v.Value})
						}

						r.FeisjyDeviceControlMachine(reqFrame)
					}

				case "deviceUpgrade":
					{
						deviceUpgradeData := MQTTFeisjyUpGradeTemplate{}

						err := json.Unmarshal(frame.Payload, &deviceUpgradeData)
						if err != nil {
							setting.ZAPS.Error("Unmarshal deviceUpgrade err", err)
							continue
						}

						setting.ZAPS.Info("发送升级信息...")
						r.ReceiveDevUpGradeChan <- deviceUpgradeData
					}

				case "fileList": //获取文件列表
					{
						fileListData := FileListFeisjyTemplate{}

						err := json.Unmarshal(frame.Payload, &fileListData)
						if err != nil {
							setting.ZAPS.Error("Unmarshal deviceUpgrade err", err)
							continue
						}

						r.ReceiveFileListChan <- fileListData
					}

				case "resultMsg":
					{
						var resultMsgData struct {
							Topic string `json:"topic"`
							Data  string `json:"data"`
							Msg   string `json:"msg"`
						}
						err := json.Unmarshal(frame.Payload, &resultMsgData)
						if err != nil {
							setting.ZAPS.Error("Unmarshal deviceControl err", err)
							continue
						}

						if resultMsgData.Msg == "no_login" || resultMsgData.Msg == "no_login_confirm" {
							if ID == r.GWParam.Param.DeviceID {
								r.GWLogin()
							} else {
								for k, v := range r.NodeList {
									if v.Param.DeviceID == ID {
										r.NodeList[k].ReportStatus = "offLine"
										/**
										todo20230615设备离线更新sqlite
										*/
										//NewRealtimeDataRepository().UpdateDeviceConnetStatus(r.NodeList[k])
									}
								}
							}
						}
						if resultMsgData.Msg == "device_offline" {
							if ID == r.GWParam.Param.DeviceID {
								//r.GWLogin()
							} else {
								for k, v := range r.NodeList {
									if v.Param.DeviceID == ID {
										r.NodeList[k].ReportStatus = "offLine"
										/**
										todo20230615设备离线更新sqlite
										*/
										//NewRealtimeDataRepository().UpdateDeviceConnetStatus(r.NodeList[k])

									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func ReportServiceFeisjyReadParamFromJson() bool {
	type ReportServiceConfigParamFeisjyTemplate struct {
		ServiceList []ReportServiceParamFeisjyTemplate `json:"ServiceList"`
	}

	configParam := ReportServiceConfigParamFeisjyTemplate{}
	data, err := utils.FileRead("./selfpara/reportServiceParamListFeisjy.json")
	if err != nil {
		setting.ZAPS.Debugf("上报服务[Feisjy]配置json文件读取失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &configParam)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[Emqx]配置json文件格式化失败")
		return false
	}
	setting.ZAPS.Debug("上报服务[Feisjy]配置json文件读取成功")

	//初始化
	for _, v := range configParam.ServiceList {
		ReportServiceParamListFeisjy.ServiceList = append(ReportServiceParamListFeisjy.ServiceList, NewReportServiceParamFeisjy(v.GWParam, v.NodeList))
	}

	return true
}

func ReportServiceFeisjyWriteParamToJson() {
	utils.DirIsExist("./selfpara")
	sJson, _ := json.Marshal(ReportServiceParamListFeisjy)
	err := utils.FileWrite("./selfpara/reportServiceParamListFeisjy.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[Feisjy]配置json文件写入失败")
		return
	}
	setting.ZAPS.Debugf("上报服务[Feisjy]配置json文件写入成功")
}

func (s *ReportServiceParamListFeisjyTemplate) AddReportService(param ReportServiceGWParamFeisjyTemplate) error {

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

	nodeList := make([]ReportServiceNodeParamFeisjyTemplate, 0)
	ReportServiceParam := NewReportServiceParamFeisjy(param, nodeList)
	s.ServiceList = append(s.ServiceList, ReportServiceParam)

	ReportServiceFeisjyWriteParamToJson()
	return nil
}

func (s *ReportServiceParamListFeisjyTemplate) ModifyReportService(param ReportServiceGWParamFeisjyTemplate) error {

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
	ReportServiceFeisjyWriteParamToJson()

	//启动新上报服务
	s.ServiceList[index] = NewReportServiceParamFeisjy(param, s.ServiceList[index].NodeList)

	return nil
}
func (s *ReportServiceParamListFeisjyTemplate) DeleteReportService(serviceName string) {

	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == serviceName {
			s.ServiceList = append(s.ServiceList[:k], s.ServiceList[k+1:]...)
			ReportServiceFeisjyWriteParamToJson()
			//协程退出
			v.CancelFunc()
			return
		}
	}
}

func (r *ReportServiceParamFeisjyTemplate) AddReportNode(param ReportServiceNodeParamFeisjyTemplate) error {

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
	param.ReportStatus = "offLine"
	r.NodeList = append(r.NodeList, param)
	ReportServiceFeisjyWriteParamToJson()

	return nil
}
func (r *ReportServiceParamFeisjyTemplate) ModifyReportNode(param ReportServiceNodeParamFeisjyTemplate) error {

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
			param.ReportStatus = "offLine"
			r.NodeList[k] = param
			ReportServiceFeisjyWriteParamToJson()

			return nil
		}
	}

	//节点不存在则新建
	return errors.New("设备名称不存在")
}

func (r *ReportServiceParamFeisjyTemplate) DeleteReportNode(name string) int {

	index := -1
	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == name {
			index = k

			//解除订阅主题
			r.MQTTFeisjyUnsubNodeTopic(v.Param.DeviceID)

			r.NodeList = append(r.NodeList[:k], r.NodeList[k+1:]...)
			ReportServiceFeisjyWriteParamToJson()
			return index
		}
	}
	return index
}

func (r *ReportServiceParamFeisjyTemplate) ExportParamToCsv() (bool, string) {

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
		{"上报服务名称", "采集接口名称", "设备名称", "设备通信地址", "上报模型", "上报服务协议", "通信编码", "产品密钥", "平台设备地址", "设备密钥"},
		{"ServiceName", "CollInterfaceName", "Name", "Addr", "UploadModel", "Protocol", "DeviceName", "ProductKey", "DeviceID", "DeviceSecret"},
	}

	for _, n := range r.NodeList {
		param := make([]string, 0)
		param = append(param, n.ServiceName)
		param = append(param, n.CollInterfaceName)
		param = append(param, n.Name)
		param = append(param, n.Addr)
		param = append(param, n.UploadModel)
		param = append(param, n.Protocol)
		param = append(param, n.Param.ProductKey)
		param = append(param, n.Param.DeviceID)
		param = append(param, n.Param.DeviceSecret)
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

func (r *ReportServiceParamFeisjyTemplate) ExportParamToXlsx() (bool, string) {

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + r.GWParam.ServiceName + ".xlsx"

	cellData := [][]string{
		{"上报服务名称", "采集接口名称", "设备名称", "设备通信地址", "上报模型", "上报服务协议", "通信编码", "产品密钥", "平台设备地址", "设备密钥"},
		{"ServiceName", "CollInterfaceName", "Name", "Addr", "UploadModel", "Protocol", "DeviceName", "ProductKey", "DeviceID", "DeviceSecret"},
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
		param = append(param, n.Param.ProductKey)
		param = append(param, n.Param.DeviceID)
		param = append(param, n.Param.DeviceSecret)
		cellData = append(cellData, param)
	}

	//写入数据
	err := setting.WriteExcel(fileName, cellData)
	if err != nil {
		return false, ""
	}

	return true, fileName
}
