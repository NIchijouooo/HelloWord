package contorl

import (
	"gateway/httpServer/model"
	"gateway/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiGetMobileModuleParam(context *gin.Context) {
	type MobileParamTemplate struct {
		Name           string                                     `json:"name"`           //模块名称
		Model          string                                     `json:"model"`          //模块型号
		RunParam       setting.MobileModuleRunParamTemplate       `json:"runParam"`       //模块运行参数
		ConfigParam    setting.MobileModuleConfigParamTemplate    `json:"configParam"`    //模块配置参数
		CommParam      setting.MobileModuleCommParamTemplate      `json:"commParam"`      //模块通信参数
		KeepAliveParam setting.MobileModuleKeepAliveParamTemplate `json:"keepAliveParam"` //模块保活参数
	}

	mobiles := make([]MobileParamTemplate, 0)
	if setting.MobileModule.Name == "" {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "",
			Data:    mobiles,
		})
		return
	}
	param := MobileParamTemplate{
		Name:           setting.MobileModule.Name,
		Model:          setting.MobileModule.Model,
		RunParam:       setting.MobileModule.RunParam,
		ConfigParam:    setting.MobileModule.ConfigParam,
		CommParam:      setting.MobileModule.CommParam,
		KeepAliveParam: setting.MobileModule.KeepAliveParam,
	}
	mobiles = append(mobiles, param)
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取移动模块参数成功",
		Data:    mobiles,
	})
}

func ApiAddMobileModuleParam(context *gin.Context) {

	param := setting.MobileModuleParamTemplate{}

	err := context.ShouldBindJSON(&param)
	if err != nil {
		setting.ZAPS.Warnf("移动模块参数JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
		return
	}

	err = setting.AddMobileModuleParam(param)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "添加移动模块参数成功",
			Data:    "",
		})
	}
}

func ApiModifyMobileModuleParam(context *gin.Context) {

	param := setting.MobileModuleParamTemplate{}

	err := context.ShouldBindJSON(&param)
	if err != nil {
		setting.ZAPS.Warnf("移动模块参数JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
		return
	}

	err = setting.ModifyMobileModuleParam(param)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "修改移动模块参数成功",
			Data:    "",
		})
	}
}

func ApiDeleteMobileModuleParam(context *gin.Context) {

	param := struct {
		Name string `json:"Name"`
	}{}

	err := context.ShouldBindJSON(&param)
	if err != nil {
		setting.ZAPS.Warnf("移动模块参数JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
		return
	}

	err = setting.DeleteMobileModuleParam(param.Name)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "删除移动模块参数成功",
			Data:    "",
		})
	}
}
