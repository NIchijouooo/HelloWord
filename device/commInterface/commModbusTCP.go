package commInterface

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/protocol/modbus"
	"gateway/setting"
	"gateway/utils"
	"math"
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

type MBTCPInterfaceParam struct { //CommunicationMBTCPTemplate
	IP       string `json:"ip"`       //从机地址
	Port     string `json:"port"`     //从机端口
	Timeout  string `json:"timeout"`  //通信超时
	Interval string `json:"interval"` //通信间隔
}

type MBTCPInterfaceParamTemplate struct {
	SlaveID      byte                 `json:"slaveID"`      //从机ID
	FuncCode     byte                 `json:"funcCode"`     //功能码
	StartRegAddr int                  `json:"startRegAddr"` //起始地址
	TotalRegCnt  int                  `json:"totalRegCnt"`  //寄存器数量
	Param        []MBTCPParamTemplate `json:"param"`        //数据内容
}

type MBTCPParamTemplate struct {
	Rule    string      `json:"rule"`    //Big-endian:ABCD(32),ABCDEFGH(64);Little-endian:DCBA(32),HGFEDCBA(64);Big-endian byte swap:BADC(32),BADCFEHG(64);Little-endian byte swap:CDAB(32),GHEFCDAB(64);
	RegAddr int         `json:"regAddr"` //寄存器地址
	Type    int         `json:"type"`    //数据类型，0表示无符号整形，1表示有符号整形，2表示double，3表示字符串
	RegCnt  int         `json:"regCnt"`  //长度
	Value   interface{} `json:"value"`   //值
}

type CommunicationMBTCPTemplate struct {
	Name   string              `json:"name"`   //接口名称
	Type   string              `json:"type"`   //接口类型,比如serial,tcp,udp,http
	Status ConnectStatus       `json:"status"` //连接状态
	Param  MBTCPInterfaceParam `json:"param"`  //接口参数
	Client modbus.Client       `json:"-"`
}

var CommunicationMBTCPMap = make([]*CommunicationMBTCPTemplate, 0)

func (c *CommunicationMBTCPTemplate) Open() bool {

	p := modbus.NewTCPClientProvider(c.Param.IP+":"+c.Param.Port, modbus.WithEnableLogger())

	c.Client = modbus.NewClient(p)
	err := c.Client.Connect()
	if err != nil {
		setting.ZAPS.Errorf("通信接口ModbusTCP[%s]打开失败 %v", c.Name, err)
		c.Status = CommIsUnConnect
		return false
	}
	c.Status = CommIsConnect
	setting.ZAPS.Infof("通信接口ModbusTCP[%s]打开成功", c.Name)
	return true
}

func (c *CommunicationMBTCPTemplate) Close() bool {

	err := c.Client.Close()
	if err != nil {
		fmt.Println("Close failed: ", err)
		return false
	}
	return true
}

func (c *CommunicationMBTCPTemplate) FloatToBytes(f float64, rule string) []byte {
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

func (c *CommunicationMBTCPTemplate) BytesToFloat(buf []byte, rule string) float64 {
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

func (c *CommunicationMBTCPTemplate) BytesToBinary(bs []byte) []byte {
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

func (c *CommunicationMBTCPTemplate) BytesToBinaryGwai(bs []byte) []byte {
	buf := bytes.NewBuffer([]byte{})

	for _, v := range bs {
		s := Reverse(fmt.Sprintf("%08b", v))
		buf.WriteString(s)
	}

	var Binary []byte
	for _, v := range buf.String() {
		Binary = append(Binary, byte(v-48))
	}

	return Binary
}

func Reverse(s string) string {
	a := []rune(s)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return string(a)
}

func (c *CommunicationMBTCPTemplate) WriteData(data []byte) ([]byte, error) {

	dataWrite := MBTCPInterfaceParamTemplate{}

	err := json.Unmarshal(data, &dataWrite)
	if err != nil {
		setting.ZAPS.Errorf("通信接口ModbusTCP[%s]写变量JSON格式化错误%v", c.Name, err)
		return nil, err
	}
	setting.ZAPS.Debugf("通信接口ModbusTCP[%s]写变量命令参数%+v", c.Name, dataWrite)

	SumWriteRegCnt := 0
	for _, WriteRegCheck := range dataWrite.Param {
		SumWriteRegCnt += WriteRegCheck.RegCnt
		if WriteRegCheck.RegAddr < dataWrite.StartRegAddr || WriteRegCheck.RegAddr+WriteRegCheck.RegCnt > dataWrite.StartRegAddr+dataWrite.TotalRegCnt {
			setting.ZAPS.Errorf("通信接口ModbusTCP[%s]寄存器地址或数量超限！--> 命令[StartRegAddr:%v, TotalRegCnt:%v] 命令参数[RegAddr:%v, RegCnt:%v]", c.Name, dataWrite.StartRegAddr, dataWrite.TotalRegCnt, WriteRegCheck.RegAddr, WriteRegCheck.RegCnt)
			return nil, errors.New(fmt.Sprintf("通信接口[%s]寄存器地址或数量超限！", c.Name))
		}
	}
	setting.ZAPS.Debugf("通信接口ModbusTCP[%s]寄存器累计个数[%d]，写命令寄存器个数[%d]", c.Name, SumWriteRegCnt, dataWrite.TotalRegCnt)
	//if SumWriteRegCnt > dataWrite.TotalRegCnt {
	//	setting.ZAPS.Errorf("通信接口ModbusTCP[%s]寄存器累计个数[%d]超过写命令寄存器个数[%d]", c.Name, SumWriteRegCnt, dataWrite.TotalRegCnt)
	//	return nil, errors.New("通信接口ModbusTCP寄存器累计个数超过写命令寄存器个数")
	//}

	params := make([]MBTCPParamTemplate, 0)
	if dataWrite.FuncCode == 15 {
		params, err = c.ProcessWriteCoils(dataWrite)
	} else if dataWrite.FuncCode == 16 {
		params, err = c.ProcessWriteMutilRegs(dataWrite)
	}
	ackParam, _ := json.Marshal(params)

	return ackParam, nil
}

func (c *CommunicationMBTCPTemplate) ProcessWriteCoils(dataWrite MBTCPInterfaceParamTemplate) ([]MBTCPParamTemplate, error) {

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
			setting.ZAPS.Errorf("通信接口ModbusTCP[%s]第%d项Rule=%v，格式错误", c.Name, i+1, dataWrite.Param[i].Rule)
			//return nil, errors.New("通信接口ModbusTCPRule格式错误")   //gwai del 2023-04-21
		}
		//if dataWrite.Param[i].Value.(bool) {   //gwai del 2023-04-21
		if dataWrite.Param[i].Value.(float64) != 0 {
			params[i] = 1
			Bin += "1"
		} else {
			params[i] = 0
			Bin += "0"
		}
	}
	setting.ZAPS.Debugf("通信接口ModbusTCP[%s]params[%v]Bin[%v]", c.Name, params, Bin)

	MultiBool := make([]byte, regCnt)
	for i := 0; i < len(params); i += 8 {
		//bytesLow := params[i]*8 + params[i+1]*4 + params[i+2]*2 + params[i+3]    //gwai del 2023-04-21
		//bytesHigh := params[i+4]*8 + params[i+5]*4 + params[i+6]*2 + params[i+7]

		bytesLow := params[i] + params[i+1]*2 + params[i+2]*4 + params[i+3]*8
		bytesHigh := params[i+4] + params[i+5]*2 + params[i+6]*4 + params[i+7]*8
		//MultiBool[i/8] = bytesLow*16 + bytesHigh   //gwai del 2023-04-21
		MultiBool[i/8] = bytesLow + bytesHigh*16
	}

	ackParam := dataWrite.Param
	err := c.Client.WriteMultipleCoils(dataWrite.SlaveID, uint16(dataWrite.StartRegAddr), uint16(dataWrite.TotalRegCnt), MultiBool)
	if err != nil {
		setting.ZAPS.Errorf("通信接口ModbusTCP[%s]写多个线圈错误 %v", c.Name, err)
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

func (c *CommunicationMBTCPTemplate) ProcessWriteMutilRegs(dataWrite MBTCPInterfaceParamTemplate) ([]MBTCPParamTemplate, error) {

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

func (c *CommunicationMBTCPTemplate) ReadData(data []byte) ([]byte, error) {

	dataRead := MBTCPInterfaceParamTemplate{}
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
	setting.ZAPS.Debugf("通信接口ModbusTCP[%s]寄存器累计个数[%d]，读命令寄存器个数[%d]", c.Name, SumReadRegCnt, dataRead.TotalRegCnt)
	//if SumReadRegCnt > dataRead.TotalRegCnt {
	//	return nil, errors.New(fmt.Sprintf("通信接口[%s]寄存器数量超出", c.Name))
	//}

	params := make([]MBTCPParamTemplate, 0)

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
		setting.ZAPS.Errorf("通信接口[%s]发送读数据命令错误 %v", c.Name, err)
		return nil, errors.New(fmt.Sprintf("通信接口[%s]发送读数据命令错误，%v", c.Name, err.Error()))
	}

	rData, _ := json.Marshal(params)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("通信接口[%s]读数据返回JSON格式化错误，%v", c.Name, err.Error()))
	}

	return rData, err
}

func (c *CommunicationMBTCPTemplate) ProcessReadCoil(dataRead MBTCPInterfaceParamTemplate) ([]MBTCPParamTemplate, error) {

	params := make([]MBTCPParamTemplate, 0)

	results, err := c.Client.ReadCoils(dataRead.SlaveID, uint16(dataRead.StartRegAddr), uint16(dataRead.TotalRegCnt)) //(dataRead.Quantity/8+1)*8
	if err != nil {
		return nil, errors.New(fmt.Sprintf("通信接口[%s]读多个线圈错误，%v", c.Name, err.Error()))
	}
	BinBytes := c.BytesToBinaryGwai(results) //gwai add
	setting.ZAPS.Debugf("BinBytes:%v", BinBytes)
	for _, v := range dataRead.Param {
		//setting.ZAPS.Debugf("-----RegAddr[%d][%d][%d]", v.RegAddr, v.RegAddr-dataRead.StartRegAddr+dataRead.StartRegAddr%8, BinBytes[v.RegAddr-dataRead.StartRegAddr+dataRead.StartRegAddr%8])
		params = append(params, MBTCPParamTemplate{
			Rule:    v.Rule,
			RegAddr: v.RegAddr,
			Type:    v.Type,
			RegCnt:  v.RegCnt,
			//Value:   BinBytes[v.RegAddr%8],    //gwai del
			Value: BinBytes[v.RegAddr-dataRead.StartRegAddr], //gwai add
		})
	}
	return params, nil
}

func (c *CommunicationMBTCPTemplate) ProcessReadDiscreteInputs(dataRead MBTCPInterfaceParamTemplate) ([]MBTCPParamTemplate, error) {

	params := make([]MBTCPParamTemplate, 0)

	results, err := c.Client.ReadDiscreteInputs(dataRead.SlaveID, uint16(dataRead.StartRegAddr), uint16(dataRead.TotalRegCnt))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("通信接口[%s]读离线输入线圈错误，%v", c.Name, err.Error()))
	}

	BinBytes := c.BytesToBinary(results)
	setting.ZAPS.Debugf("BinBytes:%v", BinBytes)
	for _, v := range dataRead.Param {
		params = append(params, MBTCPParamTemplate{
			Rule:    v.Rule,
			RegAddr: v.RegAddr,
			Type:    v.Type,
			RegCnt:  v.RegCnt,
			Value:   BinBytes[v.RegAddr%8],
		})
	}
	return params, nil
}

func (c *CommunicationMBTCPTemplate) ProcessReadHoldingRegister(dataRead MBTCPInterfaceParamTemplate) ([]MBTCPParamTemplate, error) {

	params := make([]MBTCPParamTemplate, 0)

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

		params = append(params, MBTCPParamTemplate{
			Rule:    dataRead.Param[i].Rule,
			RegAddr: dataRead.Param[i].RegAddr,
			Type:    dataRead.Param[i].Type,
			RegCnt:  dataRead.Param[i].RegCnt,
			Value:   value})
	}

	return params, nil
}

func (c *CommunicationMBTCPTemplate) ProcessReadInputRegisters(dataRead MBTCPInterfaceParamTemplate) ([]MBTCPParamTemplate, error) {

	params := make([]MBTCPParamTemplate, 0)

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
		params = append(params, MBTCPParamTemplate{
			Rule:    dataRead.Param[i].Rule,
			RegAddr: dataRead.Param[i].RegAddr,
			Type:    dataRead.Param[i].Type,
			RegCnt:  dataRead.Param[i].RegCnt,
			Value:   value})
	}

	return params, nil
}

func (c *CommunicationMBTCPTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationMBTCPTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationMBTCPTemplate) GetInterval() string {
	return c.Param.Interval
}

func (c *CommunicationMBTCPTemplate) GetType() int {
	return CommTypeModbusTCP
}

func ReadCommMBTCPInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commMBTCPInterface.json")
	if err != nil {
		setting.ZAPS.Debugf("打开通信接口[Modbus-TCP]json配置文件失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &CommunicationMBTCPMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[Modbus-TCP]json配置文件格式化失败 %v", err)
		return false
	}
	setting.ZAPS.Debugf("打开通信接口[Modbus-TCP]json配置文件成功")
	return true
}

func WriteCommMBTCPInterfaceListToJson() {
	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationMBTCPMap)
	err := utils.FileWrite("./selfpara/commMBTCPInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入通信接口[MBTCP]json配置文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入通信接口[MBTCP]json配置文件 %s", "成功")
}
