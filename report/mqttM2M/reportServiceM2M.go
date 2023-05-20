package mqttM2M

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/setting"
	"gateway/utils"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/robfig/cron"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ParamM2MTemplate struct {
	UserName string
	Password string
	ClientID string //MQTT客户端ID，用DeviceID+一个随机数(1000内)
	AppKey   string
	DeviceID string
}

// 上报网关参数结构体
type ReportServiceGWParamM2MTemplate struct {
	ServiceName       string
	IP                string
	Port              string
	ReportStatus      string
	ReportTime        int
	ReportErrCnt      int
	Protocol          string
	Param             ParamM2MTemplate
	MQTTClient        MQTT.Client         `json:"-"`
	MQTTClientOptions *MQTT.ClientOptions `json:"-"`
	OfflineChan       chan bool           `json:"-"`
	HeartBeatMark     bool                `json:"-"` //定时上报时间到了，是否需要推送心跳，如果定时过程中有其他命令已经发过数据了，就不再发心跳
	HeartBeatCnt      uint32              `json:"-"`
}

// 上报服务参数，网关参数，节点参数
type ReportServiceParamM2MTemplate struct {
	GWParam                        ReportServiceGWParamM2MTemplate
	ReceiveFrameChan               chan MQTTM2MReceiveFrameTemplate   `json:"-"`
	ReceiveLogInAckFrameChan       chan MQTTM2MLogInAckTemplate       `json:"-"` //平台下发的登陆确认报文
	ReportPropertyRequestFrameChan chan MQTTM2MReportPropertyTemplate `json:"-"`
	CancelFunc                     context.CancelFunc                 `json:"-"`
}

type ReportServiceParamListM2MTemplate struct {
	ServiceList []*ReportServiceParamM2MTemplate
}

// 实例化上报服务
var ReportServiceParamListM2M ReportServiceParamListM2MTemplate

func ReportServiceM2MInit() {

	setting.ZAPS.Infof("初始化M2M...")

	p := ParamM2MTemplate{
		UserName: "feisjy",
		Password: "feisjy2016",
		ClientID: "gwai123456",
		AppKey:   "Feisjy20190507",
		DeviceID: "gwai123456",
	}

	gp := ReportServiceGWParamM2MTemplate{
		ServiceName:       "m2m",
		IP:                "m2m.feisjy.com",
		Port:              "1883",
		ReportStatus:      "",
		ReportTime:        10,
		ReportErrCnt:      0,
		Protocol:          "",
		Param:             p,
		MQTTClient:        nil,
		MQTTClientOptions: nil,
		OfflineChan:       nil,
		HeartBeatMark:     false,
		HeartBeatCnt:      0,
	}

	ReportServiceParamListM2M.ServiceList = append(ReportServiceParamListM2M.ServiceList, NewReportServiceParamFeisjy(gp))
}

func NewReportServiceParamFeisjy(gw ReportServiceGWParamM2MTemplate) *ReportServiceParamM2MTemplate {
	feisjyParam := &ReportServiceParamM2MTemplate{
		GWParam:                        gw,
		ReceiveFrameChan:               make(chan MQTTM2MReceiveFrameTemplate, 100),
		ReceiveLogInAckFrameChan:       make(chan MQTTM2MLogInAckTemplate, 2),
		ReportPropertyRequestFrameChan: make(chan MQTTM2MReportPropertyTemplate, 50),
	}

	ctx, cancel := context.WithCancel(context.Background())
	feisjyParam.CancelFunc = cancel

	go ReportServiceM2MPoll(ctx, feisjyParam)

	return feisjyParam
}

func ReportServiceM2MPoll(ctx context.Context, r *ReportServiceParamM2MTemplate) {

	reportState := 0

	// 定义一个cron运行器
	cronProcess := cron.New()

	reportTime := fmt.Sprintf("@every %dm%ds", r.GWParam.ReportTime/60, r.GWParam.ReportTime%60)
	setting.ZAPS.Infof("上报服务[%s]定时上报周期为%v", r.GWParam.ServiceName, reportTime)

	_ = cronProcess.AddFunc(reportTime, r.ReportTimeFun)

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
						time.Sleep(5 * time.Second)
					}
					break

				case 1: //  网关登录流程
					if r.M2MLogInMachine(r.GWParam.Param.DeviceID) == true {
						r.GWParam.ReportStatus = "onLine"
						cronProcess.Start()
						reportState = 2
					} else {
						time.Sleep(10 * time.Second)
					}
					break
				case 2:
					reportGWProperty := MQTTM2MReportPropertyTemplate{
						DeviceType: "gw",
					}
					r.ReportPropertyRequestFrameChan <- reportGWProperty
					time.Sleep(10 * time.Second)
					break
				}

				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (r *ReportServiceParamM2MTemplate) ReportTimeFun() {

	if r.GWParam.ReportStatus == "offLine" {
		return
	}
	//发送网关心跳到平台
	if r.GWParam.HeartBeatMark == false { //如果有其它数据推送到平台了，可以不发心路
		r.GWParam.HeartBeatCnt++
		r.FeisjyHeartBeat(r.GWParam.Param.DeviceID, r.GWParam.HeartBeatCnt)
	}
	r.GWParam.HeartBeatMark = false

}

func (r *ReportServiceParamM2MTemplate) ProcessUpLinkFrame(ctx context.Context) {

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
				}
			}
		}
	}
}

func (r *ReportServiceParamM2MTemplate) ProcessDownLinkFrame(ctx context.Context) {

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
					ackFrame := MQTTM2MLogInAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.ZAPS.Errorf("LogIn Ack json unmarshal err")
						continue
					}
					r.ReceiveLogInAckFrameChan <- ackFrame
					break
				case "allcall":
					if ID == r.GWParam.Param.DeviceID {
						reportGWProperty := MQTTM2MReportPropertyTemplate{
							DeviceType: "gw",
						}
						r.ReportPropertyRequestFrameChan <- reportGWProperty
					}
				case "resultYx":
				case "resultYc":
				case "setting":
					var settingPayload MQTTM2MReportSettingTemplate
					err := json.Unmarshal(frame.Payload, &settingPayload)
					if err != nil {
						continue
					}
					for _, v := range settingPayload.SettingList {
						switch v.ID {
						case 200: //git方式升级设备
							setting.ZAPS.Infof("git方式升级设备[%v]", v.Value)
						case 201: //http方式升级设备
							httpStr := "http://m2m.feisjy.com/" + ID + "/update/" + "system.rar"
							setting.ZAPS.Infof("http方式升级设备[%v]", httpStr)
							exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
							filePath := exeCurDir + "/m2mFile/"

							if ok := HttpGetFile(filePath, "system.rar", httpStr, 5); ok != nil {
								setting.ZAPS.Errorf("m2m升级系统失败[%v]", ok)
							} else {
								//fileAbsoluteDir := exeCurDir + "/"
							}
						}
					}
				}
			}
		}
	}
}

func HttpGetFile(fileSavePath string, fileName string, url string, maxRetries int) error {

	// 实例化一个http client对象
	client := &http.Client{}
	// 实例化一个http Get Request对象
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("http.NewRequest Fail! %v", err))
	}

	req.Header.Add("Cookie", "SHAREJSESSIONID=f0bae696-f1be-4df8-a5e1-10ba646b7ff0")

	// 设置 Range 头字段，指定要下载的文件块的起始和结束位置
	blockSize := 1024 * 1024 // 每次下载1MB的数据
	var start, end int64
	var fileSize int64
	for retries := 0; retries < maxRetries; retries++ {
		fmt.Printf("[%d]", retries+1)

		// 获取文件大小
		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("获取文件大小失败，正在进行第 %d 次重试，错误信息：%v\n", retries+1, err)
			continue
		}
		defer res.Body.Close()

		fileSize = res.ContentLength
		if fileSize == -1 {
			return errors.New(fmt.Sprintf("无法获取文件大小!"))
		}

		// 判断是否路径是否存在，不存在则递归创建
		utils.DirIsExist(fileSavePath)

		filePath := fmt.Sprintf("%s%s", fileSavePath, fileName)

		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return errors.New(fmt.Sprintf("文件打开失败[%s]! %v", filePath, err))
		}
		defer file.Close()

		// 创建一个带缓冲区的写入器
		write := bufio.NewWriter(file)

		// 分块下载文件
		for start < fileSize {

			fmt.Printf("□□□□□□□□")

			end = start + int64(blockSize) - 1
			if end >= fileSize {
				end = fileSize - 1
			}

			// 设置 Range 头字段，指定要下载的文件块的起始和结束位置
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

			// 将http Get 请求对象发送服务器上 若发送成功则会得到响应
			res, err := client.Do(req)
			if err != nil {
				fmt.Printf("下载失败，正在进行第 %d 次重试，错误信息：%v\n", retries+1, err)
				break
			}
			defer res.Body.Close()

			// 将响应体直接写入文件中
			_, err = io.Copy(write, res.Body)
			if err != nil {
				fmt.Printf("写入文件失败，正在进行第 %d 次重试，错误信息：%v\n", retries+1, err)
				break
			}

			// 更新起始位置
			start = end + 1
		}
		fmt.Printf("\n")
		// 确保文件被正确关闭
		err = write.Flush()
		if err != nil {
			return errors.New(fmt.Sprintf("文件关闭失败[%s]! %v", filePath, err))
		}

		return nil
	}

	return errors.New(fmt.Sprintf("下载失败，已达到最大重试次数：%d", maxRetries))
}
