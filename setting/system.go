package setting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gateway/buildInfo"
	"gateway/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemStateTemplate struct {
	MemTotal         string `json:"memTotal"`
	MemUse           string `json:"memUse"`
	CPUUse           string `json:"cpuUse"`
	DiskTotal        string `json:"diskTotal"`
	DiskUse          string `json:"diskUse"`
	Name             string `json:"name"`
	SN               string `json:"SN"`
	HardVer          string `json:"hardVer"`
	SoftVer          string `json:"softVer"`
	SystemRTC        string `json:"systemRTC"`
	RunTime          string `json:"runTime"`          //累计时间
	DeviceOnline     string `json:"deviceOnline"`     //设备在线率
	DevicePacketLoss string `json:"devicePacketLoss"` //设备丢包率
	LockStatus       int    `json:"lockStatus"`       //锁定状态
	GOOS             string `json:"GOOS"`
	GOARCH           string `json:"GOARCH"`
	AuthStatus       string `json:"authStatus"`
	BuildTime        string `json:"buildTime"`
}

type DataPointTemplate struct {
	Value string `json:"value"`
	Time  string `json:"time"`
}

type DataStreamTemplate struct {
	DataPoint    []DataPointTemplate `json:"dataPoint"`
	DataPointCnt int                 `json:"dataPointCnt"`
	Legend       string              `json:"legend"` //别名
}

type LockStatusTemplate struct {
	Status int
}

const (
	LockCmdDisable = iota
	LockCmdEnable
)

var SystemState = SystemStateTemplate{
	MemTotal:         "0",
	MemUse:           "0",
	CPUUse:           "0",
	DiskTotal:        "0",
	DiskUse:          "0",
	Name:             "",
	SN:               "",
	HardVer:          "V0.0.1",
	SoftVer:          "V1.0.1",
	SystemRTC:        "2020-05-26 12:00:00",
	RunTime:          "0",
	DeviceOnline:     "0",
	DevicePacketLoss: "0",
	BuildTime:        buildInfo.BuildTime,
	LockStatus:       LockCmdDisable,
	GOOS:             runtime.GOOS,
	GOARCH:           runtime.GOARCH,
	AuthStatus:       "未授权",
}

var timeStart time.Time
var (
	CPUDataStream              *DataStreamTemplate
	MemoryDataStream           *DataStreamTemplate
	DiskDataStream             *DataStreamTemplate
	DeviceOnlineDataStream     *DataStreamTemplate
	DevicePacketLossDataStream *DataStreamTemplate
)
var LockStatus LockStatusTemplate

//var CSTZone =

func SystemInit() {
	CPUDataStream = NewDataStreamTemplate("CPU使用率")
	MemoryDataStream = NewDataStreamTemplate("内存使用率")
	DiskDataStream = NewDataStreamTemplate("硬盘使用率")
	DeviceOnlineDataStream = NewDataStreamTemplate("设备在线率")
	DevicePacketLossDataStream = NewDataStreamTemplate("通信丢包率")

	_ = ReadProductParamFromJson()

	SystemState.HardVer = GetHardVer()

	//_ = GetSystemLock()

	time.Local = time.FixedZone("CST", 8*3600)

	GetTimeStart()
}

func SystemReboot() {
	cmd := exec.Command("reboot")
	var out bytes.Buffer
	cmd.Stdout = &out
	str := out.String()
	fmt.Println(str)

	err := cmd.Run()
	if err != nil {
		ZAPS.Errorf("执行重启命令失败 %v", err)
	}
}

func SystemLock(cmd int) {
	SystemState.LockStatus = cmd

	LockStatus.Status = cmd
	sJson, _ := json.Marshal(LockStatus)

	utils.DirIsExist("./selfpara")
	err := utils.FileWrite("./selfpara/lockStatus.json", sJson)
	if err != nil {
		ZAPS.Errorf("写入锁定命令json文件 %s %v", "失败", err)
		return
	}
	ZAPS.Infof("写入锁定命令json文件 %s", "成功")
}

func GetSystemLock() error {

	data, err := utils.FileRead("./selfpara/lockStatus.json")
	if err != nil {
		ZAPS.Debugf("打开锁定命令json文件失败 %v", err)
		return err
	}
	err = json.Unmarshal(data, &LockStatus)
	if err != nil {
		ZAPS.Errorf("锁定命令json文件格式化失败 %v", err)
		return err
	}
	ZAPS.Debugf("打开锁定命令json文件成功")

	return nil
}

func SystemSetRTC(rtc string) {

	var out bytes.Buffer
	ZAPS.Debugf("系统校时 %v", rtc)

	//cmdStr := "\"" + rtc + "\""
	cmd := exec.Command("date", "-s", rtc)
	ZAPS.Debugf("系统校时参数 %v", cmd.Args)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		ZAPS.Errorf("系统校时执行date -s 失败 %v", rtc)
		return
	}
	err = cmd.Wait()
	if err != nil {
		ZAPS.Errorf("系统校时执行date -s 失败 %v", rtc)
		return
	}
	ZAPS.Debugf("系统校时执行date结果 %v", out.String())

	if strings.Compare(runtime.GOARCH, "arm") == 0 {
		//将时间写入硬件RTC中
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd = exec.CommandContext(ctx, "hwclock", "-w")
		cmd.Stdout = &out
		err = cmd.Start()
		if err != nil {
			ZAPS.Errorf("系统校时指执行hwclock -w 失败 %v", rtc)
			return
		}
		err = cmd.Wait()
		if err != nil {
			ZAPS.Errorf("系统校时执行hwclock -w 失败 %v", rtc)
			return
		}
		ZAPS.Debugf("系统校时执行hwclock -w结果 %v", out.String())
	}
}

func GetCPUState(interval time.Duration) {

	v, _ := cpu.Percent(interval, false)
	SystemState.CPUUse = fmt.Sprintf("%3.1f", v[0])
}

func GetMemState() {

	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	//log.Printf("Mem Total: %v, Free:%v, UsedPercent:%f%%\n",
	//					v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)

	SystemState.MemTotal = fmt.Sprintf("%d", v.Total/1024/1024)
	SystemState.MemUse = fmt.Sprintf("%3.1f", v.UsedPercent)
}

func GetDiskState() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	v, _ := disk.Usage(exeCurDir)

	// almost every return value is a struct
	//log.Printf("Disk Total: %v, Free:%v, UsedPercent:%f%%\n",
	//				v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)

	SystemState.DiskTotal = fmt.Sprintf("%d", v.Total/1024/1024)
	SystemState.DiskUse = fmt.Sprintf("%3.1f", v.UsedPercent)
}

func GetTimeStart() {

	//BeiJingZone := time.FixedZone("CST", 8*3600)
	timeStart = time.Now()
	ZAPS.Infof("系统当前时间 %v", time.Now().Format("2006-01-02 15:04:05"))
}

func GetRunTime() {

	elapsed := time.Since(timeStart)
	sec := int64(elapsed.Seconds())
	day := sec / 86400
	hour := sec % 86400 / 3600
	min := sec % 3600 / 60
	sec = sec % 60

	strRunTime := fmt.Sprintf("%d天%d时%d分%d秒", day, hour, min, sec)

	//BeiJingZone := time.FixedZone("CST", 8*3600)
	SystemState.SystemRTC = time.Now().Format("2006-01-02 15:04:05")
	SystemState.RunTime = strRunTime
}

func NewDataStreamTemplate(legend string) *DataStreamTemplate {

	return &DataStreamTemplate{
		DataPoint:    make([]DataPointTemplate, 0),
		DataPointCnt: 0,
		Legend:       legend,
	}
}

func (d *DataStreamTemplate) AddDataPoint(data DataPointTemplate) {

	if d.DataPointCnt < 100 {
		d.DataPoint = append(d.DataPoint, data)
		d.DataPointCnt++
	} else {
		d.DataPoint = d.DataPoint[1:]
		d.DataPoint = append(d.DataPoint, data)
	}
}

func CollectSystemParam() {

	GetCPUState(100 * time.Millisecond)
	GetMemState()
	GetRunTime()

	point := DataPointTemplate{}

	point.Value = SystemState.CPUUse
	point.Time = SystemState.SystemRTC
	CPUDataStream.AddDataPoint(point)

	point.Value = SystemState.MemUse
	point.Time = SystemState.SystemRTC
	MemoryDataStream.AddDataPoint(point)

	point.Value = SystemState.DiskUse
	point.Time = SystemState.SystemRTC
	DiskDataStream.AddDataPoint(point)

	point.Value = SystemState.DeviceOnline
	point.Time = SystemState.SystemRTC
	DeviceOnlineDataStream.AddDataPoint(point)

	point.Value = SystemState.DevicePacketLoss
	point.Time = SystemState.SystemRTC
	DevicePacketLossDataStream.AddDataPoint(point)
}

func GetSystemOS() {
	SystemState.GOOS = runtime.GOOS
	SystemState.GOARCH = runtime.GOARCH
}
