//go:generate go run gen_build_info.go

package main

import (
	"fmt"
	"gateway/buildInfo"
	"gateway/device"
	"gateway/httpServer"
	"gateway/models"
	"gateway/report"
	"gateway/rule"
	"gateway/service/job"
	"gateway/setting"
	"gateway/utils"
	"gateway/virtual"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/jasonlvhit/gocron"
)

func main() {
	// 连接SQLite数据库
	models.InitDB()
	models.InitTaosDB()
	utils.ErrorLogInit()
	defer func() {
		r := recover()
		if r != nil {
			utils.ErrLog.Printf("程序发生错误原因 %v", r)
		}
	}()

	fmt.Println(buildInfo.BuildTime)

	//time.Sleep(time.Second * 30)

	/**************获取配置文件***********************/
	setting.GetConf()
	/**************日志初始化***********************/
	setting.InitLogger()
	/**************起始***********************/
	setting.SystemInit()
	setting.ZAPS.Infof("%s %s", setting.SystemState.Name, setting.SystemState.SoftVer)
	/**************网口初始化***********************/
	setting.ReadNetworkParamFromJson()
	/************读取授权文件*************************/
	setting.ReadAuthFile()
	/**************移动模块初始化***********************/
	setting.MobileModuleInit()
	/**************变量模板初始化****************/
	device.TSLModelsInit()
	device.CollectInterfaceInit()
	/**************NTP校时初始化****************/
	setting.NTPInit()
	go func() {
		time.Sleep(time.Minute)
		err := setting.NTPGetTime()
		if err == nil {
			setting.GetTimeStart()
		}
	}()
	/**************创建定时获取网络状态的任务***********************/

	scheduler := gocron.NewScheduler()
	// 定时60秒,定时获取系统信息
	// 定时1小时,定时获取NTP服务器的时间，并校时
	_ = scheduler.Every(1).Hour().Do(setting.NTPGetTime)
	// 启动充放电量定时任务
	setting.ZAPS.Debug("注册充放电量定时任务到 GoCron")
	//_ = scheduler.Every(1).Hour().Do(job.StatisticsDayChargingAndDischarging)

	// 每天0：0重启系统
	//_ = scheduler.Every(1).Day().At("0:0").Do(setting.SystemReboot)
	scheduler.Start()
	/**************虚拟设备初始化****************/
	virtual.VirtualDeviceInit()
	/**************上报服务初始化****************/
	report.ReportServiceInit()
	/**************LED初始化****************/
	setting.LedInit()
	setting.LedFlash(setting.LEDNet)
	/**************pprof****************/
	go func() {
		_ = http.ListenAndServe("0.0.0.0:9090", nil)
	}()
	/**************Http服务初始化****************/
	httpServer.RouterWeb(setting.HttpPort)
	setting.ZAPS.Infof("gateway 初始化成功!")
	/**************rule服务初始化****************/
	rule.Init()
	//**************应用程序监听退出****************/
}
