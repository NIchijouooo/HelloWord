package commInterface

import (
	"encoding/json"
	"errors"
	"gateway/setting"
	"gateway/utils"
	"net"
	"strconv"
	"time"
)

type STLInterfaceParam struct {
	Port    string `json:"port"`
	Timeout string `json:"timeout"` //通信超时
}

type STLGWTemplate struct {
	IP   net.IP
	Port int
	SN   string
}

type CommunicationSTLTemplate struct {
	Name   string                   `json:"name"`   //接口名称
	Type   string                   `json:"type"`   //接口类型,比如serial,STL,udp,http
	Param  STLInterfaceParam        `json:"param"`  //接口参数
	Conn   net.Conn                 `json:"-"`      //通信句柄
	Status ConnectStatus            `json:"status"` //连接状态
	GW     map[string]STLGWTemplate `json:"-"`
}

var CommunicationSTLMap = make([]*CommunicationSTLTemplate, 0)

func (c *CommunicationSTLTemplate) Open() bool {
	port, _ := strconv.Atoi(c.Param.Port)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	})
	if err != nil {
		setting.ZAPS.Errorf("通信接口UDP服务端[%s]接口打开失败 %v", c.Name, err)
		c.Conn = nil
		c.Status = CommIsUnConnect
		return false
	}
	setting.ZAPS.Infof("通信接口UDP客户端[%s]接口打开成功", c.Name)
	c.Status = CommIsConnect
	c.Conn = conn

	go c.Accept(conn)

	return true
}

func (c *CommunicationSTLTemplate) Accept(conn *net.UDPConn) bool {

	for {
		rData := make([]byte, 2048)
		rCnt, rAddr, err := conn.ReadFromUDP(rData) //接受UDP数据
		if err != nil {
			setting.ZAPS.Errorf("通信接口UDP服务端[%s]接口读取数据失败 %v", c.Name, err)
			continue
		}
		setting.ZAPS.Debugf("rAddr %v,rData[%d] %X", rAddr, rCnt, rData[:rCnt])
		c.ProcessReceiveData(rAddr, rData[:rCnt])
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *CommunicationSTLTemplate) ProcessReceiveData(addr *net.UDPAddr, rData []byte) {
	var fVer byte  //协议版本
	var fType byte //报文类型
	var dLen int   //数据长度
	var dCRC int   //数据CRC
	var fCRC int   //报文CRC

	if len(rData) < 14 { //报文头长度固定14字节
		return
	}
	for k, v := range rData {
		if v != 0xC3 {
			continue
		}
		sn := string(rData[k+1 : k+5])
		fVer = rData[k+6]
		fVer = fVer
		fType = rData[k+7]
		dLen = int(rData[k+8])*256 + int(rData[k+9])
		dLen = dLen
		dCRC = int(rData[k+10])*256 + int(rData[k+11])
		dCRC = dCRC
		fCRC = int(rData[k+12])*256 + int(rData[k+13])
		tmpCRC := int(setting.CRC16(rData[k : k+12]))
		if tmpCRC != fCRC {
			setting.ZAPS.Debugf("CRC错误 fCRC %x,tmpCRC %x", fCRC, tmpCRC)
			continue
		}
		switch fType {
		case 0x01:
			{
				setting.ZAPS.Debugf("通信接口STL接收到登陆消息")
			}
		}
		gw := STLGWTemplate{
			SN:   sn,
			IP:   addr.IP,
			Port: addr.Port,
		}
		c.GW[sn] = gw
		setting.ZAPS.Debugf("通信接口STL中网关map[%+v]", c.GW)
	}
}

func (c *CommunicationSTLTemplate) Close() bool {
	if c.Conn != nil {
		err := c.Conn.Close()
		if err != nil {
			setting.ZAPS.Errorf("通信接口STL[%s]接口关闭失败 %v", c.Name, err)
			return false
		}
		setting.ZAPS.Infof("通信接口STL[%s]接口关闭成功", c.Name)
		c.Status = CommIsUnConnect
	}
	return true
}

func (c *CommunicationSTLTemplate) WriteData(data []byte) ([]byte, error) {

	if c.Conn != nil {
		_, err := c.Conn.Write(data)
		if err != nil {
			setting.ZAPS.Errorf("通信接口STL[%s]接口写失败 %v", c.Name, err)
			c.Close()
			c.Open()
			return nil, nil
		}
		return nil, nil
	} else {
		c.Open()
	}
	return nil, nil
}

func (c *CommunicationSTLTemplate) ReadData(data []byte) ([]byte, error) {

	if c.Conn != nil {
		cnt, err := c.Conn.Read(data)
		//setting.ZAPS.Debugf("%s,STL read data cnt %v", c.Name, cnt)
		if err != nil {
			//setting.ZAPS.Errorf("%s,STL read err,%v", c.Name, err)
			return nil, err
		}
		return data[:cnt], nil
	}
	return nil, errors.New("socket未打开")
}

func (c *CommunicationSTLTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationSTLTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationSTLTemplate) GetInterval() string {
	return ""
}

func (c *CommunicationSTLTemplate) GetType() int {
	return CommTypeSTL
}

func ReadCommSTLInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commSTLInterface.json")
	if err != nil {
		setting.ZAPS.Debugf("打开通信接口[Tcp客户端]通信接口配置json文件失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &CommunicationSTLMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[Tcp客户端]通信接口配置json文件格式化失败 %v", err)
		return false
	}
	for _, v := range CommunicationSTLMap {
		v.GW = make(map[string]STLGWTemplate)
	}
	setting.ZAPS.Debugf("打开通信接口[Tcp客户端]通信接口配置json文件成功")
	return true
}

func WriteCommSTLInterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationSTLMap)
	err := utils.FileWrite("./selfpara/commSTLInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入通信接口[Tcp客户端]通信接口配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入通信接口[Tcp客户端]通信接口配置json文件 %s", "成功")
}
