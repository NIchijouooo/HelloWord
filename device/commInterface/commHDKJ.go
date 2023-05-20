package commInterface

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/setting"
	"gateway/utils"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type HDKJInterfaceParam struct {
	Port           string `json:"port"`
	ReportInterval string `json:"reportInterval"`
	OfflinePeriod  string `json:"offlinePeriod"`
}

type HDKJPropertyTemplate struct {
	Index     int         `json:"index"`
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`   //变量值，不可以是字符串
	Explain   interface{} `json:"explain"` //变量值解释，必须是字符串
	TimeStamp time.Time   `json:"timeStamp"`
}
type HDKJTypesTemplate struct {
	ChannelNum byte `json:"channelNum"`
	LastFlow   uint `json:"lastFlow"`
	MultiData  uint `json:"multiData"`
}
type HDKJNodeTemplate struct {
	CommName string
	//Conn            *net.UDPConn `json:"-"`
	IMEI            string
	StartTime       string
	DataNumber      byte
	DataInterval    byte
	BatteryVoltage  uint
	ReportTimes     uint
	SignalIntensity byte
	DeviceInfo      string
	ReservedBytes   []byte
	DataTypes       []HDKJTypesTemplate
	Properties      map[string]HDKJPropertyTemplate
	TimeStamp       time.Time //成功收到数据时间
}

type CommunicationHDKJTemplate struct {
	Name   string                       `json:"name"` //接口名称
	Conn   *net.UDPConn                 `json:"-"`
	Type   string                       `json:"type"`   //接口类型,比如serial,TcpClient,HDKJ,udp,http
	Param  HDKJInterfaceParam           `json:"param"`  //接口参数
	Status ConnectStatus                `json:"status"` //连接状态
	Nodes  map[string]*HDKJNodeTemplate `json:"-"`      //属性
}

var CommunicationHDKJMap = make([]*CommunicationHDKJTemplate, 0)
var CommunicationHDKJMapLock sync.RWMutex

func (c *CommunicationHDKJTemplate) Open() bool {

	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0"+":"+c.Param.Port)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		setting.ZAPS.Errorf("通信HDKJ[%s]打开失败 %v", c.Name, err)
		c.Status = CommIsUnConnect
		return false
	}
	setting.ZAPS.Infof("通信HDKJ[%s]打开成功", c.Name)
	c.Conn = conn
	c.Status = CommIsConnect

	//defer conn.Close()

	//node := &HDKJNodeTemplate{
	//	CommName:   c.Name,
	//	Conn:       conn,
	//	IMEI:       "",
	//	Properties: make(map[string]HDKJPropertyTemplate),
	//}
	go ReadFromClient(c.Conn, c.Name)
	return true
}

func (c *CommunicationHDKJTemplate) Close() bool {
	fmt.Println("Close")
	err := c.Conn.Close()
	if err != nil {
		setting.ZAPS.Errorf("通信HDKJ[%s]关闭失败 %v", c.Name, err)
		return false
	}
	return true
}

func ReadFromClient(conn *net.UDPConn, commName string) {

	for {
		s := &HDKJNodeTemplate{
			CommName:        commName,
			IMEI:            "",
			StartTime:       "",
			DataNumber:      0,
			DataInterval:    0,
			BatteryVoltage:  0,
			ReportTimes:     0,
			SignalIntensity: 0,
			DeviceInfo:      "",
			ReservedBytes:   nil,
			DataTypes:       nil,
			Properties:      make(map[string]HDKJPropertyTemplate),
			TimeStamp:       time.Time{},
		}

		rData := make([]byte, 4096)
		rCnt, remoteAddr, err := conn.ReadFromUDP(rData)

		fmt.Printf("ReadFromClient <%s> -> rCnt:%d, err[%v]\n", remoteAddr, rCnt, err)
		if err != nil && err != io.EOF {
			setting.ZAPS.Errorf("通信HDKJ[%s]读取数据失败 err[%v]", s.CommName, err)
			break
		}
		if rCnt > 0 {
			setting.ZAPS.Debugf("通信HDKJ[%s]接收数据[%d]", s.CommName, rCnt)
			conn.WriteToUDP([]byte{0x40, 0x32, 0x32, 0x32}, remoteAddr)
			s.AnalyseReceiveData(rData[:rCnt])
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (s *HDKJNodeTemplate) AnalyseReceiveData(rData []byte) {

	for i := 0; i < len(rData); i++ {
		if (i + 3) >= len(rData) {
			return
		}
		if rData[i] != 0xAB || rData[i+1] != 0xCD {
			continue
		}
		length := i + 3 + int(rData[i+2])*256 + int(rData[i+3])
		if rData[len(rData)-3] != GetCRC(rData[i:], length) {
			setting.ZAPS.Debugf("通信HDKJ CRC错误")
			continue
		}
		if int(rData[i+2])*256+int(rData[i+3])+2+2+2 > len(rData) {
			setting.ZAPS.Debugf("通信HDKJ 数据帧长度错误")
			return
		}
		s.Properties = make(map[string]HDKJPropertyTemplate)
		s.IMEI = strconv.Itoa(int(rData[i+4])*256 + int(rData[i+5])) //不带0
		s.StartTime = fmt.Sprintf("%02X", rData[i+6:i+11])
		s.DataNumber = rData[i+11]
		s.DataInterval = rData[i+12]
		s.BatteryVoltage = uint(rData[i+13])*256 + uint(rData[i+14])
		s.ReportTimes = uint(rData[i+15])*256 + uint(rData[i+16])
		s.SignalIntensity = rData[i+17]
		s.DeviceInfo = strings.TrimSpace(string(rData[i+18 : i+38]))
		s.ReservedBytes = rData[i+38 : i+46]
		setting.ZAPS.Debugf("站号[%v], 起始时间[%v], 数据个数[%v], 数据间隔[%v], 电池电压[%v], 上发次数[%v], 信号强度[%v], 设备信息[%v], 保留字节[%v]", s.IMEI, s.StartTime, s.DataNumber, s.DataInterval, s.BatteryVoltage, s.ReportTimes, s.SignalIntensity, s.DeviceInfo, s.ReservedBytes)

		s.Properties["BatteryVoltage"] = HDKJPropertyTemplate{
			Index:     len(s.Properties),
			Name:      "BatteryVoltage",
			Value:     s.BatteryVoltage,
			Explain:   fmt.Sprintf("%v", s.BatteryVoltage),
			TimeStamp: time.Now(),
		}
		s.Properties["SignalIntensity"] = HDKJPropertyTemplate{
			Index:     len(s.Properties),
			Name:      "SignalIntensity",
			Value:     s.SignalIntensity,
			Explain:   fmt.Sprintf("%v", s.SignalIntensity),
			TimeStamp: time.Now(),
		}
		err := s.AnalyzeReadData(rData[i+46:])
		if err != nil {
			setting.ZAPS.Errorf("AnalyzeReadData Error[%v]", err)
			return
		}
		//setting.ZAPS.Debugf("s.Properties-> %v", s.Properties)
		CommunicationHDKJMapLock.Lock()
		for k, v := range CommunicationHDKJMap {
			if v.Name == s.CommName {
				s.TimeStamp = time.Now()
				CommunicationHDKJMap[k].Nodes[s.IMEI] = s
				setting.ZAPS.Debugf("CommunicationHDKJMap[%v]", CommunicationHDKJMap[k].Nodes)
				setting.ZAPS.Debugf("通信SAC009设备数量[%d]，列表 %+v", len(CommunicationHDKJMap[k].Nodes), CommunicationHDKJMap[k].Nodes)
			}
		}
		CommunicationHDKJMapLock.Unlock()
	}
}

func (s *HDKJNodeTemplate) AnalyzeReadData(data []byte) error {

	var Start, countType = 0, 0
	var tempNode = HDKJTypesTemplate{}
	//s.Properties = make(map[string]HDKJPropertyTemplate)
	for Start+2 < len(data) && data[Start+1] != 0x0D && data[Start+2] != 0x0A {
		var Type byte
		countType++
		if data[Start] >= 1 && data[Start] <= 6 {
			Type = data[Start]
			tempNode.ChannelNum = data[Start+1]
			//setting.ZAPS.Debugf("类型计数[%v], 类型[%v], 通道号[%v]", countType, Type, tempNode.ChannelNum)
		}
		dataNumber := int(s.DataNumber)
		switch Type {
		case 1:
			tempNode.LastFlow = GetValue(0, data[2:6], true)

			for i := 0; i < dataNumber; i++ {
				tempNode.MultiData = GetValue(0, data[Start+i*4+6:Start+i*4+10], true)
				s.DataTypes = append(s.DataTypes, HDKJTypesTemplate{
					ChannelNum: tempNode.ChannelNum,
					LastFlow:   tempNode.LastFlow,
					MultiData:  tempNode.MultiData,
				})
				switch tempNode.ChannelNum {
				case 1:
					//setting.ZAPS.Debugf("PositiveFlux 昨日流量[%v]", tempNode.LastFlow)
					property := HDKJPropertyTemplate{
						Index:     len(s.Properties),
						Name:      "PositiveFlux",
						Value:     tempNode.MultiData, //[len(tempNode.MultiData)-1]
						Explain:   fmt.Sprintf("%v", tempNode.MultiData),
						TimeStamp: time.Now(),
					}
					s.Properties["PositiveFlux"] = property
				case 2:
					//setting.ZAPS.Debugf("ReverseFlux 昨日流量[%v]", tempNode.LastFlow)
					property := HDKJPropertyTemplate{
						Index:     len(s.Properties),
						Name:      "ReverseFlux",
						Value:     tempNode.MultiData, //[len(tempNode.MultiData)-1]
						Explain:   fmt.Sprintf("%v", tempNode.MultiData),
						TimeStamp: time.Now(),
					}
					s.Properties["ReverseFlux"] = property
				}
			}
			Start += 1 + 1 + 4 + dataNumber*4
		case 2:
			fallthrough
		case 3:
			for i := 0; i < dataNumber; i++ {
				tempNode.MultiData = uint(data[Start+i+2])
			}
			Start += 1 + 1 + dataNumber
		case 4:
			for i := 0; i < dataNumber; i++ {
				tempNode.MultiData = GetValue(0, data[Start+i*4+2:Start+i*4+6], true)
			}
			Start += 1 + 1 + dataNumber*4
		case 5:
			Start += 1 + 1 + dataNumber*2
		case 6:
			Start += 1 + 1 + dataNumber*4
		}
		tempNode = HDKJTypesTemplate{}
	}
	setting.ZAPS.Debugf("Properties:%v", s.Properties["PositiveFlux"])
	setting.ZAPS.Debugf("ReverseFlux:%v", s.Properties["ReverseFlux"])
	Start, countType = 0, 0
	return nil
}

func (c *CommunicationHDKJTemplate) ReadData(data []byte) ([]byte, error) {

	rNode := HDKJNodeTemplate{
		Properties: make(map[string]HDKJPropertyTemplate),
	}
	err := json.Unmarshal(data, &rNode)
	if err != nil {
		return nil, errors.New("JSON格式化错误")
	}

	setting.ZAPS.Debugf("rNode.IMEI %+v", rNode.IMEI)

	setting.ZAPS.Debugf("c.Nodes %+v", c.Nodes)
	node, ok := c.Nodes[rNode.IMEI]
	if !ok {
		return nil, errors.New("设备不存在")
	}

	setting.ZAPS.Debugf("设备属性 %v", node.Properties)

	lastTime := node.TimeStamp
	latestTime := time.Now()
	interval, err := strconv.Atoi(c.Param.ReportInterval)
	offlinePeriod, err := strconv.Atoi(c.Param.OfflinePeriod)
	if err != nil {
		setting.ZAPS.Errorf("time error -> LastTime[%s], latestTime[%s], interval[%d], offlinePeriod[%d]", lastTime.Format("060102150405"), latestTime.Format("060102150405"), interval, offlinePeriod)
		return nil, errors.New("时间格式错误")
	}

	if latestTime.Sub(node.TimeStamp).Milliseconds() >= int64(interval*offlinePeriod*60000) {
		setting.ZAPS.Errorf("Offline timeout -> LastTime[%s], latestTime[%s], Duration[%f], interval[%d], offlinePeriod[%d]", lastTime.Format("060102150405"), latestTime.Format("060102150405"), latestTime.Sub(node.TimeStamp).Minutes(), interval, offlinePeriod)
		return nil, errors.New("设备离线")
	}
	rData, _ := json.Marshal(node.Properties)
	return rData, nil
}

func (c *CommunicationHDKJTemplate) WriteData(data []byte) ([]byte, error) {
	//
	//	wNode := HDKJNodeTemplate{
	//		Properties: make([]HDKJPropertyTemplate, 0),
	//	}
	//	err := json.Unmarshal(data, &wNode)
	//	if err != nil {
	//		return nil, errors.New("JSON格式化错误")
	//	}
	//
	//	node, ok := c.Nodes[wNode.IMEI]
	//	if !ok {
	//		return nil, errors.New("设备IMEI不存在")
	//	}
	//
	//	startStopCmd := -1
	//	tempSet := -1
	//	modeSet := -1
	//	airBrandID := -1
	//	for _, v := range wNode.Properties {
	//		if v.Name == "StartStopCmd" {
	//			switch v.Value.(type) {
	//			case float64:
	//				startStopCmd = int(v.Value.(float64))
	//			case string:
	//				startStopCmd, _ = strconv.Atoi(v.Value.(string))
	//			default:
	//				return nil, errors.New("数据类型错误")
	//			}
	//		} else if v.Name == "TempSet" {
	//			switch v.Value.(type) {
	//			case float64:
	//				tempSet = int(v.Value.(float64)) * 10
	//			case string:
	//				temp, _ := strconv.ParseFloat(v.Value.(string), 16)
	//				tempSet = int(temp) * 10
	//			default:
	//				return nil, errors.New("数据类型错误")
	//			}
	//		} else if v.Name == "ModeSet" {
	//			switch v.Value.(type) {
	//			case float64:
	//				modeSet = int(v.Value.(float64))
	//			case string:
	//				modeSet, _ = strconv.Atoi(v.Value.(string))
	//			default:
	//				return nil, errors.New("数据类型错误")
	//			}
	//		} else if v.Name == "AirBrandID" {
	//			switch v.Value.(type) {
	//			case float64:
	//				airBrandID = int(v.Value.(float64))
	//			case string:
	//				airBrandID, _ = strconv.Atoi(v.Value.(string))
	//			default:
	//				return nil, errors.New("数据类型错误")
	//			}
	//		}
	//	}
	//
	//	wData := make([]byte, 0)
	//	if startStopCmd != -1 {
	//		wData = append(wData, 0x01)
	//		wData = append(wData, 0x06)
	//		wData = append(wData, 0x00)
	//		wData = append(wData, 0x08)
	//		if startStopCmd == 1 {
	//			wData = append(wData, 0x00)
	//			wData = append(wData, 0x01)
	//		} else {
	//			wData = append(wData, 0x00)
	//			wData = append(wData, 0x02)
	//		}
	//	} else if tempSet != -1 && modeSet != -1 {
	//		wData = append(wData, 0x01)
	//		wData = append(wData, 0x10)
	//		wData = append(wData, 0x00)
	//		wData = append(wData, 0x09)
	//		wData = append(wData, 0x00)
	//		wData = append(wData, 0x03)
	//		wData = append(wData, 0x06)
	//
	//		wData = append(wData, 0x00)
	//		wData = append(wData, byte(modeSet))
	//		wData = append(wData, byte(tempSet/256))
	//		wData = append(wData, byte(tempSet%256))
	//		wData = append(wData, 0x00)
	//		wData = append(wData, 0x01)
	//	} else if airBrandID != -1 {
	//		wData = append(wData, 0x01)
	//		wData = append(wData, 0x06)
	//		wData = append(wData, 0x00)
	//		wData = append(wData, 0x05)
	//		wData = append(wData, byte(airBrandID/256))
	//		wData = append(wData, byte(airBrandID%256))
	//	}
	//
	//	crc16 := setting.CRC16(wData)
	//	wData = append(wData, byte(crc16%256))
	//	wData = append(wData, byte(crc16/256))
	//
	//	setting.ZAPS.Debugf("通信接口HDKJ[%s]设备[%s]写数据%X", c.Name, node.IMEI, wData)
	//	_, err = node.Conn.Write(wData)
	//	if err != nil {
	//		return nil, errors.New("写数据错误")
	//	}
	//
	//	select {
	//	case recv := <-node.ReceiveChan:
	//		{
	//			if recv.FunCode == 0x06 || recv.FunCode == 0x10 {
	//				ackProperties := make(map[string]HDKJPropertyTemplate)
	//				for _, v := range wNode.Properties {
	//					property := HDKJPropertyTemplate{
	//						Name:      v.Name,
	//						Value:     0,
	//						TimeStamp: time.Now(),
	//					}
	//					ackProperties[v.Name] = property
	//				}
	//
	//				//setting.ZAPS.Debugf("ackProperty %v", ackProperties)
	//				data, _ := json.Marshal(ackProperties)
	//				return data, nil
	//			} else {
	//				return nil, errors.New("等待数据超时")
	//			}
	//		}
	//	case <-time.After(time.Millisecond * 2000):
	//		{
	//			return nil, errors.New("等待数据超时")
	//		}
	//	}
	return nil, nil
}

func (c *CommunicationHDKJTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationHDKJTemplate) GetTimeOut() string {
	return "0"
}

func (c *CommunicationHDKJTemplate) GetInterval() string {
	return "0"
}

func (c *CommunicationHDKJTemplate) GetType() int {
	return CommTypeHDKJ
}

func ReadCommHDKJInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commHDKJInterface.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such directory") {
			setting.ZAPS.Debug("打开通信接口[HDKJ]通信接口配置json文件失败")
		} else {
			setting.ZAPS.Debugf("打开通信接口[HDKJ]通信接口配置json文件失败 %v", err)
		}
		return false
	}
	err = json.Unmarshal(data, &CommunicationHDKJMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[HDKJ]通信接口配置json文件格式化失败")
		return false
	}

	for _, v := range CommunicationHDKJMap {
		v.Nodes = make(map[string]*HDKJNodeTemplate)
	}
	setting.ZAPS.Debugf("打开通信接口[HDKJ]通信接口配置json文件成功")
	return true
}

func WriteCommHDKJInterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationHDKJMap)
	err := utils.FileWrite("./selfpara/commHDKJInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入HDKJ通信接口配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入HDKJ通信接口配置json文件 %s", "成功")
}

// GetValue reverse:true 低字节在前
func GetValue(bytes int, data []byte, reverse bool) uint {
	var result uint
	if reverse {
		for i := 0; i < len(data); i++ {
			result += uint(data[i]) * uint(math.Pow(256, float64(i)))
		}
	} else {
		for i := 0; i < len(data); i++ {
			result += uint(data[i]) * uint(math.Pow(256, float64(len(data)-1-i)))
		}
	}
	return result
}

func GetCRC(data []byte, length int) byte {
	var result byte
	for i := 0; i < length; i++ {
		result += data[i]
	}
	return result
}

func HexToFloat(data []byte) float32 {
	var value float32
	bits := binary.BigEndian.Uint32(data)
	value = math.Float32frombits(bits)
	return value
}
