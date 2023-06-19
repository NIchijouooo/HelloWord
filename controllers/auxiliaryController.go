package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/models/query"
	repositories "gateway/repositories"
	"gateway/service"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"strings"
)

//定义辅控管理的控制器
type AuxiliaryController struct {
	repo    *repositories.AuxiliaryRepository
	hisRepo *repositories.HistoryDataRepository
	emRepo  *repositories.EmRepository
}

func NewAuxiliaryController() *AuxiliaryController {
	return &AuxiliaryController{repo: repositories.NewAuxiliaryRepository(),
		hisRepo: repositories.NewHistoryDataRepository(),
		emRepo:  repositories.NewEmRepository()}
}

func (ctrl *AuxiliaryController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/v2/auxiliary/getDeviceListByDeviceType", ctrl.GetDeviceListByDeviceType)
	router.GET("/api/v2/auxiliary/getDeviceType", ctrl.GetAuxiliaryDeviceType)
	router.POST("/api/v2/auxiliary/getLastYcByDeviceIdAndCodes", ctrl.GetLastYcByDeviceIdAndCodes)
	router.POST("/api/v2/auxiliary/getHistoryYcByDeviceIdCodes", ctrl.GetHistoryYcByDeviceIdCodes)
	router.POST("/api/v2/auxiliary/getEmDeviceModelCmdParamListByDeviceId", ctrl.GetEmDeviceModelCmdParamListByDeviceId)
}

///获取设备类型下的所有设备数据
func (c *AuxiliaryController) GetDeviceListByDeviceType(ctx *gin.Context) {
	type tmpQuery struct {
		Label string `json:"label"`
	}
	var query tmpQuery
	if err := ctx.Bind(&query); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			"1",
			"error" + err.Error(),
			"",
		})
		return
	}
	deviceList, err := c.repo.GetAuxiliaryDevice(query.Label)
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

//获取最新遥测信息GetLastYcListByCode
//遥测控制器已经写有
func (c *AuxiliaryController) GetLastYcByDeviceIdAndCodes(ctx *gin.Context) {
	type ycQeury struct {
		DeviceIds []int `form:"deviceIds"`
		Codes     []int `form:"codes"`
	}
	var ycQuery ycQeury
	//将传过来的请求体解析到ycQuery中
	if err := ctx.Bind(&ycQuery); err != nil {
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

type ReturnMap struct {
	XAxisList []string          `json:"xAxisList"`
	DataMap   map[int][]float64 `json:"dataMap"`
}

//根据选择的codes返回对应时间的历史数据

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
	// 初始化x轴数据,返回x轴时间对应的历史数据分组,key=x轴,value=x轴对应的历史数据集合
	returnMap := service.GetCharData(xAxisList, ycQuery.StartTime, ycQuery.EndTime, ycQuery.Interval, ycQuery.IntervalType, ycList, ycQuery.CodeList, ycQuery.CodeNameList)

	ctx.JSON(http.StatusOK, model.ResponseData{
		"0",
		"获取信息成功！",
		returnMap,
	})
}

// 根据设备ID获取所有模型
func (c *AuxiliaryController) GetEmDeviceModelCmdParamListByDeviceId(ctx *gin.Context) {
	var tmp struct {
		DeviceId int `json:"deviceId"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var deviceCmdParamList []models.EmDeviceModelCmdParam
	deviceCmdParamList, _ = c.emRepo.GetEmDeviceModelCmdParamListByDeviceId(tmp.DeviceId)
	if deviceCmdParamList == nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "无数据",
		})
		return
	}

	//收集所有codes查询最新测点信息
	var codes []string
	mapCodes := make(map[string]string) //按name进行label分组
	for i := 0; i < len(deviceCmdParamList); i++ {
		codes = append(codes, deviceCmdParamList[i].Name)
		mapCodes[deviceCmdParamList[i].Name] = deviceCmdParamList[i].Label
	}

	ycList, err := c.hisRepo.GetLastYcListByCode(strconv.Itoa(tmp.DeviceId), strings.Join(codes, ","))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, model := range ycList {
		model.Name = mapCodes[strconv.Itoa(model.Code)]
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    ycList,
	})
}
