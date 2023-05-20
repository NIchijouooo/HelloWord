package contorl

import (
	"fmt"
	"gateway/httpServer/model"
	"gateway/report/modbusTCP"
	"gateway/setting"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func ApiAddReportMBTCPRegister(context *gin.Context) {
	regInfo := &struct {
		ServiceName string                          `json:"serviceName"` // 名称
		Register    modbusTCP.MBTCPRegisterTemplate `json:"register"`    //
	}{}

	err := context.BindJSON(&regInfo)
	if err != nil {
		setting.ZAPS.Error("增加ModbusTCP上报服务寄存器JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加ModbusTCP上报服务寄存器JSON格式化错误",
			Data:    "",
		})
		return
	}

	index := -1
	for k, v := range modbusTCP.ReportServiceMBTCPList.ServiceList {
		if v.GWParam.ServiceName == regInfo.ServiceName {
			index = k
		}
	}
	if index == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "上报服务名称不存在",
			Data:    "",
		})
		return
	}

	err = modbusTCP.ReportServiceMBTCPList.ServiceList[index].GWParam.AddRegister(regInfo.Register)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加ModbusTCP上报服务寄存器错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加ModbusTCP上报服务寄存器成功",
		Data:    "",
	})
	return
}

func ApiModifyReportMBTCPRegister(context *gin.Context) {
	regInfo := &struct {
		ServiceName string                          `json:"serviceName"` // 名称
		Register    modbusTCP.MBTCPRegisterTemplate `json:"register"`    //
	}{}

	err := context.BindJSON(&regInfo)
	if err != nil {
		setting.ZAPS.Error("修改ModbusTCP上报服务寄存器JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改ModbusTCP上报服务寄存器JSON格式化错误",
			Data:    "",
		})
		return
	}

	index := -1
	for k, v := range modbusTCP.ReportServiceMBTCPList.ServiceList {
		if v.GWParam.ServiceName == regInfo.ServiceName {
			index = k
		}
	}
	if index == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "上报服务名称不存在",
			Data:    "",
		})
		return
	}

	err = modbusTCP.ReportServiceMBTCPList.ServiceList[index].GWParam.ModifyRegister(regInfo.Register)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改ModbusTCP上报服务寄存器错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改ModbusTCP上报服务寄存器成功",
		Data:    "",
	})
	return
}

func ApiDeleteReportMBTCPRegisters(context *gin.Context) {
	regInfo := &struct {
		ServiceName string   `json:"serviceName"` // 名称
		Register    []string `json:"registers"`   //
	}{}

	err := context.BindJSON(&regInfo)
	if err != nil {
		setting.ZAPS.Error("删除ModbusTCP上报服务寄存器JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除ModbusTCP上报服务寄存器JSON格式化错误",
			Data:    "",
		})
		return
	}

	index := -1
	for k, v := range modbusTCP.ReportServiceMBTCPList.ServiceList {
		if v.GWParam.ServiceName == regInfo.ServiceName {
			index = k
		}
	}
	if index == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "上报服务名称不存在",
			Data:    "",
		})
		return
	}

	err = modbusTCP.ReportServiceMBTCPList.ServiceList[index].GWParam.DeleteRegisters(regInfo.Register)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除ModbusTCP上报服务寄存器错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除ModbusTCP上报服务寄存器成功",
		Data:    "",
	})
	return
}

func ApiGetReportMBTCPRegisters(context *gin.Context) {

	serviceName := context.Query("serviceName")

	index := -1
	for k, v := range modbusTCP.ReportServiceMBTCPList.ServiceList {
		if v.GWParam.ServiceName == serviceName {
			index = k
		}
	}
	if index == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "上报服务名称不存在",
			Data:    "",
		})
		return
	}

	regs := modbusTCP.ReportServiceMBTCPList.ServiceList[index].GWParam.GetRegisters()

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取ModbusTCP上报服务寄存器成功",
		Data:    regs,
	})
	return
}

func ApiBatchExportReportMBTCPRegistersToXlsx(context *gin.Context) {
	serviceName := context.Query("serviceName")

	index := -1
	for k, v := range modbusTCP.ReportServiceMBTCPList.ServiceList {
		if v.GWParam.ServiceName == serviceName {
			index = k
		}
	}
	if index == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "上报服务名称不存在",
			Data:    "",
		})
		return
	}

	status, name := modbusTCP.ReportServiceMBTCPList.ServiceList[index].GWParam.ExportRegistersToXlsx()
	if status == true {
		//返回文件流
		context.Writer.Header().Add("Content-Disposition",
			fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(name))))
		context.File(name) //返回文件路径，自动调用http.ServeFile方法
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "导出失败",
			Data:    "",
		})
	}
}

func ApiBatchImportReportMBTCPRegistersToXlsx(context *gin.Context) {

	//获取上报模型名称
	serviceName := context.PostForm("serviceName")

	index := -1
	for k, v := range modbusTCP.ReportServiceMBTCPList.ServiceList {
		if v.GWParam.ServiceName == serviceName {
			index = k
		}
	}
	if index == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "上报服务名称不存在",
			Data:    "",
		})
		return
	}

	// 获取文件头
	file, err := context.FormFile("fileName")
	if err != nil {
		setting.ZAPS.Errorf("modbus采集模型待导入xlsx文件不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "modbus采集模型待导入xlsx文件不存在",
			Data:    "",
		})
		return
	}

	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + file.Filename

	//保存文件到服务器本地
	err = utils.FileCreate(fileName)
	if err != nil {
		setting.ZAPS.Errorf("创建xlsx文件错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "创建xlsx文件错误",
			Data:    "",
		})
		return
	}
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		setting.ZAPS.Errorf("保存xlsx文件错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存xlsx文件错误",
			Data:    "",
		})
		return
	}

	defer os.Remove(fileName)

	err, cells := setting.ReadExcel(fileName) //标题在第2行，从第3行取数据，第2列取数据
	if err != nil {
		setting.ZAPS.Errorf("加载xlsx文件错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "加载xlsx文件错误",
			Data:    "",
		})
		return
	}

	setting.ZAPS.Debugf("cells %v", cells)
	for _, cell := range cells {
		if len(cell) < 10 {
			continue
		}

		register := modbusTCP.MBTCPRegisterTemplate{
			RegName:      setting.GetString(cell[0]),
			Label:        setting.GetString(cell[1]),
			CollName:     setting.GetString(cell[3]),
			NodeName:     setting.GetString(cell[4]),
			PropertyName: setting.GetString(cell[5]),
			RegAddr:      setting.GetInt(cell[6]),
			RegCnt:       setting.GetInt(cell[7]),
			Rule:         setting.GetString(cell[9]),
		}

		if setting.GetString(cell[2]) == "uint32" {
			register.PropertyType = 0
		} else if setting.GetString(cell[2]) == "int32" {
			register.PropertyType = 1
		} else if setting.GetString(cell[2]) == "double" {
			register.PropertyType = 2
		} else {
			register.PropertyType = 2
		}

		if setting.GetString(cell[8]) == "coilStatus" {
			register.RegType = 0
		} else if setting.GetString(cell[8]) == "inputStatus" {
			register.RegType = 1
		} else if setting.GetString(cell[8]) == "holdingRegister" {
			register.RegType = 2
		} else if setting.GetString(cell[8]) == "inputRegister" {
			register.RegType = 3
		} else {
			register.RegType = 3
		}

		err = modbusTCP.ReportServiceMBTCPList.ServiceList[index].GWParam.AddRegister(register)
		if err != nil {
			continue
		}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "modbusTCP上报服务导入xlsx成功",
		Data:    "",
	})
}
