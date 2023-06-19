package contorl

import (
	"encoding/json"
	"fmt"
	"gateway/controllers"
	"gateway/device"
	"gateway/httpServer/model"
	"gateway/setting"
	"gateway/utils"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ApiAddDeviceTSL(context *gin.Context) {

	tslInfo := &struct {
		Name  string `json:"name"`  // 名称
		Label string `json:"label"` // 标签
		Type  int    `json:"type"`  //类型
		//Plugin device.TSLLuaPluginTemplate `json:"plugin"` // 插件
	}{}

	emController := controllers.NewEMController()
	emController.AddEmDeviceModel(context)

	err := context.ShouldBindBodyWith(&tslInfo, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("添加物模型模版JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加物模型模版JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = device.AddTSLMode(tslInfo.Name, tslInfo.Label, tslInfo.Type)
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
		Message: "添加物模型模版成功",
		Data:    "",
	})
}

func ApiDeleteDeviceTSL(context *gin.Context) {

	tsl := &struct {
		Name string `json:"name"` // 名称
		Type int    `json:"type"` //类型
	}{}

	emController := controllers.NewEMController()
	emController.DeleteEmDeviceModel(context)

	err := context.ShouldBindBodyWith(&tsl, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("删除物模型模版JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除物模型模版JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = device.DeleteTSLMode(tsl.Name, tsl.Type)
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
		Message: "删除物模型模版成功",
		Data:    "",
	})
}

func ApiModifyDeviceTSL(context *gin.Context) {

	param := json.RawMessage{}
	tslInfo := &struct {
		Name  string           `json:"name"`  // 名称
		Label string           `json:"label"` // 标签
		Type  int              `json:"type"`  //类型
		Param *json.RawMessage `json:"param"`
	}{
		Param: &param,
	}

	emController := controllers.NewEMController()
	emController.UpdateEmDeviceModel(context)

	err := context.ShouldBindBodyWith(&tslInfo, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("修改物模型模版JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改物模型模版JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = device.ModifyTSLMode(tslInfo.Name, tslInfo.Label, tslInfo.Type, param)
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
		Message: "修改物模型模版成功",
		Data:    "",
	})
}

func ApiGetDeviceTSLs(context *gin.Context) {
	type TSLModelTemplate struct {
		Index int         `json:"-"`
		Name  string      `json:"name"`
		Label string      `json:"label"`
		Type  int         `json:"type"`
		Param interface{} `json:"param"`
	}

	tslModels := make([]TSLModelTemplate, 0)

	tslTypeStr := context.Query("type")
	if tslTypeStr == "" {
		for _, v := range device.TSLModels {
			tslModel := TSLModelTemplate{
				Index: v.GetTSLModelIndex(),
				Name:  v.GetTSLModelName(),
				Label: v.GetTSLModelLabel(),
				Type:  v.GetTSLModelType(),
				Param: v.GetTSLModelParam(),
			}
			tslModels = append(tslModels, tslModel)
		}
		for _, v := range device.TSLModbusMap {
			tslModel := TSLModelTemplate{
				Index: v.GetTSLModelIndex(),
				Name:  v.GetTSLModelName(),
				Label: v.GetTSLModelLabel(),
				Type:  v.GetTSLModelType(),
				Param: v.GetTSLModelParam(),
			}
			tslModels = append(tslModels, tslModel)
		}
		for _, v := range device.TSLDLT6452007Map {
			tslModel := TSLModelTemplate{
				Index: v.GetTSLModelIndex(),
				Name:  v.GetTSLModelName(),
				Label: v.GetTSLModelLabel(),
				Type:  v.GetTSLModelType(),
				Param: v.GetTSLModelParam(),
			}
			tslModels = append(tslModels, tslModel)
		}
	} else if tslTypeStr == "2" {
		tslType, err := strconv.Atoi(tslTypeStr)
		if err != nil {
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "获取物模型模板类型错误",
				Data:    tslModels,
			})
			return
		}
		for _, v := range device.TSLModbusMap {
			if v.GetTSLModelType() != tslType {
				continue
			}
			tslModel := TSLModelTemplate{
				Index: v.GetTSLModelIndex(),
				Name:  v.GetTSLModelName(),
				Label: v.GetTSLModelLabel(),
				Type:  v.GetTSLModelType(),
				Param: v.GetTSLModelParam(),
			}
			tslModels = append(tslModels, tslModel)
		}
	} else if tslTypeStr == "3" {
		tslType, err := strconv.Atoi(tslTypeStr)
		if err != nil {
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "获取物模型模板类型错误",
				Data:    tslModels,
			})
			return
		}
		for _, v := range device.TSLDLT6452007Map {

			if v.GetTSLModelType() != tslType {
				continue
			}
			tslModel := TSLModelTemplate{
				Index: v.GetTSLModelIndex(),
				Name:  v.GetTSLModelName(),
				Label: v.GetTSLModelLabel(),
				Type:  v.GetTSLModelType(),
				Param: v.GetTSLModelParam(),
			}
			tslModels = append(tslModels, tslModel)
		}
	} else {
		tslType, err := strconv.Atoi(tslTypeStr)
		if err != nil {
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "获取物模型模板类型错误",
				Data:    tslModels,
			})
			return
		}
		for _, v := range device.TSLModels {
			if v.GetTSLModelType() != tslType {
				continue
			}
			tslModel := TSLModelTemplate{
				Index: v.GetTSLModelIndex(),
				Name:  v.GetTSLModelName(),
				Label: v.GetTSLModelLabel(),
				Type:  v.GetTSLModelType(),
				Param: v.GetTSLModelParam(),
			}
			tslModels = append(tslModels, tslModel)
		}
	}

	//排序，方便前端页面显示
	sort.Slice(tslModels, func(i, j int) bool {
		iName := tslModels[i].Index
		jName := tslModels[j].Index
		return iName > jName
	})

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取物模型模版成功",
		Data:    tslModels,
	})
}

func ApiGetDeviceTSL(context *gin.Context) {

	type TSLModelTemplate struct {
		Index int         `json:"-"`
		Name  string      `json:"name"`
		Label string      `json:"label"`
		Type  int         `json:"type"`
		Param interface{} `json:"param"`
	}

	tslModels := make([]TSLModelTemplate, 0)

	tslTypeStr := context.Query("type")
	tslType, err := strconv.Atoi(tslTypeStr)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取物模型模板类型错误",
			Data:    tslModels,
		})
		return
	}

	for _, v := range device.TSLModels {
		if v.GetTSLModelType() != tslType {
			continue
		}
		tslModel := TSLModelTemplate{
			Index: v.GetTSLModelIndex(),
			Name:  v.GetTSLModelName(),
			Label: v.GetTSLModelLabel(),
			Type:  v.GetTSLModelType(),
			Param: v.GetTSLModelParam(),
		}
		tslModels = append(tslModels, tslModel)
	}

	//排序，方便前端页面显示
	sort.Slice(tslModels, func(i, j int) bool {
		iName := tslModels[i].Index
		jName := tslModels[j].Index
		return iName > jName
	})

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取物模型模版成功",
		Data:    tslModels,
	})
}

func ApiImportTSLLuaPlugin(context *gin.Context) {

	//获取文件头
	file, err := context.FormFile("fileName")
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取文件头失败",
			Data:    "",
		})
		return
	}

	fileFullName := "./plugin/" + file.Filename
	fileName := strings.TrimSuffix(file.Filename, ".zip")

	//判断plugin中是否存在文件夹，不存在就创建
	utils.DirIsExist("./plugin/" + fileName)

	//保存文件到本地
	if err := context.SaveUploadedFile(file, fileFullName); err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存文件到" + fileFullName + "失败",
			Data:    "",
		})
		return
	}

	err = device.ImportTSLPlugin(file.Filename)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "添加plugin文件成功",
			Data:    "",
		})
	}
}

func ApiGetTSLLuaPluginParam(context *gin.Context) {

	pluginName := context.Query("name")

	pluginParam, err := device.GetTSLPluginParam(pluginName)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取失败" + err.Error(),
			Data:    pluginParam,
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取成功",
		Data:    pluginParam,
	})
}

func ApiExportTSLLuaPlugin(context *gin.Context) {

	tslName := context.Query("name")

	status, name := device.ExportTSLPlugin(tslName)
	if status == true {
		//返回文件流
		context.Writer.Header().Add("Content-Disposition",
			fmt.Sprintf("attachment;filename=%s.zip", tslName))
		context.File(name) //返回文件路径，自动调用http.ServeFile方法

	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
	}
}

func ApiAddTSLProperty(context *gin.Context) {

	properties := json.RawMessage{}
	tslInfo := &struct {
		TSLName  string           `json:"name"` // 名称
		TSLType  int              `json:"type"`
		Property *json.RawMessage `json:"property"` //
	}{
		Property: &properties,
	}

	emController := controllers.NewEMController()
	emController.AddEmDevicePlcModelCmd(context)

	err := context.ShouldBindBodyWith(&tslInfo, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("增加物模型属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加物模型属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLModels[tslInfo.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加物模型属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	err = tslModel.TSLModelPropertiesAdd(properties)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加物模型属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加物模型属性成功",
		Data:    "",
	})

}

func ApiModifyTSLProperty(context *gin.Context) {

	property := json.RawMessage{}
	tslInfo := &struct {
		TSLName  string           `json:"name"`     // 名称
		Property *json.RawMessage `json:"property"` //
	}{
		Property: &property,
	}

	emController := controllers.NewEMController()
	emController.UpdateEmDevicePlcModelCmd(context)

	err := context.ShouldBindBodyWith(&tslInfo, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("修改物模型属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改物模型属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLModels[tslInfo.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改物模型属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	err = tslModel.TSLModelPropertiesModify(property)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改物模型属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改物模型属性成功",
		Data:    "",
	})
}

func ApiDeleteTSLProperties(context *gin.Context) {

	tslInfo := &struct {
		TSLName    string   `json:"name"`       // 名称
		Properties []string `json:"properties"` //
	}{}

	emController := controllers.NewEMController()
	emController.DeleteEmDevicePlcModelCmd(context)

	err := context.ShouldBindBodyWith(&tslInfo, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("删除物模型属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除物模型属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLModels[tslInfo.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除物模型属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	err = tslModel.TSLModelPropertiesDelete(tslInfo.Properties)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除物模型属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除物模型属性成功",
		Data:    "",
	})
}

func ApiGetTSLProperties(context *gin.Context) {

	tslName := context.Query("name")

	for _, v := range device.TSLModels {
		if v.GetTSLModelName() == tslName {
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "获取采集模型属性成功",
				Data:    v.TSLModelPropertiesGet(),
			})
			return
		}
	}

	for _, v := range device.TSLModbusMap {
		if v.GetTSLModelName() == tslName {
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "获取物模型属性成功",
				Data:    v.TSLModelPropertiesGet(),
			})
			return
		}
	}

	for _, v := range device.TSLDLT6452007Map {
		if v.GetTSLModelName() == tslName {
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "获取物模型属性成功",
				Data:    v.TSLModelPropertiesGet(),
			})
			return
		}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "1",
		Message: "获取采集模型类型不支持",
		Data:    "",
	})
}

func ApiAddTSLService(context *gin.Context) {

	tslInfo := &struct {
		TSLName string                       `json:"name"`    // 名称
		Service device.TSLLuaServiceTemplate `json:"service"` //
	}{}

	err := context.BindJSON(&tslInfo)
	if err != nil {
		setting.ZAPS.Error("增加物模型服务JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加物模型服务JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = device.AddTSLService(tslInfo.TSLName, tslInfo.Service)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加物模型服务错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加物模型服务成功",
		Data:    "",
	})
	return
}

func ApiModifyTSLService(context *gin.Context) {
	tslInfo := &struct {
		TSLName string                       `json:"name"`    // 名称
		Service device.TSLLuaServiceTemplate `json:"service"` //
	}{}

	err := context.BindJSON(&tslInfo)
	if err != nil {
		setting.ZAPS.Error("修改物模型服务JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改物模型服务JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = device.ModifyTSLService(tslInfo.TSLName, tslInfo.Service)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改物模型服务错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改物模型服务成功",
		Data:    "",
	})
	return
}

func ApiDeleteTSLServices(context *gin.Context) {

	tslInfo := &struct {
		TSLName      string   `json:"name"`         // 名称
		ServiceNames []string `json:"serviceNames"` //
	}{}

	err := context.BindJSON(&tslInfo)
	if err != nil {
		setting.ZAPS.Error("删除物模型服务JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除物模型服务JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = device.DeleteTSLServices(tslInfo.TSLName, tslInfo.ServiceNames)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除物模型服务错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除物模型服务成功",
		Data:    "",
	})
}

func ApiAddTSLModbusCmd(context *gin.Context) {

	cmdParam := struct {
		TSLName      string `json:"tslName"`
		Name         string `json:"name"` // 名称
		Label        string `json:"label"`
		FunCode      int    `json:"funCode"`
		StartRegAddr int    `json:"startRegAddr"`
		RegCnt       int    `json:"regCnt"`
	}{}

	emController := controllers.NewEMController()
	emController.AddEmDeviceModelCmd(context)

	err := context.ShouldBindBodyWith(&cmdParam, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("增加采集模型[Modbus]命令JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加采集模型[Modbus]命令JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLModbusMap[cmdParam.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[Modbus]命令错误" + "命令名称不存在",
			Data:    "",
		})
		return
	}

	cmd := device.TSLModbusCmdTemplate{
		Name:         cmdParam.Name,
		Label:        cmdParam.Label,
		FunCode:      cmdParam.FunCode,
		StartRegAddr: cmdParam.StartRegAddr,
		RegCnt:       cmdParam.RegCnt,
	}
	err = tslModel.TSLModelCmdAdd(cmd)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[Modbus]命令参数错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加采集模型[Modbus]命令参数成功",
		Data:    "",
	})

}

func ApiAddTSLD07Cmd(context *gin.Context) {

	cmdParam := struct {
		TSLName      string `json:"tslName"`
		Name         string `json:"name"` // 名称
		Label        string `json:"label"`
		BlockRulerId string `json:"blockRulerId"`
		BlockRead    byte   `json:"blockRead"`
	}{}

	emController := controllers.NewEMController()
	emController.AddEmDeviceModelCmd(context)

	err := context.ShouldBindBodyWith(&cmdParam, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("增加采集模型[DLT645-2007]命令JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加采集模型[DLT645-2007]命令JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLDLT6452007Map[cmdParam.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[DLT645-2007]命令错误" + "命令名称不存在",
			Data:    "",
		})
		return
	}

	cmd := device.TSLDLT6452007CmdTemplate{
		Name:         cmdParam.Name,
		Label:        cmdParam.Label,
		BlockRulerId: cmdParam.BlockRulerId,
		BlockRead:    cmdParam.BlockRead,
	}
	err = tslModel.TSLModelCmdAdd(cmd)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[DLT645-2007]命令参数错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加采集模型[DLT645-2007]命令参数成功",
		Data:    "",
	})

}

func ApiModifyTSLModbusCmd(context *gin.Context) {

	cmdParam := struct {
		TSLName      string `json:"tslName"`
		Name         string `json:"name"` // 名称
		Label        string `json:"label"`
		FunCode      int    `json:"funCode"`
		StartRegAddr int    `json:"startRegAddr"`
		RegCnt       int    `json:"regCnt"`
	}{}

	emController := controllers.NewEMController()
	emController.UpdateEmDeviceModelCmd(context)

	err := context.ShouldBindBodyWith(&cmdParam, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("修改采集模型[Modbus]命令JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[Modbus]命令JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLModbusMap[cmdParam.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[Modbus]命令错误，" + "模型名称不存在",
			Data:    "",
		})
		return
	}

	cmd := device.TSLModbusCmdTemplate{
		Name:         cmdParam.Name,
		Label:        cmdParam.Label,
		FunCode:      cmdParam.FunCode,
		StartRegAddr: cmdParam.StartRegAddr,
		RegCnt:       cmdParam.RegCnt,
	}
	err = tslModel.TSLModelCmdModify(cmd)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[Modbus]命令参数错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改采集模型[Modbus]命令参数成功",
		Data:    "",
	})
}
func ApiModifyTSLD07Cmd(context *gin.Context) {

	cmdParam := struct {
		TSLName      string `json:"tslName"`
		Name         string `json:"name"` // 名称
		Label        string `json:"label"`
		BlockRulerId string `json:"blockRulerId"`
		BlockRead    byte   `json:"blockRead"`
	}{}

	emController := controllers.NewEMController()
	emController.UpdateEmDeviceModelCmd(context)

	err := context.ShouldBindBodyWith(&cmdParam, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("修改采集模型[DLT645-2007]命令JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[DLT645-2007]命令JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLDLT6452007Map[cmdParam.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[DLT645-2007]命令错误，" + "模型名称不存在",
			Data:    "",
		})
		return
	}

	cmd := device.TSLDLT6452007CmdTemplate{
		Name:         cmdParam.Name,
		Label:        cmdParam.Label,
		BlockRulerId: cmdParam.BlockRulerId,
		BlockRead:    cmdParam.BlockRead,
	}
	err = tslModel.TSLModelCmdModify(cmd)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[DLT645-2007]命令参数错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改采集模型[DLT645-2007]命令参数成功",
		Data:    "",
	})
}

func ApiDeleteTSLModbusCmd(context *gin.Context) {

	cmdParam := struct {
		TSLName  string   `json:"tslName"` // 名称
		CmdNames []string `json:"names"`   //
	}{}

	emController := controllers.NewEMController()
	emController.DeleteEmDeviceModelCmd(context)

	err := context.ShouldBindBodyWith(&cmdParam, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("删除物模型属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除物模型属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLModbusMap[cmdParam.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[Modbus]命令错误，" + "模型名称不存在",
			Data:    "",
		})
		return
	}

	err = tslModel.TSLModelCmdDelete(cmdParam.CmdNames)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集模型[Modbus]命令错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除采集模型[Modbus]命令成功",
		Data:    "",
	})
}

func ApiDeleteTSLD07Cmd(context *gin.Context) {

	cmdParam := struct {
		TSLName  string   `json:"tslName"` // 名称
		CmdNames []string `json:"names"`   //
	}{}

	emController := controllers.NewEMController()
	emController.DeleteEmDeviceModelCmd(context)

	err := context.ShouldBindBodyWith(&cmdParam, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("删除物模型属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除物模型属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLDLT6452007Map[cmdParam.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[DLT645-2007]命令错误，" + "模型名称不存在",
			Data:    "",
		})
		return
	}

	err = tslModel.TSLModelCmdDelete(cmdParam.CmdNames)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集模型[DLT645-2007]命令错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除采集模型[DLT645-2007]命令成功",
		Data:    "",
	})
}

func ApiGetTSLModbusCmd(context *gin.Context) {

	tslName := context.Query("tslName")

	tslModel, ok := device.TSLModbusMap[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集模型[Modbus]命令错误," + "模型名称不存在",
			Data:    "",
		})
		return
	}

	cmds := make([]device.TSLModbusCmdTemplate, 0)
	for _, v := range tslModel.Cmd {
		cmds = append(cmds, *v)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取采集模型[Modbus]命令成功",
		Data:    cmds,
	})
}

func ApiGetTSLD07Cmd(context *gin.Context) {

	tslName := context.Query("tslName")

	tslModel, ok := device.TSLDLT6452007Map[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集模型[DLT645-2007]命令错误," + "模型名称不存在",
			Data:    "",
		})
		return
	}

	cmds := make([]device.TSLDLT6452007CmdTemplate, 0)
	for _, v := range tslModel.Cmd {
		cmds = append(cmds, *v)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取采集模型[DLT645-2007]命令成功",
		Data:    cmds,
	})
}

func ApiAddTSLModbusCmdProperty(context *gin.Context) {

	propertyParam := &struct {
		TSLName     string                                `json:"tslName"` // 名称
		CmdName     string                                `json:"cmdName"`
		Name        string                                `json:"name"`
		Label       string                                `json:"label"`
		AccessMode  int                                   `json:"accessMode"`
		Type        int                                   `json:"type"`
		Decimals    int                                   `json:"decimals"`
		Unit        string                                `json:"unit"`
		RegAddr     int                                   `json:"regAddr"`
		RegCnt      int                                   `json:"regCnt"`
		RuleType    string                                `json:"ruleType"`
		Formula     string                                `json:"formula"`
		BitOffsetSw bool                                  `json:"bitSwitch"` // 位偏移开关
		BitOffset   int                                   `json:"bitOffset"` // 位偏移数量
		Params      device.TSLModbusPropertyParamTemplate `json:"params"`    //ltg add 2023-06-15
		IotDataType string                                `json:"iotDataType"`
		Identity    string                                `json:"identity"` //唯一标识
	}{}

	emController := controllers.NewEMController()
	emController.AddEmDeviceModelCmdParam(context)

	err := context.ShouldBindBodyWith(propertyParam, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("增加采集模型[modbus]属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加采集模型[modbus]属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLModbusMap[propertyParam.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[modbus]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	//ltg add 2023-06-15
	param := device.TSLModbusPropertyParamTemplate{
		Min:             propertyParam.Params.Min,
		Max:             propertyParam.Params.Max,
		MinMaxAlarm:     propertyParam.Params.MinMaxAlarm,
		Step:            propertyParam.Params.Step,
		StepAlarm:       propertyParam.Params.StepAlarm,
		DataLength:      propertyParam.Params.DataLength,
		DataLengthAlarm: propertyParam.Params.DataLengthAlarm,
	}

	property := device.TSLModbusPropertyTemplate{
		Name:        propertyParam.Name,
		Label:       propertyParam.Label,
		AccessMode:  propertyParam.AccessMode,
		Type:        propertyParam.Type,
		Decimals:    propertyParam.Decimals,
		Unit:        propertyParam.Unit,
		RegAddr:     propertyParam.RegAddr,
		RegCnt:      propertyParam.RegCnt,
		RuleType:    propertyParam.RuleType,
		Formula:     propertyParam.Formula,
		BitOffsetSw: propertyParam.BitOffsetSw,
		BitOffset:   propertyParam.BitOffset,
		Params:      param,
		IotDataType: propertyParam.IotDataType,
		Identity:    propertyParam.Identity,
	}
	err = tslModel.TSLModelPropertiesAdd(propertyParam.CmdName, property)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[modbus]属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加采集模型[modbus]属性成功",
		Data:    "",
	})

}

func ApiAddTSLD07CmdProperty(context *gin.Context) {

	propertyParam := &struct {
		TSLName        string                                    `json:"tslName"` // 名称
		CmdName        string                                    `json:"cmdName"`
		Name           string                                    `json:"name"`
		Label          string                                    `json:"label"`
		RulerId        string                                    `json:"rulerId"` //数据标识
		Format         string                                    `json:"format"`  //数据格式YYMMDDhhmm,XXXXXX.XX,XX.XXXX...
		Len            int                                       `json:"len"`     //数据长度
		Unit           string                                    `json:"unit"`
		AccessMode     int                                       `json:"accessMode"`
		BlockAddOffset int                                       `json:"blockAddOffset"` //当前数据在块数据域内的偏移地址
		RulerAddOffset int                                       `json:"rulerAddOffset"` //当前变量在当前ID数据地址中的偏移地址
		Type           int                                       `json:"type"`           //float,uint32...
		Params         device.TSLDLT6452007PropertyParamTemplate `json:"params"`         //ltg add 2023-06-15
		IotDataType    string                                    `json:"iotDataType"`
		Identity       string                                    `json:"identity"` //唯一标识
	}{}

	emController := controllers.NewEMController()
	emController.AddEmDeviceModelCmdParam(context)

	err := context.ShouldBindBodyWith(propertyParam, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("增加采集模型[modbus]属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加采集模型[DLT645-2007]属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLDLT6452007Map[propertyParam.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[modbus]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	//ltg add 2023-06-15
	param := device.TSLDLT6452007PropertyParamTemplate{
		Min:             propertyParam.Params.Min,
		Max:             propertyParam.Params.Max,
		MinMaxAlarm:     propertyParam.Params.MinMaxAlarm,
		Step:            propertyParam.Params.Step,
		StepAlarm:       propertyParam.Params.StepAlarm,
		DataLength:      propertyParam.Params.DataLength,
		DataLengthAlarm: propertyParam.Params.DataLengthAlarm,
	}

	property := device.TSLDLT6452007PropertyTemplate{
		Name:           propertyParam.Name,
		Label:          propertyParam.Label,
		RulerId:        propertyParam.RulerId,
		Format:         propertyParam.Format,
		Len:            propertyParam.Len,
		Unit:           propertyParam.Unit,
		AccessMode:     propertyParam.AccessMode,
		BlockAddOffset: propertyParam.BlockAddOffset,
		RulerAddOffset: propertyParam.RulerAddOffset,
		Type:           propertyParam.Type,
		Params:         param,
		IotDataType:    propertyParam.IotDataType,
		Identity:       propertyParam.Identity,
	}
	err = tslModel.TSLModelPropertiesAdd(propertyParam.CmdName, property)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[DLT645-2007]属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "添加采集模型[DLT645-2007]属性成功",
		Data:    "",
	})

}

// ApiExportTSLModbusCmdToXlsx 导出命令列表
func ApiExportTSLModbusCmdToXlsx(context *gin.Context) {
	//获取采集模型名称
	tslName := context.Query("name")

	tslModel, ok := device.TSLModbusMap[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[modbus]命令列表错误，" + "模型名称不存在",
			Data:    "",
		})
		return
	}

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + tslName + ".xlsx"

	//创建一个新的写入文件流
	csvRecords := make([][]string, 0)

	csvRecords = [][]string{
		{"序号", "命令名称", "命令标识符", "功能码", "寄存器地址", "寄存器数量"},
		{"Index", "Name", "Label", "FunCode", "RegAddr", "RegCnt"},
	}

	k := 0
	for _, v := range tslModel.Cmd {
		record := make([]string, 0)

		k++
		record = append(record, fmt.Sprintf("%d", k))
		record = append(record, v.Name)
		record = append(record, v.Label)
		record = append(record, fmt.Sprintf("%d", v.FunCode))
		record = append(record, fmt.Sprintf("%d", v.StartRegAddr))
		record = append(record, fmt.Sprintf("%d", v.RegCnt))

		csvRecords = append(csvRecords, record)
	}

	err := setting.WriteExcel(fileName, csvRecords)
	if err != nil {
		setting.ZAPS.Errorf("保存xlsx文件错误")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存xlsx文件错误",
			Data:    "",
		})
		return
	}

	//返回文件流
	context.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(fileName))))
	context.File(fileName) //返回文件路径，自动调用http.ServeFile方法

	return

}

// ApiAddTSLModbusCmdFromXlsx 导入命令列表
func ApiAddTSLModbusCmdFromXlsx(context *gin.Context) {
	//获取采集模型名称
	tslName := context.PostForm("name")

	tslModel, ok := device.TSLModbusMap[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[modbus]命令列表错误，" + "模型名称不存在",
			Data:    "",
		})
		return
	}

	// 获取文件头
	file, err := context.FormFile("fileName")
	if err != nil {
		setting.ZAPS.Errorf("modbus采集模型待导入命令列表xlsx文件不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "modbus采集模型待导入命令列表xlsx文件不存在",
			Data:    "",
		})
		return
	}

	//创建文件
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

	setting.ZAPS.Debugf("cells %v", cells)
	for _, cell := range cells {
		if len(cell) < 5 {
			continue
		}

		cmd := device.TSLModbusCmdTemplate{
			Name:         setting.GetString(cell[1]),
			Label:        setting.GetString(cell[2]),
			FunCode:      setting.GetInt(cell[3]),
			StartRegAddr: setting.GetInt(cell[4]),
			RegCnt:       setting.GetInt(cell[5]),
		}

		err = tslModel.TSLModelCmdAdd(cmd)
		// 导入cmd写入sqlite
		emController := controllers.NewEMController()
		emController.AddEmDeviceModelCmdFromXlsx(cmd, "modbus", tslName)

		if err != nil {
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "modbus采集模型导入命令列表xlsx错误 " + err.Error(),
				Data:    "",
			})
			return
		}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "modbus采集模型导入命令列表xlsx成功",
		Data:    "",
	})

	return
}

func ApiAddTSLModbusCmdPropertyFromXlsx(context *gin.Context) {

	//获取采集模型名称
	tslName := context.PostForm("tslName")
	//获取modbus命令名称
	cmdName := context.PostForm("cmdName")

	tslModel, ok := device.TSLModbusMap[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[modbus]属性错误，" + "模型名称不存在",
			Data:    "",
		})
		return
	}

	_, ok = tslModel.Cmd[cmdName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[modbus]属性错误，" + "命令名称不存在",
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

	//创建文件
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
		if len(cell) < 20 {
			continue
		}

		property := device.TSLModbusPropertyTemplate{
			Name:        setting.GetString(cell[0]),
			Label:       setting.GetString(cell[1]),
			AccessMode:  setting.GetInt(cell[2]),
			Type:        setting.GetInt(cell[3]),
			Decimals:    setting.GetInt(cell[4]),
			Unit:        setting.GetString(cell[5]),
			RegAddr:     setting.GetInt(cell[6]),
			RegCnt:      setting.GetInt(cell[7]),
			RuleType:    setting.GetString(cell[8]),
			IotDataType: setting.GetString(cell[19]),
			Identity:   setting.GetString(cell[19]),
		}
		if setting.GetString(cell[2]) == "R" {
			property.AccessMode = 0
		} else if setting.GetString(cell[2]) == "W" {
			property.AccessMode = 1
		} else if setting.GetString(cell[2]) == "RW" {
			property.AccessMode = 2
		} else {
			property.AccessMode = 2
		}

		if setting.GetString(cell[3]) == "uint32" {
			property.Type = 0
		} else if setting.GetString(cell[3]) == "int32" {
			property.Type = 1
		} else if setting.GetString(cell[3]) == "double" {
			property.Type = 2
		} else {
			property.Type = 2
		}

		if setting.GetString(cell[9]) == "-" {
			property.Formula = ""
		} else {
			property.Formula = setting.GetString(cell[9])
		}
		if setting.GetString(cell[10]) == "-" {
			property.Formula = ""
		}

		//ltg add 2023-06-16
		if setting.GetString(cell[10]) == "ture" {
			property.BitOffsetSw = true
		} else {
			property.BitOffsetSw = false
		}
		property.BitOffset = setting.GetInt(cell[11])

		if setting.GetString(cell[12]) == "ture" {
			property.Params.MinMaxAlarm = true
		} else {
			property.Params.MinMaxAlarm = false
		}
		property.Params.Min = setting.GetString(cell[13])
		property.Params.Max = setting.GetString(cell[14])

		if setting.GetString(cell[15]) == "ture" {
			property.Params.StepAlarm = true
		} else {
			property.Params.StepAlarm = false
		}
		property.Params.Step = setting.GetString(cell[16])

		if setting.GetString(cell[17]) == "ture" {
			property.Params.DataLengthAlarm = true
		} else {
			property.Params.DataLengthAlarm = false
		}
		property.Params.DataLength = setting.GetString(cell[18])
		property.IotDataType = setting.GetString(cell[19])

		err = tslModel.TSLModelPropertiesAdd(cmdName, property)
		// 导入param写入sqlite
		emController := controllers.NewEMController()
		emController.AddEmDeviceModelCmdParamFromXlsx(property, "modbus", cmdName)

		if err != nil {
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "modbus采集模型导入xlsx错误 " + err.Error(),
				Data:    "",
			})
			return
		}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "modbus采集模型导入xlsx成功",
		Data:    "",
	})

	return
}

func ApiAddTSLD07CmdPropertyFromXlsx(context *gin.Context) {

	//获取采集模型名称
	tslName := context.PostForm("tslName")
	//获取DLT645-2007命令名称
	cmdName := context.PostForm("cmdName")

	tslModel, ok := device.TSLDLT6452007Map[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[DLT645-2007]属性错误，" + "模型名称不存在",
			Data:    "",
		})
		return
	}

	_, ok = tslModel.Cmd[cmdName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[DLT645-2007]属性错误，" + "命令名称不存在",
			Data:    "",
		})
		return
	}

	// 获取文件头
	file, err := context.FormFile("fileName")
	if err != nil {
		setting.ZAPS.Errorf("dlt645-2007采集模型待导入xlsx文件不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "DLT645-2007采集模型待导入xlsx文件不存在",
			Data:    "",
		})
		return
	}

	//创建文件
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
		if len(cell) < 18 {
			continue
		}

		property := device.TSLDLT6452007PropertyTemplate{
			Name:           setting.GetString(cell[0]),
			Label:          setting.GetString(cell[1]),
			RulerId:        setting.GetString(cell[2]),
			Format:         setting.GetString(cell[3]),
			Len:            setting.GetInt(cell[4]),
			AccessMode:     setting.GetInt(cell[5]),
			Type:           setting.GetInt(cell[6]),
			Unit:           setting.GetString(cell[7]),
			BlockAddOffset: setting.GetInt(cell[8]),
			RulerAddOffset: setting.GetInt(cell[9]),
			IotDataType:    setting.GetString(cell[17]),
			Identity:       setting.GetString(cell[17]),
		}
		if setting.GetString(cell[5]) == "R" {
			property.AccessMode = 0
		} else if setting.GetString(cell[5]) == "W" {
			property.AccessMode = 1
		} else if setting.GetString(cell[5]) == "RW" {
			property.AccessMode = 2
		} else {
			property.AccessMode = 2
		}

		if setting.GetString(cell[6]) == "uint32" {
			property.Type = 0
		} else if setting.GetString(cell[6]) == "int32" {
			property.Type = 1
		} else if setting.GetString(cell[6]) == "double" {
			property.Type = 2
		} else {
			property.Type = 2
		}

		//ltg add 2023-06-16
		if setting.GetString(cell[10]) == "ture" {
			property.Params.MinMaxAlarm = true
		} else {
			property.Params.MinMaxAlarm = false
		}
		property.Params.Min = setting.GetString(cell[11])
		property.Params.Max = setting.GetString(cell[12])

		if setting.GetString(cell[13]) == "ture" {
			property.Params.StepAlarm = true
		} else {
			property.Params.StepAlarm = false
		}
		property.Params.Step = setting.GetString(cell[14])

		if setting.GetString(cell[15]) == "ture" {
			property.Params.DataLengthAlarm = true
		} else {
			property.Params.DataLengthAlarm = false
		}
		property.Params.DataLength = setting.GetString(cell[16])
		property.IotDataType = setting.GetString(cell[17])

		err = tslModel.TSLModelPropertiesAdd(cmdName, property)

		// 导入param写入sqlite
		emController := controllers.NewEMController()
		emController.AddEmDeviceModelCmdParamFromXlsx(property, "dlt645", cmdName)

		if err != nil {
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "dlt645-2007采集模型导入xlsx错误 " + err.Error(),
				Data:    "",
			})
			return
		}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "dlt645-2007采集模型导入xlsx成功",
		Data:    "",
	})

	return
}

func ApiModifyTSLModbusCmdProperty(context *gin.Context) {

	propertyParam := &struct {
		TSLName    string `json:"tslName"` // 名称
		CmdName    string `json:"cmdName"`
		Name       string `json:"name"`
		Label      string `json:"label"`
		AccessMode int    `json:"accessMode"`
		Type       int    `json:"type"`
		Decimals   int    `json:"decimals"`
		Unit       string `json:"unit"`
		RegAddr    int    `json:"regAddr"`
		RegCnt     int    `json:"regCnt"`
		RuleType   string `json:"ruleType"`
		Formula    string `json:"formula"`

		BitOffsetSw bool                                  `json:"bitSwitch"` // 位偏移开关
		BitOffset   int                                   `json:"bitOffset"` // 位偏移数量
		Params      device.TSLModbusPropertyParamTemplate `json:"params"`    //ltg add 2023-06-15
		IotDataType string                                `json:"iotDataType"`
		Identity    string                                `json:"identity"`  //唯一标识
	}{}

	emController := controllers.NewEMController()
	emController.UpdateEmDeviceModelCmdParam(context)

	err := context.ShouldBindBodyWith(propertyParam, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("修改采集模型[modbus]属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[modbus]属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLModbusMap[propertyParam.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[modbus]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	//ltg add 2023-06-15
	param := device.TSLModbusPropertyParamTemplate{
		Min:             propertyParam.Params.Min,
		Max:             propertyParam.Params.Max,
		MinMaxAlarm:     propertyParam.Params.MinMaxAlarm,
		Step:            propertyParam.Params.Step,
		StepAlarm:       propertyParam.Params.StepAlarm,
		DataLength:      propertyParam.Params.DataLength,
		DataLengthAlarm: propertyParam.Params.DataLengthAlarm,
	}

	property := device.TSLModbusPropertyTemplate{
		Name:        propertyParam.Name,
		Label:       propertyParam.Label,
		AccessMode:  propertyParam.AccessMode,
		Type:        propertyParam.Type,
		Decimals:    propertyParam.Decimals,
		Unit:        propertyParam.Unit,
		RegAddr:     propertyParam.RegAddr,
		RegCnt:      propertyParam.RegCnt,
		RuleType:    propertyParam.RuleType,
		Formula:     propertyParam.Formula,
		BitOffsetSw: propertyParam.BitOffsetSw,
		BitOffset:   propertyParam.BitOffset,
		Params:      param,
		IotDataType: propertyParam.IotDataType,
		Identity:    propertyParam.Identity,
	}
	err = tslModel.TSLModelPropertiesModify(propertyParam.CmdName, property)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[modbus]属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改采集模型[modbus]属性成功",
		Data:    "",
	})

}

func ApiModifyTSLD07CmdProperty(context *gin.Context) {

	propertyParam := &struct {
		TSLName        string                                    `json:"tslName"` // 名称
		CmdName        string                                    `json:"cmdName"`
		Name           string                                    `json:"name"`
		Label          string                                    `json:"label"`
		RulerId        string                                    `json:"rulerId"` //数据标识
		Format         string                                    `json:"format"`  //数据格式YYMMDDhhmm,XXXXXX.XX,XX.XXXX...
		Len            int                                       `json:"len"`     //数据长度
		Unit           string                                    `json:"unit"`
		AccessMode     int                                       `json:"accessMode"`
		BlockAddOffset int                                       `json:"blockAddOffset"` //当前数据在块数据域内的偏移地址
		RulerAddOffset int                                       `json:"rulerAddOffset"` //当前变量在当前ID数据地址中的偏移地址
		Type           int                                       `json:"type"`           //float,uint32...
		Params         device.TSLDLT6452007PropertyParamTemplate `json:"params"`         //ltg add 2023-06-15
		IotDataType    string                                    `json:"iotDataType"`
		Identity       string                                    `json:"identity"`       //唯一标识
	}{}

	emController := controllers.NewEMController()
	emController.UpdateEmDeviceModelCmdParam(context)

	err := context.ShouldBindBodyWith(propertyParam, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("修改采集模型[modbus]属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[DLT645-2007]属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLDLT6452007Map[propertyParam.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[DLT645-2007]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	//ltg add 2023-06-15
	param := device.TSLDLT6452007PropertyParamTemplate{
		Min:             propertyParam.Params.Min,
		Max:             propertyParam.Params.Max,
		MinMaxAlarm:     propertyParam.Params.MinMaxAlarm,
		Step:            propertyParam.Params.Step,
		StepAlarm:       propertyParam.Params.StepAlarm,
		DataLength:      propertyParam.Params.DataLength,
		DataLengthAlarm: propertyParam.Params.DataLengthAlarm,
	}

	property := device.TSLDLT6452007PropertyTemplate{
		Name:           propertyParam.Name,
		Label:          propertyParam.Label,
		RulerId:        propertyParam.RulerId,
		Format:         propertyParam.Format,
		Len:            propertyParam.Len,
		Unit:           propertyParam.Unit,
		AccessMode:     propertyParam.AccessMode,
		BlockAddOffset: propertyParam.BlockAddOffset,
		RulerAddOffset: propertyParam.RulerAddOffset,
		Type:           propertyParam.Type,
		Params:         param,
		IotDataType:    propertyParam.IotDataType,
		Identity:       propertyParam.Identity,
	}
	err = tslModel.TSLModelPropertiesModify(propertyParam.CmdName, property)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改采集模型[DLT645-2007]属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改采集模型[DLT645-2007]属性成功",
		Data:    "",
	})

}

func ApiDeleteTSLModbusCmdProperties(context *gin.Context) {

	tslInfo := &struct {
		TSLName       string   `json:"tslName"` // 名称
		CmdName       string   `json:"cmdName"`
		PropertyNames []string `json:"propertyNames"` //
	}{}

	emController := controllers.NewEMController()
	emController.DeleteEmDeviceModelCmdParam(context)

	err := context.ShouldBindBodyWith(&tslInfo, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("删除采集模型[modbus]属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集模型[modbus]属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLModbusMap[tslInfo.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集模型[modbus]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	err = tslModel.TSLModelPropertiesDelete(tslInfo.CmdName, tslInfo.PropertyNames)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集模型[modbus]属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除采集模型[modbus]属性成功",
		Data:    "",
	})
}
func ApiDeleteTSLD07CmdProperties(context *gin.Context) {

	tslInfo := &struct {
		TSLName       string   `json:"tslName"` // 名称
		CmdName       string   `json:"cmdName"`
		PropertyNames []string `json:"propertyNames"` //
	}{}

	emController := controllers.NewEMController()
	emController.DeleteEmDeviceModelCmdParam(context)

	err := context.ShouldBindBodyWith(&tslInfo, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("删除采集模型[DLT645-2007]属性JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集模型[DLT645-2007]属性JSON格式化错误",
			Data:    "",
		})
		return
	}

	tslModel, ok := device.TSLDLT6452007Map[tslInfo.TSLName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集模型[DLT645-2007]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	err = tslModel.TSLModelPropertiesDelete(tslInfo.CmdName, tslInfo.PropertyNames)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除采集模型[DLT645-2007]属性错误 " + err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除采集模型[DLT645-2007]属性成功",
		Data:    "",
	})
}

func ApiGetTSLModbusProperties(context *gin.Context) {

	tslName := context.Query("tslName")

	tslModel, ok := device.TSLModbusMap[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取采集模型[modbus]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取物模型属性成功",
		Data:    tslModel.GetTSLModbusModelProperties(),
	})
}

func ApiGetTSLD07Properties(context *gin.Context) {

	tslName := context.Query("tslName")

	tslModel, ok := device.TSLDLT6452007Map[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取采集模型[DLT645-2007]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取物模型属性成功",
		Data:    tslModel.GetTSLDLT6452007ModelProperties(),
	})
}

func ApiGetTSLModbusCmdProperties(context *gin.Context) {

	tslName := context.Query("tslName")
	cmdName := context.Query("cmdName")

	tslModel, ok := device.TSLModbusMap[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取采集模型[modbus]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	cmd, ok := tslModel.Cmd[cmdName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取采集模型[modbus]属性错误 " + "命令名称不存在",
			Data:    "",
		})
		return
	}

	params := make([]device.TSLModbusPropertyTemplate, 0)
	for _, v := range cmd.Registers {
		params = append(params, *v)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取物模型属性成功",
		Data:    params,
	})
}

func ApiGetTSLD07CmdProperties(context *gin.Context) {

	tslName := context.Query("tslName")
	cmdName := context.Query("cmdName")

	tslModel, ok := device.TSLDLT6452007Map[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取采集模型[DLT645-2007]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	cmd, ok := tslModel.Cmd[cmdName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取采集模型[DLT645-2007]属性错误 " + "命令名称不存在",
			Data:    "",
		})
		return
	}

	params := make([]device.TSLDLT6452007PropertyTemplate, 0)
	for _, v := range cmd.Properties {
		params = append(params, *v)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取物模型属性成功",
		Data:    params,
	})
}

func ApiExportTSLModbusCmdPropertiesToXlsx(context *gin.Context) {

	tslName := context.Query("tslName")
	cmdName := context.Query("cmdName")

	tslModel, ok := device.TSLModbusMap[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取采集模型[modbus]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	cmd, ok := tslModel.Cmd[cmdName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取采集模型[modbus]属性错误 " + "命令名称不存在",
			Data:    "",
		})
		return
	}

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + tslName + cmdName + ".xlsx"

	//创建一个新的写入文件流
	csvRecords := make([][]string, 0)

	csvRecords = [][]string{
		{"属性名称", "属性标识符", "读写类型", "数据类型", "小数位", "单位", "寄存器地址", "寄存器数量", "解析规则", "计算公式", "位解析开关", "位偏移", "范围报警", "最小值", "最大值", "步长报警", "步长", "字符串长度报警", "字符串长度", "点位类型", "唯一标识"},
		{"Name", "Label", "AccessMode", "Type", "Decimals", "Unit", "RegAddr", "RegCnt", "RuleType", "Formula", "BitOffsetSw", "BitOffset", "MinMaxAlarm", "Min", "Max", "StepAlarm", "Step", "DataLengthAlarm", "DataLength", "IotDataType", "Identity"},
	}

	for _, v := range cmd.Registers {
		record := make([]string, 0)
		record = append(record, v.Name)
		record = append(record, v.Label)
		if v.AccessMode == 0 {
			record = append(record, "R")
		} else if v.AccessMode == 1 {
			record = append(record, "W")
		} else if v.AccessMode == 2 {
			record = append(record, "RW")
		}
		if v.Type == 0 {
			record = append(record, "uint32")
		} else if v.Type == 1 {
			record = append(record, "int32")
		} else if v.Type == 2 {
			record = append(record, "double")
		} else if v.Type == 3 {
			record = append(record, "string")
		}
		record = append(record, fmt.Sprintf("%d", v.Decimals))
		record = append(record, v.Unit)
		record = append(record, fmt.Sprintf("%d", v.RegAddr))
		record = append(record, fmt.Sprintf("%d", v.RegCnt))
		record = append(record, v.RuleType)
		if v.Formula == "" {
			record = append(record, "-")
		} else {
			record = append(record, v.Formula)
		}
		//ltg add 2023-06-16
		if v.BitOffsetSw == true {
			record = append(record, "true")
		} else {
			record = append(record, "false")
		}

		record = append(record, fmt.Sprintf("%d", v.BitOffset))

		if v.Params.MinMaxAlarm == true {
			record = append(record, "true")
		} else {
			record = append(record, "false")
		}

		record = append(record, v.Params.Min)
		record = append(record, v.Params.Max)

		if v.Params.StepAlarm == true {
			record = append(record, "true")
		} else {
			record = append(record, "false")
		}

		record = append(record, v.Params.Step)

		if v.Params.DataLengthAlarm == true {
			record = append(record, "true")
		} else {
			record = append(record, "false")
		}

		record = append(record, v.Params.DataLength)
		record = append(record, v.IotDataType)
		record = append(record, v.Identity)

		csvRecords = append(csvRecords, record)
	}

	err := setting.WriteExcel(fileName, csvRecords)
	if err != nil {
		setting.ZAPS.Errorf("保存xlsx文件错误")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存xlsx文件错误",
			Data:    "",
		})
		return
	}

	//返回文件流
	context.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(fileName))))
	context.File(fileName) //返回文件路径，自动调用http.ServeFile方法

	return
}

func ApiExportTSLD07CmdPropertiesToXlsx(context *gin.Context) {

	tslName := context.Query("tslName")
	cmdName := context.Query("cmdName")

	tslModel, ok := device.TSLDLT6452007Map[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取采集模型[DLT645-2007]属性错误 " + "模型名称不存在",
			Data:    "",
		})
		return
	}

	cmd, ok := tslModel.Cmd[cmdName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取采集模型[DLT645-2007]属性错误 " + "命令名称不存在",
			Data:    "",
		})
		return
	}

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + tslName + cmdName + ".xlsx"

	//创建一个新的写入文件流
	csvRecords := make([][]string, 0)
	csvRecords = [][]string{
		{"属性名称", "属性标识符", "数据标识", "数据格式", "数据长度", "读写类型", "数据类型", "单位", "数据块偏移地址", "数据标识偏移地址", "范围报警", "最小值", "最大值", "步长报警", "步长", "字符串长度报警", "字符串长度", "点位类型", "唯一标识"},
		{"Name", "Label", "RulerId", "Format", "Len", "AccessMode", "Type", "Unit", "BlockAddOffset", "RulerAddOffset", "MinMaxAlarm", "Min", "Max", "StepAlarm", "Step", "DataLengthAlarm", "DataLength", "IotDataType", "Identity"},
	}

	for _, v := range cmd.Properties {
		record := make([]string, 0)
		record = append(record, v.Name)
		record = append(record, v.Label)
		record = append(record, fmt.Sprintf("%s", v.RulerId))
		record = append(record, fmt.Sprintf("%s", v.Format))
		record = append(record, fmt.Sprintf("%d", v.Len))
		record = append(record, fmt.Sprintf("%s", v.Unit))
		record = append(record, fmt.Sprintf("%d", v.BlockAddOffset))
		record = append(record, fmt.Sprintf("%d", v.RulerAddOffset))

		if v.AccessMode == 0 {
			record = append(record, "R")
		} else if v.AccessMode == 1 {
			record = append(record, "W")
		} else if v.AccessMode == 2 {
			record = append(record, "RW")
		}
		if v.Type == 0 {
			record = append(record, "uint32")
		} else if v.Type == 1 {
			record = append(record, "int32")
		} else if v.Type == 2 {
			record = append(record, "double")
		} else if v.Type == 3 {
			record = append(record, "string")
		}
		//ltg add 2023-06-16
		if v.Params.MinMaxAlarm == true {
			record = append(record, "true")
		} else {
			record = append(record, "false")
		}

		record = append(record, v.Params.Min)
		record = append(record, v.Params.Max)

		if v.Params.StepAlarm == true {
			record = append(record, "true")
		} else {
			record = append(record, "false")
		}

		record = append(record, v.Params.Step)

		if v.Params.DataLengthAlarm == true {
			record = append(record, "true")
		} else {
			record = append(record, "false")
		}

		record = append(record, v.Params.DataLength)
		record = append(record, v.IotDataType)
		record = append(record, v.Identity)

		csvRecords = append(csvRecords, record)
	}

	err := setting.WriteExcel(fileName, csvRecords)
	if err != nil {
		setting.ZAPS.Errorf("保存xlsx文件错误")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存xlsx文件错误",
			Data:    "",
		})
		return
	}

	//返回文件流
	context.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(fileName))))
	context.File(fileName) //返回文件路径，自动调用http.ServeFile方法

	return
}

// ApiExportTSLD07CmdToXlsx 导出命令列表
func ApiExportTSLD07CmdToXlsx(context *gin.Context) {
	//获取采集模型名称
	tslName := context.Query("name")

	tslModel, ok := device.TSLDLT6452007Map[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[DLT645-2007]命令列表错误，" + "模型名称不存在",
			Data:    "",
		})
		return
	}

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + tslName + ".xlsx"

	//创建一个新的写入文件流
	csvRecords := make([][]string, 0)

	csvRecords = [][]string{
		{"序号", "命令名称", "命令标识符", "块数据标识", "块读取开关"},
		{"Index", "Name", "Label", "BlockRulerId", "BlockRead"},
	}

	k := 0
	for _, v := range tslModel.Cmd {
		record := make([]string, 0)

		k++
		record = append(record, fmt.Sprintf("%d", k))
		record = append(record, v.Name)
		record = append(record, v.Label)
		record = append(record, v.BlockRulerId)
		if v.BlockRead == 0 {
			record = append(record, "关")
		} else if v.BlockRead == 1 {
			record = append(record, "开")
		}

		csvRecords = append(csvRecords, record)
	}

	err := setting.WriteExcel(fileName, csvRecords)
	if err != nil {
		setting.ZAPS.Errorf("保存xlsx文件错误")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存xlsx文件错误",
			Data:    "",
		})
		return
	}

	//返回文件流
	context.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(fileName))))
	context.File(fileName) //返回文件路径，自动调用http.ServeFile方法

	return

}

// ApiAddTSLD07CmdFromXlsx 导入命令列表
func ApiAddTSLD07CmdFromXlsx(context *gin.Context) {
	//获取采集模型名称
	tslName := context.PostForm("name")

	tslModel, ok := device.TSLDLT6452007Map[tslName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "添加采集模型[DLT645-2007]命令列表错误，" + "模型名称不存在",
			Data:    "",
		})
		return
	}

	// 获取文件头
	file, err := context.FormFile("fileName")
	if err != nil {
		setting.ZAPS.Errorf("DLT645-2007采集模型待导入命令列表xlsx文件不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "DLT645-2007采集模型待导入命令列表xlsx文件不存在",
			Data:    "",
		})
		return
	}

	//创建文件
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
		if len(cell) < 5 {
			continue
		}

		cmd := device.TSLDLT6452007CmdTemplate{
			Name:         setting.GetString(cell[1]),
			Label:        setting.GetString(cell[2]),
			BlockRulerId: setting.GetString(cell[3]),
		}

		if setting.GetString(cell[4]) == "关" {
			cmd.BlockRead = 0
		} else if setting.GetString(cell[4]) == "开" {
			cmd.BlockRead = 1
		}

		err = tslModel.TSLModelCmdAdd(cmd)
		// 导入cmd写入sqlite
		emController := controllers.NewEMController()
		emController.AddEmDeviceModelCmdFromXlsx(cmd, "dlt645", tslName)

		if err != nil {
			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "1",
				Message: "DLT645-2007采集模型导入命令列表xlsx错误 " + err.Error(),
				Data:    "",
			})
			return
		}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "DLT645-2007采集模型导入命令列表xlsx成功",
		Data:    "",
	})

	return
}
