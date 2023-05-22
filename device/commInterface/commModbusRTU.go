package commInterface

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/pin"
	"gateway/protocol/modbus"
	"gateway/setting"
	"gateway/utils"
	"github.com/tarm/serial"
	"math"
	"strconv"
	"time"
)

//功能十进制	十六进制	       英文	                中文	          最小数据单位
//   01	     0x01	Read Coils	              读多个线圈	     bit（布尔值）
//   05	     0x05	Write Single Coil	      写单个线圈	     bit（布尔值）
//   15	     0x0F	Write Multiple Coils      写多个线圈	     bit（布尔值）
//   02	     0x02	Read Discrete Inputs      读多个离散输入	 bit（布尔值）

//   04	     0x04	Read Input Registers	  读多个输入寄存器	 16bit
//   03	     0x03	Read Holding Registers	  读多个保持寄存器	 16bit
//   06	     0x06	Write Single Register	  写单个保持寄存器	 16bit
//   16	     0x10	Write Multiple Registers  写多个保持寄存器	 16bit

type MBRTUInterfaceParam struct { //CommunicationMBRTUTemplate
	Name     string `json:"name"`
	BaudRate string `json:"baudRate"`
	DataBits string `json:"dataBits"` //数据位: 5, 6, 7 or 8 (default 8)
	StopBits string `json:"stopBits"` //停止位: 1 or 2 (default 1)
	Parity   string `json:"parity"`   //校验: N - None, E - Even, O - Odd (default E),(The use of no parity requires 2 stop bits.)
	Timeout  string `json:"timeout"`  //通信超时
	Interval string `json:"interval"` //通信间隔
}

type MBRTUInterfaceParamTemplate struct {
	SlaveID      byte                 `json:"slaveID"`      //从机ID
	FuncCode     byte                 `json:"funcCode"`     //功能码
	StartRegAddr int                  `json:"startRegAddr"` //起始地址
	TotalRegCnt  int                  `json:"totalRegCnt"`  //寄存器数量
	Param        []MBRTUParamTemplate `json:"param"`        //数据内容
}

type MBRTUParamTemplate struct {
	Rule    string      `json:"rule"`    //Big-endian:ABCD(32),ABCDEFGH(64);Little-endian:DCBA(32),HGFEDCBA(64);Big-endian byte swap:BADC(32),BADCFEHG(64);Little-endian byte swap:CDAB(32),GHEFCDAB(64);
	RegAddr int         `json:"regAddr"` //寄存器地址
	Type    int         `json:"type"`    //数据类型，0表示无符号整形，1表示有符号整形，2表示double，3表示字符串
	RegCnt  int         `json:"regCnt"`  //长度
	Value   interface{} `json:"value"`   //值
}

type CommunicationMBRTUTemplate struct {
	Name   string              `json:"name"`   //接口名称
	Type   string              `json:"type"`   //接口类型,比如serial,tcp,udp,http
	Status ConnectStatus       `json:"status"` //连接状态
	Param  MBRTUInterfaceParam `json:"param"`  //接口参数
	Client modbus.Client       `json:"-"`
}

var CommunicationMBRTUMap = make([]*CommunicationMBRTUTemplate, 0)

func (c *CommunicationMBRTUTemplate) Open() bool {

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

	serialBaud, _ := strconv.Atoi(c.Param.BaudRate)
	//serialDataBits, _ := strconv.Atoi(c.Param.DataBits)

	var serialParity serial.Parity
	switch c.Param.Parity {
	case "N":
		serialParity = serial.ParityNone
	case "O":
		serialParity = serial.ParityOdd
	case "E":
		serialParity = serial.ParityEven
	}

	var serialStop serial.StopBits
	switch c.Param.StopBits {
	case "1":
		serialStop = serial.Stop1
	case "1.5":
		serialStop = serial.Stop1Half
	case "2":
		serialStop = serial.Stop2
	}

	serialConfig := serial.Config{
		Name:        c.Param.Name,
		Baud:        serialBaud,
		Parity:      serialParity,
		StopBits:    serialStop,
		ReadTimeout: time.Millisecond * 1,
	}

	timeout, _ := strconv.Atoi(c.Param.Timeout)
	//p := modbus.NewRTUClientProvider(modbus.WithEnableLogger(),
	//	modbus.WithSerialConfig(serialConfig),
	//	modbus.WithTCPTimeout(time.Duration(timeout)*time.Millisecond))
	p := modbus.NewRTUClientProvider(modbus.WithSerialConfig(serialConfig),
		modbus.WithTCPTimeout(time.Duration(timeout)*time.Millisecond))

	c.Client = modbus.NewClient(p)
	err := c.Client.Connect()
	if err != nil {
		setting.ZAPS.Errorf("通信接口ModbusRTU[%s]打开失败 %v", c.Name, err)
		c.Status = CommIsUnConnect
		return false
	}
	c.Status = CommIsConnect
	setting.ZAPS.Infof("通信接口ModbusRTU[%s]打开成功", c.Name)
	return true
}

func (c *CommunicationMBRTUTemplate) Close() bool {

	err := c.Client.Close()
	if err != nil {
		fmt.Println("Close failed: ", err)
		return false
	}
	return true
}

func (c *CommunicationMBRTUTemplate) FloatToBytes(f float64, rule string) []byte {
	switch rule {
	case "Int_AB":
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(f))
		return buf
	case "Int_BA":
		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, uint16(f))
		return buf
	case "Long_ABCD":
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(f))
		return buf
	case "Long_BADC":
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(f))
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		return buf
	case "Long_DCBA":
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(f))
		return buf
	case "Long_CDAB":
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(f))
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		return buf
	case "Float_ABCD":
		bits := math.Float32bits(float32(f))
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, bits)
		return buf
	case "Float_BADC":
		bits := math.Float32bits(float32(f))
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, bits)
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		return buf
	case "Float_DCBA":
		bits := math.Float32bits(float32(f))
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, bits)
		return buf
	case "Float_CDAB":
		bits := math.Float32bits(float32(f))
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, bits)
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		return buf
	case "Double_ABCDEFGH":
		bits := math.Float64bits(f)
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, bits)
		return buf
	case "Double_GHEFCDAB":
		bits := math.Float64bits(f)
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, bits)
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		buf[4], buf[5] = buf[5], buf[4]
		buf[6], buf[7] = buf[7], buf[6]
		return buf
	case "Double_BADCFEHG":
		bits := math.Float64bits(f)
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, bits)
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		buf[4], buf[5] = buf[5], buf[4]
		buf[6], buf[7] = buf[7], buf[6]
		return buf
	case "Double_HGFEDCBA":
		bits := math.Float64bits(f)
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, bits)
		return buf
	default:
		return nil
	}
}

func (c *CommunicationMBRTUTemplate) BytesToFloat(buf []byte, rule string) float64 {
	var value float64

	switch rule {
	case "Int_AB":
		value = float64(binary.BigEndian.Uint16(buf))
	case "Int_BA":
		value = float64(binary.LittleEndian.Uint16(buf))
	case "Long_ABCD":
		value = float64(binary.BigEndian.Uint32(buf))
	case "Long_BADC":
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		value = float64(binary.BigEndian.Uint32(buf))
	case "Long_DCBA":
		value = float64(binary.LittleEndian.Uint32(buf))
	case "Long_CDAB":
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		value = float64(binary.LittleEndian.Uint32(buf))
	case "Float_ABCD":
		bits := binary.BigEndian.Uint32(buf)
		value = float64(math.Float32frombits(bits))
	case "Float_BADC":
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		bits := binary.BigEndian.Uint32(buf)
		value = float64(math.Float32frombits(bits))
	case "Float_DCBA":
		bits := binary.LittleEndian.Uint32(buf)
		value = float64(math.Float32frombits(bits))
	case "Float_CDAB":
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		bits := binary.LittleEndian.Uint32(buf)
		value = float64(math.Float32frombits(bits))
	case "Double_ABCDEFGH":
		bits := binary.BigEndian.Uint64(buf)
		value = math.Float64frombits(bits)
	case "Double_GHEFCDAB":
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		buf[4], buf[5] = buf[5], buf[4]
		buf[6], buf[7] = buf[7], buf[6]
		bits := binary.LittleEndian.Uint64(buf)
		value = math.Float64frombits(bits)
	case "Double_BADCFEHG":
		buf[0], buf[1] = buf[1], buf[0]
		buf[2], buf[3] = buf[3], buf[2]
		buf[4], buf[5] = buf[5], buf[4]
		buf[6], buf[7] = buf[7], buf[6]
		bits := binary.BigEndian.Uint64(buf)
		value = math.Float64frombits(bits)
	case "Double_HGFEDCBA":
		bits := binary.LittleEndian.Uint64(buf)
		value = math.Float64frombits(bits)
	}
	return value
}

func (c *CommunicationMBRTUTemplate) BytesToBinary(bs []byte) []byte {
	buf := bytes.NewBuffer([]byte{})

	for _, v := range bs {
		buf.WriteString(fmt.Sprintf("%08b", v))
	}
	b := Reverse(buf.String())
	var Binary []byte
	for _, v := range b {
		Binary = append(Binary, byte(v-48))
	}
	return Binary
}

func (c *CommunicationMBRTUTemplate) WriteData(data []byte) ([]byte, error) {

	dataWrite := MBRTUInterfaceParamTemplate{}

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

	err := json.Unmarshal(data, &dataWrite)
	if err != nil {
		setting.ZAPS.Errorf("通信接口ModbusRTU[%s]写变量JSON格式化错误%v", c.Name, err)
		return nil, err
	}
	setting.ZAPS.Debugf("通信接口ModbusRTU[%s]写变量命令参数%+v", c.Name, dataWrite)

	SumWriteRegCnt := 0
	for _, WriteRegCheck := range dataWrite.Param {
		SumWriteRegCnt += WriteRegCheck.RegCnt
		if WriteRegCheck.RegAddr < dataWrite.StartRegAddr || WriteRegCheck.RegAddr+WriteRegCheck.RegCnt > dataWrite.StartRegAddr+dataWrite.TotalRegCnt {
			setting.ZAPS.Errorf("通信接口ModbusRTU[%s]寄存器地址或数量超限！--> 命令[StartRegAddr:%v, TotalRegCnt:%v] 命令参数[RegAddr:%v, RegCnt:%v]", c.Name, dataWrite.StartRegAddr, dataWrite.TotalRegCnt, WriteRegCheck.RegAddr, WriteRegCheck.RegCnt)
			return nil, errors.New(fmt.Sprintf("通信接口[%s]寄存器地址或数量超限！", c.Name))
		}
	}
	setting.ZAPS.Debugf("通信接口ModbusRTU[%s]寄存器累计个数[%d]，写命令寄存器个数[%d]", c.Name, SumWriteRegCnt, dataWrite.TotalRegCnt)
	//if SumWriteRegCnt > dataWrite.TotalRegCnt {
	//	setting.ZAPS.Errorf("通信接口ModbusRTU[%s]寄存器累计个数[%d]超过写命令寄存器个数[%d]", c.Name, SumWriteRegCnt, dataWrite.TotalRegCnt)
	//	return nil, errors.New("通信接口ModbusRTU寄存器累计个数超过写命令寄存器个数")
	//}

	params := make([]MBRTUParamTemplate, 0)
	if dataWrite.FuncCode == 15 {
		params, err = c.ProcessWriteCoils(dataWrite)
	} else if dataWrite.FuncCode == 16 {
		params, err = c.ProcessWriteMutilRegs(dataWrite)
	}
	ackParam, _ := json.Marshal(params)

	return ackParam, nil
}

func (c *CommunicationMBRTUTemplate) ProcessWriteCoils(dataWrite MBRTUInterfaceParamTemplate) ([]MBRTUParamTemplate, error) {

	regCnt := 0
	if len(dataWrite.Param)%8 > 0 {
		regCnt = len(dataWrite.Param)/8 + 1
	} else {
		regCnt = len(dataWrite.Param) / 8
	}

	params := make([]byte, regCnt*8)
	var Bin string
	for i := 0; i < len(dataWrite.Param); i++ {
		if dataWrite.Param[i].Rule != "AB" && dataWrite.Param[i].Rule != "BA" {
			setting.ZAPS.Errorf("通信接口ModbusRTU[%s]第%d项Rule=%v，格式错误", c.Name, i+1, dataWrite.Param[i].Rule)
			return nil, errors.New("通信接口ModbusRTURule格式错误")
		}
		if dataWrite.Param[i].Value.(bool) {
			params[i] = 1
			Bin += "1"
		} else {
			params[i] = 0
			Bin += "0"
		}
	}
	setting.ZAPS.Debugf("通信接口ModbusRTU[%s]params[%v]Bin[%v]", c.Name, params, Bin)

	MultiBool := make([]byte, regCnt)
	for i := 0; i < len(params); i += 8 {
		bytesLow := params[i]*8 + params[i+1]*4 + params[i+2]*2 + params[i+3]
		bytesHigh := params[i+4]*8 + params[i+5]*4 + params[i+6]*2 + params[i+7]
		MultiBool[i/8] = bytesLow*16 + bytesHigh
	}

	ackParam := dataWrite.Param
	err := c.Client.WriteMultipleCoils(dataWrite.SlaveID, uint16(dataWrite.StartRegAddr), uint16(dataWrite.TotalRegCnt), MultiBool)
	if err != nil {
		setting.ZAPS.Errorf("通信接口ModbusRTU[%s]写多个线圈错误 %v", c.Name, err)
		for i := 0; i < len(ackParam); i++ {
			ackParam[i].Value = -1
		}
		return nil, err
	} else {
		for i := 0; i < len(ackParam); i++ {
			ackParam[i].Value = 0
		}
	}

	return ackParam, nil
}

func (c *CommunicationMBRTUTemplate) ProcessWriteMutilRegs(dataWrite MBRTUInterfaceParamTemplate) ([]MBRTUParamTemplate, error) {

	params := dataWrite.Param

	MultipleRegistersBytes := make([]byte, 0)
	for i := 0; i < len(dataWrite.Param); i++ {

		if dataWrite.Param[i].Type != 3 {
			MultiFloat64 := c.FloatToBytes(dataWrite.Param[i].Value.(float64), dataWrite.Param[i].Rule)
			MultipleRegistersBytes = append(MultipleRegistersBytes, MultiFloat64...)
		} else {
			if len(dataWrite.Param[i].Value.(string)) > int(dataWrite.Param[i].RegCnt)*2 {
				setting.ZAPS.Errorf("字符串长度：%d大于输入长度：%d\n", len(dataWrite.Param[i].Value.(string)), int(dataWrite.Param[i].RegCnt)*2)
				return nil, nil
			}
			MultiString := []byte(dataWrite.Param[i].Value.(string))
			if len(MultiString)%2 > 0 {
				MultiString = append(MultiString, 0)
			}
			MultipleRegistersBytes = append(MultipleRegistersBytes, MultiString...)
		}
	}

	err := c.Client.WriteMultipleRegistersBytes(dataWrite.SlaveID, uint16(dataWrite.StartRegAddr), uint16(dataWrite.TotalRegCnt), MultipleRegistersBytes)
	if err != nil {
		setting.ZAPS.Errorf("WriteMultipleRegistersBytes Error %v", err)
		for i := 0; i < len(dataWrite.Param); i++ {
			params[i].Value = -1
		}
	} else {
		for i := 0; i < len(dataWrite.Param); i++ {
			params[i].Value = 0
		}
	}

	return params, nil
}

func (c *CommunicationMBRTUTemplate) ReadData(data []byte) ([]byte, error) {

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

	dataRead := MBRTUInterfaceParamTemplate{}
	err := json.Unmarshal(data, &dataRead)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("通信接口[%s]读数据JSON格式化错误，%v", c.Name, err.Error()))
	}
	//setting.ZAPS.Debugf("dataRead %+v", dataRead)

	//判断命令寄存器数量与寄存器数量
	SumReadRegCnt := 0
	for _, ReadRegCheck := range dataRead.Param {
		SumReadRegCnt += ReadRegCheck.RegCnt
		if ReadRegCheck.RegAddr < dataRead.StartRegAddr || ReadRegCheck.RegAddr+ReadRegCheck.RegCnt > dataRead.StartRegAddr+dataRead.TotalRegCnt {
			setting.ZAPS.Errorf("通信接口[%s]寄存器地址或数量超限！--> 命令[StartRegAddr:%v, TotalRegCnt:%v] 命令参数[RegAddr:%v, RegCnt:%v]", c.Name, dataRead.StartRegAddr, dataRead.TotalRegCnt, ReadRegCheck.RegAddr, ReadRegCheck.RegCnt)
			return nil, errors.New(fmt.Sprintf("通信接口[%s]寄存器地址或数量超限！", c.Name))
		}
	}
	setting.ZAPS.Debugf("通信接口ModbusRTU[%s]寄存器累计个数[%d]，读命令寄存器个数[%d]", c.Name, SumReadRegCnt, dataRead.TotalRegCnt)
	//if SumReadRegCnt > dataRead.TotalRegCnt {
	//	return nil, errors.New(fmt.Sprintf("通信接口[%s]寄存器数量超出", c.Name))
	//}

	params := make([]MBRTUParamTemplate, 0)

	if dataRead.FuncCode == 1 {
		params, err = c.ProcessReadCoil(dataRead)
	} else if dataRead.FuncCode == 2 {
		params, err = c.ProcessReadDiscreteInputs(dataRead)
	} else if dataRead.FuncCode == 3 {
		params, err = c.ProcessReadHoldingRegister(dataRead)
	} else if dataRead.FuncCode == 4 {
		params, err = c.ProcessReadInputRegisters(dataRead)
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("通信接口[%s]发送读数据命令错误，%v", c.Name, err.Error()))
	}

	rData, _ := json.Marshal(params)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("通信接口[%s]读数据返回JSON格式化错误，%v", c.Name, err.Error()))
	}

	return rData, err
}

func (c *CommunicationMBRTUTemplate) ProcessReadCoil(dataRead MBRTUInterfaceParamTemplate) ([]MBRTUParamTemplate, error) {

	params := make([]MBRTUParamTemplate, 0)

	results, err := c.Client.ReadCoils(dataRead.SlaveID, uint16(dataRead.StartRegAddr), uint16(dataRead.TotalRegCnt)) //(dataRead.Quantity/8+1)*8
	if err != nil {
		return nil, errors.New(fmt.Sprintf("通信接口[%s]读多个线圈错误，%v", c.Name, err.Error()))
	}
	BinBytes := c.BytesToBinary(results)
	setting.ZAPS.Debugf("BinBytes:%v", BinBytes)
	for _, v := range dataRead.Param {
		params = append(params, MBRTUParamTemplate{
			Rule:    v.Rule,
			RegAddr: v.RegAddr,
			Type:    v.Type,
			RegCnt:  v.RegCnt,
			Value:   BinBytes[v.RegAddr%8],
		})
	}
	return params, nil
}

func (c *CommunicationMBRTUTemplate) ProcessReadDiscreteInputs(dataRead MBRTUInterfaceParamTemplate) ([]MBRTUParamTemplate, error) {

	params := make([]MBRTUParamTemplate, 0)

	results, err := c.Client.ReadDiscreteInputs(dataRead.SlaveID, uint16(dataRead.StartRegAddr), uint16(dataRead.TotalRegCnt))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("通信接口[%s]读离线输入线圈错误，%v", c.Name, err.Error()))
	}

	BinBytes := c.BytesToBinary(results)
	setting.ZAPS.Debugf("BinBytes:%v", BinBytes)
	for _, v := range dataRead.Param {
		params = append(params, MBRTUParamTemplate{
			Rule:    v.Rule,
			RegAddr: v.RegAddr,
			Type:    v.Type,
			RegCnt:  v.RegCnt,
			Value:   BinBytes[v.RegAddr%8],
		})
	}
	return params, nil
}

func (c *CommunicationMBRTUTemplate) ProcessReadHoldingRegister(dataRead MBRTUInterfaceParamTemplate) ([]MBRTUParamTemplate, error) {

	params := make([]MBRTUParamTemplate, 0)

	results, err := c.Client.ReadHoldingRegistersBytes(dataRead.SlaveID, uint16(dataRead.StartRegAddr), uint16(dataRead.TotalRegCnt))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("通信接口[%s]读保持寄存器错误，%v", c.Name, err.Error()))
	}

	var startTemp int //dataRead.Address
	var endTemp int
	var value interface{}

	for i := 0; i < len(dataRead.Param); i++ {
		startTemp = (dataRead.Param[i].RegAddr - dataRead.StartRegAddr) * 2
		endTemp = startTemp + dataRead.Param[i].RegCnt*2
		if dataRead.Param[i].Type != 3 {
			value = c.BytesToFloat(results[startTemp:endTemp], dataRead.Param[i].Rule)
		} else {
			value = string(results[startTemp:endTemp])
		}

		params = append(params, MBRTUParamTemplate{
			Rule:    dataRead.Param[i].Rule,
			RegAddr: dataRead.Param[i].RegAddr,
			Type:    dataRead.Param[i].Type,
			RegCnt:  dataRead.Param[i].RegCnt,
			Value:   value})
	}

	return params, nil
}

func (c *CommunicationMBRTUTemplate) ProcessReadInputRegisters(dataRead MBRTUInterfaceParamTemplate) ([]MBRTUParamTemplate, error) {

	params := make([]MBRTUParamTemplate, 0)

	results, err := c.Client.ReadInputRegistersBytes(dataRead.SlaveID, uint16(dataRead.StartRegAddr), uint16(dataRead.TotalRegCnt))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("通信接口[%s]读输入寄存器错误，%v", c.Name, err.Error()))
	}
	var startTemp int
	var endTemp int
	var value interface{}
	for i := 0; i < len(dataRead.Param); i++ {
		startTemp = (dataRead.Param[i].RegAddr - dataRead.StartRegAddr) * 2
		endTemp = startTemp + dataRead.Param[i].RegCnt*2
		if dataRead.Param[i].Type != 3 {
			value = c.BytesToFloat(results[startTemp:endTemp], dataRead.Param[i].Rule)
		} else {
			value = string(results[startTemp:endTemp])
		}
		params = append(params, MBRTUParamTemplate{
			Rule:    dataRead.Param[i].Rule,
			RegAddr: dataRead.Param[i].RegAddr,
			Type:    dataRead.Param[i].Type,
			RegCnt:  dataRead.Param[i].RegCnt,
			Value:   value})
	}

	return params, nil
}

func (c *CommunicationMBRTUTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationMBRTUTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationMBRTUTemplate) GetInterval() string {
	return c.Param.Interval
}

func (c *CommunicationMBRTUTemplate) GetType() int {
	return CommTypeModbusRTU
}

func ReadCommMBRTUInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commMBRTUInterface.json")
	if err != nil {
		setting.ZAPS.Debugf("打开通信接口[MBRTU]json配置文件失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &CommunicationMBRTUMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[MBRTU]json配置文件格式化失败 %v", err)
		return false
	}
	setting.ZAPS.Debugf("打开通信接口[MBRTU]json配置文件成功")
	return true
}

func WriteCommMBRTUInterfaceListToJson() {
	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationMBRTUMap)
	err := utils.FileWrite("./selfpara/commMBRTUInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入通信接口[MBRTU]json配置文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入通信接口[MBRTU]json配置文件 %s", "成功")
}
