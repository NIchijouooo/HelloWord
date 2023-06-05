package setting

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/utils"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

type NetworkNamesTemplate struct {
	Name []string `json:"name"`
}

type NetworkParamTemplate struct {
	Index       int                        `json:"-"`
	Name        string                     `json:"name"`  // e.g., "en0", "lo0", "eth0.100"
	MTU         int                        `json:"mtu"`   // maximum transmission unit
	MAC         string                     `json:"mac"`   // IEEE MAC-48, EUI-48 and EUI-64 form
	Flags       string                     `json:"flags"` // e.g., FlagUp, FlagLoopback, FlagMulticast
	IP          string                     `json:"ip"`
	Netmask     string                     `json:"netmask"`
	Gateway     string                     `json:"gateway"`
	ConfigParam NetworkConfigParamTemplate `json:"configParam"`
}

type NetworkConfigParamTemplate struct {
	Name          string `json:"name"`
	ConfigEnable  bool   `json:"configEnable"`
	DHCPEnable    bool   `json:"dhcpEnable"`
	ConfigIP      string `json:"configIP"`
	ConfigNetmask string `json:"configNetmask"`
	ConfigGateway string `json:"configGateway"`
}

var NetworkConfigParams = make([]NetworkConfigParamTemplate, 0)
var NetworkParams = make([]NetworkParamTemplate, 0)

func init() {

}

func UpdateNetworkParams() {
	param := NetworkParamTemplate{}

	if len(NetworkConfigParams) == 0 {
		return
	}
	for k, v := range NetworkConfigParams {
		param.Index = k
		param.Name = v.Name
		param.ConfigParam = v
		NetworkParams = append(NetworkParams, param)
	}

}

func GetNetworkNames() []string {

	names := make([]string, 0)

	inters, err := net.Interfaces()
	if err != nil {
		ZAPS.Errorf("获取网卡信息错误 %v", err)
		return names
	}

	for _, v := range inters {
		names = append(names, v.Name)
	}

	return names
}

// 获取当前网络参数
func GetNetworkParams() []NetworkParamTemplate {

	inters, err := net.Interfaces()
	if err != nil {
		ZAPS.Errorf("获取网卡信息错误 %v", err)
		return NetworkParams
	}

	for k, n := range NetworkParams {
		for _, v := range inters {
			if n.Name == v.Name {
				NetworkParams[k].Index = v.Index
				NetworkParams[k].Name = v.Name
				NetworkParams[k].MTU = v.MTU
				NetworkParams[k].MAC = strings.ToUpper(hex.EncodeToString(v.HardwareAddr))
				NetworkParams[k].Flags = v.Flags.String()
				NetworkParams[k].IP, NetworkParams[k].Netmask = GetIPAndMask(v.Name)
				NetworkParams[k].Gateway = GetGateway(v.Name)
			}
		}
	}

	return NetworkParams
}

// 获取系统所有网卡参数
func GetNetworkInterface() []NetworkParamTemplate {

	params := make([]NetworkParamTemplate, 0)

	inters, err := net.Interfaces()
	if err != nil {
		ZAPS.Errorf("获取网卡信息错误 %v", err)
		return params
	}

	for _, v := range inters {
		param := NetworkParamTemplate{}
		param.Index = v.Index
		param.Name = v.Name
		param.MTU = v.MTU
		param.MAC = strings.ToUpper(hex.EncodeToString(v.HardwareAddr))
		param.Flags = v.Flags.String()
		param.IP, param.Netmask = GetIPAndMask(v.Name)
		param.Gateway = GetGateway(v.Name)

		params = append(params, param)
	}

	return params
}

func GetIPAndMask(name string) (ip, netmask string) {

	inter, err := net.InterfaceByName(name)
	if err != nil {
		return "", ""
	}
	address, _ := inter.Addrs()
	for _, addr := range address {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), net.IP(ip.Mask).String()
			}
		}
	}

	return "", ""
}

func GetGateway(name string) string {
	if runtime.GOOS == "linux" {
		out, err := exec.Command("/bin/sh", "-c",
			fmt.Sprintf("route -n | grep %s | grep UG | awk '{print $2}'", name)).Output()
		if err != nil {
			return ""
		}
		return strings.Trim(string(out), "\n")
	}

	return ""
}

func AddNetworkConfigParam(configParam NetworkConfigParamTemplate) error {

	names := GetNetworkNames()
	index := -1
	for k, v := range names {
		if v == configParam.Name {
			index = k
		}
	}
	if index == -1 {
		return errors.New("网卡不存在")
	}

	index = -1
	for k, v := range NetworkConfigParams {
		if v.Name == configParam.Name {
			index = k
		}
	}
	if index != -1 {
		return errors.New("网卡偏好设置已经存在")
	}

	//保存网卡偏好配置到JSON文件中
	NetworkConfigParams = append(NetworkConfigParams, configParam)
	WriteNetworkParamToJson(NetworkConfigParams)

	//追加新增网卡偏好配置到缓存中
	param := NetworkParamTemplate{
		Index:       len(NetworkParams),
		Name:        configParam.Name,
		ConfigParam: configParam,
	}
	NetworkParams = append(NetworkParams, param)

	return nil
}

// 修改网络参数
func ModifyNetworkConfigParam(configParam NetworkConfigParamTemplate) error {

	names := GetNetworkNames()
	index := -1
	for k, v := range names {
		if v == configParam.Name {
			index = k
		}
	}
	if index == -1 {
		return errors.New("网卡不存在")
	}

	index = -1
	for k, v := range NetworkConfigParams {
		if v.Name == configParam.Name {
			index = k
		}
	}
	if index == -1 {
		return errors.New("网卡偏好设置不存在")
	}

	//保存网卡偏好配置到JSON文件中
	NetworkConfigParams[index] = configParam
	WriteNetworkParamToJson(NetworkConfigParams)

	//更新新增网卡偏好配置到缓存中
	for k, v := range NetworkParams {
		if v.Name == configParam.Name {
			NetworkParams[k].ConfigParam = configParam
		}
	}

	return nil
}

// 删除网络参数
func DeleteNetworkConfigParam(name string) error {

	names := GetNetworkNames()
	index := -1
	for k, v := range names {
		if v == name {
			index = k
		}
	}
	if index == -1 {
		return errors.New("网卡不存在")
	}

	index = -1
	for k, v := range NetworkConfigParams {
		if v.Name == name {
			index = k
		}
	}
	if index == -1 {
		return errors.New("网卡偏好设置不存在")
	}

	//保存网卡偏好配置到JSON文件中
	NetworkConfigParams = append(NetworkConfigParams[:index], NetworkConfigParams[index+1:]...)
	WriteNetworkParamToJson(NetworkConfigParams)

	//更新新增网卡偏好配置到缓存中
	index = -1
	for k, v := range NetworkParams {
		if v.Name == name {
			index = k
		}
	}
	NetworkParams = append(NetworkParams[:index], NetworkParams[index+1:]...)

	return nil
}

func (n *NetworkConfigParamTemplate) CmdSetDHCP() error {

	//非阻塞,动态获取IP有可能不成功
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("udhcpc -i %s", n.Name)).Output()
	if err != nil {
		ZAPS.Debugf("网卡[%s]动态获取IP失败 %s %v", n.Name, string(out), err)
		return err
	}
	ZAPS.Debugf("网卡[%s]动态获取IP成功 %s", n.Name, string(out))

	return nil
}

func (n *NetworkConfigParamTemplate) CmdSetStaticIP() {

	//strNetMask := "netmask " + n.Netmask
	//cmd := exec.Command("ifconfig",
	//	n.Name,
	//	n.IP,
	//	strNetMask)
	//
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//cmd.Start() //执行到此,直接往后执行

	out, err := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("ifconfig %s %s netmask %s", n.Name, n.ConfigIP, n.ConfigNetmask)).Output()
	if err != nil {
		ZAPS.Debugf("网卡[%s]设置IP[%s]Netmask[%s]失败 %s %v", n.Name, n.ConfigIP, n.ConfigNetmask, string(out), err)
	} else {
		ZAPS.Debugf("网卡[%s]设置IP[%s]Netmask[%s]成功", n.Name, n.ConfigIP, n.ConfigNetmask)
	}

	// UPDATA QJHui 2023/6/5 修改路由表，导致4G使用问题
	//out, err = exec.Command("/sbin/route", "add", "default", "gw", n.ConfigGateway).Output()
	//if err != nil {
	//	ZAPS.Debugf("网卡[%s]添加默认网关[%s]失败 %s %v", n.Name, n.ConfigGateway, string(out), err)
	//	return
	//}
	ZAPS.Debugf("网卡[%s]添加默认网关[%s]成功", n.Name, n.ConfigGateway)

}

func ReadNetworkParamFromJson() bool {
	data, err := utils.FileRead("./selfpara/networkConfigParams.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			ZAPS.Debug("打开网络配置json文件失败")
		} else {
			ZAPS.Errorf("打开网络配置json文件失败 %v", err)
		}
		return false
	}

	err = json.Unmarshal(data, &NetworkConfigParams)
	if err != nil {
		ZAPS.Errorf("格式化网络配置json文件失败 %v", err)
		return false
	}
	ZAPS.Info("打开网络配置json文件成功")

	UpdateNetworkParams()

	for _, v := range NetworkConfigParams {
		if v.ConfigEnable == false {
			continue
		}
		if v.DHCPEnable == true {
			go v.CmdSetDHCP()
		} else {
			v.CmdSetStaticIP()
		}
	}

	return true
}

func WriteNetworkParamToJson(params []NetworkConfigParamTemplate) {
	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(params)
	err := utils.FileWrite("./selfpara/networkConfigParams.json", sJson)
	if err != nil {
		ZAPS.Warnf("写入网络配置json文件失败 %v", err)
		return
	}
	ZAPS.Info("写入网络配置json文件成功")
}

func GetIPByNetCardName(name string) (ip string) {

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil && iface.Name == name {
					ip = ipnet.IP.String()
					break
				}
			}
		}
	}

	return ip
}
