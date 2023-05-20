package contorl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/httpServer/model"
	"gateway/report"
	"gateway/report/modbusTCP"
	mqttEmqx "gateway/report/mqttEMQX"
	"gateway/report/mqttFeisjy"
	"gateway/report/mqttRT"
	"gateway/report/mqttSagooIOT"
	"gateway/report/mqttThingsBoard"
	"gateway/report/reportModel"

	"gateway/setting"
	"gateway/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"

	"github.com/gin-gonic/gin"
)

func ApiGetReportProtocol(context *gin.Context) {
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取上报服务协议成功",
		Data:    report.ReportProtocols,
	})
}

func ApiAddReportGWParam(context *gin.Context) {

	type ReportServiceTemplate struct {
		ServiceName string      `json:"serviceName"`
		IP          string      `json:"ip"`
		Port        string      `json:"port"`
		ReportTime  int         `json:"reportTime"`
		Protocol    string      `json:"protocol"`
		Param       interface{} `json:"param"`
	}

	var Param json.RawMessage
	serviceParam := ReportServiceTemplate{
		Param: &Param,
	}

	//从context中取出缓存，然后再赋值给context，用于第二次校验
	contextRawData, _ := context.GetRawData()
	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(contextRawData))

	err := context.ShouldBindJSON(&serviceParam)
	if err != nil {
		setting.ZAPS.Errorf("增加上报服务参数JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加上报服务参数JSON格式化错误",
			Data:    "",
		})
		return
	}

	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(contextRawData))

	switch serviceParam.Protocol {
	//case "Aliyun.MQTT":
	//	ReportServiceGWParamAliyun := mqttAliyun.ReportServiceGWParamAliyunTemplate{}
	//	err := context.ShouldBindJSON(&ReportServiceGWParamAliyun)
	//	if err != nil {
	//		setting.ZAPS.Errorf("增加上报服务Aliyun.MQTT参数JSON格式化错误[%v]", err)
	//	}
	//	mqttAliyun.ReportServiceParamListAliyun.AddReportService(ReportServiceGWParamAliyun)

	case "FSJY.MQTT": //gwai add 2023-04-05
		ReportServiceGWParamFeisjy := mqttFeisjy.ReportServiceGWParamFeisjyTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParamFeisjy)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务FSJY参数JSON格式化错误[%v]", err)
		}
		err = mqttFeisjy.ReportServiceParamListFeisjy.AddReportService(ReportServiceGWParamFeisjy)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务FSJY参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("增加上报服务FSJY.MQTT错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}
	case "EMQX.MQTT":
		ReportServiceGWParamEmqx := mqttEmqx.ReportServiceGWParamEmqxTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParamEmqx)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务EMQX参数JSON格式化错误[%v]", err)
		}
		err = mqttEmqx.ReportServiceParamListEmqx.AddReportService(ReportServiceGWParamEmqx)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务EMQX参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("增加上报服务EMQX.MQTT错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}
	case "RT.MQTT":
		ReportServiceGWParamRT := mqttRT.ReportServiceGWParamRTTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParamRT)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务RT参数JSON格式化错误[%v]", err)
		}
		err = mqttRT.ReportServiceParamListRT.AddReportService(ReportServiceGWParamRT)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务RT参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("增加上报服务RT.MQTT错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}
	case "SagooIOT.MQTT":
		ReportServiceGWParamSagooIOT := mqttSagooIOT.ReportServiceGWParamSagooIOTTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParamSagooIOT)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务SagooIOT参数JSON格式化错误[%v]", err)
		}
		err = mqttSagooIOT.ReportServiceParamListSagooIOT.AddReportService(ReportServiceGWParamSagooIOT)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务SagooIOT参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("增加上报服务SagooIOT.MQTT错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}

	case "ModbusTCP":
		ReportServiceGWParam := modbusTCP.ReportServiceMBTCPParamTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParam)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务ModbusTCP参数JSON格式化错误[%v]", err)
		}
		err = modbusTCP.ReportServiceMBTCPList.AddReportService(ReportServiceGWParam)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务ModbusTCP参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("增加上报服务ModbusTCP错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}
	case "ThingsBoard.MQTT":
		ReportServiceGWParamThingsBoard := mqttThingsBoard.ReportServiceGWParamThingsBoardTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParamThingsBoard)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务ThingsBoard参数JSON格式化错误[%v]", err)
		}
		mqttThingsBoard.ReportServiceParamListThingsBoard.AddReportService(ReportServiceGWParamThingsBoard)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务ThingsBoard参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("增加上报服务ThingsBoard错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}

	//case "Huawei.MQTT":
	//	ReportServiceGWParamHuawei := mqttHuawei.ReportServiceGWParamHuaweiTemplate{}
	//	err := context.ShouldBindJSON(&ReportServiceGWParamHuawei)
	//	if err != nil {
	//		setting.ZAPS.Errorf("增加上报服务HuaweiYun参数JSON格式化错误[%v]", err)
	//	}
	//	mqttHuawei.ReportServiceParamListHuawei.AddReportService(ReportServiceGWParamHuawei)
	default:
		setting.ZAPS.Error("增加上报服务参数错误[未知协议]")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "未知协议",
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "增加上报服务参数成功",
		Data:    "",
	})
}

func ApiModifyReportGWParam(context *gin.Context) {

	type ReportServiceTemplate struct {
		ServiceName string      `json:"serviceName"`
		IP          string      `json:"ip"`
		Port        string      `json:"port"`
		ReportTime  int         `json:"reportTime"`
		Protocol    string      `json:"protocol"`
		Param       interface{} `json:"param"`
	}

	var Param json.RawMessage
	serviceParam := ReportServiceTemplate{
		Param: &Param,
	}

	//从context中取出缓存，然后再赋值给context，用于第二次校验
	contextRawData, _ := context.GetRawData()
	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(contextRawData))

	err := context.ShouldBindJSON(&serviceParam)
	if err != nil {
		setting.ZAPS.Errorf("修改上报服务参数JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改上报服务参数JSON格式化错误",
			Data:    "",
		})
		return
	}

	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(contextRawData))

	switch serviceParam.Protocol {
	//case "Aliyun.MQTT":
	//	ReportServiceGWParamAliyun := mqttAliyun.ReportServiceGWParamAliyunTemplate{}
	//	err := context.ShouldBindJSON(&ReportServiceGWParamAliyun)
	//	if err != nil {
	//		setting.ZAPS.Errorf("修改上报服务Aliyun.MQTT参数JSON格式化错误[%v]", err)
	//		context.JSON(http.StatusOK, model.ResponseData{
	//			Code:    "1",
	//			Message: "修改上报服务Aliyun.MQTT参数JSON格式化错误",
	//			Data:    "",
	//		})
	//		return
	//	}
	//	mqttAliyun.ReportServiceParamListAliyun.AddReportService(ReportServiceGWParamAliyun)

	case "FSJY.MQTT": //gwai add 2023-04-05
		ReportServiceGWParamFeisjy := mqttFeisjy.ReportServiceGWParamFeisjyTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParamFeisjy)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务FSJY参数JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "修改上报服务FSJY.MQTT参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		err = mqttFeisjy.ReportServiceParamListFeisjy.ModifyReportService(ReportServiceGWParamFeisjy)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务FSJY参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("修改上报服务FSJY.MQTT错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}
	case "EMQX.MQTT":
		ReportServiceGWParamEmqx := mqttEmqx.ReportServiceGWParamEmqxTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParamEmqx)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务EMQX参数JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "修改上报服务EMQX.MQTT参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		err = mqttEmqx.ReportServiceParamListEmqx.ModifyReportService(ReportServiceGWParamEmqx)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务EMQX参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("修改上报服务EMQX.MQTT错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}
	case "RT.MQTT":
		ReportServiceGWParamRT := mqttRT.ReportServiceGWParamRTTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParamRT)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务RT参数JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "修改上报服务RT.MQTT参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		err = mqttRT.ReportServiceParamListRT.ModifyReportService(ReportServiceGWParamRT)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务RT参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("修改上报服务RT.MQTT错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}
	case "SagooIOT.MQTT":
		ReportServiceGWParamSagooIOT := mqttSagooIOT.ReportServiceGWParamSagooIOTTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParamSagooIOT)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务SagooIOT参数JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "修改上报服务SagooIOT.MQTT参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		err = mqttSagooIOT.ReportServiceParamListSagooIOT.ModifyReportService(ReportServiceGWParamSagooIOT)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务SagooIOT参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("修改上报服务SagooIOT.MQTT错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}
	case "ModbusTCP":
		ReportServiceModbusTCPParam := modbusTCP.ReportServiceMBTCPParamTemplate{}
		err := context.ShouldBindJSON(&ReportServiceModbusTCPParam)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务ModbusTCP参数JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "修改上报服务ModbusTCP参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		err = modbusTCP.ReportServiceMBTCPList.ModifyReportService(ReportServiceModbusTCPParam)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务ModbusTCP参数错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: fmt.Sprintf("修改上报服务ModbusTCP错误[%v]", err.Error()),
				Data:    "",
			})
			return
		}
	case "ThingsBoard.MQTT":
		ReportServiceGWParamThingsBoard := mqttThingsBoard.ReportServiceGWParamThingsBoardTemplate{}
		err := context.ShouldBindJSON(&ReportServiceGWParamThingsBoard)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务ThingsBoard参数JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "修改上报服务ThingsBoard.MQTT参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		mqttThingsBoard.ReportServiceParamListThingsBoard.AddReportService(ReportServiceGWParamThingsBoard)
	//case "Huawei.MQTT":
	//	ReportServiceGWParamHuawei := mqttHuawei.ReportServiceGWParamHuaweiTemplate{}
	//	err := context.ShouldBindJSON(&ReportServiceGWParamHuawei)
	//	if err != nil {
	//		setting.ZAPS.Errorf("修改上报服务HuaweiYun参数JSON格式化错误[%v]", err)
	//		context.JSON(http.StatusOK, model.ResponseData{
	//			Code:    "1",
	//			Message: "修改上报服务HuaweiYun.MQTT参数JSON格式化错误",
	//			Data:    "",
	//		})
	//		return
	//	}
	//	mqttHuawei.ReportServiceParamListHuawei.AddReportService(ReportServiceGWParamHuawei)
	default:
		setting.ZAPS.Error("修改上报服务参数错误[未知协议]")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "未知协议",
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改上报服务参数成功",
		Data:    "",
	})
}

func ApiGetReportGWParam(context *gin.Context) {

	type ReportServiceTemplate struct {
		Index        int         `json:"index"`
		ServiceName  string      `json:"serviceName"`
		IP           string      `json:"ip"`
		Port         string      `json:"port"`
		ReportTime   int         `json:"reportTime"`
		ReportStatus string      `json:"reportStatus"`
		Protocol     string      `json:"protocol"`
		Param        interface{} `json:"param"`
	}

	params := make([]ReportServiceTemplate, 0)

	//for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
	//
	//	ReportService := ReportServiceTemplate{}
	//	ReportService.ServiceName = v.GWParam.ServiceName
	//	ReportService.IP = v.GWParam.IP
	//	ReportService.Port = v.GWParam.Port
	//	ReportService.ReportTime = v.GWParam.ReportTime
	//	ReportService.Protocol = v.GWParam.Protocol
	//	ReportService.Param = v.GWParam.Param
	//	ReportService.ReportStatus = v.GWParam.ReportStatus
	//
	//	params = append(params, ReportService)
	//}

	for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
		ReportService := ReportServiceTemplate{}
		ReportService.Index = len(params)
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = v.GWParam.Param
		ReportService.ReportStatus = v.GWParam.ReportStatus

		params = append(params, ReportService)
	}

	//gwai add 2023-04-05
	for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
		ReportService := ReportServiceTemplate{}
		ReportService.Index = len(params)
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = v.GWParam.Param
		ReportService.ReportStatus = v.GWParam.ReportStatus

		params = append(params, ReportService)
	}

	for _, v := range mqttThingsBoard.ReportServiceParamListThingsBoard.ServiceList {

		ReportService := ReportServiceTemplate{}
		ReportService.Index = len(params)
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = v.GWParam.Param
		ReportService.ReportStatus = v.GWParam.ReportStatus

		params = append(params, ReportService)
	}

	//for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {
	//
	//	ReportService := ReportServiceTemplate{}
	//	ReportService.ServiceName = v.GWParam.ServiceName
	//	ReportService.IP = v.GWParam.IP
	//	ReportService.Port = v.GWParam.Port
	//	ReportService.ReportTime = v.GWParam.ReportTime
	//	ReportService.Protocol = v.GWParam.Protocol
	//	ReportService.Param = v.GWParam.Param
	//	ReportService.ReportStatus = v.GWParam.ReportStatus
	//
	//	params = append(params, ReportService)
	//}

	for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {

		ReportService := ReportServiceTemplate{}
		ReportService.Index = len(params)
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = v.GWParam.Param
		ReportService.ReportStatus = v.GWParam.ReportStatus

		params = append(params, ReportService)
	}

	for _, v := range mqttSagooIOT.ReportServiceParamListSagooIOT.ServiceList {
		ReportService := ReportServiceTemplate{}
		ReportService.Index = len(params)
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = v.GWParam.Param
		ReportService.ReportStatus = v.GWParam.ReportStatus

		params = append(params, ReportService)
	}

	for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
		ReportService := ReportServiceTemplate{}
		ReportService.Index = len(params)
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = v.GWParam.Param
		ReportService.ReportStatus = v.GWParam.ReportStatus

		params = append(params, ReportService)
	}

	for _, v := range modbusTCP.ReportServiceMBTCPList.ServiceList {

		param := struct {
			SlaveID                  int `json:"slaveID"`
			CoilStatusRegisterStart  int `json:"coilStatusRegStart"`
			CoilStatusRegisterCnt    int `json:"coilStatusRegCnt"`
			InputStatusRegisterStart int `json:"inputStatusRegStart"`
			InputStatusRegisterCnt   int `json:"inputStatusRegCnt"`
			HoldingRegisterStart     int `json:"holdingRegStart"`
			HoldingRegisterCnt       int `json:"holdingRegCnt"`
			InputRegisterStart       int `json:"inputRegStart"`
			InputRegisterCnt         int `json:"inputRegCnt"`
		}{
			SlaveID:                  v.GWParam.Param.SlaveID,
			CoilStatusRegisterStart:  v.GWParam.Param.CoilStatusRegisterStart,
			CoilStatusRegisterCnt:    v.GWParam.Param.CoilStatusRegisterCnt,
			InputStatusRegisterStart: v.GWParam.Param.InputStatusRegisterStart,
			InputStatusRegisterCnt:   v.GWParam.Param.InputStatusRegisterCnt,
			HoldingRegisterStart:     v.GWParam.Param.HoldingRegisterStart,
			HoldingRegisterCnt:       v.GWParam.Param.HoldingRegisterCnt,
			InputRegisterStart:       v.GWParam.Param.InputRegisterStart,
			InputRegisterCnt:         v.GWParam.Param.InputRegisterCnt,
		}

		ReportService := ReportServiceTemplate{}
		ReportService.Index = len(params)
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = param
		ReportService.ReportStatus = v.GWParam.ReportStatus

		params = append(params, ReportService)
	}

	sort.Slice(params, func(i, j int) bool {
		i = params[i].Index
		j = params[i].Index
		return i > j
	})

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取上报服务参数成功",
		Data:    params,
	})
}

func ApiDeleteReportGWParam(context *gin.Context) {

	param := struct {
		Name string `json:"serviceName"`
	}{}

	err := context.ShouldBindJSON(&param)
	if err != nil {
		setting.ZAPS.Errorf("删除上报服务参数JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除上报服务参数JSON格式化错误",
			Data:    "",
		})
		return
	}

	status := false

	////查看Aliyun
	//for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
	//	if v.GWParam.ServiceName == param.Name {
	//		mqttAliyun.ReportServiceParamListAliyun.DeleteReportService(param.Name)
	//
	//		status = true
	//	}
	//}

	//查看Feisjy   gwai add 2023-04-05
	for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
		if v.GWParam.ServiceName == param.Name {
			mqttFeisjy.ReportServiceParamListFeisjy.DeleteReportService(param.Name)

			status = true
		}
	}

	//查看Emqx
	for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
		if v.GWParam.ServiceName == param.Name {
			mqttEmqx.ReportServiceParamListEmqx.DeleteReportService(param.Name)

			status = true
		}
	}

	//查看RT.MQTT
	for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
		if v.GWParam.ServiceName == param.Name {
			mqttRT.ReportServiceParamListRT.DeleteReportService(param.Name)
			status = true
		}
	}

	//查看SagooIOT.MQTT
	for _, v := range mqttSagooIOT.ReportServiceParamListSagooIOT.ServiceList {
		if v.GWParam.ServiceName == param.Name {
			mqttSagooIOT.ReportServiceParamListSagooIOT.DeleteReportService(param.Name)

			status = true
		}
	}

	//查看ModbusTCP
	for _, v := range modbusTCP.ReportServiceMBTCPList.ServiceList {
		if v.GWParam.ServiceName == param.Name {
			modbusTCP.ReportServiceMBTCPList.DeleteReportService(param.Name)

			status = true
		}
	}

	//查看ThingsBoard
	for _, v := range mqttThingsBoard.ReportServiceParamListThingsBoard.ServiceList {
		if v.GWParam.ServiceName == param.Name {
			mqttThingsBoard.ReportServiceParamListThingsBoard.DeleteReportService(param.Name)

			status = true
		}
	}

	//for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {
	//	if v.GWParam.ServiceName == param.Name {
	//		mqttHuawei.ReportServiceParamListHuawei.DeleteReportService(param.Name)
	//
	//		status = true
	//	}
	//}

	if status == true {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "删除上报服务参数成功",
			Data:    "",
		})
		return
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "上报服务名称不存在",
			Data:    "",
		})
		return
	}
}

func ApiAddReportModel(context *gin.Context) {

	modelInfo := &struct {
		Name  string `json:"name"`  // 名称
		Label string `json:"label"` // 标签
		Code  string `json:"code"`  //编码
	}{}

	err := context.BindJSON(&modelInfo)
	if err != nil {
		setting.ZAPS.Error("添加上报模型JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加上报模型JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = reportModel.AddReportModel(modelInfo.Name, modelInfo.Label, modelInfo.Code)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加上报模型成功",
		Data:    "",
	})
}

func ApiDeleteReportModel(context *gin.Context) {

	modelInfo := &struct {
		Name string `json:"name"` // 名称
	}{}

	err := context.BindJSON(&modelInfo)
	if err != nil {
		setting.ZAPS.Error("删除上报模型JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除上报模型JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = reportModel.DeleteReportModel(modelInfo.Name)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除上报模型成功",
		Data:    "",
	})
}

func ApiModifyReportModel(context *gin.Context) {

	modelInfo := &struct {
		Name  string `json:"name"`  // 名称
		Label string `json:"label"` // 标签
		Code  string `json:"code"`  //编码
	}{}

	err := context.BindJSON(&modelInfo)
	if err != nil {
		setting.ZAPS.Error("修改上报模型JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改上报模型JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = reportModel.ModifyReportModel(modelInfo.Name, modelInfo.Label, modelInfo.Code)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改上报模型成功",
		Data:    "",
	})
}

func ApiGetReportModel(context *gin.Context) {

	type ModelTemplate struct {
		Name  string `json:"name"`
		Label string `json:"label"`
		Code  string `json:"code"` //编码
	}

	models := make([]ModelTemplate, 0)
	for _, v := range reportModel.ReportModels {
		rModel := ModelTemplate{
			Name:  v.Name,
			Label: v.Label,
			Code:  v.Code,
		}
		models = append(models, rModel)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取上报模型成功",
		Data:    models,
	})
}

func ApiAddReportModelProperty(context *gin.Context) {

	modelInfo := &struct {
		ModelName string                                  `json:"name"`     // 名称
		Property  reportModel.ReportModelPropertyTemplate `json:"property"` //
	}{}

	err := context.BindJSON(&modelInfo)
	if err != nil {
		setting.ZAPS.Error("增加上报模型属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加上报模型属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = reportModel.AddReportModelProperty(modelInfo.ModelName, &modelInfo.Property)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加上报模型属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加上报模型属性成功",
		Data:    "",
	})
	return
}

func ApiModifyReportModelProperty(context *gin.Context) {

	modelInfo := &struct {
		ModelName string                                  `json:"name"`     // 名称
		Property  reportModel.ReportModelPropertyTemplate `json:"property"` //
	}{}

	err := context.BindJSON(&modelInfo)
	if err != nil {
		setting.ZAPS.Error("修改上报模型属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改上报模型属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = reportModel.ModifyReportModelProperty(modelInfo.ModelName, &modelInfo.Property)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改上报模型属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改上报模型属性成功",
		Data:    "",
	})
}

func ApiDeleteReportModelProperties(context *gin.Context) {

	modelInfo := &struct {
		ModelName     string   `json:"name"`       // 模型名称
		PropertyNames []string `json:"properties"` // 模型属性名称
	}{}

	err := context.BindJSON(&modelInfo)
	if err != nil {
		setting.ZAPS.Error("删除上报模型属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除上报模型属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = reportModel.DeleteReportModelProperties(modelInfo.ModelName, modelInfo.PropertyNames)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除上报模型属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除上报模型属性成功",
		Data:    "",
	})
}

func ApiGetReportModelProperties(context *gin.Context) {

	type ModelPropertyTemplate struct {
		Index      int                                          `json:"index"`
		Name       string                                       `json:"name"` //属性名称，只可以是字母和数字的组合
		UploadName string                                       `json:"uploadName"`
		Label      string                                       `json:"label"`    //属性解释
		Type       int                                          `json:"type"`     //类型 uint32 int32 double string
		Decimals   int                                          `json:"decimals"` //小数位数
		Unit       string                                       `json:"unit"`
		Params     reportModel.ReportModelPropertyParamTemplate `json:"params"`
	}

	modelName := context.Query("name")

	rModel, ok := reportModel.ReportModels[modelName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "上报模型名称不存在",
			Data:    "",
		})
		return
	}

	properties := make([]ModelPropertyTemplate, 0)
	property := ModelPropertyTemplate{}
	for _, v := range rModel.Properties {
		property.Index = v.Index
		property.Name = v.Name
		property.Type = v.Type
		property.Label = v.Label
		property.UploadName = v.UploadName
		property.Decimals = v.Decimals
		property.Unit = v.Unit
		property.Params = v.Params
		properties = append(properties, property)
	}

	//排序，方便前端页面显示
	sort.Slice(properties, func(i, j int) bool {
		iName := properties[i].Index
		jName := properties[j].Index
		return iName > jName
	})

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取上报模型属性成功",
		Data:    properties,
	})
}

func ApiAddReportNodeWParam(context *gin.Context) {

	type ReportServiceNodeTemplate struct {
		ServiceName       string      `json:"serviceName"`
		CollInterfaceName string      `json:"collInterfaceName"`
		DeviceName        string      `json:"deviceName"`
		DeviceAddr        string      `json:"deviceAddr"`
		DeviceLabel       string      `json:"deviceLabel"`
		UploadModel       string      `json:"uploadModel"`
		Protocol          string      `json:"protocol"`
		Param             interface{} `json:"param"`
	}

	var Param json.RawMessage
	param := ReportServiceNodeTemplate{
		Param: &Param,
	}

	//从context中取出缓存，然后再赋值给context，用于第二次校验
	contextRawData, _ := context.GetRawData()
	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(contextRawData))

	err := context.ShouldBindJSON(&param)
	if err != nil {
		setting.ZAPS.Errorf("增加上报服务设备JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加上报服务设备JSON格式化错误",
			Data:    "",
		})
		return
	}

	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(contextRawData))

	switch param.Protocol {
	//case "Aliyun.MQTT":
	//	ReportServiceNodeParamAliyun := mqttAliyun.ReportServiceNodeParamAliyunTemplate{}
	//	err := context.ShouldBindJSON(&ReportServiceNodeParamAliyun)
	//	if err != nil {
	//		setting.ZAPS.Errorf("增加上报服务Aliyun.MQTT设备JSON格式化错误[%v]", err)
	//		context.JSON(http.StatusOK, model.ResponseData{
	//			Code:    "1",
	//			Message: "增加上报服务[Aliyun.MQTT]设备参数JSON格式化错误",
	//			Data:    "",
	//		})
	//		return
	//	}
	//	for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
	//		if v.GWParam.ServiceName == param.ServiceName {
	//			v.AddReportNode(ReportServiceNodeParamAliyun)
	//		}
	//	}
	//	setting.ZAPS.Debugf("ParamListAliyun %v", mqttAliyun.ReportServiceParamListAliyun.ServiceList)
	case "FSJY.MQTT": //gwai add 2023-04-05
		ReportServiceNodeParamFeisjy := mqttFeisjy.ReportServiceNodeParamFeisjyTemplate{}
		err := context.ShouldBindJSON(&ReportServiceNodeParamFeisjy)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务FSJY.MQTT设备JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "增加上报服务[FSJY.MQTT]设备参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				v.AddReportNode(ReportServiceNodeParamFeisjy)
			}
		}
		setting.ZAPS.Debugf("ParamListEmqx %+v", mqttEmqx.ReportServiceParamListEmqx.ServiceList)
	case "EMQX.MQTT":
		ReportServiceNodeParamEmqx := mqttEmqx.ReportServiceNodeParamEmqxTemplate{}
		err := context.ShouldBindJSON(&ReportServiceNodeParamEmqx)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务EMQX.MQTT设备JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "增加上报服务[EMQX.MQTT]设备参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				v.AddReportNode(ReportServiceNodeParamEmqx)
			}
		}
		setting.ZAPS.Debugf("ParamListEmqx %+v", mqttEmqx.ReportServiceParamListEmqx.ServiceList)
	case "RT.MQTT":
		ReportServiceNodeParamRT := mqttRT.ReportServiceNodeParamRTTemplate{}
		err := context.ShouldBindJSON(&ReportServiceNodeParamRT)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务RT.MQTT设备JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "增加上报服务[RT.MQTT]设备参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				_ = v.AddReportNode(ReportServiceNodeParamRT)
			}
		}
		setting.ZAPS.Debugf("ParamListRT %+v", mqttRT.ReportServiceParamListRT.ServiceList)
	case "ThingsBoard.MQTT":
		ReportServiceNodeParamThingsBoard := mqttThingsBoard.ReportServiceNodeParamThingsBoardTemplate{}
		err := context.ShouldBindJSON(&ReportServiceNodeParamThingsBoard)
		if err != nil {
			setting.ZAPS.Errorf("增加上报服务ThingsBoard.MQTT设备JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "增加上报服务[Thingsboard.MQTT]设备参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		for _, v := range mqttThingsBoard.ReportServiceParamListThingsBoard.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				v.AddReportNode(ReportServiceNodeParamThingsBoard)
			}
		}
		setting.ZAPS.Debugf("ParamListThingsBoard %v", mqttThingsBoard.ReportServiceParamListThingsBoard.ServiceList)
	//case "Huawei.MQTT":
	//	ReportServiceNodeParamHuawei := mqttHuawei.ReportServiceNodeParamHuaweiTemplate{}
	//	err := context.ShouldBindJSON(&ReportServiceNodeParamHuawei)
	//	if err != nil {
	//		setting.ZAPS.Errorf("增加上报服务Huawei.MQTT设备JSON格式化错误[%v]", err)
	//		context.JSON(http.StatusOK, model.ResponseData{
	//			Code:    "1",
	//			Message: "增加上报服务[Huawei.MQTT]设备参数JSON格式化错误",
	//			Data:    "",
	//		})
	//		return
	//	}
	//	for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {
	//		if v.GWParam.ServiceName == param.ServiceName {
	//			v.AddReportNode(ReportServiceNodeParamHuawei)
	//		}
	//	}
	default:
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加上报服务设备错误[未知协议]",
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "增加上报服务设备成功",
		Data:    "",
	})
}

func ApiModifyReportNodeWParam(context *gin.Context) {

	type ReportServiceNodeTemplate struct {
		ServiceName       string      `json:"serviceName"`
		CollInterfaceName string      `json:"collInterfaceName"`
		DeviceName        string      `json:"deviceName"`
		DeviceLabel       string      `json:"deviceLabel"`
		DeviceAddr        string      `json:"deviceAddr"`
		UploadModel       string      `json:"uploadModel"`
		Protocol          string      `json:"protocol"`
		Param             interface{} `json:"param"`
	}

	var Param json.RawMessage
	param := ReportServiceNodeTemplate{
		Param: &Param,
	}

	//从context中取出缓存，然后再赋值给context，用于第二次校验
	contextRawData, _ := context.GetRawData()
	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(contextRawData))

	err := context.ShouldBindJSON(&param)
	if err != nil {
		setting.ZAPS.Errorf("修改上报服务设备JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改上报服务设备JSON格式化错误",
			Data:    "",
		})
		return
	}

	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(contextRawData))

	switch param.Protocol {
	//case "Aliyun.MQTT":
	//	ReportServiceNodeParamAliyun := mqttAliyun.ReportServiceNodeParamAliyunTemplate{}
	//	err := context.ShouldBindJSON(&ReportServiceNodeParamAliyun)
	//	if err != nil {
	//		setting.ZAPS.Errorf("修改上报服务Aliyun.MQTT设备JSON格式化错误[%v]", err)
	//		context.JSON(http.StatusOK, model.ResponseData{
	//			Code:    "1",
	//			Message: "修改上报服务[Aliyun.MQTT]设备参数JSON格式化错误",
	//			Data:    "",
	//		})
	//		return
	//	}
	//	for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
	//		if v.GWParam.ServiceName == param.ServiceName {
	//			v.AddReportNode(ReportServiceNodeParamAliyun)
	//		}
	//	}
	//	setting.ZAPS.Debugf("ParamListAliyun %v", mqttAliyun.ReportServiceParamListAliyun.ServiceList)
	case "FSJY.MQTT": //gwai add 2023-04-05
		ReportServiceNodeParamFeisjy := mqttFeisjy.ReportServiceNodeParamFeisjyTemplate{}
		err := context.ShouldBindJSON(&ReportServiceNodeParamFeisjy)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务FSJY.MQTT设备JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "修改上报服务[FSJY.MQTT]设备参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				err := v.ModifyReportNode(ReportServiceNodeParamFeisjy)
				if err != nil {
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: fmt.Sprintf("修改上报服务[FSJY.MQTT]设备错误[%v]", err.Error()),
						Data:    "",
					})
					return
				}
			}
		}
		setting.ZAPS.Debugf("ParamListFeisjy %v", mqttFeisjy.ReportServiceParamListFeisjy.ServiceList)
	case "EMQX.MQTT":
		ReportServiceNodeParamEmqx := mqttEmqx.ReportServiceNodeParamEmqxTemplate{}
		err := context.ShouldBindJSON(&ReportServiceNodeParamEmqx)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务EMQX.MQTT设备JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "修改上报服务[EMQX.MQTT]设备参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				err := v.ModifyReportNode(ReportServiceNodeParamEmqx)
				if err != nil {
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: fmt.Sprintf("修改上报服务[EMQX.MQTT]设备错误[%v]", err.Error()),
						Data:    "",
					})
					return
				}
			}
		}
		setting.ZAPS.Debugf("ParamListEmqx %v", mqttEmqx.ReportServiceParamListEmqx.ServiceList)
	case "RT.MQTT":
		ReportServiceNodeParamRT := mqttRT.ReportServiceNodeParamRTTemplate{}
		err := context.ShouldBindJSON(&ReportServiceNodeParamRT)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务RT.MQTT设备JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "修改上报服务[RT.MQTT]设备参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				err := v.ModifyReportNode(ReportServiceNodeParamRT)
				if err != nil {
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: fmt.Sprintf("修改上报服务[RT.MQTT]设备错误[%v]", err.Error()),
						Data:    "",
					})
					return
				}
			}
		}
		setting.ZAPS.Debugf("ParamListRT %v", mqttRT.ReportServiceParamListRT.ServiceList)
	case "ThingsBoard.MQTT":
		ReportServiceNodeParamThingsBoard := mqttThingsBoard.ReportServiceNodeParamThingsBoardTemplate{}
		err := context.ShouldBindJSON(&ReportServiceNodeParamThingsBoard)
		if err != nil {
			setting.ZAPS.Errorf("修改上报服务ThingsBoard.MQTT设备JSON格式化错误[%v]", err)
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "修改上报服务[Thingsboard.MQTT]设备参数JSON格式化错误",
				Data:    "",
			})
			return
		}
		for _, v := range mqttThingsBoard.ReportServiceParamListThingsBoard.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				v.AddReportNode(ReportServiceNodeParamThingsBoard)
			}
		}
		setting.ZAPS.Debugf("ParamListThingsBoard %v", mqttThingsBoard.ReportServiceParamListThingsBoard.ServiceList)
	//case "Huawei.MQTT":
	//	ReportServiceNodeParamHuawei := mqttHuawei.ReportServiceNodeParamHuaweiTemplate{}
	//	err := context.ShouldBindJSON(&ReportServiceNodeParamHuawei)
	//	if err != nil {
	//		setting.ZAPS.Errorf("修改上报服务Huawei.MQTT设备JSON格式化错误[%v]", err)
	//		context.JSON(http.StatusOK, model.ResponseData{
	//			Code:    "1",
	//			Message: "修改上报服务[Huawei.MQTT]设备参数JSON格式化错误",
	//			Data:    "",
	//		})
	//		return
	//	}
	//	for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {
	//		if v.GWParam.ServiceName == param.ServiceName {
	//			v.AddReportNode(ReportServiceNodeParamHuawei)
	//		}
	//	}
	default:
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改上报服务设备错误[未知协议]",
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改上报服务设备成功",
		Data:    "",
	})
}

func ApiBatchAddReportNodeParam(context *gin.Context) {

	aParam := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	// 获取文件头
	file, err := context.FormFile("file")
	if err != nil {
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + file.Filename

	//保存文件到服务器本地
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		aParam.Code = "1"
		aParam.Message = "save File Error"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	result := setting.LoadCsvCfg(fileName, 1, 2, 0) //标题在第2行，从第3行取数据，第1列取数据
	if result == nil {
		return
	}

	for _, record := range result.Records {
		setting.ZAPS.Debugf("record %v", record)
		protocol := record.GetString("Protocol")
		switch protocol {
		//case "Aliyun.MQTT":
		//	{
		//		ReportServiceNodeParamAliyun := mqttAliyun.ReportServiceNodeParamAliyunTemplate{}
		//		ReportServiceNodeParamAliyun.ServiceName = record.GetString("ServiceName")
		//		ReportServiceNodeParamAliyun.CollInterfaceName = record.GetString("CollInterfaceName")
		//		ReportServiceNodeParamAliyun.Name = record.GetString("Name")
		//		ReportServiceNodeParamAliyun.Addr = record.GetString("Addr")
		//		ReportServiceNodeParamAliyun.Protocol = record.GetString("Protocol")
		//		ReportServiceNodeParamAliyun.Param.ProductKey = record.GetString("ProductKey")
		//		ReportServiceNodeParamAliyun.Param.DeviceName = record.GetString("DeviceName")
		//		ReportServiceNodeParamAliyun.Param.DeviceSecret = record.GetString("DeviceSecret")
		//
		//		for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
		//			if v.GWParam.ServiceName == ReportServiceNodeParamAliyun.ServiceName {
		//				v.AddReportNode(ReportServiceNodeParamAliyun)
		//			}
		//		}
		//	}
		case "FSJY.MQTT":
			{
				ReportServiceNodeParamFeisjy := mqttFeisjy.ReportServiceNodeParamFeisjyTemplate{}
				ReportServiceNodeParamFeisjy.ServiceName = record.GetString("ServiceName")
				ReportServiceNodeParamFeisjy.CollInterfaceName = record.GetString("CollInterfaceName")
				ReportServiceNodeParamFeisjy.Name = record.GetString("Name")
				ReportServiceNodeParamFeisjy.Addr = record.GetString("Addr")
				ReportServiceNodeParamFeisjy.UploadModel = record.GetString("UploadModel")
				ReportServiceNodeParamFeisjy.Protocol = record.GetString("Protocol")

				for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
					if v.GWParam.ServiceName == ReportServiceNodeParamFeisjy.ServiceName {
						v.AddReportNode(ReportServiceNodeParamFeisjy)
					}
				}
			}
		case "EMQX.MQTT":
			{
				ReportServiceNodeParamEmqx := mqttEmqx.ReportServiceNodeParamEmqxTemplate{}
				ReportServiceNodeParamEmqx.ServiceName = record.GetString("ServiceName")
				ReportServiceNodeParamEmqx.CollInterfaceName = record.GetString("CollInterfaceName")
				ReportServiceNodeParamEmqx.Name = record.GetString("Name")
				ReportServiceNodeParamEmqx.Addr = record.GetString("Addr")
				ReportServiceNodeParamEmqx.UploadModel = record.GetString("UploadModel")
				ReportServiceNodeParamEmqx.Protocol = record.GetString("Protocol")
				ReportServiceNodeParamEmqx.Param.DeviceCode = record.GetString("DeviceCode")

				for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
					if v.GWParam.ServiceName == ReportServiceNodeParamEmqx.ServiceName {
						v.AddReportNode(ReportServiceNodeParamEmqx)
					}
				}
			}
		case "RT.MQTT":
			{
				ReportServiceNodeParamRT := mqttRT.ReportServiceNodeParamRTTemplate{}
				ReportServiceNodeParamRT.ServiceName = record.GetString("ServiceName")
				ReportServiceNodeParamRT.CollInterfaceName = record.GetString("CollInterfaceName")
				ReportServiceNodeParamRT.Name = record.GetString("Name")
				ReportServiceNodeParamRT.Label = record.GetString("Label")
				ReportServiceNodeParamRT.Addr = record.GetString("Addr")
				ReportServiceNodeParamRT.UploadModel = record.GetString("UploadModel")
				ReportServiceNodeParamRT.Protocol = record.GetString("Protocol")
				ReportServiceNodeParamRT.Param.DeviceCode = record.GetString("DeviceCode")

				for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
					if v.GWParam.ServiceName == ReportServiceNodeParamRT.ServiceName {
						v.AddReportNode(ReportServiceNodeParamRT)
					}
				}
			}
		}

	}

	aParam.Code = "0"
	aParam.Message = "批量导入设备成功！"

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func ApiBatchAddReportNodeParamFromXlsx(context *gin.Context) {

	aParam := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	// 获取文件头
	file, err := context.FormFile("file")
	if err != nil {
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + file.Filename

	//保存文件到服务器本地
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		aParam.Code = "1"
		aParam.Message = "save File Error"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	err, cells := setting.ReadExcel(fileName) //标题在第2行，从第3行取数据，第2列取数据
	if err != nil {
		aParam.Code = "1"
		aParam.Message = "读取Excel文件错误"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))

		return
	}

	for _, cell := range cells {
		if len(cell) < 8 {
			continue
		}
		protocol := setting.GetString(cell[6])
		switch protocol {
		//case "Aliyun.MQTT":
		//	{
		//		ReportServiceNodeParamAliyun := mqttAliyun.ReportServiceNodeParamAliyunTemplate{}
		//		ReportServiceNodeParamAliyun.ServiceName = record.GetString("ServiceName")
		//		ReportServiceNodeParamAliyun.CollInterfaceName = record.GetString("CollInterfaceName")
		//		ReportServiceNodeParamAliyun.Name = record.GetString("Name")
		//		ReportServiceNodeParamAliyun.Addr = record.GetString("Addr")
		//		ReportServiceNodeParamAliyun.Protocol = record.GetString("Protocol")
		//		ReportServiceNodeParamAliyun.Param.ProductKey = record.GetString("ProductKey")
		//		ReportServiceNodeParamAliyun.Param.DeviceName = record.GetString("DeviceName")
		//		ReportServiceNodeParamAliyun.Param.DeviceSecret = record.GetString("DeviceSecret")
		//
		//		for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
		//			if v.GWParam.ServiceName == ReportServiceNodeParamAliyun.ServiceName {
		//				v.AddReportNode(ReportServiceNodeParamAliyun)
		//			}
		//		}
		//	}
		case "FSJY.MQTT": //gwai add 2023-04-05
			{
				ReportServiceNodeParamFeisjy := mqttFeisjy.ReportServiceNodeParamFeisjyTemplate{}
				ReportServiceNodeParamFeisjy.ServiceName = setting.GetString(cell[0])
				ReportServiceNodeParamFeisjy.CollInterfaceName = setting.GetString(cell[1])
				ReportServiceNodeParamFeisjy.Name = setting.GetString(cell[2])
				ReportServiceNodeParamFeisjy.Addr = setting.GetString(cell[3])
				ReportServiceNodeParamFeisjy.UploadModel = setting.GetString(cell[4])
				ReportServiceNodeParamFeisjy.Protocol = setting.GetString(cell[5])

				for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
					if v.GWParam.ServiceName == ReportServiceNodeParamFeisjy.ServiceName {
						v.AddReportNode(ReportServiceNodeParamFeisjy)
					}
				}
			}
		case "EMQX.MQTT":
			{
				ReportServiceNodeParamEmqx := mqttEmqx.ReportServiceNodeParamEmqxTemplate{}
				ReportServiceNodeParamEmqx.ServiceName = setting.GetString(cell[0])
				ReportServiceNodeParamEmqx.CollInterfaceName = setting.GetString(cell[1])
				ReportServiceNodeParamEmqx.Name = setting.GetString(cell[2])
				ReportServiceNodeParamEmqx.Addr = setting.GetString(cell[3])
				ReportServiceNodeParamEmqx.UploadModel = setting.GetString(cell[4])
				ReportServiceNodeParamEmqx.Protocol = setting.GetString(cell[5])
				ReportServiceNodeParamEmqx.Param.DeviceCode = setting.GetString(cell[6])

				for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
					if v.GWParam.ServiceName == ReportServiceNodeParamEmqx.ServiceName {
						v.AddReportNode(ReportServiceNodeParamEmqx)
					}
				}
			}
		case "RT.MQTT":
			{
				ReportServiceNodeParamRT := mqttRT.ReportServiceNodeParamRTTemplate{}
				ReportServiceNodeParamRT.ServiceName = setting.GetString(cell[0])
				ReportServiceNodeParamRT.CollInterfaceName = setting.GetString(cell[1])
				ReportServiceNodeParamRT.Name = setting.GetString(cell[2])
				ReportServiceNodeParamRT.Label = setting.GetString(cell[3])
				ReportServiceNodeParamRT.Addr = setting.GetString(cell[4])
				ReportServiceNodeParamRT.UploadModel = setting.GetString(cell[5])
				ReportServiceNodeParamRT.Protocol = setting.GetString(cell[6])
				ReportServiceNodeParamRT.Param.DeviceCode = setting.GetString(cell[7])

				for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
					if v.GWParam.ServiceName == ReportServiceNodeParamRT.ServiceName {
						v.AddReportNode(ReportServiceNodeParamRT)
					}
				}
			}
		}

	}

	aParam.Code = "0"
	aParam.Message = "批量导入设备成功！"

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func ApiGetReportNodeWParam(context *gin.Context) {

	type ReportServiceNodeTemplate struct {
		Index             int         `json:"index"`
		ServiceName       string      `json:"serviceName"`
		CollInterfaceName string      `json:"collInterfaceName"`
		DeviceName        string      `json:"deviceName"`
		DeviceLabel       string      `json:"deviceLabel"`
		DeviceAddr        string      `json:"deviceAddr"`
		CommStatus        string      `json:"commStatus"`
		ReportStatus      string      `json:"reportStatus"`
		UploadModel       string      `json:"uploadModel"`
		Protocol          string      `json:"protocol"`
		Param             interface{} `json:"param"`
	}

	ServiceName := context.Query("serviceName")

	nodes := make([]ReportServiceNodeTemplate, 0)
	//for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
	//	if v.GWParam.ServiceName == ServiceName {
	//		ReportServiceNode := ReportServiceNodeTemplate{}
	//		for _, d := range v.NodeList {
	//			ReportServiceNode.ServiceName = d.ServiceName
	//			ReportServiceNode.CollInterfaceName = d.CollInterfaceName
	//			ReportServiceNode.DeviceName = d.Name
	//			ReportServiceNode.DeviceAddr = d.Addr
	//			ReportServiceNode.Protocol = d.Protocol
	//			ReportServiceNode.CommStatus = d.CommStatus
	//			ReportServiceNode.ReportStatus = d.ReportStatus
	//			ReportServiceNode.Param = d.Param
	//			nodes = append(nodes, ReportServiceNode)
	//		}
	//	}
	//}

	//gwai add 2023-04-05
	for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			ReportServiceNode := ReportServiceNodeTemplate{}
			for _, d := range v.NodeList {
				ReportServiceNode.ServiceName = d.ServiceName
				ReportServiceNode.CollInterfaceName = d.CollInterfaceName
				ReportServiceNode.DeviceName = d.Name
				ReportServiceNode.DeviceAddr = d.Addr
				ReportServiceNode.UploadModel = d.UploadModel
				ReportServiceNode.Protocol = d.Protocol
				ReportServiceNode.CommStatus = d.CommStatus
				ReportServiceNode.ReportStatus = d.ReportStatus
				ReportServiceNode.Param = d.Param
				nodes = append(nodes, ReportServiceNode)
			}
		}
	}

	for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			ReportServiceNode := ReportServiceNodeTemplate{}
			for _, d := range v.NodeList {
				ReportServiceNode.ServiceName = d.ServiceName
				ReportServiceNode.CollInterfaceName = d.CollInterfaceName
				ReportServiceNode.DeviceName = d.Name
				ReportServiceNode.DeviceAddr = d.Addr
				ReportServiceNode.UploadModel = d.UploadModel
				ReportServiceNode.Protocol = d.Protocol
				ReportServiceNode.CommStatus = d.CommStatus
				ReportServiceNode.ReportStatus = d.ReportStatus
				ReportServiceNode.Param = d.Param
				nodes = append(nodes, ReportServiceNode)
			}
		}
	}
	for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			ReportServiceNode := ReportServiceNodeTemplate{}
			for _, d := range v.NodeList {
				ReportServiceNode.Index = d.Index
				ReportServiceNode.ServiceName = d.ServiceName
				ReportServiceNode.CollInterfaceName = d.CollInterfaceName
				ReportServiceNode.DeviceName = d.Name
				ReportServiceNode.DeviceAddr = d.Addr
				ReportServiceNode.DeviceLabel = d.Label
				ReportServiceNode.UploadModel = d.UploadModel
				ReportServiceNode.Protocol = d.Protocol
				ReportServiceNode.CommStatus = d.CommStatus
				ReportServiceNode.ReportStatus = d.ReportStatus
				ReportServiceNode.Param = d.Param
				nodes = append(nodes, ReportServiceNode)
			}
		}
	}

	for _, v := range mqttThingsBoard.ReportServiceParamListThingsBoard.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			ReportServiceNode := ReportServiceNodeTemplate{}
			for _, d := range v.NodeList {
				ReportServiceNode.ServiceName = d.ServiceName
				ReportServiceNode.CollInterfaceName = d.CollInterfaceName
				ReportServiceNode.DeviceName = d.Name
				ReportServiceNode.DeviceAddr = d.Addr
				ReportServiceNode.Protocol = d.Protocol
				ReportServiceNode.CommStatus = d.CommStatus
				ReportServiceNode.ReportStatus = d.ReportStatus
				ReportServiceNode.Param = d.Param
				nodes = append(nodes, ReportServiceNode)
			}
		}
	}

	//for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {
	//	if v.GWParam.ServiceName == ServiceName {
	//		ReportServiceNode := ReportServiceNodeTemplate{}
	//		for _, d := range v.NodeList {
	//			ReportServiceNode.ServiceName = d.ServiceName
	//			ReportServiceNode.CollInterfaceName = d.CollInterfaceName
	//			ReportServiceNode.DeviceName = d.Name
	//			ReportServiceNode.DeviceAddr = d.Addr
	//			ReportServiceNode.Protocol = d.Protocol
	//			ReportServiceNode.CommStatus = d.CommStatus
	//			ReportServiceNode.ReportStatus = d.ReportStatus
	//			ReportServiceNode.Param = d.Param
	//			nodes = append(nodes, ReportServiceNode)
	//		}
	//	}
	//}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取上报服务设备成功",
		Data:    nodes,
	})
}

func ApiBatchExportReportNodeWParam(context *gin.Context) {

	ServiceName := context.Query("serviceName")

	for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			status, name := v.ExportParamToCsv()
			if status == true {
				//返回文件流
				context.Writer.Header().Add("Content-Disposition",
					fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(name))))
				context.File(name) //返回文件路径，自动调用http.ServeFile方法
				return
			} else {
				aParam := struct {
					Code    string
					Message string
					Data    string
				}{
					Code:    "1",
					Message: "",
					Data:    "",
				}
				sJson, _ := json.Marshal(aParam)
				context.String(http.StatusOK, string(sJson))
				return
			}
		}
	}

	for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			status, name := v.ExportParamToCsv()
			if status == true {
				//返回文件流
				context.Writer.Header().Add("Content-Disposition",
					fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(name))))
				context.File(name) //返回文件路径，自动调用http.ServeFile方法
				return
			} else {
				aParam := struct {
					Code    string
					Message string
					Data    string
				}{
					Code:    "1",
					Message: "",
					Data:    "",
				}
				sJson, _ := json.Marshal(aParam)
				context.String(http.StatusOK, string(sJson))
				return
			}
		}
	}

	aParam := struct {
		Code    string
		Message string
		Data    string
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func ApiBatchExportReportNodeWParamToXlsx(context *gin.Context) {

	ServiceName := context.Query("serviceName")
	setting.ZAPS.Debugf("ServiceName %v", ServiceName)

	//gwai add 2023-04-05
	for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			status, name := v.ExportParamToXlsx()
			if status == true {
				//返回文件流
				context.Writer.Header().Add("Content-Disposition",
					fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(name))))
				context.File(name) //返回文件路径，自动调用http.ServeFile方法
				return
			} else {
				aParam := struct {
					Code    string
					Message string
					Data    string
				}{
					Code:    "1",
					Message: "",
					Data:    "",
				}
				sJson, _ := json.Marshal(aParam)
				context.String(http.StatusOK, string(sJson))
				return
			}
		}
	}

	for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			status, name := v.ExportParamToXlsx()
			if status == true {
				//返回文件流
				context.Writer.Header().Add("Content-Disposition",
					fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(name))))
				context.File(name) //返回文件路径，自动调用http.ServeFile方法
				return
			} else {
				aParam := struct {
					Code    string
					Message string
					Data    string
				}{
					Code:    "1",
					Message: "",
					Data:    "",
				}
				sJson, _ := json.Marshal(aParam)
				context.String(http.StatusOK, string(sJson))
				return
			}
		}
	}
	for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			status, name := v.ExportParamToXlsx()
			if status == true {
				//返回文件流
				context.Writer.Header().Add("Content-Disposition",
					fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(name))))
				context.File(name) //返回文件路径，自动调用http.ServeFile方法
				return
			} else {
				aParam := struct {
					Code    string
					Message string
					Data    string
				}{
					Code:    "1",
					Message: "xlsx导入失败",
					Data:    "",
				}
				sJson, _ := json.Marshal(aParam)
				context.String(http.StatusOK, string(sJson))
				return
			}
		}
	}

	aParam := struct {
		Code    string
		Message string
		Data    string
	}{
		Code:    "1",
		Message: "上报服务名称不存在",
		Data:    "",
	}
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func ApiDeleteReportNodeWParam(context *gin.Context) {

	param := struct {
		ServiceName string   `json:"serviceName"`
		DeviceNames []string `json:"deviceNames"`
	}{}
	err := context.ShouldBindJSON(&param)
	if err != nil {
		setting.ZAPS.Errorf("删除上报服务设备参数JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除上报服务设备参数JSON格式化错误",
			Data:    "",
		})
		return
	}
	// 定义删除失败列表
	var errDeviceList []string
	var index int
	//查看Aliyun
	//for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
	//	for _, n := range v.NodeList {
	//		if n.ServiceName == param.ServiceName {
	//			for _, name := range param.DeviceNames {
	//				if name == n.Name {
	//					index = v.DeleteReportNode(n.Name)
	//					if index == -1 {
	//						errDeviceList = append(errDeviceList, name)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}

	//查看Feisjy  gwai add 2023-04-05
	for _, name := range param.DeviceNames {
		for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
			for _, n := range v.NodeList {
				if n.ServiceName == param.ServiceName {
					if name == n.Name {
						index = v.DeleteReportNode(n.Name)
						if index == -1 {
							errDeviceList = append(errDeviceList, name)
						}
					}
				}
			}
		}
	}

	//查看Emqx
	for _, name := range param.DeviceNames {
		for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
			for _, n := range v.NodeList {
				if n.ServiceName == param.ServiceName {
					if name == n.Name {
						index = v.DeleteReportNode(n.Name)
						if index == -1 {
							errDeviceList = append(errDeviceList, name)
						}
					}
				}
			}
		}
	}
	//查看RT.MQTT
	for _, name := range param.DeviceNames {
		for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
			for _, n := range v.NodeList {
				if n.ServiceName == param.ServiceName {
					if name == n.Name {
						index = v.DeleteReportNode(n.Name)
						if index == -1 {
							errDeviceList = append(errDeviceList, name)
						}
					}
				}
			}
		}
	}

	//查看ThingsBoard
	for _, v := range mqttThingsBoard.ReportServiceParamListThingsBoard.ServiceList {
		for _, n := range v.NodeList {
			if n.ServiceName == param.ServiceName {
				for _, name := range param.DeviceNames {
					if name == n.Name {
						v.DeleteReportNode(n.Name)
						index = v.DeleteReportNode(n.Name)
						if index == -1 {
							errDeviceList = append(errDeviceList, name)
						}
					}
				}
			}
		}
	}

	//for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {
	//	for _, n := range v.NodeList {
	//		if n.ServiceName == param.ServiceName {
	//			for _, name := range param.DeviceNames {
	//				if name == n.Name {
	//					v.DeleteReportNode(n.Name)
	//					index = v.DeleteReportNode(n.Name)
	//					if index == -1 {
	//						errDeviceList = append(errDeviceList, name)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}
	if len(errDeviceList) == 0 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "删除上报服务设备成功",
			Data:    errDeviceList,
		})
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除上报服务设备失败",
			Data:    errDeviceList,
		})
	}

}

func ApiSetReportNodeReport(context *gin.Context) {

	param := struct {
		ServiceName string   `json:"serviceName"`
		DeviceNames []string `json:"deviceNames"`
	}{}
	err := context.ShouldBindJSON(&param)
	if err != nil {
		setting.ZAPS.Errorf("删除上报服务设备参数JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除上报服务设备参数JSON格式化错误",
			Data:    "",
		})
		return
	}

	status := false

	//查看Feisjy gwai add 2023-04-05
	for _, name := range param.DeviceNames {
		for _, v := range mqttFeisjy.ReportServiceParamListFeisjy.ServiceList {
			for _, n := range v.NodeList {
				if n.ServiceName == param.ServiceName {
					if name == n.Name {
						//v.ReportNode(n.Name)   //gwai
						status = true

					}
				}
			}
		}
	}

	//查看Emqx
	for _, name := range param.DeviceNames {
		for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
			for _, n := range v.NodeList {
				if n.ServiceName == param.ServiceName {
					if name == n.Name {
						v.ReportNode(n.Name)
						status = true

					}
				}
			}
		}
	}

	//查看RT.MQTT
	for _, name := range param.DeviceNames {
		for _, v := range mqttRT.ReportServiceParamListRT.ServiceList {
			for _, n := range v.NodeList {
				if n.ServiceName == param.ServiceName {
					if name == n.Name {
						v.ReportNodeNoCheck(n.Name)
						status = true
					}
				}
			}
		}
	}
	if status == true {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "手动设置设备立即上报成功",
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "手动设置设备立即上报失败",
			Data:    "",
		})
	}

}
