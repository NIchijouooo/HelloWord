package commInterface

import (
	"encoding/json"
	"fmt"
	"gateway/setting"
	"gateway/utils"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type DR504InterfaceParam struct {
	Port    string `json:"port"`
	Timeout string `json:"timeout"` //通信超时
}

type DR504NodeTemplate struct {
	CommName    string
	SN          string
	Conn        net.Conn `json:"-"`
	ReceiveData []byte   `json:"-"`
}

type CommunicationDR504Template struct {
	Name     string                        `json:"name"`   //接口名称
	Type     string                        `json:"type"`   //接口类型,比如serial,TcpClient,DR504,udp,http
	Param    DR504InterfaceParam           `json:"param"`  //接口参数
	Status   ConnectStatus                 `json:"status"` //连接状态
	Nodes    map[string]*DR504NodeTemplate `json:"-"`      //属性
	NodeAddr string                        `json:"-"`      //设备地址
}

var CommunicationDR504Map = make([]*CommunicationDR504Template, 0)
var CommunicationDR504MapLock sync.RWMutex

func NewDR504Node(commName string, conn net.Conn) *DR504NodeTemplate {
	return &DR504NodeTemplate{
		CommName:    commName,
		Conn:        conn,
		ReceiveData: make([]byte, 0),
	}
}

func (c *CommunicationDR504Template) Open() bool {
	listener, err := net.Listen("tcp", "0.0.0.0"+":"+c.Param.Port)
	if err != nil {
		setting.ZAPS.Errorf("通信DR504[%s]打开失败 %v", c.Name, err)
		c.Status = CommIsUnConnect
		return false
	}
	setting.ZAPS.Infof("通信DR504[%s]打开成功", c.Name)
	c.Status = CommIsConnect

	go c.Accept(listener)

	return true
}

func (c *CommunicationDR504Template) Accept(listener net.Listener) bool {

	for {
		conn, err := listener.Accept()
		if err != nil {
			setting.ZAPS.Errorf("通信DR504[%s]监听失败 %v", c.Name, err)
			continue
		}
		node := NewDR504Node(c.Name, conn)
		go node.ReadFromClient()
	}
}

func (s *DR504NodeTemplate) ReadFromClient() {

	for {
		rData := make([]byte, 2048)
		rCnt, err := s.Conn.Read(rData)
		if err != nil && err != io.EOF {
			setting.ZAPS.Errorf("通信DR504[%s]读取数据失败 %v", s.CommName, err)
			return
		}
		if rCnt > 0 {
			setting.ZAPS.Debugf("源IP[%v]", s.Conn.RemoteAddr())
			setting.ZAPS.Debugf("通信DR504[%s]接收数据[%d]%02X", s.CommName, rCnt, rData[:rCnt])
			s.AnalyseReceiveData(rData[:rCnt])
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (s *DR504NodeTemplate) AnalyseReceiveData(rData []byte) {

	rCnt := len(rData)
	for i := 0; i < rCnt; i++ {
		//注册包
		if (i + 1) >= rCnt {
			return
		}
		if rData[i] == 0x55 && rData[i+1] == 0xAA {
			if rCnt < 10 {
				return
			}
			sn := ""
			for index := 0; index < 10; index++ {
				sn = sn + fmt.Sprintf("%02X", rData[i+index])
			}
			CommunicationDR504MapLock.Lock()
			for _, v := range CommunicationDR504Map {
				if v.Name == s.CommName {
					s.SN = sn
					v.Nodes[sn] = s
					setting.ZAPS.Debugf("通信DR504设备数量[%d]，列表 %+v", len(v.Nodes), v.Nodes)
				}
			}
			CommunicationDR504MapLock.Unlock()
		} else {
			s.ReceiveData = append(s.ReceiveData, rData...)
			setting.ZAPS.Debugf("通信DR504[%s]设备[%s]接收数据[%d]%02X", s.CommName, s.SN, len(s.ReceiveData), s.ReceiveData)
			return
		}
	}
}

func (c *CommunicationDR504Template) Close() bool {

	return true
}

func (c *CommunicationDR504Template) WriteData(data []byte) ([]byte, error) {

	cmd := struct {
		SN   string
		Data []byte
	}{}

	err := json.Unmarshal(data, &cmd)
	if err != nil {
		setting.ZAPS.Debugf("通信接口DR504[%s]写数据JSON格式化错误 %v", c.Name, err)
		return nil, err
	}
	node, ok := c.Nodes[cmd.SN]
	if !ok {
		setting.ZAPS.Debugf("通信接口DR504[%s]写数据地址不存在", c.Name)
		return nil, err
	}
	c.NodeAddr = cmd.SN
	_, err = node.Conn.Write(cmd.Data)
	if err != nil {
		setting.ZAPS.Debugf("通信接口DR504[%s]写数据错误 %v", c.Name, err)
		return nil, err
	}

	return nil, nil
}

func (c *CommunicationDR504Template) ReadData(data []byte) ([]byte, error) {

	node, ok := c.Nodes[c.NodeAddr]
	if !ok {
		//setting.ZAPS.Debugf("通信接口DR504[%s]写数据地址不存在", c.Name)
		return nil, nil
	}

	rData := node.ReceiveData
	node.ReceiveData = node.ReceiveData[0:0]
	return rData, nil
}

func (c *CommunicationDR504Template) GetName() string {
	return c.Name
}

func (c *CommunicationDR504Template) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationDR504Template) GetInterval() string {
	return "0"
}

func (c *CommunicationDR504Template) GetType() int {
	return CommTypeDR504
}

func ReadCommDR504InterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commDR504Interface.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such directory") {
			setting.ZAPS.Debug("打开通信接口[DR504]通信接口配置json文件失败")
		} else {
			setting.ZAPS.Debugf("打开通信接口[DR504]通信接口配置json文件失败 %v", err)
		}
		return false
	}
	err = json.Unmarshal(data, &CommunicationDR504Map)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[DR504]通信接口配置json文件格式化失败")
		return false
	}

	for _, v := range CommunicationDR504Map {
		v.Nodes = make(map[string]*DR504NodeTemplate)
	}
	setting.ZAPS.Debugf("打开通信接口[DR504]通信接口配置json文件成功")
	return true
}

func WriteCommDR504InterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationDR504Map)
	err := utils.FileWrite("./selfpara/commDR504Interface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入DR504通信接口配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入DR504通信接口配置json文件 %s", "成功")
}
