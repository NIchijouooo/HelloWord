package commInterface

import (
	"encoding/json"
	"errors"
	"net"
	"strings"
	"gateway/setting"
	"gateway/utils"
)

type TcpServerInterfaceParam struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Timeout  string `json:"timeout"`  //通信超时
	Interval string `json:"interval"` //通信间隔
}

type CommunicationTcpServerTemplate struct {
	Name   string                  `json:"name"`   //接口名称
	Type   string                  `json:"type"`   //接口类型,比如serial,TcpClient,TcpServer,udp,http
	Param  TcpServerInterfaceParam `json:"param"`  //接口参数
	Conn   net.Conn                `json:"-"`      //通信句柄
	Status ConnectStatus           `json:"status"` //连接状态
}

var CommunicationTcpServerMap = make([]*CommunicationTcpServerTemplate, 0)

func (c *CommunicationTcpServerTemplate) Open() bool {
	listener, err := net.Listen("tcp", "0.0.0.0"+":"+c.Param.Port)
	if err != nil {
		setting.ZAPS.Errorf("通信TcpServer[%s]打开失败 %v", c.Name, err)
		c.Status = CommIsUnConnect
		return false
	}
	setting.ZAPS.Infof("通信TcpServer[%s]打开成功", c.Name)
	c.Status = CommIsConnect

	go c.Accept(listener)

	return true
}

func (c *CommunicationTcpServerTemplate) Accept(listener net.Listener) bool {

	for {
		conn, err := listener.Accept()
		if err != nil {
			setting.ZAPS.Errorf("通信TcpServer[%s]监听失败 %v", c.Name, err)
			continue
		}
		//c.Close()
		c.Conn = conn
		setting.ZAPS.Infof("通信TcpServer[%s]监听成功", c.Name)
	}
}

func (c *CommunicationTcpServerTemplate) Close() bool {
	if c.Conn != nil {
		err := c.Conn.Close()
		if err != nil {
			setting.ZAPS.Errorf("通信TcpServer[%s]关闭失败 %v", c.Name, err)
			return false
		}
		setting.ZAPS.Infof("通信TcpServer[%s]关闭成功", c.Name)
		c.Status = CommIsUnConnect
	}
	c.Conn = nil
	return true
}

func (c *CommunicationTcpServerTemplate) WriteData(data []byte) ([]byte, error) {

	if c.Conn != nil {
		_, err := c.Conn.Write(data)
		if err != nil {
			//setting.ZAPS.Errorf("%s TcpServer write err %v", c.Name, err)
			return nil, nil
		}
		return nil, nil
	}
	return nil, nil
}

func (c *CommunicationTcpServerTemplate) ReadData(data []byte) ([]byte, error) {

	if c.Conn != nil {
		cnt, err := c.Conn.Read(data)
		//setting.ZAPS.Debugf("%s,TcpServer read data cnt %v", c.Name, cnt)
		if err != nil {
			//setting.ZAPS.Errorf("%s,TcpServer read err,%v", c.Name, err)
			return nil, err
		}
		return data[:cnt], nil
	}
	return nil, errors.New("socket未打开")
}

func (c *CommunicationTcpServerTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationTcpServerTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationTcpServerTemplate) GetInterval() string {
	return c.Param.Interval
}

func (c *CommunicationTcpServerTemplate) GetType() int {
	return CommTypeTcpServer
}

func ReadCommTcpServerInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commTcpServerInterface.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such directory") {
			setting.ZAPS.Debug("打开通信接口[TcpServer]通信接口配置json文件失败")
		} else {
			setting.ZAPS.Debugf("打开通信接口[TcpServer]通信接口配置json文件失败 %v", err)
		}
		return false
	}
	err = json.Unmarshal(data, &CommunicationTcpServerMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[TcpServer]通信接口配置json文件格式化失败")
		return false
	}
	setting.ZAPS.Debugf("打开通信接口[TcpServer]通信接口配置json文件成功")
	return true
}

func WriteCommTcpServerInterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationTcpServerMap)
	err := utils.FileWrite("./selfpara/commTcpServerInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入TcpServer通信接口配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入TcpServer通信接口配置json文件 %s", "成功")
}
