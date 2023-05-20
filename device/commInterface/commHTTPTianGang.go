package commInterface

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"gateway/setting"
	"gateway/utils"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type HTTPTianGangInterfaceParam struct {
	Timeout string `json:"timeout"` //通信超时
}

type CommunicationHTTPTianGangTemplate struct {
	Name   string                     `json:"name"`   //接口名称
	Type   string                     `json:"type"`   //接口类型
	Param  HTTPTianGangInterfaceParam `json:"param"`  //接口参数
	Conn   net.Conn                   `json:"-"`      //通信句柄
	Status ConnectStatus              `json:"status"` //连接状态
}

var CommunicationHTTPTianGangMap = make([]*CommunicationHTTPTianGangTemplate, 0)

func (c *CommunicationHTTPTianGangTemplate) Open() bool {
	//conn, err := net.DialTimeout("tcp", c.Param.IP+":"+c.Param.Port, 500*time.Millisecond)
	//if err != nil {
	//	setting.ZAPS.Errorf("通信HTTP[%s]接口打开失败 %v", c.Name, err)
	//	c.Conn = nil
	//	return false
	//}
	//setting.ZAPS.Debugf("通信HTTPTianGang[%s]接口打开成功", c.Name)
	//
	//c.Conn = conn
	c.Status = CommIsConnect
	return true
}

func (c *CommunicationHTTPTianGangTemplate) Close() bool {
	//if c.Conn != nil {
	//	err := c.Conn.Close()
	//	if err != nil {
	//		setting.ZAPS.Errorf("通信HTTPTianGang[%s]接口关闭失败 %v", c.Name, err)
	//		return false
	//	}
	//	setting.ZAPS.Debugf("通信HTTPTianGang[%s]接口关闭成功", c.Name)
	//}
	c.Status = CommIsUnConnect
	return true
}

func (c *CommunicationHTTPTianGangTemplate) WriteData(data []byte) ([]byte, error) {

	//if c.Conn != nil {
	//	cnt, err := c.Conn.Write(data)
	//	if err != nil {
	//		setting.ZAPS.Errorf("通信HTTPTianGang[%s]接口写失败 %v", c.Name, err)
	//		c.Close()
	//		c.Open()
	//		return 0
	//	}
	//	return cnt
	//} else {
	//	c.Open()
	//}
	return nil, nil
}

func (c *CommunicationHTTPTianGangTemplate) ReadData(data []byte) ([]byte, error) {

	queryVariantData := struct {
		DeviceCode string `json:"deviceCode"`
	}{}

	err := json.Unmarshal(data, &queryVariantData)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[%s]设备地址JSON格式化错误%v", c.Name, err)
		return nil, err
	}

	type RequestParamTemplate struct {
		Code       string `json:"code"`
		Param      string `json:"param"`
		EncryptKey string `json:"encryptkey"`
	}

	md5Sum := md5.New()
	//timeStr := time.Now().Format("2006010215") + "Pl_Wuxi"
	//setting.ZAPS.Debugf("timeStr %s", timeStr)
	md5Sum.Write([]byte(time.Now().Format("2006010215") + "Pl_Wuxi"))
	key := hex.EncodeToString(md5Sum.Sum(nil))
	//setting.ZAPS.Debugf("key %s", key)
	reqParam := RequestParamTemplate{
		Code:       "303_E",
		Param:      "0_" + time.Now().Format("2006010215") + "_" + queryVariantData.DeviceCode,
		EncryptKey: key,
	}
	reqParamJson, err := json.Marshal(&reqParam)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[%s]JSON格式化错误%v", c.Name, err)
		return nil, err
	}
	reqBody := strings.NewReader(string(reqParamJson))

	url := "http://" + "221.2.162.151" + ":" + "52007" + "/External/" + "WaterDataList"
	//setting.ZAPS.Debugf("url %s", url)
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[%s]构建请求包错误%v", c.Name, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[%s]发送请求包错误%v", c.Name, err)
		return nil, err
	}
	defer resp.Body.Close()

	//setting.ZAPS.Debugf("应答包回应Status %v", resp.Status)
	//setting.ZAPS.Debugf("应答包回应Header %v", resp.Header)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[%s]读取应答包内容错误%v", c.Name, err)
	}
	//setting.ZAPS.Debugf("应答包回应Body %v", string(respBody))

	return respBody, nil
}

func (c *CommunicationHTTPTianGangTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationHTTPTianGangTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationHTTPTianGangTemplate) GetInterval() string {
	return "0"
}

func (c *CommunicationHTTPTianGangTemplate) GetType() int {
	return CommTypeHTTPTianGang
}

func ReadCommHTTPTianGangInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commHTTPTianGangInterface.json")
	if err != nil {
		setting.ZAPS.Debugf("打开通信接口[HTTPTianGang]通信接口配置json文件失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &CommunicationHTTPTianGangMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[HTTPTianGang]通信接口配置json文件格式化失败 %v", err)
		return false
	}
	setting.ZAPS.Debugf("打开通信接口[HTTPTianGang]通信接口配置json文件成功")
	return true
}

func WriteCommHTTPTianGangInterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationHTTPTianGangMap)
	err := utils.FileWrite("./selfpara/commHTTPTianGangInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入通信接口[HTTPTianGang]通信接口配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入通信接口[HTTPTianGang]通信接口配置json文件 %s", "成功")
}
