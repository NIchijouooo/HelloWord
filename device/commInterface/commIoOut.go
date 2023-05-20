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

type IoOutInterfaceParam struct {
	Name string   `json:"name"`
	FD   *os.File `json:"-"`
}

type CommunicationIoOutTemplate struct {
	Name   string              `json:"name"`   //接口名称
	Type   string              `json:"type"`   //接口类型,比如serial,IoOut,udp,http
	Param  IoOutInterfaceParam `json:"param"`  //接口参数
	Status ConnectStatus       `json:"status"` //连接状态
}

var CommunicationIoOutMap = make([]*CommunicationIoOutTemplate, 0)

func (c *CommunicationIoOutTemplate) Open() bool {

	//fd, err := os.OpenFile(c.Param.Name, os.O_RDWR, 0666)
	//if err != nil {
	//	setting.ZAPS.Errorf("通信接口[开关量输出][%s]打开失败 %v", c.Param.Name, err)
	//	c.Status = CommIsUnConnect
	//	return false
	//}
	//setting.ZAPS.Debugf("通信接口[开关量输出][%s]打开成功", c.Param.Name)
	//c.Param.FD = fd
	c.Status = CommIsConnect
	return true
}

func (c *CommunicationIoOutTemplate) Close() bool {

	if c.Param.FD != nil {
		err := c.Param.FD.Close()
		if err != nil {
			setting.ZAPS.Errorf("通信接口[开关量输出][%s]关闭失败 %v", c.Param.Name, err)
			return false
		}
		setting.ZAPS.Infof("通信接口[开关量输出][%s]关闭成功", c.Param.Name)
		c.Status = CommIsUnConnect
	}

	return true
}

func (c *CommunicationIoOutTemplate) WriteData(data []byte) ([]byte, error) {

	ioOut := struct {
		Name  string `json:"name"`
		Value []byte `json:"value"`
	}{}

	_ = json.Unmarshal(data, &ioOut)

	setting.ZAPS.Debugf("开关量输出通信接口名称[%v],设定值%+v", ioOut.Name, ioOut)
	fd, err := os.OpenFile(ioOut.Name, os.O_RDWR, 0777)
	if err != nil {
		setting.ZAPS.Errorf("开关量输出通信接口[%s]打开失败 %v", c.Param.Name, err)
		return nil, err
	}
	defer fd.Close()

	_, err = fd.Write(ioOut.Value)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[开关量输出][%s]写入失败 %v", c.Param.Name, err)
	}
	return nil, nil
}

func (c *CommunicationIoOutTemplate) ReadData(data []byte) ([]byte, error) {

	ioOut := struct {
		Name  string `json:"name"`
		Value []byte `json:"value"`
	}{}

	_ = json.Unmarshal(data, &ioOut)

	//setting.ZAPS.Debugf("开关量输出通信接口名称[%v]", ioOut.Name)

	rData := make(chan []byte)
	go func(name string) {
		data, err := ioutil.ReadFile(name)
		if err != nil {
			if err != io.EOF {
				setting.ZAPS.Errorf("开关量输出通信接口[%s]读取失败 %v", c.Param.Name, err)
			} else {
				//setting.ZAPS.Debugf("开关量输出通信接口[%s]读取完成 %v", c.Param.Name, data[:cnt])
			}
		}
		rData <- data
	}(ioOut.Name)

	select {
	case buf := <-rData:
		{
			return buf, nil
		}
	case <-time.After(100 * time.Millisecond):
		{
			//setting.ZAPS.Errorf("开关量输出通信接口[%s]读取超时", c.Param.Name)
			return make([]byte, 0), errors.New("读取文件超时")
		}
	}
}

func (c *CommunicationIoOutTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationIoOutTemplate) GetTimeOut() string {
	return ""
}

func (c *CommunicationIoOutTemplate) GetInterval() string {
	return ""
}

func (c *CommunicationIoOutTemplate) GetType() int {
	return CommTypeIoOut
}

func ReadCommIoOutInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commIoOutInterface.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such directory") {
			setting.ZAPS.Debug("打开通信接口[开关量输出]配置json文件失败")
		} else {
			setting.ZAPS.Debugf("打开通信接口[开关量输出]配置json文件失败 %v", err)
		}

		return false
	}
	err = json.Unmarshal(data, &CommunicationIoOutMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[开关量输出]配置json文件格式化失败 %v", err)
		return false
	}
	setting.ZAPS.Debugf("通信接口[开关量输出]配置json文件格式化成功")
	return true
}

func WriteCommIoOutInterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationIoOutMap)
	err := utils.FileWrite("./selfpara/commIoOutInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入通信接口[开关量输出]配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入通信接口[开关量输出]配置json文件 %s", "成功")
}
