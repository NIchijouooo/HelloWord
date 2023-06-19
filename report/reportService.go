package report

import (
	"encoding/json"
	"gateway/report/modbusTCP"
	mqttEmqx "gateway/report/mqttEMQX"
	"gateway/report/mqttFeisjy"
	"gateway/report/mqttZxJs"
	"gateway/report/reportModel"
	"gateway/setting"
	"gateway/utils"
	"strings"
)

type ReportServiceAPI interface {
	GWLogIn()
	GWLogOut()
	NodesLogIn()
	NodesLogOut()
	GWPropertyReport()
	NodesPropertyReport()
}

type ProtocolTemplate struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

var ReportProtocols = make([]ProtocolTemplate, 0)

func init() {

}

func ReportServiceInit() {

	ReadReportProtocol()

	reportModel.ReportModelInit()

	mqttEmqx.ReportServiceEmqxInit()

	mqttFeisjy.ReportServiceFeisjyInit()

	mqttZxJs.ReportServiceZxjsInit()

	//mqttThingsBoard.ReportServiceThingsBoardInit()

	//mqttSagooIOT.ReportServiceSagooIOTInit()

	modbusTCP.ReportServiceMBTCPInit()

}

func ReadReportProtocol() {
	data, err := utils.FileRead("./selfpara/reportProtocol.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			setting.ZAPS.Debug("打开上报服务协议配置json文件失败[如果未配置，可以忽略]")
		} else {
			setting.ZAPS.Errorf("打开上报服务协议配置json文件失败 %v", err)
		}
		return
	}
	err = json.Unmarshal(data, &ReportProtocols)
	if err != nil {
		setting.ZAPS.Errorf("上报服务协议配置json文件格式化失败 %v", err)
		return
	}
	setting.ZAPS.Info("打开上报服务协议配置json文件成功")
}
