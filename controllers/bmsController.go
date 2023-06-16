package controllers

import (
	"gateway/httpServer/model"
	"gateway/models/query"
	repositories "gateway/repositories"
	"gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"strings"
)

type BmsController struct {
	hisRepo *repositories.HistoryDataRepository
}

func NewBmsController() *BmsController {
	return &BmsController{hisRepo: repositories.NewHistoryDataRepository()}
}
func (ctrl *BmsController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/bms/getYcLastByDeviceIdAndDict", ctrl.GetYcLastByDeviceIdAndDict)
	router.POST("/api/v2/bms/getYcLogById", ctrl.GetYcLogById)
	router.POST("/api/v2/bms/getHistoryYcByDeviceIdCodes", ctrl.GetHistoryYcByDeviceIdCodes)
}

// 获取设备点位最新的一条非空数据
func (c *BmsController) GetYcLastByDeviceIdAndDict(ctx *gin.Context) {
	//1.根据设备id去查询字典
	//2.将查询出来的所有code拼接成字符串
	deviceId := ctx.Value("deviceId")
	//codes := ctx.Value("codes")
	print(deviceId)
	codes := "2005,2002"
	//3.然后查询涛思数据库
	ycList, err := c.hisRepo.GetLastYcListByCode(codes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		ycList,
	})
}

//批量获取遥信息GetYcLogById,用于·历史数据
func (c *BmsController) GetYcLogById(ctx *gin.Context) {
	type ycQeury struct {
		DeviceId  int    `form:"DeviceId"`
		Codes     string `form:"codes"`
		StartTime int64  `form:"StartTime"`
		EndTime   int64  `form:"EndTime"`
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
	ycLog, err := c.hisRepo.GetYcLogById(ycQuery.DeviceId, ycQuery.Codes, ycQuery.StartTime, ycQuery.EndTime)
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

//根据选择的codes返回对应的历史数据
func (c *BmsController) GetHistoryYcByDeviceIdCodes(ctx *gin.Context) {
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
	strArr := make([]string, len(ycQuery.CodeList))
	for i, v := range ycQuery.CodeList {
		strArr[i] = strconv.Itoa(v)
	}

	//拼接codelist sql语句
	ycQuery.Codes = strings.Join(strArr, ",")
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
