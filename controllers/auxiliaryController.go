package controllers

import (
	"gateway/httpServer/model"
	"gateway/models"
	"gateway/models/ReturnModel"
	"gateway/models/query"
	repositories "gateway/repositories"
	"gateway/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// AuxiliaryController 定义辅控管理的控制器
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
	router.POST("/api/v2/auxiliary/getDeviceControlPointList", ctrl.GetDeviceControlPointList)

}

// GetDeviceListByDeviceType /获取设备类型下的所有设备数据
func (ctrl *AuxiliaryController) GetDeviceListByDeviceType(ctx *gin.Context) {
	type tmpQuery struct {
		DeviceType string `json:"deviceType"`
	}
	var queryParam tmpQuery
	if err := ctx.Bind(&queryParam); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	deviceList, err := ctrl.repo.GetAuxiliaryDevice(queryParam.DeviceType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取信息成功！",
		Data:    deviceList,
	})
	return
}

// GetAuxiliaryDeviceType 获取所有设备类型
func (ctrl *AuxiliaryController) GetAuxiliaryDeviceType(ctx *gin.Context) {
	deviceTypeList, err := ctrl.repo.GetAuxiliaryDeviceType()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取信息成功！",
		Data:    deviceTypeList,
	})
	return
}

// GetLastYcByDeviceIdAndCodes 获取最新遥测信息GetLastYcListByCode
// 遥测控制器已经写有
func (ctrl *AuxiliaryController) GetLastYcByDeviceIdAndCodes(ctx *gin.Context) {
	type ycQueryData struct {
		DeviceIds []int `form:"deviceIds"`
		Codes     []int `form:"codes"`
	}
	var ycQuery ycQueryData
	//将传过来的请求体解析到ycQuery中
	if err := ctx.Bind(&ycQuery); err != nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "error" + err.Error(),
			Data:    "",
		})
		return
	}
	//查询历史数据
	//将ycQuery.DeviceIds转换成字符串
	deviceIds := utils.IntArrayToString(ycQuery.DeviceIds, ",")
	codes := utils.IntArrayToString(ycQuery.Codes, ",")
	ycLog, err := ctrl.hisRepo.GetLastYcListByCode(deviceIds, codes)
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

type ReturnMap struct {
	XAxisList []string          `json:"xAxisList"`
	DataMap   map[int][]float64 `json:"dataMap"`
}

//根据选择的codes返回对应时间的历史数据

func (ctrl *AuxiliaryController) GetHistoryYcByDeviceIdCodes(ctx *gin.Context) {
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
	returnMap, err := ctrl.hisRepo.GetCharData(ycQuery)
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

type DeviceModel struct {
	Property struct {
		AccessMode int `json:"accessMode"`
	} `json:"property"`
}

// GetEmDeviceModelCmdParamListByDeviceId 根据设备ID获取遥测列表
func (ctrl *AuxiliaryController) GetEmDeviceModelCmdParamListByDeviceId(ctx *gin.Context) {
	var tmp struct {
		DeviceId int `json:"deviceId"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	array := []string{"yc"}
	likeQuery := "param.data  like '%\"accessMode\":0%'" //等于0说明是不可控测点
	var deviceCmdParamList []models.EmDeviceModelCmdParam
	deviceCmdParamList, _ = ctrl.emRepo.GetCodesListByDeviceIdAndYxYc(tmp.DeviceId, array, likeQuery)
	if deviceCmdParamList == nil {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "无数据",
		})
		return
	}

	//收集所有codes查询最新测点信息
	var codes []string
	for i := 0; i < len(deviceCmdParamList); i++ {
		codes = append(codes, deviceCmdParamList[i].Name)
	}

	//拿到所有不可控测点的最新数据
	ycList, err := ctrl.hisRepo.GetLastYcListByCode(strconv.Itoa(tmp.DeviceId), strings.Join(codes, ","))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ycMap := make(map[string]*models.YcData) //收集遥测信息
	for i := 0; i < len(ycList); i++ {
		ycMap[strconv.Itoa(ycList[i].Code)] = ycList[i]
	}

	var resYcData []ReturnModel.AuxYcData
	for _, deviceParam := range deviceCmdParamList { //遍历结果 重新拼接，需要返回name，单位字段
		var ycData ReturnModel.AuxYcData
		var tmpMap = ycMap[deviceParam.Name]
		ycData.Name = deviceParam.Label                 //从模型取
		ycData.Unit = deviceParam.Unit                  //从模型取
		ycData.ParamId = deviceParam.Id                 //从模型取
		ycData.Code, _ = strconv.Atoi(deviceParam.Name) //从模型取
		if tmpMap != nil {
			ycData.DeviceId = tmpMap.DeviceId //从测点取
			ycData.Ts = tmpMap.Ts             //从测点取
			ycData.Value = tmpMap.Value       //从测点取
		}
		resYcData = append(resYcData, ycData)
	}

	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "成功",
		Data:    resYcData,
	})
}

// GetDeviceControlPointList 获取实时控制遥控遥调列表
func (ctrl *AuxiliaryController) GetDeviceControlPointList(ctx *gin.Context) {
	type Res struct {
		YcData []ReturnModel.AuxYcData `json:"ycData"`
		YxData []ReturnModel.AuxYxData `json:"yxData"`
	}
	var tmp struct {
		DeviceId int `json:"deviceId"`
	}
	if err := ctx.ShouldBindBodyWith(&tmp, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	array := []string{"yc", "yx"}
	likeQuery := "param.data not like '%\"accessMode\":0%'" //不在0里面为可控测点
	YkYtList, err := ctrl.emRepo.GetCodesListByDeviceIdAndYxYc(tmp.DeviceId, array, likeQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var result Res
	if len(YkYtList) == 0 {
		ctx.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "无数据",
			Data:    result,
		})
		return
	}
	//收集所有codes查询最新测点信息
	var YxCodes []string
	var YcCodes []string
	nameMap := make(map[string]models.EmDeviceModelCmdParam) //存如map
	var ycData []ReturnModel.AuxYcData                       //存储可控遥测返回值
	var yxData []ReturnModel.AuxYxData                       //存储可控遥信返回值
	for _, item := range YkYtList {
		nameMap[item.Name] = item
		if item.IotDataType == "yc" {
			var tmpYcData ReturnModel.AuxYcData
			tmpYcData.DeviceId = tmp.DeviceId //设备id
			tmpYcData.Name = item.Label
			tmpYcData.Unit = item.Unit
			tmpYcData.ParamId = item.Id                 //参数id用于控制
			tmpYcData.Code, _ = strconv.Atoi(item.Name) //转成int
			ycData = append(ycData, tmpYcData)          //收集遥测数据
			YcCodes = append(YcCodes, item.Name)        //收集遥测code
		} else if item.IotDataType == "yx" {
			var tmpYxData ReturnModel.AuxYxData
			tmpYxData.DeviceId = tmp.DeviceId //设备id
			tmpYxData.Name = item.Label
			tmpYxData.Unit = item.Unit
			tmpYxData.ParamId = item.Id
			tmpYxData.Code, _ = strconv.Atoi(item.Name) //转成int
			yxData = append(yxData, tmpYxData)
			YxCodes = append(YxCodes, item.Name)
		}
	}
	//收集成map
	ycMap := make(map[int]*models.YcData)
	if len(YcCodes) > 0 {
		ycList, err := ctrl.hisRepo.GetLastYcListByCode(strconv.Itoa(tmp.DeviceId), strings.Join(YcCodes, ",")) //拿到所有可控测点的最新数据
		if err == nil && len(ycList) > 0 {
			for i := 0; i < len(ycList); i++ {
				ycMap[ycList[i].Code] = ycList[i]
			}
			result.YcData = ycData
		}
	}
	//收集成map
	yxMap := make(map[int]*models.YxData)
	if len(YxCodes) > 0 {
		yxList, err := ctrl.hisRepo.GetLastYxListByCode(tmp.DeviceId, strings.Join(YxCodes, ","))
		if err == nil && len(yxList) > 0 {
			for i := 0; i < len(yxList); i++ {
				yxMap[yxList[i].Code] = yxList[i] //存入map
			}

		}

	}
	//遍历ycData，yxData，如果有对应的code，就赋值
	for i := range ycData {
		if ycMap[ycData[i].Code] != nil {
			ycData[i].Value = ycMap[ycData[i].Code].Value
			ycData[i].Ts = ycMap[ycData[i].Code].Ts
		}
	}
	for i := range yxData {
		if yxMap[yxData[i].Code] != nil {
			yxData[i].Value = yxMap[yxData[i].Code].Value
			yxData[i].Ts = yxMap[yxData[i].Code].Ts
		}
	}
	result.YcData = ycData
	result.YxData = yxData
	ctx.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取数据成功！",
		Data:    result,
	})
	return
}
