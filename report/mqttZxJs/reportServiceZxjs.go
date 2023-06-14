package mqttZxJs

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
)

// 上报节点参数结构体
type ReportServiceNodeParamZxjsTemplate struct {
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
		ProductSn    string
		DeviceSn     string
		DeviceSecret string
	}
}

// 上报网关参数结构体
type ReportServiceGWParamZxjsTemplate struct {
	ServiceName  string
	IP           string
	Port         string
	ReportStatus string
	ReportTime   int
	ReportErrCnt int
	Protocol     string
	Param        struct {
		ProductSn string
		DeviceSn  string
		DevicePwd string
	}
	MQTTClient        MQTT.Client         `json:"-"`
	MQTTClientOptions *MQTT.ClientOptions `json:"-"`
	OfflineChan       chan bool           `json:"-"`
	HeartBeatMark     bool                `json:"-"` //定时上报时间到了，是否需要推送心跳，如果定时过程中有其他命令已经发过数据了，就不再发心跳
	HeartBeatCnt      uint32              `json:"-"`
}

// 上报服务参数，网关参数，节点参数
type ReportServiceParamZxjsTemplate struct {
	GWParam                        ReportServiceGWParamZxjsTemplate
	NodeList                       []ReportServiceNodeParamZxjsTemplate
	ReceiveFrameChan               chan MQTTZxjsReceiveFrameTemplate   `json:"-"`
	ReceiveLogInAckFrameChan       chan MQTTZxjsLogInAckTemplate       `json:"-"` //平台下发的登陆确认报文
	ReportPropertyRequestFrameChan chan MQTTZxjsReportPropertyTemplate `json:"-"`
	CancelFunc                     context.CancelFunc                  `json:"-"`
}

type ReportServiceParamListZxjsTemplate struct {
	ServiceList []*ReportServiceParamZxjsTemplate
}

// 实例化上报服务
var ReportServiceParamListZxjs ReportServiceParamListZxjsTemplate

func ReportServiceZxjsInit() {
	ReportServiceZxjsReadParamFromJson()
}

func NewReportServiceParamZxjs(gw ReportServiceGWParamZxjsTemplate, nodeList []ReportServiceNodeParamZxjsTemplate) *ReportServiceParamZxjsTemplate {
	ZxjsParam := &ReportServiceParamZxjsTemplate{
		GWParam:                        gw,
		NodeList:                       nodeList,
		ReceiveFrameChan:               make(chan MQTTZxjsReceiveFrameTemplate, 100),
		ReceiveLogInAckFrameChan:       make(chan MQTTZxjsLogInAckTemplate, 2),
		ReportPropertyRequestFrameChan: make(chan MQTTZxjsReportPropertyTemplate, 50),
	}

	ctx, cancel := context.WithCancel(context.Background())
	ZxjsParam.CancelFunc = cancel

	go ReportServiceZxjsPoll(ctx, ZxjsParam)
	return ZxjsParam
}

func ReportServiceZxjsPoll(ctx context.Context, r *ReportServiceParamZxjsTemplate) {

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
	r.GWLogin()

}

func (r *ReportServiceParamZxjsTemplate) ReportTimeFun() {

	if r.GWParam.ReportStatus == "offLine" {
		return
	}

}

func (r *ReportServiceParamZxjsTemplate) ProcessCollEvent(ctx context.Context, sub eventBus.Sub) {
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

				if subMsg.Topic == "onLine" || subMsg.Topic == "offLine" || subMsg.Topic == "update" {
					{
						nodeName = append(nodeName, subMsg.NodeName)
						if r.NodeList[index].CommStatus == "onLine" {
							reportNodeProperty := MQTTZxjsReportPropertyTemplate{
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

func (r *ReportServiceParamZxjsTemplate) ProcessUpLinkFrame(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case reqFrame := <-r.ReportPropertyRequestFrameChan:
			{
				if reqFrame.DeviceType == "node" {
					r.NodePropertyPost(reqFrame.DeviceName)
				}
			}
		}
	}
}

func (r *ReportServiceParamZxjsTemplate) ProcessDownLinkFrame(ctx context.Context) {

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
				case "deviceControl": //平台下发的deviceControl报文
					{
						type CmdItemsZxjsTemplate struct {
							Code  int         `json:"code"`
							Value interface{} `json:"value"`
						}

						var deviceControlData struct {
							Uuid       string                 `json:"uuid"`
							DeviceAddr string                 `json:"deviceAddr"`
							CmdType    string                 `json:"cmdType"`
							Time       int64                  `json:"time"`
							CmdItems   []CmdItemsZxjsTemplate `json:"cmdItems"`
						}

						setting.ZAPS.Infof("%s", frame.Payload)
						err := json.Unmarshal(frame.Payload, &deviceControlData)
						if err != nil {
							setting.ZAPS.Error("Unmarshal deviceControl err", err)
							continue
						}

						reqFrame := MQTTZxjsWritePropertyTemplate{
							CmdType:    deviceControlData.CmdType,
							Uuid:       deviceControlData.Uuid,
							DeviceAddr: deviceControlData.DeviceAddr,
							Properties: nil,
						}

						for _, v := range deviceControlData.CmdItems {
							reqFrame.Properties = append(reqFrame.Properties, MQTTZxjsWritePropertyRequestParamPropertyTemplate{fmt.Sprintf("%d", v.Code), v.Value})
						}

						r.ZxjsDeviceControlMachine(reqFrame)
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
							if ID == r.GWParam.Param.DeviceSn {
								r.GWLogin()
							} else {
								for k, v := range r.NodeList {
									if v.Param.DeviceSn == ID {
										r.NodeList[k].ReportStatus = "offLine"
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

func ReportServiceZxjsReadParamFromJson() bool {
	type ReportServiceConfigParamZxjsTemplate struct {
		ServiceList []ReportServiceParamZxjsTemplate `json:"ServiceList"`
	}

	configParam := ReportServiceConfigParamZxjsTemplate{}
	data, err := utils.FileRead("./selfpara/reportServiceParamListZxjs.json")
	if err != nil {
		setting.ZAPS.Debugf("上报服务[Zxjs]配置json文件读取失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &configParam)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[Emqx]配置json文件格式化失败")
		return false
	}
	setting.ZAPS.Debug("上报服务[Zxjs]配置json文件读取成功")

	//初始化
	for _, v := range configParam.ServiceList {
		ReportServiceParamListZxjs.ServiceList = append(ReportServiceParamListZxjs.ServiceList, NewReportServiceParamZxjs(v.GWParam, v.NodeList))
	}

	return true
}

func ReportServiceZxjsWriteParamToJson() {
	utils.DirIsExist("./selfpara")
	sJson, _ := json.Marshal(ReportServiceParamListZxjs)
	err := utils.FileWrite("./selfpara/reportServiceParamListZxjs.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[Zxjs]配置json文件写入失败")
		return
	}
	setting.ZAPS.Debugf("上报服务[Zxjs]配置json文件写入成功")
}

func (s *ReportServiceParamListZxjsTemplate) AddReportService(param ReportServiceGWParamZxjsTemplate) error {

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

	nodeList := make([]ReportServiceNodeParamZxjsTemplate, 0)
	ReportServiceParam := NewReportServiceParamZxjs(param, nodeList)
	s.ServiceList = append(s.ServiceList, ReportServiceParam)

	ReportServiceZxjsWriteParamToJson()
	return nil
}

func (s *ReportServiceParamListZxjsTemplate) ModifyReportService(param ReportServiceGWParamZxjsTemplate) error {

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
	ReportServiceZxjsWriteParamToJson()

	//启动新上报服务
	s.ServiceList[index] = NewReportServiceParamZxjs(param, s.ServiceList[index].NodeList)

	return nil
}
func (s *ReportServiceParamListZxjsTemplate) DeleteReportService(serviceName string) {

	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == serviceName {
			s.ServiceList = append(s.ServiceList[:k], s.ServiceList[k+1:]...)
			ReportServiceZxjsWriteParamToJson()
			//协程退出
			v.CancelFunc()
			return
		}
	}
}

func (r *ReportServiceParamZxjsTemplate) AddReportNode(param ReportServiceNodeParamZxjsTemplate) error {

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
	ReportServiceZxjsWriteParamToJson()

	return nil
}
func (r *ReportServiceParamZxjsTemplate) ModifyReportNode(param ReportServiceNodeParamZxjsTemplate) error {

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
			ReportServiceZxjsWriteParamToJson()

			return nil
		}
	}

	//节点不存在则新建
	return errors.New("设备名称不存在")
}

func (r *ReportServiceParamZxjsTemplate) DeleteReportNode(name string) int {

	index := -1
	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == name {
			index = k

			//解除订阅主题
			r.MQTTZxjsUnsubNodeTopic(v.Param.DeviceSn)

			r.NodeList = append(r.NodeList[:k], r.NodeList[k+1:]...)
			ReportServiceZxjsWriteParamToJson()
			return index
		}
	}
	return index
}

func (r *ReportServiceParamZxjsTemplate) ExportParamToCsv() (bool, string) {

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
		{"上报服务名称", "采集接口名称", "设备名称", "设备通信地址", "上报模型", "上报服务协议", "产品序列号", "设备序列号", "设备密码"},
		{"ServiceName", "CollInterfaceName", "Name", "Addr", "UploadModel", "Protocol", "ProductSn", "DeviceSn", "DeviceSecret"},
	}

	for _, n := range r.NodeList {
		param := make([]string, 0)
		param = append(param, n.ServiceName)
		param = append(param, n.CollInterfaceName)
		param = append(param, n.Name)
		param = append(param, n.Addr)
		param = append(param, n.UploadModel)
		param = append(param, n.Protocol)
		param = append(param, n.Param.ProductSn)
		param = append(param, n.Param.DeviceSn)
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

func (r *ReportServiceParamZxjsTemplate) ExportParamToXlsx() (bool, string) {

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + r.GWParam.ServiceName + ".xlsx"

	cellData := [][]string{
		{"上报服务名称", "采集接口名称", "设备名称", "设备通信地址", "上报模型", "上报服务协议", "产品序列号", "设备序列号", "设备密码"},
		{"ServiceName", "CollInterfaceName", "Name", "Addr", "UploadModel", "Protocol", "ProductSn", "DeviceSn", "DeviceSecret"},
	}

	for _, n := range r.NodeList {
		param := make([]string, 0)
		param = append(param, n.ServiceName)
		param = append(param, n.CollInterfaceName)
		param = append(param, n.Name)
		param = append(param, n.Addr)
		param = append(param, n.UploadModel)
		param = append(param, n.Protocol)
		param = append(param, n.Param.ProductSn)
		param = append(param, n.Param.DeviceSn)
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
