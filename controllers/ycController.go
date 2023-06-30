package controllers

import (
	"gateway/httpServer/model"
	"gateway/models/query"
	"gateway/repositories"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义辅控管理的控制器
type YcController struct {
	hisRepo *repositories.HistoryDataRepository
}

func NewYcController() *YcController {
	return &YcController{
		hisRepo: repositories.NewHistoryDataRepository()}
}
func (c *YcController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/yc/getLastYcByDeviceIdAndCodes", c.GetLastYcByDeviceIdsAndCodes)
	router.POST("/api/v2/yc/batchYcHistoryListByDeviceIdAndCodes", c.BatchYcHistoryListByDeviceIdAndCodes)
}

// GetLastYcByDeviceIdsAndCodes 获取最新遥测信息GetLastYcListByCode
/*
{
    "deviceIds":[36559],  //设备id 必填
    "codes":[2002,2005]  //遥测编码  必填
}
*/
func (c *YcController) GetLastYcByDeviceIdsAndCodes(ctx *gin.Context) {
	type YcQeury struct {
		DeviceIds []int `form:"deviceIds"`
		Codes     []int `form:"codes"`
	}
	var ycQuery YcQeury
	//将传过来的请求体解析到ycQuery中
	if err := ctx.ShouldBindJSON(&ycQuery); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	//查询历史数据
	//将ycQuery.DeviceIds转换成字符串
	deviceIds := utils.IntArrayToString(ycQuery.DeviceIds, ",")
	codes := utils.IntArrayToString(ycQuery.Codes, ",")
	ycLog, err := c.hisRepo.GetLastYcListByCode(deviceIds, codes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取信息成功！",
		Data:    ycLog,
	})
}

// BatchYcHistoryListByDeviceIdAndCodes 根据时间获取多个遥测的历史数据信息
/*{
"deviceIds":36559,  //设备id集合 必填
"codeList":[2002,2005],  //遥测编码集合  必填
"startTime":946685445000, //开始时间 必填
"endTime": 947519999000,  //结束时间  必填
"interval":10,  //间隔时间  必填
"intervalType":2 //间隔类型 1秒 2分 3时 4天 必填
}*/
func (c *YcController) BatchYcHistoryListByDeviceIdAndCodes(ctx *gin.Context) {
	var ycQuery query.QueryTaoData
	//解析json
	if err := ctx.Bind(&ycQuery); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	returnMap, err := c.hisRepo.GetCharData(ycQuery)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    nil,
		})
	} else {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "获取信息成功！",
			Data:    returnMap,
		})
	}
}
