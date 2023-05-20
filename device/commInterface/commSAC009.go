package commInterface

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/setting"
	"gateway/utils"
	"io"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SAC009InterfaceParam struct {
	Port    string `json:"port"`
	Timeout string `json:"timeout"` //通信超时
}

type SAC009PropertyTemplate struct {
	Index     int         `json:"index"`
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`   //变量值，不可以是字符串
	Explain   interface{} `json:"explain"` //变量值解释，必须是字符串
	TimeStamp time.Time   `json:"timeStamp"`
}

type SAC009NodeTemplate struct {
	CommName    string
	Conn        net.Conn `json:"-"`
	IMEI        string
	SN          string
	Properties  map[string]SAC009PropertyTemplate
	ReceiveChan chan ReceiveFrameTemplate `json:"-"`
}

type ReceiveFrameTemplate struct {
	FunCode int    `json:"funCode"`
	Frame   []byte `json:"frame"`
}

type CommunicationSAC009Template struct {
	Name   string                         `json:"name"`   //接口名称
	Type   string                         `json:"type"`   //接口类型,比如serial,TcpClient,SAC009,udp,http
	Param  SAC009InterfaceParam           `json:"param"`  //接口参数
	Status ConnectStatus                  `json:"status"` //连接状态
	Nodes  map[string]*SAC009NodeTemplate `json:"-"`      //属性
}

var CommunicationSAC009Map = make([]*CommunicationSAC009Template, 0)
var CommunicationSAC009MapLock sync.RWMutex

func NewSAC009Node(commName string, conn net.Conn) *SAC009NodeTemplate {
	return &SAC009NodeTemplate{
		CommName:    commName,
		Conn:        conn,
		IMEI:        "",
		Properties:  make(map[string]SAC009PropertyTemplate),
		ReceiveChan: make(chan ReceiveFrameTemplate, 5),
	}
}

func (c *CommunicationSAC009Template) Open() bool {
	listener, err := net.Listen("tcp", "0.0.0.0"+":"+c.Param.Port)
	if err != nil {
		setting.ZAPS.Errorf("通信SAC009[%s]打开失败 %v", c.Name, err)
		c.Status = CommIsUnConnect
		return false
	}
	setting.ZAPS.Infof("通信SAC009[%s]打开成功", c.Name)
	c.Status = CommIsConnect

	go c.Accept(listener)

	return true
}

func (c *CommunicationSAC009Template) Accept(listener net.Listener) bool {
	defer func() {
		r := recover()
		if r != nil {
			//loger.Printf("主程序发生错误 %s", debug.Stack())
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			utils.ErrLog.Println("<==========================>")
			utils.ErrLog.Printf("%s\n", string(buf[:n]))
			//os.Exit(1)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			setting.ZAPS.Errorf("通信SAC009[%s]监听失败 %v", c.Name, err)
			time.Sleep(1 * time.Second)
			continue
		}
		node := NewSAC009Node(c.Name, conn)

		go node.ReadFromClient(c.Name)
	}
}

func (s *SAC009NodeTemplate) ReadFromClient(commName string) {

	defer func() {
		r := recover()
		if r != nil {
			//loger.Printf("主程序发生错误 %s", debug.Stack())
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			utils.ErrLog.Println("<==========================>")
			utils.ErrLog.Printf("%s\n", string(buf[:n]))
		}
	}()
	for {
		rData := make([]byte, 2048)
		if s.Conn == nil {
			return
		}
		rCnt, err := s.Conn.Read(rData)
		if err != nil && err != io.EOF {
			setting.ZAPS.Errorf("通信SAC009[%s]读取设备[%s]数据失败 %v", s.CommName, s.IMEI, err)
			if strings.Contains(err.Error(), "reset by peer") || strings.Contains(err.Error(), "broken pipe") {
				setting.ZAPS.Infof("通信SAC009[%s]设备[%s]读数据协程退出", s.CommName, s.IMEI)
				_ = s.Conn.Close()
				return
			}
		}
		if rCnt > 0 {
			//setting.ZAPS.Debugf("通信SAC009[%s]设备[%s]接收数据[%d]%02X", s.CommName, s.IMEI, rCnt, rData[:rCnt])
			s.AnalyseReceiveData(commName, rData[:rCnt])
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (s *SAC009NodeTemplate) AnalyseReceiveData(commName string, rData []byte) {
	for i := 0; i < len(rData); i++ {
		if rData[i] != 0x01 {
			continue
		}
		if (i + 1) >= len(rData) {
			return
		}
		funCode := rData[i+1]
		switch funCode {
		case 0x50: //主动上报
			{
				if len(rData) < (42 + 5) {
					continue
				}
				if (i + 41 + 5) >= len(rData) {
					return
				}
				dataLen := rData[i+41+5]
				totalLen := int(42 + 5 + dataLen + 2)
				//setting.ZAPS.Debugf("dataLen %v,totalLen %v", dataLen, totalLen)
				if len(rData) < (totalLen + 2) {
					continue
				}

				rIMEI := ""
				for index := 2; index < 17; index++ {
					rIMEI = rIMEI + fmt.Sprintf("%X", rData[index])
				}
				//setting.ZAPS.Debugf("IMEI %v", rIMEI)
				if s.IMEI == "" {
					s.IMEI = rIMEI
					CommunicationSAC009MapLock.Lock()
					for k, v := range CommunicationSAC009Map {
						if v.Name == s.CommName {
							CommunicationSAC009Map[k].Nodes[rIMEI] = s
							setting.ZAPS.Debugf("通信接口SAC009[%s]设备数量[%d]，列表 %+v", commName, len(CommunicationSAC009Map[k].Nodes), CommunicationSAC009Map[k].Nodes)
						}
					}
					CommunicationSAC009MapLock.Unlock()
					setting.ZAPS.Infof("通信接口SAC009[%s]设备[%s]上线成功", commName, rIMEI)
				}
				if (i + 47) > len(rData) {
					return
				}
				_ = s.AnalyzeReadData(rData[i+47:])
				return
			}
		case 0x03: //读数据返回
			{
				if (i + 2) > len(rData) {
					return
				}
				dataLen := int(rData[i+2])
				if (i + 3 + dataLen) > len(rData) {
					return
				}
				frame := ReceiveFrameTemplate{
					FunCode: 0x03,
					Frame:   rData[i+3 : i+3+dataLen],
				}
				s.ReceiveChan <- frame
				return
			}
		case 0x06: //写数据返回
			{
				if (i + 3) > len(rData) {
					return
				}
				frame := ReceiveFrameTemplate{
					FunCode: 0x06,
					Frame:   rData[i+3:],
				}
				s.ReceiveChan <- frame
				return
			}
		case 0x10: //写数据返回
			{
				if (i + 3) > len(rData) {
					return
				}
				frame := ReceiveFrameTemplate{
					FunCode: 0x10,
					Frame:   rData[i+3:],
				}
				s.ReceiveChan <- frame
				return
			}
		}
	}
}

func (c *CommunicationSAC009Template) Close() bool {

	return true
}

func (c *CommunicationSAC009Template) WriteData(data []byte) ([]byte, error) {

	wNode := SAC009NodeTemplate{
		Properties: make(map[string]SAC009PropertyTemplate),
	}
	err := json.Unmarshal(data, &wNode)
	if err != nil {
		return nil, errors.New("JSON格式化错误")
	}

	node, ok := c.Nodes[wNode.IMEI]
	if !ok {
		return nil, errors.New("设备IMEI不存在")
	}

	startStopCmd := -1
	tempSet := -1
	modeSet := -1
	airBrandID := -1
	for _, v := range wNode.Properties {
		if v.Name == "StartStopCmd" {
			switch v.Value.(type) {
			case float64:
				startStopCmd = int(v.Value.(float64))
			case string:
				startStopCmd, _ = strconv.Atoi(v.Value.(string))
			default:
				return nil, errors.New("数据类型错误")
			}
		} else if v.Name == "TempSet" {
			switch v.Value.(type) {
			case float64:
				tempSet = int(v.Value.(float64)) * 10
			case string:
				temp, _ := strconv.ParseFloat(v.Value.(string), 16)
				tempSet = int(temp) * 10
			default:
				return nil, errors.New("数据类型错误")
			}
		} else if v.Name == "ModeSet" {
			switch v.Value.(type) {
			case float64:
				modeSet = int(v.Value.(float64))
			case string:
				modeSet, _ = strconv.Atoi(v.Value.(string))
			default:
				return nil, errors.New("数据类型错误")
			}
		} else if v.Name == "AirBrandID" {
			switch v.Value.(type) {
			case float64:
				airBrandID = int(v.Value.(float64))
			case string:
				airBrandID, _ = strconv.Atoi(v.Value.(string))
			default:
				return nil, errors.New("数据类型错误")
			}
		}
	}

	wData := make([]byte, 0)
	if startStopCmd != -1 {
		wData = append(wData, 0x01)
		wData = append(wData, 0x06)
		wData = append(wData, 0x00)
		wData = append(wData, 0x08)
		if startStopCmd == 1 {
			wData = append(wData, 0x00)
			wData = append(wData, 0x01)
		} else {
			wData = append(wData, 0x00)
			wData = append(wData, 0x02)
		}
	} else if tempSet != -1 && modeSet != -1 {
		wData = append(wData, 0x01)
		wData = append(wData, 0x10)
		wData = append(wData, 0x00)
		wData = append(wData, 0x09)
		wData = append(wData, 0x00)
		wData = append(wData, 0x03)
		wData = append(wData, 0x06)

		wData = append(wData, 0x00)
		wData = append(wData, byte(modeSet))
		wData = append(wData, byte(tempSet/256))
		wData = append(wData, byte(tempSet%256))
		wData = append(wData, 0x00)
		wData = append(wData, 0x01)
	} else if airBrandID != -1 {
		wData = append(wData, 0x01)
		wData = append(wData, 0x06)
		wData = append(wData, 0x00)
		wData = append(wData, 0x05)
		wData = append(wData, byte(airBrandID/256))
		wData = append(wData, byte(airBrandID%256))
	}

	crc16 := setting.CRC16(wData)
	wData = append(wData, byte(crc16%256))
	wData = append(wData, byte(crc16/256))

	setting.ZAPS.Debugf("通信接口SAC009[%s]设备[%s]写数据%X", c.Name, node.IMEI, wData)
	_, err = node.Conn.Write(wData)
	if err != nil {
		return nil, errors.New("写数据错误")
	}

	select {
	case recv := <-node.ReceiveChan:
		{
			if recv.FunCode == 0x06 || recv.FunCode == 0x10 {
				ackProperties := make(map[string]SAC009PropertyTemplate)
				for _, v := range wNode.Properties {
					property := SAC009PropertyTemplate{
						Name:      v.Name,
						Value:     0,
						TimeStamp: time.Now(),
					}
					ackProperties[v.Name] = property
				}

				//setting.ZAPS.Debugf("ackProperty %v", ackProperties)
				data, _ := json.Marshal(ackProperties)
				return data, nil
			} else {
				return nil, errors.New("等待数据超时")
			}
		}
	case <-time.After(time.Millisecond * 2000):
		{
			return nil, errors.New("等待数据超时")
		}
	}
}

func InsertUint16(index int, name string, data []byte) SAC009PropertyTemplate {
	property := SAC009PropertyTemplate{}

	property.Index = index
	property.Name = name
	property.Value = binary.BigEndian.Uint16(data)
	property.TimeStamp = time.Now()

	return property
}

func InsertUint32(index int, name string, data []byte) SAC009PropertyTemplate {
	property := SAC009PropertyTemplate{}

	property.Index = index
	property.Name = name
	property.Value = binary.BigEndian.Uint32(data)
	property.TimeStamp = time.Now()

	return property
}

func (s *SAC009NodeTemplate) AnalyzeReadData(data []byte) error {

	//setting.ZAPS.Debugf("dataLen %v", len(data))
	if len(data) < 76 {
		return errors.New("数据长度不够")
	}

	dataIndex := 0

	s.Properties["CSQ"] = InsertUint16(dataIndex, "CSQ", data[dataIndex:])
	dataIndex += 2

	s.Properties["RoomHumi"] = InsertUint16(dataIndex, "RoomHumi", data[dataIndex:])
	dataIndex += 2

	dataIndex += 2

	s.Properties["RoomTemp"] = InsertUint16(dataIndex, "RoomTemp", data[dataIndex:])
	dataIndex += 2

	dataIndex += 2 * 3

	nameArray := []string{
		"StartStopStatus", //开关状态
		"ModeSet",         //设定模式
		"TempSet",         //设定温度
		"AirSet",          //设定风速
	}

	for i := 0; i < 4; i++ {
		s.Properties[nameArray[i]] = InsertUint16(dataIndex, nameArray[i], data[dataIndex:])
		dataIndex += 2
	}

	dataIndex += 2 * 15

	s.Properties["TotalE"] = InsertUint32(dataIndex, "TotalE", data[dataIndex:])
	dataIndex += 4

	nameArray = []string{
		"Voltage", //电压
		"Current", //电流
		"Power",   //功率
	}
	for i := 0; i < 3; i++ {
		s.Properties[nameArray[i]] = InsertUint16(dataIndex, nameArray[i], data[dataIndex:])
		dataIndex += 2
	}

	//setting.ZAPS.Debugf("nodeIMEI %s", s.IMEI)
	//if strings.Contains(s.IMEI, "040307") || strings.Contains(s.IMEI, "080309") {
	//	setting.ZAPS.Debugf("StartStopStatus %v", s.Properties["StartStopStatus"])
	//	setting.ZAPS.Debugf("ModeSet %v", s.Properties["ModeSet"])
	//setting.ZAPS.Debugf("TempSet %v", s.Properties["TempSet"])
	//	setting.ZAPS.Debugf("voltage %v", s.Properties["Voltage"])
	//	setting.ZAPS.Debugf("Current %v", s.Properties["Current"])
	//}

	return nil
}

func (c *CommunicationSAC009Template) ReadData(data []byte) ([]byte, error) {

	rNode := SAC009NodeTemplate{
		Properties: make(map[string]SAC009PropertyTemplate),
	}
	err := json.Unmarshal(data, &rNode)
	if err != nil {
		return nil, errors.New("JSON格式化错误")
	}

	setting.ZAPS.Debugf("rNode %+v", rNode)
	node, ok := c.Nodes[rNode.IMEI]
	if !ok {
		return nil, errors.New("设备不存在")
	}

	wData := make([]byte, 0)
	wData = append(wData, 0x01)
	wData = append(wData, 0x03)
	wData = append(wData, 0x00)
	wData = append(wData, 0x05)
	wData = append(wData, 0x00)
	wData = append(wData, 0x01)
	crc16 := setting.CRC16(wData)
	wData = append(wData, byte(crc16%256))
	wData = append(wData, byte(crc16/256))

	_, err = node.Conn.Write(wData)
	if err != nil {
		return nil, errors.New("写数据错误")
	}

	select {
	case recv := <-node.ReceiveChan:
		{
			if len(recv.Frame) < 2 {
				return nil, errors.New("接收解析错误")
			}
			node.Properties["AirBrandID"] = InsertUint16(0, "AirBrandID", recv.Frame[0:])
		}
	case <-time.After(time.Millisecond * 2000):
		{
			return nil, errors.New("等待数据超时")
		}
	}

	setting.ZAPS.Debugf("设备属性 %v", node.Properties)
	rData, _ := json.Marshal(node.Properties)
	return rData, nil
}

func (c *CommunicationSAC009Template) GetName() string {
	return c.Name
}

func (c *CommunicationSAC009Template) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationSAC009Template) GetInterval() string {
	return "0"
}

func (c *CommunicationSAC009Template) GetType() int {
	return CommTypeSAC009
}

func ReadCommSAC009InterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commSAC009Interface.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such directory") {
			setting.ZAPS.Debug("打开通信接口[SAC009]通信接口配置json文件失败")
		} else {
			setting.ZAPS.Debugf("打开通信接口[SAC009]通信接口配置json文件失败 %v", err)
		}
		return false
	}
	err = json.Unmarshal(data, &CommunicationSAC009Map)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[SAC009]通信接口配置json文件格式化失败")
		return false
	}

	for _, v := range CommunicationSAC009Map {
		v.Nodes = make(map[string]*SAC009NodeTemplate)
	}
	setting.ZAPS.Debugf("打开通信接口[SAC009]通信接口配置json文件成功")
	return true
}

func WriteCommSAC009InterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationSAC009Map)
	err := utils.FileWrite("./selfpara/commSAC009Interface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入SAC009通信接口配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入SAC009通信接口配置json文件 %s", "成功")
}
