/*////go:generate go run gen_build_info.go*/

package main

import (
	"fmt"
	"gateway/buildInfo"
	"gateway/device"
	"gateway/httpServer"
	"gateway/protocol/dlt645"
	"gateway/report"
	"gateway/setting"
	"gateway/utils"
	"gateway/virtual"
	"github.com/jasonlvhit/gocron"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var foo string

func main1111() {

	//outRulerInfo1, eer1 := dlt645.GetD07RulerInfo(0x0201ff00)
	//fmt.Printf("outRulerInfo<%v> eer<%v>  \r\n", outRulerInfo1, eer1)
	//
	//outRulerInfo2, eer2 := dlt645.GetD07RulerInfo(0x0202ff00)
	//fmt.Printf("outRulerInfo<%v> eer<%v>  \r\n", outRulerInfo2, eer2)
	//
	//outRulerInfo3, eer3 := dlt645.GetD07RulerInfo(0x0203ff00)
	//fmt.Printf("outRulerInfo<%v> eer<%v>  \r\n", outRulerInfo3, eer3)
	//
	//outRulerInfo4, eer4 := dlt645.GetD07RulerInfo(0x02010100)
	//fmt.Printf("outRulerInfo<%v> eer<%v>  \r\n", outRulerInfo4, eer4)
	//
	//outRulerInfo5, eer5 := dlt645.GetD07RulerInfo(0x02010200)
	//fmt.Printf("outRulerInfo<%v> eer<%v>  \r\n", outRulerInfo5, eer5)
	//
	//outRulerInfo6, eer6 := dlt645.GetD07RulerInfo(0x02010300)
	//fmt.Printf("outRulerInfo<%v> eer<%v>  \r\n", outRulerInfo6, eer6)

	fmt.Println("------------------------------0")
	//floatData := 0.0
	//timeData := ""
	//var frameData = []byte{0x44, 0x55, 0x06, 0x77, 0x88}
	//
	//d := dlt645.TransD07DataTemplate{
	//	TransDir: dlt645.ED07TransF2u,
	//	User:     floatData,
	//	Frame:    frameData,
	//}
	//
	//fmt.Println(dlt645.UnpackD07ByFormat("X.XXX", &d))
	//fmt.Println(dlt645.UnpackD07ByFormat("XX.XXXX", &d))
	//fmt.Println(dlt645.UnpackD07ByFormat("XXX.XXX", &d))
	//fmt.Println(dlt645.UnpackD07ByFormat("XXX.X", &d))
	//fmt.Println(dlt645.UnpackD07ByFormat("XXXXXX.XX", &d))
	//
	//fmt.Println("------------------------------000000")
	//
	//d.User = timeData
	//fmt.Println(dlt645.UnpackD07ByFormat("YYMMDDhhmm", &d))
	//
	//fmt.Println(dlt645.UnpackD07ByFormat("YYMMDDhhmm", &dlt645.TransD07DataTemplate{TransDir: dlt645.ED07TransU2f, User: "2305090940", Frame: []byte{}}))
	//
	//fmt.Println(dlt645.UnpackD07ByFormat("X.XXX", &dlt645.TransD07DataTemplate{TransDir: dlt645.ED07TransU2f, User: 2.211, Frame: []byte{}}))
	//fmt.Println(dlt645.UnpackD07ByFormat("XX.XXXX", &dlt645.TransD07DataTemplate{TransDir: dlt645.ED07TransU2f, User: 33.2211, Frame: []byte{}}))
	//fmt.Println(dlt645.UnpackD07ByFormat("XXX.XXX", &dlt645.TransD07DataTemplate{TransDir: dlt645.ED07TransU2f, User: 332.211, Frame: []byte{}}))
	//fmt.Println(dlt645.UnpackD07ByFormat("XXX.X", &dlt645.TransD07DataTemplate{TransDir: dlt645.ED07TransU2f, User: 221.1, Frame: []byte{}}))
	//fmt.Println(dlt645.UnpackD07ByFormat("XXXXXX.XX", &dlt645.TransD07DataTemplate{TransDir: dlt645.ED07TransU2f, User: 443322.11, Frame: []byte{}}))
	//
	//fmt.Println("------------------------------1")

	//buff := dlt645.D07PackFrame{
	//	RulerId:  "000a0b02",
	//	CtrlCode: 0x11,
	//	DataLen:  0x08,
	//	Address:  "0123456789ab",
	//	Data:     []byte{0x00, 0x01, 0x02, 0x03},
	//}
	//buff.PackD07FrameByData()
	//
	//fmt.Println("------------------------------2")
	//bcddata := []byte{0x20, 0x01, 0x12, 0x12, 0x15, 0x20}
	//str, ok := dlt645.D07BCD2Str(bcddata, 6)
	//fmt.Printf("D07BCD2Str-->bcd[%x] str[%s] ok[%v]\r\n", bcddata, str, ok)
	//
	//str1 := "201512120120"
	//bcd, ok := dlt645.D07Str2BCD(str1, 6)
	//fmt.Printf("D07Str2BCD-->str[%s] bcd[%x] ok[%v]\r\n", str1, bcd, ok)
	//
	//fmt.Println("------------------------------3")
	//dd := dlt645.TransD07DataTemplate{
	//	TransDir: dlt645.ED07TransU2f,
	//	User:     float32(2.16),
	//	Frame:    []byte{0x05, 0x06},
	//}
	//
	//fmt.Println(dd)
	//rr, err := dd.TransD07DataX_XXX()
	//fmt.Printf("---->E_D07_TRANS_U2F  user[%v] frame[%v][%x] err[%v]\r\n", rr.User, rr.Frame, rr.Frame, err)
	//
	//dd.TransDir = dlt645.ED07TransF2u
	//dd.User = float32(0.2)
	//dd.Frame = []byte{0x49, 0x35}
	//
	//fmt.Println(dd)
	//ff, err1 := dd.TransD07DataX_XXX()
	//fmt.Printf("---->E_D07_TRANS_F2U  user[%v] frame[%v][%x] err[%v]\n", ff.User, ff.Frame, ff.Frame, err1)

	//0xFE, 0xFE, 0xFE, 0xFE, 0x68, 0x30, 0x94, 0x39, 0x00, 0x00, 0x00, 0x68, 0x91, 0x09, 0x33, 0x32, 0x34, 0x33, 0x74, 0xA6, 0x36, 0x33, 0x34, 0xEA, 0x16
	d, ok1 := dlt645.UnpackD07Frame([]byte{0xFE, 0xFE, 0xFE, 0xFE, 0x68, 0x30, 0x94, 0x39, 0x00, 0x00, 0x00, 0x68, 0x91, 0x0D, 0x33, 0x33, 0x34, 0x34, 0x86, 0x66, 0x37, 0x7B, 0x4A, 0x3C, 0x38, 0x56, 0x34, 0x1F, 0x16})
	fmt.Println(d, ok1)

	i, ok2 := dlt645.GetD07RulerInfo(d.RulerID)
	fmt.Println(i, ok2)

	v1, _ := dlt645.UnpackD07ByFormat(i[0].Format, &dlt645.TransD07DataTemplate{TransDir: dlt645.ED07TransF2u, User: 0.0, Frame: d.Data[0:]})
	v2, _ := dlt645.UnpackD07ByFormat(i[1].Format, &dlt645.TransD07DataTemplate{TransDir: dlt645.ED07TransF2u, User: "", Frame: d.Data[i[1].RulerAddOffset-1:]})

	fmt.Printf("%s=%f\r\n", i[0].Name, v1.User)
	fmt.Printf("%s=%s\r\n", i[1].Name, v2.User)

}

func main() {
	utils.ErrorLogInit()
	defer func() {
		r := recover()
		if r != nil {
			utils.ErrLog.Printf("程序发生错误原因 %v", r)
		}
	}()

	if foo == "" {
		fmt.Println("foo is empty.")
	} else {
		fmt.Printf("foo=%s\n", foo)
	}

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
	_ = scheduler.Every(60).Second().Do(setting.CollectSystemParam)
	// 定时1小时,定时获取NTP服务器的时间，并校时
	_ = scheduler.Every(1).Hour().Do(setting.NTPGetTime)
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
	//**************应用程序监听退出****************/

}
