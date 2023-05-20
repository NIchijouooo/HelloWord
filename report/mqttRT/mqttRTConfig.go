package mqttRT

import (
	"encoding/json"
	"gateway/setting"
	"gateway/utils"
	"time"
)

var writeTimer *time.Timer

func ReportServiceRTConfigInit() {
	writeTimer = time.AfterFunc(time.Second, func() {
		ReportServiceRTWriteParamJsonToFile()
	})
	writeTimer.Stop()
}

func ReportServiceRTReadParamFromJson() bool {
	type ReportServiceConfigParamRTTemplate struct {
		ServiceList []ReportServiceParamRTTemplate `json:"ServiceList"`
	}

	configParam := ReportServiceConfigParamRTTemplate{}
	data, err := utils.FileRead("./selfpara/reportServiceParamListRT.json")
	if err != nil {
		setting.ZAPS.Debugf("上报服务[RT]配置json文件读取失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &configParam)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[RT]配置json文件格式化失败")
		return false
	}
	setting.ZAPS.Debug("上报服务[RT]配置json文件读取成功")

	//初始化
	for _, v := range configParam.ServiceList {
		v.GWParam.ReportStatus = "offLine"
		for _, n := range v.NodeList {
			n.CommStatus = "offLine"
			n.ReportStatus = "offLine"
		}
		ReportServiceParamListRT.ServiceList = append(ReportServiceParamListRT.ServiceList, NewReportServiceParamRT(v.GWParam, v.NodeList))
	}

	return true
}

func ReportServiceRTWriteParamToJson() {
	writeTimer.Reset(time.Second)
}

func ReportServiceRTWriteParamJsonToFile() {
	utils.DirIsExist("./selfpara")
	sJson, _ := json.Marshal(ReportServiceParamListRT)
	err := utils.FileWrite("./selfpara/reportServiceParamListRT.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[RT]配置json文件写入失败")
		return
	}
	setting.ZAPS.Debugf("上报服务[RT]配置json文件写入成功")
}
