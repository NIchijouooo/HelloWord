package commInterface

import (
	"encoding/json"
	"errors"
	"gateway/setting"
	"gateway/utils"
	"net"
	"time"
)

type TcpClientInterfaceParam struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Timeout  string `json:"timeout"`  //通信超时
	Interval string `json:"interval"` //通信间隔
}

type CommunicationTcpClientTemplate struct {
	Name   string                  `json:"name"`   //接口名称
	Type   string                  `json:"type"`   //接口类型,比如serial,TcpClient,udp,http
	Param  TcpClientInterfaceParam `json:"param"`  //接口参数
	Conn   net.Conn                `json:"-"`      //通信句柄
	Status ConnectStatus           `json:"status"` //连接状态
}

var CommunicationTcpClientMap = make([]*CommunicationTcpClientTemplate, 0)

func (c *CommunicationTcpClientTemplate) Open() bool {
	conn, err := net.DialTimeout("tcp", c.Param.IP+":"+c.Param.Port, 500*time.Millisecond)
	if err != nil {
		setting.ZAPS.Errorf("通信TCP客户端[%s]接口打开失败 %v", c.Name, err)
		c.Conn = nil
		c.Status = CommIsUnConnect
		return false
	}
	setting.ZAPS.Infof("通信TCP客户端[%s]接口打开成功", c.Name)
	c.Status = CommIsConnect

	c.Conn = conn
	return true
}

func (c *CommunicationTcpClientTemplate) Close() bool {
	if c.Conn != nil {
		err := c.Conn.Close()
		if err != nil {
			setting.ZAPS.Errorf("通信TCP客户端[%s]接口关闭失败 %v", c.Name, err)
			return false
		}
		setting.ZAPS.Infof("通信TCP客户端[%s]接口关闭成功", c.Name)
		c.Status = CommIsUnConnect
	}
	return true
}

func (c *CommunicationTcpClientTemplate) WriteData(data []byte) ([]byte, error) {

	if c.Conn != nil {
		_, err := c.Conn.Write(data)
		if err != nil {
			setting.ZAPS.Errorf("通信TCP客户端[%s]接口写失败 %v", c.Name, err)
			//c.Close()
			c.Open()
			return nil, nil
		}
		return nil, nil
	} else {
		c.Open()
	}
	return nil, nil
}

func (c *CommunicationTcpClientTemplate) ReadData(data []byte) ([]byte, error) {

	if c.Conn != nil {
		cnt, err := c.Conn.Read(data)
		//setting.ZAPS.Debugf("%s,TcpClient read data cnt %v", c.Name, cnt)
		if err != nil {
			//setting.ZAPS.Errorf("%s,TcpClient read err,%v", c.Name, err)
			return nil, err
		}
		return data[:cnt], nil
	}
	return nil, errors.New("socket未打开")
}

func (c *CommunicationTcpClientTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationTcpClientTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationTcpClientTemplate) GetInterval() string {
	return c.Param.Interval
}

func (c *CommunicationTcpClientTemplate) GetType() int {
	return CommTypeTcpClient
}

func ReadCommTcpClientInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commTcpClientInterface.json")
	if err != nil {
		setting.ZAPS.Debugf("打开通信接口[Tcp客户端]通信接口配置json文件失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &CommunicationTcpClientMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[Tcp客户端]通信接口配置json文件格式化失败 %v", err)
		return false
	}
	setting.ZAPS.Debugf("打开通信接口[Tcp客户端]通信接口配置json文件成功")
	return true
}

func WriteCommTcpClientInterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationTcpClientMap)
	err := utils.FileWrite("./selfpara/commTcpClientInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入通信接口[Tcp客户端]通信接口配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入通信接口[Tcp客户端]通信接口配置json文件 %s", "成功")
}
