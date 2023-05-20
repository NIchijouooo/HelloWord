package commInterface

import (
	"encoding/json"
	"errors"
	"gateway/setting"
	"gateway/utils"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type IoInInterfaceParam struct {
	Name string   `json:"name"`
	FD   *os.File `json:"-"`
}

type CommunicationIoInTemplate struct {
	Name   string             `json:"name"`   //接口名称
	Type   string             `json:"type"`   //接口类型,比如serial,IoIn,udp,http
	Param  IoInInterfaceParam `json:"param"`  //接口参数
	Status ConnectStatus      `json:"status"` //连接状态
}

var CommunicationIoInMap = make([]*CommunicationIoInTemplate, 0)

func (c *CommunicationIoInTemplate) Open() bool {

	//fd, err := os.OpenFile(c.Param.Name, os.O_RDWR, 0777)
	//if err != nil {
	//	setting.ZAPS.Errorf("开关量输入通信接口[%s]打开失败 %v", c.Param.Name, err)
	//	return false
	//}
	//setting.ZAPS.Debugf("开关量输入通信接口[%s]打开成功", c.Param.Name)
	//c.Param.FD = fd
	c.Status = CommIsConnect
	return true
}

func (c *CommunicationIoInTemplate) Close() bool {

	if c.Param.FD != nil {
		err := c.Param.FD.Close()
		if err != nil {
			setting.ZAPS.Errorf("开关量输入通信接口[%s]关闭失败 %v", c.Param.Name, err)
			return false
		}
		setting.ZAPS.Infof("开关量输入通信接口[%s]关闭成功", c.Param.Name)
		c.Status = CommIsUnConnect
	}

	return true
}

func (c *CommunicationIoInTemplate) WriteData(data []byte) ([]byte, error) {

	return nil, nil
}

func (c *CommunicationIoInTemplate) ReadData(data []byte) ([]byte, error) {

	//setting.ZAPS.Debugf("开关量输入通信接口名称[%v]", string(data))

	rData := make(chan []byte)
	go func(name string) {
		data, err := ioutil.ReadFile(name)
		if err != nil {
			if err != io.EOF {
				setting.ZAPS.Errorf("开关量输入通信接口[%s]读取失败 %v", c.Param.Name, err)
			} else {
				//setting.ZAPS.Debugf("开关量输入通信接口[%s]读取完成 %v", c.Param.Name, data[:cnt])
			}
		}
		rData <- data
	}(string(data))

	select {
	case buf := <-rData:
		{
			return buf, nil
		}
	case <-time.After(100 * time.Millisecond):
		{
			//setting.ZAPS.Errorf("开关量输入通信接口[%s]读取超时", c.Param.Name)
			return make([]byte, 0), errors.New("读取文件超时")
		}
	}

	//setting.ZAPS.Debugf("开关量输入通信接口[%s]读取完成 %s", c.Param.Name, rData[:rCnt])
}

func (c *CommunicationIoInTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationIoInTemplate) GetTimeOut() string {
	return "0"
}

func (c *CommunicationIoInTemplate) GetInterval() string {
	return ""
}

func (c *CommunicationIoInTemplate) GetType() int {
	return CommTypeIoIn
}

func ReadCommIoInInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commIoInInterface.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such directory") {
			setting.ZAPS.Debug("打开通信接口[开关量输入]配置json失败，文件不存在")
		} else {
			setting.ZAPS.Debugf("打开通信接口[开关量输入]配置json文件失败 %v", err)
		}
		return false
	}
	err = json.Unmarshal(data, &CommunicationIoInMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[开关量输入]配置json文件格式化失败 %v", err)
		return false
	}
	setting.ZAPS.Info("通信接口[开关量输入]配置json文件格式化成功")
	return true
}

func WriteCommIoInInterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationIoInMap)
	err := utils.FileWrite("./selfpara/commIoInInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入通信接口[开关量输入]配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入通信接口[开关量输入]配置json文件 %s", "成功")
}
