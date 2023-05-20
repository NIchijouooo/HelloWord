package commInterface

import (
	"encoding/binary"
	"encoding/json"
	"gateway/setting"
	"gateway/utils"
	"github.com/pkg/errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/robinson/gos7"
)

type S7ParamTemplate struct {
	Name    string      `json:"name"`
	Address string      `json:"address"` //PLC写入地址
	Start   string      `json:"start"`   //PLC写入起始位
	Type    int         `json:"type"`    //PLC写入类型
	Value   interface{} `json:"value"`   //PLC写入值
}

type S7ReadVariableTemplate struct {
	Address string `json:"address"` //变量地址
	Start   string `json:"start"`   //变量起始位
	Type    string `json:"type"`    //变量类型
}

type S7InterfaceParam struct {
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Timeout string `json:"timeout"` //通信超时
	Connect string `json:"connect"` //连接时间
	PLCType int    `json:"plcType"` //
}

type CommunicationS7Template struct {
	Name   string                 `json:"name"`  //接口名称
	Type   string                 `json:"type"`  //接口类型,比如serial,S7,udp,http
	Param  S7InterfaceParam       `json:"param"` //接口参数
	Conn   *gos7.TCPClientHandler `json:"-"`     //通信句柄
	Client gos7.Client            `json:"-"`
	Status ConnectStatus          `json:"status"` //连接状态
}

type PLCVariableType int

const (
	VariableTypeuint8 int = iota
	VariableTypeint8
	VariableTypeuint16
	VariableTypeint16
	VariableTypeuint32
	VariableTypeint32
	VariableTypefloat
	VariableTypedouble
	VariableTypebool
)

const (
	PLCTypeSmart200 int = 0
)

var CommunicationS7Map = make([]*CommunicationS7Template, 0)

func (c *CommunicationS7Template) Open() bool {

	switch c.Param.PLCType {
	case PLCTypeSmart200:
		c.Conn = gos7.NewTCPClientHandler(c.Param.IP, 0, 1)
	default:
		c.Conn = gos7.NewTCPClientHandler(c.Param.IP, 0, 1)
	}
	//c.Conn.Timeout, _ = time.ParseDuration("1")
	timeOut, _ := time.ParseDuration(c.Param.Timeout)
	c.Conn.Timeout = timeOut * time.Millisecond
	c.Conn.IdleTimeout = 200 * time.Second
	//c.Conn.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	go c.Connect()

	return true
}

func (c *CommunicationS7Template) Connect() {
	err := c.Conn.Connect()
	if err != nil {
		setting.ZAPS.Errorf("通信S7客户端[%s]接口打开失败 %v", c.Name, err)
		return
	}

	c.Status = CommIsConnect
	c.Client = gos7.NewClient(c.Conn)
	//c.Conn.Timeout, _ = time.ParseDuration(c.Param.Timeout)
}

func (c *CommunicationS7Template) ReConnect() {
	err := c.Conn.Close()
	if err != nil {
		setting.ZAPS.Errorf("通信S7客户端[%s]接口关闭失败 %v", c.Name, err)
	}

	err = c.Conn.Connect()
	if err != nil {
		setting.ZAPS.Errorf("通信S7客户端[%s]接口打开失败 %v", c.Name, err)
		return
	}

	c.Status = CommIsConnect
	c.Client = gos7.NewClient(c.Conn)
	//c.Conn.Timeout, _ = time.ParseDuration(c.Param.Timeout)
}

func (c *CommunicationS7Template) Close() bool {
	if c.Conn != nil {
		err := c.Conn.Close()
		if err != nil {
			setting.ZAPS.Errorf("通信S7客户端[%s]接口关闭失败 %v", c.Name, err)
			c.Status = CommIsUnConnect
			return false
		}
		c.Status = CommIsUnConnect
		setting.ZAPS.Infof("通信S7客户端[%s]接口关闭成功", c.Name)
	}
	return true
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, bits)
	return bytes
}
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)
	return bytes
}

func GetValueFromParam(param S7ParamTemplate) (int, int, []byte) {

	start, err := strconv.Atoi(param.Start)
	if err != nil {
		return -1, -1, nil
	}

	switch param.Value.(type) {
	case string:
		value, _ := strconv.ParseFloat(param.Value.(string), 64)
		param.Value = value
	}

	length := 0
	buffer := make([]byte, 0)
	switch param.Type {
	case VariableTypeuint8:
		length = 1
		value := int8(param.Value.(float64))
		buffer = append(buffer, uint8(value))
	case VariableTypeint8:
		length = 1
		value := int8(param.Value.(float64))
		buffer = append(buffer, byte(value))
	case VariableTypeint16:
		length = 2
		value := int16(param.Value.(float64))
		buffer = append(buffer, byte(value/256))
		buffer = append(buffer, byte(value%256))
	case VariableTypeuint16:
		length = 2
		value := uint16(param.Value.(float64))
		buffer = append(buffer, byte(value/256))
		buffer = append(buffer, byte(value%256))
	case VariableTypeint32:
		length = 4
		value := int32(param.Value.(float64))
		buffer = append(buffer, byte(value>>24))
		buffer = append(buffer, byte(value>>16))
		buffer = append(buffer, byte(value>>8))
		buffer = append(buffer, byte(value))
	case VariableTypeuint32:
		length = 4
		value := uint32(param.Value.(float64))
		buffer = append(buffer, byte(value>>24))
		buffer = append(buffer, byte(value>>16))
		buffer = append(buffer, byte(value>>8))
		buffer = append(buffer, byte(value))
	case VariableTypefloat:
		length = 4
		value := float32(param.Value.(float64))
		buffer = Float32ToByte(value)
	case VariableTypedouble:
		length = 8
		value := param.Value.(float64)
		buffer = Float64ToByte(value)
	default:
		return -1, -1, nil
	}
	return start, length, buffer
}

func (c *CommunicationS7Template) WriteData(data []byte) ([]byte, error) {

	param := S7ParamTemplate{}
	err := json.Unmarshal(data, &param)
	if err != nil {
		setting.ZAPS.Errorf("通信S7客户端[%s]接口写变量JSON格式错误 %v", c.Name, err)
		return nil, err
	}

	if strings.Contains(param.Address, "DB") && len(param.Address) >= 3 {
		DBNum, _ := strconv.Atoi(param.Address[2:])
		if err != nil {
			setting.ZAPS.Errorf("通信S7客户端[%s]接口写变量地址格式错误 %v", c.Name, err)
			param.Value = -1
			return nil, err
		}
		if param.Type == VariableTypebool {
			//写入数据的字节二位数组
			//buffers := [][]byte{
			//	make([]byte, 2),
			//	make([]byte, 2),
			//	make([]byte, 4),
			//	make([]byte, 256),
			//	make([]byte, 512),
			//}
			buffer := make([]byte, 1)

			VB, errVB := strconv.Atoi(strings.Split(param.Start, ".")[0])
			if errVB != nil {
				setting.ZAPS.Errorf("通信S7客户端[%s]接口读bool变量起始地址错误 %v", c.Name, errVB)
				return nil, errVB
			}

			bit, errBit := strconv.Atoi(strings.Split(param.Start, ".")[1])
			if errBit != nil || bit < 0 {
				setting.ZAPS.Errorf("通信S7客户端[%s]接口读bool变量起始地址位数错误 %v", c.Name, errBit)
				return nil, errBit
			}
			//生成需要写入的变量的数组
			value, err := strconv.Atoi(param.Value.(string))
			if err != nil || bit < 0 {
				setting.ZAPS.Errorf("通信S7客户端[%s]接口读bool变量写入值错误 %v", c.Name, errBit)
				return nil, errBit
			}

			err = c.Client.AGReadDB(DBNum, VB, 1, buffer)
			if err != nil {
				setting.ZAPS.Errorf("通信S7客户端[%s]接口读变量错误 %v", c.Name, err)
				return nil, err
			} else {
				param.Value = 0
			}

			setting.ZAPS.Debugf("buffer:[%d]", buffer[0])
			switch value {
			case 0:
				buffer[0] = buffer[0] &^ (1 << bit)
				setting.ZAPS.Debugf("case 1,bit[%d]buffer[%d]", bit, buffer[0])
			case 1:
				buffer[0] = buffer[0] | (1 << bit)
				setting.ZAPS.Debugf("case 0,bit[%d]buffer[%d]", bit, buffer[0])
			}
			//s7.SetValueAt(buffers[1], 0, uint16(66))
			//s7.SetRealAt(buffers[2], 0, float32(33.33))
			err = c.Client.AGWriteDB(DBNum, VB, 1, buffer)
			if err != nil {
				setting.ZAPS.Errorf("通信S7客户端[%s]接口写变量错误 %v", c.Name, err)
				return nil, err
			} else {
				param.Value = 0
			}
			//获取批量写入的DataItem
			//datas := []gos7.S7DataItem{
			//	{
			//		Area:     0x84,
			//		WordLen:  0x01,
			//		DBNumber: 1,
			//		Start:    0,
			//		Amount:   1,
			//		Data:     buffer,
			//	},
			//}
			//err = c.Client.AGWriteMulti(datas, len(datas))
			//if err != nil {
			//	return nil, err
			//}

		} else {
			start, length, buffer := GetValueFromParam(param)
			setting.ZAPS.Debugf("通信S7客户端[%s]接口写变量地址[%d]偏移地址[%d]变量类型[%d]变量内容[%d]",
				c.Name, DBNum, start, length, buffer)
			err = c.Client.AGWriteDB(DBNum, start, length, buffer)
			if err != nil {
				setting.ZAPS.Errorf("通信S7客户端[%s]接口写变量错误 %v", c.Name, err)
				param.Value = -1
			} else {
				param.Value = 0
			}
		}
	} else {
		param.Value = -1
	}

	aBuffer, _ := json.Marshal(param)

	return aBuffer, nil
}

func GetBufferLength(param S7ParamTemplate) (int, []byte) {
	buffer := make([]byte, 0)
	length := 0
	switch param.Type {
	case VariableTypebool:
		length = 1
		buffer = make([]byte, 2)
	case VariableTypeuint8:
		fallthrough
	case VariableTypeint8:
		length = 1
		buffer = make([]byte, 2)
	case VariableTypeuint16:
		fallthrough
	case VariableTypeint16:
		length = 2
		buffer = make([]byte, 2)
	case VariableTypeuint32:
		fallthrough
	case VariableTypeint32:
		fallthrough
	case VariableTypefloat:
		length = 4
		buffer = make([]byte, 4)
	case VariableTypedouble:
		length = 8
		buffer = make([]byte, 8)
	default:
		return -1, nil
	}
	return length, buffer
}

func (c *CommunicationS7Template) ReadData(data []byte) ([]byte, error) {

	if c.Status == CommIsUnConnect {
		return nil, errors.New("通信S7客户端未打开")
	}

	s7Param := S7ParamTemplate{}
	aDatas := make([]byte, 0)
	err := json.Unmarshal(data, &s7Param)
	if err != nil {
		setting.ZAPS.Errorf("通信S7客户端[%s]接口读变量JSON格式错误 %v", c.Name, err)
		return nil, err
	}

	length, buffer := GetBufferLength(s7Param)
	if length == -1 {
		setting.ZAPS.Errorf("通信S7客户端[%s]接口读变量变量类型错误 %v", c.Name, err)
		return nil, err
	}

	if strings.Contains(s7Param.Address, "DB") && len(s7Param.Address) >= 3 {
		DBNum, err := strconv.Atoi(s7Param.Address[2:])
		if err != nil {
			setting.ZAPS.Errorf("通信S7客户端[%s]接口读变量数据块地址错误 %v", c.Name, err)
			return nil, err
		}
		//setting.ZAPS.Debugf("buffer %v", len(buffer))

		var bit int
		if s7Param.Type == VariableTypebool {
			VB, err := strconv.Atoi(strings.Split(s7Param.Start, ".")[0])
			bit, err = strconv.Atoi(strings.Split(s7Param.Start, ".")[1])
			if err != nil {
				setting.ZAPS.Errorf("通信S7客户端[%s]接口读bool变量起始地址错误 %v", c.Name, err)
				return nil, err
			}
			err = c.Client.AGReadDB(DBNum, VB, length, buffer)
			if err != nil {
				setting.ZAPS.Errorf("通信S7客户端[%s]接口读变量数据块错误 %v", c.Name, err)
				if strings.Contains(err.Error(), "write: broken pipe") {
					c.ReConnect()
				}
				return nil, err
			}
		} else {
			VWorVD, err := strconv.Atoi(s7Param.Start)
			if err != nil {
				setting.ZAPS.Errorf("通信S7客户端[%s]接口读变量变量起始地址错误 %v", c.Name, err)
				return nil, err
			}
			err = c.Client.AGReadDB(DBNum, VWorVD, length, buffer)
			if err != nil {
				setting.ZAPS.Errorf("通信S7客户端[%s]接口读变量数据块错误 %v", c.Name, err)
				if strings.Contains(err.Error(), "write: broken pipe") {
					c.ReConnect()
				}
				return nil, err
			}
		}

		s7 := gos7.Helper{}
		switch s7Param.Type {
		case VariableTypeuint8:
			var result uint8
			s7.GetValueAt(buffer, 0, &result)
			s7Param.Value = result
		case VariableTypeint8:
			var result int8
			s7.GetValueAt(buffer, 0, &result)
			s7Param.Value = result
		case VariableTypeuint16:
			var result uint16
			s7.GetValueAt(buffer, 0, &result)
			s7Param.Value = result
		case VariableTypeint16:
			var result int16
			s7.GetValueAt(buffer, 0, &result)
			s7Param.Value = result
		case VariableTypeuint32:
			var result uint32
			s7.GetValueAt(buffer, 0, &result)
			s7Param.Value = result
		case VariableTypeint32:
			var result int32
			s7.GetValueAt(buffer, 0, &result)
			s7Param.Value = result
		case VariableTypefloat:
			var result float32
			s7.GetValueAt(buffer, 0, &result)
			s7Param.Value = result
		case VariableTypebool:
			result := ByteToBool(buffer[0])[bit]
			//setting.ZAPS.Debugf("buffer:%d,bit:%d,result: %v", buffer, bit, result)
			s7Param.Value = result
		}
		//setting.ZAPS.Debugf("通信S7客户端[%s]接口读变量数据块数据 %v", c.Name, s7Param.Value)
	} else {
		return nil, err
	}

	aDatas, err = json.Marshal(s7Param)
	if err != nil {
		setting.ZAPS.Errorf("通信S7客户端[%s]接口读变量数据块JSON格式化错误 %v", c.Name, err)
	}

	return aDatas, err
}

// ByteToBool 字节转bool数组（大端）
func ByteToBool(data byte) [8]byte {
	var res [8]byte
	for i := 0; i < 8; i++ {
		res[i] = data & 1
		data = data >> 1
	}
	return res
}

func (c *CommunicationS7Template) GetName() string {
	return c.Name
}

func (c *CommunicationS7Template) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationS7Template) GetInterval() string {
	return "0"
}

func (c *CommunicationS7Template) GetType() int {
	return CommTypeS7
}

func ReadCommS7InterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commS7Interface.json")
	if err != nil {
		setting.ZAPS.Debugf("打开通信接口[S7客户端]通信接口配置json文件失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &CommunicationS7Map)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[S7客户端]通信接口配置json文件格式化失败 %v", err)
		return false
	}
	setting.ZAPS.Debugf("打开通信接口[S7客户端]通信接口配置json文件成功")
	return true
}

func WriteCommS7InterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationS7Map)
	err := utils.FileWrite("./selfpara/commS7Interface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入通信接口[S7客户端]通信接口配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入通信接口[S7客户端]通信接口配置json文件 %s", "成功")
}
