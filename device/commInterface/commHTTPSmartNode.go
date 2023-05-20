package commInterface

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"gateway/setting"
	"gateway/utils"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type HTTPSmartNodeInterfaceParam struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Timeout  string `json:"timeout"`  //通信超时
	Interval string `json:"interval"` //通信间隔
	AppCode  string `json:"appCode"`  //应用编码
	AppKey   string `json:"appKey"`   //应用秘钥
}

type CommunicationHTTPSmartNodeTemplate struct {
	Name   string                      `json:"name"`   //接口名称
	Type   string                      `json:"type"`   //接口类型
	Param  HTTPSmartNodeInterfaceParam `json:"param"`  //接口参数
	Conn   net.Conn                    `json:"-"`      //通信句柄
	Status ConnectStatus               `json:"status"` //连接状态
}

var CommunicationHTTPSmartNodeMap = make([]*CommunicationHTTPSmartNodeTemplate, 0)

func (c *CommunicationHTTPSmartNodeTemplate) Open() bool {

	c.Status = CommIsConnect
	return true
}

func (c *CommunicationHTTPSmartNodeTemplate) Close() bool {

	c.Status = CommIsUnConnect
	return true
}

func (c *CommunicationHTTPSmartNodeTemplate) WriteData(data []byte) ([]byte, error) {

	unixStr := strconv.FormatInt(time.Now().Unix(), 10)
	appCodeStr := c.Param.AppCode
	appKeyStr := c.Param.AppKey

	//对time,appCode,appKey字符串进行升序排序
	sortStr := []string{unixStr, appCodeStr, appKeyStr}
	sort.Strings(sortStr)
	//字符串拼接
	str := sortStr[0] + sortStr[1] + sortStr[2]
	//md5校验
	sign := md5.Sum([]byte(str))
	signStr := fmt.Sprintf("%X", sign)

	urls := url.Values{}
	urls.Add("time", unixStr)
	urls.Add("sign", signStr)
	urlsParam := "http://" + c.Param.IP + ":" + c.Param.Port + "/app/" + c.Param.AppCode + "/api/protected/writeDeviceVariantData?" + urls.Encode()
	//setting.ZAPS.Debugf("urlsParam %s", urlsParam)

	//setting.ZAPS.Debugf("urlData %v", string(data))
	req, err := http.NewRequest("POST", urlsParam, bytes.NewBuffer(data))
	if err != nil {
		setting.ZAPS.Errorf("通信接口[%s]构建请求包错误%v", c.Name, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[%s]发送请求包错误%v", c.Name, err)
		return nil, nil
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

func (c *CommunicationHTTPSmartNodeTemplate) ReadData(data []byte) ([]byte, error) {

	unixStr := strconv.FormatInt(time.Now().Unix(), 10)
	appCodeStr := c.Param.AppCode
	appKeyStr := c.Param.AppKey

	//对time,appCode,appKey字符串进行升序排序
	sortStr := []string{unixStr, appCodeStr, appKeyStr}
	sort.Strings(sortStr)
	//字符串拼接
	str := sortStr[0] + sortStr[1] + sortStr[2]
	//md5校验
	sign := md5.Sum([]byte(str))
	signStr := fmt.Sprintf("%X", sign)

	urls := url.Values{}
	urls.Add("time", unixStr)
	urls.Add("sign", signStr)
	urlsParam := "http://" + c.Param.IP + ":" + c.Param.Port + "/app/" + appCodeStr + "/api/protected/queryDeviceVariantData?" + urls.Encode()
	//setting.ZAPS.Debugf("urlsParam %s", urlsParam)

	//setting.ZAPS.Debugf("urlData %v", string(data))
	req, err := http.NewRequest("POST", urlsParam, bytes.NewBuffer(data))
	if err != nil {
		setting.ZAPS.Errorf("通信接口[%s]构建请求包错误%v", c.Name, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[%s]发送请求包错误%v", c.Name, err)
		return nil, nil
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

func (c *CommunicationHTTPSmartNodeTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationHTTPSmartNodeTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationHTTPSmartNodeTemplate) GetInterval() string {
	return c.Param.Interval
}

func (c *CommunicationHTTPSmartNodeTemplate) GetType() int {
	return CommTypeHTTPSmartNode
}

func ReadCommHTTPSmartNodeInterfaceListFromJson() bool {

	data, err := utils.FileRead("./selfpara/commHTTPSmartNodeInterface.json")
	if err != nil {
		setting.ZAPS.Debugf("打开通信接口[HTTPSmartNode]通信接口配置json文件失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &CommunicationHTTPSmartNodeMap)
	if err != nil {
		setting.ZAPS.Errorf("通信接口[HTTPSmartNode]通信接口配置json文件格式化失败 %v", err)
		return false
	}
	setting.ZAPS.Debugf("打开通信接口[HTTPSmartNode]通信接口配置json文件成功")
	return true
}

func WriteCommHTTPSmartNodeInterfaceListToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(CommunicationHTTPSmartNodeMap)
	err := utils.FileWrite("./selfpara/commHTTPSmartNodeInterface.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("写入通信接口[HTTPSmartNode]通信接口配置json文件 %s %v", "失败", err)
		return
	}
	setting.ZAPS.Infof("写入通信接口[HTTPSmartNode]通信接口配置json文件 %s", "成功")
}
