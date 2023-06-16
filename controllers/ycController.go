package controllers

import (
	"gateway/httpServer/model"
	"gateway/models/query"
	"gateway/repositories"
	"gateway/service"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

//定义辅控管理的控制器
type YcController struct {
	hisRepo *repositories.HistoryDataRepository
}

func NewYcController() *YcController {
	return &YcController{
		hisRepo: repositories.NewHistoryDataRepository()}
}
func (ctrl *YcController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/yc/getLastYcByDeviceIdAndCodes", ctrl.GetLastYcByDeviceIdsAndCodes)
	router.POST("/api/v2/yc/batchYcHistoryListByDeviceIdAndCodes", ctrl.BatchYcHistoryListByDeviceIdAndCodes)
}

//获取最新遥测信息GetLastYcListByCode
/*
{
    "deviceIds":[36559],  //设备id 必填
    "codes":[2002,2005]  //遥测编码  必填
}
*/
func (c *YcController) GetLastYcByDeviceIdsAndCodes(ctx *gin.Context) {
	type ycQeury struct {
		DeviceIds []int `form:"deviceIds"`
		Codes     []int `form:"codes"`
	}
	var ycQuery ycQeury
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
		"0",
		"获取信息成功！",
		ycLog,
	})
}

//根据时间获取多个遥测的历史数据信息
/*{
"deviceId":36559,  //设备id 必填
"codeList":[2002,2005],  //遥测编码  必填
"startTime":946685445000, //开始时间 必填
"endTime": 947519999000,  //结束时间  必填
"interval":10,  //间隔时间  必填
"intervalType":2 //间隔类型 1秒 2分 3时 4天 必填
}*/
func (c *YcController) BatchYcHistoryListByDeviceIdAndCodes(ctx *gin.Context) {
	var ycQuery *query.QueryTaoData
	//解析json
	if err := ctx.ShouldBindBodyWith(&ycQuery, binding.JSON); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	//必要条件校验
	if ycQuery.StartTime == 0 || ycQuery.EndTime == 0 || ycQuery.Interval == 0 || ycQuery.IntervalType == 0 || len(ycQuery.CodeList) == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error",
			"缺少必要参数",
		})
		return
	}
	//间隔时间段
	intervalStr := "s"
	if ycQuery.IntervalType == 2 {
		intervalStr = "m"
	} else if ycQuery.IntervalType == 3 {
		intervalStr = "h"
	} else if ycQuery.IntervalType > 3 {
		intervalStr = "d"
	}
	// 将 int 数组转换为字符串数组
	ycQuery.Codes = utils.IntArrayToString(ycQuery.CodeList, ",")

	//查询历史数据
	ycList, err := c.hisRepo.GetLastYcHistoryByDeviceIdAndCodeList(ycQuery.DeviceId, ycQuery.Codes, ycQuery.StartTime, ycQuery.EndTime, strconv.Itoa(ycQuery.Interval)+intervalStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var xAxisList []string
	// 初始化x轴数据,返回x轴时间对应的历史数据分组,key=x轴,value=x轴对应的历史数据集合
	returnMap := service.GetCharData(xAxisList, ycQuery.StartTime, ycQuery.EndTime, ycQuery.Interval, ycQuery.IntervalType, ycList, ycQuery.CodeList)

	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		returnMap,
	})
}
