package commInterface

import (
	"encoding/json"
	"errors"
	"gateway/pin"
	"gateway/setting"
	"gateway/utils"
	"log"
	"strconv"
	"time"

	"github.com/tarm/serial"
)

type SerialInterfaceParam struct {
	Name     string `json:"name"`
	BaudRate string `json:"baudRate"`
	DataBits string `json:"dataBits"` //数据位: 5, 6, 7 or 8 (default 8)
	StopBits string `json:"stopBits"` //停止位: 1 or 2 (default 1)
	Parity   string `json:"parity"`   //校验: N - None, E - Even, O - Odd (default E),(The use of no parity requires 2 stop bits.)
	Timeout  string `json:"timeout"`  //通信超时
	Interval string `json:"interval"` //通信间隔
}

type CommunicationSerialTemplate struct {
	Name   string               `json:"name"`   //接口名称
	Type   string               `json:"type"`   //接口类型,比如serial,tcp,udp,http
	Param  SerialInterfaceParam `json:"param"`  //接口参数
	Port   *serial.Port         `json:"-"`      //通信句柄
	Status ConnectStatus        `json:"status"` //连接状态
}

var CommunicationSerialMap = make([]*CommunicationSerialTemplate, 0)

func (c *CommunicationSerialTemplate) Open() bool {

	//if c.Param.Name == "/dev/ttyS4" {
	//	setting.Exec_shell("echo 95 > /sys/class/gpio/export")
	//} else if c.Param.Name == "/dev/ttyS5" {
	//	setting.Exec_shell("echo 94 > /sys/class/gpio/export")
	//}

	if c.Param.Name == "/dev/ttyS4" {
		if pin.Rs485PinInit(pin.Rs4851Pin) == false {
			setting.ZAPS.Errorf("[%s]初始化RS485控制脚失败", c.Param.Name)
		}
	} else if c.Param.Name == "/dev/ttyS5" {
		if pin.Rs485PinInit(pin.Rs4852Pin) == false {
			setting.ZAPS.Errorf("[%s]初始化RS485控制脚失败", c.Param.Name)
		}
	}

	serialParam := c.Param
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
	c.Port, err = serial.OpenPort(serialConfig)
	if err != nil {
		setting.ZAPS.Errorf("通信串口接口[%s]打开失败 %v", c.Param.Name, err)
		c.Status = CommIsUnConnect
		return false
	}
	setting.ZAPS.Debugf("通信串口接口[%s]打开成功", c.Param.Name)
	c.Status = CommIsConnect

	return true
}

func (c *CommunicationSerialTemplate) Close() bool {

	if c.Port == nil {
		setting.ZAPS.Errorf("通信串口接口[%s] 文件句柄不存在", c.Param.Name)
		return false
	}
	err := c.Port.Close()
	if err != nil {
		setting.ZAPS.Errorf("通信串口接口[%s]关闭失败 %v", c.Param.Name, err)
		return false
	}
	setting.ZAPS.Infof("通信串口接口[%s]关闭成功", c.Param.Name)
	c.Status = CommIsUnConnect
	return true
}

func (c *CommunicationSerialTemplate) WriteData(data []byte) ([]byte, error) {

	if c.Port == nil {
		setting.ZAPS.Errorf("通信串口接口[%s] 文件句柄不存在", c.Param.Name)
		return nil, nil
	}

	//if c.Param.Name == "/dev/ttyS4" {
	//	setting.Exec_shell("echo high > /sys/class/gpio/gpio95/direction")
	//} else if c.Param.Name == "/dev/ttyS5" {
	//	setting.Exec_shell("echo high > /sys/class/gpio/gpio94/direction")
	//}

	if c.Param.Name == "/dev/ttyS4" {
		pin.Rs485xTX(pin.Rs4851Pin)
	} else if c.Param.Name == "/dev/ttyS5" {
		pin.Rs485xTX(pin.Rs4852Pin)
	}

	_, err := c.Port.Write(data)
	if err != nil {
		log.Println(err)
	}

	return nil, nil
}

func (c *CommunicationSerialTemplate) ReadData(data []byte) ([]byte, error) {

	if c.Port == nil {
		return nil, errors.New("端口不存在")
	}

	//if c.Param.Name == "/dev/ttyS4" {
	//	setting.Exec_shell("echo low > /sys/class/gpio/gpio95/direction")
	//} else if c.Param.Name == "/dev/ttyS5" {
	//	setting.Exec_shell("echo low > /sys/class/gpio/gpio94/direction")
	//}

	if c.Param.Name == "/dev/ttyS4" {
		pin.Rs485xRX(pin.Rs4851Pin)
	} else if c.Param.Name == "/dev/ttyS5" {
		pin.Rs485xRX(pin.Rs4852Pin)
	}

	cnt, err := c.Port.Read(data)

	return data[:cnt], err
}

func (c *CommunicationSerialTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationSerialTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationSerialTemplate) GetInterval() string {
	return c.Param.Interval
}

func (c *CommunicationSerialTemplate) GetType() int {
	return CommTypeSerial
}

//func NewCommunicationSerialTemplate(commName, commType string, param SerialInterfaceParam) *CommunicationSerialTemplate {
//
//	return &CommunicationSerialTemplate{
//		Param: param,
//		CommunicationTemplate: CommunicationTemplate{
//			Name: commName,
//			Type: commType,
//		},
//	}
//}

func ReadCommSerialInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commSerialInterface.json")
	if err != nil {
		setting.ZAPS.Debugf("打开通信接口[串口]配置json文件失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &CommunicationSerialMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[串口]配置json文件格式化失败 %v", err)
		return false
	}
	setting.ZAPS.Debugf("打开通信接口[串口]配置json文件成功")
	return true
}

func WriteCommSerialInterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationSerialMap)
	err := utils.FileWrite("./selfpara/commSerialInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入通信接口[串口配置]json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入通信接口[串口配置]json文件 %s", "成功")
}
