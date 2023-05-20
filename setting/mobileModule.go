package setting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/utils"
	"strconv"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"

	"github.com/tarm/serial"
)

const (
	OperatorTypeCTCC     int = iota //0，电信
	OperatorTypeCMCCCUCC            //1，移动或联通
)

type MobileModuleCommParamTemplate struct {
	Name       string `json:"name"`
	BaudRate   string `json:"baudRate"`
	DataBits   string `json:"dataBits"`   //数据位: 5, 6, 7 or 8 (default 8)
	StopBits   string `json:"stopBits"`   //停止位: 1 or 2 (default 1)
	Parity     string `json:"parity"`     //校验: N - None, E - Even, O - Odd (default E),(The use of no parity requires 2 stop bits.)
	Timeout    string `json:"timeout"`    //通信超时（毫秒）
	Interval   string `json:"interval"`   //通信间隔（毫秒）
	PollPeriod string `json:"pollPeriod"` //轮询周期（秒）
}

type MobileModuleRunParamTemplate struct {
	ICCID     string  `json:"iccid"` //SIM卡号
	IMEI      string  `json:"imei"`  //模块IMEI编号
	CSQ       string  `json:"csq"`   //信号强度
	Flow      string  `json:"flow"`  //手机流量
	CurFlow   float32 `json:"-"`
	Lac       string  `json:"lac"`                //基站定位数据
	Ci        string  `json:"ci"`                 //基站定位数据
	SIMInsert bool    `json:"simInsert"`          //手机卡是否插入
	Network   bool    `json:"netRegister"`        //网络注册状态
	Operator  int     `json:"operator,omitempty"` //运营商
}

type MobileModuleConfigParamTemplate struct {
	FlowAlarmEnable bool   `json:"flowAlarm"`      //手机流量告警使能
	FlowAlarm       string `json:"flowAlarmValue"` //手机流量告警阈值
}

type MobileModuleKeepAliveParamTemplate struct {
	IPMaster   string `json:"ipMaster"`   //保活检测IP1
	IPSlave    string `json:"ipSlave"`    //保活检测IP2
	PollPeriod int    `json:"pollPeriod"` //轮询周期（秒）
	OfflineCnt int    `json:"offlineCnt"` //离线判断次数
	CurrentCnt int    `json:"-"`
	ResetCnt   int    `json:"-"`
}

type MobileModuleTemplate struct {
	Name           string                             `json:"name"`                //模块名称
	Model          string                             `json:"model"`               //模块型号
	CommParam      MobileModuleCommParamTemplate      `json:"commParam"`           //通信参数
	RunParam       MobileModuleRunParamTemplate       `json:"runParam"`            //模块参数
	ConfigParam    MobileModuleConfigParamTemplate    `json:"configParam"`         //配置参数
	KeepAliveParam MobileModuleKeepAliveParamTemplate `json:"keepAliveParam"`      //保活检测参数
	GPIOParam      MobileModuleGPIOParamTemplate      `json:"GPIOParam,omitempty"` //GPIO接口
	ReadBuffer     chan []byte                        `json:"-"`                   //通信接收缓存
	Port           *serial.Port                       `json:"-"`                   //通信句柄
	Status         bool                               `json:"status"`              //开关状态
	CancelFunc     context.CancelFunc                 `json:"-"`
}

type MobileModuleGPIOValueParamTemplate struct {
	GPIO       string `json:"gpio"`
	SetValue   string `json:"setValue"`
	ResetValue string `json:"resetValue"`
}

type MobileModuleGPIOParamTemplate struct {
	PowerCtr      MobileModuleGPIOValueParamTemplate `json:"powerCtr"`      //电源管脚
	ResetCtr      MobileModuleGPIOValueParamTemplate `json:"resetCtr"`      //复位管脚
	RSSILedHigh   MobileModuleGPIOValueParamTemplate `json:"RSSILedHigh"`   //信号强度高指示灯
	RSSILedMiddle MobileModuleGPIOValueParamTemplate `json:"RSSILedMiddle"` //信号强度中指示灯
	RSSILedLow    MobileModuleGPIOValueParamTemplate `json:"RSSILedLow"`    //信号强度低指示灯
}

type MobileModuleParamTemplate struct {
	Name           string                             `json:"name"`                //模块名称
	Model          string                             `json:"model"`               //模块型号
	ConfigParam    MobileModuleConfigParamTemplate    `json:"configParam"`         //模块配置参数
	CommParam      MobileModuleCommParamTemplate      `json:"commParam"`           //模块通信参数
	KeepAliveParam MobileModuleKeepAliveParamTemplate `json:"keepAliveParam"`      //保活检测参数
	GPIOParam      MobileModuleGPIOParamTemplate      `json:"GPIOParam,omitempty"` //GPIO接口
}

var MobileModule *MobileModuleTemplate
var MobileModuleParam MobileModuleParamTemplate

func MobileModuleInit() {
	//获取移动模块配置参数
	if ReadMobileModuleParamFromJson() == false {
		MobileModule = &MobileModuleTemplate{}
		return
	}
	//实例化
	MobileModule = NewMobileModule(MobileModuleParam)

}

func NewMobileModule(param MobileModuleParamTemplate) *MobileModuleTemplate {
	mobile := &MobileModuleTemplate{
		Name:           param.Name,
		Model:          param.Model,
		CommParam:      param.CommParam,
		ConfigParam:    param.ConfigParam,
		KeepAliveParam: param.KeepAliveParam,
		GPIOParam:      param.GPIOParam,
		ReadBuffer:     make(chan []byte, 1024),
		Status:         false,
	}

	ctx, cancel := context.WithCancel(context.Background())
	mobile.CancelFunc = cancel

	mobile.MobileModuleOpen()

	go mobile.MobileModuleReadATCommand(ctx)

	go mobile.MobileModulePowerOn()

	go mobile.MobileModuleScheduler(ctx)

	go mobile.MobileModuleKeepAliveScheduler(ctx)

	return mobile
}

func (m *MobileModuleTemplate) MobileModuleOpen() bool {
	serialParam := m.CommParam
	serialBaud, _ := strconv.Atoi(serialParam.BaudRate)

	var serialParity serial.Parity
	switch serialParam.Parity {
	case "N":
		serialParity = serial.ParityNone
	case "O":
		serialParity = serial.ParityOdd
	case "E":
		serialParity = serial.ParityEven
	}

	var serialStop serial.StopBits
	switch serialParam.StopBits {
	case "1":
		serialStop = serial.Stop1
	case "1.5":
		serialStop = serial.Stop1Half
	case "2":
		serialStop = serial.Stop2
	}

	serialConfig := &serial.Config{
		Name:        serialParam.Name,
		Baud:        serialBaud,
		Parity:      serialParity,
		StopBits:    serialStop,
		ReadTimeout: time.Millisecond * 1,
	}

	var err error
	m.Port, err = serial.OpenPort(serialConfig)
	if err != nil {
		ZAPS.Errorf("移动模块串口接口[%s]打开失败 %v", m.CommParam.Name, err)
		return false
	}
	ZAPS.Debugf("移动模块串口接口[%s]打开成功", m.CommParam.Name)
	m.Status = true

	return true
}

func (m *MobileModuleTemplate) MobileModuleClose() bool {
	if m.Port == nil {
		ZAPS.Debugf("移动模块串口接口[%s] 文件句柄不存在", m.CommParam.Name)
		return false
	}
	err := m.Port.Close()
	if err != nil {
		ZAPS.Errorf("移动模块串口接口[%s]关闭失败 %v", m.CommParam.Name, err)
		return false
	}
	ZAPS.Infof("移动模块串口接口[%s]关闭成功", m.CommParam.Name)
	m.Status = false
	return true
}

func (m *MobileModuleTemplate) MobileModuleWriteData(data string) int {

	if m.Port == nil {
		ZAPS.Debugf("移动模块串口%s文件句柄不存在", m.CommParam.Name)
		m.MobileModuleOpen()
		return 0
	}

	cnt, err := m.Port.Write([]byte(data))
	if err != nil {
		ZAPS.Debugf("移动模块写数据失败 %v", err)
	}

	return cnt
}

func (m *MobileModuleTemplate) MobileModuleReadData(data []byte) int {

	if m.Port == nil {
		return 0
	}

	cnt, _ := m.Port.Read(data)

	return cnt
}

func (m *MobileModuleTemplate) MobileModuleScheduler(ctx context.Context) {
	scheduler := gocron.NewScheduler()
	err := scheduler.Every(30).Second().Do(m.MobileModuleGetParam)
	if err != nil {
		ZAPS.Errorf("移动模块周期任务调用失败 %v", err)
		return
	}
	scheduler.Start()
	ZAPS.Info("移动模块周期任务开始")

	for {
		select {
		case <-ctx.Done():
			{
				scheduler.Clear()
				ZAPS.Info("移动模块周期任务结束")
				return
			}
		}
	}
}

func (m *MobileModuleTemplate) MobileModuleKeepAliveScheduler(ctx context.Context) {
	scheduler := gocron.NewScheduler()

	err := scheduler.Every(uint64(m.KeepAliveParam.PollPeriod)).Second().Do(m.MobileModuleKeepAlive)
	if err != nil {
		ZAPS.Errorf("移动模块保活检测周期任务调用失败 %v", err)
		return
	}
	scheduler.Start()
	ZAPS.Info("移动模块保活检测周期任务开始")

	for {
		select {
		case <-ctx.Done():
			{
				scheduler.Clear()
				ZAPS.Info("移动模块保活检测周期任务结束")
				return
			}
		}
	}
}

func (m *MobileModuleTemplate) MobileModuleReadATCommand(ctx context.Context) {
	//阻塞读
	rxBuf := make([]byte, 1024)
	rxBufCnt := 0
	ZAPS.Info("移动模块数据接收任务开始")
	for {
		select {
		case <-ctx.Done():
			{
				ZAPS.Info("移动模块数据接收任务结束")
				m.MobileModuleClose()
				return
			}
		default:
			{
				//阻塞读
				rxBufCnt = m.MobileModuleReadData(rxBuf)
				if rxBufCnt > 0 {
					//ZAPS.Debugf("%s:curRxBufCnt %v", collName, rxBufCnt)
					//ZAPS.Debugf("%s:CurRxBuf %X", collName, rxBuf[:rxBufCnt])

					//追加接收的数据到接收缓冲区
					m.ReadBuffer <- rxBuf[:rxBufCnt]
					//清除本次接收数据
					rxBufCnt = 0
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func (m *MobileModuleTemplate) MobileModuleWriteATCommand(txCommand, rxCommand string, timeOut int) (error, string) {

	rxBuf := make([]byte, 0)
	rxBufLen := 0

	//判断串口是否打开
	if m.Status == false {
		m.MobileModuleOpen()
	}
	//发送AT指令
	m.MobileModuleWriteData(txCommand)
	timerOut := time.NewTimer(time.Duration(timeOut) * time.Millisecond)
	for {
		select {
		case <-timerOut.C: //等待超时
			{
				rxStr := string(rxBuf[:rxBufLen])
				ZAPS.Debugf("移动模块接收超时,接收数据[%s]", rxStr)
				return errors.New("移动模块接收超时"), ""
			}
		case buf := <-m.ReadBuffer:
			{
				bufLen := len(buf)
				if bufLen > 0 {
					rxBuf = append(rxBuf, buf[:bufLen]...)
					rxBufLen += bufLen
				}
				rxStr := string(rxBuf[:rxBufLen])
				if strings.Contains(rxStr, rxCommand) == false {
					continue
				}
				//ZAPS.Debugf("移动模块接收成功，接收数据[%s]", rxStr)
				timerOut.Stop()
				return nil, rxStr
			}
		}
	}
}

func (m *MobileModuleTemplate) MobileModulePowerOn() {

	err, _ := m.MobileModuleWriteATCommand("AT\r\n", "OK", 200)
	//AT命令检测不通过，不进行下一步操作
	if err != nil {
		return
	}
	m.MobileModuleCheckSIM()
	m.MobileModuleGetICCID()
	m.MobileModuleGetIMEI()
	m.MobileModuleGetCSQ()
	m.MobileModuleGetLocation(1)
}

func (m *MobileModuleTemplate) MobileModuleGetParam() {

	m.MobileModuleCheckSIM()
	m.MobileModuleGetICCID()
	m.MobileModuleGetIMEI()
	m.MobileModuleGetCSQ()
	m.MobileModuleGetLocation(0)
}

func (m *MobileModuleTemplate) MobileModuleCheckSIM() {
	err, _ := m.MobileModuleWriteATCommand("AT+CPIN?\r\n", "READY", 1000)
	if err != nil {
		m.RunParam.SIMInsert = false
		ZAPS.Debug("移动模块SIM卡未检测到")
		return
	}
	m.RunParam.SIMInsert = true
	ZAPS.Debug("移动模块SIM卡检测到")
}

func (m *MobileModuleTemplate) MobileModuleGetICCID() {
	err, rxStr := m.MobileModuleWriteATCommand("AT+QCCID\r\n", "+QCCID:", 1000)
	if err != nil {
		m.RunParam.ICCID = ""
		ZAPS.Debug("移动模块ICCID未获取到")
		return
	}

	rxStrSlice := strings.Split(rxStr, "\r\n")
	for _, v := range rxStrSlice {
		if strings.Contains(v, "+QCCID: ") == true {
			index := strings.Index(v, "+QCCID: ") + 8
			m.RunParam.ICCID = v[index : index+19]
			ZAPS.Debugf("移动模块ICCID[%v]", m.RunParam.ICCID)
			break
		}
	}
}

func (m *MobileModuleTemplate) MobileModuleGetIMEI() {
	err, rxStr := m.MobileModuleWriteATCommand("AT+GSN\r\n", "OK", 1000)
	if err != nil {
		return
	}
	rxStrSlice := strings.Split(rxStr, "\r\n")
	for _, v := range rxStrSlice {
		if len(v) > 10 {
			m.RunParam.IMEI = v
			ZAPS.Debugf("移动模块IMEI[%v]", m.RunParam.IMEI)
		}
	}
}

func (m *MobileModuleTemplate) MobileModuleGetCSQ() {
	err, rxStr := m.MobileModuleWriteATCommand("AT+CSQ\r\n", "OK", 1000)
	if err != nil {
		return
	}
	rxStrSlice := strings.Split(rxStr, "\r\n")
	for _, v := range rxStrSlice {
		if strings.Contains(v, "+CSQ:") == true {
			index1 := strings.Index(rxStrSlice[1], "CSQ: ") + 5
			index2 := strings.Index(rxStrSlice[1], ",")
			m.RunParam.CSQ = rxStrSlice[1][index1:index2]

			csq, _ := strconv.Atoi(m.RunParam.CSQ)
			ZAPS.Debugf("移动模块信号强度[%v]", csq)
			if csq < 10 {
				GPIOWriteCmd(m.GPIOParam.RSSILedHigh.GPIO, m.GPIOParam.RSSILedHigh.ResetValue)
				GPIOWriteCmd(m.GPIOParam.RSSILedMiddle.GPIO, m.GPIOParam.RSSILedMiddle.ResetValue)
				GPIOWriteCmd(m.GPIOParam.RSSILedLow.GPIO, m.GPIOParam.RSSILedLow.SetValue)
			} else if csq >= 10 && csq < 25 {
				GPIOWriteCmd(m.GPIOParam.RSSILedHigh.GPIO, m.GPIOParam.RSSILedHigh.ResetValue)
				GPIOWriteCmd(m.GPIOParam.RSSILedMiddle.GPIO, m.GPIOParam.RSSILedMiddle.SetValue)
				GPIOWriteCmd(m.GPIOParam.RSSILedLow.GPIO, m.GPIOParam.RSSILedLow.SetValue)
			} else {
				GPIOWriteCmd(m.GPIOParam.RSSILedHigh.GPIO, m.GPIOParam.RSSILedHigh.SetValue)
				GPIOWriteCmd(m.GPIOParam.RSSILedMiddle.GPIO, m.GPIOParam.RSSILedMiddle.SetValue)
				GPIOWriteCmd(m.GPIOParam.RSSILedLow.GPIO, m.GPIOParam.RSSILedLow.SetValue)
			}
		}
	}
}

func (m *MobileModuleTemplate) MobileModuleGetLocation(cmd int) {

	setupFlag := false

	err, _ := m.MobileModuleWriteATCommand("AT+CREG=2\r\n", "OK", 1000)
	if err != nil {
		return
	}

	err, rxStr := m.MobileModuleWriteATCommand("AT+CREG?\r\n", "OK", 1000)
	if err != nil {
		return
	}

	rxStrSlice := strings.Split(rxStr, "\r\n")
	for _, v := range rxStrSlice {
		if strings.Contains(v, "+CREG") == true {
			tmpSlice := strings.Split(v, ",")
			if len(tmpSlice) == 5 {
				//判断网络注册状态
				if tmpSlice[1] == "1" || tmpSlice[1] == "5" {
					//从未注册变成注册
					if m.RunParam.Network == false {
						setupFlag = true
					}
					m.RunParam.Network = true
					ZAPS.Debug("移动模块网络状态[已注册]")
				} else {
					m.RunParam.Network = false
					ZAPS.Debug("移动模块网络状态[未注册]")
				}
				m.RunParam.Lac = strings.ReplaceAll(tmpSlice[2], "\"", "")
				m.RunParam.Ci = strings.ReplaceAll(tmpSlice[3], "\"", "")
			}
		}
	}

	//重新拨号
	if setupFlag == true || cmd == 1 {
		m.MobileModuleGetOperator()
		m.MobileModuleSetupDataCall()
	}
}

/*
获取运营商信息
*/
func (m *MobileModuleTemplate) MobileModuleGetOperator() {
	err, rxStr := m.MobileModuleWriteATCommand("AT+COPS?\r\n", "+COPS:", 1000)
	if err != nil {
		return
	}

	rxStrSlice := strings.Split(rxStr, "\r\n")
	for _, v := range rxStrSlice {
		if strings.Contains(v, "+COPS:") == true {
			tmpSlice := strings.Split(v, ",")
			if len(tmpSlice) == 4 {
				ZAPS.Debugf("Operator %v", tmpSlice[2])
				//判断网络注册状态
				if strings.Contains(tmpSlice[2], "CHN-UNICOM") { //移动、联通
					m.RunParam.Operator = OperatorTypeCMCCCUCC
					ZAPS.Debug("移动模块运营商[移动]")
				} else if strings.Contains(tmpSlice[2], "CHN-CT") { //电信
					m.RunParam.Operator = OperatorTypeCTCC
					ZAPS.Debug("移动模块运营商[电信]")
				} else {

				}
			}
		}
	}
}

func (m *MobileModuleTemplate) MobileModuleSetupDataCall() {
	ZAPS.Debugf("物联网卡运营商 %v", m.RunParam.Operator)
	err, _ := m.MobileModuleWriteATCommand("AT+QNETDEVCTL=1,1,1\r\n", "OK", 1000)
	if err != nil {
		ZAPS.Debugf("移动网络发送AT拨号指令失败 %v", err.Error())
		return
	}
	//非阻塞,动态获取有可能不成功
	err, errString, outString := RunShellCmd("udhcpc", "-fnq", "-i", "usb0")
	if err != nil {
		ZAPS.Debugf("移动网络拨号失败 %v %s", err.Error(), errString)
	} else {
		ZAPS.Debugf("移动网络拨号成功 %v", outString)
	}
}

func (m *MobileModuleTemplate) MobileModuleFlowAdd(flow int) {

	m.RunParam.CurFlow += float32(flow)
	if m.RunParam.CurFlow <= 1024 {
		m.RunParam.Flow = fmt.Sprintf("%f Byte", m.RunParam.CurFlow)
	} else if m.RunParam.CurFlow > 1024 && m.RunParam.CurFlow <= 1048576 {
		m.RunParam.Flow = fmt.Sprintf("%4.2fK Byte", m.RunParam.CurFlow/1024)
	} else if m.RunParam.CurFlow > 1048576 && m.RunParam.CurFlow < 1073741824 {
		m.RunParam.Flow = fmt.Sprintf("%4.2fM Byte", m.RunParam.CurFlow/1048576)
	} else if m.RunParam.CurFlow > 1099511627776 {
		m.RunParam.Flow = fmt.Sprintf("%4.2fG Byte", m.RunParam.CurFlow/1099511627776)
	}
}

func (m *MobileModuleTemplate) MobileModuleKeepAlive() {

	var err error
	if m.KeepAliveParam.IPMaster != "" {
		err, _ = SendPing(m.KeepAliveParam.IPMaster)

	} else if m.KeepAliveParam.IPSlave != "" {
		err, _ = SendPing(m.KeepAliveParam.IPSlave)
	}

	if err != nil {
		m.KeepAliveParam.CurrentCnt++
		if m.KeepAliveParam.CurrentCnt >= m.KeepAliveParam.OfflineCnt {
			m.KeepAliveParam.CurrentCnt = 0
			m.MobileModuleSetupDataCall()
			m.KeepAliveParam.ResetCnt++
			if m.KeepAliveParam.ResetCnt > 3 {
				m.KeepAliveParam.ResetCnt = 0
				m.MobileModuleClose()
				GPIOWriteCmd(m.GPIOParam.ResetCtr.GPIO, m.GPIOParam.ResetCtr.SetValue)
				time.Sleep(time.Second * 2)
				GPIOWriteCmd(m.GPIOParam.ResetCtr.GPIO, m.GPIOParam.ResetCtr.ResetValue)
				time.Sleep(time.Second * 10)
				m.MobileModuleOpen()
			}
		}
	} else {
		m.KeepAliveParam.CurrentCnt = 0
		m.KeepAliveParam.ResetCnt = 0
	}
}

func ReadMobileModuleParamFromJson() bool {

	data, err := utils.FileRead("./selfpara/mobileModule.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			ZAPS.Debug("打开移动模块配置json文件失败[如果未配置，可以忽略]")
		} else {
			ZAPS.Errorf("打开移动模块配置json文件失败 %v", err)
		}
		return false
	}
	err = json.Unmarshal(data, &MobileModuleParam)
	if err != nil {
		ZAPS.Errorf("移动模块配置json文件格式化失败 %v", err)
		return false
	}
	ZAPS.Debugf("打开移动模块配置json文件成功")

	return true
}

func AddMobileModuleParam(param MobileModuleParamTemplate) error {

	if MobileModuleParam.Name == param.Name {
		return errors.New("移动模块名字已经存在")
	}

	//保存配置参数
	MobileModuleParam = param
	WriteMobileModuleParamToJson()

	//实例化
	MobileModule = NewMobileModule(MobileModuleParam)

	return nil
}

func ModifyMobileModuleParam(param MobileModuleParamTemplate) error {

	if MobileModuleParam.Name != param.Name {
		return errors.New("移动模块名字不存在")
	}

	//保存配置参数
	MobileModuleParam = param
	WriteMobileModuleParamToJson()

	//重新实例化
	MobileModule.CancelFunc()
	MobileModule = NewMobileModule(MobileModuleParam)

	return nil
}

func DeleteMobileModuleParam(name string) error {

	if MobileModuleParam.Name != name {
		return errors.New("移动模块名字不存在")
	}

	//删除配置文件
	utils.DirIsExist("./selfpara")
	err := utils.FileRemove("/selfpara/mobileModule.json")
	if err != nil {
		ZAPS.Errorf("删除移动模块配置json文件 %s %v", "失败", err)
		return err
	}

	//关闭协程
	MobileModule.CancelFunc()

	MobileModule = &MobileModuleTemplate{}

	return nil
}

func WriteMobileModuleParamToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(MobileModuleParam)
	err := utils.FileWrite("./selfpara/mobileModule.json", sJson)
	if err != nil {
		ZAPS.Errorf("写入移动模块配置json文件 %s %v", "失败", err)
		return
	}
	ZAPS.Infof("写入移动模块json文件 %s", "成功")
}
