package controllers

import (
	"fmt"
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/models/query"
	repositories "gateway/repositories"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//定义辅控管理的控制器
type AuxiliaryController struct {
	repo    *repositories.AuxiliaryRepository
	hisRepo *repositories.HistoryDataRepository
}

func NewAuxiliaryController() *AuxiliaryController {
	return &AuxiliaryController{repo: repositories.NewAuxiliaryRepository(),
		hisRepo: repositories.NewHistoryDataRepository()}
}

func (ctrl *AuxiliaryController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/api/v2/auxiliary/getDeviceListByDeviceType", ctrl.GetDeviceListByDeviceType)
	router.GET("/api/v2/auxiliary/getDeviceType", ctrl.GetAuxiliaryDeviceType)
	router.POST("/api/v2/auxiliary/getYcLogById", ctrl.GetYcLogById)
	router.POST("/api/v2/auxiliary/getHistoryYcByDeviceIdCodes", ctrl.GetHistoryYcByDeviceIdCodes)
}

///获取设备类型下的所有设备数据
func (c *AuxiliaryController) GetDeviceListByDeviceType(ctx *gin.Context) {
	label := ctx.Query("label")
	deviceList, err := c.repo.GetAuxiliaryDevice(label)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		deviceList,
	})
	return
}

//获取所有设备类型
func (c *AuxiliaryController) GetAuxiliaryDeviceType(ctx *gin.Context) {
	deviceTypeList, err := c.repo.GetAuxiliaryDeviceType()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		deviceTypeList,
	})
	return
}

//获取遥信息GetYcLogById
func (c *AuxiliaryController) GetYcLogById(ctx *gin.Context) {
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

type MapWrapper struct {
	ycMap     map[int][]string
	xAxisList []string
}

type ReturnMap struct {
	XAxisList []string          `json:"xAxisList"`
	DataMap   map[int][]float64 `json:"dataMap"`
}

//根据选择的codes返回对应的历史数据
func (c *AuxiliaryController) GetHistoryYcByDeviceIdCodes(ctx *gin.Context) {
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
	dateHistoryMap := make(map[string][]models.YcData)
	// 初始化x轴数据,返回x轴时间对应的历史数据分组,key=x轴,value=x轴对应的历史数据集合
	xAxisList, dateHistoryMap = initXAxisList(xAxisList, ycQuery.StartTime, ycQuery.EndTime, ycQuery.Interval, ycQuery.IntervalType, ycList)
	println(dateHistoryMap)
	////TODO 根据codes查询字典，获取出codes对于的名称
	////for循环ycList按code进行分组
	//var ycMap = make(map[int][]string)
	//for _, yc := range ycList {
	//	ycMap[yc.Code] = append(ycMap[yc.Code], yc.Value)
	//}
	////TODO 拼接返回时间段
	dataMap := make(map[int][]float64)
	for _, code := range ycQuery.CodeList { //按code分组
		var valList []float64 //存储结果值
		for _, xAxis := range xAxisList {
			ycHistoryList, exists := dateHistoryMap[xAxis] //根据时间获取值
			var val float64
			if exists { //如果值不为空，再按code进行分组
				//将ycHistoruList按code进行分组收集
				codeMap := make(map[int][]models.YcData)
				// 按照 编码进行分组
				for _, person := range ycHistoryList {
					group := person.Code
					codeMap[group] = append(codeMap[group], person)
				}
				ycModels := codeMap[code] //将code对应的值取出来
				if len(ycModels) > 0 {
					//统计值
					val = utils.YcValueSum(ycModels)
					fmt.Println(val)
				}
			}
			valList = append(valList, val)
		}
		dataMap[code] = valList //将结果存入map
	}
	var returnMap ReturnMap
	returnMap.XAxisList = xAxisList
	returnMap.DataMap = dataMap
	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		returnMap,
	})
}

func initXAxisList(xAxisList []string, beginDt int64, endDt int64, interval int, intervalType int, historyList []*models.YcData) ([]string, map[string][]models.YcData) {
	historyGroupMap := make(map[string][]models.YcData)
	//创建Calendar对象并设置时间
	calendar := time.Unix(0, beginDt*int64(time.Millisecond))
	//初始化Calendar对象的时间部分为0
	calendar = time.Date(calendar.Year(), calendar.Month(), calendar.Day(), 0, 0, 0, 0, calendar.Location())
	// 历史数据分组 map，key=x轴，value=x轴对应的历史数据
	//格式化日期后面遍历存储用
	dataFormat := utils.GetIntervalDateFormat(intervalType)
	//获取当前日期
	for i := beginDt; i <= endDt; {
		//先将毫秒转换成秒，再转成t.time对象
		t := time.Unix(0, i*int64(time.Millisecond))
		//转换日期作为前端展示
		format := t.Format(dataFormat)        //格式化日期
		xAxisList = append(xAxisList, format) //添加到x轴列表

		//计算long长度，增长
		var intervalLong int64
		//计算曾长长度
		intervalLong = utils.GetIntervalTime(calendar, intervalType, interval)
		var list []models.YcData
		intervalStart := t                                                   //开始时间等于当前遍历到的i时间                                          //当前时间
		intervalEnd := t.Add(time.Duration(intervalLong) * time.Millisecond) //当前时间加上长度，等于结束时间，用于后面遍历使用
		/*
		   获取当前时间间隔内的历史数据,有的数据不在x轴整点内,算到上个时间间隔里
		   如按两小时间隔查询历史数据,则x轴为0h,2h,4h...
		   历史集合里没有整点0点的数据,但是有1点的数据,将1点的数据算到0h里
		*/
		var newHistoryList []*models.YcData
		for _, item := range historyList { //遍历历史数据
			if item.Ts.After(intervalStart) && item.Ts.Before(intervalEnd) { //大于开始时间，小于结束时间
				//将符合条件的数据添加到list
				list = append(list, *item)
			} else {
				//如果不符合条件，那么就加进新的list后面重新赋值旧的historyList重新遍历 目的加快查询速度
				newHistoryList = append(newHistoryList, item)
			}
		}
		historyList = newHistoryList
		//添加到结果集
		historyGroupMap[format] = list
		//计算下一个时间段
		i += intervalLong
	}
	return xAxisList, historyGroupMap
}
