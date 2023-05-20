package mqttSagooIOT

import (
	"encoding/json"
	"gateway/setting"
	"gateway/utils"
	"time"
)

var writeTimer *time.Timer

func ReportServiceSagooIOTConfigInit() {
	writeTimer = time.AfterFunc(time.Second, func() {
		ReportServiceSagooIOTWriteParamJsonToFile()
	})
	writeTimer.Stop()
}

func ReportServiceSagooIOTReadParamFromJson() bool {
	type ReportServiceConfigParamSagooIOTTemplate struct {
		ServiceList []ReportServiceParamSagooIOTTemplate `json:"ServiceList"`
	}

	configParam := ReportServiceConfigParamSagooIOTTemplate{}
	data, err := utils.FileRead("./selfpara/reportServiceParamListSagooIOT.json")
	if err != nil {
		setting.ZAPS.Debugf("上报服务[SagooIOT]配置json文件读取失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &configParam)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[SagooIOT]配置json文件格式化失败")
		return false
	}
	setting.ZAPS.Debug("上报服务[SagooIOT]配置json文件读取成功")

	//初始化
	for _, v := range configParam.ServiceList {
		v.GWParam.ReportStatus = "offLine"
		for _, n := range v.NodeList {
			n.CommStatus = "offLine"
			n.ReportStatus = "offLine"
		}
		ReportServiceParamListSagooIOT.ServiceList = append(ReportServiceParamListSagooIOT.ServiceList, NewReportServiceParamSagooIOT(v.GWParam, v.NodeList))
	}

	return true
}

func ReportServiceSagooIOTWriteParamToJson() {
	writeTimer.Reset(time.Second)
}

func ReportServiceSagooIOTWriteParamJsonToFile() {
	utils.DirIsExist("./selfpara")
	sJson, _ := json.Marshal(ReportServiceParamListSagooIOT)
	err := utils.FileWrite("./selfpara/reportServiceParamListSagooIOT.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("上报服务[SagooIOT]配置json文件写入失败")
		return
	}
	setting.ZAPS.Debugf("上报服务[SagooIOT]配置json文件写入成功")
}
