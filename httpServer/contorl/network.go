package contorl

import (
	"gateway/httpServer/model"
	"gateway/setting"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func ApiAddNetwork(context *gin.Context) {

	networkParam := setting.NetworkConfigParamTemplate{}

	err := context.ShouldBindJSON(&networkParam)
	if err != nil {
		setting.ZAPS.Errorf("增加网卡偏好设置JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加网卡偏好设置JSON格式化错误",
			Data:    "",
		})
		return
	}

	err = setting.AddNetworkConfigParam(networkParam)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "增加网卡偏好设置成功",
			Data:    "",
		})
	}
}

func ApiModifyNetwork(context *gin.Context) {

	networkParam := setting.NetworkConfigParamTemplate{}

	err := context.ShouldBindJSON(&networkParam)
	if err != nil {
		setting.ZAPS.Errorf("修改网卡偏好设置JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
		return
	}

	err = setting.ModifyNetworkConfigParam(networkParam)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "修改网卡偏好设置成功",
			Data:    "",
		})
	}
}

func ApiDeleteNetwork(context *gin.Context) {

	param := struct {
		Name string `json:"name"`
	}{}

	err := context.ShouldBindJSON(&param)
	if err != nil {
		setting.ZAPS.Errorf("删除网卡偏好设置JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
		return
	}

	err = setting.DeleteNetworkConfigParam(param.Name)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "删除网卡偏好设置成功",
			Data:    "",
		})
	}
}

func ApiGetNetworkParams(context *gin.Context) {

	params := setting.GetNetworkParams()

	//排序，方便前端页面显示
	sort.Slice(params, func(i, j int) bool {
		iName := params[i].Index
		jName := params[j].Index
		return iName < jName
	})

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取网卡信息成功",
		Data:    params,
	})
}

func ApiGetNetworkNames(context *gin.Context) {

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取成功",
		Data:    setting.GetNetworkNames(),
	})
}
